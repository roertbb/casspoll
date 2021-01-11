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
partition start
Cannot vote in the poll after it's due time
Cannot vote in the poll after it's due time
Cannot vote in the poll after it's due time
Cannot vote in the poll after it's due time
Cannot vote in the poll after it's due time
Cannot vote in the poll after it's due time
Cannot vote in the poll after it's due time
Cannot vote in the poll after it's due time
Cannot vote in the poll after it's due time
Cannot vote in the poll after it's due time
Cannot vote in the poll after it's due time
Cannot vote in the poll after it's due time
Cannot vote in the poll after it's due time
Cannot vote in the poll after it's due time
Cannot vote in the poll after it's due time
Cannot vote in the poll after it's due time
finish voting
Cannot vote in the poll after it's due time
PollID: 508967d1-2552-46df-92d9-3e5cde6c18ac | AnswerID: bbadc9fd-4a0d-42b7-abff-dc1e61108bea | wrong answer difference is - 22 / expected - 24
PollID: 508967d1-2552-46df-92d9-3e5cde6c18ac | AnswerID: bbdf081c-d1f8-4986-aa17-00f35195f97d | wrong answer difference is - 17 / expected - 24
PollID: 508967d1-2552-46df-92d9-3e5cde6c18ac | AnswerID: 8fa7364d-0d00-4654-af43-9b49adc81af9 | wrong answer difference is - 19 / expected - 24
PollID: 508967d1-2552-46df-92d9-3e5cde6c18ac | AnswerID: b6372f3d-5628-4275-bc93-a7df52889d6b | wrong answer difference is - 22 / expected - 24
PollID: 06bd0364-896c-4ced-a99a-05b6107a23c8 | AnswerID: 1bdfba9a-849e-4035-a76d-3381c5ea4023 | wrong answer difference is - 3 / expected - 4
PollID: 06bd0364-896c-4ced-a99a-05b6107a23c8 | AnswerID: 5a34916c-0253-43c6-b45a-554fe944bacd | wrong answer difference is - 3 / expected - 5
PollID: 06bd0364-896c-4ced-a99a-05b6107a23c8 | wrong poll winner | is - 1bdfba9a-849e-4035-a76d-3381c5ea4023 / expected - 5a34916c-0253-43c6-b45a-554fe944bacd
PollID: 14cb502f-79f5-40c5-a824-a6c7d6cfd80d | AnswerID: e508a596-e293-43b3-868b-a3113083a05c | wrong answer difference is - 18 / expected - 28
PollID: 14cb502f-79f5-40c5-a824-a6c7d6cfd80d | AnswerID: fed5470e-8a42-4421-a4db-8becb0ca5149 | wrong answer difference is - 22 / expected - 28
PollID: 14cb502f-79f5-40c5-a824-a6c7d6cfd80d | AnswerID: 3d599f47-fc9d-4b73-9c87-c33ff232f16a | wrong answer difference is - 19 / expected - 28
PollID: 14cb502f-79f5-40c5-a824-a6c7d6cfd80d | AnswerID: 6f6b65db-9df4-41e2-bc56-bc7d5c0057e5 | wrong answer difference is - 16 / expected - 28
PollID: 14cb502f-79f5-40c5-a824-a6c7d6cfd80d | AnswerID: 9a2cdb3c-263d-4464-ba39-5e535d85a11f | wrong answer difference is - 19 / expected - 28
PollID: 14cb502f-79f5-40c5-a824-a6c7d6cfd80d | AnswerID: b470e21f-7f34-4b1c-bd7e-939b7a025844 | wrong answer difference is - 19 / expected - 28
PollID: 14cb502f-79f5-40c5-a824-a6c7d6cfd80d | AnswerID: b72006e7-eff2-4941-be86-d803ec55a115 | wrong answer difference is - 20 / expected - 28
PollID: 14cb502f-79f5-40c5-a824-a6c7d6cfd80d | AnswerID: c9d2cec4-cec8-42bd-888c-d17cdc565953 | wrong answer difference is - 19 / expected - 28
PollID: 14cb502f-79f5-40c5-a824-a6c7d6cfd80d | wrong poll winner | is - fed5470e-8a42-4421-a4db-8becb0ca5149 / expected - e508a596-e293-43b3-868b-a3113083a05c
PollID: 2f228f3a-ba5c-4ca9-a5df-4e1cee1f0002 | AnswerID: 82d3c09a-ef5a-457f-91b4-b15bd78ca626 | wrong answer difference is - 23 / expected - 27
PollID: 2f228f3a-ba5c-4ca9-a5df-4e1cee1f0002 | AnswerID: ba0ecb85-caa1-4e0c-9f5a-be6cb0b30840 | wrong answer difference is - 24 / expected - 27
PollID: 2f228f3a-ba5c-4ca9-a5df-4e1cee1f0002 | AnswerID: d4e31e87-1686-47f0-82ee-52fb99c26c0c | wrong answer difference is - 24 / expected - 27
PollID: 2f228f3a-ba5c-4ca9-a5df-4e1cee1f0002 | AnswerID: 1acc8c68-146b-4cf7-8b7d-f65cfad20bab | wrong answer difference is - 21 / expected - 27
PollID: 2f228f3a-ba5c-4ca9-a5df-4e1cee1f0002 | AnswerID: 6e2169e6-2561-4345-98b6-6462a4ba66dd | wrong answer difference is - 21 / expected - 27
PollID: 2f228f3a-ba5c-4ca9-a5df-4e1cee1f0002 | AnswerID: 828a7cf2-4223-40b0-b11f-95529594835e | wrong answer difference is - 22 / expected - 27
PollID: 2f228f3a-ba5c-4ca9-a5df-4e1cee1f0002 | wrong poll winner | is - ba0ecb85-caa1-4e0c-9f5a-be6cb0b30840 / expected - 82d3c09a-ef5a-457f-91b4-b15bd78ca626
PollID: b1d57480-d3f1-4ba9-9c22-340f721d27b7 | AnswerID: 57771a5e-5183-4e16-a744-375d5965fbd5 | wrong answer difference is - 5 / expected - 7
PollID: b1d57480-d3f1-4ba9-9c22-340f721d27b7 | AnswerID: c04aeb3f-efa9-47a3-9a7f-a20246b404fe | wrong answer difference is - 2 / expected - 4
PollID: b1d57480-d3f1-4ba9-9c22-340f721d27b7 | AnswerID: f455779a-ecde-4e6d-84d9-707580c0495a | wrong answer difference is - 6 / expected - 7
PollID: b1d57480-d3f1-4ba9-9c22-340f721d27b7 | wrong poll winner | is - f455779a-ecde-4e6d-84d9-707580c0495a / expected - 57771a5e-5183-4e16-a744-375d5965fbd5
PollID: 0888c09f-12ed-42e2-9d3a-50d15bc8c5ac | AnswerID: 0f7ecf10-5f89-45eb-ad6e-383c5e02bf88 | wrong answer difference is - 7 / expected - 8
PollID: 0888c09f-12ed-42e2-9d3a-50d15bc8c5ac | AnswerID: af24f1a6-08a9-4e3b-b662-18427ac12b14 | wrong answer difference is - 4 / expected - 5
PollID: 013601a6-0e48-484a-a247-bb9ae32fe1f2 | AnswerID: 4085cb00-ce3f-44af-9076-3811f6509770 | wrong answer difference is - 1 / expected - 2
PollID: 013601a6-0e48-484a-a247-bb9ae32fe1f2 | AnswerID: 5200a1f5-ffd1-4c58-a344-f154e2e5a846 | wrong answer difference is - 3 / expected - 4
PollID: 49185db7-21e2-4b4a-ac1f-98b36f46185d | AnswerID: 590349e8-68ac-40b1-aab8-19d78cb0525c | wrong answer difference is - 5 / expected - 8
PollID: 49185db7-21e2-4b4a-ac1f-98b36f46185d | AnswerID: 62a0e888-3dd5-4a60-a8b0-7c2b25201d32 | wrong answer difference is - 2 / expected - 4
PollID: 49185db7-21e2-4b4a-ac1f-98b36f46185d | AnswerID: 702e75d5-a1cf-4674-a8d3-5e960fc3a18a | wrong answer difference is - 3 / expected - 4
PollID: 49185db7-21e2-4b4a-ac1f-98b36f46185d | AnswerID: 2f4642bd-b37c-457d-bc8c-088d9166c870 | wrong answer difference is - 2 / expected - 4
PollID: 49185db7-21e2-4b4a-ac1f-98b36f46185d | AnswerID: 556ad992-a965-4207-a909-45b22bfd56dd | wrong answer difference is - 4 / expected - 5
PollID: fc1ad310-8da6-4bb4-afe9-beef39f88451 | AnswerID: 99bca9ee-f652-44d0-8653-1f17e7b8fc01 | wrong answer difference is - 26 / expected - 34
PollID: fc1ad310-8da6-4bb4-afe9-beef39f88451 | AnswerID: d62ad8f7-5cf8-4d2c-9734-c63db0484448 | wrong answer difference is - 25 / expected - 34
PollID: fc1ad310-8da6-4bb4-afe9-beef39f88451 | AnswerID: 483a2712-3b22-423a-be35-a0910ebe38ae | wrong answer difference is - 27 / expected - 34
PollID: fc1ad310-8da6-4bb4-afe9-beef39f88451 | AnswerID: 7af62dcd-9be8-4099-bd0e-3738f47e7fe6 | wrong answer difference is - 24 / expected - 34
PollID: fc1ad310-8da6-4bb4-afe9-beef39f88451 | AnswerID: 8098b4f5-44d3-4006-9eee-5eda3c115f73 | wrong answer difference is - 27 / expected - 34
PollID: fc1ad310-8da6-4bb4-afe9-beef39f88451 | AnswerID: 81e77ab8-f865-41e1-9d45-5aaa6171bbc9 | wrong answer difference is - 28 / expected - 34
PollID: fc1ad310-8da6-4bb4-afe9-beef39f88451 | AnswerID: 8c193f61-ee7b-4d34-a83a-95dc6f31ce24 | wrong answer difference is - 21 / expected - 34
PollID: fc1ad310-8da6-4bb4-afe9-beef39f88451 | wrong poll winner | is - 81e77ab8-f865-41e1-9d45-5aaa6171bbc9 / expected - 99bca9ee-f652-44d0-8653-1f17e7b8fc01
PollID: bc4cd7bf-636d-4497-805c-6ffc5cb06b0f | AnswerID: ae8a4c22-30e1-4637-a4e9-15e6ebb051ca | wrong answer difference is - 17 / expected - 19
PollID: bc4cd7bf-636d-4497-805c-6ffc5cb06b0f | AnswerID: 0ba56b58-b382-4c07-80af-1a616d54fd14 | wrong answer difference is - 17 / expected - 19
PollID: bc4cd7bf-636d-4497-805c-6ffc5cb06b0f | AnswerID: 0fe5db43-1216-4b31-99b0-2345be2ca8e0 | wrong answer difference is - 15 / expected - 19
PollID: bc4cd7bf-636d-4497-805c-6ffc5cb06b0f | AnswerID: 1a942f3a-bef1-49ef-afe2-09b0db1308e0 | wrong answer difference is - 15 / expected - 19
PollID: bc4cd7bf-636d-4497-805c-6ffc5cb06b0f | AnswerID: 1f2876d4-d438-402f-abe1-2e8cd5a6191e | wrong answer difference is - 14 / expected - 19
PollID: bc4cd7bf-636d-4497-805c-6ffc5cb06b0f | AnswerID: 2086a92d-8acb-4ded-b334-f79daf77df17 | wrong answer difference is - 15 / expected - 19
PollID: bc4cd7bf-636d-4497-805c-6ffc5cb06b0f | AnswerID: 5aeafb2b-ee99-444c-9af0-bb0705587a96 | wrong answer difference is - 15 / expected - 19
PollID: bc4cd7bf-636d-4497-805c-6ffc5cb06b0f | AnswerID: 968cec84-47d0-4b02-84a5-72aec827a519 | wrong answer difference is - 16 / expected - 19
after voting time finished
wrongAnswersCount 47
wrongPollWinner 5
after removing partition and syncing
wrongAnswersCount 0
wrongPollWinner 0
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
