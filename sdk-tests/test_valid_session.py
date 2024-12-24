from unittest import TestCase, IsolatedAsyncioTestCase

from xtu_ems.ems.handler.valid_session import SessionValidator
from xtu_ems.ems.session import Session


class TestSessionValidator(TestCase):
    def test_handler(self):
        """测试验证session"""
        from common_data import session
        validator = SessionValidator()
        res = validator.handler(session=session)
        self.assertTrue(res)

    def test_handler_with_invalid_session(self):
        """测试验证无效session"""
        validator = SessionValidator()
        session = Session(token="6258ADBA22E73B91E06BD45EBFA9AEBD")
        res = validator.handler(session=session)
        self.assertFalse(res)


class TestAsyncSessionValidator(IsolatedAsyncioTestCase):
    def test_async_handler(self):
        """测试异步验证session"""
        from common_data import session
        validator = SessionValidator()
        res = validator.async_handler(session=session)
        self.assertTrue(res)
