FROM golang:1.17-stretch as build
WORKDIR /go/src/tztask
COPY . ./
ENV GOPROXY="https://goproxy.cn"
RUN go mod download
RUN go build -o server
FROM alpine:3.16
COPY --from=build /go/src/tztask/server  /opt
COPY --from=build /go/src/tztask/config.yaml  /opt
COPY --from=build /go/src/tztask/jobs.yaml  /opt

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
ENV SERVICE_NAME="tztask_server"
WORKDIR /opt
CMD ["./server"]
