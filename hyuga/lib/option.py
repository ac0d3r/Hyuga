# coding: utf-8
import os

from hyuga.core.log import setup_hyuga_logger

from .settings import config, BaseConfig

CONFIG_LOG_LEVEL = {
    "production": "INFO",
    "development": "DEBUG",
    "testing": "DEBUG",
}
LOG_LEVEL: str = None
CONFIG: BaseConfig = None


def init(config_env=None):
    global LOG_LEVEL
    global CONFIG
    if config_env is None:
        config_env = os.getenv('APP_ENV', 'development')
    LOG_LEVEL = CONFIG_LOG_LEVEL[config_env]
    setup_hyuga_logger(LOG_LEVEL)
    CONFIG = config[config_env]
