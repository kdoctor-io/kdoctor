include Makefile.defs

.PHONY: all
all:


# ========================== build image


define BUILD_BIN
for SUBCMD_BIN_DIR in $(CMD_BIN_DIR); do  \
	BIN_NAME=`basename $${SUBCMD_BIN_DIR}` ; \
    echo "begin to build $${BIN_NAME} under $${SUBCMD_BIN_DIR}" ; \
    mkdir -p $(DESTDIR_BIN) ; \
	rm -f $(DESTDIR_BIN)/$${BIN_NAME} ; \
	$(GO_BUILD) -o $(DESTDIR_BIN)/$${BIN_NAME}  $${SUBCMD_BIN_DIR}/main.go ; \
	(($$?!=0)) && echo "error, failed to build $${BIN_NAME}" && exit 1 ; \
	echo "succeeded to build '$${BIN_NAME}' to $(DESTDIR_BIN)/$${BIN_NAME}" ; \
done
endef

.PHONY: build_all_bin
build_all_bin:
	make build_controller_bin
	make build_agent_bin


.PHONY: build_controller_bin
build_controller_bin: CMD_BIN_DIR := $(ROOT_DIR)/cmd/controller
build_controller_bin:
	$(BUILD_BIN)

.PHONY: build_agent_bin
build_agent_bin: CMD_BIN_DIR := $(ROOT_DIR)/cmd/agent
build_agent_bin:
	$(BUILD_BIN)

# ------------

define BUILD_FINAL_IMAGE
echo "Build Image $(IMAGE_NAME):$(IMAGE_TAG)" ; \
		sed -i '2 a \ARG TARGETPLATFORM' $(DOCKERFILE_PATH) ; \
		sed -i '2 a \ARG BUILDPLATFORM' $(DOCKERFILE_PATH) ; \
		docker build  \
				--build-arg RACE=1 \
				--build-arg NOSTRIP=1 \
				--build-arg NOOPT=1 \
				--build-arg GIT_COMMIT_VERSION=$(GIT_COMMIT_VERSION) \
				--build-arg GIT_COMMIT_TIME=$(GIT_COMMIT_TIME) \
				--build-arg VERSION=$(GIT_COMMIT_VERSION) \
				--build-arg BUILDPLATFORM="linux/$(TARGETARCH)" \
				--build-arg TARGETPLATFORM="linux/$(TARGETARCH)" \
				--build-arg TARGETARCH=$(TARGETARCH) \
				--build-arg TARGETOS=linux \
				--file $(DOCKERFILE_PATH) \
				--tag ${IMAGE_NAME}:$(IMAGE_TAG) . ; \
		echo "build success for ${IMAGE_NAME}:$(IMAGE_TAG) " ; \
		sed -i '3 d' $(DOCKERFILE_PATH) ; \
		sed -i '3 d' $(DOCKERFILE_PATH)
endef


.PHONY: build_local_image
build_local_image: build_local_agent_image build_local_controller_image

.PHONY: build_local_agent_image
build_local_agent_image: IMAGE_NAME := ${REGISTER}/${GIT_REPO}-agent
build_local_agent_image: DOCKERFILE_PATH := $(ROOT_DIR)/images/agent/Dockerfile
build_local_agent_image: IMAGE_TAG := $(GIT_COMMIT_VERSION)
build_local_agent_image:
	$(BUILD_FINAL_IMAGE)

.PHONY: build_local_controller_image
build_local_controller_image: IMAGE_NAME := ${REGISTER}/${GIT_REPO}-controller
build_local_controller_image: DOCKERFILE_PATH := $(ROOT_DIR)/images/controller/Dockerfile
build_local_controller_image: IMAGE_TAG := $(GIT_COMMIT_VERSION)
build_local_controller_image:
	$(BUILD_FINAL_IMAGE)

#---------

