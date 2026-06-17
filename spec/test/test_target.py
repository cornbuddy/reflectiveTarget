import pytest
from playwright.sync_api import expect

from constants import URL
from pages.base import LOCATORS
from pages.pages import NewTargetPage, TargetsListPage, UpdateTargetPage

TARGET_NAME = "totally unique target"


def test_has_proper_components(new_target_page: NewTargetPage):
    locators = [
        (LOCATORS["form"]),
        (LOCATORS["canvas"]),
        (LOCATORS["submit"]),
        (LOCATORS["add_question"]),
        (LOCATORS["name"]),
        (LOCATORS["question"](0)),
    ]
    for locator in locators:
        expect(new_target_page.page.locator(locator)).to_have_count(1)


def test_canvas_is_properly_sized(new_target_page: NewTargetPage):
    size = new_target_page.page.locator(LOCATORS["canvas"]).bounding_box()
    assert abs(size["height"] - size["width"]) < 1


def test_target_can_have_many_questions(new_target_page: NewTargetPage):
    name, questions = "kek?", ["kek1", "kek2", "kek3"]
    target_id = new_target_page.create_target(name, questions)
    update_page = UpdateTargetPage(new_target_page.page, target_id)
    inputs = update_page.page.locator(LOCATORS["questions"])
    expect(inputs).to_have_count(len(questions))
    expect(update_page.name_input).to_have_attribute("value", name)
    for i, qstn_input in enumerate(inputs.all()):
        expect(qstn_input).to_have_attribute("value", questions[i])


def test_user_should_be_able_to_create_target(new_target_page: NewTargetPage):
    new_target_page.create_target(TARGET_NAME, ["kek"])
    page = new_target_page.page
    expect(page).to_have_url(f"{URL}/targets")
    expect(page.get_by_text(TARGET_NAME)).to_have_count(1)


@pytest.mark.order(after=test_user_should_be_able_to_create_target.__name__)
def test_user_should_be_able_to_edit_its_target(
        targets_list_page: TargetsListPage,
):
    page = targets_list_page.page
    target_id = targets_list_page.get_id_by_target_name(TARGET_NAME)
    update_page = UpdateTargetPage(page, target_id)
    want_question = "new question value"
    update_page.update_target(TARGET_NAME, [want_question])
    page.get_by_text(TARGET_NAME).click()
    question_input = page.locator(LOCATORS["question"](0))
    expect(question_input).to_have_attribute("value", want_question)


@pytest.mark.skip
def test_user_should_not_be_able_to_edit_another_users_target():
    raise NotImplementedError


@pytest.mark.skip
@pytest.mark.order(after=test_user_should_be_able_to_create_target.__name__)
def test_anonymous_should_be_able_to_shoot_target():
    raise NotImplementedError
