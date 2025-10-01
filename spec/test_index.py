from os import environ
from random import randint
from shutil import which

import pytest
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.firefox.options import Options
from selenium.webdriver.firefox.service import Service

ALLOWED_SHOTS = 4
URL = f"http://localhost:{environ['PORT']}"


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
    _driver.get(URL)
    _driver.implicitly_wait(10)
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
