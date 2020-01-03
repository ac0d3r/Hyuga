#!/usr/bin/env python3
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


if __name__ == "__main__":
    main()
