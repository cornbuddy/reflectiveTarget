from os import environ
from random import randint

import pytest
from selenium import webdriver
from selenium.webdriver.common.by import By

ALLOWED_SHOTS = 4
URL = f"http://localhost:{environ['PORT']}"


@pytest.fixture
def driver():
    _driver = webdriver.Chrome()
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


def test_should_handle_mutliple_shots(driver):
    submit = driver.find_element(By.ID, "send")
    assert not submit.is_displayed()

    target = driver.find_element(By.TAG_NAME, "img")
    for _ in range(ALLOWED_SHOTS):
        target.click()
    assert submit.is_displayed()

    shots = driver.find_elements(By.CSS_SELECTOR, "div.shot")
    assert len(shots) == ALLOWED_SHOTS


def test_should_ignore_excessive_shots(driver):
    target = driver.find_element(By.TAG_NAME, "img")
    for _ in range(ALLOWED_SHOTS + randint(1, 10)):
        target.click()
    shots = driver.find_elements(By.CSS_SELECTOR, "div.shot")
    assert len(shots) == ALLOWED_SHOTS
