# 基础镜像
FROM python:3.11-slim-bullseye AS base

LABEL authors="leo"

# 设置工作目录
WORKDIR /app

ARG PROXY="https://mirrors.cloud.tencent.com/pypi/simple"

# 安装依赖，仅在 requirements 文件发生变化时重新安装
COPY ./requirements*.txt ./
RUN pip install --upgrade pip -i $PROXY --no-cache\
    && pip config set global.index-url $PROXY \
    && pip install -r requirements.txt --no-cache\
    && rm -f requirements*.txt

# 应用构建阶段
FROM base

LABEL authors="leo"

WORKDIR /app

RUN mkdir /logs

# 创建一个非特权用户，避免以 root 身份运行
RUN useradd -m fastapi_user && \
    chown -R fastapi_user:fastapi_user /app /logs

#RUN chmod -R u+w /logs

# 复制源码文件
COPY ./ems-sdk /app

# 切换到非 root 用户
USER fastapi_user

# 暴露端口
EXPOSE 8080

# 入口命令
ENTRYPOINT ["uvicorn", "app:api", "--host", "0.0.0.0", "--port", "8080", "--log-config", "log_config.json"]