define BUILD_BASE_IMAGE
IMAGE_DIR=` dirname $(DOCKERFILE_PATH) ` \
		sed -i '2 a \ARG TARGETPLATFORM' $(DOCKERFILE_PATH) ; \
		sed -i '2 a \ARG BUILDPLATFORM' $(DOCKERFILE_PATH) ; \
		TAG=` git ls-tree --full-tree HEAD -- $${IMAGE_DIR} | awk '{ print $$3 }' ` ; \
		echo "Build base image $(BASE_IMAGE_NAME):$${TAG}" ; \
		docker build  \
				--build-arg USE_PROXY_SOURCE=true \
				--build-arg TARGETPLATFORM="linux/$(TARGETARCH)" \
				--build-arg BUILDPLATFORM="linux/$(TARGETARCH)" \
				--build-arg TARGETARCH=$(TARGETARCH) \
				--build-arg TARGETOS=linux \
				--file $(DOCKERFILE_PATH) \
				--output type=docker \
				--tag $(BASE_IMAGE_NAME):$${TAG}   $${IMAGE_DIR} ; \
		(($$?==0)) || { echo "error , failed to build base image" ; exit 1 ;} ; \
		echo "build success $(BASE_IMAGE_NAME):$${TAG} " ; \
		sed -i '3 d' $(DOCKERFILE_PATH) ; \
		sed -i '3 d' $(DOCKERFILE_PATH)
endef

.PHONY: build_local_base_image
build_local_base_image: build_local_agent_base_image build_local_controller_base_image

.PHONY: build_local_agent_base_image
build_local_agent_base_image: DOCKERFILE_PATH := $(ROOT_DIR)/images/agent-base/Dockerfile
build_local_agent_base_image: BASE_IMAGE_NAME := ${REGISTER}/${GIT_REPO}-agent-base
build_local_agent_base_image:
	$(BUILD_BASE_IMAGE)

.PHONY: build_local_controller_base_image
build_local_controller_base_image: DOCKERFILE_PATH := $(ROOT_DIR)/images/controller-base/Dockerfile
build_local_controller_base_image: BASE_IMAGE_NAME := ${REGISTER}/${GIT_REPO}-controller-base
build_local_controller_base_image:
	$(BUILD_BASE_IMAGE)


#================= update golang

## Update Go version for all the components
.PHONY: update_go_version
update_go_version: update_images_dockerfile_golang update_mod_golang update_workflow_golang


.PHONY: update_images_dockerfile_golang
update_images_dockerfile_golang:
	GO_VERSION=$(GO_VERSION) $(ROOT_DIR)/tools/images/update-golang-image.sh


# Update Go version for GitHub workflow
.PHONY: update_workflow_golang
update_workflow_golang:
	$(QUIET) for fl in $(shell find .github/workflows -name "*.yaml" -print) ; do \
  			sed -i 's/go-version: .*/go-version: ${GO_IMAGE_VERSION}/g' $$fl ; \
  			done
	@echo "Updated go version in GitHub Actions to $(GO_IMAGE_VERSION)"


# Update Go version in go.mod
.PHONY: update_mod_golang
update_mod_golang:
	$(QUIET) sed -i -E 's/^go .*/go '$(GO_MAJOR_AND_MINOR_VERSION)'/g' go.mod
	@echo "Updated go version in go.mod to $(GO_VERSION)"


.PHONY: update_gofmt
update_gofmt: ## Run gofmt on Go source files in the repository.
	$(QUIET)for pkg in $(GOFILES); do $(GO) fmt $$pkg; done


.PHONY: lint_code_spell
lint_code_spell:
	$(QUIET) if ! which codespell &> /dev/null ; then \
  				echo "try to install codespell" ; \
  				if ! pip3 install codespell ; then \
  					echo "error, miss tool codespell, install it: pip3 install codespell" ; \
  					exit 1 ; \
  				fi \
  			fi ;\
  			codespell --config .github/codespell-config

