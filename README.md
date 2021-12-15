<div align="center" >
    <img src="./docs/hyuga.png" width="280" alt="Hyuga" />
</div>
<p align="center">
    <a href="https://github.com/Buzz2d0/Hyuga">
        <img alt="Hyuga" src="https://img.shields.io/badge/Hyuga-3.0.1-yellow"/>
    </a>
    <img src="https://img.shields.io/badge/Language-Golang-blue" alt="Language" />
    <a href="https://github.com/Buzz2d0/Hyuga/blob/master/LICENSE">
        <img alt="License" src="https://img.shields.io/github/license/Buzz2d0/Hyuga"/>
    </a>
    <a href="https://github.com/Buzz2d0/Hyuga/stargazers">
        <img alt="stars" src="https://img.shields.io/github/stars/Buzz2d0/Hyuga"/>
    </a>
 </p>

⚡️Hyuga 是一个用来检测带外(Out-of-Band)流量的监控工具。

---
## 🎉 项目简介

DEMO 主页：http://hyuga.icu

项目地址：https://github.com/Buzz2d0/Hyuga

## 📷 预览
![demo.png](./docs/demo.png)


## 🌀 oob
- dns
    - 记录dns查询记录(query name, remote address)
    - 支持 dns-rebinding [#🔗](https://github.com/Buzz2d0/Hyuga#-dns-rebinding)
- http 
    - 记录 http 请求记录(url, method, remote address, raw request)

## 👀 其他
- 部署参见 [DEPLOY.md](./DEPLOY.md)
- 📝 更新日志[CHANGELOG.md](./CHANGELOG.md)

### 🚀 查询 API
- `GET` - http://`<hyuga.io>`/api/record/list?type=`<dns|http>`&token=`<token>`&filter=`<filter>`
    - `type`: 查询类型 `dns|http`
    - `token`: 域名 token
    - `filter`: 过滤字符

### 🪓 DNS Rebinding
查询 `r.xxx.hyuga.io` 时根据访问次数依次返回所设置的dns（无缓存时）。

**e.g.** ip 为 `1.1.1.1`；dns 设置如下：

<img src="https://user-images.githubusercontent.com/26270009/146206555-49450822-44b7-46f4-8942-b6bf831d76f8.png" width="420"/>

查询 `r.8q56.hyuga.io` 根据访问次数计算依次返回：`1.1.1.1` -> `127.0.0.1` -> `1.1.1.1`...
