FROM golang:1.21.0

RUN apt-get update
RUN apt-get install -y graphviz
RUN go install github.com/cosmtrek/air@latest