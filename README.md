# go-kafka-example

## Install

### Install Kafka

- with Docker:
```shell
docker-compose up -d
```

OR

- with Unix:
1. click [Kafka](https://kafka.apache.org/downloads) Download kafka
2. unzip kafka or tar kafka
    ```shell
    mkdir /home/kafka
    tar -xzvf ./kafka_<vsersion>.tgz -C /home/kafka
    ```
3. cd `kafka`
    ```shell
    vi /home/kafka/kafka_<vsersion>/config/kraft/server.properties
    ```
4. edit `server.properties`
    ```shell
    node.id=1
    process.roles=broker,controller
    listeners=PLAINTEXT://zimug1:9092,CONTROLLER://zimug1:9093
    advertised.listeners = PLAINTEXT://:9092
    controller.quorum.voters=1@zimug1:9093,2@zimug2:9093,3@zimug3:9093
    log.dirs=/home/kafka/data/kafka3
    ```
5. Formatted storage directory
    ```shell
    /home/kafka/kafka_<vsersion>/bin/kafka-storage.sh random-uuid
    ```
6. Confirm the configuration and path
    ```shell
    /home/kafka/kafka_<vsersion>/bin/kafka-storage.sh format \
    -t SzIhECn-QbCLzIuNxk1A2A \
    -c /home/kafka/kafka_<vsersion>/config/kraft/server.properties
    ```
7. Confirm configuration
    ```shell
   cat /home/kafka/data/kafka3/meta.properties
    ```
8. Start Kafka
    ```shell
    nohup /home/kafka/kafka_<vsersion>/bin/kafka-server-start.sh /home/kafka/kafka_<vsersion>/config/kraft/server.properties 1>/dev/null 2>&1 &
    ```
   
### Install Go
Mac:
```shell
brew install go
```

Linux
- Ubuntu:
```shell
sudo apt-get install golang
```

- CentOS:
```shell
sudo yum install golang
```

Windows:
```shell
choco install go
```

### Install Go Kafka

cd `project`
```shell
go mod tidy
```

## Use

### Write Kafka Message
```shell
go run write_message.go
```

### Read Kafka Message
```shell
go run read_message.go
```
