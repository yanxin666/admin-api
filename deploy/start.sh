#!/bin/bash

# 颜色集合
setup_color() {
    FMT_RAINBOW=$(printf '\033[38;2;245;0;172m') # 深粉
    FMT_RAINBOW_LIGHT=$(printf '\033[38;5;163m') # 浅粉
    FMT_RED=$(printf '\033[31m')                 # 红色
    FMT_GREEN=$(printf '\033[32m')               # 绿色
    FMT_YELLOW=$(printf '\033[33m')              # 黄色
    FMT_BLUE=$(printf '\033[34m')                # 蓝色
    FMT_BOLD=$(printf '\033[1m')                 # 白色加粗
    FMT_RESET=$(printf '\033[0m')                # 重置
}

# 获取根目录
base_path() {
    BASE_PATH=$(
        cd $(dirname $0)
        cd ..
        pwd
    )
    cd $BASE_PATH
}


# 获取模板内容
get_temp() {
    cat << DESC

🤫  Tips：若目标分支并非当前分支，请在项目目录下执行 ${FMT_BOLD} sh ./deploy/start.sh you_env ${FMT_RESET}

环境部署说明 👉 ${FMT_YELLOW}请仔细阅读此说明！！${FMT_RESET}：
    dev：    开发环境，端口为 8886 -> [docker部署]
    test：   测试环境，端口为 8886 -> [docker部署]
    pre：    预发环境，coding上执行构建计划，之后会自动部署 -> [k8s部署]
    pro：    生产环境，coding上执行构建计划，之后需要联系相关权限者手动部署 -> [k8s部署]

✍️  如有疑问，请移步下方文档中详细操作
${FMT_BOLD}https://douyuxingchen.feishu.cn/docx/UYsUdXf5houE2Lx0n42cqZKznqc${FMT_RESET}

DESC
}

# 设置分支
set_branch() {
    if [ "$1" != "" ]; then
        branch="$1"
    fi
}

# 获取git内容
git_info() {
    env="dev"
    git_current_branch=$(git rev-parse --abbrev-ref HEAD)

    if [ "$branch" != "" ]; then
        git_current_branch=${branch}
    fi

    if ! git show-ref --verify --quiet "refs/heads/$git_current_branch"; then
        echo -e "${FMT_RAINBOW_LIGHT} $git_current_branch 分支切换不存在，请检查后重试  ${FMT_RESET} \n"
        exit 1
    fi

    printf "${FMT_BOLD}正在拉取当前分支最新记录：${FMT_RESET} \n"

    # 切回到对应分支拉取最新代码
    git checkout $git_current_branch
    git pull origin $git_current_branch
    printf "\n"

    now_time=$(date "+%Y-%m-%d %H:%M:%S")
    git_current_commit=$(git rev-parse HEAD)
    git_current_commit_msg=$(git log -1 --pretty=%B)
    git_current_log_history=$(git log -1 --pretty=format:"%h%x09%an%x09%ar%x09%s" | awk -F $'\t' '{printf "%-10s%-15s%-15s%s\n", $1, $2, $3, $4}')
    git_log_history=$(git log -n 4 --skip=1 --pretty=format:"%h%x09%an%x09%ar%x09%s" | awk -F $'\t' '{printf "%-10s%-15s%-15s%s\n", $1, $2, $3, $4}')

    # 获取当前路径
    path=$(pwd)
    # 使用basename命令获取路径的最后一个部分，即项目名称
    project_name=$(basename $path)
}

# 获取git核心内容
get_content() {
    # 设置字符串的最大宽度
    max_width=30

    printf "${FMT_BOLD}发布详情：${FMT_RESET} \n"
    # 打印表格
    printf "+-------------+--------------------------------------------------------------+\n"
    printf "|   部署项目  |   %-*s\n" $max_width "${FMT_BLUE}$project_name${FMT_RESET} "
    printf "+-------------+--------------------------------------------------------------+\n"
    printf "|   部署环境  |   %-*s\n" $max_width "${FMT_RED}$env${FMT_RESET} "
    printf "+-------------+--------------------------------------------------------------+\n"
    printf "|   分支名称  |   %-*s\n" $max_width "${FMT_BOLD}$git_current_branch${FMT_RESET} "
    printf "+-------------+--------------------------------------------------------------+\n"
    printf "|   部署时间  |   %-*s\n" $max_width "${FMT_GREEN}$now_time${FMT_RESET} "
    printf "+-------------+--------------------------------------------------------------+\n"
    printf "|   提交ID    |   %-*s\n" $max_width "${git_current_commit:0:$max_width} "
    printf "+-------------+--------------------------------------------------------------+\n"
    printf "|   最新注释  |   %-*s\n" $max_width "${FMT_YELLOW}$git_current_commit_msg${FMT_RESET} "
    printf "+-------------+--------------------------------------------------------------+\n"
}

