FROM golang:1.17.6-alpine AS BUILDER
COPY . /app
WORKDIR /app
RUN cd /app \
    go build -o phala-monitor