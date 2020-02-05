#!/usr/bin/env python3
from wsgiref import simple_server

import click

from hyuga.lib.option import init

BANNER = r"""
 __                            
|  |--.--.--.--.--.-----.---.-.
|     |  |  |  |  |  _  |  _  |
|__|__|___  |_____|___  |___._|
      |_____|     |_____|      
"""
init()


@click.group()
def main():
    click.echo(BANNER)


@main.command()
@click.option("--username", "username", prompt=True)
@click.option("--password", "password", prompt=True, hide_input=True,
              confirmation_prompt=True)
def createsuperuser(username, password):
    from hyuga.models.user import User
    click.echo("[CMD] create super user...")
    User.create(
        username=username,
        password=password,
        identify="admin.hyuga.co",
        token="",
        administrator=True
    )
    click.echo("[CMD] create super user success...")


@main.command()
def createtables():
    from hyuga.database import database
    from hyuga.models.user import User
    click.echo("[CMD] create tables...")
    with database:
        database.create_tables([User], safe=True)
    click.echo("[CMD] create tables success...")


@main.command()
@click.option("--host", "host", default="127.0.0.1")
@click.option("--port", "port", default=8080)
def runweb(host, port):
    """Runs the application on a local development server.
    """
    from hyuga.app import create
    try:
        app = create(testing=True)
        click.echo("[TESTING] %s" % ("="*50))
        click.echo("[TESTING] API Server run on: http://%s:%s" %
                   (host, port))
        click.echo("[TESTING] %s" % ("="*50))
        httpd = simple_server.make_server(host, int(port), app)
        httpd.serve_forever()
    except KeyboardInterrupt:
        click.echo("\nBye👋")
        httpd.shutdown()


@main.command()
def rundns():
    from hyuga.dns import main
    main()


if __name__ == "__main__":
    main()
