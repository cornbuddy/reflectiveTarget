import pytest


@pytest.hookimpl(wrapper=True)
def pytest_runtest_makereport(item, call):
    raise RuntimeError("not implemented")
