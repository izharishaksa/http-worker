# HTTP Worker with Kafka Integration

This is a Go application that serves as an HTTP worker with Kafka integration. It is designed to consume messages from a Kafka topic, process them as HTTP requests, and retry if necessary. Below is an overview of the key components and functionality of this application.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Configuration](#configuration)
- [Usage](#usage)
- [Components](#components)
  - [Kafka Integration](#kafka-integration)
  - [Message Processing](#message-processing)
- [Retry Mechanism](#retry-mechanism)

## Prerequisites

Before running the application, ensure you have the following prerequisites:

- Go programming environment set up.
- Kafka broker accessible and configured.
- A valid configuration file (check Configuration section).
- Dependencies installed using `go get`:

```shell
go get github.com/google/uuid
go get github.com/izharishaksa/http-worker/config
go get github.com/segmentio/kafka-go
```

## Configuration
The application uses a configuration file to set various parameters. Configuration is loaded using the config package. Ensure you have a valid configuration file with the following parameters:

- **BrokerHost**: Kafka broker host address.
- **BrokerPort**: Kafka broker port.
- **TopicName**: The Kafka topic from which messages are consumed.
- **DlqTopicName**: The Kafka topic for dead-letter queue (DLQ) messages.
- **GroupId**: Kafka consumer group ID (optional).
- **MaxRetryAttempt**: Maximum number of retry attempts for processing a message.
- **OffsetTime**: Offset time for Kafka consumer (used when GroupId is not specified).
- **InvokingMethod**: HTTP request method (e.g., `POST`).
- **InvokingUrl**: URL to which HTTP requests are sent.

## Usage
To run the application, execute the following steps:

Ensure the prerequisites are met, and the configuration file is in place.

Build and run the application using the Go tool:

```shell
go build
./http-worker
```

The application will start consuming messages from the configured Kafka topic and processing them as HTTP requests.

You can gracefully shut down the application by sending a `SIGINT`, `SIGTERM`, `SIGQUIT`, or `SIGHUP` signal.

## Components
### Kafka Integration
The application uses the `segmentio/kafka-go` library for Kafka integration. It has a Kafka writer for publishing messages to a topic and a Kafka reader for consuming messages from a topic.

### Message Processing
Each Kafka message is processed as follows:

1. An HTTP request is created based on the configuration settings.
2. The message payload is set as the request body.
3. Headers, including a unique request ID, are added to the HTTP request.
4. The request is sent to the specified URL.
5. The response status code is checked, and if it's not 200 OK, an error is logged.

### Retry Mechanism
The application implements a retry mechanism for failed message processing. If a message processing attempt fails, it retries up to the configured maximum retry attempts with an increasing delay between retries. If the maximum retries are exhausted, the message is published to the Dead-Letter Queue (DLQ) topic for further analysis.