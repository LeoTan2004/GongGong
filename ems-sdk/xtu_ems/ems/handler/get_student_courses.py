from bs4 import BeautifulSoup, Tag

from xtu_ems.ems.config import XTUEMSConfig
from xtu_ems.ems.handler import EMSPoster
from xtu_ems.ems.model import CourseInfo, CourseTable, _get_day_name, CourseList


class StudentCourseGetter(EMSPoster[CourseList]):

    def _data(self):
        return {'xnxq01id': XTUEMSConfig.get_current_term()}

    def url(self):
        return XTUEMSConfig.XTU_EMS_STUDENT_COURSE_URL

    def _extra_info(self, soup: BeautifulSoup):
        class_table = soup.find(id="kbtable")
        return self._extra_student_courses(class_table).to_list()

    def _extra_student_courses(self, class_table: BeautifulSoup):
        try:
            table = CourseTable()
            for time, row in enumerate(class_table.find_all("tr")):
                for week, td in enumerate(row.find_all("td")):

                    if td.text.strip() == "":
                        courses = []
                    else:
                        course_content = td.find(class_="kbcontent")
                        if course_content is None:
                            continue
                        courses = self._extra_courses(course_content, day=week, start=(time - 1) * 2 + 1)
                    table[week].append(courses)
            return table
        except Exception as e:
            raise e

    def _extra_courses(self, td: Tag, course_name=None, day=0, start=1) -> list[CourseInfo]:
        """提起某一天的课程信息"""
        course_name = course_name or td.contents[0].text
        teacher = td.find_next(title='老师')
        weeks = teacher.find_next(title='周次(节次)')
        classroom = weeks.find_next(title='教室') or Tag(name='p')
        duration = 2
        for i, c in enumerate(classroom.next_siblings):
            if i < 2:
                if '上课节次' in c.text:
                    duration = int(c.text.split('：')[1].split('节')[0])
                    break
            else:
                break

        course = CourseInfo(name=course_name.strip(),
                            teacher=teacher.text.strip(),
                            weeks=weeks.text.split('(')[0].strip(),
                            classroom=classroom.text.strip(),
                            start_time=start,
                            duration=duration,
                            day=_get_day_name(day))
        next_c = classroom.find_next_sibling(string='---------------------')
        ret = [course]
        if next_c is not None:
            next_course_name = next_c.next_sibling.next_sibling.text
            ret += self._extra_courses(next_c, next_course_name, day=day, start=start)
        return ret
