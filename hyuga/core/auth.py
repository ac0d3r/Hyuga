import functools

import jwt

from hyuga.core.errors import UnauthorizedError, UserNotExistsError
from hyuga.lib.option import CONFIG
from hyuga.models.user import User


def authenticated(method):
    @functools.wraps(method)
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
                method(self, req, resp, *args, **kwargs)
            except User.DoesNotExist:
                raise UserNotExistsError()
        except jwt.ExpiredSignatureError:
            raise UnauthorizedError('Expired Signature')
        except jwt.InvalidSignatureError:
            raise UnauthorizedError('Invalid Signature')
        except jwt.DecodeError:
            raise UnauthorizedError('JWToken Decode Error')
    return wrapper
