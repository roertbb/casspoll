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



### What can or does go wrong

 - Too many requests to Cassandra
    - it does tend to fail at around 100 concurrent polls run by the bot
 - Wysyp API HTTP - error handling, zgubione głosy
 - Podział na partycje - głosy się nie zgadzają
 - Nie da się zagłosować bo jest po dueTime (przesunięcie tworzenia polla)

### Possible improvements

- Oher types of voting than majority of votes

### Takeaways

- Using Docker took longer than it should have, but we have learned a lot in the process
