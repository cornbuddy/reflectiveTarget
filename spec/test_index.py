from os import environ
from random import randint

from selenium.webdriver.common.by import By

ALLOWED_SHOTS = 4
URL = f"http://localhost:{environ['PORT']}"


def test_should_submit_shots_only_once(driver):
    driver.get(URL)
    driver.implicitly_wait(10)

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
