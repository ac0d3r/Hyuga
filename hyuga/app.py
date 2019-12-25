import falcon

from hyuga.api.common import base
from hyuga.api.v1 import records, users
from hyuga.core.errors import AppError
from hyuga.core.log import logger
from hyuga.middleware import GlobalFilter, HandleCORS, PeeweeConnection


class App(falcon.API):
    def __init__(self, *args, **kwargs):
        super(App, self).__init__(*args, **kwargs)
        logger.info('API Server is starting')
        # index
        self.add_route('/', base.BaseResource())
        # users
        self.add_route('/v1/users', users.Users())
        self.add_route('/v1/users/{user_id:int}', users.UsersItem())
        # users self
        self.add_route('/v1/users/self', users.UsersSelf())
        self.add_route('/v1/users/self/login', users.UsersSelfOperation())
        # records
        self.add_route('/v1/records', records.Records())
        self.add_route('/v1/users/self/records', records.UsersSelfRecords())

        self.add_error_handler(AppError, AppError.handle)


def create(testing=False):
    middlewares = [HandleCORS(), PeeweeConnection()]
    if testing is False:
        middlewares.append(GlobalFilter())
    return App(middleware=middlewares)
