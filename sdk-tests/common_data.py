import os

username = os.getenv("XTU_USERNAME")
password = os.getenv("XTU_PASSWORD")
from xtu_ems.ems.account import AuthenticationAccount
from xtu_ems.ems.ems import QZEducationalManageSystem

account = AuthenticationAccount(username=username,
                                password=password)
ems = QZEducationalManageSystem()
session = ems.login(account)
