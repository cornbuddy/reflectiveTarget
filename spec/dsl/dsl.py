from time import sleep
import logging

from selenium.webdriver.remote.webdriver import WebDriver, WebElement
from selenium.webdriver.common.by import By

from constants import ANIMATION_DURATION_SECS

log = logging.getLogger(__name__)


class DSL:
    """represents domain specific language for the application"""

    def __init__(self, driver: WebDriver, base_url: str):
        self.driver = driver
        self._url = base_url
        self.driver.get(self._url)

    @property
    def sidebar_toggler(self) -> WebElement:
        """button to toggle navigation sidebar"""
        return self.driver.find_element(By.ID, "sidebar-toggler")

    @property
    def sidebar(self) -> WebElement:
        """navigation sidebar"""
        return self.driver.find_element(By.TAG_NAME, "nav")

    @property
    def submit(self) -> WebElement:
        """button to submit form"""
        return self.driver.find_element(By.XPATH, "//button[@type='submit']")

    def go_to_new_target(self):
        """opens 'New target' view"""
        self.ensure_navigation_opened()
        self.driver.find_element(By.LINK_TEXT, "Targets").click()
        self.driver.find_element(By.LINK_TEXT, "New target").click()

    def signup(self, username: str, password: str, confirm: str = None):
        if confirm is None:
            confirm = password

        self.ensure_navigation_opened()
        self.driver.find_element(By.LINK_TEXT, "Signup").click()
        self.driver.find_element(By.NAME, "username").send_keys(username)
        self.driver.find_element(By.NAME, "password").send_keys(password)
        self.driver.find_element(By.NAME, "confirmation").send_keys(confirm)
        self.submit.click()
        sleep(ANIMATION_DURATION_SECS)

    def login(self, username: str, password: str):
        self.ensure_navigation_opened()
        self.driver.find_element(By.LINK_TEXT, "Login").click()
        self.driver.find_element(By.NAME, "username").send_keys(username)
        self.driver.find_element(By.NAME, "password").send_keys(password)
        self.submit.click()
        sleep(ANIMATION_DURATION_SECS)

    def logout(self):
        self.ensure_navigation_opened()
        self.driver.find_element(By.LINK_TEXT, "Logout").click()

    def ensure_navigation_opened(self):
        sidebar_hidden = not self.sidebar.is_displayed()
        log.info("sidebar hidden: %s", sidebar_hidden)
        if sidebar_hidden:
            log.info("show sidebar")
            self.sidebar_toggler.click()
            sleep(ANIMATION_DURATION_SECS)
            # for some reason, it doesn't work nice after redirects, so doing
            # a bit of recursion to 100% ensure sidebar is displayed
            self.ensure_navigation_opened()
