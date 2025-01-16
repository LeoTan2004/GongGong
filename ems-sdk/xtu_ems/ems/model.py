"""校务系统信息"""
from dataclasses import field
from datetime import datetime, date as ddate
from typing import Tuple, Literal, TypeVar, Generic, Union

from pydantic import BaseModel

from xtu_ems.ems.config import classroom_prefix_category

T = TypeVar('T')


class InformationPackage(BaseModel, Generic[T]):
    """信息包"""

    student_id: str
    """学生id"""
    data: T
    """数据"""
    update_time: datetime = datetime.now()
    """更新时间"""


class StudentBasicInfo(BaseModel):
    """学生基本信息"""
    student_id: str
    """学号"""
    name: str
    """姓名"""
    gender: str
    """性别"""
    birthday: str
    """出生日期"""
    major: str
    """专业"""
    class_: str
    """班级"""
    entrance_day: str
    """入学日期"""
    college: str
    """学院"""


class CourseInfo(BaseModel):
    """课程基本信息"""

    name: str
    """课程名称"""
    teacher: str
    """老师"""
    classroom: str
    """教室"""
    weeks: str
    """周次(节次)"""
    start_time: int
    """开始上课节次"""
    duration: int
    """结束上课节次"""
    day: str
    """星期"""


def _get_day_no(weekday: str):
    """
    获取星期几的课程表
    Args:
        weekday: 星期几
    Returns:
        返回星期几的课程表
    """
    return ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'].index(weekday)


def _get_day_name(day: int):
    """
    获取星期几的课程表
    Args:
        day: 星期几

    Returns:
        返回星期几的课程表
    """
    return ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'][day]


class CourseList(BaseModel):
    """课程列表"""

    courses: list[CourseInfo] = field(default_factory=list)
    """课程列表"""


class CourseTable(BaseModel):
    """课程表数据结构"""
    Sunday: list[list[CourseInfo]] = field(default_factory=list)
    Monday: list[list[CourseInfo]] = field(default_factory=list)
    Tuesday: list[list[CourseInfo]] = field(default_factory=list)
    Wednesday: list[list[CourseInfo]] = field(default_factory=list)
    Thursday: list[list[CourseInfo]] = field(default_factory=list)
    Friday: list[list[CourseInfo]] = field(default_factory=list)
    Saturday: list[list[CourseInfo]] = field(default_factory=list)

    def __getitem__(self, item):
        # 判断item是否是一个int，如果为int，则返回相应的课程表
        if isinstance(item, int):
            return self.__getattribute__(_get_day_name(item))
        # 如果不是int，则返回相应的课程表
        return self.__getattribute__(item)

    def to_list(self):
        """从课程表中提取课程列表"""
        ret = []
        for day in self.__dict__.values():
            for course in day:
                ret += course
        return CourseList(courses=ret)


class Score(BaseModel):
    """成绩信息"""

    name: str
    """课程名称"""
    score: str
    """成绩。"""
    credit: str
    """学分"""
    type: Literal['必修', '选修', '跨学科选修']
    """成绩类型"""
    term: int
    """学期"""


class ScoreBoard(BaseModel):
    """成绩信息"""
    student_id: str = ""
    """学号"""
    name: str = ""
    """姓名"""
    college: str = ""
    """学院"""
    major: str = ""
    """专业"""

    scores: list[Score] = field(default_factory=list)
    """成绩列表"""
    total_credit: Tuple[str, str] = ('0', '0')
    """总学分"""
    elective_credit: Tuple[str, str] = ('0', '0')
    """选修课学分"""
    compulsory_credit: Tuple[str, str] = ('0', '0')
    """必修课学分"""
    cross_course_credit: Tuple[str, str] = ('0', '0')
    """跨学科选修学分"""
    average_score: str = '0'
    """平均分"""
    gpa: str = '0'
    """总绩点"""

    cet4: str = None
    """CET4 成绩"""
    cet6: str = None
    """CET6 成绩"""


class RankInfo(BaseModel):
    """排名信息"""

    average_score: str = ""
    """平均分"""
    gpa: str = ""
    """成绩"""
    class_rank: int = -1
    """班级排名"""
    major_rank: int = -1
    """专业排名"""
    terms: list[str] = ""
    """学期"""


class ExamInfo(BaseModel):
    """考试信息"""

    name: str = ""
    """考试名称"""
    start_time: Union[datetime, str] = ''
    """开始时间"""
    end_time: Union[datetime, str] = ''
    """结束时间"""
    location: str = ''
    """考试地点"""
    type: str = '考试'
    """考核方式"""


class ExamInfoList(BaseModel):
    """考试信息列表"""

    exams: list[ExamInfo] = field(default_factory=list)
    """考试信息"""


class ClassroomStatus(BaseModel):
    """教室信息"""

    name: str = ""
    """教室名称"""
    status: list[str] = field(default_factory=list)
    """教室状态"""


class CategoryClassroomBoard(BaseModel):
    """分类的教室信息"""
    classrooms: dict[str, list[ClassroomStatus]]
    """教室信息"""
    date: ddate = datetime.now().date()
    """日期"""


class ClassroomBoard(BaseModel):
    """教室信息"""

    classrooms: list[ClassroomStatus] = field(default_factory=list)
    """教室信息"""
    date: ddate = datetime.now().date()
    """日期"""

    def to_category(self):
        """分类"""
        ret = CategoryClassroomBoard(classrooms={}, date=self.date)
        for classroom in self.classrooms:
            for prefix in classroom_prefix_category:
                if classroom.name.startswith(prefix):
                    ret.classrooms.setdefault(classroom_prefix_category[prefix], []).append(classroom)
                    # 将教室名称中的前缀去除
                    classroom.name = classroom.name[len(prefix):]
                    break
            else:
                ret.classrooms.setdefault('其他', []).append(classroom)
        return ret


class TeachingCalendar(BaseModel):
    """教学日历"""

    start: ddate = None
    """开始时间"""
    weeks: int = None
    """本学期周数"""
    term_id: str = ""
    """学期"""
