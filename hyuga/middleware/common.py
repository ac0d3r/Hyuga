# coding: utf-8
import falcon
from falcon.http_status import HTTPStatus

from hyuga.database import database


class HandleCORS:
    def process_request(self, req, resp):
        resp.set_header('Access-Control-Allow-Origin', '*')
        resp.set_header('Access-Control-Allow-Methods', '*')
        resp.set_header('Access-Control-Allow-Headers', '*')
        resp.set_header('Access-Control-Max-Age', 24 * 3600)
        if req.method == 'OPTIONS':
            raise HTTPStatus(falcon.HTTP_200, body='\n')



class PeeweeConnection:
    def process_request(self, req, resp):
        database.connect()

    def process_response(self, req, resp, resource, req_succeeded):
        if not database.is_closed():
            database.close()
