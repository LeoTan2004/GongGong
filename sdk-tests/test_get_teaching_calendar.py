from unittest import TestCase
from unittest.async_case import IsolatedAsyncioTestCase

from common_data import session
from xtu_ems.ems.handler import SessionInvalidException
from xtu_ems.ems.session import Session


class TestTeachingCalendarGetter(TestCase):
    def test_handler(self):
        """测试获取教学日历"""
        from xtu_ems.ems.handler.get_teaching_calendar import TeachingCalendarGetter
        handler = TeachingCalendarGetter()
        resp = handler.handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    def test_handler_with_invalid_session(self):
        """测试无效的session"""
        from xtu_ems.ems.handler.get_teaching_calendar import TeachingCalendarGetter
        handler = TeachingCalendarGetter()
        with self.assertRaises(SessionInvalidException):
            handler.handler(Session(token="invalid_token"))


class TestAsyncTeachingCalendarGetter(IsolatedAsyncioTestCase):
    async def test_async_handler(self):
        """测试异步获取教学日历"""
        from xtu_ems.ems.handler.get_teaching_calendar import TeachingCalendarGetter
        handler = TeachingCalendarGetter()
        resp = await handler.async_handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    async def test_async_handler_with_invalid_session(self):
        """测试异步无效的session"""
        from xtu_ems.ems.handler.get_teaching_calendar import TeachingCalendarGetter
        handler = TeachingCalendarGetter()
        with self.assertRaises(SessionInvalidException):
            await handler.async_handler(Session(token="invalid_token"))
