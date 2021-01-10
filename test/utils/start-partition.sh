#!/bin/bash
# ip l s eth0 down
echo "start partition"
iptables -A INPUT -p tcp --destination-port 7000 -j DROP
iptables -A INPUT -p tcp --destination-port 7001 -j DROP