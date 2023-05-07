<div align="center" >
    <img src="./docs/hyuga.png" width="200" alt="Hyuga" />
</div>
<p align="center">
    <a href="https://github.com/ac0d3r/Hyuga">
        <img alt="Hyuga" src="https://img.shields.io/badge/Hyuga-dev-yellow"/>
    </a>
    <img src="https://img.shields.io/badge/Language-Golang-blue" alt="Language" />
    <a href="https://github.com/ac0d3r/Hyuga/blob/master/LICENSE">
        <img alt="License" src="https://img.shields.io/github/license/ac0d3r/Hyuga"/>
    </a>
    <a href="https://github.com/ac0d3r/Hyuga/stargazers">
        <img alt="stars" src="https://img.shields.io/github/stars/ac0d3r/Hyuga"/>
    </a>
 </p>

Hyuga 是一个用来监控带外(Out-of-Band)流量的工具。🪤

## 🎉 项目简介

项目地址：https://github.com/ac0d3r/Hyuga

## 📷 预览
<img width="1511" alt="image" src="https://user-images.githubusercontent.com/26270009/236437388-a862b25f-e049-420d-aaf8-b06e4aa66ccc.png">

## 🎉 功能

### 🌀 oob
- dns
    - dns查询记录(query name, remote address)
    - 支持 dns-rebinding [#🔗](#-dns-rebinding)
- http 
    - http 请求记录(url, method, remote address, raw request)
- ldap & rmi
    - ldap&rmi 请求记录(protocol, remote address, path) 

### 🪃 实时推送 
- 通过 websocket 将结果推送到前端。
- 支持第三方推送到Bark、Lark、钉钉、飞书、Sever酱。
    - thx: https://github.com/moonD4rk/notifier

### 🔦 单文件部署
- github action 自动发布 [Releases](https://github.com/ac0d3r/Hyuga/releases)

### 🔐 支持HTTPS
1. [安装Caddy](https://caddyserver.com/docs/install)
2. 配置 `/etc/caddy/Caddyfile` & 重启 `systemctl restart caddy`
```caddyfile
// Example
zznq.hyuga.icu {
    reverse_proxy localhost:8080
}
:80 {
    reverse_proxy localhost:8080
}
```

### 🚀 查询 API
- `GET` - `https://{hyuga.io}/api/v2/record/all?token={token}&type={type}&filter={filter}`
    - `type`: 查询类型 `dns|http|ldap|rmi`
    - `token`: 域名 token
    - `filter`: 过滤字符
- 支持重置 API Token

    <img width="250" alt="image" src="https://user-images.githubusercontent.com/26270009/236441871-60e51cf3-e0dc-4786-a6a8-869655b31a07.png">


## 👀 其他

### 🪓 DNS Rebinding
假设DNS Rebinding的域名为 `r.b34s.hyuga.io`， 公网IP为 `2.3.3.3`，dns的配置如下图：

<img width="420" alt="image" src="https://user-images.githubusercontent.com/26270009/236439602-09e1222f-09b5-4cee-b10b-d8e23b384464.png">

那么查询 `r.b34s.hyuga.io` 时根据访问次数依次返回所设置的dns（无缓存时）：`2.3.3.3` -> `127.0.0.1` -> `2.3.3.3`...
