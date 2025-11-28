from selenium.webdriver.remote.webdriver import WebDriver
from selenium.webdriver.common.by import By

from constants import SESSION_TOKEN


class DSL:
    """represents domain specific language for the application"""

    def __init__(self, driver: WebDriver, base_url: str):
        self.driver = driver
        self._url = base_url

    def signup(self, username: str, password: str) -> WebDriver:
        self.driver.get(f"{self._url}/signup")
        self.driver.find_element(By.NAME, "username").send_keys(username)
        self.driver.find_element(By.NAME, "password").send_keys(password)
        self.driver.find_element(By.XPATH, "//button[@type='submit']").click()
        session = self.driver.get_cookie(SESSION_TOKEN)
        assert session is not None

        return self.driver

    def login(self, username: str, password: str) -> WebDriver:
        self.driver.get(f"{self._url}/login")
        self.driver.find_element(By.NAME, "username").send_keys(username)
        self.driver.find_element(By.NAME, "password").send_keys(password)
        self.driver.find_element(By.XPATH, "//button[@type='submit']").click()
        session = self.driver.get_cookie(SESSION_TOKEN)
        assert session is not None

        return self.driver

    def logout(self) -> WebDriver:
        self.driver.find_element(By.ID, "sidebar-toggler").click()
        sidebar = self.driver.find_element(By.TAG_NAME, "nav")
        assert sidebar.is_displayed()

        before_logout = self.driver.get_cookie(SESSION_TOKEN)
        self.driver.find_element(By.LINK_TEXT, "Logout").click()
        after_logout = self.driver.get_cookie(SESSION_TOKEN)
        assert before_logout != after_logout

        return self.driver
