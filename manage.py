#!/usr/bin/env python3
from wsgiref import simple_server

import click

from hyuga.app import create_app
from hyuga.core.database import database
from hyuga.dns import dnsserver
from hyuga.models.user import User

BANNER = r"""
 __                            
|  |--.--.--.--.--.-----.---.-.
|     |  |  |  |  |  _  |  _  |
|__|__|___  |_____|___  |___._|
      |_____|     |_____|      
"""


@click.group()
def manage():
    click.echo(BANNER)


@manage.command()
@click.option("--username", "username", prompt=True)
@click.option("--password", "password", prompt=True, hide_input=True,
              confirmation_prompt=True)
def createsuperuser(username, password):
    """Create superuser.
    """
    click.echo("[CMD] create super user...")
    User.create(
        username=username,
        password=password,
        identify="admin.hyuga.co",
        token="",
        administrator=True
    )
    click.echo("[CMD] create super user success...")


@manage.command()
def createtables():
    """Create tables.
    """
    click.echo("[CMD] create tables...")
    with database:
        database.create_tables([User], safe=True)
    click.echo("[CMD] create tables success...")


@manage.command()
@click.option("--host", "host", default="127.0.0.1")
@click.option("--port", "port", default=8080)
def runweb(host, port):
    """Runs the application on a local development server.
    """
    try:
        app = create_app(testing=True)
        click.echo("[TESTING] %s" % ("="*50))
        click.echo("[TESTING] API Server run on: http://%s:%s" %
                   (host, port))
        click.echo("[TESTING] %s" % ("="*50))
        httpd = simple_server.make_server(host, int(port), app)
        httpd.serve_forever()
    except KeyboardInterrupt:
        click.echo("\nBye👋")
        httpd.shutdown()


@manage.command()
def rundns():
    """Runs the simple dns server.
    """
    dnsserver()


if __name__ == "__main__":
    manage()
