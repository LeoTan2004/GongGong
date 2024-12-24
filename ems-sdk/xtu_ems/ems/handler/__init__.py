import logging
from abc import ABC, abstractmethod
from typing import Generic, TypeVar

import requests
from aiohttp import ClientSession
from bs4 import BeautifulSoup

from xtu_ems.ems.config import RequestConfig
from xtu_ems.ems.ems import QZEducationalManageSystem
from xtu_ems.ems.session import Session

_R = TypeVar("_R")
"""返回值类型"""

logger = logging.getLogger('xtu-ems.handler')


class Handler(ABC, Generic[_R]):
    @abstractmethod
    def handler(self, session: Session, *args, **kwargs) -> _R:
        """同步函数处理"""
        pass

    @abstractmethod
    async def async_handler(self, session: Session, *args, **kwargs) -> _R:
        """异步处理"""
        pass

    def get_session(self, session: Session):
        sess = requests.session()
        sess.cookies.set(QZEducationalManageSystem.SESSION_NAME, session.token)
        sess.headers.setdefault('User-Agent',
                                'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) '
                                'Chrome/111.0.0.0 Safari/537.36 Edg/111.0.1661.41')
        return sess

    def get_async_session(self, session: Session):
        sess = ClientSession(cookies={QZEducationalManageSystem.SESSION_NAME: session.token})
        sess.headers.setdefault('User-Agent',
                                'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) '
                                'Chrome/111.0.0.0 Safari/537.36 Edg/111.0.1661.41')
        return sess


class EMSGetter(Handler[_R]):
    def handler(self, session: Session, *args, **kwargs) -> _R:
        """获取学生信息"""
        with self.get_session(session) as ems_session:
            logger.debug(f'[{self.__class__.__name__}] 正在获取数据-{self.url()}')
            resp = ems_session.get(self.url(), timeout=RequestConfig.XTU_EMS_REQUEST_TIMEOUT, allow_redirects=False)
            soup = BeautifulSoup(resp.text, 'html.parser')
            return self._extra_info(soup)

    async def async_handler(self, session: Session, *args, **kwargs) -> _R:
        """异步获取学生信息"""
        async with ClientSession(cookies={QZEducationalManageSystem.SESSION_NAME: session.token}) as ems_session:
            logger.debug(f'[{self.__class__.__name__}] 正在异步获取数据-{self.url()}')
            resp = await ems_session.get(self.url(), timeout=RequestConfig.XTU_EMS_REQUEST_TIMEOUT,
                                         allow_redirects=False)
            soup = BeautifulSoup(await resp.text(), 'html.parser')
            return self._extra_info(soup)

    @abstractmethod
    def url(self):
        pass

    @abstractmethod
    def _extra_info(self, soup: BeautifulSoup):
        pass


class EMSPoster(EMSGetter[_R]):
    def handler(self, session: Session, *args, **kwargs) -> _R:
        """获取学生信息"""
        with self.get_session(session) as ems_session:
            logger.debug(f'[{self.__class__.__name__}] 正在获取数据-{self.url()}')
            resp = ems_session.post(url=self.url(), data=self._data(), timeout=RequestConfig.XTU_EMS_REQUEST_TIMEOUT,
                                    allow_redirects=False)
            soup = BeautifulSoup(resp.text, 'html.parser')
            return self._extra_info(soup)

    async def async_handler(self, session: Session, *args, **kwargs) -> _R:
        """异步获取学生信息"""
        async with self.get_async_session(session) as ems_session:
            logger.debug(f'[{self.__class__.__name__}] 正在异步获取数据-{self.url()}')
            resp = await ems_session.post(url=self.url(), data=self._data(),
                                          timeout=RequestConfig.XTU_EMS_REQUEST_TIMEOUT, allow_redirects=False)
            soup = BeautifulSoup(await resp.text(), 'html.parser')
            return self._extra_info(soup)

    @abstractmethod
    def _data(self):
        pass
