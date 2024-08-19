# image-syncer
`image-syncer` 是阿里云开发的容器镜像同步工具，可用来进行多对多的镜像仓库同步，支持目前绝大多数主流的 docker 镜像仓库服务
阿里云官方[image-syncer](https://github.com/AliyunContainerService/image-syncer) 地址

## Image-syncer
该项目是受image-syncer启发同样用go语言写的`轻量`(go还使用的不是太熟练😅)的镜像同步工具。

### 手动编译

#### 当前系统架构编译
```bash
git clone https://gitlab.ayou.ink/IabSDocker/Image-syncer.git
cd Image-syncer
go build -o image-syncer
```

#### 所有常用系统架构编译
```bash
git clone https://gitlab.ayou.ink/IabSDocker/Image-syncer.git
cd Image-syncer
sh build.sh
```

### 命令用例
##### namespace与源仓库一致
```bash
[root@luckly ~]# cat config.yaml
# 原仓库
source_repo: "docker.io"
# 目标仓库
target_repo: "localhost:5000"
# 要同步的镜像列表
images:
  - "iabsdocker/registry:latest"  
  - "nginx:latest"
[root@luckly ~]# ./image-syncer --config config.yaml
[root@luckly ~]# docker image ls
REPOSITORY                                  TAG               IMAGE ID       CREATED         SIZE
localhost:5000/nginx                        latest            5ef79149e0ec   4 days ago      188MB
nginx                                       latest            5ef79149e0ec   4 days ago      188MB
iabsdocker/registry                         latest            f7b1a4c78949   5 weeks ago     45.3MB
localhost:5000/registry                     latest            f7b1a4c78949   5 weeks ago     45.3MB
```

##### 自定义namespace字段
```bash
[root@luckly ~]# cat config.yaml
# 原仓库
source_repo: "docker.io"
# 目标仓库
target_repo: "localhost:5000"
namespace: "mytest"
# 要同步的镜像列表
images:
  - "iabsdocker/registry:latest"  
  - "nginx:latest"
[root@luckly ~]# ./image-syncer --config config.yaml
[root@luckly ~]# docker image ls
REPOSITORY                                  TAG               IMAGE ID       CREATED         SIZE
localhost:5000/mytest/nginx                 latest            5ef79149e0ec   4 days ago      188MB
nginx                                       latest            5ef79149e0ec   4 days ago      188MB
iabsdocker/registry                         latest            f7b1a4c78949   5 weeks ago     45.3MB
localhost:5000/mytest/iabsdocker/registry   latest            f7b1a4c78949   5 weeks ago     45.3MB
```

### 日志相关
#### 镜像同步完成后会统计结果，所有同步失败的镜像会在Failures字段中给出报错信息！！！
```bash
[root@luckly Image-syncer]# cat config/config.yaml 
source_repo: "docker.io"
target_repo: "localhost:5000"
images:
  - "nginx:1.29.0"
  - "nginx:1.30.0"
  - "nginx:latest"
  - "alpine:3.12"
[root@luckly Image-syncer]# ./syncer-linux-amd64 --config config/config.yaml 
INFO[2024-08-19T16:33:25+08:00] Pulling image: docker.io/nginx:latest        
INFO[2024-08-19T16:33:25+08:00] Pulling image: docker.io/nginx:1.30.0        
INFO[2024-08-19T16:33:25+08:00] Pulling image: docker.io/nginx:1.29.0        
INFO[2024-08-19T16:33:25+08:00] Pulling image: docker.io/alpine:3.12         
INFO[2024-08-19T16:33:26+08:00] Tagging image: docker.io/alpine:3.12 as localhost:5000/alpine:3.12 
INFO[2024-08-19T16:33:26+08:00] Pushing image: localhost:5000/alpine:3.12    
INFO[2024-08-19T16:33:26+08:00] Tagging image: docker.io/nginx:latest as localhost:5000/nginx:latest 
INFO[2024-08-19T16:33:26+08:00] Pushing image: localhost:5000/nginx:latest   
INFO[2024-08-19T16:33:26+08:00] Successfully synced alpine:3.12              
ERRO[2024-08-19T16:33:26+08:00] Failed to pull image docker.io/nginx:1.30.0: exit status 1, Output: Error response from daemon: manifest for nginx:1.30.0 not found: manifest unknown: manifest unknown 
ERRO[2024-08-19T16:33:26+08:00] Failed to pull image docker.io/nginx:1.29.0: exit status 1, Output: Error response from daemon: manifest for nginx:1.29.0 not found: manifest unknown: manifest unknown 
INFO[2024-08-19T16:33:26+08:00] Successfully synced nginx:latest             
ERRO[2024-08-19T16:33:26+08:00] Failures: 
Failed to pull image docker.io/nginx:1.30.0: exit status 1, Output: Error response from daemon: manifest for nginx:1.30.0 not found: manifest unknown: manifest unknown

Failed to pull image docker.io/nginx:1.29.0: exit status 1, Output: Error response from daemon: manifest for nginx:1.29.0 not found: manifest unknown: manifest unknown
 
INFO[2024-08-19T16:33:26+08:00] Syncing completed. Successfully synced 2 images, Failed to sync 2 images. 
```