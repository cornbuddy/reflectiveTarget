from os import environ
from random import randint

import pytest
from selenium import webdriver
from selenium.webdriver.common.by import By
import chromedriver_autoinstaller

ALLOWED_SHOTS = 4
URL = f"http://localhost:{environ['PORT']}"


@pytest.fixture
def driver():
    chromedriver_autoinstaller.install()
    chrome_options = webdriver.ChromeOptions()
    options = [
        "--ignore-certificate-errors",
        # These flags BELOW are recommended for stability when running Chrome
        # in headless or containerized environments (such as GitHub Actions).
        "--disable-gpu",
        "--no-sandbox",
        "--disable-dev-shm-usage",
        "--remote-debugging-port=9222",
    ]
    for option in options:
        chrome_options.add_argument(option)
    _driver = webdriver.Chrome(options=chrome_options)
    _driver.get(URL)
    yield _driver

    _driver.quit()


def test_should_submit_shots_only_once(driver):
    target = driver.find_element(By.TAG_NAME, "img")
    for _ in range(ALLOWED_SHOTS):
        target.click()
    submit = driver.find_element(By.ID, "send")
    submit.click()
    driver.refresh()
    shots = driver.find_elements(By.CSS_SELECTOR, "div.shot")
    assert len(shots) == ALLOWED_SHOTS

    submit = driver.find_element(By.ID, "send")
    target = driver.find_element(By.TAG_NAME, "img")
    shots = driver.find_elements(By.CSS_SELECTOR, "div.shot")
    for _ in range(randint(1, 10)):
        target.click()
    submit.click()
    assert len(shots) == ALLOWED_SHOTS
