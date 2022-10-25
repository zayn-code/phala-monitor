FROM golang:1.17.6-alpine AS BUILDER
RUN apk add --no-cache gcc g++ git openssh-client && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone
COPY . /app
WORKDIR /app
RUN cd /app && go env -w GOPROXY=https://goproxy.io,direct && go build -o phala-monitor