from hyuga.lib.option import init
init("production")
from hyuga.app import create
app = create()
