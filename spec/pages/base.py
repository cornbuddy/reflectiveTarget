import logging

from playwright.sync_api import Page, Locator


class Locators:
    """contains key locators of the application"""

    @classmethod
    def sidebar_toggler(cls) -> str:
        """button to toggle sidebar"""
        return "#sidebar-toggler"

    @classmethod
    def sidebar(cls) -> str:
        """sidebar menu"""
        return "nav"

    @classmethod
    def form(cls) -> str:
        """http form"""
        return "form"

    @classmethod
    def submit(cls) -> str:
        """button to submit the form"""
        return "//button[@type='submit']"

    @classmethod
    def target(cls) -> str:
        """canvas to draw target"""
        return "canvas"

    @classmethod
    def target_name(cls) -> str:
        """target name form input"""
        return "//input[@name='name']"

    @classmethod
    def add_question(cls) -> str:
        """button to add question to the form"""
        return "//button[text()='Add question']"

    @classmethod
    def questions(cls) -> str:
        """list of all questions of the target"""
        return "//input[starts-with(@placeholder, 'Question')]"

    @classmethod
    def question(cls, i: int) -> str:
        """question with the given index i"""
        return f"//input[@name='question_{i}_value']"

    @classmethod
    def username(cls) -> str:
        """username form input"""
        return "input#username"

    @classmethod
    def password(cls) -> str:
        """password form input"""
        return "input#password"

    @classmethod
    def confirmation(cls) -> str:
        """confirmation form input"""
        return "input#confirmation"


class _Model:
    """wraps playwright abstractions"""

    def __init__(self, page: Page):
        self._page = page
        self._log = logging.getLogger(type(self).__name__)
        self.log.setLevel(logging.DEBUG)

    @property
    def log(self) -> logging.Logger:
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
        return self.page.locator(Locators.sidebar_toggler())

    @property
    def sidebar(self):
        """sidebar with navigation"""
        return self.page.locator(Locators.sidebar())

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
        return self.page.locator(Locators.form())

    @property
    def submit(self) -> Locator:
        """button to submit form"""
        return self.page.locator(Locators.submit())
