import uuid

from peewee import BooleanField, CharField, ForeignKeyField

from hyuga.lib.utils import get_shortuuid

from .base import BaseModel, PasswordField


class User(BaseModel):
    username = CharField(max_length=30, unique=True, null=False)
    nickname = CharField(max_length=255, default="路人甲")
    password = PasswordField(null=False)
    identify = CharField(max_length=32, unique=True, null=False)
    token = CharField(max_length=32, unique=True, null=False)
    administrator = BooleanField(default=False)

    def model_to_dict(self):
        return {
            'username': self.username,
            'nickname': self.nickname,
            'identify': self.identify,
            'token': self.token
        }


def __get_unique_identify(length=6, max_times=3) -> str:
    has_times = 0
    while True:
        identify = get_shortuuid(length)
        try:
            User.get(User.identify == identify)
            has_times += 1
            if has_times >= max_times:
                has_times = 0
                length += 1
        except User.DoesNotExist:
            break
    return identify


def create_user(username, password, nickname="路人甲"):
    token = uuid.uuid1().hex
    identify = __get_unique_identify(6, 3)
    User.create(
        username=username,
        password=password,
        nickname=nickname,
        identify=identify,
        token=token
    )