.PHONY: fix_code_spell
fix_code_spell:
	$(QUIET) if ! which codespell &> /dev/null ; then \
  				echo "try to install codespell" ; \
  				if ! pip3 install codespell ; then \
  					echo "error, miss tool codespell, install it: pip3 install codespell" ; \
  					exit 1 ;\
  				fi \
  			fi; \
  			codespell --config .github/codespell-config  --write-changes

#================== chart

.PHONY: chart_package
chart_package: lint_chart_format lint_chart_version
	-@rm -rf $(DESTDIR_CHART)
	-@mkdir -p $(DESTDIR_CHART)
	cd $(DESTDIR_CHART) ; \
   		echo "package chart " ; \
   		helm package  $(CHART_DIR) ; \
   		echo "package app chart " ; \
   		helm package  $(APP_CHART_DIR) ; \

.PHONY: update_chart_version
update_chart_version:
	VERSION=`cat VERSION | tr -d '\n' ` ; [ -n "$${VERSION}" ] || { echo "error, wrong version" ; exit 1 ; } ; \
		echo "update chart version to $${VERSION}" ; \
		CHART_VERSION=`echo $${VERSION} | tr -d 'v' ` ; \
		sed -E -i 's?^version: .*?version: '$${CHART_VERSION}'?g' $(CHART_DIR)/Chart.yaml &>/dev/null  ; \
		sed -E -i 's?^appVersion: .*?appVersion: "'$${CHART_VERSION}'"?g' $(CHART_DIR)/Chart.yaml &>/dev/null  ; \
		sed -E -i 's?^version: .*?version: '$${CHART_VERSION}'?g' $(APP_CHART_DIR)/Chart.yaml &>/dev/null  ; \
		sed -E -i 's?^appVersion: .*?appVersion: "'$${CHART_VERSION}'"?g' $(APP_CHART_DIR)/Chart.yaml &>/dev/null  ; \
   		echo "version of all chart is right"


.PHONY: lint_chart_format
lint_chart_format:
	mkdir -p $(DESTDIR_CHART) ; \
   			echo "check chart" ; \
   			helm lint --with-subcharts $(CHART_DIR) ;\
    		echo "check app chart" ; \
    		helm lint --with-subcharts $(APP_CHART_DIR)


.PHONY: lint_chart_version
lint_chart_version:
	VERSION=`cat VERSION | tr -d '\n' ` ; [ -n "$${VERSION}" ] || { echo "error, wrong version" ; exit 1 ; } ; \
		echo "check chart version $${VERSION}" ; \
		CHART_VERSION=`echo $${VERSION} | tr -d 'v' ` ; \
			grep -E "^version: $${CHART_VERSION}" $(CHART_DIR)/Chart.yaml &>/dev/null || { echo "error, wrong version in Chart.yaml" ; exit 1 ; } ; \
			grep -E "^appVersion: \"$${CHART_VERSION}\"" $(CHART_DIR)/Chart.yaml &>/dev/null || { echo "error, wrong appVersion in Chart.yaml" ; exit 1 ; } ; \
			grep -E "^version: $${CHART_VERSION}" $(APP_CHART_DIR)/Chart.yaml &>/dev/null || { echo "error, wrong version in Chart.yaml" ; exit 1 ; } ; \
			grep -E "^appVersion: \"$${CHART_VERSION}\"" $(APP_CHART_DIR)/Chart.yaml &>/dev/null || { echo "error, wrong appVersion in Chart.yaml" ; exit 1 ; } ; \
   		echo "version of all chart is right"


.PHONY: lint_chart_trivy
lint_chart_trivy:
	@ docker run --rm \
 		  -v /tmp/trivy:/root/trivy.cache/  \
          -v $(ROOT_DIR):/tmp/src  \
          aquasec/trivy:0.36.0 config --exit-code 1  --severity $(LINT_TRIVY_SEVERITY_LEVEL) /tmp/src/charts  ; \
      (($$?==0)) || { echo "error, failed to check chart trivy" && exit 1 ; } ; \
      echo "chart trivy check: pass"


