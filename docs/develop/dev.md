# develop

## local develop

1. `make build_local_image`

2. `make e2e_init`

3. `make e2e_run`

4. check proscope, browser vists http://NodeIP:4040

5. apply cr

        cat <<EOF > mybook.yaml
        apiVersion: kdoctor.io/v1beta1
        kind: Mybook
        metadata:
          name: test
        spec:
          ipVersion: 4
          subnet: "1.0.0.0/8"
        EOF
        kubectl apply -f mybook.yaml

## chart develop

helm repo add rock https://kdoctor-io.github.io/kdoctor/
