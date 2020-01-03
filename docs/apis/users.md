# Get Users

简要描述：
- 获取到所有注册的用户（仅限管理员）

请求URL：
- `http://api.[doname]/v1/users`


请求方式：
- `GET`

参数：
- 无

请求头：

|  参数名   |  必选  |  说明 |
| :------  | :----- | :----: |
| JWToken  | 是     | 登录后的 jwtoken|

返回示例：

```json
{
    "meta": {
        "code": 200,
        "message": "OK"
    },
    "data": [
        {
            "username": "admin",
            "nickname": "路人甲",
            "identify": "admin.hyuga.io",
            "token": ""
        },
        {
            "username": "buzz",
            "nickname": "路人甲",
            "identify": "233333",
            "token": "2333334c26dc2333331facde48233333"
        }
    ]
}
```

# Create User

简要描述：
- 用户注册

请求URL：
- `http://api.[doname]/v1/users`


请求方式：
- `POST`

参数：
|  参数名   |  必选  | 类型 | 说明 |
| :------  | :----- | :----- | -------  |
| username | 是     | String | 用户名 |
| password | 是 | String | 密码 |
| nickname | 否 | string | 昵称 |

请求参数示例：
- `{"username":"buzz", "password":"12345678", "nickname":"路人甲"}`

请求头：

|  参数名   |  必选  |  说明 |
| :------  | :----- | :----: |
| Content-Type  | 是     | application/json|


返回示例：

```json
{
    "meta": {
        "code": 200,
        "message": "OK"
    },
    "data": null
}
```
