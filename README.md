# GongGong

[![Static Badge](https://img.shields.io/badge/sky31%20studio-red)](https://github.com/sky31studio) ![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/sky31studio/GongGong/xtu-ems-sdk-test.yml)  ![GitHub commit activity](https://img.shields.io/github/commit-activity/w/sky31studio/GongGong)![GitHub top language](https://img.shields.io/github/languages/top/sky31studio/GongGong)![GitHub License](https://img.shields.io/github/license/sky31studio/GongGong)

拱拱是一个基于网络爬虫的湘潭大学校园APP。本项目是GongGong的后端部分。

## 项目功能

- **网络爬虫**
    - 查询个人课表
    - 查询个人基本信息
    - 查询个人成绩
    - 查询个人排名
    - 查询考试安排
    - 查询空闲教室
    - 查询教学周历

- **平台服务**
    - 登陆账户
  - 获取相关信息
- 反馈服务
    - 反馈数据

## 快速上手

### 平台服务

#### 下载源码

```bash
git clone https://github.com/sky31studio/GongGong.git
```

#### 使用Docker环境启动

切换到根目录，即可

```bash
sudo docker-compose up -d
```

默认端口映射在***8000***端口上，可以通过`http://<host>:<port>/docs`查看接口文档。

我们在***8080***端口上还添加了使用反馈的接口，在`POST http://<host>:<port>/feedback`可以使用。此功能与主服务独立，如果不需要可以在
`docker-compose.yaml`文件中删除该服务。

> [!Important]
> 生产环境下，建议您将该项目的OpenAPI文档关闭。以避免被恶意攻击。
> 你可以通过设置环境变量 `ENV=prod` 来启用生产环境下的服务。

#### 使用Python环境启动

> [!Note]
>
> 开发环境为**Python 3.11**，建议使用**Python 3.10+**的环境。

1. 安装Python依赖包

   ```bash
   pip install -r requirements.txt
   ```

2. 检查端口

   服务将会使用***8000***端口，在启动之前建议先检查一下端口的使用情况，如果端口正在使用嗯，你可以在`ems-sdk/app.py`
   代码中将端口改成一个可用的端口号。

3. 启动服务

   ```bash
   cd ems-sdk && python app.py
   ```

### 网络爬虫SDK

> [!Tip]
>
> 湘潭大学教务系统是由[强智科技](https://www.qzdatasoft.com/)公司开发的大学教务系统，其业务逻辑部分大体相同，如果你想为其他学校的校务系统进行开发，你需要更改
`ems-sdk/xtu_ems/ems/config.py`中的文件。
>
> 如果你在使用的过程中发现了问题，欢迎通过 [Issues](https://github.com/sky31studio/GongGong/issues) 向我们反映。

#### 下载安装

网络爬虫部分代码在 `ems-sdk/xtu_ems` 目录下。我们采用SDK的方式来允许其他人在本项目的基础上进行二次开发。您可以在releases中下载对应版本的SDK并安装。

```shell
pip install ./xtu_ems-**.whl
```

也可以直接将项目**clone**下来

```bash
git clone https://github.com/sky31studio/GongGong.git
```

然后将ems-sdk设置为源码根目录（**PyCharm**）,或者在**PYTHONPATH**中添加该目录。

#### 如何使用

> [!TIP]
>
> - 您需要有一个可用的[湘潭大学校务系统](https://jwxt.xtu.edu.cn/jsxsd)学生账号
>- 您需要可以正常访问[湘潭大学校务系统](https://jwxt.xtu.edu.cn/jsxsd)

在您成功安装版本之后，你便可以通过该SDK获取你在湘潭大学校务系统上的信息。SDK中主要使用到几个概念：**账号、校务系统、
会话、操作**。

```mermaid
flowchart LR

Account((账号))--登陆-->EMS[教务系统]--返回会话-->Session((会话))--基于会话进行操作-->Handle[操作]
```

我们这里以获取基本用户信息为例：

```python
username = "XTU_USERNAME"  # 你的校务系统账号
password = "XTU_PASSWORD"  # 你的教务系统密码

from xtu_ems.ems.account import AuthenticationAccount
from xtu_ems.ems.ems import QZEducationalManageSystem
from xtu_ems.ems.handler.get_student_info import StudentInfoGetter

# 创建一个校务系统账号
account = AuthenticationAccount(username=username,
                                password=password)
ems = QZEducationalManageSystem()
# 登陆校务系统
session = ems.login(account)

handler = StudentInfoGetter()
# 执行校务系统爬虫操作
resp = handler.handler(session)
print(resp.model_dump_json(indent=4))
```

我们也提供了**异步执行**的方式，你可以利用异步函数加速代码的执行效率（这在多个任务执行时极大的体现出差异）

```python
username = "XTU_USERNAME"  # 你的校务系统账号
password = "XTU_PASSWORD"  # 你的教务系统密码

from xtu_ems.ems.account import AuthenticationAccount
from xtu_ems.ems.ems import QZEducationalManageSystem
from xtu_ems.ems.handler.get_student_info import StudentInfoGetter

# 创建一个校务系统账号
account = AuthenticationAccount(username=username,
                                password=password)
ems = QZEducationalManageSystem()
handler = StudentInfoGetter()


async def main():
    # 登陆校务系统
    session = await ems.async_login(account)
    # 执行校务系统爬虫操作
    resp = await handler.async_handler(session)
    print(resp.model_dump_json(indent=4))


if __name__ == '__main__':
    import asyncio

    asyncio.run(main())
```

