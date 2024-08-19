package syncer

import (
    "fmt"
    "os/exec"
    "sync"
    "image-syncer/utils"
    "github.com/sirupsen/logrus"
)

type Syncer struct {
    SourceRepo  string
    TargetRepo  string
    Images      []string
    Logger      *logrus.Logger
    successCount int
    failureCount int
    failures     []string
}

func NewSyncer(config *utils.Config, logger *logrus.Logger) *Syncer {
    return &Syncer{
        SourceRepo: config.SourceRepo,
        TargetRepo: config.TargetRepo,
        Images:     config.Images,
        Logger:     logger,
        failures:   []string{},
    }
}

func (s *Syncer) syncImage(image string, wg *sync.WaitGroup, ch chan<- string) {
    defer wg.Done()

    sourceImage := fmt.Sprintf("%s/%s", s.SourceRepo, image)
    targetImage := fmt.Sprintf("%s/%s", s.TargetRepo, image)

    // Pull the image from the source repository
    s.Logger.Infof("Pulling image: %s", sourceImage)
    cmd := exec.Command("docker", "pull", sourceImage)
    if output, err := cmd.CombinedOutput(); err != nil {
        s.Logger.Errorf("Failed to pull image %s: %v, Output: %s", sourceImage, err, string(output))
        s.failureCount++
        s.failures = append(s.failures, fmt.Sprintf("Failed to pull image %s: %v, Output: %s", sourceImage, err, string(output)))
        return
    }

    // Tag the image for the target repository
    s.Logger.Infof("Tagging image: %s as %s", sourceImage, targetImage)
    cmd = exec.Command("docker", "tag", sourceImage, targetImage)
    if output, err := cmd.CombinedOutput(); err != nil {
        s.Logger.Errorf("Failed to tag image %s as %s: %v, Output: %s", sourceImage, targetImage, err, string(output))
        s.failureCount++
        s.failures = append(s.failures, fmt.Sprintf("Failed to tag image %s as %s: %v, Output: %s", sourceImage, targetImage, err, string(output)))
        return
    }

    // Push the image to the target repository
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