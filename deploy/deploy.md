# Deploy

## 修改环境变量

- 写入 `.env` 文件或者写入环境变量：

```
APP_ENV=[production/development]

REDIS_SERVER=[localhost]
REDIS_PORT=[6379]

DOMAIN=[hyuga.io]
NS1_DOMAIN=[ns1.app.io]
NS2_DOMAIN=[ns2.app.io]
SERVER_IP=[1.1.1.1]
```

- 修改 `nginx-hyuga.conf` 域名
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

## 编译 hyuga.go
> CentOS 编译命令如下：

```shell
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" hyuga.go
```

## 需要的服务 services
- redis
- nginx [nginx-hyuga.conf](./nginx-hyuga.conf)
- supervisor [supervisor-hyuga.ini](./supervisor-hyuga.ini)


## 部署
1. `bash deploy/build.sh` (简单编写 centos 系统上的部署脚本)

2. `vim /etc/nginx/nginx.conf` 需要注释掉主配置文件中的内容：

```nginx
...
# listen       80 default_server;
# listen       [::]:80 default_server;
# server_name  _;
# root         /usr/share/nginx/html;
...
```
