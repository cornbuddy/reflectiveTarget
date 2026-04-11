from shutil import which

import pytest
from selenium import webdriver
from selenium.webdriver.firefox.options import Options
from selenium.webdriver.firefox.service import Service
from selenium.webdriver.remote.webdriver import WebDriver

from dsl.dsl import DSL
from constants import DEBUG, URL, USERNAME, PASSWORD


@pytest.fixture(scope="session")
def anon():
    """starts anonymous user session"""
    browser = start_browser()
    yield DSL(browser, URL)

    browser.quit()


@pytest.fixture(scope="session")
def user():
    """starts authorized user session"""
    browser = start_browser()
    dsl = DSL(browser, URL)
    dsl.signup(USERNAME, PASSWORD)
    yield dsl

    browser.quit()


@pytest.hookimpl(hookwrapper=True)
def pytest_runtest_makereport(item):
    """makes screenshots for each test case"""
    driver = None
    test_args = item.funcargs
    if test_args.get("user", False):
        driver = test_args["user"].driver
    elif test_args.get("anon", False):
        driver = test_args["anon"].driver
    else:
        raise RuntimeError(f"failed to get driver from args: {test_args}")

    outcome = yield
    test_report = outcome.get_result()
    if test_report.when == "call":
        screenshot = driver.get_screenshot_as_base64()
        pytest_html = item.config.pluginmanager.getplugin("html")
        extras = getattr(test_report, "extra", [])
        extras.append(pytest_html.extras.image(screenshot))
        test_report.extras = extras


def start_browser() -> WebDriver:
    """configures and runs selenium driver"""
    opts = list(filter(None, [
        "--disable-gpu",
        "--no-sandbox",
        "--disable-dev-shm-usage",
        "--headless" if not DEBUG else "",
    ]))
    options = Options()
    for opt in opts:
        options.add_argument(opt)
    path = which("firefox.geckodriver")
    service = Service(executable_path=path)
    browser = webdriver.Firefox(options=options, service=service)
    browser.set_window_size(1920, 1080)
    browser.implicitly_wait(10)
    return browser
