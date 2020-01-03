# Records

简要描述：
- 通过 token 查询 [dns|http] 记录

请求URL：
- `http://api.[domain]/v1/records?token=`{token}`&type=`{dns | http}`&filter=`{filter}`

请求方式：
- `GET`

参数：
|  参数名   |  必选  | 类型 | 说明 |
| :------  | :----- | :----- | -------  |
| token | 是     | String | 用户的api token |
| type | 是 | String | 查询类型 "dns" 或 "http"         |
| filter | 否 | String | 匹配名称规则，筛选器最大长度为20 |

请求示例：

- `http://api.hyuga.io/v1/records?token=6a6f9a4c26dc11eab41facde41323&type=http`


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