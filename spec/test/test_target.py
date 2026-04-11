import pytest
from selenium.webdriver.common.by import By


@pytest.mark.parametrize("by,value", [
    (By.TAG_NAME, "form"),
    (By.TAG_NAME, "canvas"),
    (By.XPATH, "//button[@type='submit']"),
    (By.XPATH, "//button[text()='Add question']"),
    (By.XPATH, "//input[@type='text'][@name='name']"),
    (By.XPATH, "//input[@type='text'][starts-with(@name, 'question_')]"),
])
def test_has_proper_components(user, by, value):
    user.go_to_new_target()
    elems = user.driver.find_elements(by, value)
    assert len(elems) == 1


def test_canvas_is_properly_sized(user):
    canvas = user.driver.find_element(By.TAG_NAME, "canvas")
    assert canvas.size["height"] == canvas.size["width"]


def test_question_can_be_added(user):
    user.go_to_new_target()
    driver = user.driver
    add_btn = driver.find_element(By.XPATH, "//button[text()='Add question']")
    add_btn.click()
    inputs_xpath = "//input[@type='text'][starts-with(@name, 'question_')]"
    inputs = driver.find_elements(By.XPATH, inputs_xpath)
    assert len(inputs) == 2


@pytest.mark.skip
def test_user_should_be_able_to_create_target(user):
    raise RuntimeError("not implemented")


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
