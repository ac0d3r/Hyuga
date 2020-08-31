# Deploy
## 修改环境变量
- 修改 `.env` 文件或者写入环境变量

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
```
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

## 需要的服务 services
- redis
- nginx
- supervisor

`deploy/build.sh` 只简简单编写 centos 系统的脚本

### Nginx

1. `vim /etc/nginx/nginx.conf` 需要注释掉主配置文件中的内容：
```conf
# listen       80 default_server;
# listen       [::]:80 default_server;
# server_name  _;
# root         /usr/share/nginx/html;
```

2. 拷贝 `deploy/nginx-hyuga.conf` 配置文件到 `/etc/nginx/conf.d/` 目录下.
3. 将 `frontend/*` 下的前端文件拷贝到 `nginx` 配置文件 `root` 目录下。
