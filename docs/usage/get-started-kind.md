# Kind Quick Start

[**简体中文**](./get-started-kind-zh_CN.md) | **English**

Kind is a tool for running local Kubernetes clusters using Docker container nodes. kdoctor provides scripts for installing Kind clusters to quickly build a Kind cluster with nginx, ingress, and loadbalancer, which you can use for kdoctor testing and experience.

## Prerequisites

* [Go](https://go.dev/) is installed.

## Clustering kdoctor on Kind

1. Clone the kdoctor repository to your localhost and go to the root directory of your kdoctor project.

    ```bash
    https://github.com/kdoctor-io/kdoctor.git && cd kdoctor
    cd kdoctor
    ```

2. Execute ``make checkBin`` to check if the development tools on the localhost meet the requirements for deploying a Kind cluster with kdoctor, and if the components are missing, they will be installed for you.

3. Get the latest image of kdoctor by using the following method.

    ```bash
    ~# KDOCTOR_LATEST_IMAGE_TAG=$(curl -s https://api.github.com/repos/kdoctor-io/kdoctor/releases | jq -r '. [].tag_name | select(("^v1.[0-9]*. [0-9]*$"))' | head -n 1)
    ```

4. Execute the following command to create the Kind cluster and install metallb, contour, nginx, and kdoctor for you.

    ```bash
    ~# make e2e_init -e PROJECT_IMAGE_VERSION=KDOCTOR_LATEST_IMAGE_TAG
    ```

    Note: If you are a domestic user, you can use the following command to avoid pulling the image failures.

    ```bash
    ~# make e2e_init -e E2E_SPIDERPOOL_TAG=$SPIDERPOOL_LATEST_IMAGE_TAG -e E2E_CHINA_IMAGE_REGISTRY=true
    ```

## Verify the installation

Configure KUBECONFIG for the Kind cluster for kubectl by executing the following command in the root directory of the kdoctor project.

   ```bash
   ~# export KUBECONFIG=$(pwd)/test/runtime/kubeconfig_kdoctor.config
   ```

You can see the following output:

   ```bash
   ~# kubectl get nodes 
   NAME STATUS ROLES AGE VERSION
   kdoctor-control-plane Ready control-plane 7m3s v1.27.1
   kdoctor-worker Ready <none> 6m42s v1.27.1 
   ~# kubectll get po -n kdoctor
   NAME READY STATUS RESTARTS AGE
   kdoctor-controller-686b75d6d7-ktctx 1/1 Running 0 2m33s
   ```

Next you can set up tasks as you see fit [AppHttpHealthy](./apphttphealthy.md), [NetReach](./netreach.md), [NetDns](./netdns.md)

## Uninstalling

* To uninstall the Kind cluster
    Execute `make e2e_clean` to uninstall the Kind cluster.
