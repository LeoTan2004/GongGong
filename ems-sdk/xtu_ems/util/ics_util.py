import time
from datetime import date as ddate, timedelta, datetime
from typing import Tuple

from xtu_ems.ems.model import CourseInfo, _get_day_no, ExamInfoList, ExamInfo
from xtu_ems.util.icalendar import BaseEvent, BaseRepeatRule, BaseAlarm


class ExamIcalendarUtil:
    """考试日历工具"""
    DEFAULT_ALARM = [BaseAlarm(trigger=timedelta(days=-7), description="考试提醒，距离考试只有7天了！"),
                     BaseAlarm(trigger=timedelta(days=-3), description="考试提醒，距离考试只有3天了！"),
                     BaseAlarm(trigger=timedelta(hours=-1), description="考试提醒，距离考试只有1小时了！")
                     ]

    def convert_exams_to_events(self, exam_list: ExamInfoList) -> list[BaseEvent]:
        """将考试转换为事件"""
        events = []
        for exam in exam_list.exams:
            if exam.start_time:
                events.append(self.convert_exam_to_event(exam))
        return events

    def convert_exam_to_event(self, exam: ExamInfo) -> BaseEvent:
        """将考试转换为事件"""
        e = BaseEvent(summary=f'【考试】{exam.name}',
                      location=exam.location,
                      start_time=exam.start_time,
                      end_time=exam.end_time,
                      description=f"【{exam.name}】 {exam.location}  ——《拱拱》",
                      alarm=self.DEFAULT_ALARM)
        return e


class CourseIcalendarUtil:
    _initialized = False

    @staticmethod
    def initialize():
        """初始化方法"""
        CourseIcalendarUtil.SUMMER_SAVING_TIME = [
            (datetime.strptime(start, "%H:%M").time(), datetime.strptime(end, "%H:%M").time())
            for start, end in CourseIcalendarUtil.SUMMER_SAVING_TIME]
        CourseIcalendarUtil.WINTER_SAVING_TIME = [
            (datetime.strptime(start, "%H:%M").time(), datetime.strptime(end, "%H:%M").time())
            for start, end in CourseIcalendarUtil.WINTER_SAVING_TIME]

    def __init__(self):
        if not CourseIcalendarUtil._initialized:
            CourseIcalendarUtil._initialized = True
            CourseIcalendarUtil.initialize()

    SUMMER_SAVING_TIME = [("8:00", "8:45"),
                          ("8:55", "9:40"),
                          ("10:10", "10:55"),
                          ("11:05", "11:50"),
                          ("14:30", "15:15"),
                          ("15:25", "16:10"),
                          ("16:40", "17:25"),
                          ("17:35", "18:20"),
                          ("19:30", "20:15"),
                          ("20:25", "21:10"),
                          ("21:20", "22:05")]
    """夏令时"""
    WINTER_SAVING_TIME = [("8:00", "8:45"),
                          ("8:55", "9:40"),
                          ("10:10", "10:55"),
                          ("11:05", "11:50"),
                          ("14:00", "14:45"),
                          ("14:55", "15:40"),
                          ("16:10", "16:55"),
                          ("17:05", "17:50"),
                          ("19:00", "19:45"),
                          ("19:55", "20:40"),
                          ("20:50", "21:35")]
    """冬令时"""

    DEFAULT_ALARM = BaseAlarm(trigger=timedelta(minutes=-25), description="课程提醒")

    def get_time_table(self, base_date) -> Tuple[int, list, list]:
        """
        获取时间表和分割周次
        Args:
            base_date: 开始日期

        Returns:
            切换冬夏令时的周次，，前半学期的时间表，后半学期的时间表
        """
        if base_date.month < 5:
            sep_week = (ddate(base_date.year, month=5, day=1) - base_date).days // 7 + 1
            """切换冬夏令时的周次"""
            former = self.WINTER_SAVING_TIME
            later = self.SUMMER_SAVING_TIME
        else:
            sep_week = (ddate(base_date.year, month=10, day=1) - base_date).days // 7 + 1
            former = self.SUMMER_SAVING_TIME
            later = self.WINTER_SAVING_TIME
        return sep_week, former, later

    def convert_courses_to_events(self, courses: list[CourseInfo], base_date: ddate) -> list[BaseEvent]:
        """将课程转换为事件"""
        events = []
        for course in courses:
            events.extend(self.convert_course_to_event(course, base_date))
        return events

    def convert_course_to_event(self, course: CourseInfo, base_date: ddate) -> list[BaseEvent]:
        """将课程转换为事件"""
        events = []
        sep_week, former, later = self.get_time_table(base_date)
        weeks = course.weeks.split(',')
        for week in weeks:
            if '-' in week:
                start, end = week.split('-')
            else:
                start = end = week
            start = int(start)
            end = int(end)
            """星期，Monday = 0"""
            if end <= sep_week:
                # 前半学期
                events.append(self.convert_single_course(course, base_date, start, end, former))
            elif start <= sep_week < end:
                # 横跨一个变换周次，需切分成两个事件处理
                former_events = self.convert_single_course(course, base_date, start, sep_week, former)
                later_events = self.convert_single_course(course, base_date, sep_week + 1, end, later)
                events.append(former_events)
                events.append(later_events)
            else:
                # 不横跨一个变换周次，直接处理
                events.append(self.convert_single_course(course, base_date, start, end, later))
        return events

    def convert_single_course(self, course: CourseInfo, base_date: ddate, start_week: int, end_week: int,
                              time_table: list[time.struct_time]) -> BaseEvent:
        """将单个课程转换为事件"""
        count = end_week - start_week + 1
        start_date = base_date + timedelta(weeks=start_week - 1, days=_get_day_no(course.day))
        start_time = time_table[course.start_time - 1][0]
        end_time = time_table[course.start_time + course.duration - 2][1]
        return BaseEvent(summary=f'【课程】{course.name}',
                         location=course.classroom,
                         description=f"【{course.teacher}】 {course.duration}节 ——《拱拱》",
                         start_time=datetime.combine(start_date, start_time),
                         end_time=datetime.combine(start_date, end_time),
                         rrule=BaseRepeatRule(freq="WEEKLY", count=count),
                         alarm=self.DEFAULT_ALARM
                         )
