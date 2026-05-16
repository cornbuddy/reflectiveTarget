from playwright.sync_api import Page

from constants import URL
from .base import _Form, _Page


class LoginPage(_Form):
    """represents user login page"""

    def __init__(self, page: Page):
        super().__init__(page, f"{URL}/login")
        self.username_input = self.form.locator("input#username")
        self.password_input = self.form.locator("input#password")

    def login(self, username: str, password: str):
        self.username_input.fill(username)
        self.password_input.fill(password)
        self.submit.click()

    def logout(self):
        self._page.goto(f"{URL}/logout")


class NewTargetPage(_Form):
    """represents new target form"""

    def __init__(self, page: Page):
        super().__init__(page, f"{URL}/targets/new")


class SignupPage(_Form):
    """represents user signup page"""

    def __init__(self, page: Page):
        super().__init__(page, f"{URL}/signup")
        self.username_input = self.form.locator("input#username")
        self.password_input = self.form.locator("input#password")
        self.confirmation_input = self.form.locator("input#confirmation")

    def signup(self, username: str, password: str, confirmation: str = None):
        if confirmation is None:
            confirmation = password
        self.layout.ensure_navigation_opened()
        self.username_input.fill(username)
        self.password_input.fill(password)
        self.confirmation_input.fill(confirmation)
        self.submit.click()


class TargetsListPage(_Page):
    """represents list of targets"""

    def __init__(self, page: Page):
        super().__init__(page, f"{URL}/targets")
        self.targets_list = self.page.locator("ul")


class UpdateTargetPage(_Form):
    """represents form to update target"""

    def __init__(self, page: Page, target_id: int):
        super().__init__(page, f"{URL}/targets/{target_id}")
