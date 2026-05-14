from time import sleep

import pytest

from constants import ANIMATION_DURATION_SECS


def test_sidebar_should_be_toggled_on_click(anon):
    anon.sidebar_toggler.click()
    sleep(ANIMATION_DURATION_SECS)
    assert anon.sidebar.is_visible()

    anon.sidebar_toggler.click()
    sleep(ANIMATION_DURATION_SECS)
    assert not anon.sidebar.is_visible()


@pytest.mark.parametrize("link", ["Login", "Signup"])
def test_sidebar_should_have_links_for_anonymous_uesr(anon, link):
    anon.ensure_navigation_opened()
    assert anon.sidebar.get_by_text(link) is not None


@pytest.mark.parametrize("link", ["Logout", "Targets"])
def test_sidebar_should_have_links_for_authorized_user(user, link):
    user.ensure_navigation_opened()
    assert user.sidebar.get_by_text(link) is not None
