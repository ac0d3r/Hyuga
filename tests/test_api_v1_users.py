import unittest

from falcon import testing

import hyuga.app
from hyuga.lib.option import CONFIG
from hyuga.models.user import *


class TestApiUser(testing.TestCase):

    def login_users(self):
        data = {"username": "test_0", "password": "password"}

        res = self.simulate_post(
            '/v1/users/self/login',
            host=self.host,
            headers=self.json_headers,
            json=data)

        self.logined_user_jwtoken = res.json['data']['jwtoken']

        data = {"username": "test_4", "password": "password"}
        res = self.simulate_post(
            '/v1/users/self/login',
            host=self.host,
            headers=self.json_headers,
            json=data)

        self.logined_superuser_jwtoken = res.json['data']['jwtoken']

    def setUp(self):
        super(TestApiUser, self).setUp()
        self.app = hyuga.app.create_app()
        self.host = CONFIG.API_DOMAIN
        self.json_headers = {"Content-Type": "application/json"}
        self.jwt_headers = lambda value: {"JWToken": value}
        # creat test user
        passwd = "password"
        for i in range(5):
            uname = "test_%s" % i
            if i == 4:
                create_superuser(uname, passwd)
                break
            create_user(uname, passwd)
        # test login
        self.login_users()

    def tearDown(self):
        """clear test users
        """
        for user in User.objects.all():
            if user:
                user.delete()
        for ut in UserToken.objects.all():
            if ut:
                ut.delete()


class TestUsers(TestApiUser):
    def test_get_not_carry_jwt(self):
        # 不携带 jwtoken
        doc = {"meta": {"code": 99, "message": "Authentication Required",
                        "description": "JWToken Not Exists"}}
        res = self.simulate_get('/v1/users', host=self.host)

        self.assertEqual(res.json, doc)

    def test_get_error_jwt(self):
        doc = {'meta': {'code': 99, 'message': 'Authentication Required',
                        'description': 'Expired Signature'}}
        res = self.simulate_get(
            '/v1/users', host=self.host, headers=self.jwt_headers("abcd123"))
        self.assertEqual(res.json, doc)

    def test_get_low_authority(self):
        doc = {"meta": {"code": 10, "message": "Not Supported",
                        "description": "method: GET, url: /v1/users"}}
        res = self.simulate_get(
            '/v1/users', host=self.host, headers=self.jwt_headers(self.logined_user_jwtoken))

        self.assertEqual(doc, res.json)

    def test_get(self):
        res = self.simulate_get(
            '/v1/users', host=self.host, headers=self.jwt_headers(self.logined_superuser_jwtoken))

        self.assertEqual({"code": 200, "message": "OK"}, res.json["meta"])
        self.assertEqual(5, len(res.json["data"]))

    def test_post(self):
        data = {"username": "test_5", "password": "password"}
        doc = {"meta": {"code": 200, "message": "OK"}, "data": None}

        res = self.simulate_post(
            '/v1/users', json=data, headers=self.json_headers, host=self.host)
        self.assertEqual(res.json, doc)

        # 重复创建一个用户
        data = {"username": "test_5", "password": "password"}
        doc = {"meta": {"code": 25, "message": "Users Already Exists"}}
        res = self.simulate_post(
            '/v1/users', json=data, headers=self.json_headers, host=self.host)
        self.assertEqual(res.json, doc)

    def test_post_data_validate(self):
        # 数据验证
        data = {"username": None, "password": "password"}
        doc = {'meta': {'code': 88, 'description': {'username': [
            'null value not allowed']}, 'message': 'Invalid Parameter'}}
        res = self.simulate_post(
            '/v1/users', json=data, headers=self.json_headers, host=self.host)
        self.assertEqual(res.json, doc)


class TestUsersItem(TestApiUser):
    def test_get(self):
        user = User.objects.filter(username="test_0").first()
        if user:
            res = self.simulate_get('/v1/users/%s' % user.id, host=self.host,
                                    headers=self.jwt_headers(self.logined_superuser_jwtoken))
            self.assertEqual({"code": 200, "message": "OK"}, res.json["meta"])
            self.assertIn("username", res.json["data"])

    def test_delete(self):
        user = User.objects.filter(username="test_0").first()
        if user:
            doc = {"meta": {"code": 200, "message": "OK"},
                   "data": {"id": int(user.id)}}
            res = self.simulate_delete('/v1/users/%s' % user.id, host=self.host,
                                       headers=self.jwt_headers(self.logined_superuser_jwtoken))
            self.assertEqual(res.json, doc)


class TestUsersSelf(TestApiUser):
    def test_on_get(self):
        res = self.simulate_get('/v1/users/self', host=self.host,
                                headers=self.jwt_headers(self.logined_user_jwtoken))

        self.assertEqual({"code": 200, "message": "OK"}, res.json["meta"])
        self.assertEqual("test_0", res.json["data"]["username"])


class TestUsersSelfOperation(TestApiUser):
    def test_login(self):
        username = "test_0"
        password = "password"
        data = {"username": username, "password": password}

        res = self.simulate_post(
            '/v1/users/self/login',
            host=self.host,
            headers=self.json_headers,
            json=data)

        self.assertEqual({"code": 200, "message": "OK"}, res.json["meta"])
        self.assertIn(username, res.json["data"]["username"])
        self.assertIn("jwtoken", res.json["data"].keys())

    def test_login_error_password(self):
        username = "test_0"
        password = "passw0rd"
        data = {"username": username, "password": password}
        doc = {"meta": {"code": 99, "message": "Authentication Required"}}

        res = self.simulate_post(
            '/v1/users/self/login',
            host=self.host,
            headers=self.json_headers,
            json=data)

        self.assertEqual(doc, res.json)

    def test_logout(self):
        doc = {"meta": {"code": 200, "message": "OK"}, "data": None}
        res = self.simulate_post('/v1/users/self/logout', host=self.host,
                                 headers=self.jwt_headers(self.logined_user_jwtoken))

        self.assertEqual(doc, res.json)


if __name__ == '__main__':
    unittest.main()
