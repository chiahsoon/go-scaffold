FROM golang:latest
LABEL maintainer="Chiah Soon <chiahsoon18@gmail.com>"
COPY ./ ./app
WORKDIR ./app
RUN go get -d -v ./...

ENV GO111MODULE=on

# docker build -t chiahsoon/go_scaffold:0.0.1 .
# docker run -it chiahsoon/go_scaffold go run cmd/api/run.go