from playwright.sync_api import Page

from constants import URL
from .base import _Form, _Page, LOCATORS


class LoginPage(_Form):
    """represents user login page"""

    def __init__(self, page: Page):
        super().__init__(page, f"{URL}/login")
        self.username_input = page.locator(LOCATORS["username"])
        self.password_input = page.locator(LOCATORS["password"])

    def login(self, username: str, password: str):
        self.log.info("going to login as %s", username)
        self.username_input.fill(username)
        self.password_input.fill(password)
        with self.page.expect_response(self.url) as response:
            self.submit.click()
        if response.value.ok:
            self.log.info("user %s logged in successfully", username)
        else:
            self.log.warning(
                "user %s is not logged in, status: %d, body: %s",
                username, response.value.status, response.value.text(),
            )

    def logout(self):
        self.page.goto(f"{URL}/logout")


class SignupPage(_Form):
    """represents user signup page"""

    def __init__(self, page: Page):
        super().__init__(page, f"{URL}/signup")
        self.username_input = page.locator(LOCATORS["username"])
        self.password_input = page.locator(LOCATORS["password"])
        self.confirmation_input = self.form.locator(LOCATORS["confirmation"])

    def signup(self, username: str, password: str, confirmation: str = None):
        self.log.info("signing up as %s", username)
        if confirmation is None:
            confirmation = password
        self.username_input.fill(username)
        self.password_input.fill(password)
        self.confirmation_input.fill(confirmation)
        with self.page.expect_response(self.url) as response:
            self.submit.click()
        if response.value.ok:
            self.log.info("user %s is signed up successfully", username)
        else:
            self.log.warning(
                "user %s is not signed up, status: %d, body: %s",
                username, response.value.status, response.value.text(),
            )


class TargetsListPage(_Page):
    """represents list of targets"""

    @staticmethod
    def get_id_of_target(name: str, page: Page) -> int:
        """returns id of the target with the given name"""
        a = page.locator("ul").get_by_text(name)
        target_id = a.get_attribute("href").split("/")[-1]
        return int(target_id)

    def __init__(self, page: Page):
        super().__init__(page, f"{URL}/targets")
        self.targets_list = self.page.locator("ul")

    def get_id_by_target_name(self, name: str) -> int:
        """returns id of the target with the given name"""
        return TargetsListPage.get_id_of_target(name, self.page)


class NewTargetPage(_Form):
    """represents new target form"""

    def __init__(self, page: Page):
        super().__init__(page, f"{URL}/targets/new")
        self.name_input = page.locator(LOCATORS["name"])
        self.add_question = page.locator(LOCATORS["add_question"])

    def create_target(self, name: str, questions: list) -> int:
        """creates target with given parameters. returns id of the target"""
        self.log.info("creating target name=%s", name)
        self.name_input.fill(name)
        for i, question in enumerate(questions):
            self.log.debug("creating question i=%d text=%s", i, question)
            if i > 0:
                self.log.debug("adding question input")
                self.add_question.click()
            self.log.debug("looking for input question locator")
            inpt = self.page.locator(LOCATORS["question"](i))
            self.log.debug("filling input text=%s input=%s", question, inpt)
            inpt.fill(question)
            self.log.debug("done filling")
        with self.page.expect_response(self.url):
            self.log.debug("submitting target")
            self.submit.click()
            self.log.debug("target submitted")
        self.log.info("target created name=%s", name)
        self.log.debug("attempt to fetch new target id url=%s", self.page.url)
        return TargetsListPage.get_id_of_target(name, self.page)


class UpdateTargetPage(_Form):
    """represents form to update target"""

    def __init__(self, page: Page, target_id: int):
        super().__init__(page, f"{URL}/targets/{target_id}")
        self.target_id = target_id
        self.name_input = page.locator(LOCATORS["name"])

    def update_target(self, name: str, questions: list):
        self.log.info("updating target name=%s id=%d", name, self.target_id)
        self.name_input.clear()
        self.name_input.fill(name)
        for i, question in enumerate(questions):
            self.log.debug("creating question i=%d text=%s", i, question)
            inpt = self.page.locator(LOCATORS["question"](i))
            self.log.debug("filling question input")
            inpt.clear()
            inpt.fill(question)
        with self.page.expect_response(self.url):
            self.log.debug("submitting target")
            self.submit.click()
            self.log.debug("target submitted")
        self.log.info("target updated name=%s id=%d", name, self.target_id)
