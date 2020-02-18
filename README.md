<div align="center" >
    <img src="./docs/hyuga.png" width="280" alt="Hyuga" />
</div>
<p align="center">
    <a href="https://github.com/Buzz2d0/Hyuga/blob/master/Pipfile">
        <img alt="Python" src="https://img.shields.io/badge/python-3.7-blue">
    </a>
    <a href="https://github.com/Buzz2d0/Hyuga/blob/master/LICENSE">
        <img alt="License" src="https://img.shields.io/github/license/Buzz2d0/Hyuga">
    </a>
    <a href="https://github.com/Buzz2d0/Hyuga/stargazers">
        <img alt="stars" src="https://img.shields.io/github/stars/Buzz2d0/Hyuga">
    </a>
 </p>
⚡️Hyuga 是一个用来检测带外(Out-of-Band)流量(DNS查询和HTTP请求)的监控工具。

---

## 🎉 项目简介

DEMO 主页：http://hyuga.co/

项目地址：https://github.com/Buzz2d0/Hyuga

Hyuga 与 `ceye` 相似，只因兴趣使然决定写出这款工具将其开源。

Hyuga 的名字来自《火影忍者》中的日向一族的名称，日向一族是火影忍者中的火之国木叶忍者村的一个氏族，此氏族拥有血继限界白眼。

## 🚀 上手指南

在 [DEMO](http://hyuga.co/) 上简单注册账号后就可以使用了。

`📌tips`：目前没有使用邮箱或者手机号进行账号注册导致 Forgot password 功能没有实现。

**🐋Docker 安装：**

先阅读此文档：[deploy.md](./docs/deploy.md) 。进行简单设置后，就可以使用 docker 运行。

**🏝Others**

- 目前 Hyuga API 功能都已经放入此目录下了
  [apis/](https://github.com/Buzz2d0/Hyuga/blob/master/docs/apis/)
- API 接口查询：[records.md](./docs/apis/records.md)

## 👏 主要框架

- [Falcon](https://github.com/falconry/falcon) 用于构建快速 WEB API 和应用程序后端的极简 WSGI 库。
- [cerberus](https://github.com/pyeve/cerberus) 轻量级、可扩展的数据验证库。
- [peewee](https://github.com/coleifer/peewee) 简单小巧的 ORM。
- [redisco](https://github.com/chen2aaron/redisco) 简单好用的 Redis ORM。
- [dnslib](https://pypi.org/project/dnslib/) 用于对 DNS 数据包进行编码/解码。
- [click](https://github.com/pallets/click) 组合命令行界面工具包。
- [React](https://github.com/facebook/react) 用于构建用户界面的 JavaScript 库。

## ⌛ 后续计划

**Frontend**

- [x] 注册成功路由跳转登录
- [x] 登录和注册时对返回的错误消息显示
- [ ] 添加管理员可用的功能(查看、删除其他用户等等)

## 🙏 参考

- [DNSLog](https://github.com/BugScanTeam/DNSLog)
- [ceye.io](http://ceye.io)
