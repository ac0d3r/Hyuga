import re
import uuid

import falcon
import jwt
import redis

from hyuga.api.common import FIELDS, BaseResource, BaseValidate
from hyuga.core import errors
from hyuga.core.auth import authenticated
from hyuga.core.log import _api_logger
from hyuga.lib.option import CONFIG
from hyuga.models.user import PasswordHash, User, UserToken, create_user


class Users(BaseResource):
    """
    Handle for endpoint: /v1/users
    """
    postSchema = {
        "username": FIELDS["username"],
        "password": FIELDS["password"],
        "nickname": FIELDS["nickname"]
    }

    @authenticated(verifyadmin=True)
    def on_get(self, req, resp):
        try:
            resp_data = [user.part_attr_to_dict()
                         for user in User.objects.all()]
            self.on_success(resp, resp_data)
        except redis.exceptions.ConnectionError:
            raise errors.DatabaseError(errors.ERR_DATABASE_CONNECTION)
        except Exception as e:
            _api_logger.info("Users on_get ERROR: %s" % e)
            self.on_error(resp, errors.ERR_UNKNOWN)

    @falcon.before(BaseValidate(postSchema).validate)
    def on_post(self, req, resp):
        req_data = req.context["data"]
        username = req_data["username"]
        password = req_data["password"]

        try:
            user = User.objects.filter(username=username).first()
        except redis.exceptions.ConnectionError:
            raise errors.DatabaseError(errors.ERR_DATABASE_CONNECTION)
        except Exception as e:
            _api_logger.info("Users on_post ERROR: %s" % e)
            self.on_error(resp, errors.ERR_UNKNOWN)

        if user:
            self.on_error(resp, errors.ERR_USER_ALREADY_EXISTS)
            return
        kwargs = {"username": username, "password": password}
        if "nickname" in req_data.keys():
            kwargs["nickname"] = req_data["nickname"]
        try:
            r = create_user(**kwargs)
        except errors.UserUnableIdentifyError:
            self.on_error(resp, errors.ERR_CREATE_USER)
        except redis.exceptions.ConnectionError:
            raise errors.DatabaseError(errors.ERR_DATABASE_CONNECTION)
        else:
            if not r is True:
                infos = "".join(["field name %s %s" %
                                 (err[0], err[1]) for err in r])
                raise errors.InvalidParameterError("Create User %s" % infos)
        self.on_success(resp)


class UsersItem(BaseResource):
    """
    Handle for endpoint: /v1/users/{user_id}
    """

    @authenticated(verifyadmin=True)
    def on_get(self, req, resp, user_id):
        try:
            user = User.objects.get_by_id(user_id)
        except redis.exceptions.ConnectionError:
            raise errors.DatabaseError(errors.ERR_DATABASE_CONNECTION)
        except Exception as e:
            _api_logger.info("UsersItem on_get ERROR: %s" % e)
            self.on_error(resp, errors.ERR_UNKNOWN)
        else:
            if user:
                self.on_success(resp, user.part_attr_to_dict())
            else:
                raise errors.UserNotExistsError()

    @authenticated(verifyadmin=True)
    def on_delete(self, req, resp, user_id):
        try:
            user = User.objects.get_by_id(user_id)
        except redis.exceptions.ConnectionError:
            raise errors.DatabaseError(errors.ERR_DATABASE_CONNECTION)
        except Exception as e:
            _api_logger.info("UsersItem on_delete ERROR: %s" % e)
            self.on_error(resp, errors.ERR_UNKNOWN)
        else:
            if user:
                user.delete()
                self.on_success(resp, {"id": user_id})
            else:
                raise errors.UserNotExistsError()


class UsersSelf(BaseResource):
    """
    Handle for endpoint: /v1/users/self
    """
    @authenticated()
    def on_get(self, req, resp):
        self.on_success(resp, self.current_user.part_attr_to_dict())


class UsersSelfOperation(BaseResource):
    """
    Handle for endpoint: /v1/users/self/[operation]
    """
    LOGIN = "login"
    LOGOUT = "logout"
    MODIFY = "modify"
    postSchema = {
        "username": FIELDS["username"],
        "password": FIELDS["password"]
    }

    def on_post(self, req, resp):
        cmd = re.split("\\W+", req.path)[-1:][0]
        if cmd == self.LOGIN:
            self.process_login(req, resp)
        elif cmd == self.MODIFY:
            self.process_modify(req, resp)
        elif cmd == self.LOGOUT:
            self.process_logout(req, resp)

    @falcon.before(BaseValidate(postSchema).validate)
    def process_login(self, req, resp):
        req_data = req.context["data"]
        username = req_data["username"]
        password = req_data["password"]
        try:
            user = User.objects.filter(username=username).first()
        except redis.exceptions.ConnectionError:
            raise errors.DatabaseError(errors.ERR_DATABASE_CONNECTION)
        except Exception as e:
            _api_logger.info("UsersSelfOperation login ERROR: %s" % e)
            self.on_error(resp, errors.ERR_UNKNOWN)
        _api_logger.debug("password: %s, %s" % (password, user.password))
        if not user:
            raise errors.UnauthorizedError()
        if not PasswordHash.py_value(user.password).check_password(password):
            raise errors.UnauthorizedError()

        jwtoken = jwt.encode(
            payload={"id": user.id,
                     "username": user.username
                     },
            key=CONFIG.SECRET_KEY,
            algorithm=CONFIG.JWT_ALGORITHM).decode("utf8")

        user_token = UserToken(token=uuid.uuid1().hex, jwt=jwtoken)
        if user_token.save():
            user_token.expire(CONFIG.JWT_EXPIRE)
        self.on_success(
            resp, data={"username": username, "jwtoken": user_token.token})

    def process_modify(self, req, resp):
        pass

    @authenticated()
    def process_logout(self, req, resp):
        try:
            user_token = UserToken.objects.filter(
                token=req.get_header("JWToken")).first()
            if user_token:
                user_token.delete()
            self.on_success(resp)
        except redis.exceptions.ConnectionError:
            raise errors.DatabaseError(errors.ERR_DATABASE_CONNECTION)
        except Exception as e:
            _api_logger.info("UsersSelfOperation login ERROR: %s" % e)
            self.on_error(resp, errors.ERR_UNKNOWN)
