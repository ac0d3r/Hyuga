# Deploy

## 修改环境变量以及配置

- 环境变量写入项目根目录下的 `.env` 文件：

```
APP_ENV=[production/development]

REDIS_SERVER=[localhost/redis]
REDIS_PORT=[6379]

DOMAIN=[hyuga.io]
NS1_DOMAIN=[ns1.app.io]
NS2_DOMAIN=[ns2.app.io]
SERVER_IP=[1.1.1.1]
```

- 修改 [nginx-hyuga.conf](../deploy/nginx/nginx-hyuga.conf) 中的 `[hyuga.io]` 域名
```nginx
server {
    listen 80;
    server_name [hyuga.io];
...
}
server {
    listen 80;
    server_name *.[hyuga.io];
....
```

## Running with docker
```bash
$ cd Hyuga
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" hyuga.go # 编译
$ docker-compose build
$ docker-compose up -d
```