import pytest

from constants import PASSWORD, USERNAME, URL


def test_should_signup(dsl):
    driver = dsl.signup(USERNAME, PASSWORD)
    cookie = driver.get_cookie("session-token")
    assert driver.current_url == f"{URL}/"
    assert cookie is not None


@pytest.mark.order(after="test_should_signup")
def test_should_login(dsl):
    cookie = dsl.driver.get_cookie("session-token")
    assert cookie is None

    driver = dsl.login(USERNAME, PASSWORD)
    cookie = driver.get_cookie("session-token")
    assert cookie is not None
    assert driver.current_url == f"{URL}/"
