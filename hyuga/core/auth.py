import functools

import jwt

from hyuga.core.errors import UnauthorizedError, UserNotExistsError, NotSupportedError
from hyuga.lib.option import CONFIG
from hyuga.models.user import User


def authenticated(verifyadmin=False):
    """Verify that the user is logged in
    :param verifyadmin: verify the current_user is an administrator (if true)
    """
    def decorate(func):
        @functools.wraps(func)
        def wrapper(self, req, resp, *args, **kwargs):
            self.current_user = None
            jwt_token = req.get_header("JWToken", default=None)
            if not jwt_token:
                raise UnauthorizedError('JWToken Not Exists')
            try:
                user_id = jwt.decode(
                    jwt_token, CONFIG.SECRET_KEY, algorithms=CONFIG.JWT_ALGORITHM, leeway=CONFIG.JWT_EXPIRE, options={"verify_exp": True})["id"]
                try:
                    self.current_user = User.get(User.id == user_id)
                except User.DoesNotExist:
                    raise UserNotExistsError()
                if verifyadmin and (not self.current_user.administrator):
                    raise NotSupportedError(
                        method=req.method, url=req.path)
                func(self, req, resp, *args, **kwargs)
            except jwt.ExpiredSignatureError:
                raise UnauthorizedError('Expired Signature')
            except jwt.InvalidSignatureError:
                raise UnauthorizedError('Invalid Signature')
            except jwt.DecodeError:
                raise UnauthorizedError('JWToken Decode Error')
        return wrapper
    return decorate
