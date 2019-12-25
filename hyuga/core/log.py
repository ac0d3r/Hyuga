# coding: utf-8
import logging
import sys

INFO_FORMAT = '[%(asctime)s] [%(process)d] [%(levelname)s] %(message)s'
DEBUG_FORMAT = INFO_FORMAT + ' [in %(pathname)s:%(lineno)d]'
TIMESTAMP_FORMAT = '%Y-%m-%d %H:%M:%S %z'

logger: logging.Logger = logging.getLogger('Hyuga')
logger.propagate = False


def setup_hyuga_logger(level=None, verbose=False, very_verbose=False):
    """设置 Hyuga 全局日志
        Args:
            very_verbose: 启动调试模式
    """
    # 设置 log format
    Format = INFO_FORMAT
    if verbose or very_verbose or (level != "INFO"):
        Format = DEBUG_FORMAT
    # 清除 logging 默认日志输出
    logging.getLogger().handlers.clear()
    formatter = logging.Formatter(Format, TIMESTAMP_FORMAT)
    output_handler = logging.StreamHandler(sys.stdout)
    output_handler.setFormatter(formatter)
    # 清除 handlers
    logger.handlers.clear()
    logger.addHandler(output_handler)

    if level is not None:
        logger.setLevel(level)
    else:
        if very_verbose:
            logger.setLevel(logging.DEBUG)
        elif verbose:
            logger.setLevel(logging.INFO)
        else:
            logger.setLevel(logging.WARNING)
    return logger


class Logger:
    def _process_msg(self, msg):
        return msg

    def debug(self, msg, *args, **kwargs):
        logger.debug(self._process_msg(msg), *args, **kwargs)

    def info(self, msg, *args, **kwargs):
        logger.info(self._process_msg(msg), *args, **kwargs)

    def warning(self, msg, *args, **kwargs):
        logger.warning(self._process_msg(msg), *args, **kwargs)

    def warn(self, msg, *args, **kwargs):
        logger.warning(self._process_msg(msg), *args, **kwargs)

    def error(self, msg, *args, **kwargs):
        logger.error(self._process_msg(msg), *args, **kwargs)

    def exception(self, msg, *args, exc_info=True, **kwargs):
        kwargs['exc_info'] = exc_info
        logger.exception(self._process_msg(msg), *args, **kwargs)

    def critical(self, msg, *args, **kwargs):
        logger.critical(self._process_msg(msg), *args, **kwargs)


class PrefixedLogger(Logger):
    def __init__(self, prefix: str):
        Logger.__init__(self)
        self.__prefix = prefix

    def _process_msg(self, msg):
        return '%s %s' % (self.__prefix, msg)


_api_logger = PrefixedLogger("[API]")
_dns_logger = PrefixedLogger("[DNS]")
