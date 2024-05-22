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
    redis-cli --cluster create localhost:7001 localhost:7002 localhost:7003 localhost:7004 localhost:7005 localhost:7006 --cluster-replicas 1
    ```

4. Start the server by running the following command:
    ```shell
    go run cmd/main.go
    ```

## Redis Cluster

Redis scales horizontally with a deployment topology called Redis Cluster. Redis Cluster provides a way to run a Redis
installation where data is automatically sharded across multiple Redis nodes. Redis Cluster also provides some degree of
availability during partitions—in practical terms, the ability to continue operations when some nodes fail or are unable
to communicate.

### Sharding strategies

Redis uses a hash-based sharding strategy to distribute keys across multiple nodes. This is implemented using the
following strategies:

1. Hash slot sharding:

    - Redis cluster uses a fixed set of 16384 hash slots.
    - Keys are mapped to hash slots using the CRC16 algorithm.
    - Each node in the cluster is responsible for a subset of these hash slots.
    - Example: Key foo hashes to slot 12182, so it will be stored in the node responsible for that slot.

2. Key tagging:

    - Key tags are used to ensure related keys are stored in the same hash slot.
    - Example: {user1}:name and {user1}:email will be hashed to the same slot because of the {user1} tag.

### Resharding process

Resharding basically means to move hash slots from a set of nodes to another set of nodes. Like cluster creation, it is
accomplished using the redis-cli utility.

1. Connect to any node in the cluster

    ```shell
    docker compose exec -it redis-node-1 sh
    redis-cli -c -h localhost -p <node-port>
    ```

2. Start the resharding process

   Enter the Redis CLI in cluster mode and use the CLUSTER command to reshard.

    ```shell
    redis-cli --cluster reshard localhost:<any-node-port>
    ```

3. Enter the total number of slots to move:

   When prompted, specify the number of hash slots you want to move.

    ```text
    How many slots do you want to move (from 1 to 16384)? <number-of-slots>
    ```

4. Specify the target node ID:

   Provide the ID of the node to which the slots should be moved. You can get the node IDs using:

    ```shell
    redis-cli -c -h localhost -p <node-port> cluster nodes
    ```

   Then, enter the target node ID when prompted.

    ```text
    What is the receiving node ID? <target-node-id>
    ```

5. Specify source node IDs:

   You can specify the source node IDs manually or let the system pick the source nodes automatically.

    ```text
    Source node #1: <source-node-id>
    ```

6. Confirm resharding:

    ```text
    Do you want to proceed with the proposed reshard plan (yes/no)? yes
    ```

   The resharding will then begin, and you will see output indicating the progress of slot migration.