## 部署流程

### 设计

docker-compose和镜像之间的关系要有明确分工

- ```build/package``` 目录只负责构建生成镜像
- ```deploy/docker-compose/docker-compose.yml``` 只负责编排镜像, 禁止直接build镜像

通过上述分工，可以方便的拆分docker-compose和Dockerfile之间的关系，便于后续调整

### 线下环境
通过docker-compose进行编排运行

### 设计
通过k8s进行编排运行
