import re

from redisco import connection_setup, models

from hyuga.lib.option import CONFIG

connection_setup(host=CONFIG.REDIS_SERVER, port=int(
    CONFIG.REDIS_PROT), db=int(CONFIG.REDIS_DB))


def not_domain_name(field_name, value):
    if isinstance(value, str):
        if not re.match(r'^(?=^.{3,255}$)[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+$', value):
            return ((field_name, 'Not domain name'),)


class BaseModel(models.Model):
    created = models.DateTimeField(auto_now_add=True)


class DnsRecord(BaseModel):
    uidentify = models.Attribute(required=True, indexed=True)
    name = models.Attribute(required=True, validator=not_domain_name)
    remote_addr = models.Attribute(required=True)


class HttpRecord(BaseModel):
    uidentify = models.Attribute(required=True, indexed=True)
    name = models.Attribute(required=True)
    method = models.Attribute(required=True)
    data = models.Attribute()
    user_agent = models.Attribute()
    content_type = models.Attribute()
    remote_addr = models.Attribute()
