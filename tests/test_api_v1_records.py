import unittest

from falcon import testing

import hyuga.app
from hyuga.lib.option import CONFIG
from hyuga.models.record import *
from hyuga.models.user import *


class TestApiRecord(testing.TestCase):

    def create_users(self):
        passwd = "password"
        for i in range(5):
            uname = "test_%s" % i
            if i == 4:
                create_superuser(uname, passwd)
                break
            create_user(uname, passwd)
        user = User.objects.filter(username="test_0").first()
        self.user_token = user.token
        self.create_test_records(uid=user.identify)

    def login_user(self):
        data = {"username": "test_0", "password": "password"}

        res = self.simulate_post(
            '/v1/users/self/login',
            host=self.host,
            headers=self.json_headers,
            json=data)
        self.logined_user_jwtoken = res.json['data']['jwtoken']

    def create_test_records(self, uid):
        for i in range(3):
            name = "test%s.%s" % (i, CONFIG.DOMAIN)
            DnsRecord.objects.create(
                uidentify=uid,
                name=name,
                remote_addr="127.0.0.1"
            )

        for i in range(3):
            if i == 2:
                break
            name = "http://test%s.%s" % (i, CONFIG.DOMAIN)
            HttpRecord.objects.create(
                uidentify=uid,
                name=name,
                method="GET",
                date=None,
                user_agent=None,
                content_type=None,
                remote_addr="127.0.0.1"
            )

    def setUp(self):
        super(TestApiRecord, self).setUp()
        self.app = hyuga.app.create_app()
        self.host = CONFIG.API_DOMAIN
        self.json_headers = {"Content-Type": "application/json"}
        self.jwt_headers = lambda value: {"JWToken": value}
        # creat test user
        self.create_users()
        # test login
        self.login_user()

    def tearDown(self):
        """clear test users
        """
        for user in User.objects.all():
            if user:
                user.delete()
        for ut in UserToken.objects.all():
            if ut:
                ut.delete()
        for dns in DnsRecord.objects.all():
            if dns:
                dns.delete()
        for http in HttpRecord.objects.all():
            if http:
                http.delete()


class TestUsersSelfRecords(TestApiRecord):
    def test_get_dns(self):
        res = self.simulate_get('/v1/users/self/records?type=dns',
                                headers=self.jwt_headers(
                                    self.logined_user_jwtoken),
                                host=CONFIG.API_DOMAIN)

        self.assertEqual({"code": 200, "message": "OK"}, res.json["meta"])
        self.assertTrue(len(res.json["data"]) == 3)

    def test_get_http(self):
        res = self.simulate_get('/v1/users/self/records?type=http',
                                headers=self.jwt_headers(
                                    self.logined_user_jwtoken),
                                host=CONFIG.API_DOMAIN)

        self.assertEqual({"code": 200, "message": "OK"}, res.json["meta"])
        self.assertTrue(len(res.json["data"]) == 2)

    def test_get_error_type(self):
        doc = {'meta': {'code': 88, 'message': 'Invalid Parameter', 'description': {
            'type': ["value does not match regex '(dns|http)'"]}}}
        res = self.simulate_get('/v1/users/self/records?type=https',
                                headers=self.jwt_headers(
                                    self.logined_user_jwtoken),
                                host=CONFIG.API_DOMAIN)
        self.assertTrue(doc == res.json)

    def test_delete(self):
        headers = {}
        headers.update(self.json_headers)
        headers.update(self.jwt_headers(self.logined_user_jwtoken))
        doc = {"meta": {"code": 200, "message": "OK"}, "data": None}
        res = self.simulate_delete('/v1/users/self/records',
                                   json={"type": "dns"},
                                   headers=headers,
                                   host=CONFIG.API_DOMAIN)
        self.assertTrue(doc == res.json)


class TestRecords(TestApiRecord):
    def test_get(self):
        res = self.simulate_get('/v1/records?type=dns&token=%s' % self.user_token,
                                headers=self.jwt_headers(
                                    self.logined_user_jwtoken),
                                host=CONFIG.API_DOMAIN)

        self.assertEqual({"code": 200, "message": "OK"}, res.json["meta"])
        self.assertTrue(len(res.json["data"]) == 3)

    def test_get_filter(self):
        res = self.simulate_get('/v1/records?type=dns&token=%s&filter=test0' % self.user_token,
                                headers=self.jwt_headers(
                                    self.logined_user_jwtoken),
                                host=CONFIG.API_DOMAIN)

        self.assertEqual({"code": 200, "message": "OK"}, res.json["meta"])
        self.assertTrue(len(res.json["data"]) == 1)


if __name__ == '__main__':
    unittest.main()
