# coding: utf-8
import json
import re

import falcon
import redis
from falcon.http_status import HTTPStatus

from hyuga.api.common import BaseResource
from hyuga.core import errors, log
from hyuga.lib.option import CONFIG
from hyuga.models.record import HttpRecord
from hyuga.models.user import User


class GlobalFilter:
    """全局过滤
    """

    def process_request(self, req, resp):
        log._api_logger.debug(
            "middleware filter GlobalFilter - host: %s path: %s" % (req.host, req.path))
        # filter out ip and others domain
        if not req.host.endswith(CONFIG.DOMAIN):
            raise errors.NotSupportedError(method=req.method, url=req.url)

        # api
        if req.host == CONFIG.API_DOMAIN:
            if not req.method in CONFIG.ALLOW_METHODS:
                raise errors.NotSupportedError(method=req.method, url=req.url)

            log._api_logger.debug("Method: %s, ContentType: %s" %
                                  (req.method, req.content_type))
            if "application/json" != req.content_type:
                req.context["data"] = None
                return
            raw_json = req.bounded_stream.read()
            if not raw_json:
                raise errors.InvalidParameterError(
                    "A valid JSON document is required")
            try:
                req.context["data"] = json.loads(raw_json.decode("utf-8"))
            except UnicodeDecodeError:
                raise errors.InvalidParameterError(
                    "Cannot be decoded by utf-8")
            except ValueError:
                raise errors.InvalidParameterError(
                    "No JSON object could be decoded or Malformed JSON")

        # record *.`CONFIG.DOMAIN`(http)
        elif req.host != CONFIG.API_DOMAIN:
            host = re.search(r"([^\.]+)\.%s" % CONFIG.DOMAIN, req.host)
            if not host:
                raise errors.NotSupportedError(method=req.method, url=req.url)

            uid = host.group(1)
            if not User.objects.filter(identify=uid):
                raise errors.NotSupportedError(method=req.method, url=req.url)

            str_data = req.bounded_stream.read().decode("utf-8").rstrip("")
            try:
                http_record = HttpRecord(
                    uidentify=uid,
                    name=req.url,
                    method=req.method,
                    data=str_data or None,
                    user_agent=req.user_agent or None,
                    content_type=req.content_type or None,
                    remote_addr=req.access_route[0] if req.access_route else None
                )
                http_record.save()
                http_record.expire(CONFIG.RECORDS_EXPIRE)
                raise HTTPStatus(
                    falcon.HTTP_200, body=BaseResource.on_record_http_success())

            except redis.exceptions.ConnectionError:
                raise errors.DatabaseError(errors.ERR_DATABASE_CONNECTION)
