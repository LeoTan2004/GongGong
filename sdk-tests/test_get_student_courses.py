from unittest import TestCase
from unittest.async_case import IsolatedAsyncioTestCase

from common_data import session, username
from xtu_ems.ems.config import RequestConfig, XTUEMSConfig
from xtu_ems.ems.handler import SessionInvalidException
from xtu_ems.ems.handler.get_student_courses import StudentCourseGetter
from xtu_ems.ems.model import InformationPackage
from xtu_ems.ems.session import Session


class TestStudentCourseGetter(TestCase):
    def test_handler(self):
        """测试获取学生课程"""
        handler = StudentCourseGetter()
        resp = handler.handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    def test_extra_student_courses(self):
        """测试解析课程"""
        handler = StudentCourseGetter()
        url = handler.url()
        with handler.get_session(session) as ems_session:
            resp = ems_session.post(url=url, data={"xnxq01id": XTUEMSConfig.get_current_term()},
                                    timeout=RequestConfig.XTU_EMS_REQUEST_TIMEOUT)
        import bs4
        li = handler._extra_info(bs4.BeautifulSoup(resp.text, "html.parser"))
        info = InformationPackage(student_id=username, data=li)
        print(info.model_dump_json(indent=4))
        self.assertIsNotNone(li)

    def test_handler_with_invalid_session(self):
        """测试无效的session"""
        handler = StudentCourseGetter()
        with self.assertRaises(SessionInvalidException):
            handler.handler(Session(token="invalid_token"))


class TestAsyncStudentCourseGetter(IsolatedAsyncioTestCase):

    async def test_async_handler(self):
        """测试异步获取学生课程"""
        handler = StudentCourseGetter()
        resp = await handler.async_handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    async def test_async_handler_with_invalid_session(self):
        """测试异步无效的session"""
        handler = StudentCourseGetter()
        with self.assertRaises(SessionInvalidException):
            await handler.async_handler(Session(token="invalid_token"))
