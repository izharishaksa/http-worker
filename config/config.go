package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	workerName      string
	instanceId      string
	groupId         string
	topicName       string
	dlqTopicName    string
	maxRetryAttempt int
	brokerHost      string
	brokerPort      int
	invokingMethod  string
	invokingUrl     string
	elasticHost     string
	elasticUsername string
	elasticPassword string
	offsetTime      time.Time
}

func NewConfig(timeGenerator func() time.Time) (*Config, error) {
	workerName := os.Getenv("WORKER_NAME")
	instanceId := os.Getenv("INSTANCE_ID")
	groupId := os.Getenv("GROUP_ID")
	topicName := os.Getenv("TOPIC_NAME")
	dlqTopicName := os.Getenv("DLQ_TOPIC_NAME")
	maxRetryAttemptStr := os.Getenv("MAX_RETRY_ATTEMPT")
	brokerHostStr := os.Getenv("BROKER_HOST")
	brokerPortStr := os.Getenv("BROKER_PORT")
	invokingMethod := os.Getenv("INVOKING_METHOD")
	invokingUrl := os.Getenv("INVOKING_URL")
	offsetTimeStr := os.Getenv("OFFSET_TIME")
	elasticHost := os.Getenv("ELASTIC_HOST")
	elasticUsername := os.Getenv("ELASTIC_USERNAME")
	elasticPassword := os.Getenv("ELASTIC_PASSWORD")

	if workerName == "" || instanceId == "" || groupId == "" || topicName == "" || dlqTopicName == "" || maxRetryAttemptStr == "" ||
		brokerHostStr == "" || brokerPortStr == "" || invokingMethod == "" || invokingUrl == "" {
		return nil, errors.New("missing required environment variable")
	}

	maxRetryAttempt, err := strconv.Atoi(maxRetryAttemptStr)
	if err != nil {
		return nil, errors.New("invalid value for MAX_RETRY_ATTEMPT, must be an integer")
	}

	brokerPort, err := strconv.Atoi(brokerPortStr)
	if err != nil {
		return nil, errors.New("invalid value for BROKER_PORT, must be an integer")
	}

	offsetTime := timeGenerator()
	if offsetTimeStr != "" {
		offsetTime, err = time.Parse("2006-01-02:15:04:05", offsetTimeStr)
		if err != nil {
			return nil, errors.New("invalid value for OFFSET_TIME, must be in format YYYY-MM-DD:HH:MM:SS")
		}
	}

	return &Config{
		workerName:      workerName,
		instanceId:      instanceId,
		groupId:         groupId,
		topicName:       topicName,
		dlqTopicName:    dlqTopicName,
		maxRetryAttempt: maxRetryAttempt,
		brokerHost:      brokerHostStr,
		brokerPort:      brokerPort,
		invokingMethod:  invokingMethod,
		invokingUrl:     invokingUrl,
		offsetTime:      offsetTime,
		elasticHost:     elasticHost,
		elasticUsername: elasticUsername,
		elasticPassword: elasticPassword,
	}, nil
}

func (c Config) WorkerName() string {
	return c.workerName
}

func (c Config) InstanceId() string {
	return c.instanceId
}

func (c Config) GroupId() string {
	return c.groupId
}

func (c Config) TopicName() string {
	return c.topicName
}

func (c Config) DlqTopicName() string {
	return c.dlqTopicName
}

func (c Config) MaxRetryAttempt() int {
	return c.maxRetryAttempt
}

func (c Config) BrokerHost() string {
	return c.brokerHost
}

func (c Config) BrokerPort() int {
	return c.brokerPort
}

func (c Config) InvokingMethod() string {
	return c.invokingMethod
}

func (c Config) InvokingUrl() string {
	return c.invokingUrl
}

func (c Config) OffsetTime() time.Time {
	return c.offsetTime
}

func (c Config) ElasticHost() string {
	return c.elasticHost
}

func (c Config) ElasticUsername() string {
	return c.elasticUsername
}

func (c Config) ElasticPassword() string {
	return c.elasticPassword
}
