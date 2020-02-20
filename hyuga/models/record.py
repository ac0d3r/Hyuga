import re

from redisco import models


def domain_name_validator(field_name, value):
    if not re.match(r'^(?=^.{3,255}$)[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+$', value):
        return ((field_name, 'Not domain name'),)


class BaseModel(models.Model):
    created = models.DateTimeField(auto_now_add=True)


class DnsRecord(BaseModel):
    # User.identify
    uidentify = models.Attribute(required=True, indexed=True)
    # 请求解析的域名
    name = models.Attribute(required=True, validator=domain_name_validator)
    remote_addr = models.Attribute(required=True)


class HttpRecord(BaseModel):
    # User.identify
    uidentify = models.Attribute(required=True, indexed=True)
    # url
    name = models.Attribute(required=True)
    method = models.Attribute(required=True)
    data = models.Attribute()
    user_agent = models.Attribute()
    content_type = models.Attribute()
    remote_addr = models.Attribute()
