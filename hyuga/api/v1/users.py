import re
from datetime import datetime

import falcon
import jwt

from hyuga.api.common import FIELDS, BaseResource, BaseValidate
from hyuga.core.auth import authenticated
from hyuga.core.errors import (ERR_UNKNOWN, ERR_USER_ALREADY_EXISTS,
                               InvalidParameterError, NotSupportedError,
                               UnauthorizedError, UserNotExistsError)
from hyuga.core.log import _api_logger
from hyuga.lib.option import CONFIG
from hyuga.models.user import User, create_user


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
            resp_data = [user.model_to_dict() for user in User.select()]
            self.on_success(resp, resp_data)
        except Exception as e:
            _api_logger.info("Users on_get ERROR: %s" % e)
            self.on_error(resp, ERR_UNKNOWN)

    @falcon.before(BaseValidate(postSchema).validate)
    def on_post(self, req, resp):
        req_data = req.context["data"]
        if not req_data:
            raise InvalidParameterError(req.context["data"])

        username = req_data["username"]
        password = req_data["password"]
        try:
            User.get(User.username == username)
            self.on_error(resp, ERR_USER_ALREADY_EXISTS)
        except User.DoesNotExist:
            kwargs = {"username": username, "password": password}
            if "nickname" in req_data.keys():
                kwargs["nickname"] = req_data["nickname"]
            create_user(**kwargs)
            self.on_success(resp)


class UsersItem(BaseResource):
    """
    Handle for endpoint: /v1/users/{user_id}
    """

    @authenticated(verifyadmin=True)
    def on_get(self, req, resp, user_id):
        try:
            user = User.get(User.id == user_id)
            self.on_success(resp, user.model_to_dict())
        except User.DoesNotExist:
            raise UserNotExistsError()

    @authenticated(verifyadmin=True)
    def on_delete(self, req, resp, user_id):
        try:
            user = User.get(User.id == user_id)
            user.delete_instance()
            self.on_success(resp, {"id": user_id})
        except User.DoesNotExist:
            raise UserNotExistsError()
        except Exception as e:
            _api_logger.info("UsersItem on_delete ERROR: %s" % e)
            self.on_error(resp, ERR_UNKNOWN)


class UsersSelf(BaseResource):
    """
    Handle for endpoint: /v1/users/self
    """
    @authenticated()
    def on_get(self, req, resp):
        self.on_success(resp, self.current_user.model_to_dict())


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
        if not req_data:
            raise InvalidParameterError(req.context["data"])

        username = req_data["username"]
        password = req_data["password"]
        try:
            user = User.get(User.username == username)
            _api_logger.debug(f"password: {password}, {user.password}")
            if not user.password.check_password(password):
                raise UnauthorizedError()

            payload = {
                "id": user.id,
                "username": user.username,
                "exp": datetime.utcnow()
            }
            jwtoken = jwt.encode(
                payload, CONFIG.SECRET_KEY, algorithm=CONFIG.JWT_ALGORITHM).decode("utf8")
            self.on_success(resp, {"username": username, "jwtoken": jwtoken})
        except User.DoesNotExist:
            raise UnauthorizedError()

    def process_modify(self, req, resp):
        pass

    def process_logout(self, req, resp):
        pass
