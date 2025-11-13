import json
from os import environ
from shutil import which

import pytest
from selenium import webdriver
from selenium.webdriver.firefox.options import Options
from selenium.webdriver.firefox.service import Service

from dsl.dsl import DSL
from constants import URL


@pytest.fixture
def dsl(driver):
    return DSL(driver, URL)


@pytest.fixture
def driver():
    debug = json.loads(environ.get("DEBUG", "false").lower())
    opts = list(filter(None, [
        "--disable-gpu",
        "--no-sandbox",
        "--disable-dev-shm-usage",
        "--headless" if not debug else "",
    ]))
    options = Options()
    for opt in opts:
        options.add_argument(opt)
    path = which("firefox.geckodriver")
    service = Service(executable_path=path)
    _driver = webdriver.Firefox(options=options, service=service)
    _driver.set_window_size(1920, 1080)
    _driver.implicitly_wait(10)
    yield _driver

    _driver.quit()


@pytest.hookimpl(hookwrapper=True)
def pytest_runtest_makereport(item):
    driver = item.funcargs.get("driver", None)
    outcome = yield
    test_report = outcome.get_result()
    if driver and test_report.when == "call":
        screenshot = driver.get_screenshot_as_base64()
        pytest_html = item.config.pluginmanager.getplugin("html")
        extras = getattr(test_report, "extra", [])
        extras.append(pytest_html.extras.image(screenshot))
        test_report.extras = extras
