
# environment
- Kubenetes: `v1.25.4`
- container runtime: `containerd 1.6.12`
- OS: `CentOS Linux 8`
- kernel: `4.18.0-348.7.1.el8_5.x86_64`

| Node     | Role          | CPU  | Memory |
| -------- | ------------- | ---- | ------ |
| master1  | control-plane | 4C   | 8Gi    |
| master2  | control-plane | 4C   | 8Gi    |
| master3  | control-plane | 4C   | 8Gi    |
| worker4  |               | 3C   | 8Gi    |
| worker5  |               | 3C   | 8Gi    |
| worker6  |               | 3C   | 8Gi    |
| worker7  |               | 3C   | 8Gi    |
| worker8  |               | 3C   | 8Gi    |
| worker9  |               | 3C   | 8Gi    |
| worker10 |               | 3C   | 8Gi    |

# Nethttp

In a pod with a CPU of 1 

The test server is a server that sleeps for one second and then returns
## Http1.1

| client       | time | requests | qps     | Memory |
|--------------|------|----------|---------|--------|
| kdoctor | 0.5m | 81912    | 2988.67 | 210Mb  |
| hey          | 0.5m | 58423    | 1947.42 | 210Mb  |

| client       | time | requests | qps     | Memory |
|--------------|------|----------|---------|--------|
| kdoctor | 1m   | 179634   | 2,993.9 | 210Mb  |
| hey          | 1m   | 118452   | 1974.2  | 220Mb  |

| client       | time | requests | qps     | Memory |
|--------------|------|----------|---------|--------|
| kdoctor | 5m   | 897979   | 2993.26 | 210Mb  |
| hey          | 5m   | 596077   | 1986.92 | 270Mb  |


## Http2

| client       | time | requests | qps     | Memory |
|--------------|------|----------|---------|--------|
| kdoctor | 0.5m | 238787   | 7959.57 | 350Mb  |
| hey          | 0.5m | 7213     | 240.44  | 110Mb  |

| client       | time | requests  | qps      | Memory |
|--------------|------|-----------|----------|--------|
| kdoctor | 1m   | 481070    | 8017.83  | 370Mb  |
| hey          | 1m   | 14665     | 244.42   | 120Mb  |

| client       | time | requests | qps      | Memory |
|--------------|------|----------|----------|--------|
| kdoctor | 5m   | 2419874  | 8066.25  | 390Mb  |
| hey          | 5m   | 74776    | 249.25   | 130Mb  |


# Netdns

In a pod with a CPU of 1

| client       | time | requests | qps        | Memory |
|--------------|------|----------|------------|--------|
| kdoctor | 1m   | 1855511  | 30925.18   | 23Mb   |
| dnsperf      | 1m   | 1728086  | 28800.406  | 8Mb    |

| client       | time | requests | qps      | Memory |
|--------------|------|----------|----------|--------|
| kdoctor | 5m   | 9171699  | 30572.33 | 100Mb  |
| dnsperf      | 5m   | 8811137  | 29370.34 | 8Mb    |

| client       | time | requests  | qps       | Memory |
|--------------|------|-----------|-----------|--------|
| kdoctor | 10m  | 18561282  | 30935.47  | 173Mb  |
| dnsperf      | 10m  | 17260779  | 28767.666 | 8Mb    |
