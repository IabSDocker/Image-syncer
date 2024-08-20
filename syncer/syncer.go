package syncer

import (
    "fmt"
    "os/exec"
    "sync"
    "image-syncer/utils"
    "github.com/sirupsen/logrus"
)

type Syncer struct {
    SourceRepo   string
    TargetRepo   string
    Namespace    string
    Images       []string
    Logger       *logrus.Logger
    successCount int
    failureCount int
    failures     []string
}

func NewSyncer(config *utils.Config, logger *logrus.Logger) *Syncer {
    return &Syncer{
        SourceRepo: config.SourceRepo,
        TargetRepo: config.TargetRepo,
        Namespace:  config.Namespace,
        Images:     config.Images,
        Logger:     logger,
        failures:   []string{},
    }
}

func (s *Syncer) syncImage(image string, wg *sync.WaitGroup, ch chan<- string) {
    defer wg.Done()

    // 构建源镜像路径
    sourceImage := fmt.Sprintf("%s/%s", s.SourceRepo, image)

    // 构建目标镜像路径
    var targetImage string
    if s.Namespace == "" {
        // 不启用 namespace 字段，目标镜像路径和源镜像路径一致
        targetImage = fmt.Sprintf("%s/%s", s.TargetRepo, image)
    } else {
        // 启用 namespace 字段，将 namespace 加入到目标镜像路径
        // 提取镜像名称（去掉仓库路径）
        imageName := image
        if idx := findFirstSlash(image); idx != -1 {
            imageName = image[idx+1:]
        }
        targetImage = fmt.Sprintf("%s/%s/%s", s.TargetRepo, s.Namespace, imageName)
    }

    s.Logger.Infof("Source image: %s", sourceImage)
    s.Logger.Infof("Target image: %s", targetImage)

    // 拉取镜像
    s.Logger.Infof("Pulling image: %s", sourceImage)
    cmd := exec.Command("docker", "pull", sourceImage)
    if output, err := cmd.CombinedOutput(); err != nil {
        s.Logger.Errorf("Failed to pull image %s: %v, Output: %s", sourceImage, err, string(output))
        s.failureCount++
        s.failures = append(s.failures, fmt.Sprintf("Failed to pull image %s: %v, Output: %s", sourceImage, err, string(output)))
        return
    }

    // 标记镜像
    s.Logger.Infof("Tagging image: %s as %s", sourceImage, targetImage)
    cmd = exec.Command("docker", "tag", sourceImage, targetImage)
    if output, err := cmd.CombinedOutput(); err != nil {
        s.Logger.Errorf("Failed to tag image %s as %s: %v, Output: %s", sourceImage, targetImage, err, string(output))
        s.failureCount++
        s.failures = append(s.failures, fmt.Sprintf("Failed to tag image %s as %s: %v, Output: %s", sourceImage, targetImage, err, string(output)))
        return
    }

    // 推送镜像
    s.Logger.Infof("Pushing image: %s", targetImage)
    cmd = exec.Command("docker", "push", targetImage)
    if output, err := cmd.CombinedOutput(); err != nil {
        s.Logger.Errorf("Failed to push image %s: %v, Output: %s", targetImage, err, string(output))
        s.failureCount++
        s.failures = append(s.failures, fmt.Sprintf("Failed to push image %s: %v, Output: %s", targetImage, err, string(output)))
        return
    }

    s.successCount++
    ch <- fmt.Sprintf("Successfully synced %s", image)
    
    // 删除源镜像
    exec.Command("docker", "rmi", sourceImage).Run()

    // 删除目标镜像
    exec.Command("docker", "rmi", targetImage).Run()
}

func (s *Syncer) StartSync() {
    var wg sync.WaitGroup
    ch := make(chan string, len(s.Images))

    // 使用固定数量的 goroutines 并发处理镜像同步
    maxConcurrency := 5
    sem := make(chan struct{}, maxConcurrency)

    for _, image := range s.Images {
        wg.Add(1)
        sem <- struct{}{}
        go func(img string) {
            defer func() { <-sem }()
            s.syncImage(img, &wg, ch)
        }(image)
    }

    go func() {
        wg.Wait()
        close(ch)
    }()

    for msg := range ch {
        s.Logger.Info(msg)
    }

    // 输出失败信息
    if s.failureCount > 0 {
        s.Logger.Errorf("Failures: \n%s", formatFailures(s.failures))
    }

    // 输出同步统计信息
    s.Logger.Infof("Syncing completed. Successfully synced %d images, Failed to sync %d images.", s.successCount, s.failureCount)
}

// 格式化失败信息的函数
func formatFailures(failures []string) string {
    result := ""
    for _, failure := range failures {
        result += failure + "\n"
    }
    return result
}

// 查找第一个 '/' 的索引
func findFirstSlash(image string) int {
    for i, ch := range image {
        if ch == '/' {
            return i
        }
    }
    return -1
}
