import pytest

from constants import PASSWORD, USERNAME, URL


def test_user_should_be_able_to_logout(user):
    before_logout = user.page.context.cookies()
    user.logout()
    after_logout = user.page.context.cookies()
    assert before_logout != after_logout
    assert user.page.url == f"{URL}/"


@pytest.mark.parametrize("username,message", [
    (USERNAME, "user already exists"),
])
def test_username_should_be_validated_upon_signup(anon, username, message):
    anon.signup(username, PASSWORD)
    assert message in anon.page.content()


@pytest.mark.parametrize("password,message", [
    ("kek", "password should be at least 8 characters long"),
    ("kekekekeke", "password should contain at least one digit"),
    ("kekekekek1", "password should contain at least one special character"),
])
def test_password_should_be_validated_upon_signup(anon, password, message):
    anon.signup(USERNAME, password)
    assert message in anon.page.content()


@pytest.mark.parametrize("confirmation,message", [
    ("kek", "passwords should match"),
])
def test_confirmation_should_be_validated_upon_signup(
        anon, confirmation, message,
):
    anon.signup(USERNAME, PASSWORD, confirmation)
    assert message in anon.page.content()
