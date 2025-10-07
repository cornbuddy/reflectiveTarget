import pytest


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
