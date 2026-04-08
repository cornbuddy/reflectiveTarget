import pytest
from selenium.webdriver.common.by import By


def test_target_view_has_proer_components(user):
    user.go_to_new_target()
    components = {
        "form": {
            "by": By.TAG_NAME,
            "selector": "form",
        },
        "canvas": {
            "by": By.TAG_NAME,
            "selector": "canvas",
        },
        "submit": {
            "by": By.XPATH,
            "selector": "//button[@type='submit']",
        },
    }
    for _, select in components.items():
        elem = user.driver.find_element(select["by"], select["selector"])
        select["element"] = elem
        assert elem is not None
    size = components["canvas"]["element"].size
    assert size["height"] == size["width"]


def test_user_should_be_able_to_create_target(user):
    user.go_to_new_target()
    assert not user.submit.is_enabled()


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
