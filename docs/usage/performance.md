
# Kdoctor Performance

## Environment

- Kubenetes: `v1.28.2`
- Container runtime: `containerd 1.6.25`
- OS: `Ubuntu 23.04`
- Kernel: `6.2.0-36-generic.x86_64`
- CPU: `Intel(R) Xeon(R) CPU E5-2680 v4 @ 2.40GHz`

| Node    | Role          | CPU | Memory |
|---------| ------------- |-----|--------|
| master1 | control-plane | 56C | 128Gi  |
| worker1 |               | 56C | 128Gi  |

## Nethttp

The following resource use cases for the kdcotor agent and other pressure testing tools, 
including memory and cpu overhead, can be used as a reference for direct deployment

### Http1.1

| Client  | Time | CPU | QPS  | Memory |
|---------|------|-----|------|--------|
| kdoctor | 1m   | 1C  | 7570 | 50Mb   |
| ab      | 1m   | 1C  | 9045 | 15Mb   |
| wrk     | 1m   | 1C  | 8920 | 10Mb   |
| hey     | 1m   | 1C  | 4637 | 115Mb  |

| Client  | Time | CPU | QPS   | Memory |
|---------|------|-----|-------|--------|
| kdoctor | 1m   | 2C  | 18888 | 55Mb   |
| ab      | 1m   | 2C  | 21031 | 20Mb   |
| wrk     | 1m   | 2C  | 20860 | 20Mb   |
| hey     | 1m   | 2C  | 10774 | 140Mb  |

| Client  | Time | CPU | QPS   | Memory |
|---------|------|-----|-------|--------|
| kdoctor | 1m   | 3C  | 28879 | 60Mb   |
| ab      | 1m   | 3C  | 35310 | 30Mb   |
| wrk     | 1m   | 3C  | 34445 | 28Mb   |
| hey     | 1m   | 3C  | 17174 | 167Mb  |


### Http2

| Client  | Time | CPU | QPS  | Memory |
|---------|------|-----|------|--------|
| kdoctor | 1m   | 1C  | 9733 | 77Mb   |
| hey     | 1m   | 1C  | 6100 | 140Mb  |

| Client  | Time | CPU | QPS   | Memory |
|---------|------|-----|-------|--------|
| kdoctor | 1m   | 2C  | 20943 | 78Mb   |
| hey     | 1m   | 2C  | 12600 | 167Mb  |

| Client  | Time | CPU | QPS   | Memory |
|---------|------|-----|-------|--------|
| kdoctor | 1m   | 3C  | 31524 | 79Mb   |
| hey     | 1m   | 3C  | 15300 | 230Mb  |

## Netdns

Use two replicas of coredns as the test server.

| Client  | Time | CPU | QPS   | Memory |
|---------|------|-----|-------|--------|
| kdoctor | 1m   | 1C  | 10064 | 82Mb   |
| dnsperf | 1m   | 1C  | 16800 | 5Mb    |

| Client  | Time | CPU | QPS   | Memory |
|---------|------|-----|-------|--------|
| kdoctor | 1m   | 2C  | 21137 | 91Mb   |
| dnsperf | 1m   | 2C  | 28769 | 8Mb    |

| Client  | Time | CPU | QPS   | Memory |
|---------|------|-----|-------|--------|
| kdoctor | 1m   | 3C  | 29987 | 103Mb  |
| dnsperf | 1m   | 3C  | 35800 | 9Mb    |