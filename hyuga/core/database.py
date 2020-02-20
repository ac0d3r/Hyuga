import os
from peewee import MySQLDatabase

import redisco
from hyuga.lib.option import CONFIG
from hyuga.core.log import logger


def db_setup():
    global database
    logger.debug('Set up to Mysql database..')
    database = MySQLDatabase(CONFIG.MYSQL_DB, user=CONFIG.MYSQL_USER,
                             password=CONFIG.MYSQL_PASSWORD, host=CONFIG.MYSQL_SERVER, port=int(CONFIG.MYSQL_PROT))
    logger.debug('Set up to Redis database..')
    redisco.connection_setup(host=CONFIG.REDIS_SERVER, port=int(
        CONFIG.REDIS_PROT), db=int(CONFIG.REDIS_DB))


db_setup()
