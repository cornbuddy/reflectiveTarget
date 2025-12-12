from os import environ


URL = f"http://localhost:{environ.get('PORT', 8080)}"
SESSION_TOKEN = "session-token"

USERNAME = "username"
PASSWORD = "p@ssword1"
