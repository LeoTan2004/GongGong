from unittest import TestCase
from unittest.async_case import IsolatedAsyncioTestCase

from xtu_ems.ems.account import AuthenticationAccount
from xtu_ems.ems.ems import QZEducationalManageSystem


class TestQZEducationalManageSystem(TestCase):
    def test_login(self):
        """测试登录"""
        from common_data import account
        ems = QZEducationalManageSystem()
        session = ems.login(account)
        self.assertIsNotNone(session)

    def test_invalid_account(self):
        """测试无效账户"""
        from xtu_ems.ems.ems import QZEducationalManageSystem
        invalid_account = AuthenticationAccount(username='invalid_username',
                                                password='invalid_password')
        ems = QZEducationalManageSystem()
        from xtu_ems.ems.ems import InvalidAccountException
        with self.assertRaises(InvalidAccountException):
            ems.login(invalid_account)


class TestAsyncQZEducationalManageSystem(IsolatedAsyncioTestCase):
    async def test_async_login(self):
        """测试异步登录"""
        from common_data import account
        ems = QZEducationalManageSystem()
        res = await ems.async_login(account)
        self.assertIsNotNone(res)

    async def test_invalid_async_login(self):
        """异步测试无效账户"""
        from xtu_ems.ems.account import AuthenticationAccount
        from xtu_ems.ems.ems import QZEducationalManageSystem
        invalid_account = AuthenticationAccount(username='invalid_username',
                                                password='invalid_password')
        ems = QZEducationalManageSystem()
        from xtu_ems.ems.ems import InvalidAccountException
        with self.assertRaises(InvalidAccountException):
            await ems.async_login(invalid_account)
