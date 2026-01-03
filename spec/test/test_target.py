import pytest


def test_user_should_be_able_to_create_target(authorized_user):
    raise RuntimeError("not implemented")


@pytest.mark.order(after=test_user_should_be_able_to_create_target.__name__)
def test_anonymous_should_be_able_to_shoot_target(anonymous_user):
    raise RuntimeError("not implemented")


@pytest.mark.order(after=test_user_should_be_able_to_create_target.__name__)
def test_user_should_be_able_to_edit_its_target(authorized_user):
    raise RuntimeError("not implemented")


def test_user_should_not_be_able_to_edit_another_users_target(authorized_user):
    raise RuntimeError("not implemented")
