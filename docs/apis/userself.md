# User Self

简要描述：

- 获取到用户自身的信息

请求 URL：

- `http://api.[doname]/v1/users/self`

请求方式：

- `GET`

参数：

- 无

请求头：

| 参数名  | 必选 |       说明       |
| :------ | :--- | :--------------: |
| JWToken | 是   | 登录后的 jwtoken |

返回示例：

```json
{
  "meta": {
    "code": 200,
    "message": "OK"
  },
  "data": {
    "username": "buzz",
    "nickname": "路人甲",
    "identify": "233333",
    "token": "2333334c26dc2333331facde48233333"
  }
}
```
