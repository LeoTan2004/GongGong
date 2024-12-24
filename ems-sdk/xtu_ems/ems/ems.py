"""教务系统模块"""
import json
import logging
from abc import ABC, abstractmethod

import requests
from aiohttp import ClientSession

from xtu_ems.ems.account import AuthenticationAccount
from xtu_ems.ems.config import XTUEMSConfig, RequestConfig
from xtu_ems.ems.session import Session
from xtu_ems.util.captcha import ImageDetector

logger = logging.getLogger('xtu-ems.login')


class EducationalManageSystem(ABC):
    """教务系统类"""

    @abstractmethod
    def login(self, account: AuthenticationAccount, retry_time=2) -> Session:
        """教务系统登陆，返回登陆session"""
        pass

    @abstractmethod
    async def async_login(self, account: AuthenticationAccount, retry_time=3) -> Session:
        """
        登陆教务系统
        Args:
            account: 账户信息
            retry_time: 重试次数，如果账号密码错误会直接抛出异常，不会重试

        Returns:
            登陆后的session
        Raises:
            InvalidAccountException: 账号密码错误
            InvalidCaptchaException: 多次尝试仍然验证码错误
            UninitializedPasswordException: 未初始化密码
        """
        pass


class InvalidAccountException(Exception):
    """无效账户异常"""

    def __init__(self, message="账号名或者密码错误"):
        self.message = message
        super().__init__(self.message)


class InvalidCaptchaException(Exception):
    """无效验证码异常"""

    def __init__(self, message="验证码错误"):
        self.message = message
        super().__init__(self.message)


# 未初始化密码
class UninitializedPasswordException(Exception):
    """未初始化密码异常"""

    def __init__(self, message="未初始化密码"):
        self.message = message
        super().__init__(self.message)


