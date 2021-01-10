#!/bin/bash
ADDRESS=$CASSANDRA_SERVER KEYSPACE=casspoll /build/server &

# docker exec cass1 /bin/bash -c "/utils/start-server.sh"