from datetime import datetime

from bcrypt import gensalt, hashpw
from peewee import BlobField, DateTimeField, Model

from hyuga.core.database import database


class BaseModel(Model):
    created = DateTimeField(default=datetime.utcnow())
    modified = DateTimeField(default=datetime.utcnow())

    class Meta:
        database = database


class PasswordHash(bytes):
    def check_password(self, password) -> bool:
        password = password.encode('utf-8')
        return hashpw(password, self) == self


class PasswordField(BlobField):
    def __init__(self, iterations=12, *args, **kwargs):
        if None in (hashpw, gensalt):
            raise ValueError(
                'Missing library required for PasswordField: bcrypt')
        self.bcrypt_iterations = iterations
        self.raw_password = None
        super(PasswordField, self).__init__(*args, **kwargs)

    def db_value(self, value):
        """Convert the python value for storage in the database."""
        if isinstance(value, PasswordHash):
            return bytes(value)

        if isinstance(value, str):
            value = value.encode('utf-8')
        salt = gensalt(self.bcrypt_iterations)
        return value if value is None else hashpw(value, salt)

    def python_value(self, value):
        """Convert the database value to a pythonic value."""
        if isinstance(value, str):
            value = value.encode('utf-8')

        return PasswordHash(value)
