# coding: utf-8
import os


class BaseConfig:
    BRAND_NAME = "Hyuga"
    # JWT
    SECRET_KEY = os.getenv("SECRET_KEY", "a secret string")
    JWT_ALGORITHM = "HS256"
    JWT_EXPIRE = 7 * 24 * 3600
    # records
    RECORDS_EXPIRE = 6 * 3600
    RECORDS_MAX_NUM = 100
    # MySQL
    MYSQL_SERVER = os.getenv("MYSQL_SERVER", "localhost")
    MYSQL_PROT = os.getenv("MYSQL_PROT", 3306)
    MYSQL_USER = os.getenv("MYSQL_USER", "root")
    MYSQL_PASSWORD = os.getenv("MYSQL_PASSWORD", "12345678")
    MYSQL_DB = os.getenv("MYSQL_DB", "Hyuga")
    # redis
    REDIS_SERVER = os.getenv("REDIS_SERVER", "localhost")
    REDIS_PROT = os.getenv("REDIS_PROT", 6379)
    REDIS_DB = os.getenv("REDIS_DB", 0)
    # dns
    DOMAIN = os.getenv("DOMAIN", "hyuga.io")
    API_DOMAIN = os.getenv("API_DOMAIN", "api."+DOMAIN)
    # NS域名
    NS1_DOMAIN = os.getenv("NS1_DOMAIN", "ns1.a.com")
    NS2_DOMAIN = os.getenv("NS2_DOMAIN", "ns2.a.com")
    # 服务器外网地址
    SERVER_IP = os.getenv("SERVER_IP", "1.1.1.1")
    # dns 保留域名
    DNS_IGNORE_DOMAINS = (
        DOMAIN,
        API_DOMAIN,
        "admin." + DOMAIN
    )
    # 允许的请求的 http method
    ALLOW_METHODS = ("POST", "GET", "DELETE", "PUT")
    # http record 数据长度限制
    DATA_MAX_LENGHT = 512


class DevelopmentConfig(BaseConfig):
    pass


class ProductionConfig(BaseConfig):
    pass


class TestingConfig(BaseConfig):
    TESTING = True


config = {
    "development": DevelopmentConfig,
    "production": ProductionConfig,
    "testing": TestingConfig
}
