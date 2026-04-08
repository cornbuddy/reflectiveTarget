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
    canvas_size = components["canvas"]["element"].size
    assert canvas_size["height"] == canvas_size["width"]
    assert not user.submit.is_enabled()


def test_user_should_be_able_to_create_target(user):
    user.go_to_new_target()
    driver = user.driver
    inputs_xpath = "//input[@type='text'][starts-with(@name, 'question')]"
    inputs = driver.find_elements(By.XPATH, inputs_xpath)
    assert len(inputs) == 0

    add_btn = driver.find_element(By.XPATH, "//button[text()='Add question']")
    assert add_btn.is_enabled()

    add_btn.click()
    inputs = driver.find_elements(By.XPATH, inputs_xpath)
    assert len(inputs) == 1


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