# 获取git页尾记录
get_content_foot() {
    echo "👏👏👏"
    printf "${FMT_YELLOW}$git_current_log_history${FMT_RESET}\n"
    printf "${FMT_GREEN}$git_log_history${FMT_RESET}\n"
    printf "+-------------+--------------------------------------------------------------+\n"
}

# 检查分支
check_branch() {
    if [ "$1" != "" ]; then
        env="$1"
        if [ "$env" != "dev" ] && [ "$env" != "test" ] && [ "$env" != "pro" ] && [ "$env" != "pre" ]; then
            echo -e "\n 参数输入错误，请在 「dev、test、pro、pre」中选择一个进行发布。当前参数: ${env}"
            exit 1
        fi
    fi

    # 服务启动
    if [ "$env" = "pro" ] || [ "$env" = "pre" ]; then
        runonline
    fi

}

# 开始构建
go_build() {
    read -p "上述信息确认无误后，进行发布 👉 ${FMT_BOLD}$git_current_branch${FMT_RESET}  (Y/N): " choice

    # 根据用户输入执行相应操作
    if [[ "$choice" = "Y" || "$choice" = "y" ]]; then

        # 追加日志
        printf "${FMT_RAINBOW}=======> 发布了新版本：${FMT_RESET}\n" >> "/tmp/logs/deployment.log"
        get_content >> "/tmp/logs/deployment.log"
        echo -e "\n${FMT_BLUE}此发布记录已追加到日志文件中${FMT_RESET} ==============> ${FMT_BOLD}/tmp/logs/deployment.log${FMT_RESET}\n"

        # 开始运行
        runoffine

    elif
        [[ "$choice" = "N" || "$choice" = "n" ]]
    then
        echo "您退出了发布操作。"
    else
        echo "无效的输入。请键入 Y 或 N。"
    fi
}

# 线上环境使用k8s部署
# 线上环境包括pre/online
runonline() {
    echo "线上环境使用k8s部署, 目前暂不支持线上环境部署"
    exit 1
}

# 线下环境使用docker-compose部署
# 线下环境包括dev/test
runoffine() {
    suffix="-dev"

    # 项目容器
    echo "step1 停止容器 ===> stop container swagger-muse-admin${suffix} && muse-admin${suffix}"
    docker stop muse-admin${suffix}
    docker stop swagger-muse-admin${suffix}

    echo "step2 删除容器 ===> rm container swagger-muse-ability${suffix} && muse-ability${suffix}"
    docker rm muse-admin${suffix}
    docker rm swagger-muse-admin${suffix}

    echo "step3 删除镜像 ===> rm image muse-admin${suffix}:latest"
    docker rmi muse-admin:latest

    echo "step4 构建镜像===> bash build/package/build.sh"
    bash build/package/build.sh

    echo "step5 启动服务===> build & run muse-admin${suffix}"
    docker-compose -p muse-admin${suffix} -f "${BASE_PATH}/deploy/docker-compose/docker-compose.${env}.yml" up -d

    # 检查容器是否成功启动
    if [ "$(docker ps -q -f name=muse-admin${suffix})" ]; then
        echo -e "\n 🎉🎉🎉 容器部署成功 🥳🥳🥳"
    else
        echo -e "\n 👻👻👻 容器部署失败 🤕🤕🤕"
    fi
}

main() {
    setup_color
    base_path
    get_temp
    set_branch "$1"
    git_info
    get_content
    get_content_foot
#    check_branch "$1"
    go_build
}

main "$@"
