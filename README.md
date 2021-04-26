<div align="center" >
    <img src="./docs/hyuga.png" width="280" alt="Hyuga" />
</div>
<p align="center">
    <a href="https://github.com/Buzz2d0/Hyuga">
        <img alt="Hyuga" src="https://img.shields.io/badge/Hyuga-2.1.0-yellow"/>
    </a>
    <img src="https://img.shields.io/badge/Language-Golang-blue" alt="Language" />
    <a href="https://github.com/Buzz2d0/Hyuga/blob/master/LICENSE">
        <img alt="License" src="https://img.shields.io/github/license/Buzz2d0/Hyuga"/>
    </a>
    <a href="https://github.com/Buzz2d0/Hyuga/stargazers">
        <img alt="stars" src="https://img.shields.io/github/stars/Buzz2d0/Hyuga"/>
    </a>
 </p>
⚡️Hyuga 是一个用来检测带外(Out-of-Band)流量(DNS查询和HTTP请求)的监控工具。

---
## 🎉 项目简介
> `Hyuga` 名字来自《火影忍者》中的日向一族。

DEMO 主页：http://hyuga.co/

项目地址：https://github.com/Buzz2d0/Hyuga

### 🚀查询 API
- `GET` - `http://api.<hyuga.io>/v1/records?type=<dns|http>&token=<token>&filter=<filter>`
    - `type`: 查询类型 `dns|http`
    - `token`: 域名 token
    - `filter`: 过滤字符


## 📷 预览
![demo.png](./docs/demo.png)

## 👏 主要框架

- [echo](https://github.com/labstack/echo/)
- [redis](https://github.com/go-redis/redis/)
- [dns](https://github.com/miekg/dns/)
- [Vue.js](https://cn.vuejs.org)

## 📝 更新日志

[CHANGELOG.md](./CHANGELOG.md)

## ⌛ 未来计划

 - [x] 前端重构
 - [ ] 增加其他好用的小工具

## 🙏 致谢

- [DNSLog](http://dnslog.cn)：旧版前端样式借鉴
- [PockyRayzz](https://github.com/PockyRayzz)：使用 Vue 重构前端

## 👀 其他
- 部署参见 [DEPLOY.md](./DEPLOY.md)
- 保留之前的 [Python](https://github.com/Buzz2d0/Hyuga/tree/python) 版本
