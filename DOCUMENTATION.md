# 🪐 Documentation 📝

`Author:` **Miguel Ángel González Martín**

## Requirements
In order to run the **lunar-backend-engineer-challenge** application, you need to install the Golang language, Docker, and, in my case, Orbstack on Mac (to deploy the necessary Docker containers in the app, such as Redis and MySQL database, because the application runs directly locally).

## Run app 🚀
A Makefile file has been created to facilitate the various services and tests. You can view the description of the different commands typing in the terminal
```bash
make help
```

#### 1.- Install dependencies
```bash
go mod download
```
#### 2.- Start services
```bash
make start
```
You can check if the services are started correctly
```bash
docker compose ps -a
```
#### 3.- Run different endpoints
I've implemented the following endpoints:
1. Store rocket messages `POST localhost:8080/rockets`
2. Get all rockets stored ordered by `type` `GET localhost:8080/rockets`
3. Get rocket info by id `GET localhost:8080/rockets/{rocket-id}`

Open three terminals:
1. Run main.go with `go run cmd/rockets/web/main.go`
2. Check rockets stored `curl localhost:8080/rockets` or `curl localhost:8080/rockets/{rocket-id}`
3. Launch test program `./exe/darwin_arm64/rockets launch "http://localhost:8080/rockets" --message-delay=100ms --concurrency-level=10`

#### 4. Run tests
I`ve implemented unitary and acceptance tests, you can run them typing:
1. Unit tests `make unit-tests`
2. Acceptance tests `make acceptance-tests`

Or run both tests with `make test`

## Design considerations 🧑‍💻
When analyzing the solution to be implemented, a hexagonal architecture was applied that promotes the use of DDD (Domain Driven Design) and, with it, SOLID principles and clean code techniques. This gives us significant advantages by focusing the application logic on the domain with low coupling between services, resulting in high testability, tolerance to change, and high code reuse.

We decided to apply the CQRS (Command Query Responsibility Segregation) architecture pattern to separate write operations (commands) from read operations (queries) and thus be able to optimize and scale each data model independently.

We opted for a synchronous implementation to reduce complexity and meet the challenge deadline.
A more appropriate alternative for high-volume systems would have been to implement it asynchronously using DE (domain events) and event consumers.
In this way, for each message we received, a domain event or message would be published to an SNS topic and stored in an SQS queue subscribed to that topic, pending asynchronous consumption by a message consumer that would update the rocket's status at all times.
This would improve performance, resilience, and fault tolerance.

Note: to speed up implementation, I have used a series of libraries and utilities that I'm using in personal projects located in `/pkg` folder