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

2. Start the Docker services.
    ```shell
    make up
    ```

3. Start the server by running the following command:
    ```shell
    go run cmd/main.go
    ```