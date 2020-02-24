import falcon

from hyuga.api.common import base
from hyuga.api.v1 import records, users
from hyuga.core.errors import AppError
from hyuga.middleware import GlobalFilter, HandleCORS


def init_routes(app):
    # index
    app.add_route('/', base.BaseResource())
    # users
    app.add_route('/v1/users', users.Users())
    app.add_route('/v1/users/{user_id:int}', users.UsersItem())
    # users self
    app.add_route('/v1/users/self', users.UsersSelf())
    app.add_route('/v1/users/self/login', users.UsersSelfOperation())
    app.add_route('/v1/users/self/logout', users.UsersSelfOperation())
    # records
    app.add_route('/v1/records', records.Records())
    app.add_route('/v1/users/self/records', records.UsersSelfRecords())

    app.add_error_handler(AppError, AppError.handle)


def create_app(testing=False):
    middlewares = [HandleCORS()]
    if testing is False:
        middlewares.append(GlobalFilter())
    app = falcon.App(middleware=middlewares)
    init_routes(app)
    return app
