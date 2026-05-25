import pytest
from playwright.sync_api import Page, expect

from constants import PASSWORD, USERNAME, URL
from pages.pages import LoginPage, SignupPage


def test_user_should_be_able_to_logout(authorized_page: Page):
    login_page = LoginPage(authorized_page)
    before_logout = login_page.page.context.cookies()
    login_page.logout()
    after_logout = login_page.page.context.cookies()
    assert before_logout != after_logout
    assert login_page.page.url == f"{URL}/"


@pytest.mark.parametrize("username,message", [
    (USERNAME, "user already exists"),
])
def test_username_should_be_validated_upon_signup(
        signup_page: SignupPage, username: str, message: str,
):
    signup_page.signup(username, PASSWORD)
    expect(signup_page.form).to_contain_text(message)


@pytest.mark.parametrize("password,message", [
    ("kek", "password should be at least 8 characters long"),
    ("kekekekeke", "password should contain at least one digit"),
    ("kekekekek1", "password should contain at least one special character"),
])
def test_password_should_be_validated_upon_signup(
        signup_page: SignupPage, password: str, message: str,
):
    signup_page.signup(USERNAME, password)
    assert message in signup_page.page.content()


@pytest.mark.parametrize("confirmation,message", [
    ("kek", "passwords should match"),
])
def test_confirmation_should_be_validated_upon_signup(
        signup_page: SignupPage, confirmation: str, message: str,
):
    signup_page.signup(USERNAME, PASSWORD, confirmation)
    assert message in signup_page.page.content()
