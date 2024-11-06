from plat.repository.d_basic import KVRepository, SimpleKVRepository
from plat.service.entity import Account, AccountStatus
from xtu_ems.ems.account import AuthenticationAccount
from xtu_ems.ems.ems import QZEducationalManageSystem


class ExpiredAccountException(Exception):
    """账户已过期"""

    def __init__(self, username):
        self.username = username


class BannedAccountException(Exception):
    """账户已被封禁"""

    def __init__(self, username):
        self.username = username


class AccountService:
    ems = QZEducationalManageSystem()

    def __init__(self, account_repository: KVRepository[str, Account],
                 token_repository: KVRepository[str, Account] = SimpleKVRepository()):
        """
        账户服务类
        Args:
            account_repository: 账户存储库
            token_repository: token存储库
        """
        self.account_repository = account_repository
        self.token_repository = token_repository

    async def login(self, username: str, password: str):
        """
        登陆教务系统

        当登陆校务系统后，之前的token会被覆盖，并且返回一个新的token
        Args:
            username: 用户名
            password: 密码

        Returns:

        """
        account = AuthenticationAccount(username=username, password=password)
        session = await self.ems.async_login(account)
        authed_account = await self.account_repository.async_get_item(username)
        if authed_account:
            authed_account.session = session.session_id
            authed_account.status = AccountStatus.NORMAL
            token = authed_account.token
            await self.token_repository.async_del_item(token)
            authed_account.refresh_token()
        else:
            authed_account = Account(student_id=username,
                                     password=password,
                                     session=session.session_id,
                                     status=AccountStatus.NORMAL)
        authed_account = await self.save_new_account(authed_account)
        return authed_account

    async def save_new_account(self, account: Account):
        """
        保存新的用户

        Args:
            account: 用户信息

        Returns:
            用户信息
        """
        while await self.token_repository.async_get_item(account.token):
            account.refresh_token()
        await self.token_repository.async_set_item(account.token, account)
        await self.account_repository.async_set_item(account.student_id, account)
        return account

    async def auth_with_token(self, token: str):
        """
        用token验证用户
        Args:
            token: 用户凭证

        Returns:

        """
        account = await self.token_repository.async_get_item(token)
        if account is not None:
            if account.status == AccountStatus.EXPIRED:
                raise ExpiredAccountException(account.student_id)
            elif account.status == AccountStatus.BANNED:
                raise BannedAccountException(account.student_id)
            else:
                return account
        else:
            return None

    async def expire_account(self, username: str):
        """
        标记过期用户

        Args:
            username: 用户名

        Returns:

        """
        account = await self.account_repository.async_get_item(username)
        if account:
            account.status = AccountStatus.EXPIRED
            await self.account_repository.async_set_item(username, account)
            return account
        else:
            return None
