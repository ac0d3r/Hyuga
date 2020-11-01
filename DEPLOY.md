# Deploy

## 域名与公网 IP 准备

搭建并使用 `Hyuga`，准备两个域名，一个域名作为 NS 服务器域名(例：ns.cn)，一个用于记录域名(例：hyuga.io)，一个公网 IP(例：1.1.1.1 )。

注意：ns.cn 的域名提供商需要支持自定义 NS 记录, hyuga.io 则无要求。

1. 在 ns.cn 中设置两条 A 记录：
- ns1.ns.cn  A 记录指向  1.1.1.1        
- ns2.ns.cn  A 记录指向  1.1.1.1

修改 hyuga.io 的 NS 记录为 `ns1.ns.cn`, `ns2.ns.cn` 

## 搭建运行 `Hyuga`

### 修改环境变量以及配置

1. 环境变量写入项目根目录下的 `.env` 文件：

```bash
APP_ENV=production

REDIS_SERVER=redis
REDIS_PORT=6379

DOMAIN=hyuga.io # 记录域名
NS1_DOMAIN=ns1.app.io # NS域名
NS2_DOMAIN=ns2.app.io # NS域名
SERVER_IP=1.1.1.1 # 公网IP
```

2. 修改 [nginx-hyuga.conf](./deploy/nginx/nginx-hyuga.conf) 中的 `server_name`
```nginx
server {
    listen 80;
    server_name hyuga.io;
...
}
...
server {
    listen 80;
    server_name *.hyuga.io;
....
```

### 前端
#### 1. 修改 [config.js](./ui/src/utils/conf.js) API 接口

修改 `api.hyuga.io:5000` 为记录域名，例：
```JavaScript
// before: const apihost = "http://api.hyuga.io:5000";
const apihost = "http://api.hyuga.io;
```
#### 2. 构建前端文件

```bash
$ cd Hyuga/ui
$ yarn build
$ rm -r ../frontend
$ mv dist/ ../frontend
```



## Running with docker
```bash
$ cd Hyuga
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" hyuga.go # 编译
$ docker-compose build
$ docker-compose up -d
```