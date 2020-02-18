# Get User

简要描述：

- 获取某个用户信息（仅限管理员）

请求 URL：

- `http://api.[doname]/v1/users/[user id]`

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
    "nickname": "路人家",
    "identify": "233333",
    "token": "2333334c26dc2333331facde48233333"
  }
}
```

# Delete User

简要描述：

- 删除某个用户信息（仅限管理员）

请求 URL：

- `http://api.[doname]/v1/users/[user id]`

请求方式：

- `DELETE`

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
    "id": 6
  }
}
```

返回参数说明:

| 参数名 | 类型 | 类型            |
| :----- | :--- | :-------------- |
| id     | Int  | 被删除的用户 id |
