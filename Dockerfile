# syntax=docker/dockerfile:1

FROM golang:1.18-alpine
ENV GOOS=linux
RUN mkdir -p /Typing/src/Services/Inout
WORKDIR /src
RUN export GOPATH=/Typing/src/Services/Inout
ADD Services/Inout/ /Typing/src/Services/Inout
WORKDIR /Typing/src/Services/Inout
EXPOSE 5006
RUN go build -o main .
CMD ["/Typing/src/Services/Inout/main"]