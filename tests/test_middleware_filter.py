import unittest

from falcon import testing

import hyuga.app
from hyuga.lib.option import CONFIG
from hyuga.models.record import *
from hyuga.models.user import *


class TestMiddlewareFilter(testing.TestCase):
    def create_user(self):
        create_user("test0", "password")
        user = User.objects.filter(username="test0").first()
        if user:
            return user.identify

    def setUp(self):
        super(TestMiddlewareFilter, self).setUp()
        self.app = hyuga.app.create_app()
        self.host = CONFIG.API_DOMAIN
        self.json_headers = {"Content-Type": "application/json"}
        self.test_uid = self.create_user()

    def tearDown(self):
        """clear test users
        """
        for user in User.objects.all():
            if user:
                user.delete()
        for http in HttpRecord.objects.all():
            if http:
                http.delete()


class TestUsersSelfRecords(TestMiddlewareFilter):
    def test_filter_out_others_domain(self):
        doc = {"meta": {"code": 10, "message": "Not Supported",
                        "description": "method: GET, url: http://abc.com/"}}
        res = self.simulate_get('/', host="abc.com")

        self.assertEqual(doc, res.json)

    def test_filter_out_not_in_allow_methods(self):
        doc = {"meta": {"code": 10, "message": "Not Supported",
                        "description": "method: PATCH, url: http://abc.com/"}}
        res = self.simulate_patch('/', host="abc.com")
        self.assertEqual(doc, res.json)

    def test_filter_out_empty_json(self):
        doc = {"meta": {"code": 88, "message": "Invalid Parameter",
                        "description": "A valid JSON document is required"}}
        res = self.simulate_post(
            '/v1/users/self/login', host=self.host, headers=self.json_headers)

        self.assertEqual(doc, res.json)

    def test_filter_out_error_json(self):
        doc = {"meta": {"code": 88, "message": "Invalid Parameter",
                        "description": "No JSON object could be decoded or Malformed JSON"}}
        data = '{"username": niao, "password": 123123}'
        res = self.simulate_post(
            '/v1/users/self/login', host=self.host, headers=self.json_headers, body=data)

        self.assertEqual(doc, res.json)

    def test_filter_out_not_subdomain_name(self):
        doc = {"meta": {"code": 10, "message": "Not Supported",
                        "description": "method: GET, url: http://%s/" % CONFIG.DOMAIN}}
        res = self.simulate_get('/', host=CONFIG.DOMAIN)
        self.assertEqual(doc, res.json)

    def test_nonexist_user_record_http(self):
        doc = {'meta':
               {'code': 10, 'description': 'method: GET, url: http://test1.hyuga.io/',
                'message': 'Not Supported'}}
        res = self.simulate_get('/', host="test1." + CONFIG.DOMAIN)

        self.assertEqual(doc, res.json)

    def test_record_http(self):
        doc = {"meta": {"code": 201, "message": "HTTP Record Insert Success"}}
        res = self.simulate_get('/', host="%s.%s" %
                                (self.test_uid, CONFIG.DOMAIN))

        self.assertEqual(doc, res.json)

        http = HttpRecord.objects.filter(
            uidentify=self.test_uid)
        self.assertTrue(len(http) == 1)


if __name__ == '__main__':
    unittest.main()
