import json
from os import environ


ANIMATION_DURATION_SECS = 1
DEBUG = json.loads(environ.get("DEBUG", "false").lower())
BROWSER = environ.get("BROWSER", "firefox")

URL = f"http://localhost:{environ.get('PORT', 8080)}"
SESSION_TOKEN = "session-token"

USERNAME = "username"
PASSWORD = "p@ssword1"
