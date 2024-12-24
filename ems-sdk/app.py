import logging
from typing import Generic, TypeVar

import requests.exceptions
from fastapi import FastAPI, Body, Header
from pydantic import BaseModel

from xtu_ems.ems.account import AuthenticationAccount
from xtu_ems.ems.ems import QZEducationalManageSystem, InvalidCaptchaException, InvalidAccountException, \
    UninitializedPasswordException
from xtu_ems.ems.handler import Handler
from xtu_ems.ems.handler.get_classroom_status import TodayClassroomStatusGetter, TomorrowClassroomStatusGetter
from xtu_ems.ems.handler.get_student_courses import StudentCourseGetter
from xtu_ems.ems.handler.get_student_exam import StudentExamGetter
from xtu_ems.ems.handler.get_student_info import StudentInfoGetter
from xtu_ems.ems.handler.get_students_transcript import StudentTranscriptGetter, \
    StudentTranscriptGetterForAcademicMinor, StudentRankGetter, StudentRankGetterForCompulsory
from xtu_ems.ems.handler.get_teaching_calendar import TeachingCalendarGetter
from xtu_ems.ems.session import Session

api = FastAPI()

ems = QZEducationalManageSystem()
"""校务系统"""

today_classroom_status_getter = TodayClassroomStatusGetter()
"""当日教室状态获取"""

tomorrow_classroom_status_getter = TomorrowClassroomStatusGetter()
"""次日教室状态获取"""

courses_table_getter = StudentCourseGetter()
"""课程表获取"""

exams_getter = StudentExamGetter()
"""考试安排获取"""

info_getter = StudentInfoGetter()
"""基本信息获取"""

major_scores_getter = StudentTranscriptGetter()
"""主修成绩获取"""

minor_scores_getter = StudentTranscriptGetterForAcademicMinor()
"""辅修成绩获取"""

major_total_rank_getter = StudentRankGetter()
"""主修总排名获取"""

major_compulsory_rank_getter = StudentRankGetterForCompulsory()
"""主修必修排名获取"""

calendar_getter = TeachingCalendarGetter()
"""教学周历获取"""

T = TypeVar("T")
logger = logging.getLogger("api")


class Resp(BaseModel, Generic[T]):
    """统一返回"""
    code: int = 0
    msg: str = ""
    data: T = None

    @staticmethod
    def success(msg: str = "success", data: T = None):
        return Resp(code=1, msg=msg, data=data)

    @staticmethod
    def fail(msg: str = "failed"):
        return Resp(code=0, msg=msg, data={})

    @staticmethod
    def error(msg: str = "error"):
        return Resp(code=-1, msg=msg, data={})


@api.post("/login")
async def login(username: str = Body(description="学号"), password: str = Body(description="密码")):
    logger.debug(f"【{username}】开始登陆")
    try:
        session = await ems.async_login(AuthenticationAccount(username, password),
                                        retry_time=3)
        logger.info(f"【{username}】登陆成功")
        return Resp.success(msg=f"{username}-登陆成功", data=session)
    except InvalidCaptchaException as captcha_exc:
        logger.warning(f"【{username}】登陆时验证码识别失败")
        return Resp.fail(f"【{username}】登陆时验证码识别错误")
    except InvalidAccountException as account_exc:
        logger.warning(f"【{username}】登陆失败，账户或者密码错误")
        return Resp.fail(f"【{username}】登陆失败，账户或者密码错误")
    except UninitializedPasswordException as exc:
        logger.warning(f"【{username}】登陆失败，账户未初始化")
        return Resp.fail(f"【{username}】登陆失败，需要先在教务系统中认证")
    except requests.exceptions.Timeout as exc:
        logger.warning(f"【{username}】登陆时超时")
        return Resp.error("远程连接错误")
    except Exception as e:
        logger.exception(f"【{username}】登陆时错误")
        return Resp.error("未知错误")


async def _run_handler(handler: Handler, token: str):
    session = Session(token=token)
    try:
        return Resp.success(data=await handler.async_handler(session))
    except requests.exceptions.Timeout as e:
        logger.exception(f"【{handler.__name__}】执行时超时")
        return Resp.error("远程连接错误")
    except Exception as e:
        logger.exception(f"【{handler.__name__}】执行时错误")
        return Resp.error("未知错误")


@api.get("/courses")
async def get_courses(token: str = Header(description="用户凭证")):
    """获取课表"""
    return await _run_handler(courses_table_getter, token)


@api.get("/info")
async def get_info(token: str = Header(description="用户凭证")):
    """获取用户信息"""
    return await _run_handler(info_getter, token)


@api.get("/scores")
async def get_scores(token: str = Header(description="用户凭证")):
    """获取成绩"""
    return await _run_handler(major_scores_getter, token)


@api.get("/minor/scores")
async def get_minor_score(token: str = Header(description="用户凭证")):
    """获取辅修成绩"""
    return await _run_handler(minor_scores_getter, token)


@api.get("/exams")
async def get_exams(token: str = Header(description="用户凭证")):
    """获取考试"""
    return await _run_handler(exams_getter, token)


@api.get("/rank")
async def get_major_rank(token: str = Header(description="用户凭证")):
    """获取主修排名"""
    return await _run_handler(major_total_rank_getter, token)


@api.get("/classroom/today")
async def get_today_classroom(token: str = Header(description="用户凭证")):
    """获取今天教室"""
    return await _run_handler(today_classroom_status_getter, token)


@api.get("/classroom/tomorrow")
async def get_tomorrow_classroom(token: str = Header(description="用户凭证")):
    """获取明天教室"""
    return await _run_handler(tomorrow_classroom_status_getter, token)


@api.get("/calendar")
async def get_calendar(token: str = Header(description="用户凭证")):
    """获取校历"""
    return await _run_handler(calendar_getter, token)


if __name__ == '__main__':
    import uvicorn

    uvicorn.run(app=api, host="0.0.0.0", port=8000, log_config="log_config.json")
