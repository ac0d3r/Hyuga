import re

from .base import BaseModel, models
from hyuga.lib.option import CONFIG


def hyuga_domain_validator(field_name, value):
    domain = CONFIG.DOMAIN.replace(".", r"\.")
    regex = r"([a-zA-Z0-9-]+\.)+%s" % domain
    if not re.match(regex, value):
        return ((field_name, 'Not DnsRecord name'),)


def hyuga_url_validator(field_name, value):
    domain = CONFIG.DOMAIN.replace(".", r"\.")
    regex = r"https?:\/\/([a-zA-Z0-9-]+\.)+{domain}\/?[-a-zA-Z0-9()@:%_\+.~#?&//=]*".format(
        domain=domain)
    if not re.match(regex, value):
        return ((field_name, 'Not HttpRecord name'),)


class DnsRecord(BaseModel):
    # User.identify
    uidentify = models.Attribute(required=True, indexed=True)
    # 请求解析的域名
    name = models.Attribute(required=True, validator=hyuga_domain_validator)
    remote_addr = models.Attribute(required=True)


class HttpRecord(BaseModel):
    # User.identify
    uidentify = models.Attribute(required=True, indexed=True)
    # url
    name = models.Attribute(required=True, validator=hyuga_url_validator)
    method = models.Attribute(required=True)
    remote_addr = models.Attribute(required=True)
    data = models.Attribute()
    user_agent = models.Attribute()
    content_type = models.Attribute()
