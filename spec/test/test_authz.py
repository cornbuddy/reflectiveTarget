from constants import PASSWORD, USERNAME, URL, SESSION_TOKEN


def test_should_login(authorized_user):
    driver = authorized_user.login(USERNAME, PASSWORD)
    session = driver.get_cookie(SESSION_TOKEN)
    assert session is not None
    assert driver.current_url == f"{URL}/"
