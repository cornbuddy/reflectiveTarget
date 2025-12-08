from time import sleep

from selenium.webdriver.remote.webdriver import WebDriver, WebElement
from selenium.webdriver.common.by import By

from constants import SESSION_TOKEN


class DSL:
    """represents domain specific language for the application"""

    def __init__(self, driver: WebDriver, base_url: str):
        self.driver = driver
        self._url = base_url
        self.driver.get(self._url)

    @property
    def sidebar_toggler(self) -> WebElement:
        return self.driver.find_element(By.ID, "sidebar-toggler")

    @property
    def sidebar(self) -> WebElement:
        return self.driver.find_element(By.TAG_NAME, "nav")

    def signup(self, username: str, password: str, confirmation: str = None):
        if confirmation is None:
            confirmation = password

        driver = self.driver
        driver.get(f"{self._url}/signup")
        driver.find_element(By.NAME, "username").send_keys(username)
        driver.find_element(By.NAME, "password").send_keys(password)
        driver.find_element(By.NAME, "confirmation").send_keys(confirmation)
        driver.find_element(By.XPATH, "//button[@type='submit']").click()

    def login(self, username: str, password: str):
        driver = self.driver
        driver.get(f"{self._url}/login")
        driver.find_element(By.NAME, "username").send_keys(username)
        driver.find_element(By.NAME, "password").send_keys(password)
        driver.find_element(By.XPATH, "//button[@type='submit']").click()

    def logout(self):
        self.ensure_navigation_opened()
        before_logout = self.driver.get_cookie(SESSION_TOKEN)
        self.driver.find_element(By.LINK_TEXT, "Logout").click()
        after_logout = self.driver.get_cookie(SESSION_TOKEN)
        assert before_logout != after_logout

    def toggle_navigation(self):
        self.sidebar_toggler.click()
        sleep(1)

    def ensure_navigation_opened(self):
        if not self.sidebar.is_displayed():
            self.sidebar_toggler.click()
