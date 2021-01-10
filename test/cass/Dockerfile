FROM cassandra:3.11.8

USER root

RUN apt-get update
RUN apt-get install -y iputils-ping
RUN apt-get install -y iptables

CMD ["cassandra", "-f", "-R"]