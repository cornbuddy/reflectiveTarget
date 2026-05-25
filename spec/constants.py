import json
from os import environ as env


ANIMATION_DURATION_SECS = 1
DEBUG = json.loads(env.get("DEBUG", "false").lower())
BROWSER = env.get("BROWSER", "firefox")
SLOW_MO_MS = float(env["SLOW_MO_MS"]) if env.get("SLOW_MO_MS") else 0

URL = f"http://localhost:{env.get('PORT', 8080)}"
SESSION_TOKEN = "session-token"

USERNAME = "username"
PASSWORD = "p@ssword1"
