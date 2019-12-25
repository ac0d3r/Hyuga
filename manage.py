#!/usr/bin/env python3
import click

import init

BANNER = r"""
 __                            
|  |--.--.--.--.--.-----.---.-.
|     |  |  |  |  |  _  |  _  |
|__|__|___  |_____|___  |___._|
      |_____|     |_____|      
"""


@click.group()
def main():
    click.echo(BANNER)


@main.command()
@click.option('--username', 'username', prompt=True)
@click.option('--password', 'password', prompt=True, hide_input=True,
              confirmation_prompt=True)
def createsuperuser(username, password):
    from hyuga.models.user import User
    click.echo("[CMD] create super user...")
    User.create(
        username=username,
        password=password,
        identify='admin.hyuga.co',
        token='',
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
def rundns():
    from hyuga.dns import main
    main()


@main.command()
def runapi():
    try:
        from wsgiref import simple_server
        import wsgi
        host = ("127.0.0.1", 5000)
        click.echo("[CMD] Test Listen %s:%s ..." % (host[0], host[1]))
        httpd = simple_server.make_server(
            host[0], host[1], wsgi.app)
        httpd.serve_forever()
    except KeyboardInterrupt:
        click.echo("Bye..")
        httpd.server_close()


if __name__ == "__main__":
    main()
