# CassPoll

It's like StrawPoll, but with Cassandra

Written in Golang and bash

Robert Banaszak, Marcin Ławniczak 2020-2021

functionalities:

 - single-choice polls
 - multiple-choice polls

### Project structure

 - `/bin` - generated binaries of programs
 - `/cmd` - go utilities to run
    - `bot` - script utilizing the HTTP API to demonstrate failure
    - `console` - console client used for initial testing
    - `server` - HTTP API server, an interface to Cassandra
 - `poll` - models and functions used to interface with Cassandra
 - `repo` -  Definition of Poll Repo for Cassandra, allows for swappable backend in the future
 - `schema` - Cassandra table structure
 - `test` - Docker-related files
    - `cass` - Container Dockerfile
    - `utils` - Utilities to run the experiment - migrate, server and bot
    
### Running it

A sample setup using Docker is provided.
The docker-compose file defines two containers, each containing a Cassandra node.
```
$ cd test
$ docker-compose up
```
*Do not close this terminal window, so you can see the Cassandras inner workings.*


Once the containers are running, Cassandra needs to be seeded with a schema
```
$ docker exec cass1 /bin/bash "utils/migrate.sh"
```

The HTTP runs on Cassandra nodes. It needs to be compiled and started.

```
$ cd cmd/server
$ make build-docker
$ docker exec cass1 /bin/bash "utils/start-server.sh"
$ docker exec cass2 /bin/bash "utils/start-server.sh"
```

The bot can now be run.

```
$ cd cmd/bot
$ make build
$ cd ../../test/utils
$ ./run-bot.sh
```

The bot starts 10 polls, votes in them for 20 seconds.
In the middle of voting, after 10 seconds, a partition is started by
blocking ports used by Cassandra for internode communication (simulated network failure).
Results are retrieved and compared with what was sent.
Then the connection is restored, and after 30 seconds results are retrieved again,
this time correct - including all the votes. All above values are configurable in 
`cmd/bot/main.go` (lines 25-29) - it needs to be recompiled after the values are changed.

An example output is provided below. We have been able to overload the servers using
higher concurrent poll count.

```

```





### What can or does go wrong

 - Too many requests to Cassandra
    - it does tend to fail at around 25 concurrent polls run by the bot on laptop hardware
 - We have encountered a "lost vote" problem without the partition during testing.
   It was caused by improper error handling which caused HTTP API requests to fail.
   It was easy to fix due to Golang strict typing
 - Once the connection between two servers is severed during voting, the vote counts do not match on both servers
 - The bot creates polls with dueTime (expiry time) 2 minutes into the future. Then it starts sending votes.
   After 2 minutes it stops, but the votes at the end aren't counted, because there is a delay between the poll starting
   and the bot sending votes. The last votes arrive after dueTime.

### Possible improvements

- Oher types of voting than majority of votes

### Takeaways

- Using Docker took longer than it should have, but we have learned a lot in the process
