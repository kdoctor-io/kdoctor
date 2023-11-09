# Kind Quick Start

[**English**](./get-started-kind.md) | **简体中文**

Kind 是一个使用 Docker 容器节点运行本地 Kubernetes 集群的工具。kdoctor 提供了安装 Kind 集群的脚本，能快速搭建一套配备 nginx、ingress、loadbalancer 的 Kind 集群，您可以使用它来进行 kdoctor 的测试与体验。

## 先决条件

* 已安装 [Go](https://go.dev/)

## 在 Kind 集群上部署 kdoctor

1. 克隆 kdoctor 代码仓库到本地主机上，并进入 kdoctor 工程的根目录。
  
    ```bash
    https://github.com/kdoctor-io/kdoctor.git && cd kdoctor
    ```

2. 执行 `make checkBin`，检查本地主机上的开发工具是否满足部署 Kind 集群与 kdoctor 的条件，如果缺少组件会为您自动安装。

3. 通过以下方式获取 kdoctor 的最新镜像。

    ```bash
    ~# KDOCTOR_LATEST_IMAGE_TAG=$(curl -s https://api.github.com/repos/kdoctor-io/kdoctor/releases | jq -r '.[].tag_name | select(("^v1.[0-9]*.[0-9]*$"))' | head -n 1)
    ```

4. 执行以下命令，创建 Kind 集群，并为您安装 metallb、contour、nginx、kdoctor。

    ```bash
    ~# make e2e_init -e PROJECT_IMAGE_VERSION=KDOCTOR_LATEST_IMAGE_TAG
    ```

!!! note

    如果您是国内用户，您可以使用如下命令，避免拉取镜像失败。

    ```bash
    ~# make e2e_init -e E2E_SPIDERPOOL_TAG=$SPIDERPOOL_LATEST_IMAGE_TAG -e E2E_CHINA_IMAGE_REGISTRY=true
    ```

## 验证安装

在 kdoctor 工程的根目录下执行如下命令，为 kubectl 配置 Kind 集群的 KUBECONFIG。

   ```bash
   ~# export KUBECONFIG=$(pwd)/test/runtime/kubeconfig_kdoctor.config
   ```

您可以看到如下的内容输出：

   ```bash
   ~# kubectl get nodes 
   NAME                    STATUS   ROLES           AGE     VERSION
   kdoctor-control-plane   Ready    control-plane   7m3s    v1.27.1
   kdoctor-worker          Ready    <none>          6m42s   v1.27.1
   
   ~# kubectll get po -n kdoctor
   NAME                                  READY   STATUS    RESTARTS   AGE
   kdoctor-controller-686b75d6d7-ktctx   1/1     Running   0          2m33s
   ```

接下来您可以根据您的需要进行任务的布置 [AppHttpHealthy](./apphttphealthy-zh_CN.md)、[NetReach](./netreach-zh_CN.md)、[NetDns](./netdns-zh_CN.md)

## 卸载

* 卸载 Kind 集群

    执行 `make e2e_clean` 卸载 Kind 集群。

