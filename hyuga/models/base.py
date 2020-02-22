import bcrypt

from redisco import models
from redisco.models import Attribute


class BaseModel(models.Model):
    created = models.DateTimeField(auto_now_add=True)


class PasswordHash(bytes):

    @property
    def db_store(self) -> str:
        """convert to character for store in redis
        """
        return self.decode("utf-8")

    def check_password(self, password: str) -> bool:
        password = password.encode("utf-8")
        return bcrypt.checkpw(password, self)

    @staticmethod
    def hash_password(password: str, salt=12):
        password = password.encode("utf-8")
        return PasswordHash(bcrypt.hashpw(password, bcrypt.gensalt(salt)))

    @staticmethod
    def py_value(hasded_password: str):
        return PasswordHash(hasded_password.encode("utf-8"))
