import os
from peewee import MySQLDatabase

from hyuga.lib.option import CONFIG
from hyuga.core.log import logger


def get_mysql_db():
    logger.debug('Connecting to Mysql database..')
    return MySQLDatabase(CONFIG.MYSQL_DB, user=CONFIG.MYSQL_USER, password=CONFIG.MYSQL_PASSWORD, host=CONFIG.MYSQL_SERVER, port=int(CONFIG.MYSQL_PROT))


database = get_mysql_db()
