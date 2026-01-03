import pytest
from selenium.webdriver.common.by import By


def test_sidebar_should_be_toggled_on_click(anon):
    anon.toggle_navigation()
    assert anon.sidebar.is_displayed()

    anon.toggle_navigation()
    assert not anon.sidebar.is_displayed()


@pytest.mark.parametrize("link", ["Login", "Signup"])
def test_sidebar_should_have_links_for_anonymous_uesr(anon, link):
    anon.ensure_navigation_opened()
    assert anon.sidebar.find_element(By.LINK_TEXT, link) is not None


@pytest.mark.parametrize("link", ["Logout"])
def test_sidebar_should_have_links_for_authorized_user(user, link):
    user.ensure_navigation_opened()
    assert user.sidebar.find_element(By.LINK_TEXT, link) is not None
