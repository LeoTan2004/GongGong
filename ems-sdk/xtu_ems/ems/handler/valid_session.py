"""鉴定session是否仍然有效"""
from bs4 import BeautifulSoup

from xtu_ems.ems.config import XTUEMSConfig
from xtu_ems.ems.handler import EMSGetter


class SessionValidator(EMSGetter[bool]):
    """会话有效性验证器，你也可以用这个刷新Session在校务系统的有效期"""

    def _extra_info(self, soup: BeautifulSoup):
        """校务系统并不会对于无效的会话并不会重定向，二试直接返回一个登陆页面，我们可以通过返回的页面标题来判断会话是否有效"""
        return soup.find('title').text.strip() != XTUEMSConfig.XTU_EMS_SESSION_VALIDATOR_TITLE

    def url(self):
        return XTUEMSConfig.XTU_EMS_SESSION_VALIDATOR_URL
