import datetime
from unittest import TestCase

from xtu_ems.ems.model import CourseInfo, ExamInfo
from xtu_ems.util.icalendar import BaseCalendar
from xtu_ems.util.ics_util import CourseIcalendarUtil, ExamIcalendarUtil


class TestCourseIcalendarUtil(TestCase):
    """测试ics工具"""

    def test_convert_course_to_event(self):
        """测试课程转化为事件"""
        course = CourseInfo(
            name="网络安全协议分析",
            teacher="周维副教授",
            classroom="北山二阶梯",
            weeks="2-6,8-14",
            start_time=5,
            duration=2,
            day="Monday",
        )
        util = CourseIcalendarUtil()
        events = util.convert_course_to_event(course, datetime.date(year=2024, month=8, day=26))
        calendar = BaseCalendar()
        for event in events:
            calendar.add_event(event)
        ics = calendar.to_ical()
        print(ics)
        self.assertIsNotNone(ics)


class TestExamIcalendarUtil(TestCase):
    """测试考试转化为事件"""

    def test_convert_exam_to_event(self):
        """测试考试转化为事件"""
        exam = ExamInfo(name="网络安全协议分析",
                        location="兴教楼A203",
                        start_time=datetime.datetime(year=2024, month=12, day=30, hour=10, minute=30, second=0),
                        end_time=datetime.datetime(year=2024, month=12, day=30, hour=12, minute=30, second=0))
        util = ExamIcalendarUtil()
        event = util.convert_exam_to_event(exam)
        calendar = BaseCalendar()
        calendar.add_event(event)
        ics = calendar.to_ical()
        print(ics)
        self.assertIsNotNone(ics)
