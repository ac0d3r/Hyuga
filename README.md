<div align="center" >
    <img src="./docs/hyuga.png" width="280" alt="Hyuga" />
</div>
<p align="center">
    <a href="https://github.com/Buzz2d0/Hyuga">
        <img alt="Hyuga" src="https://img.shields.io/badge/Hyuga-1.0.0-yellow"/>
    </a>
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
> `Hyuga` 的名字来自《火影忍者》中的日向一族的名称。


**通知📣!!!**
> 用 golang 重写了此项目，但保留之前的 [Python](https://github.com/Buzz2d0/Hyuga/tree/python) 版本。
> 参考了 [DNSLog.cn](http://dnslog.cn/) 平台, 丢弃了一些冗余的功能，保留查询 API。

项目地址：https://github.com/Buzz2d0/Hyuga

DEMO 主页：http://hyuga.co/
![demo.png](./docs/demo.png)

## 👏 主要框架

- [echo](https://github.com/labstack/echo/)
- [redis](https://github.com/go-redis/redis/)
- [dns](https://github.com/miekg/dns/)

## ⌛ Others

- **Deploy：** 阅读此文档：[deploy.md](./docs/deploy.md) 。
- **FrontEnd：**  使用了 dnslog.cn UI 样式。

## 🙏 Thx

- http://dnslog.cn/

## 更新

- 2020.10.27 支持 dns-rebinding (todo: ui 需要重构)