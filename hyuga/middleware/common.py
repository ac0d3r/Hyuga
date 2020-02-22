# coding: utf-8
import falcon
from falcon.http_status import HTTPStatus

from hyuga.core.log import logger
from hyuga.lib.option import LOG_LEVEL, CONFIG


class HandleCORS:
    def process_request(self, req, resp):
        if LOG_LEVEL != "info":  # allow all orgin when huyga is debug
            resp.set_header('Access-Control-Allow-Origin', '*')
        else:
            resp.set_header('Access-Control-Allow-Origin',
                            "http://%s" % CONFIG.DOMAIN)

        resp.set_header('Access-Control-Allow-Methods', '*')
        resp.set_header('Access-Control-Allow-Headers', '*')
        resp.set_header('Access-Control-Max-Age', 24 * 3600)

        if req.method == 'OPTIONS':
            raise HTTPStatus(falcon.HTTP_200, body='\n')
