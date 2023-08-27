start-worker:
	@GO111MODULE=on \
  WORKER_NAME="worker-name" \
  INSTANCE_ID="alobhj" \
  GROUP_ID="group-id" \
  TOPIC_NAME="main-topic" \
  DLQ_TOPIC_NAME="dlq-topic" \
  MAX_RETRY_ATTEMPT="3" \
  BROKER_HOST="localhost" \
  BROKER_PORT="9092" \
  INVOKING_METHOD="POST" \
  INVOKING_URL="http://example.com:3001/endpoint" \
  OFFSET_TIME="2006-01-02:15:04:05" \
  go run main.go
