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