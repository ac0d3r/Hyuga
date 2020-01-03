# User Login

简要描述：
- 用户登录

请求URL：
- `http://api.[doname]/v1/users/self/login`


请求方式：
- `POST`

参数：
|  参数名   |  必选  | 类型 | 说明 |
| :------  | :----- | :----- | -------  |
| username | 是     | String | 用户名 |
| password | 是 | String | 密码 |

请求参数示例：
- `{"username":"buzz", "password":"12345678"}`

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
    "data": {
        "username": "buzz",
        "jwtoken": "eyJ0eXAiOiJKV1iOiJIUzI1NiJ9.eycm5hbWUiOiJidXp6IiwiZXhwIjoxNTc4MDI5MjA2fQ.PBKYAqeS6PyLLA5aFwJDCxLBUF5o"
    }
}
```

返回参数说明:

|  参数名   |  类型  | 类型 |
| -------  | :----- | :----- |
| username | String | 登录成功的用户名 |
| jwtoken | String | jwtoken |


