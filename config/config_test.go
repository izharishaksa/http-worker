package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

type ConfigSuite struct {
	suite.Suite
	timeGenerator func() time.Time
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}

func (suite *ConfigSuite) SetupTest() {
	os.Setenv("WORKER_NAME", "worker")
	os.Setenv("INSTANCE_ID", "instance")
	os.Setenv("GROUP_ID", "group")
	os.Setenv("TOPIC_NAME", "topic")
	os.Setenv("DLQ_TOPIC_NAME", "dlq_topic")
	os.Setenv("MAX_RETRY_ATTEMPT", "3")
	os.Setenv("BROKER_HOST", "localhost")
	os.Setenv("BROKER_PORT", "9092")
	os.Setenv("INVOKING_METHOD", "GET")
	os.Setenv("INVOKING_URL", "https://example.com")
	os.Setenv("OFFSET_TIME", "2023-01-01:00:00:00")
	os.Setenv("ELASTIC_HOST", "localhost")
	os.Setenv("ELASTIC_USERNAME", "elastic")
	os.Setenv("ELASTIC_PASSWORD", "password")

	suite.timeGenerator = time.Now
}

func (suite *ConfigSuite) TearDownTest() {
	os.Clearenv()
}

func (suite *ConfigSuite) TestSuccess() {
	suite.Run("Success With Offset Time", func() {
		config, err := NewConfig(suite.timeGenerator)
		assert.Nil(suite.T(), err)
		assert.NotNil(suite.T(), config)
		assert.Equal(suite.T(), "worker", config.WorkerName())
		assert.Equal(suite.T(), "instance", config.InstanceId())
		assert.Equal(suite.T(), "group", config.GroupId())
		assert.Equal(suite.T(), "topic", config.TopicName())
		assert.Equal(suite.T(), "dlq_topic", config.DlqTopicName())
		assert.Equal(suite.T(), 3, config.MaxRetryAttempt())
		assert.Equal(suite.T(), "localhost", config.BrokerHost())
		assert.Equal(suite.T(), 9092, config.BrokerPort())
		assert.Equal(suite.T(), "GET", config.InvokingMethod())
		assert.Equal(suite.T(), "https://example.com", config.InvokingUrl())
		assert.Equal(suite.T(), "localhost", config.ElasticHost())
		assert.Equal(suite.T(), "elastic", config.ElasticUsername())
		assert.Equal(suite.T(), "password", config.ElasticPassword())

		expectedTime, _ := time.Parse("2006-01-02:15:04:05", "2023-01-01:00:00:00")
		assert.Equal(suite.T(), expectedTime, config.OffsetTime())
	})
}

func (suite *ConfigSuite) TestMissingEnvVars() {
	os.Clearenv()

	conf, err := NewConfig(suite.timeGenerator)
	assert.Nil(suite.T(), conf)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "missing required environment variable", err.Error())
}

func (suite *ConfigSuite) TestNewConfigWithInvalidValues() {
	suite.Run("Invalid value for MAX_RETRY_ATTEMPT", func() {
		os.Setenv("MAX_RETRY_ATTEMPT", "invalid")
		conf, err := NewConfig(suite.timeGenerator)
		assert.Nil(suite.T(), conf)
		assert.NotNil(suite.T(), err)
		assert.Equal(suite.T(), "invalid value for MAX_RETRY_ATTEMPT, must be an integer", err.Error())
	})

	suite.Run("Invalid value for BROKER_PORT", func() {
		os.Setenv("MAX_RETRY_ATTEMPT", "5")
		os.Setenv("BROKER_PORT", "invalid")
		conf, err := NewConfig(suite.timeGenerator)
		assert.Nil(suite.T(), conf)
		assert.NotNil(suite.T(), err)
		assert.Equal(suite.T(), "invalid value for BROKER_PORT, must be an integer", err.Error())
	})

	suite.Run("Invalid value for OFFSET_TIME", func() {
		os.Setenv("BROKER_PORT", "9092")
		os.Setenv("OFFSET_TIME", "invalid")
		conf, err := NewConfig(suite.timeGenerator)
		assert.Nil(suite.T(), conf)
		assert.NotNil(suite.T(), err)
		assert.Equal(suite.T(), "invalid value for OFFSET_TIME, must be in format YYYY-MM-DD:HH:MM:SS", err.Error())
	})
}
