import logging

from xtu_ems.common.config import app_list_url
from xtu_ems.common.exception import *
from xtu_ems.common.sess import HttpSessionHolder
from xtu_ems.qz_ems.config import *

logger = logging.getLogger(__name__)


def _pick_app(app_list: list[dict]):
    """
    从应用列表中选择合适的应用

    目前的策略是选择第一个名称包含“教务”的应用

    :param app_list: 应用列表
    :return: 选中的应用字典，如果没有找到合适的应用则返回 None
    """
    for app in app_list:
        if "教务系统（本科生）" in app.get("mc", ""):
            return app
    return None


async def sso_auth(session: HttpSessionHolder) -> HttpSessionHolder:
    """
    使用统一身份认证系统进行登录验证

    :param session: HttpSessionHolder 对象，用于存储会话信息
    :return: 更新后的 HttpSessionHolder 对象
    """
    async with session.to_aiohttp_session() as http_session:
        async with http_session.post(app_list_url) as response:
            logger.warning(f"SSO auth response status: {response.status}")
            if response.status != 200:
                raise SessionInvalidException(f"SSO Authentication Failed, status code: {response.status}")
            # 这里可以根据需要处理返回的数据
            data = await response.json()
        app = _pick_app(data.get("data", []))
        if not app:
            raise ServiceUnavailableException("No suitable application found")
        target_url = app.get("ywxturl")
        if not target_url:
            raise ServiceUnavailableException("Application URL not found")
        async with http_session.get(target_url) as response:
            logger.warning(f"Accessing application URL response status: {response.status}")
            # 查看请求最后一次重定向的 URL
            final_url = str(response.url)
            content = await response.text()
            if "此用户信息不存在,非法登录" in content:
                raise QzAccountNotFoundException("Account not found in QZ EMS")
            if not final_url.startswith(homepage_url_prefix):
                logger.warning(f"Final redirected URL: {final_url}")
                raise SessionInvalidException("Session invalid, please re-authenticate")
            return HttpSessionHolder.from_aiohttp_session(http_session)
