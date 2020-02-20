# coding: utf-8
import json
import re

import falcon
import redis
from falcon.http_status import HTTPStatus

from hyuga.api.common import BaseResource
from hyuga.core.errors import (ERR_DATABASE_CONNECTION, DatabaseError,
                               InvalidParameterError, NotSupportedError)
from hyuga.core.log import _api_logger
from hyuga.lib.option import CONFIG, LOG_LEVEL
from hyuga.models.record import HttpRecord


class GlobalFilter:
    """全局过滤
    """

    def __init__(self, debug=None):
        self.debug = debug
        if debug is None:
            self.debug = False
            if LOG_LEVEL != "info":
                self.debug = True

    def process_request(self, req, resp):
        _api_logger.debug(
            'middleware filter GlobalFilter - host: %s path: %s' % (req.host, req.path))
        # filter out ip and others domain
        if not self.debug and \
                not req.host.endswith(CONFIG.DOMAIN):
            raise NotSupportedError(method=req.method, url=req.url)

        # api
        if req.host == CONFIG.API_DOMAIN:
            if not req.method in CONFIG.ALLOW_METHODS:
                raise NotSupportedError(method=req.method, url=req.url)

            _api_logger.debug('Method: %s, ContentType: %s' %
                              (req.method, req.content_type))
            if 'application/json' != req.content_type:
                req.context['data'] = None
                return
            try:
                raw_json = req.stream.read()
            except Exception:
                raise falcon.HTTPBadRequest('Bad request', 'Read Error')
            try:
                req.context['data'] = json.loads(raw_json.decode('utf-8'))
            except UnicodeDecodeError:
                raise InvalidParameterError('Cannot be decoded by utf-8')
            except ValueError:
                raise InvalidParameterError(
                    'No JSON object could be decoded or Malformed JSON')

        # record *.`CONFIG.DOMAIN`(http)
        elif req.host != CONFIG.API_DOMAIN:
            if not self.debug:
                host = re.search(r'([^\.]+)\.%s' % CONFIG.DOMAIN, req.host)
                if not host:
                    raise NotSupportedError(method=req.method, url=req.url)

            str_data = req.stream.read().decode("utf-8").rstrip("")
            try:
                http_record = HttpRecord(
                    uidentify=host.group(1),
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
                raise DatabaseError(ERR_DATABASE_CONNECTION)
