FROM golang:1.17

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
	GOPROXY="https://goproxy.cn,direct"

WORKDIR /B2C/LeastMall

ADD . .

RUN go build -o app .

EXPOSE 8080

CMD ["./app"]