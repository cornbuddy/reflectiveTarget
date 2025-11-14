import pytest

from constants import PASSWORD, USERNAME, URL


def test_should_signup(dsl):
    driver = dsl.signup(USERNAME, PASSWORD)
    session = driver.get_cookie("session-token")
    assert driver.current_url == f"{URL}/"
    assert session is not None


@pytest.mark.order(after="test_should_signup")
def test_should_login(dsl):
    session = dsl.driver.get_cookie("session-token")
    assert session is None

    driver = dsl.login(USERNAME, PASSWORD)
    session = driver.get_cookie("session-token")
    assert session is not None
    assert driver.current_url == f"{URL}/"
