from flask import Flask

app = Flask(__name__)

from jam import jam # noqa: F401, E402
