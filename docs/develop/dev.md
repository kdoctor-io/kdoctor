# develop

## local develop

1. `make build_local_image`

2. setup kind 

        make e2e_init
            -----------------------------------------------------------------------------------------------------
             succeeded to setup cluster spider
             you could use following command to access the cluster
                export KUBECONFIG=$(pwd)/test/runtime/kubeconfig_kdoctor.config
                kubectl get nodes
            -----------------------------------------------------------------------------------------------------

    for chian developer 

        make e2e_init -e E2E_CHINA_IMAGE_REGISTRY=true -e E2E_HELM_HTTP_PROXY=http://xxxx

3. `make e2e_run`

4. check proscope, browser vists http://NodeIP:4040

5. `make e2e_clean`

## chart develop

helm repo add rock https://kdoctor-io.github.io/kdoctor/
