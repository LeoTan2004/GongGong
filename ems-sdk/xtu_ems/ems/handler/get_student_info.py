from bs4 import BeautifulSoup

from xtu_ems.ems.config import XTUEMSConfig
from xtu_ems.ems.handler import EMSGetter
from xtu_ems.ems.model import StudentBasicInfo


def _extra_student_info(soup: BeautifulSoup) -> StudentBasicInfo:
    """从表格中提取学生的基本信息"""
    student_info = {}
    rows = soup.find(id='xjkpTable').find_all('tr')
    # 提取学生的基本信息
    for row in rows:
        skip = 0
        for cell in row.find_all('td'):
            if cell.text.__contains__('：'):
                key, value = cell.text.split('：')
            else:
                if skip == 1:
                    skip = 0
                    continue
                key = cell.text or ''
                value = cell.find_next('td')
                value = value.text if value else ''
                skip = 1
            if key not in student_info:
                student_info[key] = value
    info = {
        'student_id': student_info.get('学号').strip(),
        'name': student_info.get('姓名').strip(),
        'college': student_info.get('院系', '').strip(),
        'major': student_info.get('专业', '').strip(),
        'class_': student_info.get('班级', '').strip(),
        'gender': student_info.get('性别', '').strip(),
        'birthday': student_info.get('出生日期', '').strip(),
        'entrance_day': student_info.get('入学日期', '').strip(),
    }
    return StudentBasicInfo(**info)


class StudentInfoGetter(EMSGetter[StudentBasicInfo]):
    """获取学生信息"""

    def url(self):
        return XTUEMSConfig.XTU_EMS_STUDENT_INFO_URL

    def _extra_info(self, soup: BeautifulSoup):
        """从表格中提取学生的基本信息"""
        return _extra_student_info(soup)
