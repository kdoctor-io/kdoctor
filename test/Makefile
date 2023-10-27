include ../Makefile.defs


.PHONY: init_kind_env
init_kind_env:
	make init_one_kind -e KIND_CONFIG_PATH=./kindconfig/global-kind.yaml -e KIND_CLUSTER_NAME=$(E2E_KIND_CLUSTER_NAME) -e KIND_KUBECONFIG=$(E2E_KIND_KUBECONFIG_PATH)
	make install_proscope


.PHONY: init_one_kind
init_one_kind: KIND_CONFIG_PATH ?=
init_one_kind: KIND_CLUSTER_NAME ?=
init_one_kind: KIND_KUBECONFIG ?=
init_one_kind: checkBin clean
	@echo "================== init kind cluster $(KIND_CLUSTER_NAME) KIND_CONFIG_PATH=$(KIND_CONFIG_PATH) KIND_KUBECONFIG=$(KIND_KUBECONFIG) E2E_IP_FAMILY=$(E2E_IP_FAMILY)"
	[ -n $(KIND_CLUSTER_NAME) ] || { echo "error, miss KIND_CLUSTER_NAME " ; exit 1 ; }
	[ -f $(KIND_CONFIG_PATH) ] || { echo "error, miss file KIND_CONFIG_PATH=$(KIND_CONFIG_PATH)" ; exit 1 ; }
	- mkdir -p $(E2E_RUNTIME_DIR) || true
	NEW_KIND_YAML=$(E2E_RUNTIME_DIR)/kind_config_$(KIND_CLUSTER_NAME).yaml ;\
		INSERT_LINE=` grep "insert subnet inform" $(KIND_CONFIG_PATH) -n | awk -F':' '{print $$1}' ` ; \
		echo "insert after line $${INSERT_LINE}" ;\
		sed  ''"$${INSERT_LINE}"' a \  ipFamily: $(E2E_IP_FAMILY)' $(KIND_CONFIG_PATH) > $${NEW_KIND_YAML} ; \
		if [ "$(E2E_IP_FAMILY)" == "ipv4" ] ; then \
			sed -i  ''"$${INSERT_LINE}"' a \  podSubnet: "$(E2E_KIND_IPV4_POD_CIDR)"' $${NEW_KIND_YAML} ;\
			sed -i  ''"$${INSERT_LINE}"' a \  serviceSubnet: "$(E2E_KIND_IPV4_SERVICE_CIDR)"' $${NEW_KIND_YAML} ;\
		elif [ "$(E2E_IP_FAMILY)" == "ipv6" ] ; then \
			sed -i  ''"$${INSERT_LINE}"' a \  podSubnet: "$(E2E_KIND_IPV6_POD_CIDR)"' $${NEW_KIND_YAML} ; \
			sed -i  ''"$${INSERT_LINE}"' a \  serviceSubnet: "$(E2E_KIND_IPV6_SERVICE_CIDR)"' $${NEW_KIND_YAML} ; \
		else \
			sed -i  ''"$${INSERT_LINE}"' a \  podSubnet: "$(E2E_KIND_IPV4_POD_CIDR),$(E2E_KIND_IPV6_POD_CIDR)"' $${NEW_KIND_YAML}  ; \
			sed -i  ''"$${INSERT_LINE}"' a \  serviceSubnet: "$(E2E_KIND_IPV4_SERVICE_CIDR),$(E2E_KIND_IPV6_SERVICE_CIDR)"' $${NEW_KIND_YAML}  ; \
  		fi
	- sysctl -w net.ipv6.conf.all.disable_ipv6=0 || true
	- sysctl -w fs.inotify.max_user_watches=524288 || true
	- sysctl -w fs.inotify.max_user_instances=8192  || true
	- kind delete cluster --name  $(KIND_CLUSTER_NAME)
	KIND_OPTION="" ; \
		[ -n "$(E2E_KIND_NODE_IMAGE)" ] && KIND_OPTION=" --image $(E2E_KIND_NODE_IMAGE) " && echo "setup kind with E2E_KIND_NODE_IMAGE=$(E2E_KIND_NODE_IMAGE)"; \
		kind create cluster --name  $(KIND_CLUSTER_NAME) --config $(E2E_RUNTIME_DIR)/kind_config_$(KIND_CLUSTER_NAME).yaml --kubeconfig $(KIND_KUBECONFIG) $${KIND_OPTION}
	- kubectl --kubeconfig $(KIND_KUBECONFIG) taint nodes --all node-role.kubernetes.io/master- || true
	- kubectl --kubeconfig $(KIND_KUBECONFIG) taint nodes --all node-role.kubernetes.io/control-plane- || true
	@echo "===================== deploy prometheus CRD ========== "
	# https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/main/example/prometheus-operator-crd/monitoring.coreos.com_servicemonitors.yaml ; } \
	kubectl apply --kubeconfig $(KIND_KUBECONFIG)  -f ./yaml/monitoring.coreos.com_servicemonitors.yaml
	# https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/main/example/prometheus-operator-crd/monitoring.coreos.com_podmonitors.yaml ; } \
	kubectl apply --kubeconfig $(KIND_KUBECONFIG) -f ./yaml/monitoring.coreos.com_podmonitors.yaml
	# https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/main/example/prometheus-operator-crd/monitoring.coreos.com_prometheusrules.yaml ; } \
	kubectl apply --kubeconfig $(KIND_KUBECONFIG) -f ./yaml/monitoring.coreos.com_prometheusrules.yaml
	# https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/main/example/prometheus-operator-crd/monitoring.coreos.com_probes.yaml  ; } \
	kubectl apply --kubeconfig $(KIND_KUBECONFIG) -f ./yaml/monitoring.coreos.com_probes.yaml
	# https://raw.githubusercontent.com/grafana-operator/grafana-operator/master/deploy/manifests/latest/crds.yaml  ; } \
	kubectl apply --kubeconfig $(KIND_KUBECONFIG) -f ./yaml/grafanadashboards.yaml
	echo "show kubernetes node image " && docker ps
	@echo "========================================================"
	@echo "   deploy kind cluster $(KIND_CLUSTER_NAME)             "
	@echo "   export KUBECONFIG=$(KIND_KUBECONFIG)                 "
	@echo "   kubectl get pod -o wide -A                           "
	@echo "========================================================"


.PHONY: checkBin
checkBin:
	$(ROOT_DIR)/test/scripts/installE2eTools.sh

.PHONY: install_proscope
install_proscope:
	if [ -n "$(PYROSCOPE_LOCAL_PORT)" ] ; then \
  		echo "install proscope " ; \
		docker stop $(PYROSCOPE_CONTAINER_NAME) &>/dev/null || true ; \
		docker rm $(PYROSCOPE_CONTAINER_NAME) &>/dev/null || true ; \
		ServerAddress=$$(docker network inspect kind -f {{\(index\ $$.IPAM.Config\ 0\).Gateway}}) ; \
		echo "setup pyroscope on $${ServerAddress}:$(PYROSCOPE_LOCAL_PORT)" ; \
		docker run -d --name $(PYROSCOPE_CONTAINER_NAME) -p $(PYROSCOPE_LOCAL_PORT):4040 $(PYROSCOPE_IMAGE_NAME) server ; \
		echo "finish setuping pyroscope " ; \
      fi



#==================

