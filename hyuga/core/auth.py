import functools

import jwt
import redis

from hyuga.core import errors, log
from hyuga.lib.option import CONFIG
from hyuga.models.user import User, UserToken


def authenticated(verifyadmin=False):
    """Verify that the user is logged in
    :param verifyadmin: verify the current_user is an administrator (if true)
    """
    def decorate(func):
        @functools.wraps(func)
        def wrapper(self, req, resp, *args, **kwargs):
            self.current_user = None
            _token = req.get_header("JWToken", default=None)
            if not _token:
                raise errors.UnauthorizedError('JWToken Not Exists')
            user_token = UserToken.objects.filter(token=_token).first()
            if not user_token:
                raise errors.UnauthorizedError('Expired Signature')
            # decode jwt
            try:
                user_id = jwt.decode(
                    jwt=user_token.jwt,
                    key=CONFIG.SECRET_KEY,
                    algorithms=CONFIG.JWT_ALGORITHM)["id"]
            except jwt.InvalidSignatureError:
                raise errors.UnauthorizedError('Invalid Signature')
            except jwt.DecodeError:
                raise errors.UnauthorizedError('JWToken Decode Error')
            # get authenticated user
            try:
                self.current_user = User.objects.get_by_id(user_id)
            except redis.exceptions.ConnectionError:
                raise errors.DatabaseError(errors.ERR_DATABASE_CONNECTION)
            except Exception as e:
                log._api_logger.info("get authenticated user ERROR: %s" % e)
                self.on_error(resp, errors.ERR_UNKNOWN)
            else:
                if not self.current_user:
                    raise errors.UserNotExistsError()
                if verifyadmin and (not self.current_user.administrator):
                    raise errors.NotSupportedError(
                        method=req.method, url=req.path)
                func(self, req, resp, *args, **kwargs)
        return wrapper
    return decorate
