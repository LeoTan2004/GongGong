"""账号模块"""
from abc import ABC


class Account(ABC):
    """账号类基类"""

    def __init__(self, user_id: str) -> None:
        self.user_id = user_id

    def __str__(self) -> str:
        return f"username: {self.user_id}"


class AuthenticationAccount(Account):
    """可鉴权账号类"""

    def __init__(self, username: str, password: str) -> None:
        super().__init__(user_id=username)
        self.password = password

    @property
    def username(self) -> str:
        return self.user_id

    @username.setter
    def username(self, username: str) -> None:
        self.user_id = username


class StudentAccount(AuthenticationAccount):
    """学生账号类"""

    def __init__(self, student_id: str, password: str) -> None:
        super().__init__(student_id, password)

    @property
    def student_id(self) -> str:
        return self.user_id

    @student_id.setter
    def student_id(self, student_id: str) -> None:
        self.user_id = student_id