# this will auto tag github ci image : agent:xxx -> github.com/kdoctor-io/kdoctor/agent:xxx
.PHONY: check_images_ready
check_images_ready:
	echo "check image  " ; \
	IMAGE_LIST=` helm template test $(ROOT_DIR)/charts --set global.imageTagOverride=$(PROJECT_IMAGE_VERSION)  | grep " image: " | tr -d '"'| awk '{print $$2}' ` ; \
	if [ -z "$${IMAGE_LIST}" ] ; then \
		echo "warning, failed to find image from chart " ; \
		exit 1 ;\
	else \
		echo "find image from chart : $${IMAGE_LIST} " ; \
		for IMAGE in $${IMAGE_LIST} ; do \
		  	echo "try to find image $${IMAGE} " ; \
			EXIST=` docker images | awk '{printf("%s:%s\n",$$1,$$2)}' | grep "$${IMAGE}" ` || true ; \
			if [ -z "$${EXIST}" ] ; then \
					CI_IMAGE=$${IMAGE##*/} ; \
			  		echo "try to find github CI image $${CI_IMAGE} " ; \
			  		EXIST=` docker images | awk '{printf("%s:%s\n",$$1,$$2)}' | grep "$${CI_IMAGE}" ` || true ; \
			  		if [ -z "$${EXIST}" ] ; then \
			  			echo "error, failed to find image $${IMAGE}" ; \
			  			echo "error, failed to find image $${CI_IMAGE}" ; \
			  			exit 1 ; \
			  		fi ; \
			  		docker tag $${CI_IMAGE} $${IMAGE} ; \
			fi ;\
			echo "image exists: $${IMAGE}" ; \
		done ; \
		docker images ; \
	fi


# install kdoctor on global cluster
.PHONY: deploy_project
deploy_project: KIND_KUBECONFIG ?= $(E2E_KIND_KUBECONFIG_PATH)
deploy_project: KIND_CLUSTER_NAME ?= $(E2E_KIND_CLUSTER_NAME)
deploy_project:
	echo "try to load local image tag $(PROJECT_IMAGE_VERSION) " ; \
	IMAGE_LIST=` helm template test $(ROOT_DIR)/charts --set global.imageTagOverride=$(PROJECT_IMAGE_VERSION)  | grep ' image: ' | tr -d '"' | awk '{print $$2}'  | sort | uniq | tr '\n' ' '  ` ; \
	if [ -z "$${IMAGE_LIST}" ] ; then \
		echo "warning, failed to find image from chart " ; \
	else \
		echo "found image from chart : $${IMAGE_LIST} " ; \
		for IMAGE in $${IMAGE_LIST} ; do \
			EXIST=` docker images | awk '{printf("%s:%s\n",$$1,$$2)}' | grep "$${IMAGE}" ` || true ; \
			if [ -z "$${EXIST}" ] ; then \
			  echo "docker pull $${IMAGE} to local" ; \
			  docker pull $${IMAGE} ; \
			fi ;\
			echo "load local image $${IMAGE} " ; \
			kind load docker-image $${IMAGE}  --name $(KIND_CLUSTER_NAME)  ; \
		done ; \
	fi
	- helm --kubeconfig=$(KIND_KUBECONFIG) uninstall -n $(E2E_INSTALL_NAMESPACE) project || true
	HELM_OPTION="" ; \
    	if [ -n "$(PYROSCOPE_LOCAL_PORT)" ] ; then \
			echo "add env" ; \
			ServerAddress=$$(docker network inspect kind -f {{\(index\ $$.IPAM.Config\ 0\).Gateway}}) ; \
			HELM_OPTION+=" --set kdoctorAgent.extraEnv[0].name=ENV_PYROSCOPE_PUSH_SERVER_ADDRESS  --set kdoctorAgent.extraEnv[0].value=http://$${ServerAddress}:$(PYROSCOPE_LOCAL_PORT) " ; \
			HELM_OPTION+=" --set kdoctorController.extraEnv[0].name=ENV_PYROSCOPE_PUSH_SERVER_ADDRESS  --set kdoctorController.extraEnv[0].value=http://$${ServerAddress}:$(PYROSCOPE_LOCAL_PORT) " ; \
		fi ; \
		if [ "$(E2E_IP_FAMILY)" == "ipv4" ] ; then  \
			HELM_OPTION+=" --set feature.enableIPv4=true --set feature.enableIPv6=false " ; \
		elif [ "$(E2E_IP_FAMILY)" == "ipv6" ] ; then  \
			HELM_OPTION+=" --set feature.enableIPv4=false --set feature.enableIPv6=true " ; \
		else  \
			HELM_OPTION+=" --set feature.enableIPv4=true --set feature.enableIPv6=true " ; \
		fi ;\
		if [ "$(E2E_ENABLE_MULTUS)" == "true" ] ; then \
		    HELM_OPTION+=" --set kdoctorAgent.podAnnotations.v1\.multus-cni\.io/default-network=kube-system/k8s-pod-network " ; \
		    HELM_OPTION+=" --set kdoctorAgent.podAnnotations.k8s\.v1\.cni\.cncf\.io/networks=kube-system/ptp " ; \
		fi ; \
		if [ "$(E2E_ENABLE_METALLB)" == "true" ] ; then \
		    HELM_OPTION+=" --set kdoctorAgent.service.type=LoadBalancer " ; \
		else \
		    HELM_OPTION+=" --set kdoctorAgent.service.type=ClusterIP " ; \
		fi ; \
		if ( [ "$(E2E_ENABLE_CONTOUR)" == "true" ] || "$(E2E_ENABLE_NGINX)" == "true" ] ) && [ "$(E2E_IP_FAMILY)" != "ipv6" ]  ; then \
		    HELM_OPTION+=" --set kdoctorAgent.ingress.enable=true " ; \
		else \
		    HELM_OPTION+=" --set kdoctorAgent.ingress.enable=false " ; \
		fi ; \
		HELM_OPTION+=" --set feature.aggregateReport.enabled=true " ; \
		HELM_OPTION+=" --set feature.nethttp_defaultRequest_MaxQps=10000 " ; \
		HELM_OPTION+=" --set kdoctorAgent.resources.requests.cpu=200m " ; \
		HELM_OPTION+=" --set kdoctorAgent.resources.requests.memory=256Mi " ; \
		HELM_OPTION+=" --set feature.aggregateReport.controller.reportHostPath=/var/run/kdoctor/controller " ; \
		HELM_OPTION+=" --set kdoctorAgent.debug.logLevel=debug --set kdoctorController.debug.logLevel=debug " ; \
		HELM_OPTION+=" --set kdoctorAgent.prometheus.enabled=true --set kdoctorController.prometheus.enabled=true  " ; \
		HELM_OPTION+=" --set feature.aggregateReport.controller.pvc.enabled=true " ; \
		helm --kubeconfig=$(KIND_KUBECONFIG) install project $(ROOT_DIR)/charts \
				-n $(E2E_INSTALL_NAMESPACE) --create-namespace --wait --debug \
				--set global.imageTagOverride=$(PROJECT_IMAGE_VERSION) \
				$${HELM_OPTION} \
				|| { KIND_CLUSTER_NAME=$(KIND_CLUSTER_NAME) ./scripts/debugCluster.sh $(KIND_KUBECONFIG) "detail" $(E2E_INSTALL_NAMESPACE)  ; exit 1 ; } ; \
		exit 0


#=========================

.PHONY: deploy_multus
deploy_multus: KIND_KUBECONFIG ?= $(E2E_KIND_KUBECONFIG_PATH)
deploy_multus: KIND_CLUSTER_NAME ?= $(E2E_KIND_CLUSTER_NAME)
deploy_multus:
	if [ "$(E2E_ENABLE_MULTUS)" == "true" ] ; then \
  		  echo "install multus" ; \
  		  make install_multus -e KIND_KUBECONFIG=$(KIND_KUBECONFIG) -e KIND_CLUSTER_NAME=$(KIND_CLUSTER_NAME) ; \
  		else \
  		  echo "ignore multus" ; \
  		fi

