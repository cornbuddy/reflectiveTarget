import unittest
from os import environ

from selenium import webdriver
from selenium.webdriver.common.by import By


class TestIndex(unittest.TestCase):
    def setUp(self):
        self.driver = webdriver.Chrome()
        self.url = f"http://localhost:{environ['PORT']}"

    def tearDown(self):
        self.driver.quit()

    def test_index(self):
        self.driver.get(self.url)
        img = self.driver.find_element(by=By.TAG_NAME, value="img")
        self.assertIsNotNone(img)
