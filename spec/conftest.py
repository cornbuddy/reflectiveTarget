import logging
from base64 import b64encode

import pytest
from playwright.sync_api import Page, sync_playwright

from pages.base import Layout
from pages.pages import SignupPage, LoginPage, NewTargetPage, TargetsListPage
from constants import BROWSER, DEBUG, USERNAME, PASSWORD, URL, SLOW_MO_MS

log = logging.getLogger(__name__)


@pytest.fixture
def authorized_page(page: Page) -> Page:
    """returns authorized session for default user"""
    login = LoginPage(page)
    login.login(USERNAME, PASSWORD)
    yield login.page
    return login.page.close()


@pytest.fixture
def targets_list_page(authorized_page: Page) -> TargetsListPage:
    """returns targets list POM"""
    targets_list = TargetsListPage(authorized_page)
    yield targets_list


@pytest.fixture
def new_target_page(authorized_page: Page) -> NewTargetPage:
    """returns new target POM"""
    new_target = NewTargetPage(authorized_page)
    yield new_target


@pytest.fixture
def signup_page(page: Page) -> SignupPage:
    """returns signup POM"""
    signup = SignupPage(page)
    yield signup
    return signup.page.close()


@pytest.fixture
def layout(page: Page) -> Layout:
    """returns layout POM"""
    layout = Layout(page)
    yield layout
    return layout.page.close()


@pytest.fixture(scope="session")
def browser_context_args(browser_context_args):
    return {
        **browser_context_args,
        "base_url": URL,
    }


@pytest.fixture(scope="session")
def browser_type_launch_args(browser_type_launch_args):
    return {
        **browser_type_launch_args,
        "headless": not DEBUG,
        "slow_mo": SLOW_MO_MS,
    }


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
        browser = getattr(playwright, BROWSER).launch(
            headless=not DEBUG,
            slow_mo=SLOW_MO_MS,
        )
        signup = SignupPage(browser.new_context().new_page())
        signup.signup(USERNAME, PASSWORD)
        log.info("created user `%s`", USERNAME)
        browser.close()
