from shutil import which

import pytest
from selenium import webdriver
from selenium.webdriver.firefox.options import Options
from selenium.webdriver.firefox.service import Service


@pytest.fixture
def driver():
    options = Options()
    opts = [
        "--disable-gpu",
        "--no-sandbox",
        "--disable-dev-shm-usage",
        "--headless",
    ]
    for opt in opts:
        options.add_argument(opt)
    path = which("firefox.geckodriver")
    service = Service(executable_path=path)
    _driver = webdriver.Firefox(options=options, service=service)
    _driver.set_window_size(1920, 1080)
    yield _driver

    _driver.quit()


def open_url(url):

    @pytest.fixture
    def fixture(driver):
        driver.get(url)
        driver.implicitly_wait(10)
        return driver

    return fixture


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
        test_report.extra = extras
