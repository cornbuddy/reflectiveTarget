from constants import PASSWORD, USERNAME, URL


def test_user_should_be_able_to_logout(authorized_user):
    authorized_user.logout()
    assert authorized_user.driver.current_url == f"{URL}/"


def test_user_should_be_able_to_login(authorized_user):
    authorized_user.login(USERNAME, PASSWORD)
    assert authorized_user.driver.current_url == f"{URL}/"


def test_password_should_be_validated_upon_signup(anonymous_user):
    raise RuntimeError("not implemented")


def test_username_should_be_validated_upon_signup(anonymous_user):
    raise RuntimeError("not implemented")
