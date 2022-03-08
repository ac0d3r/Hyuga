FROM node:lts-alpine as frontend-builder

WORKDIR /opt/frontend

COPY frontend .
RUN npm install -g npm@8.3.0
RUN npm install
RUN npm run build

FROM golang:1.16 as hyuga-builder

WORKDIR /opt/hyuga

ENV CGO_ENABLED='0'
ENV GO111MODULE='on'
ENV GOPROXY=https://goproxy.io

COPY hyuga .
RUN go mod tidy && go build -ldflags "-s -w" -o main main.go

FROM alpine 

WORKDIR /opt/hyuga
RUN mkdir -p /opt/hyuga/dist
COPY --from=frontend-builder /opt/frontend/dist /opt/hyuga/dist
COPY --from=hyuga-builder /opt/hyuga/config.yaml .
COPY --from=hyuga-builder /opt/hyuga/main .

EXPOSE 8000
EXPOSE 52/udp

RUN chmod +x /opt/hyuga/main

CMD ["/opt/hyuga/main"]