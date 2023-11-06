# Contribution

***

## unitest

run the following command to check unitest

`make unitest_tests`

## setup cluster and run test

1. check required developing tools on you local host. If something missing, please run 'test/scripts/install-tools.sh' to install them

        # make -C test  checkBin
         pass   'kubectl' installed:   GitVersion:"v1.24.4"
         pass   'kind' installed:  kind version 0.19.0
         pass   'helm' installed:  Version:"v3.12.0"
         pass   'docker' installed:  Docker version 24.0.6, build ed223bc

2. run the e2e

        # make e2e

   if your run it for the first time, it will download some images, you could set the http proxy

        # ADDR=10.6.0.1
        # export https_proxy=http://${ADDR}:7890 http_proxy=http://${ADDR}:7890
        # make e2e

   run a specified case

        # make e2e -e E2E_GINKGO_LABELS="lable1,label2"

3. you could do it step by step with the follow

    if you are in China, it could add `-e E2E_CHINA_IMAGE_REGISTRY=true` to pull images from china image registry, add `-e HTTP_PROXY=http://${ADDR}` to get chart

    build the image

        # do some coding

        $ git add .
        $ git commit -s -m 'message'

        # !!! images is built by commit sha, so make sure the commit is submit locally
        $ make build_local_image

    setup the cluster

        # setup the kind cluster of dual-stack
        # !!! images is tested by commit sha, so make sure the commit is submit locally
        $ make e2e_init

    run the e2e test

        # run all e2e test on dual-stack cluster
        $ make e2e_run

        # run all e2e test on ipv4-only cluster
        $ make e2e_run -e E2E_IP_FAMILY=ipv4

        # run all e2e test on ipv6-only cluster
        $ make e2e_run -e E2E_IP_FAMILY=ipv6

        $ ls e2ereport.json

        $ make e2e_clean

5.clean `make e2e_clean`

***

## Submit Pull Request

A pull request will be checked by following workflow, which is required for merging.

### Action: your PR should be signed off

When you commit your modification, add `-s` in your commit command `git commit -s`

### Action: check yaml files

If this check fails, see the [yaml rule](https://yamllint.readthedocs.io/en/stable/rules.html).

Once the issue is fixed, it could be verified on your local host by command `make lint-yaml`.

Note: To ignore a yaml rule, you can add it into `.github/yamllint-conf.yml`.

### Action: check golang source code

It checks the following items against any updated golang file.

* Mod dependency updated, golangci-lint, gofmt updated, go vet, use internal lock pkg

* Comment `// TODO` should follow the format: `// TODO (AuthorName) ...`, which easy to trace the owner of the remaining job

* Unitest and upload coverage to codecov

* Each golang test file should mark ginkgo label

### Action: check licenses

Any golang or shell file should be licensed correctly.

### Action: check markdown file

### Action: lint yaml file

If it fails, see <https://yamllint.readthedocs.io/en/stable/rules.html> for reasons.

You can test it on your local machine with the command `make lint-yaml`.

### Action: lint chart

You can test it on your local machine with the command `make lint_chart_format`.

### Action: lint openapi.yaml

### Action: check code spell

Any code spell error of golang files will be checked.

You can check it on your local machine with the command `make lint-code-spell`.

It could be automatically fixed on your local machine with the command `make fix-code-spell`.

If you believe it can be ignored, edit `.github/codespell-ignorewords` and make sure all letters are lower-case.

## Changelog

How to automatically generate changelogs:

1. All PRs should be labeled with "pr/release/***" and can be merged.

2. When you add the label, the changelog will be created automatically.

   The changelog contents include:

   * New Features: it includes all PRs labeled with "pr/release/feature-new"

   * Changed Features: it includes all PRs labeled with "pr/release/feature-changed"

   * Fixes: it includes all PRs labeled with "pr/release/bug"

   * All historical commits within this version

3. The changelog will be attached to Github RELEASE and submitted to /changelogs of branch 'github_pages'.
