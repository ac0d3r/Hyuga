FROM golang:1.15.8-alpine as builder

WORKDIR /opt/hyuga

ENV GO111MODULE='on'
ENV GOPROXY=https://goproxy.io

COPY . .
RUN go build -ldflags "-s -w" -o main hyuga.go

#------------------
FROM alpine
WORKDIR /opt/hyuga

COPY --from=builder /opt/hyuga/config.yml .
COPY --from=builder /opt/hyuga/main .

EXPOSE 5000
EXPOSE 52/udp

RUN chmod +x /opt/hyuga/main

CMD ["/opt/hyuga/main"]