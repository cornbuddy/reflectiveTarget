import pytest

from constants import PASSWORD, USERNAME, URL


def test_user_should_be_able_to_logout(authorized_user):
    authorized_user.logout()
    assert authorized_user.driver.current_url == f"{URL}/"


def test_user_should_be_able_to_login(authorized_user):
    authorized_user.login(USERNAME, PASSWORD)
    assert authorized_user.driver.current_url == f"{URL}/"


@pytest.mark.parametrize("username,message", [
    (USERNAME, "user already exists"),
])
def test_username_should_be_validated_upon_signup(
        anonymous_user, username, message,
):
    anonymous_user.signup(username, PASSWORD)
    assert message in anonymous_user.driver.page_source


@pytest.mark.parametrize("password,message", [
    ("kek", "password should be at least 8 characters long"),
    ("kekekekeke", "password should contain at least one digit"),
    ("kekekekek1", "password should contain at least one special character"),
])
def test_password_should_be_validated_upon_signup(
        anonymous_user, password, message,
):
    anonymous_user.signup(USERNAME, password)
    assert message in anonymous_user.driver.page_source


@pytest.mark.parametrize("confirmation,message", [
    ("kek", "passwords should match"),
])
def test_confirmation_should_be_validated_upon_signup(
        anonymous_user, confirmation, message,
):
    anonymous_user.signup(USERNAME, PASSWORD, confirmation)
    assert message in anonymous_user.driver.page_source
