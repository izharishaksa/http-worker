FROM golang:1.18-alpine3.14 AS builder

WORKDIR /app
COPY . ./

RUN apk update && apk add --no-cache git openssh-client
RUN go mod download
RUN CGO_ENABLED=0 go build -o /usr/bin/http-worker ./main.go

FROM alpine

RUN apk add tzdata
RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime

WORKDIR ./app
COPY --from=builder /usr/bin/http-worker /usr/bin/http-worker

CMD ["http-worker"]
