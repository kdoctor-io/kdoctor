# develop new crd

1. define CRD in pkg/k8s/apis/kdoctor.io/v1beta1/xx_types.go
   add role to pkg/k8s/apis/kdoctor.io/v1beta1/rbac.go

2. make update_openapi_sdk

3. add crd to MutatingWebhookConfiguration and ValidatingWebhookConfiguration in charts/templates/tls.yaml 

4. add your crd to charts/template/role.yaml

5. implement the interface pkg/pluginManager/types in pkg/plugins/xxxx
   register your interface in pkg/pluginManager/types/manager.go

the plugin manager will auto help plugins to finish following jobs:

1. schedule task and call plugin to implement each round task

2. collect all report and save to controller disc

3. summarize each round result and update to CRD
