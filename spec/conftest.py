import logging
from base64 import b64encode

import pytest
from playwright.sync_api import Page, sync_playwright

from dsl.dsl import DSL
from constants import BROWSER, DEBUG, URL, USERNAME, PASSWORD

log = logging.getLogger(__name__)


@pytest.fixture(scope="function")
def anon(page: Page):
    """starts anonymous user session"""
    yield DSL(page, URL)

    page.close()


@pytest.fixture(scope="function")
def user(page: Page):
    """starts authorized user session"""
    dsl = DSL(page, URL)
    dsl.login(USERNAME, PASSWORD)
    yield dsl

    page.close()


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
        page = browser.new_context().new_page()
        dsl = DSL(page, URL)
        dsl.signup(USERNAME, PASSWORD)
        log.info("created user `%s`", USERNAME)
        browser.close()
