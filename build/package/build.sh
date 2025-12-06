#!/bin/bash

BASE_PATH=$(cd `dirname $0`; cd ../../; pwd)
cd $BASE_PATH

# Usage:  docker build [OPTIONS] PATH | URL | -
# 最后一个传递的PATH是DockerFile构建的上下文context, 是必须传递的参数
# 通常使用`.`即可, 表示当前项目根目录, 这样无论Dockerfile文件在哪个目录下，使用ADD等指令操作文件的时候, 都可以像在根目录下操作文件一样。
docker build -t muse-admin:latest -f $BASE_PATH/build/package/Dockerfile .
