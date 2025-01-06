from unittest import TestCase, IsolatedAsyncioTestCase

from common_data import session
from xtu_ems.ems.handler import SessionInvalidException
from xtu_ems.ems.handler.get_student_exam import StudentExamGetter
from xtu_ems.ems.session import Session


class TestStudentExamGetter(TestCase):
    def test_handler(self):
        """测试获取学生考试信息"""
        handler = StudentExamGetter()
        resp = handler.handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    def test_handler_with_invalid_session(self):
        """测试无效的session"""
        handler = StudentExamGetter()
        with self.assertRaises(SessionInvalidException):
            handler.handler(Session(token="invalid_token"))


class TestAsyncStudentExamGetter(IsolatedAsyncioTestCase):
    async def test_async_handler(self):
        """测试异步获取学生考试信息"""
        handler = StudentExamGetter()
        resp = await handler.async_handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    async def test_async_handler_with_invalid_session(self):
        """测试异步无效的session"""
        handler = StudentExamGetter()
        with self.assertRaises(SessionInvalidException):
            await handler.async_handler(Session(token="invalid_token"))
