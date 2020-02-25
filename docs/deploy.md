# Deploy

## Settings

### 添加环境变量

- 写入`.env`文件

```bash
APP_ENV=[production/development/testing]
SECRET_KEY=[secret key]

REDIS_SERVER=redis
REDIS_PORT=6379

DOMAIN=[一个用于记录的域名]
NS1_DOMAIN=[NS 域名]
NS2_DOMAIN=[NS 域名]
SERVER_IP=[记录的域名服务器的公网 IP]
```

### 配置 Nginx

[hyuga.conf](../deploy/nginx/conf.d/hyuga.conf)

将 `hyuga.io` 和 `*.hyuga.io` 设置为你自己配置好的域名

#### 关于域名设置

> from https://github.com/BugScanTeam/DNSLog/blob/master/README.md#安装

搭建并使用 `Hyuga`，需要拥有两个域名，一个域名作为 `NS` 服务器域名(例: `ns.com` )，一个用于记录域名(例: `hyuga.io` )。还需要有一个公网 IP 地址(如：`1.1.1.1`)。

**注意：`hyuga.io` 的域名提供商需要支持自定义 NS 记录, `ns.com` 则无要求。**

1. 在 `ns.com` 中设置两条 A 记录：

```
ns1.ns.com  A 记录指向  1.1.1.1
ns2.ns.com  A 记录指向  1.1.1.1
```

2. 修改 `hyuga.io` 的 NS 记录为 1 中设定的两个域名

---

### Install thirdparty redisco

```
cd Hyuga/thirdparty
git clone https://github.com/chen2aaron/redisco.git
```

### Front end

> 还在持续更新 ing👆, react 只会 🤏

Front end >>>> [Hyuga-react-README](https://github.com/Buzz2d0/Hyuga-react)

## Running with docker

```bash
cd Hyuga
docker-compose build
docker-compose up -d
```
