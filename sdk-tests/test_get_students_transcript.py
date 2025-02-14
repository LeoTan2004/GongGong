from unittest import TestCase
from unittest.async_case import IsolatedAsyncioTestCase

from common_data import session
from xtu_ems.ems.handler import SessionInvalidException
from xtu_ems.ems.session import Session


class TestStudentTranscriptGetter(TestCase):
    def test_handler(self):
        """测试获取学生成绩"""
        from xtu_ems.ems.handler.get_students_transcript import StudentTranscriptGetter
        handler = StudentTranscriptGetter()
        resp = handler.handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    def test_handler_with_invalid_session(self):
        """测试无效的session"""
        from xtu_ems.ems.handler.get_students_transcript import StudentTranscriptGetter
        handler = StudentTranscriptGetter()
        with self.assertRaises(SessionInvalidException):
            handler.handler(Session(token="invalid_token"))


class TestAsyncStudentTranscriptGetter(IsolatedAsyncioTestCase):
    async def test_async_handler(self):
        """测试异步获取学生成绩"""
        from xtu_ems.ems.handler.get_students_transcript import StudentTranscriptGetter
        handler = StudentTranscriptGetter()
        resp = await handler.async_handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    async def test_async_handler_with_invalid_session(self):
        """测试异步无效的session"""
        from xtu_ems.ems.handler.get_students_transcript import StudentTranscriptGetter
        handler = StudentTranscriptGetter()
        with self.assertRaises(SessionInvalidException):
            await handler.async_handler(Session(token="invalid_token"))


class TestStudentRankGetter(TestCase):
    def test_handler(self):
        """测试获取全部学生排名"""
        from xtu_ems.ems.handler.get_students_transcript import StudentRankGetter
        handler = StudentRankGetter()
        resp = handler.handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    def test_handler_with_invalid_session(self):
        """测试无效的session"""
        from xtu_ems.ems.handler.get_students_transcript import StudentRankGetter
        handler = StudentRankGetter()
        with self.assertRaises(SessionInvalidException):
            handler.handler(Session(token="invalid_token"))


class TestAsyncStudentRankGetter(IsolatedAsyncioTestCase):
    async def test_async_handler(self):
        """测试异步获取全部学生排名"""
        from xtu_ems.ems.handler.get_students_transcript import StudentRankGetter
        handler = StudentRankGetter()
        resp = await handler.async_handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    async def test_async_handler_with_invalid_session(self):
        """测试异步无效的session"""
        from xtu_ems.ems.handler.get_students_transcript import StudentRankGetter
        handler = StudentRankGetter()
        with self.assertRaises(SessionInvalidException):
            await handler.async_handler(Session(token="invalid_token"))


class TestStudentRankGetterForCompulsory(TestCase):
    def test_handler(self):
        """测试获必修取学生排名"""
        from xtu_ems.ems.handler.get_students_transcript import StudentRankGetterForCompulsory
        handler = StudentRankGetterForCompulsory()
        resp = handler.handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    def test_handler_with_invalid_session(self):
        """测试无效的session"""
        from xtu_ems.ems.handler.get_students_transcript import StudentRankGetterForCompulsory
        handler = StudentRankGetterForCompulsory()
        with self.assertRaises(SessionInvalidException):
            handler.handler(Session(token="invalid_token"))


class AsyncTestStudentRankGetterForCompulsory(IsolatedAsyncioTestCase):
    async def test_async_handler(self):
        """测试异步获取必修学生排名"""
        from xtu_ems.ems.handler.get_students_transcript import StudentRankGetterForCompulsory
        handler = StudentRankGetterForCompulsory()
        resp = await handler.async_handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    async def test_async_handler_with_invalid_session(self):
        """测试异步无效的session"""
        from xtu_ems.ems.handler.get_students_transcript import StudentRankGetterForCompulsory
        handler = StudentRankGetterForCompulsory()
        with self.assertRaises(SessionInvalidException):
            await handler.async_handler(Session(token="invalid_token"))


class TestStudentTranscriptGetterForAcademicMinor(TestCase):
    def test_handler(self):
        """测试获取辅修成绩单"""
        from xtu_ems.ems.handler.get_students_transcript import StudentTranscriptGetterForAcademicMinor
        handler = StudentTranscriptGetterForAcademicMinor()
        resp = handler.handler(session)
        print(resp.model_dump_json(indent=4))
        self.assertIsNotNone(resp)

    def test_handler_with_invalid_session(self):
        """测试无效的session"""
        from xtu_ems.ems.handler.get_students_transcript import StudentTranscriptGetterForAcademicMinor
        handler = StudentTranscriptGetterForAcademicMinor()
        with self.assertRaises(SessionInvalidException):
            handler.handler(Session(token="invalid_token"))
