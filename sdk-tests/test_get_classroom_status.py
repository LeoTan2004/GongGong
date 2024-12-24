from unittest import TestCase, IsolatedAsyncioTestCase

from common_data import session
from xtu_ems.ems.handler.get_classroom_status import TodayClassroomStatusGetter, TomorrowClassroomStatusGetter


class TestTodayClassroomStatusGetter(TestCase):
    def test_handler(self):
        """测试获取今日空教室"""
        handler = TodayClassroomStatusGetter()
        resp = handler.handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)


class TestAsyncTodayClassroomStatusGetter(IsolatedAsyncioTestCase):

    async def test_async_handler(self):
        """测试异步获取今日空教室"""
        handler = TomorrowClassroomStatusGetter()
        resp = await handler.async_handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)


class TestTomorrowClassroomStatusGetter(TestCase):
    def test_handler(self):
        """测试获取明日空教室"""
        handler = TomorrowClassroomStatusGetter()
        resp = handler.handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)


class TestAsyncTomorrowClassroomStatusGetter(IsolatedAsyncioTestCase):
    async def test_async_handler(self):
        """测试异步获取明日空教室"""
        handler = TomorrowClassroomStatusGetter()
        resp = await handler.async_handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)
