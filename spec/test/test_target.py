import pytest
from selenium.webdriver.common.by import By

from constants import URL

TARGET_NAME = "totally unique target"
TARGET_INPUTS = {
    "name": "//input[@name='name']",
    "question_0": "//input[@name='question_0']",
    "add_question": "//button[text()='Add question']",
}


@pytest.mark.parametrize("by,value", [
    (By.TAG_NAME, "form"),
    (By.TAG_NAME, "canvas"),
    (By.XPATH, "//button[@type='submit']"),
    (By.XPATH, TARGET_INPUTS["add_question"]),
    (By.XPATH, TARGET_INPUTS["name"]),
    (By.XPATH, TARGET_INPUTS["question_0"]),
])
def test_has_proper_components(user, by, value):
    user.go_to_new_target()
    elems = user.driver.find_elements(by, value)
    assert len(elems) == 1


def test_canvas_is_properly_sized(user):
    user.go_to_new_target()
    canvas = user.driver.find_element(By.TAG_NAME, "canvas")
    assert canvas.size["height"] == canvas.size["width"]


def test_question_can_be_added(user):
    user.go_to_new_target()
    driver = user.driver
    add_btn = driver.find_element(By.XPATH, TARGET_INPUTS["add_question"])
    add_btn.click()
    inputs_xpath = "//input[starts-with(@name, 'question_')]"
    inputs = driver.find_elements(By.XPATH, inputs_xpath)
    assert len(inputs) == 2


def test_user_should_be_able_to_create_target(user):
    user.go_to_new_target()
    driver = user.driver
    driver.find_element(By.XPATH, TARGET_INPUTS["name"]).send_keys(TARGET_NAME)
    driver.find_element(By.XPATH, TARGET_INPUTS["question_0"]).send_keys("kek")
    user.submit.click()
    assert driver.current_url == f"{URL}/targets"
    assert len(driver.find_elements(By.LINK_TEXT, TARGET_NAME)) == 1


@pytest.mark.order(after=test_user_should_be_able_to_create_target.__name__)
def test_user_should_be_able_to_edit_its_target(user):
    driver = user.driver
    driver.get(f"{URL}/targets")
    driver.find_element(By.LINK_TEXT, TARGET_NAME).click()
    want_question = "new question value"
    question = driver.find_element(By.XPATH, TARGET_INPUTS["question_0"])
    question.clear()
    question.send_keys(want_question)
    user.submit.click()

    driver.find_element(By.LINK_TEXT, TARGET_NAME).click()
    got_question = driver.find_element(
        By.XPATH, TARGET_INPUTS["question_0"],
    ).get_attribute("value")
    assert got_question == want_question


@pytest.mark.skip
def test_user_should_not_be_able_to_edit_another_users_target(user):
    raise RuntimeError("not implemented")


@pytest.mark.skip
@pytest.mark.order(after=test_user_should_be_able_to_create_target.__name__)
def test_anonymous_should_be_able_to_shoot_target(anon):
    raise RuntimeError("not implemented")
