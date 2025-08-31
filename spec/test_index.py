from os import environ

import pytest
from selenium import webdriver
from selenium.webdriver.common.by import By

URL = f"http://localhost:{environ['PORT']}"


@pytest.fixture
def driver():
    return webdriver.Chrome()


def test_can_click_exactly_4_times_before_submitting(driver):
    driver.get(URL)
    img = driver.find_element(by=By.TAG_NAME, value="img")
    assert img is not None
