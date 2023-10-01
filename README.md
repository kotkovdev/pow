# Proof of Work concept

This project impelements proof or work concept for DDOS service protection and returns wind of wisdom quotes if request is allowed.

## Requirements
 - TCP server should be protected from DDOS attacks with the [Proof of Work](https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
 - The choice of the POW algorithm should be explained.
 - After Proof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
 - Docker file should be provided both for the server and for the client that solves the POW challenge


## Explanation
In this impelentation using challenge-response algorithm based on sha256 hashed puzzles. Client get part of soruce hash and target hash, so client must generate ending of soruce hash duringh sha256 of generated hash is not equal target hash.
Sha256 was selected because is it fast, secured and popular asynch hashing algorithm which has implementation on many programming languages.
[Clients puzzles protocol](https://en.wikipedia.org/wiki/Client_Puzzle_Protocol) was selected because this is simple and agile controlled workload client task.
For example while big attack we can change complexity of task to x, so algorithm compexity will be O(n<sup>x</sup>)

## Running
```
docker-compsoe up -d server # starts server
docker-compose up client # starts container with client and makes reqeust
```

## Evoltion
In the next evolution step it might be [Guided tour puzzle protocol](https://en.wikipedia.org/wiki/Guided_tour_puzzle_protocol) but it's so over engineered for this solution.