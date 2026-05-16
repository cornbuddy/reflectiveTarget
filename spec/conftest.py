import logging
from base64 import b64encode

import pytest
from playwright.sync_api import Page, sync_playwright

from pages.base import Layout
from pages.pages import SignupPage, LoginPage
from constants import BROWSER, DEBUG, USERNAME, PASSWORD

log = logging.getLogger(__name__)


@pytest.fixture
def layout(page: Page) -> Layout:
    layout = Layout(page)
    yield layout.page
    return page.close()


@pytest.fixture
def user(page: Page) -> Page:
    """returns authorized session for default user"""
    login = LoginPage(page)
    login.login(USERNAME, PASSWORD)
    yield login.page
    return login.page.close()


@pytest.hookimpl(hookwrapper=True)
def pytest_runtest_makereport(item):
    """makes screenshots for each test case"""
    outcome = yield
    test_report = outcome.get_result()
    if test_report.when == "call":
        page = item.funcargs.get("page")
        screenshot = b64encode(page.screenshot()).decode("utf-8")
        pytest_html = item.config.pluginmanager.getplugin("html")
        extras = getattr(test_report, "extras", [])
        extras.append(pytest_html.extras.image(screenshot))
        test_report.extras = extras


def pytest_sessionstart():
    """called before the first test runs"""
    log.info("creating test user...")
    with sync_playwright() as playwright:
        headless = not DEBUG
        browser = getattr(playwright, BROWSER).launch(headless=headless)
        page = SignupPage(browser.new_context().new_page())
        page.signup(USERNAME, PASSWORD)
        log.info("created user `%s`", USERNAME)
        browser.close()
