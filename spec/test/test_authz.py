import pytest
from selenium.webdriver.common.by import By

from .constants import PASSWORD, USERNAME, URL


LOGIN_URL = f"{URL}/login"
SIGNUP_URL = f"{URL}/signup"


@pytest.mark.parametrize("url", [SIGNUP_URL])
def test_should_signup(driver):
    driver.find_element(By.NAME, "username").send_keys(USERNAME)
    driver.find_element(By.NAME, "password").send_keys(PASSWORD)
    driver.find_element(By.XPATH, "//button[@type='submit']").click()
    assert driver.current_url == URL


@pytest.mark.order(after="test_should_register")
@pytest.mark.parametrize("url", [LOGIN_URL])
def test_should_login(driver):
    driver.find_element(By.NAME, "username").send_keys(USERNAME)
    driver.find_element(By.NAME, "password").send_keys(PASSWORD)
    driver.find_element(By.XPATH, "//button[@type='submit']").click()
    assert driver.current_url == URL
