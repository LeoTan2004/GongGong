from unittest import TestCase, IsolatedAsyncioTestCase

from common_data import session
from xtu_ems.ems.handler.get_student_exam import StudentExamGetter


class TestStudentExamGetter(TestCase):
    def test_handler(self):
        """测试获取学生考试信息"""
        handler = StudentExamGetter()
        resp = handler.handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)


class TestAsyncStudentExamGetter(IsolatedAsyncioTestCase):
    async def test_async_handler(self):
        """测试异步获取学生考试信息"""
        handler = StudentExamGetter()
        resp = await handler.async_handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)
