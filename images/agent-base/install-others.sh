#!/bin/bash

# Copyright 2023 Authors of kdoctor-io
# SPDX-License-Identifier: Apache-2.0

set -x

set -o xtrace
set -o errexit
set -o pipefail
set -o nounset

packages=(
  # Additional iproute2 runtime dependencies
  libelf1
  libmnl0
  #bash-completion
  iptables
  # ss
  iproute2
  # netstat
  net-tools
  arping
  iftop
  conntrack
  ipvsadm
  lsof
  iputils-ping
  iputils-tracepath
  tcpdump
  telnet
  # ssh / ssh-client
  curl
  #netcat
  # after ubuntu 22.10
  netcat-traditional
  socat
  # nping
  nmap
  ssmping
  ethtool
  jq
  stress-ng
  # too big
  #sysstat
  pciutils
  iperf3
  netperf
  dnsutils
  dnsperf
  nftables
  tcptraceroute
  tshark
  openssl
  wget
  ndisc6
)

TARGETARCH=${1:-""}
echo "TARGETARCH=$TARGETARCH"

export DEBIAN_FRONTEND=noninteractive
apt-get update
ln -fs /usr/share/zoneinfo/UTC /etc/localtime
apt-get install -y --no-install-recommends "${packages[@]}"
apt-get purge --auto-remove
apt-get clean
rm -rf /var/lib/apt/lists/*


WEBSOCAT_VERSION=v1.10.0
if [ "$TARGETARCH"x == "arm64"x ] ; then
    wget --no-check-certificate https://github.com/vi/websocat/releases/download/${WEBSOCAT_VERSION}/websocat.aarch64-unknown-linux-musl -O /usr/sbin/websocat
else
    wget --no-check-certificate https://github.com/vi/websocat/releases/download/${WEBSOCAT_VERSION}/websocat.x86_64-unknown-linux-musl -O /usr/sbin/websocat
fi
chmod +x /usr/sbin/websocat


#========= verify

# maybe fail to call on building machine
#iptables-legacy --version
#iptables-nft --version
#ip6tables-legacy --version
#ip6tables-nft --version
which iptables-legacy
which iptables-nft
which ip6tables-legacy
which ip6tables-nft

ss -v
tc -V
netstat --version
arping -h
iftop  -h
conntrack --version
# Can't initialize ipvs: Permission denied (you must be root)
which ipvsadm
which telnet
which nc
which ssmping
lsof -h
ping -V
ping6 -V
tracepath -V
tcpdump --version
curl -V
socat -V
nping -V
nmap -V
ethtool --version
jq -V
stress-ng -V
lspci --version
iperf3 -v
netperf -V
which tcptraceroute

#
#echo 'ENABLED="true"' > /etc/default/sysstat
#service sysstat restart
#sar -V
#

dig -v
dnsperf -h

exit 0
