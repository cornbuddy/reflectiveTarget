from logging import getLogger, Logger

from playwright.sync_api import Page, Locator


# holds the information about most used locators
LOCATORS = {
    # layout
    "sidebar_toggler": "#sidebar-toggler",
    "sidebar": "nav",
    "form": "form",
    "submit": "//button[@type='submit']",
    # target form
    "canvas": "canvas",
    "name": "//input[@name='name']",
    "add_question": "//button[text()='Add question']",
    "questions": "//input[starts-with(@name, 'question_')]",
    "question": lambda i: f"//input[@name='question_{i}']",
    # authz form
    "username": "input#username",
    "password": "input#password",
    "confirmation": "input#confirmation",
}


class _Model:
    """wraps playwright abstractions"""

    def __init__(self, page: Page):
        self._page = page
        self._log = getLogger(type(self).__name__)

    @property
    def log(self) -> Logger:
        """class logger"""
        return self._log

    @property
    def page(self) -> Page:
        """link to the page object"""
        return self._page


class Layout(_Model):
    """contains common page components"""

    @property
    def sidebar_toggler(self):
        """button to toggle sidebar"""
        return self.page.locator(LOCATORS["sidebar_toggler"])

    @property
    def sidebar(self):
        """sidebar with navigation"""
        return self.page.locator(LOCATORS["nav"])

    def ensure_navigation_opened(self):
        """opens navigation bar if not opened"""
        sidebar_hidden = self.sidebar.is_hidden()
        self.log.info("sidebar hidden: %s", sidebar_hidden)
        if sidebar_hidden:
            self.log.info("togging sidebar")
            self.sidebar_toggler.click()


class _Page(_Model):
    """represents basic page components"""

    def __init__(self, page: Page, url: str):
        super().__init__(page)
        self._layout = Layout(page)
        self._url = url
        self.navigate()

    @property
    def layout(self) -> Layout:
        """layout of the page"""
        return self._layout

    @property
    def url(self) -> str:
        """address of the page"""
        return self._url

    def navigate(self):
        """goes to the page address"""
        self.log.info("going to %s", self.url)
        self.page.goto(self.url)


class _Form(_Page):
    """represents submitable page"""

    @property
    def form(self) -> Locator:
        """html form"""
        return self.page.locator(LOCATORS["form"])

    @property
    def submit(self) -> Locator:
        """button to submit form"""
        return self.page.locator(LOCATORS["submit"])
