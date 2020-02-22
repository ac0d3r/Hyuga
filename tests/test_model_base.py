import unittest
from hyuga.models.base import *


class TestModelsBase(unittest.TestCase):
    """Test hyuga.models.base"""

    def test_hash_password(self):
        password = "123456"
        hashed = PasswordHash.hash_password(password)
        self.assertTrue(isinstance(hashed, PasswordHash))
        self.assertTrue(hashed.check_password(password))

    def test_check_password(self):
        password = "123456"
        hashed = PasswordHash.hash_password(password)
        self.assertTrue(hashed.check_password(password))
        self.assertFalse(hashed.check_password("1234561"))

    def test_db_store(self):
        password = "123456"
        self.assertTrue(isinstance(
            PasswordHash.hash_password(password).db_store, str))

    def test_py_value(self):
        password = "123456"
        db_store_value = PasswordHash.hash_password(password).db_store
        self.assertTrue(PasswordHash.py_value(
            db_store_value).check_password(password))


if __name__ == '__main__':
    unittest.main()
