from constants import PASSWORD, USERNAME


def test_should_login(authorized_user):
    authorized_user.login(USERNAME, PASSWORD)
    authorized_user.assert_authorized()