.PHONY: update_crd_sdk
update_crd_sdk:
	@ echo "update crd manifest" && ./tools/golang/crdControllerGen.sh
	@ echo "update crd sdk" && ./tools/golang/crdSdkGen.sh


.PHONY: validate_crd_sdk
validate_crd_sdk:
	@ echo "validate crd manifest"
	make update_crd_sdk ; \
		if ! test -z "$$(git status --porcelain)"; then \
  			echo "please run 'make update_crd_sdk' to update crd code" ; \
  			exit 1 ; \
  		fi ; echo "succeed to check crd sdk"


#=============== lint

.PHONY: lint_golang_everything
lint_golang_everything: lint_golang_lock lint_test_label lint_golang_format


define lint_go_format
	data=` find . ! \( -path './vendor' -prune \) ! \( -path './_build' -prune \) ! \( -path './.git' -prune \) ! \( -path '*.validate.go' -prune \) \
        -type f -name '*.go' | xargs gofmt -d -l -s ` ; \
	if [ -n "$${data}" ]; then \
		echo "Unformatted Go source code:" ;\
		echo "$${data}" ;\
		exit 1 ; \
	fi ; \
	echo "format of Go source code is right"
endef

.PHONY: lint_golang_format
lint_golang_format:
	@ $(lint_go_format)
	$(QUIET) $(GO_VET)  ./...
	$(QUIET) golangci-lint run
	export GOPROXY="https://goproxy.io|https://goproxy.cn|direct"  ; go mod tidy ; go mod vendor ; \
		if ! test -z "$$(git status --porcelain)"; then \
  			echo "please run 'go mod tidy && go mod vendor', and submit your changes" ; \
  			exit 1 ; \
  		fi ; echo "succeed to check golang vendor"

.PHONY: lint_golang_lock
lint_golang_lock:
	@ BAD="" ; \
 	 for l in sync.Mutex sync.RWMutex; do \
  		DATA=` grep -r --exclude-dir={.git,_build,vendor,externalversions,lock,contrib} -i --include \*.go "$${l}" . ` || true ; \
	    if [ -n "$${DATA}" ] ; then \
	   		 echo "Found $${l} usage. Please use pkg/lock instead to improve deadlock detection"; \
	   		 echo "$${DATA}" ; \
	    	 BAD="true" ;\
	    fi ; \
	  done; \
	  if [ -n "$${BAD}" ] ; then \
	    exit 1  ; \
	  fi


# should label for each test file
.PHONY: lint_test_label
lint_test_label:
	@ALL_TEST_FILE=` find  ./  -name "*_test.go" -not -path "./vendor/*" ` ; FAIL="false" ; \
		for ITEM in $$ALL_TEST_FILE ; do \
			[[ "$$ITEM" == *_suite_test.go ]] && continue  ; \
			! grep 'Label(' $${ITEM} &>/dev/null && FAIL="true" && echo "error, miss Label in $${ITEM}" ; \
		done ; \
		[ "$$FAIL" == "true" ] && echo "error, label check fail" && exit 1 ; \
		echo "each test go file is labeled right"


.PHONY: lint_yaml
lint_yaml:
	@$(CONTAINER_ENGINE) container run --rm \
		--entrypoint sh -v $(ROOT_DIR):/data cytopia/yamllint \
		-c '/usr/bin/yamllint -c /data/.github/yamllint-conf.yml /data' ; \
		if (($$?==0)) ; then echo "congratulations ,all pass" ; else echo "error, pealse refer <https://yamllint.readthedocs.io/en/stable/rules.html> " ; fi


.PHONY: lint_dockerfile_trivy
lint_dockerfile_trivy:
	@ docker run --rm \
 		  -v /tmp/trivy:/root/trivy.cache/  \
          -v $(ROOT_DIR):/tmp/src  \
          aquasec/trivy:0.36.0 config --exit-code 1  --severity $(LINT_TRIVY_SEVERITY_LEVEL) /tmp/src/images  ; \
      (($$?==0)) || { echo "error, failed to check dockerfile trivy" && exit 1 ; } ; \
      echo "dockerfile trivy check: pass"


