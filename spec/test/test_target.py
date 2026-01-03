import pytest
from selenium.webdriver.common.by import By


def test_user_should_be_able_to_create_target(user):
    user.ensure_navigation_opened()
    user.driver.find_element(By.LINK_TEXT, "Targets").click()
    user.driver.find_element(By.LINK_TEXT, "New").click()
    target_form = user.driver.find_element(By.TAG_NAME, "form")
    assert target_form is not None


@pytest.mark.order(after=test_user_should_be_able_to_create_target.__name__)
@pytest.mark.skip
def test_anonymous_should_be_able_to_shoot_target(anon):
    raise RuntimeError("not implemented")


@pytest.mark.order(after=test_user_should_be_able_to_create_target.__name__)
@pytest.mark.skip
def test_user_should_be_able_to_edit_its_target(user):
    raise RuntimeError("not implemented")


@pytest.mark.skip
def test_user_should_not_be_able_to_edit_another_users_target(user):
    raise RuntimeError("not implemented")
