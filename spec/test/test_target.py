import pytest

from constants import URL

TARGET_NAME = "totally unique target"
TARGET_INPUTS = {
    "add_question": "//button[text()='Add question']",
    "name": "//input[@name='name']",
    "questions": "//input[starts-with(@name, 'question_')]",
    "question_0": "//input[@name='question_0']",
}


@pytest.mark.parametrize("locator", [
    ("form"),
    ("canvas"),
    ("//button[@type='submit']"),
    (TARGET_INPUTS["add_question"]),
    (TARGET_INPUTS["name"]),
    (TARGET_INPUTS["question_0"]),
])
def test_has_proper_components(user, locator):
    user.go_to_new_target()
    elems = user.page.locator(locator)
    assert elems.count() == 1


def test_canvas_is_properly_sized(user):
    user.go_to_new_target()
    size = user.page.locator("canvas").bounding_box()
    assert abs(size["height"] - size["width"]) < 1


def test_question_can_be_added(user):
    user.go_to_new_target()
    page = user.page
    page.locator(TARGET_INPUTS["add_question"]).click()
    assert page.locator(TARGET_INPUTS["questions"]).count() == 2


def test_user_should_be_able_to_create_target(user):
    user.go_to_new_target()
    user.page.locator(TARGET_INPUTS["name"]).fill(TARGET_NAME)
    user.page.locator(TARGET_INPUTS["question_0"]).fill("kek")
    user.submit.click()
    want_url = f"{URL}/targets"
    user.page.wait_for_url(want_url)
    assert user.page.url == want_url
    assert user.page.get_by_text(TARGET_NAME).count() == 1


@pytest.mark.order(after=test_user_should_be_able_to_create_target.__name__)
def test_user_should_be_able_to_edit_its_target(user):
    user.page.goto(f"{URL}/targets")
    user.page.get_by_text(TARGET_NAME).click()
    want_question = "new question value"
    question_input = user.page.locator(TARGET_INPUTS["question_0"])
    question_input.clear()
    question_input.fill(want_question)
    user.submit.click()
    user.page.get_by_text(TARGET_NAME).click()
    got_question = question_input.get_attribute("value")
    assert got_question == want_question


@pytest.mark.skip
def test_user_should_not_be_able_to_edit_another_users_target(user):
    raise RuntimeError("not implemented")


@pytest.mark.skip
@pytest.mark.order(after=test_user_should_be_able_to_create_target.__name__)
def test_anonymous_should_be_able_to_shoot_target(anon):
    raise RuntimeError("not implemented")