.PHONY: install_multus
install_multus: KIND_KUBECONFIG ?= $(E2E_KIND_KUBECONFIG_PATH)
install_multus: KIND_CLUSTER_NAME ?= $(E2E_KIND_CLUSTER_NAME)
install_multus:
	@echo "beign to install multus"
	- helm repo remove $(MULTUS_REPO_NAME)
	[ -z "$(E2E_HELM_HTTP_PROXY)" ] || export https_proxy=$(E2E_HELM_HTTP_PROXY) ; \
		helm repo add $(MULTUS_REPO_NAME) $(MULTUS_CHART_REPO) ; \
		helm repo update $(MULTUS_REPO_NAME) ; \
		echo "try to load local image for multus $(MULTUS_CHART_VERSION) " ; \
		IMAGE_LIST=` helm template --version=$(MULTUS_CHART_VERSION)  $(MULTUS_REPO_NAME)/$(MULTUS_CHART_NAME) | grep ' image: ' | tr -d '"' | awk '{print $$2}'  | sort | uniq | tr '\n' ' '  ` ; \
		if [ -z "$${IMAGE_LIST}" ] ; then \
			echo "warning, failed to find image from chart template" ; \
		else \
			echo "found image from chart template: $${IMAGE_LIST} " ; \
			for IMAGE in $${IMAGE_LIST} ; do \
				EXIST=` docker images | awk '{printf("%s:%s\n",$$1,$$2)}' | grep "$${IMAGE}" ` || true ; \
				if [ -z "$${EXIST}" ] ; then \
				  echo "docker pull $${IMAGE} to local" ; \
				  docker pull $${IMAGE} ;\
				fi ;\
				echo "load local image $${IMAGE} " ; \
				kind load docker-image $${IMAGE}  --name $(KIND_CLUSTER_NAME)  ; \
			done ; \
		fi
	- helm --kubeconfig=$(KIND_KUBECONFIG) uninstall -n kube-system multus || true
	[ -z "$(E2E_HELM_HTTP_PROXY)" ] || export https_proxy=$(E2E_HELM_HTTP_PROXY) ; \
			helm install $(MULTUS_CHART_NAME)  $(MULTUS_REPO_NAME)/$(MULTUS_CHART_NAME)  \
			-n kube-system --kubeconfig=$(KIND_KUBECONFIG)  \
			--version=$(MULTUS_CHART_VERSION) --wait \
			--set sriov.manifests.enable=false \
			--set multus.config.cni_conf.clusterNetwork=k8s-pod-network
	kubectl --kubeconfig=$(KIND_KUBECONFIG) apply -f ./yaml/multus.yaml
	sleep 10
	#======= install cni and route on nodes
	- mkdir -p $(E2E_MULTUS_CNI_CONFIG_DIR)
	- rm -rf $(E2E_MULTUS_CNI_CONFIG_DIR)/*
	echo "install cni and route on nodes" ; \
		TMP=` docker exec $(E2E_KIND_MASTER_NODE_NAME) ip a s eth0 ` || { echo "error, failed to find master ip" ; exit 1 ; } ; \
		MASTER_IPV4_ADDR=` grep -oP '(?<=inet\s)[0-9]+(\.[0-9]+){3}' <<< "$${TMP}" ` || true ; \
		MASTER_IPV6_ADDR=` echo "$${TMP}" |  grep -v "scope link" | grep -oP '(?<=inet6\s)[0-9a-f:]+' ` || true ; \
		TMP=` docker exec $(E2E_KIND_WORKER_NODE_NAME) ip a s eth0 ` || { echo "error, failed to find worker ip" ; exit 1 ; } ; \
		WORKER_IPV4_ADDR=` grep -oP '(?<=inet\s)[0-9]+(\.[0-9]+){3}' <<< "$${TMP}" ` || true ; \
		WORKER_IPV6_ADDR=` echo "$${TMP}" |  grep -v "scope link" | grep -oP '(?<=inet6\s)[0-9a-f:]+' ` || true ; \
		echo "master ipv4=$${MASTER_IPV4_ADDR}, ipv6=$${MASTER_IPV6_ADDR} " ; \
		echo "worker ipv4=$${WORKER_IPV4_ADDR}, ipv6=$${WORKER_IPV6_ADDR} " ; \
		if [ "$(E2E_IP_FAMILY)" == "ipv4" ] ; then \
			  echo "install ipv4-only case" ; \
		  	  [ -n "$${MASTER_IPV4_ADDR}" ] || { echo "error, failed to find master ipv4" ; exit 1 ; } ; \
		  	  [ -n "$${WORKER_IPV4_ADDR}" ] || { echo "error, failed to find worker ipv4" ; exit 1 ; } ; \
			  sed 's?<<IP_RANGE>>?\[\{ "subnet": "172.80.1.0/24" \}\]?' ./yaml/cni.conflist | sed 's?<<NAME>>?k8s-pod-network?g' | sed 's?<<ROUTE>>?\{ "dst": "0.0.0.0/0" \}?' > $(E2E_MULTUS_CNI_CONFIG_DIR)/20-master-cni1.conflist ; \
			  sed 's?<<IP_RANGE>>?\[\{ "subnet": "172.90.1.0/24" \}\]?' ./yaml/cni.conflist | sed 's?<<NAME>>?ptp?g' | sed 's?<<ROUTE>>??' > $(E2E_MULTUS_CNI_CONFIG_DIR)/20-master-cni2.conflist ; \
			  sed 's?<<IP_RANGE>>?\[\{ "subnet": "172.80.2.0/24" \}\]?' ./yaml/cni.conflist | sed 's?<<NAME>>?k8s-pod-network?g' | sed 's?<<ROUTE>>?\{ "dst": "0.0.0.0/0" \}?' > $(E2E_MULTUS_CNI_CONFIG_DIR)/20-worker-cni1.conflist ; \
			  sed 's?<<IP_RANGE>>?\[\{ "subnet": "172.90.2.0/24" \}\]?' ./yaml/cni.conflist | sed 's?<<NAME>>?ptp?g' | sed 's?<<ROUTE>>??' > $(E2E_MULTUS_CNI_CONFIG_DIR)/20-worker-cni2.conflist ; \
			  docker cp $(E2E_MULTUS_CNI_CONFIG_DIR)/20-master-cni1.conflist $(E2E_KIND_MASTER_NODE_NAME):/etc/cni/net.d ; \
			  docker cp $(E2E_MULTUS_CNI_CONFIG_DIR)/20-master-cni2.conflist $(E2E_KIND_MASTER_NODE_NAME):/etc/cni/net.d ; \
			  docker cp $(E2E_MULTUS_CNI_CONFIG_DIR)/20-worker-cni1.conflist $(E2E_KIND_WORKER_NODE_NAME):/etc/cni/net.d ; \
			  docker cp $(E2E_MULTUS_CNI_CONFIG_DIR)/20-worker-cni2.conflist $(E2E_KIND_WORKER_NODE_NAME):/etc/cni/net.d ; \
			  docker exec $(E2E_KIND_MASTER_NODE_NAME) ip r add 172.80.2.0/24 via $${WORKER_IPV4_ADDR} ; \
			  docker exec $(E2E_KIND_MASTER_NODE_NAME) ip r add 172.90.2.0/24 via $${WORKER_IPV4_ADDR} ; \
			  docker exec $(E2E_KIND_WORKER_NODE_NAME) ip r add 172.80.1.0/24 via $${MASTER_IPV4_ADDR} ; \
			  docker exec $(E2E_KIND_WORKER_NODE_NAME) ip r add 172.90.1.0/24 via $${MASTER_IPV4_ADDR} ; \
		elif [ "$(E2E_IP_FAMILY)" == "ipv6" ] ; then \
			  echo "install ipv6-only case" ; \
		  	  [ -n "$${MASTER_IPV6_ADDR}" ] || { echo "error, failed to find master ipv6" ; exit 1 ; } ; \
		  	  [ -n "$${WORKER_IPV6_ADDR}" ] || { echo "error, failed to find worker ipv6" ; exit 1 ; } ; \
			  sed 's?<<IP_RANGE>>?\[\{ "subnet": "fd80:1::/64" \}\]?' ./yaml/cni.conflist | sed 's?<<NAME>>?k8s-pod-network?' | sed 's?<<ROUTE>>?\{ "dst": "::/0" \}?' > $(E2E_MULTUS_CNI_CONFIG_DIR)/20-master-cni1.conflist ; \
			  sed 's?<<IP_RANGE>>?\[\{ "subnet": "fd90:1::/64" \}\]?' ./yaml/cni.conflist | sed 's?<<NAME>>?ptp?' | sed 's?<<ROUTE>>??' > $(E2E_MULTUS_CNI_CONFIG_DIR)/20-master-cni1.conflist ; \
			  sed 's?<<IP_RANGE>>?\[\{ "subnet": "fd80:2::/64" \}\]?' ./yaml/cni.conflist | sed 's?<<NAME>>?k8s-pod-network?' | sed 's?<<ROUTE>>?\{ "dst": "::/0" \}?' > $(E2E_MULTUS_CNI_CONFIG_DIR)/20-worker-cni1.conflist ; \
			  sed 's?<<IP_RANGE>>?\[\{ "subnet": "fd90:2::/64" \}\]?' ./yaml/cni.conflist | sed 's?<<NAME>>?ptp?' | sed 's?<<ROUTE>>??' > $(E2E_MULTUS_CNI_CONFIG_DIR)/20-worker-cni2.conflist ; \
			  docker cp $(E2E_MULTUS_CNI_CONFIG_DIR)/20-master-cni1.conflist $(E2E_KIND_MASTER_NODE_NAME):/etc/cni/net.d ; \
			  docker cp $(E2E_MULTUS_CNI_CONFIG_DIR)/20-master-cni2.conflist $(E2E_KIND_MASTER_NODE_NAME):/etc/cni/net.d ; \
			  docker cp $(E2E_MULTUS_CNI_CONFIG_DIR)/20-worker-cni1.conflist $(E2E_KIND_WORKER_NODE_NAME):/etc/cni/net.d ; \
			  docker cp $(E2E_MULTUS_CNI_CONFIG_DIR)/20-worker-cni2.conflist $(E2E_KIND_WORKER_NODE_NAME):/etc/cni/net.d ; \
			  docker exec $(E2E_KIND_MASTER_NODE_NAME) ip -6 r add fd80:2::/64 via $${WORKER_IPV6_ADDR} ; \
			  docker exec $(E2E_KIND_MASTER_NODE_NAME) ip -6 r add fd90:2::/64 via $${WORKER_IPV6_ADDR} ; \
			  docker exec $(E2E_KIND_WORKER_NODE_NAME) ip -6 r add fd80:1::/64 via $${MASTER_IPV6_ADDR} ; \
			  docker exec $(E2E_KIND_WORKER_NODE_NAME) ip -6 r add fd90:1::/64 via $${MASTER_IPV6_ADDR} ; \
  		else \
			  echo "install dual-stack case" ; \
		  	  [ -n "$${MASTER_IPV4_ADDR}" ] || { echo "error, failed to find master ipv4" ; exit 1 ; } ; \
		  	  [ -n "$${WORKER_IPV4_ADDR}" ] || { echo "error, failed to find worker ipv4" ; exit 1 ; } ; \
		  	  [ -n "$${MASTER_IPV6_ADDR}" ] || { echo "error, failed to find master ipv6" ; exit 1 ; } ; \
		  	  [ -n "$${WORKER_IPV6_ADDR}" ] || { echo "error, failed to find worker ipv6" ; exit 1 ; } ; \
			  sed 's?<<IP_RANGE>>?\[\{ "subnet": "172.80.1.0/24" \}\],\[\{ "subnet": "fd80:1::/64" \}\]?' ./yaml/cni.conflist | sed 's?<<NAME>>?k8s-pod-network?' | sed 's?<<ROUTE>>?\{ "dst": "0.0.0.0/0" \},\{ "dst": "::/0" \}?' > $(E2E_MULTUS_CNI_CONFIG_DIR)/20-master-cni1.conflist ; \
			  sed 's?<<IP_RANGE>>?\[\{ "subnet": "172.90.1.0/24" \}\],\[\{ "subnet": "fd90:1::/64" \}\]?' ./yaml/cni.conflist | sed 's?<<NAME>>?ptp?' | sed 's?<<ROUTE>>??' > $(E2E_MULTUS_CNI_CONFIG_DIR)/20-master-cni2.conflist ; \
			  sed 's?<<IP_RANGE>>?\[\{ "subnet": "172.80.2.0/24" \}\],\[\{ "subnet": "fd80:2::/64" \}\]?' ./yaml/cni.conflist | sed 's?<<NAME>>?k8s-pod-network?' | sed 's?<<ROUTE>>?\{ "dst": "0.0.0.0/0" \},\{ "dst": "::/0" \}?' > $(E2E_MULTUS_CNI_CONFIG_DIR)/20-worker-cni1.conflist ; \
			  sed 's?<<IP_RANGE>>?\[\{ "subnet": "172.90.2.0/24" \}\],\[\{ "subnet": "fd90:2::/64" \}\]?' ./yaml/cni.conflist | sed 's?<<NAME>>?ptp?' | sed 's?<<ROUTE>>??' > $(E2E_MULTUS_CNI_CONFIG_DIR)/20-worker-cni2.conflist ; \
			  docker cp $(E2E_MULTUS_CNI_CONFIG_DIR)/20-master-cni1.conflist $(E2E_KIND_MASTER_NODE_NAME):/etc/cni/net.d ; \
			  docker cp $(E2E_MULTUS_CNI_CONFIG_DIR)/20-master-cni2.conflist $(E2E_KIND_MASTER_NODE_NAME):/etc/cni/net.d ; \
			  docker cp $(E2E_MULTUS_CNI_CONFIG_DIR)/20-worker-cni1.conflist $(E2E_KIND_WORKER_NODE_NAME):/etc/cni/net.d ; \
			  docker cp $(E2E_MULTUS_CNI_CONFIG_DIR)/20-worker-cni2.conflist $(E2E_KIND_WORKER_NODE_NAME):/etc/cni/net.d ; \
			  docker exec $(E2E_KIND_MASTER_NODE_NAME) ip r add 172.80.2.0/24 via $${WORKER_IPV4_ADDR} ; \
			  docker exec $(E2E_KIND_MASTER_NODE_NAME) ip r add 172.90.2.0/24 via $${WORKER_IPV4_ADDR} ; \
			  docker exec $(E2E_KIND_WORKER_NODE_NAME) ip r add 172.80.1.0/24 via $${MASTER_IPV4_ADDR} ; \
			  docker exec $(E2E_KIND_WORKER_NODE_NAME) ip r add 172.90.1.0/24 via $${MASTER_IPV4_ADDR} ; \
			  docker exec $(E2E_KIND_MASTER_NODE_NAME) ip -6 r add fd80:2::/64 via $${WORKER_IPV6_ADDR} ; \
			  docker exec $(E2E_KIND_MASTER_NODE_NAME) ip -6 r add fd90:2::/64 via $${WORKER_IPV6_ADDR} ; \
			  docker exec $(E2E_KIND_WORKER_NODE_NAME) ip -6 r add fd80:1::/64 via $${MASTER_IPV6_ADDR} ; \
			  docker exec $(E2E_KIND_WORKER_NODE_NAME) ip -6 r add fd90:1::/64 via $${MASTER_IPV6_ADDR} ; \
  		fi

#=========================

.PHONY: deploy_metallb
deploy_metallb: KIND_KUBECONFIG ?= $(E2E_KIND_KUBECONFIG_PATH)
deploy_metallb: KIND_CLUSTER_NAME ?= $(E2E_KIND_CLUSTER_NAME)
deploy_metallb:
	[ "$(E2E_ENABLE_METALLB)" == "true" ] || { echo "ignore metallb" ; exit 0 ; } ; \
		make install_metallb -e KIND_KUBECONFIG=$(KIND_KUBECONFIG) -e KIND_CLUSTER_NAME=$(KIND_CLUSTER_NAME)


.PHONY: install_metallb
install_metallb: KIND_KUBECONFIG ?= $(E2E_KIND_KUBECONFIG_PATH)
install_metallb: KIND_CLUSTER_NAME ?= $(E2E_KIND_CLUSTER_NAME)
install_metallb:
	@echo "beign to install metallb"
	- helm repo remove $(METALLB_REPO_NAME)
	- helm --kubeconfig=$(KIND_KUBECONFIG) uninstall -n kube-system $(METALLB_CHART_NAME) || true
	[ -z "$(E2E_HELM_HTTP_PROXY)" ] || export https_proxy=$(E2E_HELM_HTTP_PROXY) ; \
		helm repo add $(METALLB_REPO_NAME) $(METALLB_CHART_REPO) ; \
		helm repo update $(METALLB_REPO_NAME) ; \
		HELM_OPTIONS="	--set controller.image.repository=quay.$(E2E_PROXY_REGISTER)io/metallb/controller \
                      	--set speaker.image.repository=quay.$(E2E_PROXY_REGISTER)io/metallb/speaker \
                      	--set speaker.frr.image.repository=quay.$(E2E_PROXY_REGISTER)io/frrouting/frr  --set speaker.frr.enabled=true " ; \
		IMAGE_LIST=` helm template --version=$(METALLB_CHART_VERSION)  $(METALLB_REPO_NAME)/$(METALLB_CHART_NAME) $${HELM_OPTIONS} | grep ' image: ' | tr -d '"' | awk '{print $$2}'  | sort | uniq | tr '\n' ' '  ` ; \
		if [ -z "$${IMAGE_LIST}" ] ; then \
			echo "warning, failed to find image from chart template" ; \
		else \
			echo "found image from chart template: $${IMAGE_LIST} " ; \
			for IMAGE in $${IMAGE_LIST} ; do \
				EXIST=` docker images | awk '{printf("%s:%s\n",$$1,$$2)}' | grep "$${IMAGE}" ` || true ; \
				if [ -z "$${EXIST}" ] ; then \
				  echo "docker pull $${IMAGE} to local" ; \
				  docker pull $${IMAGE} ;\
				fi ;\
				echo "load local image $${IMAGE} " ; \
				kind load docker-image $${IMAGE}  --name $(KIND_CLUSTER_NAME)  ; \
			done ; \
		fi ; \
		helm install -n kube-system $(METALLB_CHART_NAME)  $(METALLB_REPO_NAME)/$(METALLB_CHART_NAME) \
			-n kube-system --kubeconfig=$(KIND_KUBECONFIG) \
			--version=$(METALLB_CHART_VERSION) --wait  --timeout 5m \
			--set speaker.nodeSelector."kubernetes\.io/hostname"=$(E2E_KIND_CLUSTER_NAME)-control-plane \
			$${HELM_OPTIONS}
	echo "apply ippool" ; \
		docker network inspect kind ; \
		if [ "$(E2E_IP_FAMILY)" == "ipv4" ] ; then  \
			Subnet1=$$(docker network inspect kind -f {{\(index\ $$.IPAM.Config\ 0\).Subnet}}) ; \
			Subnet1=$${Subnet1%%/*} ; \
			if ! grep -iE "[a-f]" <<< "$${Subnet1}" ; then \
			  	IPPOO1="$${Subnet1%0}" ; \
			 	IPPOO1="$${IPPOO1}50-$${IPPOO1}90" ; \
			else \
				Subnet2=$$(docker network inspect kind -f {{\(index\ $$.IPAM.Config\ 1\).Subnet}}) ; \
				Subnet2=$${Subnet2%%/*} ; \
				if ! grep -iE "[a-f]" <<< "$${Subnet2}" ; then \
				  	IPPOO1="$${Subnet2%0}" ; \
					IPPOO1="$${IPPOO1}50-$${IPPOO1}90" ; \
				else \
					echo "failed to find node ipv4 subnet" ; \
					exit 1 ; \
				fi ; \
  			fi ; \
			echo "IPPOO1: $${IPPOO1}" ; \
		  	cat ./yaml/metallb.yaml \
		  		| sed '/addresses:/ a\    - '"$${IPPOO1}"''  \
		  		| kubectl --kubeconfig=$(KIND_KUBECONFIG) apply -f - ; \
		elif [ "$(E2E_IP_FAMILY)" == "ipv6" ] ; then  \
			Subnet1=$$(docker network inspect kind -f {{\(index\ $$.IPAM.Config\ 0\).Subnet}}) ; \
			Subnet1=$${Subnet1%%/*} ; \
			if grep -iE "[a-f]" <<< "$${Subnet1}" ; then \
			 	IPPOO1="$${Subnet1}50-$${Subnet1}90" ; \
			else \
				Subnet2=$$(docker network inspect kind -f {{\(index\ $$.IPAM.Config\ 1\).Subnet}}) ; \
				Subnet2=$${Subnet2%%/*} ; \
				if grep -iE "[a-f]" <<< "$${Subnet2}" ; then \
					IPPOO1="$${Subnet2}50-$${Subnet2}90" ; \
				else \
					echo "failed to find node ipv6 subnet" ; \
					exit 1 ; \
				fi ; \
  			fi ; \
			echo "IPPOO1: $${IPPOO1}" ; \
		  	cat ./yaml/metallb.yaml \
		  		| sed '/addresses:/ a\    - '"$${IPPOO1}"''  \
		  		| kubectl --kubeconfig=$(KIND_KUBECONFIG) apply -f - ; \
		else  \
			Subnet1=$$(docker network inspect kind -f {{\(index\ $$.IPAM.Config\ 0\).Subnet}}) ; \
			Subnet1=$${Subnet1%%/*} ; \
			if grep -iE "[a-f]" <<< "$${Subnet1}" ; then \
			 	IPPOO1="$${Subnet1}50-$${Subnet1}90" ; \
			else \
			 	IPPOO1="$${Subnet1%0}" ; \
			 	IPPOO1="$${IPPOO1}50-$${IPPOO1}90" ; \
  			fi ; \
			Subnet2=$$(docker network inspect kind -f {{\(index\ $$.IPAM.Config\ 1\).Subnet}}) ; \
			Subnet2=$${Subnet2%%/*} ; \
			if grep -iE "[a-f]" <<< "$${Subnet2}" ; then \
			 	IPPOO2="$${Subnet2}50-$${Subnet2}90" ; \
			else \
			 	IPPOO2="$${Subnet2%0}" ; \
			 	IPPOO2="$${IPPOO2}50-$${IPPOO2}90" ; \
  			fi ; \
			echo "IPPOO1: $${IPPOO1}" ; \
			echo "IPPOO2: $${IPPOO2}" ; \
		  	cat ./yaml/metallb.yaml \
		  		| sed '/addresses:/ a\    - '"$${IPPOO1}"''  \
		  		| sed '/addresses:/ a\    - '"$${IPPOO2}"'' ; \
		  	cat ./yaml/metallb.yaml \
		  		| sed '/addresses:/ a\    - '"$${IPPOO1}"''  \
		  		| sed '/addresses:/ a\    - '"$${IPPOO2}"''  \
		  		| kubectl --kubeconfig=$(KIND_KUBECONFIG) apply -f - ; \
		fi

#=============================

# contour does not support ipv6-only
.PHONY: deploy_contour
deploy_contour: KIND_KUBECONFIG ?= $(E2E_KIND_KUBECONFIG_PATH)
deploy_contour: KIND_CLUSTER_NAME ?= $(E2E_KIND_CLUSTER_NAME)
deploy_contour:
	[ "$(E2E_ENABLE_CONTOUR)" == "true" ] || { echo "ignore contour" ; exit 0 ; } ; \
		if [ "$(E2E_IP_FAMILY)" == "ipv6" ]  ; then \
				echo "contour does not support ipv6 only " ;  \
				exit 0 ; \
		fi ; \
		make install_contour -e KIND_KUBECONFIG=$(KIND_KUBECONFIG) -e KIND_CLUSTER_NAME=$(KIND_CLUSTER_NAME)


.PHONY: install_contour
install_contour: KIND_KUBECONFIG ?= $(E2E_KIND_KUBECONFIG_PATH)
install_contour: KIND_CLUSTER_NAME ?= $(E2E_KIND_CLUSTER_NAME)
install_contour:
	@echo "beign to install contour"
	- helm repo remove $(CONTOUR_REPO_NAME)
	- helm --kubeconfig=$(KIND_KUBECONFIG) uninstall -n kube-system $(CONTOUR_CHART_NAME) || true
	[ -z "$(E2E_HELM_HTTP_PROXY)" ] || export https_proxy=$(E2E_HELM_HTTP_PROXY) ; \
		helm repo add $(CONTOUR_REPO_NAME) $(CONTOUR_CHART_REPO) ; \
		helm repo update $(CONTOUR_REPO_NAME) ; \
		HELM_OPTIONS=" --set global.imageRegistry=docker.$(E2E_PROXY_REGISTER)io  "  ; \
		IMAGE_LIST=` helm template --version=$(CONTOUR_CHART_VERSION)  $(CONTOUR_REPO_NAME)/$(CONTOUR_CHART_NAME) $${HELM_OPTIONS} | grep ' image: ' | tr -d '"' | awk '{print $$2}'  | sort | uniq | tr '\n' ' '  ` ; \
		if [ -z "$${IMAGE_LIST}" ] ; then \
			echo "warning, failed to find image from chart template" ; \
		else \
			echo "found image from chart template: $${IMAGE_LIST} " ; \
			for IMAGE in $${IMAGE_LIST} ; do \
				EXIST=` docker images | awk '{printf("%s:%s\n",$$1,$$2)}' | grep "$${IMAGE}" ` || true ; \
				if [ -z "$${EXIST}" ] ; then \
				  echo "docker pull $${IMAGE} to local" ; \
				  docker pull $${IMAGE} ;\
				fi ;\
				echo "load local image $${IMAGE} " ; \
				kind load docker-image $${IMAGE}  --name $(KIND_CLUSTER_NAME)  ; \
			done ; \
		fi ; \
		helm install -n kube-system $(CONTOUR_CHART_NAME)  $(CONTOUR_REPO_NAME)/$(CONTOUR_CHART_NAME) \
			-n kube-system --kubeconfig=$(KIND_KUBECONFIG) \
			--version=$(CONTOUR_CHART_VERSION) --wait  --timeout 5m  $${HELM_OPTIONS} \
			  --set envoy.kind=daemonset \
			  --set envoy.service.type=LoadBalancer \
			  --set envoy.service.externalTrafficPolicy=Local \
			  --set envoy.hostNetwork=false \
			  --set envoy.dnsPolicy=ClusterFirst \
			  --set ingress.enabled=true \
			  --set envoy.nodeSelector."kubernetes\.io/hostname"=$(E2E_KIND_CLUSTER_NAME)-control-plane


# contour does not support ipv6-only
.PHONY: deploy_nginx
deploy_nginx: KIND_KUBECONFIG ?= $(E2E_KIND_KUBECONFIG_PATH)
deploy_nginx: KIND_CLUSTER_NAME ?= $(E2E_KIND_CLUSTER_NAME)
deploy_nginx:
	[ "$(E2E_ENABLE_NGINX)" == "true" ] || { echo "ignore nginx" ; exit 0 ; } ; \
	if [ "$(E2E_IP_FAMILY)" == "ipv6" ]  ; then \
  			echo "nginx does not support ipv6 only " ; \
  			exit 0 ; \
  	  fi ; \
	  make install_nginx -e KIND_KUBECONFIG=$(KIND_KUBECONFIG) -e KIND_CLUSTER_NAME=$(KIND_CLUSTER_NAME)


.PHONY: install_nginx
install_nginx: KIND_KUBECONFIG ?= $(E2E_KIND_KUBECONFIG_PATH)
install_nginx: KIND_CLUSTER_NAME ?= $(E2E_KIND_CLUSTER_NAME)
install_nginx:
	@echo "beign to install nginx"
	- helm repo remove $(NGINX_REPO_NAME)
	- helm --kubeconfig=$(KIND_KUBECONFIG) uninstall -n kube-system $(NGINX_CHART_NAME) || true
	[ -z "$(E2E_HELM_HTTP_PROXY)" ] || export https_proxy=$(E2E_HELM_HTTP_PROXY) ; \
		helm repo add $(NGINX_REPO_NAME) $(NGINX_CHART_REPO) ; \
		helm repo update $(NGINX_REPO_NAME) ; \
		HELM_OPTION="  --set controller.image.registry=k8s.$(E2E_PROXY_REGISTER)io \
                       --set defaultBackend.image.registry=k8s.$(E2E_PROXY_REGISTER)io \
                       --set controller.admissionWebhooks.enabled=false \
                       --set controller.admissionWebhooks.patch.image.registry=k8s.$(E2E_PROXY_REGISTER)io \
                       --set controller.admissionWebhooks.patch.image.digest='' "; \
		IMAGE_LIST=` helm template --version=$(NGINX_CHART_VERSION)  $(NGINX_REPO_NAME)/$(NGINX_CHART_NAME) $${HELM_OPTION} | grep ' image: ' | tr -d '"' | awk '{print $$2}'  | sort | uniq | tr '\n' ' '  ` ; \
		if [ -z "$${IMAGE_LIST}" ] ; then \
			echo "warning, failed to find image from chart template" ; \
		else \
			echo "found image from chart template: $${IMAGE_LIST} " ; \
			for IMAGE in $${IMAGE_LIST} ; do \
				EXIST=` docker images | awk '{printf("%s:%s\n",$$1,$$2)}' | grep "$${IMAGE}" ` || true ; \
				if [ -z "$${EXIST}" ] ; then \
				  echo "docker pull $${IMAGE} to local" ; \
				  docker pull $${IMAGE} ;\
				fi ;\
				echo "load local image $${IMAGE} " ; \
				kind load docker-image $${IMAGE}  --name $(KIND_CLUSTER_NAME)  ; \
			done ; \
		fi ; \
		helm install -n kube-system $(NGINX_CHART_NAME)  $(NGINX_REPO_NAME)/$(NGINX_CHART_NAME) \
			-n kube-system --kubeconfig=$(KIND_KUBECONFIG) \
			--version=$(NGINX_CHART_VERSION) --wait  --timeout 5m \
			  --set controller.kind=DaemonSet \
			  --set controller.nodeSelector."kubernetes\.io/hostname"=kdoctor-control-plane \
			  --set controller.admissionWebhooks.patch.enabled=true \
			  --set controller.hostNetwork=false \
			  --set controller.service.type=LoadBalancer \
			  --set controller.service.internal.externalTrafficPolicy=Local \
			  --set controller.ingressClassResource.default=true \
			  $${HELM_OPTION}

.PHONY: clean_kdoctor_crd
clean_kdoctor_crd: KIND_KUBECONFIG ?= $(E2E_KIND_KUBECONFIG_PATH)
clean_kdoctor_crd:
	@echo -e "\033[35m [clean kdoctor crd] \033[0m"
	@kubectl delete crd netdnses.kdoctor.io --kubeconfig $(KIND_KUBECONFIG)
	@kubectl delete crd apphttphealthies.kdoctor.io --kubeconfig $(KIND_KUBECONFIG)
	@kubectl delete crd netreaches.kdoctor.io --kubeconfig $(KIND_KUBECONFIG)
	@ALL_CRD=`kubectl get crd -n $(E2E_INSTALL_NAMESPACE) --kubeconfig $(KIND_KUBECONFIG)  | awk '{print $$1 } ' ` ; FAIL="false" ; FOUND="false"; \
	for CRD in $$ALL_CRD ; do \
		echo $${CRD} | grep -q 'kdoctor.io' && FOUND="true" && echo "error, failed to clean kdoctor crd: $${CRD}" ; \
	    if [ "$$FOUND" == "true" ] ; then \
           echo -e "\033[35m [===== describe crd: $${CRD} =====] \033[0m" ; \
		   kubectl describe crd $${CRD} --kubeconfig $(KIND_KUBECONFIG) ; \
		   echo -e "\033[35m [===== get crd: $${CRD} -ojson =====] \033[0m" ; \
           kubectl get crd $${CRD} -ojson --kubeconfig $(KIND_KUBECONFIG) ; \
           FAIL="true" ; \
           FOUND="false" ; \
		fi ; \
	done ; \
	[ "$$FAIL" == "true" ] && echo "kdoctor crd cleanup failed" && exit 1 ; \
	echo -e "\033[35m [done] all kdoctor crd cleaned up !!! \033[0m"

.PHONY: helm_uninstall_kdoctor
helm_uninstall_kdoctor: KIND_KUBECONFIG ?= $(E2E_KIND_KUBECONFIG_PATH)
helm_uninstall_kdoctor:
	@echo -e "\033[35m [helm uninstall kdoctor] \033[0m"
	helm uninstall project --wait --debug -n $(E2E_INSTALL_NAMESPACE) \
	     --kubeconfig $(KIND_KUBECONFIG)  || { KIND_CLUSTER_NAME=$(KIND_CLUSTER_NAME) ./scripts/debugCluster.sh $(KIND_KUBECONFIG) "detail"   ; exit 1 ; } ; \

.PHONY: uninstall_kdoctor
uninstall_kdoctor:
	@echo -e "\033[35m [uninstall kdoctor] \033[0m"
	@ make helm_uninstall_kdoctor
	@ make clean_kdoctor_crd


#==========================

.PHONY: install_apiserver_token
install_apiserver_token: KIND_KUBECONFIG ?= $(E2E_KIND_KUBECONFIG_PATH)
install_apiserver_token: KIND_CLUSTER_NAME ?= $(E2E_KIND_CLUSTER_NAME)
install_apiserver_token:
	@echo "---------- install apiserver token"
	sed 's?<<TEST_NAMESPACE>>?$(E2E_INSTALL_NAMESPACE)?' yaml/kdoctortoken.yaml | kubectl --kubeconfig=$(KIND_KUBECONFIG) apply -f -


# http and dns server for e2e test
.PHONY: deploy_test_server
deploy_test_server: KIND_KUBECONFIG ?= $(E2E_KIND_KUBECONFIG_PATH)
deploy_test_server: KIND_CLUSTER_NAME ?= $(E2E_KIND_CLUSTER_NAME)
deploy_test_server:
	echo "try to load local image tag $(PROJECT_IMAGE_VERSION) " ; \
	IMAGE_LIST=` helm template test $(APP_CHART_DIR) --set image.tag=$(PROJECT_IMAGE_VERSION)  | grep ' image: ' | tr -d '"' | awk '{print $$2}'  | sort | uniq | tr '\n' ' '  ` ; \
	if [ -z "$${IMAGE_LIST}" ] ; then \
		echo "warning, failed to find image from chart " ; \
	else \
		echo "found image from chart : $${IMAGE_LIST} " ; \
		for IMAGE in $${IMAGE_LIST} ; do \
			EXIST=` docker images | awk '{printf("%s:%s\n",$$1,$$2)}' | grep "$${IMAGE}" ` || true ; \
			if [ -z "$${EXIST}" ] ; then \
			  echo "docker pull $${IMAGE} to local" ; \
			  docker pull $${IMAGE} ; \
			fi ;\
			echo "load local image $${IMAGE} " ; \
			kind load docker-image $${IMAGE}  --name $(KIND_CLUSTER_NAME)  ; \
		done ; \
	fi
	- helm --kubeconfig=$(KIND_KUBECONFIG) uninstall -n $(E2E_INSTALL_NAMESPACE) server || true
		helm --kubeconfig=$(KIND_KUBECONFIG) install server $(APP_CHART_DIR) \
				-n $(E2E_INSTALL_NAMESPACE) --create-namespace --wait --debug \
				--set image.tag=$(PROJECT_IMAGE_VERSION) \
				|| { KIND_CLUSTER_NAME=$(KIND_CLUSTER_NAME) ./scripts/debugCluster.sh $(KIND_KUBECONFIG) "detail" $(E2E_INSTALL_NAMESPACE)  ; exit 1 ; } ; \
		exit 0


.PHONY: clean
clean:
	-@ kind delete cluster --name $(E2E_KIND_CLUSTER_NAME)
	-@ rm -rf $(E2E_RUNTIME_DIR)
	-@ docker stop $(PYROSCOPE_CONTAINER_NAME) &>/dev/null
	-@ docker rm $(PYROSCOPE_CONTAINER_NAME) &>/dev/null



#============ e2e ====================
.PHONY: e2e_test
e2e_test: KIND_CLUSTER_NAME ?= $(E2E_KIND_CLUSTER_NAME)
e2e_test: KIND_KUBECONFIG ?= $(E2E_KIND_KUBECONFIG_PATH)
e2e_test:
	@echo -e "\033[35m Run e2e test on the cluster $(KIND_CLUSTER_NAME) \033[0m "
	@ echo -e "\033[35m [E2E] Run E2E with ginkgo label=$(E2E_GINKGO_LABELS) , timeout=$(E2E_TIMEOUT) GINKGO_OPTION=$(E2E_GINKGO_OPTION) \033[0m"
	@  NODE_LIST=` docker ps | egrep " kindest/node.* $(KIND_CLUSTER_NAME)-(control|worker)" | awk '{print $$NF }' ` ; \
		[ -n "$$NODE_LIST" ] || { echo "error, failed to find any kind nodes, please setup kind cluster $(KIND_CLUSTER_NAME) first" ; exit 1 ; } ; \
		NODE_LIST=` echo "$${NODE_LIST}" | tr -d ' ' | tr '\n' ',' ` ; \
		NODE_LIST=$${NODE_LIST%%,} ; \
		echo "find cluster node: $${NODE_LIST}" ; \
		export E2E_KIND_CLUSTER_NODE_LIST="$${NODE_LIST}" ; \
		export E2E_CLUSTER_NAME=$(KIND_CLUSTER_NAME) ; \
		export APISERVER=`kubectl config view  --kubeconfig $(KIND_KUBECONFIG) |grep server | awk '{print $$2}'` ;\
		export E2E_INSTALL_NAMESPACE=$(E2E_INSTALL_NAMESPACE) ;\
		export GIT_COMMIT_VERSION=$(GIT_COMMIT_VERSION) ;\
		export APP_CHART_DIR=$(APP_CHART_DIR) ;\
		export E2E_TEST_SERVER_NAME=$(E2E_TEST_SERVER_NAME) ;\
		export E2E_KIND_KUBECONFIG_PATH=$(E2E_KIND_KUBECONFIG_PATH) ;\
		export AGENT_REGISTER=$(REGISTER)/$(GIT_REPO)-agent ; \
		if [ "$(E2E_IP_FAMILY)" == "ipv4" ] ; then \
			export E2E_IPV4_ENABLED=true ; export E2E_IPV6_ENABLED=false ; \
		elif [ "$(E2E_IP_FAMILY)" == "ipv6" ] ; then \
			export E2E_IPV4_ENABLED=false ; export E2E_IPV6_ENABLED=true ; \
		else \
			export E2E_IPV4_ENABLED=true ; export E2E_IPV6_ENABLED=true ; \
		fi ; \
		export E2E_KUBECONFIG_PATH=$(KIND_KUBECONFIG) ; [ -f "$(KIND_KUBECONFIG)" ] || { echo "error, does not exist KUBECONFIG $(E2E_KUBECONFIG)" ; exit 1 ; } ; \
		rm -f $(E2E_LOG_FILE) || true ; \
		echo "=========== before test `date` ===========" >> $(E2E_LOG_FILE) ; \
		./scripts/debugCluster.sh $(KIND_KUBECONFIG) "system" "$(E2E_LOG_FILE)"  $(E2E_INSTALL_NAMESPACE) ; \
		RESULT=0 ; \
		$(ROOT_DIR)/tools/golang/ginkgo.sh \
			--race --timeout=$(E2E_TIMEOUT) --output-interceptor-mode=none  --slow-spec-threshold=15s \
			--json-report e2ereport.json --output-dir $(ROOT_DIR) --procs $(E2E_GINKGO_PROCS) \
			--label-filter="$(E2E_GINKGO_LABELS)" -randomize-suites -randomize-all  -vv --fail-fast  $(E2E_GINKGO_OPTION) \
			-r e2e/*  || RESULT=1  ; \
		echo "=========== after test `date` ===========" >> $(E2E_LOG_FILE) ; \
		./scripts/debugCluster.sh $(KIND_KUBECONFIG) "system" "$(E2E_LOG_FILE)"  $(E2E_INSTALL_NAMESPACE) ; \
		KIND_CLUSTER_NAME=$(KIND_CLUSTER_NAME) ./scripts/debugCluster.sh $(KIND_KUBECONFIG) "detail" "$(E2E_LOG_FILE)"  $(E2E_INSTALL_NAMESPACE) ; \
		./scripts/debugCluster.sh $(KIND_KUBECONFIG) "error" "$(E2E_LOG_FILE)"  $(E2E_INSTALL_NAMESPACE) || { echo "error, found error log, datarace/pacni/longlock !!!" ; RESULT=1 ; } ; \
		if (($${RESULT} != 0)) ; then \
		   echo "failed to run e2e test"  ; \
		   exit 1 ; \
		fi ; \
		echo "" ; \
		echo "============================================" ; \
		echo "succeeded to run all test" ; \
		echo "output report to e2ereport.json" ; \
		echo "output env log to $(E2E_LOG_FILE) "


.PHONY: ci_threshold
ci_threshold: KIND_CLUSTER_NAME ?= $(E2E_KIND_CLUSTER_NAME)
ci_threshold: KIND_KUBECONFIG ?= $(E2E_KIND_KUBECONFIG_PATH)
ci_threshold:
	@echo -e "\033[35m Run ci threshold test on the cluster $(KIND_CLUSTER_NAME) \033[0m "
	@ echo -e "\033[35m [E2E] Run ci threshold with ginkgo label=$(E2E_GINKGO_LABELS) , timeout=$(E2E_TIMEOUT) GINKGO_OPTION=$(E2E_GINKGO_OPTION) \033[0m"
	@  NODE_LIST=` docker ps | egrep " kindest/node.* $(KIND_CLUSTER_NAME)-(control|worker)" | awk '{print $$NF }' ` ; \
		[ -n "$$NODE_LIST" ] || { echo "error, failed to find any kind nodes, please setup kind cluster $(KIND_CLUSTER_NAME) first" ; exit 1 ; } ; \
		NODE_LIST=` echo "$${NODE_LIST}" | tr -d ' ' | tr '\n' ',' ` ; \
		NODE_LIST=$${NODE_LIST%%,} ; \
		echo "find cluster node: $${NODE_LIST}" ; \
		export E2E_KIND_CLUSTER_NODE_LIST="$${NODE_LIST}" ; \
		export E2E_CLUSTER_NAME=$(KIND_CLUSTER_NAME) ; \
		export APISERVER=`kubectl config view  --kubeconfig $(KIND_KUBECONFIG) |grep server | awk '{print $$2}'` ;\
		export E2E_INSTALL_NAMESPACE=$(E2E_INSTALL_NAMESPACE) ;\
		export GIT_COMMIT_VERSION=$(GIT_COMMIT_VERSION) ;\
		export APP_CHART_DIR=$(APP_CHART_DIR) ;\
		export E2E_TEST_SERVER_NAME=$(E2E_TEST_SERVER_NAME) ;\
		export E2E_KIND_KUBECONFIG_PATH=$(E2E_KIND_KUBECONFIG_PATH) ;\
		export AGENT_REGISTER=$(REGISTER)/$(GIT_REPO)-agent ; \
		if [ "$(E2E_IP_FAMILY)" == "ipv4" ] ; then \
			export E2E_IPV4_ENABLED=true ; export E2E_IPV6_ENABLED=false ; \
		elif [ "$(E2E_IP_FAMILY)" == "ipv6" ] ; then \
			export E2E_IPV4_ENABLED=false ; export E2E_IPV6_ENABLED=true ; \
		else \
			export E2E_IPV4_ENABLED=true ; export E2E_IPV6_ENABLED=true ; \
		fi ; \
		export E2E_KUBECONFIG_PATH=$(KIND_KUBECONFIG) ; [ -f "$(KIND_KUBECONFIG)" ] || { echo "error, does not exist KUBECONFIG $(E2E_KUBECONFIG)" ; exit 1 ; } ; \
		rm -f $(E2E_LOG_FILE) || true ; \
		echo "=========== before test `date` ===========" >> $(E2E_LOG_FILE) ; \
		./scripts/debugCluster.sh $(KIND_KUBECONFIG) "system" "$(E2E_LOG_FILE)"  $(E2E_INSTALL_NAMESPACE) ; \
		RESULT=0 ; \
		$(ROOT_DIR)/tools/golang/ginkgo.sh \
			--race --timeout=$(E2E_TIMEOUT) --output-interceptor-mode=none  --slow-spec-threshold=15s \
			--json-report ci_threshold_report.json --output-dir $(ROOT_DIR) --procs 1 \
			--label-filter="$(E2E_GINKGO_LABELS)" -randomize-suites -randomize-all  -vv --fail-fast  $(E2E_GINKGO_OPTION) \
			-r ci-threshold/*  || RESULT=1  ; \
		echo "=========== after test `date` ===========" >> $(E2E_LOG_FILE) ; \
		./scripts/debugCluster.sh $(KIND_KUBECONFIG) "system" "$(E2E_LOG_FILE)"  $(E2E_INSTALL_NAMESPACE) ; \
		KIND_CLUSTER_NAME=$(KIND_CLUSTER_NAME) ./scripts/debugCluster.sh $(KIND_KUBECONFIG) "detail" "$(E2E_LOG_FILE)"  $(E2E_INSTALL_NAMESPACE) ; \
		./scripts/debugCluster.sh $(KIND_KUBECONFIG) "error" "$(E2E_LOG_FILE)"  $(E2E_INSTALL_NAMESPACE) || { echo "error, found error log, datarace/pacni/longlock !!!" ; RESULT=1 ; } ; \
		if (($${RESULT} != 0)) ; then \
		   echo "failed to run e2e test"  ; \
		   exit 1 ; \
		fi ; \
		echo "" ; \
		echo "============================================" ; \
		echo "succeeded to run all test" ; \
		echo "output report to e2ereport.json" ; \
		echo "output env log to $(E2E_LOG_FILE) "