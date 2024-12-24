"""链接会话模块，用于存储会话信息"""
from pydantic import BaseModel


class Session(BaseModel):
    """会话类，用于存储会话信息"""
    token: str
