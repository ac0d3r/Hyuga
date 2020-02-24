import typing
import uuid

import redis

from hyuga.core.errors import CanNotCreateUserError
from hyuga.lib.utils import get_shortuuid

from .base import BaseModel, PasswordHash, models


class User(BaseModel):
    username = models.CharField(
        max_length=30, unique=True, required=True, indexed=True)
    nickname = models.CharField(max_length=255, default="路人甲")
    password = models.CharField(required=True)
    identify = models.CharField(
        required=True, max_length=32, unique=True, indexed=True)
    token = models.CharField(required=True, max_length=32, unique=True)
    administrator = models.BooleanField(default=False)

    def part_attr_to_dict(self):
        return {
            'username': self.username,
            'nickname': self.nickname,
            'identify': self.identify,
            'token': self.token
        }


class UserToken(BaseModel):
    # 返回的客户端的key
    token = models.CharField(unique=True, required=True)
    # jwt
    jwt = models.CharField(required=True)


def _get_unique_identify(length=6, max_times=(3, 9)) -> typing.Union[str, None]:
    """get unique identify
    :param length: identify length
    :param max_times: 
        0: maximum number of single attempts
        1: maximum number of attempts
    """
    _all_the_times = 0
    _single_times = 0
    while True:
        identify = get_shortuuid(length)
        try:
            user = User.objects.filter(identify=identify).first()
        except redis.exceptions.ConnectionError:
            return None
        if not user:
            break
        _single_times += 1
        _all_the_times += 1
        if _all_the_times > max_times[1]:
            return None
        if _single_times >= max_times[0]:
            _single_times = 0
            length += 1
    return identify


def create_user(username, password, nickname="路人甲"):
    token = uuid.uuid1().hex
    identify = _get_unique_identify(6, 3)
    if identify is None:
        raise CanNotCreateUserError()
    User.objects.create(
        username=username,
        password=PasswordHash.hash_password(password).db_store,
        nickname=nickname,
        identify=identify,
        token=token
    )


def create_superuser(username, password, nickname="管理员"):
    token = uuid.uuid1().hex
    identify = _get_unique_identify(6, 3)
    if identify is None:
        raise CanNotCreateUserError()
    User.objects.create(
        username=username,
        password=PasswordHash.hash_password(password).db_store,
        nickname=nickname,
        identify=identify,
        token=token,
        administrator=True
    )
