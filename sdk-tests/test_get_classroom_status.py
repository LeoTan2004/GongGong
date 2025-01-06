from unittest import TestCase, IsolatedAsyncioTestCase

from common_data import session
from xtu_ems.ems.handler import SessionInvalidException
from xtu_ems.ems.handler.get_classroom_status import TodayClassroomStatusGetter, TomorrowClassroomStatusGetter, \
    AssignedClassroomStatusGetter
from xtu_ems.ems.session import Session


class TestTodayClassroomStatusGetter(TestCase):
    def test_handler(self):
        """测试获取今日空教室"""
        handler = TodayClassroomStatusGetter()
        resp = handler.handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    def test_handler_with_invalid_session(self):
        """测试无效的session"""
        handler = TodayClassroomStatusGetter()
        with self.assertRaises(SessionInvalidException):
            handler.handler(Session(token="invalid_token"))


class TestAsyncTodayClassroomStatusGetter(IsolatedAsyncioTestCase):

    async def test_async_handler(self):
        """测试异步获取今日空教室"""
        handler = TomorrowClassroomStatusGetter()
        resp = await handler.async_handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    async def test_async_handler_with_invalid_session(self):
        """测试异步无效的session"""
        handler = TomorrowClassroomStatusGetter()
        with self.assertRaises(SessionInvalidException):
            await handler.async_handler(Session(token="invalid_token"))


class TestTomorrowClassroomStatusGetter(TestCase):
    def test_handler(self):
        """测试获取明日空教室"""
        handler = TomorrowClassroomStatusGetter()
        resp = handler.handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    def test_handler_with_invalid_session(self):
        """测试无效的session"""
        handler = TomorrowClassroomStatusGetter()
        with self.assertRaises(SessionInvalidException):
            handler.handler(Session(token="invalid_token"))


class TestAsyncTomorrowClassroomStatusGetter(IsolatedAsyncioTestCase):
    async def test_async_handler(self):
        """测试异步获取明日空教室"""
        handler = TomorrowClassroomStatusGetter()
        resp = await handler.async_handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    async def test_async_handler_with_invalid_session(self):
        """测试异步无效的session"""
        handler = TomorrowClassroomStatusGetter()
        with self.assertRaises(SessionInvalidException):
            await handler.async_handler(Session(token="invalid_token"))


class TestAssignedClassroomStatusGetter(TestCase):
    def test_handler(self):
        """测试获取指定日期空教室"""
        handler = AssignedClassroomStatusGetter(day=3)
        resp = handler.handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    def test_today(self):
        """测试获取今日空教室"""
        handler = AssignedClassroomStatusGetter.today()
        resp = handler.handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    def test_tomorrow(self):
        """测试获取明日空教室"""
        handler = AssignedClassroomStatusGetter.tomorrow()
        resp = handler.handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    def test_handler_with_invalid_session(self):
        """测试无效的session"""
        handler = AssignedClassroomStatusGetter(day=3)
        with self.assertRaises(SessionInvalidException):
            handler.handler(Session(token="invalid_token"))


class TestAsyncAssignedClassroomStatusGetter(IsolatedAsyncioTestCase):
    async def test_async_handler(self):
        """测试异步获取指定日期空教室"""
        handler = AssignedClassroomStatusGetter(day=3)
        resp = await handler.async_handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    async def test_async_today(self):
        """测试异步获取今日空教室"""
        handler = AssignedClassroomStatusGetter.today()
        resp = await handler.async_handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    async def test_async_tomorrow(self):
        """测试异步获取明日空教室"""
        handler = AssignedClassroomStatusGetter.tomorrow()
        resp = await handler.async_handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    async def test_async_handler_with_invalid_session(self):
        """测试异步无效的session"""
        handler = AssignedClassroomStatusGetter(day=3)
        with self.assertRaises(SessionInvalidException):
            await handler.async_handler(Session(token="invalid_token"))
