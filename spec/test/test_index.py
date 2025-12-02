import pytest
from selenium.webdriver.common.by import By


def test_sidebar_should_be_toggled_on_click(anonymous_user):
    anonymous_user.toggle_navigation()
    assert anonymous_user.sidebar.is_displayed()

    anonymous_user.toggle_navigation()
    assert not anonymous_user.sidebar.is_displayed()


@pytest.mark.parametrize("link", ["Login", "Signup"])
def test_sidebar_should_have_links_for_anonymous_uesr(anonymous_user, link):
    anonymous_user.ensure_navigation_opened()
    sidebar = anonymous_user.sidebar
    assert sidebar.find_element(By.LINK_TEXT, link) is not None


@pytest.mark.parametrize("link", ["Logout"])
def test_sidebar_should_have_links_for_authorized_user(authorized_user, link):
    authorized_user.ensure_navigation_opened()
    sidebar = authorized_user.sidebar
    assert sidebar.find_element(By.LINK_TEXT, link) is not None
