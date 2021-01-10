#!/bin/bash
echo "stop partition"
iptables -D INPUT -p tcp --destination-port 7000 -j DROP
iptables -D INPUT -p tcp --destination-port 7001 -j DROP