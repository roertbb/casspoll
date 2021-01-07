FROM cassandra:3.11.8

USER root

RUN apt-get update
RUN apt-get install -y iputils-ping

COPY docker-entrypoint.sh /usr/local/bin/

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["cassandra", "-f", "-R"]