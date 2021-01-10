#!/bin/bash
# ip l s eth0 down

iptables -A INPUT -p tcp --destination-port 7000 -j DROP
iptables -A INPUT -p tcp --destination-port 7001 -j DROP