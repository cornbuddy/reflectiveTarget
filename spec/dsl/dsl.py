from time import sleep
import logging

from playwright.sync_api import Page, Locator

from constants import ANIMATION_DURATION_SECS

log = logging.getLogger(__name__)


class DSL:
    """represents domain specific language for the application"""

    def __init__(self, page: Page, base_url: str):
        self.page = page
        self._url = base_url
        self.page.goto(self._url)

    @property
    def sidebar_toggler(self) -> Locator:
        """button to toggle navigation sidebar"""
        return self.page.locator("#sidebar-toggler")

    @property
    def sidebar(self) -> Locator:
        """navigation sidebar"""
        return self.page.locator("nav")

    @property
    def submit(self) -> Locator:
        """button to submit form"""
        return self.page.locator("//button[@type='submit']")

    def go_to_new_target(self):
        """opens 'New target' view"""
        self.ensure_navigation_opened()
        self.page.get_by_text("Targets").click()
        self.page.get_by_text("New target").click()

    def signup(self, username: str, password: str, confirmation: str = None):
        if confirmation is None:
            confirmation = password

        self.ensure_navigation_opened()
        self.page.get_by_text("Signup").click()
        self.page.locator("input#username").fill(username)
        self.page.locator("input#password").fill(password)
        self.page.locator("input#confirmation").fill(confirmation)
        self.submit.click()
        sleep(ANIMATION_DURATION_SECS)

    def login(self, username: str, password: str):
        self.ensure_navigation_opened()
        self.page.get_by_text("Login").click()
        self.page.locator("input#username").fill(username)
        self.page.locator("input#password").fill(password)
        self.submit.click()
        sleep(ANIMATION_DURATION_SECS)

    def logout(self):
        self.ensure_navigation_opened()
        self.page.get_by_text("Logout").click()

    def ensure_navigation_opened(self):
        sidebar_hidden = self.sidebar.is_hidden()
        log.info("sidebar hidden: %s", sidebar_hidden)
        if sidebar_hidden:
            log.info("show sidebar")
            self.sidebar_toggler.click()
            sleep(ANIMATION_DURATION_SECS)
            # for some reason, it doesn't work nice after redirects, so doing
            # a bit of recursion to 100% ensure sidebar is displayed
            self.ensure_navigation_opened()
