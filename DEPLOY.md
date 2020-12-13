# Deploy

## 域名与公网 IP 准备

搭建并使用 `Hyuga`，准备两个域名，一个域名作为 NS 服务器域名(例：ns.cn)，一个用于记录域名(例：hyuga.io)，一个公网 IP(例：<server_ip> )。

注意：ns.cn 的域名提供商需要支持自定义 NS 记录, hyuga.io 则无要求。

1. 在 ns.cn 中设置两条 A 记录：
- ns1.ns.cn  A 记录指向  <server_ip>        
- ns2.ns.cn  A 记录指向  <server_ip>

修改 hyuga.io 的 NS 记录为 `ns1.ns.cn`, `ns2.ns.cn` 

具体可参见 [DNSLog搭建](http://lovebear.top/2020/12/13/DNSLog/)

## 搭建运行 `Hyuga`

### 环境准备
1. OS：Ubuntu，其他系统相似，只不过安装软件命令有区别
2. 安装Redis

```
apt install redis-server
service redis start
```

3. 安装nginx

```
apt install nginx
service nginx start
```

4. 安装 yarn

```
curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | sudo apt-key add -
echo "deb https://dl.yarnpkg.com/debian/ stable main" | sudo tee /etc/apt/sources.list.d/yarn.list
apt update
apt install yarn
```

5. 安装 node

```
curl -sL https://deb.nodesource.com/setup_lts.x | bash -
apt install -y nodejs
```

6. 安装cnpm (国内npm慢)

```
npm install -g cnpm --registry=https://registry.npm.taobao.org
```

7. 安装vue

```
cnpm install -g @vue/cli
```

8. 安装vue-cli-service

```
cnpm install @vue/cli-service -g
```

### 修改环境变量以及配置

1. 环境变量写入项目根目录下的 .env 文件：

```
APP_ENV=production

REDIS_SERVER=<redis_ip> # 本地就localhost
REDIS_PORT=6379

DOMAIN=<hyuga.io> # 记录域名
NS1_DOMAIN=ns1.<ns.cn> # NS域名
NS2_DOMAIN=ns2.<ns.cn> # NS域名
SERVER_IP=<server_ip> # 公网IP
```

2. 修改 deploy/nginx/nginx-hyuga.conf 中的 server_name, proxy_pass

```
server {
    listen 80;
    server_name <hyuga.io>;
...
}
...
server {
    listen 80;
    server_name *.<hyuga.io>;
....
	location / {
        proxy_set_header Host $http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_pass http://<hyuga.io>:5000;
    }
```

### 前端
1. 修改 ui/src/utils/conf.js API 接口
修改 api.<hyuga.io>:5000 为记录域名
2. 构建前端文件

```
$ cd Hyuga/ui
$ yarn build
$ rm -r ../frontend
$ mv dist/ ../frontend
```

### 部署
1. 前端
a. 配置Nginx

```
cp deploy/nginx/nginx-hyuga.conf  /etc/nginx/conf.d/
```

b. 配置前端项目

```
rm -r /var/www/*
cp Hyuga/frontend/* /var/www/ -R
```

c. 重启nginx

```
service nginx restart
```

2. 后端
a. 编译go代码

```
go build hyuga.go
```

b. 运行

```
nohup ./hyuga &
```

P.S. 云服务器记得配置安全策略，打开80，53，5000端口



## Running with docker

```bash
$ cd Hyuga
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" hyuga.go # 编译
$ docker-compose build
$ docker-compose up -d
```