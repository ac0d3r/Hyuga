# Get User Self Records

简要描述：

- 用户查询自己的 [dns|http] 记录

请求 URL：

- `http://api.[domain]/v1/users/self/records?type=`{dns | http}

请求方式：

- `GET`

参数：

| 参数名 | 必选 | 类型   | 说明                     |
| :----- | :--- | :----- | ------------------------ |
| type   | 是   | String | 查询类型 "dns" 或 "http" |

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
  "data": [
    {
      "created": "Fri Jan  3 05:38:38 2020",
      "uidentify": "wd55te",
      "name": "http://wd55te.hyuga.io/",
      "method": "GET",
      "data": null,
      "user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36",
      "content_type": null,
      "remote_addr": "127.0.0.1",
      "id": "1"
    },
    {
      "created": "Fri Jan  3 05:38:38 2020",
      "uidentify": "wd55te",
      "name": "http://wd55te.hyuga.io/robots.txt",
      "method": "GET",
      "data": null,
      "user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.88 Safari/537.36",
      "content_type": null,
      "remote_addr": "127.0.0.1",
      "id": "2"
    }
  ]
}
```

# Delete User Self Records

简要描述：

- 用户删除自己的 [dns|http] 记录

请求 URL：

- `http://api.[domain]/v1/users/self/records`

请求方式：

- `DELETE`

参数：

| 参数名 | 必选 | 类型   | 说明                     |
| :----- | :--- | :----- | ------------------------ |
| type   | 是   | String | 查询类型 "dns" 或 "http" |

请求头：

| 参数名       | 必选 |       说明       |
| :----------- | :--- | :--------------: |
| Content-Type | 是   | application/json |
| JWToken      | 是   | 登录后的 jwtoken |

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