.PHONY: lint_image_trivy
lint_image_trivy: IMAGE_NAME ?=
lint_image_trivy:
	@ [ -n "$(IMAGE_NAME)" ] || { echo "error, please input IMAGE_NAME" && exit 1 ; }
	@ docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
 		  -v /tmp/trivy:/root/trivy.cache/  \
          aquasec/trivy:0.36.0 image --exit-code 1  --severity $(LINT_TRIVY_SEVERITY_LEVEL)  $(IMAGE_NAME) ; \
      (($$?==0)) || { echo "error, failed to check dockerfile trivy", $(IMAGE_NAME)  && exit 1 ; } ; \
      echo "trivy check: $(IMAGE_NAME) pass"



#=========== unit test

.PHONY: unitest_tests
unitest_tests: UNITEST_DIR := pkg cmd
unitest_tests:
	-@rm -rf $(UNITEST_OUTPUT)
	-@mkdir -p $(UNITEST_OUTPUT)
	@echo "run unitest tests"
	$(ROOT_DIR)/tools/golang/ginkgo.sh   \
		--cover --coverprofile=./coverage.out --covermode set  \
		--json-report unitestreport.json \
		-randomize-suites -randomize-all --keep-going  --timeout=1h  -p   --slow-spec-threshold=120s \
		-vv  -r   $(UNITEST_DIR) \
		&& mv ./coverage.out  $(UNITEST_OUTPUT)/coverage.out \
		&& mv ./unitestreport.json  $(UNITEST_OUTPUT)/unitestreport.json
	go tool cover -html=$(UNITEST_OUTPUT)/coverage.out -o $(UNITEST_OUTPUT)/coverage-all.html
	@ echo "output coverage to $(UNITEST_OUTPUT)/coverage.out "
	@ echo "output unitestreport to $(UNITEST_OUTPUT)/unitestreport.json "
	@ echo "output coverage-all.html to $(UNITEST_OUTPUT)/coverage-all.html "


# ================ e2e

.PHONY: e2e
e2e:
	make -C test check_images_ready
	make -C test e2e

.PHONY: e2e_init
e2e_init:
	make -C test check_images_ready
	make -C test init_kind_env
	make -C test deploy_multus
	make -C test deploy_metallb
	make -C test deploy_contour
	make -C test deploy_project
	make -C test install_apiserver_token




.PHONY: e2e_run
e2e_run:
	make -C test e2e_test

.PHONY: e2e_clean
e2e_clean:
	make -C test clean

.PHONY: clean_e2e_kdoctor
clean_e2e_kdoctor:
	make -C test uninstall_kdoctor
#============ doc

.PHONY: preview_doc
preview_doc: PROJECT_DOC_DIR := ${ROOT_DIR}/docs
preview_doc:
	-docker stop doc_previewer &>/dev/null
	-docker rm doc_previewer &>/dev/null
	@echo "set up preview http server  "
	@echo "you can visit the website on browser with url 'http://127.0.0.1:8000' "
	[ -f "docs/mkdocs.yml" ] || { echo "error, miss docs/mkdocs.yml "; exit 1 ; }
	docker run --rm  -p 8000:8000 --name doc_previewer -v $(PROJECT_DOC_DIR):/host/docs \
        --entrypoint sh \
        --stop-timeout 3 \
        --stop-signal "SIGKILL" \
        squidfunk/mkdocs-material:8.5.11  -c "cd /host ; cp docs/mkdocs.yml ./ ;  mkdocs serve -a 0.0.0.0:8000"
	#sleep 10 ; if curl 127.0.0.1:8000 &>/dev/null  ; then echo "succeeded to set up preview server" ; else echo "error, failed to set up preview server" ; docker stop doc_previewer ; exit 1 ; fi


