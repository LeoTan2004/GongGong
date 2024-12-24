from datetime import datetime

from bs4 import BeautifulSoup

from xtu_ems.ems.config import XTUEMSConfig
from xtu_ems.ems.handler import EMSPoster
from xtu_ems.ems.model import ExamInfo, ExamInfoList


class StudentExamGetter(EMSPoster[ExamInfoList]):
    """获取学生考试信息"""

    def _data(self):
        return {
            "xnxqid": XTUEMSConfig.get_current_term()
        }

    def _extra_info(self, soup: BeautifulSoup):
        exam_list = soup.find(id="dataList").find_all('tr')[1:]
        exam_list = [self._extra_exam_info(row) for row in exam_list]
        return ExamInfoList(exams=exam_list)

    def _extra_exam_info(self, row: BeautifulSoup) -> ExamInfo:
        """从表格的某一行中提取学生的考试信息"""
        tds = row.find_all('td')
        time = tds[5].text.strip().split('~')
        if len(time) == 2:
            start_time, end_time = time
            data = start_time.split(' ')[0]
            end_time = f'{data} {end_time}'
            start_time = datetime.strptime(start_time, '%Y-%m-%d %H:%M')
            end_time = datetime.strptime(end_time, '%Y-%m-%d %H:%M')
        else:
            start_time, end_time = '', ''
        return ExamInfo(name=tds[2].text.strip(),
                        type=tds[3].text.strip(),
                        start_time=start_time,
                        end_time=end_time,
                        location=tds[6].text.strip())

    def url(self):
        return XTUEMSConfig.XTU_EMS_STUDENT_EXAM_URL
