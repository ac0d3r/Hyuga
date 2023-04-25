<div align="center" >
    <img src="./docs/hyuga.png" width="280" alt="Hyuga" />
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

⚡️Hyuga 是一个用来检测带外(Out-of-Band)流量的监控工具。

---
## 🎉 项目简介

DEMO 主页：TODO

项目地址：https://github.com/ac0d3r/Hyuga

## 📷 预览
<img width="1511" alt="image" src="https://user-images.githubusercontent.com/26270009/157201907-0c6d62b9-4232-457c-a7b4-9dcaee429bd1.png">



## 🌀 oob
- dns
    - dns查询记录(query name, remote address)
    - 支持 dns-rebinding [#🔗](https://github.com/ac0d3r/Hyuga#-dns-rebinding)
- http 
    - http 请求记录(url, method, remote address, raw request)
- ldap & rmi
    - ldap&rmi 请求记录(protocol, remote address, path) 
    > thx: [浅谈Log4j2不借助dnslog的检测](https://4ra1n.love/post/I_AYmmK2J/)

## 👀 其他
- 部署参见 [DEPLOY.md](./DEPLOY.md)
- 📝 更新日志[CHANGELOG.md](./CHANGELOG.md)

### 🚀 查询 API
- `GET` - http://`<hyuga.io>`/api/record/list?type=`<dns|http>`&token=`<token>`&filter=`<filter>`
    - `type`: 查询类型 `dns|http|jndi`
    - `token`: 域名 token
    - `filter`: 过滤字符

### 🪓 DNS Rebinding
查询 `r.xxx.hyuga.io` 时根据访问次数依次返回所设置的dns（无缓存时）。

**e.g.** ip 为 `1.1.1.1`；dns 设置如下：

<img width="420" alt="image" src="https://user-images.githubusercontent.com/26270009/157200281-06a3752b-5b48-45df-b0c4-864d7fc81b13.png">

查询 `r.8q56.hyuga.io` 根据访问次数计算依次返回：`1.1.1.1` -> `127.0.0.1` -> `1.1.1.1`...