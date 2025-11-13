import requests

from constants import URL


def test_should_expose_health_endpoint():
    url = f"{URL}/api/health"
    resp = requests.get(url, timeout=1)
    body = resp.json()
    assert resp.status_code == 200
    assert body["connected"]
