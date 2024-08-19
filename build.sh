#!/bin/bash

# 定义目标操作系统和架构的数组
declare -A targets
targets=(
    ["linux-amd64"]="linux amd64"
    ["linux-arm"]="linux arm"
    ["linux-arm64"]="linux arm64"
    ["windows-amd64"]="windows amd64"
    ["darwin-amd64"]="darwin amd64"
    ["mips"]="linux mips"
    ["mips64"]="linux mips64"
)

# 遍历每一个目标平台并进行编译
for target in "${!targets[@]}"; do
    os_arch=${targets[$target]}
    IFS=' ' read -r os arch <<< "$os_arch"
    
    echo "Building for OS: $os, ARCH: $arch"
    
    # 设定 GOARCH 和 GOOS
    GOOS=$os GOARCH=$arch go build -o "syncer-${target}" main.go

    # 打包为 tar.gz 文件
    tar -czvf "syncer-${target}.tar.gz" "syncer-${target}"
done

echo "所有平台的编译和打包完成。"