.PHONY: build_doc
build_doc: PROJECT_DOC_DIR := ${ROOT_DIR}/docs
build_doc: OUTPUT_TAR := site.tar.gz
build_doc:
	-@rm -rf $(DOC_OUTPUT)
	-@mkdir -p $(DOC_OUTPUT)
	-docker stop doc_builder &>/dev/null
	-docker rm doc_builder &>/dev/null
	[ -f "docs/mkdocs.yml" ] || { echo "error, miss docs/mkdocs.yml "; exit 1 ; }
	-@ rm -f ./docs/$(OUTPUT_TAR)
	@echo "build doc html " ; \
		docker run --rm --name doc_builder  \
		-v ${PROJECT_DOC_DIR}:/host/docs \
        --entrypoint sh \
		squidfunk/mkdocs-material:8.5.11 -c "cd /host && cp ./docs/mkdocs.yml ./ && mkdocs build 2>&1 && cd site && tar -czvf site.tar.gz * && mv ${OUTPUT_TAR} ../docs/"	@ [ -f "$(PROJECT_DOC_DIR)/$(OUTPUT_TAR)" ] || { echo "failed to build site to $(PROJECT_DOC_DIR)/$(OUTPUT_TAR) " ; exit 1 ; }
	@ mv $(PROJECT_DOC_DIR)/$(OUTPUT_TAR) $(DOC_OUTPUT)/$(OUTPUT_TAR)
	@ echo "succeeded to build site to $(DOC_OUTPUT)/$(OUTPUT_TAR) "


.PHONY: check_doc
check_doc: PROJECT_DOC_DIR := ${ROOT_DIR}/docs
check_doc: OUTPUT_TAR := site.tar.gz
check_doc:
	-docker stop doc_builder &>/dev/null
	-docker rm doc_builder &>/dev/null
	[ -f "docs/mkdocs.yml" ] || { echo "error, miss docs/mkdocs.yml "; exit 1 ; }
	-@ rm -f ./docs/$(OUTPUT_TAR)
	@echo "check doc" ; \
		MESSAGE=`docker run --rm --name doc_builder  \
		-v ${PROJECT_DOC_DIR}:/host/docs \
        --entrypoint sh \
        squidfunk/mkdocs-material:8.5.11 -c "cd /host && cp ./docs/mkdocs.yml ./ && mkdocs build 2>&1 && cd site && tar -czvf site.tar.gz * && mv ${OUTPUT_TAR} ../docs/" 2>&1` ; \
        if (( $$? !=0 )) ; then \
        	echo "!!! error, failed to build doc" ; \
        	exit 1 ; \
        fi ; \
        if grep -E "WARNING .* which is not found" <<< "$${MESSAGE}" ; then  \
        	echo "!!! error, some link is bad" ; \
        	exit 1 ; \
        fi
	@ [ -f "$(PROJECT_DOC_DIR)/$(OUTPUT_TAR)" ] || { echo "failed to build site to $(PROJECT_DOC_DIR)/$(OUTPUT_TAR) " ; exit 1 ; }
	-@ rm -f ./docs/$(OUTPUT_TAR)
	@ echo "all doc is ok "



# ==============================

.PHONY: validate_openapi_sdk
validate_openapi_sdk:
	tools/golang/goSwagger.sh validate ${ROOT_DIR}/api/v1/controllerServer
	tools/golang/goSwagger.sh validate ${ROOT_DIR}/api/v1/agentServer

.PHONY: update_openapi_sdk
update_openapi_sdk:
	tools/golang/goSwagger.sh clean ${ROOT_DIR}/api/v1/controllerServer
	tools/golang/goSwagger.sh generate ${ROOT_DIR}/api/v1/controllerServer
	tools/golang/goSwagger.sh clean ${ROOT_DIR}/api/v1/agentServer
	tools/golang/goSwagger.sh generate ${ROOT_DIR}/api/v1/agentServer


