# Deploy

## 域名与公网 IP 准备

搭建并使用 `Hyuga` 需要准备两个域名，一个域名作为 NS 服务器域名(例：ns.cn)，一个用于记录域名(例：hyuga.io)，一个公网 IP(例：1.1.1.1 )。

注意：ns.cn 的域名提供商需要支持自定义 NS 记录, hyuga.io 则无要求。

1. 在 ns.cn 中设置两条 A 记录：
- ns1.ns.cn  A 记录指向  1.1.1.1        
- ns2.ns.cn  A 记录指向  1.1.1.1

修改 hyuga.io 的 NS 记录为 `ns1.ns.cn`, `ns2.ns.cn` 

## 🐳 使用 Docker 部署

### 修改配置文件

1. 修改 `config.yml` 文件：

```yml
app:
  env: production # development/production
  recordExpirationDays: 7
redis: redis:6379
domain:
  main: hyuga.io  # 修改记录域名
  ns: [ns1.app.io, ns2.app.io]  # 修改NS域名
  ip: 127.0.0.1 # 修改公网IP
```

1. 替换 [nginx.conf](./ui/nginx.conf) 中的 `server_name`
```nginx
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
```

### 前端
修改 [conf.js](./ui/src/utils/conf.js) API 接口

修改 `api.<hyuga.io:5000>` 为记录域名，例：
```JavaScript
const apihost = "http://api.hyuga.io;
...
```

### 运行
```bash
$ cd Hyuga
$ docker-compose build
$ docker-compose up -d
```

## Example
- [lovebear - DNSLog搭建](http://lovebear.top/2020/12/13/DNSLog/)