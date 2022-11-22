FROM golang:1.19 as builder

ADD . /app

WORKDIR /app

ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn

RUN go mod tidy && go build -o a-line-cli


FROM ubuntu:latest
RUN apt update &&\
    apt install -y curl git &&\
    curl -fsSL https://get.docker.com | bash -s docker --mirror Aliyun

COPY --from=builder /app/a-line-cli /usr/local/bin/

EXPOSE 8080

CMD ["/usr/local/bin/a-line-cli","daemon"]
