from fastapi.testclient import TestClient

from app import api

client = TestClient(app=api)


def test_login():
    """测试正常登陆"""
    from common_data import username, password
    result = client.post("/login", json={"username": username, "password": password})
    assert result.status_code == 200
    assert result.json()["code"] == 1


def test_login_with_wrong_password():
    """测试错误密码登陆"""
    from common_data import username
    result = client.post("/login", json={"username": username, "password": "wrong_password"})
    assert result.status_code == 200
    assert result.json()["code"] == 0


def test_get_courses():
    """测试获取课程表"""
    from common_data import session
    result = client.get("/courses", headers={'token': session.token})
    assert result.status_code == 200
    assert result.json()["code"] == 1


def test_get_info():
    """测试获取基本信息"""
    from common_data import session
    result = client.get("/info", headers={'token': session.token})
    assert result.status_code == 200
    assert result.json()["code"] == 1


def test_get_scores():
    """测试获取主修成绩"""
    from common_data import session
    result = client.get("/scores", headers={'token': session.token})
    assert result.status_code == 200
    assert result.json()["code"] == 1


def test_get_minor_scores():
    """测试获取辅修成绩"""
    from common_data import session
    result = client.get("/minor/scores", headers={'token': session.token})
    assert result.status_code == 200
    assert result.json()["code"] == 1


def test_get_major_rank():
    """测试获取主修排名"""
    from common_data import session
    result = client.get("/rank", headers={'token': session.token})
    assert result.status_code == 200
    assert result.json()["code"] == 1


def test_get_today_classroom():
    """测试获取今天教室"""
    from common_data import session
    result = client.get("/classroom/today", headers={'token': session.token})
    assert result.status_code == 200
    assert result.json()["code"] == 1


def test_get_exams():
    """测试获取考试安排"""
    from common_data import session
    result = client.get("/exams", headers={'token': session.token})
    assert result.status_code == 200
    assert result.json()["code"] == 1


def test_get_tomorrow_classroom():
    """测试获取明天教室"""
    from common_data import session
    result = client.get("/classroom/tomorrow", headers={'token': session.token})
    assert result.status_code == 200
    assert result.json()["code"] == 1


def test_get_calendar():
    """测试获取教学周历"""
    from common_data import session
    result = client.get("/calendar", headers={'token': session.token})
    assert result.status_code == 200
    assert result.json()["code"] == 1
