from unittest import TestCase
from unittest.async_case import IsolatedAsyncioTestCase

from common_data import session


class TestStudentInfoHandler(TestCase):
    def test_handler(self):
        """测试获取学生信息"""
        from xtu_ems.ems.handler.get_student_info import StudentInfoGetter
        handler = StudentInfoGetter()
        resp = handler.handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)


class TestAsyncStudentInfoHandler(IsolatedAsyncioTestCase):
    async def test_async_handler(self):
        """测试异步获取学生信息"""

        from xtu_ems.ems.handler.get_student_info import StudentInfoGetter
        handler = StudentInfoGetter()
        resp = await handler.async_handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)
