## 安装
- 删除无效的 docker image \
  docker images|grep none|awk '{print $3}'|xargs docker rmi


- 进入项目目录 \
  $ cd dapr 


- 删除原有安装 \
  $ dapr uninstall -k


- 创建命名空间 \
  $ kubectl create namesapce dapr-system


- 设置环境变量 \
  $ export DAPR_REGISTRY=192.168.64.12 \
  $ export DAPR_TAG=dev \
  $ export  TARGET_OS=linux \
  $ export  TARGET_ARCH=arm64 


- 编译二进制文件 \
  linux-arm64 \
  $ make build-linux GOOS=linux GOARCH=arm64 \
  
  windows-amd64 \
  $ make build GOOS=windows GOARCH=amd64 \

  mac-arm64 \
  $ make build GOOS=darwin GOARCH=arm64 \

- 修改docker/docker.rm文件， 增加参数 --load
  
    
    docker-build: check-docker-env check-arch
        $(info Building $(DOCKER_IMAGE_TAG) docker image ...)
    ifeq ($(TARGET_ARCH),amd64)
        $(DOCKER) build PKG_FILES=* -f $(DOCKERFILE_DIR)/$(DOCKERFILE) $(BIN_PATH) -t $(DOCKER_IMAGE_TAG)-$(TARGET_OS)-$(TARGET_ARCH)
        $(DOCKER) build PKG_FILES=daprd -f $(DOCKERFILE_DIR)/$(DOCKERFILE) $(BIN_PATH) -t $(DAPR_RUNTIME_DOCKER_IMAGE_TAG)-$(TARGET_OS)-$(TARGET_ARCH)
        $(DOCKER) build PKG_FILES=placement -f $(DOCKERFILE_DIR)/$(DOCKERFILE) $(BIN_PATH) -t $(DAPR_PLACEMENT_DOCKER_IMAGE_TAG)-$(TARGET_OS)-$(TARGET_ARCH)
        $(DOCKER) build PKG_FILES=sentry -f $(DOCKERFILE_DIR)/$(DOCKERFILE) $(BIN_PATH) -t $(DAPR_SENTRY_DOCKER_IMAGE_TAG)-$(TARGET_OS)-$(TARGET_ARCH)
    else
        -$(DOCKER) buildx create --use --name daprbuild
        -$(DOCKER) run --rm --privileged multiarch/qemu-user-static --reset -p yes
        $(DOCKER) buildx build --load  --build-arg PKG_FILES=*         --platform $(DOCKER_IMAGE_PLATFORM) -f $(DOCKERFILE_DIR)/$(DOCKERFILE) $(BIN_PATH) -t $(DOCKER_IMAGE_TAG)-$(TARGET_OS)-$(TARGET_ARCH)
        $(DOCKER) buildx build --load  --build-arg PKG_FILES=daprd     --platform $(DOCKER_IMAGE_PLATFORM) -f $(DOCKERFILE_DIR)/$(DOCKERFILE) $(BIN_PATH) -t $(DAPR_RUNTIME_DOCKER_IMAGE_TAG)-$(TARGET_OS)-$(TARGET_ARCH)
        $(DOCKER) buildx build --load  --build-arg PKG_FILES=placement --platform $(DOCKER_IMAGE_PLATFORM) -f $(DOCKERFILE_DIR)/$(DOCKERFILE) $(BIN_PATH) -t $(DAPR_PLACEMENT_DOCKER_IMAGE_TAG)-$(TARGET_OS)-$(TARGET_ARCH)
        $(DOCKER) buildx build --load  --build-arg PKG_FILES=sentry    --platform $(DOCKER_IMAGE_PLATFORM) -f $(DOCKERFILE_DIR)/$(DOCKERFILE) $(BIN_PATH) -t $(DAPR_SENTRY_DOCKER_IMAGE_TAG)-$(TARGET_OS)-$(TARGET_ARCH)
    endif

- 生成 docker image \
  $ sudo make docker-build TARGET_OS=linux TARGET_ARCH=arm64 DAPR_REGISTRY=192.168.64.12 DAPR_TAG=dapr
	

- 查看是否生成 \
  $ docker images 
	

    REPOSITORY                    TAG                     IMAGE ID          CREATED             SIZE
    192.168.64.12/sentry          dapr-linux-arm64        8f909cf60e46      13 minutes ago      37MB
    192.168.64.12/placement       dapr-linux-arm64        04eb59389523      13 minutes ago      16.3MB
    192.168.64.12/daprd           dapr-linux-arm64        936608e34a01      13 minutes ago      105MB
    192.168.64.12/dapr            dapr-linux-arm64        207d8890f756      13 minutes ago      286MB
	
	
- 推送image文件到私库中 \
    $ docker push 192.168.64.12/sentry:dapr-linux-arm64 \
    $ docker push 192.168.64.12/placement:dapr-linux-arm64 \
    $ docker push 192.168.64.12/daprd:dapr-linux-arm64 \
    $ docker push 192.168.64.12/dapr:dapr-linux-arm64


- 安装方法 \
  $ make docker-deploy-k8s TARGET_OS=linux TARGET_ARCH=arm64 DAPR_REGISTRY=192.168.64.12  DAPR_TAG=dapr


- 查询pod \
  $ kubectl get pod -n dapr-system 


      NAMESPACE     NAME                                       READY     STATUS     RESTARTS     AGE 
      dapr-system   dapr-dashboard-c8dd8d969-xgdtx             1/1       Running    0            18m 
      dapr-system   dapr-operator-7c7d5f887b-b5dsn             1/1       Running    0            18m 
      dapr-system   dapr-placement-server-0                    1/1       Running    0            18m 
      dapr-system   dapr-sentry-6cc6687cc8-sqrcr               1/1       Running    0            18m 
      dapr-system   dapr-sidecar-injector-698756d5d7-46fpd     1/1       Running    0            18m 
      
