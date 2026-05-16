from playwright.sync_api import Page, Locator

from constants import URL


class _Model:
    """wraps playwright abstractions"""

    def __init__(self, page: Page, url: str):
        self._page = page
        self._url = url
        self.navigate()

    @property
    def page(self) -> Page:
        """link to the page object"""
        return self._page

    @property
    def url(self) -> str:
        """address of the page"""
        return self._url

    def navigate(self):
        """goes to the page address"""
        self.page.goto(self.url)


class Layout(_Model):
    """contains common page components"""

    def __init__(self, page: Page):
        super().__init__(page, URL)

    @property
    def sidebar_toggler(self):
        return self.page.locator("#sidebar-toggler")

    @property
    def sidebar(self):
        return self.page.locator("nav")

    def ensure_navigation_opened(self):
        """opens navigation bar if not opened"""
        sidebar_hidden = self.sidebar.is_hidden()
        if sidebar_hidden:
            self.sidebar_toggler.click()


class _Page(_Model):
    """represents basic page components"""

    def __init__(self, page: Page, url: str):
        super().__init__(page, url)
        self._layout = Layout(page)

    @property
    def layout(self) -> Layout:
        """layout of the page"""
        return self._layout


class _Form(_Page):
    """represents submitable page"""

    @property
    def form(self) -> Locator:
        """html form"""
        return self.page.locator("form")

    @property
    def submit(self) -> Locator:
        """button to submit form"""
        return self.page.locator("//button[@type='submit']")
