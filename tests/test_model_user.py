import unittest
from hyuga.models.user import _get_unique_identify
from hyuga.models.user import *


class TestModelsUser(unittest.TestCase):
    """Test hyuga.models.user"""

    def setUp(self):
        """create some test users
        """
        passwd = "password"
        for i in range(5):
            uname = "test_%s" % i
            if i == 4:
                create_superuser(uname, passwd)
                break
            create_user(uname, passwd)

    def tearDown(self):
        """clear test users
        """
        for i in range(5):
            user = User.objects.filter(username="test_%s" % i).first()
            if user:
                user.delete()

    def test_get_unique_identify(self):
        for _ in range(50):
            self.assertNotEqual(_get_unique_identify(), _get_unique_identify())

    def test_user_part_attr_to_dict(self):
        username = "test_0"
        user = User.objects.filter(username=username).first()
        if user:
            _tmp = user.part_attr_to_dict()
            print("get the username '%s' - part_attr_dict: " % username, _tmp)
            self.assertTrue(isinstance(_tmp, dict))
            self.assertEqual(
                ["username", "nickname", "identify", "token"], list(_tmp.keys()))
            self.assertEqual(username, _tmp["username"])

    def test_create_superuser(self):
        users = User.objects.filter(administrator=True)
        self.assertTrue(len(users) == 1)
        self.assertTrue(users[0].administrator == True)


if __name__ == '__main__':
    unittest.main()