class QZEducationalManageSystem(EducationalManageSystem):
    """
    强智教务系统
    """
    SESSION_NAME = "JSESSIONID"
    """Session名"""

    _HEADER = {'Content-Type': 'application/x-www-form-urlencoded'}
    """公共头部"""

    def _data(self, username: str, password: str, encode: str, random_code: str):
        return {
            "USERNAME": username,
            "PASSWORD": password,
            "encoded": encode,
            "RANDOMCODE": random_code
        }

    def __init__(self) -> None:
        super().__init__()
        self.captcha = ImageDetector()

    def pre_check(self, account: AuthenticationAccount):
        """
        检查账户是否合法
        Args:
            account: 账户信息

        Returns:
            检查结果

        """
        if (account.username is None
            or account.username.strip() == ''
            or account.password is None):
            raise Exception("用户名或密码为空")

    def login(self, account: AuthenticationAccount, retry_time=3) -> Session:
        """
        登陆教务系统
        Args:
            account: 账户信息
            retry_time: 重试次数，如果账号密码错误会直接抛出异常，不会重试

        Returns:
            登陆后的session
        Raises:
            InvalidAccountException: 账号密码错误
            InvalidCaptchaException: 多次验证码错误
            UninitializedPasswordException: 未初始化密码
        """
        self.pre_check(account)
        err = None
        if retry_time <= 0:
            raise Exception("登陆失败")
        for i in range(retry_time):
            try:
                session = self._login(account)
                logger.info(f'登陆成功-{account.username}')
                return session
            except InvalidAccountException as e:
                err = e
                logger.info(f'错误的账号密码-{account.username}')
                break
            except InvalidCaptchaException as e:
                err = e
                logger.debug(f'正在重试{i}/{retry_time}次登陆失败-验证码错误')
                continue
            except UninitializedPasswordException as e:
                err = e
                logger.debug(f'[{account.username}]-未初始化密码')
                break
            except TimeoutError as e:
                err = e
                logger.debug(f'正在重试{i}/{retry_time}次登陆失败-网络超时')
                continue
        raise err

    async def async_post_process(self, resp, session):
        """
        处理登陆后的响应
        Args:
            resp: 登陆后的响应
            session: 登陆后的session
        Returns:
            处理后的响应
        Raises:
            InvalidAccountException: 账号密码错误
            InvalidCaptchaException: 验证码错误
            UninitializedPasswordException: 未初始化密码
        """
        if resp.status != 302:
            content = await resp.text()
            if "用户名或密码错误,请联系本院教务老师!" in content:
                raise InvalidAccountException()
            if "验证码错误!!" in content:
                raise InvalidCaptchaException()
            raise Exception("登陆失败")
        else:
            location = resp.headers.get("location")
            if location.endswith(XTUEMSConfig.XTU_EMS_UPDATE_PASSWORD_URL):
                raise UninitializedPasswordException()
            return session

    def post_process(self, resp, session):
        """
        处理登陆后的响应
        Args:
            resp: 登陆后的响应
            session: 登陆后的session
        Returns:
            处理后的响应
        Raises:
            InvalidAccountException: 账号密码错误
            InvalidCaptchaException: 验证码错误
            UninitializedPasswordException: 未初始化密码
        """
        if resp.status_code != 302:
            content = resp.text
            if "用户名或密码错误,请联系本院教务老师!" in content:
                raise InvalidAccountException()
            if "验证码错误!!" in content:
                raise InvalidCaptchaException()
            raise Exception("登陆失败")
        else:
            location = resp.headers.get("location")
            if location.endswith(XTUEMSConfig.XTU_EMS_UPDATE_PASSWORD_URL):
                raise UninitializedPasswordException()
            return session

    async def async_login(self, account: AuthenticationAccount, retry_time=3) -> Session:
        """
        登陆教务系统
        Args:
            account: 账户信息
            retry_time: 重试次数，如果账号密码错误会直接抛出异常，不会重试

        Returns:
            登陆后的session
        Raises:
            InvalidAccountException: 账号密码错误
            InvalidCaptchaException: 多次尝试仍然验证码错误
            UninitializedPasswordException: 未初始化密码
        """
        err = None
        self.pre_check(account)
        if retry_time <= 0:
            raise Exception("登陆失败")
        for i in range(retry_time):
            try:
                session = await self._async_login(account)
                logger.info(f'登陆成功-{account.username}')
                return session
            except InvalidAccountException as e:
                err = e
                logger.debug('错误的账号密码')
                break
            except InvalidCaptchaException as e:
                err = e
                logger.debug(f'正在重试{i}/{retry_time}次登陆失败-验证码错误')
                continue
            except UninitializedPasswordException as e:
                err = e
                logger.debug(f'[{account.username}]-未初始化密码')
                break
            except TimeoutError as e:
                err = e
                logger.debug(f'正在重试{i}/{retry_time}次登陆失败-网络超时')
                continue
        raise err

    async def _async_login(self, account: AuthenticationAccount) -> Session:
        """
        异步登陆教务系统
        Args:
            account: 账户信息
        Returns:
            登陆后的session
        Raises:
            InvalidAccountException: 账号密码错误
            InvalidCaptchaException: 多次验证码错误
            UninitializedPasswordException: 未初始化密码
        """
        async with ClientSession() as ems_session:
            resp = await ems_session.get(url=XTUEMSConfig.XTU_EMS_CAPTCHA_URL,
                                         timeout=RequestConfig.XTU_EMS_REQUEST_TIMEOUT)
            session_id = resp.cookies.get(QZEducationalManageSystem.SESSION_NAME).value
            session = Session(token=session_id)
            captcha = self.captcha.verify(await resp.read())
            resp = await ems_session.post(url=XTUEMSConfig.XTU_EMS_SIG_URL,
                                          timeout=RequestConfig.XTU_EMS_REQUEST_TIMEOUT)
            text = await resp.text()
            signature = json.loads(text).get("data")
            encoded = self._signature(account.username, account.password, signature)
            data = self._data(username=account.username,
                              password=account.password,
                              encode=encoded,
                              random_code=captcha)
            async with ems_session.post(url=XTUEMSConfig.XTU_EMS_LOGIN_URL,
                                        data=data,
                                        headers=QZEducationalManageSystem._HEADER,
                                        allow_redirects=False,
                                        timeout=RequestConfig.XTU_EMS_REQUEST_TIMEOUT) as resp:
                return await self.async_post_process(resp, session) or session

    def _login(self, account: AuthenticationAccount) -> Session:
        """
        登陆教务系统
        Args:
            account: 账户信息

        Returns:
            登陆后的session

        Raises:
            InvalidAccountException: 账号密码错误
            InvalidCaptchaException: 多次验证码错误
            UninitializedPasswordException: 未初始化密码
        """
        with requests.session() as ems_session:
            resp = ems_session.get(url=XTUEMSConfig.XTU_EMS_CAPTCHA_URL,
                                   timeout=RequestConfig.XTU_EMS_REQUEST_TIMEOUT)
            session_id = resp.cookies.get(QZEducationalManageSystem.SESSION_NAME)
            session = Session(token=session_id)
            captcha = self.captcha.verify(resp.content)
            resp = ems_session.post(url=XTUEMSConfig.XTU_EMS_SIG_URL,
                                    timeout=RequestConfig.XTU_EMS_REQUEST_TIMEOUT)
            signature = json.loads(resp.content).get("data")
            encoded = self._signature(account.username, account.password, signature)
            data = self._data(username=account.username,
                              password=account.password,
                              encode=encoded,
                              random_code=captcha)

            resp = ems_session.post(url=XTUEMSConfig.XTU_EMS_LOGIN_URL,
                                    data=data,
                                    headers=QZEducationalManageSystem._HEADER,
                                    allow_redirects=False,
                                    timeout=RequestConfig.XTU_EMS_REQUEST_TIMEOUT)

            return self.post_process(resp, session) or session

    def _signature(self, username: str, password: str, signature: str):
        """
        智强教务系统签名算法
        Args:
            username: 用户名
            password: 密码
            signature: 加密签名

        Returns:
            返回签名后的密码

        """
        # 将data按照'#'分割成列表
        split = signature.split("#")
        # 创建code字符串
        code = username + "%%%" + password
        # 初始化编码结果的StringBuilder（在Python中使用列表模拟）
        encoded = []
        # 获取code的长度
        length = len(code)
        # 初始化偏移量b
        b = 0
        # 遍历code中的每个字符
        for i in range(length):
            if i < 20:
                # 追加code中的当前字符
                encoded.append(code[i])
                # 追加split[0]中从b开始的由split[1][i]指定数量的字符
                for _ in range(ord(split[1][i]) - ord('0')):
                    encoded.append(split[0][b])
                    b += 1
            else:
                # 追加code从索引20开始的剩余部分
                encoded.extend(code[20:])
                break

        # 将列表转换为字符串并返回
        return ''.join(encoded)
