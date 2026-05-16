from time import sleep

from pages.base import Layout
from constants import ANIMATION_DURATION_SECS


def test_sidebar_should_be_toggled_on_click(layout):
    layout.sidebar_toggler.click()
    sleep(ANIMATION_DURATION_SECS)
    assert layout.sidebar.is_visible()

    layout.sidebar_toggler.click()
    sleep(ANIMATION_DURATION_SECS)
    assert not layout.sidebar.is_visible()


def test_sidebar_should_have_links_for_anonymous_uesr(layout):
    layout.ensure_navigation_opened()
    for link in ["Login", "Signup"]:
        assert layout.sidebar.get_by_text(link) is not None


def test_sidebar_should_have_links_for_authorized_user(user):
    layout = Layout(user)
    layout.ensure_navigation_opened()
    for link in ["Logout", "Targets"]:
        assert layout.sidebar.get_by_text(link) is not None
