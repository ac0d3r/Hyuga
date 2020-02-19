import json
from datetime import datetime

import falcon

from hyuga.core.errors import NotSupportedError
from hyuga.lib.option import CONFIG

try:
    from collections import OrderedDict
except ImportError:
    OrderedDict = dict


class BaseResource:
    HELLO_WORLD = {
        'server': CONFIG.BRAND_NAME
    }

    def __init__(self):
        self.current_user = None

    def on_error(self, resp, error=None):
        resp.status = error['status']
        meta = OrderedDict()
        meta['code'] = error['code']
        meta['message'] = error['message']

        obj = OrderedDict()
        obj['meta'] = meta
        resp.body = BaseResource.to_json(obj)

    def on_success(self, resp, data=None):
        resp.status = falcon.HTTP_200
        meta = OrderedDict()
        meta['code'] = 200
        meta['message'] = 'OK'

        obj = OrderedDict()
        obj['meta'] = meta
        if data:
            BaseResource.find_datetime_to_ctime(data)
        obj['data'] = data
        resp.body = BaseResource.to_json(obj)

    def on_get(self, req, resp):
        if req.path == '/':
            resp.status = falcon.HTTP_200
            resp.body = BaseResource.to_json(self.HELLO_WORLD)
        else:
            raise NotSupportedError(method='GET', url=req.path)

    def on_post(self, req, resp):
        raise NotSupportedError(method='POST', url=req.path)

    def on_put(self, req, resp):
        raise NotSupportedError(method='PUT', url=req.path)

    def on_delete(self, req, resp):
        raise NotSupportedError(method='DELETE', url=req.path)

    @staticmethod
    def find_datetime_to_ctime(data: list):
        """set the value of the datetime object in data to the value of its ctime method

            datetime object canont serialize to json
        """
        def datetime_to_ctime_in_dict(ditem):
            for k, v in ditem.items():
                if isinstance(v, datetime):
                    ditem[k] = v.ctime()

        if isinstance(data, list):
            for item in data:
                if isinstance(item, dict):
                    datetime_to_ctime_in_dict(item)
        elif isinstance(data, dict):
            datetime_to_ctime_in_dict(data)

    @staticmethod
    def to_json(body: dict):
        return json.dumps(body)

    @staticmethod
    def on_record_http_success():
        """set response body when the 'httplog' is recorded successfully.
        """
        meta = OrderedDict()
        meta['code'] = 201
        meta['message'] = "HTTP Record Insert Success"
        obj = OrderedDict()
        obj['meta'] = meta
        return BaseResource.to_json(obj)
