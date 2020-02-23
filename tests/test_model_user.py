import unittest
from hyuga.models.user import _get_unique_identify
from hyuga.models.user import *


class TestModelsUser(unittest.TestCase):
    """Test hyuga.models.user"""

    def setUp(self):
        """create some test users
        """
        for i in range(5):
            create_user("test_%s" % i, "password")

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


if __name__ == '__main__':
    unittest.main()
