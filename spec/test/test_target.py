import pytest
from playwright.sync_api import expect

from constants import URL
from pages.base import LOCATORS
from pages.pages import NewTargetPage, TargetsListPage, UpdateTargetPage

TARGET_NAME = "totally unique target"


@pytest.mark.parametrize("locator", [
    (LOCATORS["form"]),
    (LOCATORS["canvas"]),
    (LOCATORS["submit"]),
    (LOCATORS["add_question"]),
    (LOCATORS["name"]),
    (LOCATORS["question"](0)),
])
def test_has_proper_components(new_target_page: NewTargetPage, locator: str):
    expect(new_target_page.page.locator(locator)).to_have_count(1)


def test_canvas_is_properly_sized(new_target_page: NewTargetPage):
    size = new_target_page.page.locator(LOCATORS["canvas"]).bounding_box()
    assert abs(size["height"] - size["width"]) < 1


def test_question_can_be_added(new_target_page: NewTargetPage):
    new_target_page.page.locator(LOCATORS["add_question"]).click()
    expect(
        new_target_page.page.locator(LOCATORS["questions"]),
    ).to_have_count(2)


@pytest.mark.problem
def test_user_should_be_able_to_create_target(new_target_page: NewTargetPage):
    page = new_target_page.page
    new_target_page.create_target(TARGET_NAME, ["kek"])
    expect(page).to_have_url(f"{URL}/targets")
    expect(page.get_by_text(TARGET_NAME)).to_have_count(1)


@pytest.mark.problem
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
    expect(
        page.locator(LOCATORS["question"](0)),
    ).to_have_attribute("value", want_question)


@pytest.mark.skip
def test_user_should_not_be_able_to_edit_another_users_target():
    raise NotImplementedError


@pytest.mark.skip
@pytest.mark.order(after=test_user_should_be_able_to_create_target.__name__)
def test_anonymous_should_be_able_to_shoot_target():
    raise NotImplementedError
