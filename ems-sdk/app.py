import logging
import requests.exceptions
from aiohttp import ClientConnectorDNSError
from fastapi import FastAPI, Body, Header
from fastapi.params import Path
from pydantic import BaseModel
from starlette.responses import PlainTextResponse
from typing import Generic, TypeVar

from xtu_ems.ems.account import AuthenticationAccount
from xtu_ems.ems.ems import QZEducationalManageSystem, InvalidCaptchaException, InvalidAccountException, \
    UninitializedPasswordException
from xtu_ems.ems.handler import Handler, SessionInvalidException
from xtu_ems.ems.handler.get_classroom_status import TodayClassroomStatusGetter, TomorrowClassroomStatusGetter, \
    AssignedClassroomStatusGetter
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
    message: str = ""
    data: T = None

    @staticmethod
    def success(msg: str = "success", data: T = None):
        resp = Resp(code=1, message=msg, data=data)
        return resp

    @staticmethod
    def unauthorized(msg: str = "unauthorized"):
        """账户密码错误，或者token失效"""
        return PlainTextResponse(status_code=401, content=msg)

    @staticmethod
    def not_initialized(msg: str = "your account was not initialized"):
        """账户登陆成功，但是可能未初始化，需要在教务系统中认证"""
        return PlainTextResponse(status_code=409, content=msg)

    @staticmethod
    def ems_request_failed(msg: str = "something wrong"):
        """教务系统请求失败或者验证码识别错误"""
        return PlainTextResponse(status_code=503, content=msg)

    @staticmethod
    def error(msg: str = "something wrong"):
        """未知错误"""
        return PlainTextResponse(status_code=500, content=msg)


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
        return Resp.ems_request_failed(f"【{username}】登陆时验证码识别错误")
    except InvalidAccountException as account_exc:
        logger.warning(f"【{username}】登陆失败，账户或者密码错误")
        return Resp.unauthorized(f"【{username}】登陆失败，账户或者密码错误")
    except UninitializedPasswordException as exc:
        logger.warning(f"【{username}】登陆失败，账户未初始化")
        return Resp.not_initialized(f"【{username}】登陆失败，需要先在教务系统中认证")
    except requests.exceptions.Timeout as exc:
        logger.warning(f"【{username}】登陆时超时")
        return Resp.ems_request_failed("远程连接错误")
    except ClientConnectorDNSError as e:
        logger.exception(f"无法访问服务")
        return Resp.ems_request_failed("远程无法访问")
    except Exception as e:
        logger.exception(f"【{username}】登陆时错误")
        return Resp.error("未知错误")


async def _run_handler(handler: Handler, token: str):
    session = Session(token=token)
    try:
        return Resp.success(data=await handler.async_handler(session))
    except requests.exceptions.Timeout as e:
        logger.exception(f"【{handler.__class__.__name__}】执行时超时")
        return Resp.ems_request_failed("远程连接错误")
    except SessionInvalidException as e:
        logger.warning(f"【{handler.__class__.__name__}】执行时session失效")
        return Resp.unauthorized("session失效")
    except ClientConnectorDNSError as e:
        logger.exception(f"【{handler.__class__.__name__}】无法访问服务")
        return Resp.ems_request_failed("远程无法访问")
    except Exception as e:
        logger.exception(f"【{handler.__class__.__name__}】执行时错误")
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


@api.get("/classroom/{day}")
async def get_classroom(
    day: str = Path(description="举例当天的时间，例如今天为0,明天为1,或者也可以今天today,明天为tomorrow"),
    token: str = Header(description="用户凭证")):
    """获取指定日期教室"""
    if day == "0" or day == "today":
        return await _run_handler(today_classroom_status_getter, token)
    elif day == "1" or day == "tomorrow":
        return await _run_handler(tomorrow_classroom_status_getter, token)
    else:
        return await _run_handler(AssignedClassroomStatusGetter(int(day)), token)


@api.get("/compulsory/rank")
async def get_compulsory_rank(token: str = Header(description="用户凭证")):
    """获取必修排名"""
    return await _run_handler(major_compulsory_rank_getter, token)


@api.get("/calendar")
async def get_calendar(token: str = Header(description="用户凭证")):
    """获取校历"""
    return await _run_handler(calendar_getter, token)


if __name__ == '__main__':
    import uvicorn

    uvicorn.run(app=api, host="0.0.0.0", port=8080, log_config="log_config.json")
