"""初始化环境变量
"""
import os

from dotenv import load_dotenv


# init
dotenv_path = os.path.join(os.path.dirname(__file__), '.env')
if os.path.exists(dotenv_path):
    load_dotenv(dotenv_path)
from hyuga.lib.option import init
init("production")
