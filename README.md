# Highload Social Messenger

## Description

This project is a part of [Highload Social](https://github.com/eydeveloper/highload-social). It extends the main project
with messenger functionality.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- You have installed Docker
- You have installed Go

## Setup

To set up this project, follow these steps:

1. Clone the repository:
    ```shell
    git clone https://github.com/eydeveloper/highload-social-messenger.git
    ```

2. Start the Docker services:
    ```shell
    make up
    ```

3. Configue Redis Cluster.
    ```shell
    docker compose exec -it redis-node-1 sh
    redis-cli --cluster create 172.22.0.4:7001 172.22.0.7:7002 172.22.0.3:7003 172.22.0.6:7004 172.22.0.5:7005 172.22.0.2:7006 --cluster-replicas 1
    
    ```

4. Start the server by running the following command:
    ```shell
    go run cmd/main.go
    ```