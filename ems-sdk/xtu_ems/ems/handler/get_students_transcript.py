from functools import cache
from io import BytesIO

from bs4 import BeautifulSoup
from pdfplumber import PDF

from xtu_ems.ems.config import XTUEMSConfig, RequestConfig
from xtu_ems.ems.handler import Handler, _R, EMSPoster, logger
from xtu_ems.ems.model import ScoreBoard, Score, RankInfo
from xtu_ems.ems.session import Session

_data = {
    "xs0101id": "",
    "cjtype": "",
    "sjlj": ["110", "120"],
    "sjcj": "2",
    "bblx": "all"
}


def pre_proc(res: str):
    if not isinstance(res, str):
        return res
    return res.replace('\n', '').replace('\r', '').replace(' ', '').replace('\t', '')


def extract_field(text, start_key, end_key=None):
    start = text.find(start_key) + len(start_key)
    if end_key:
        end = text.find(end_key, start)
        return text[start:end].strip()
    return text[start:].strip()


class StudentTranscriptGetter(Handler[ScoreBoard]):
    """通过教务系统获取成绩单，并且解析成结构化数据"""

    async def async_handler(self, session: Session, *args, **kwargs) -> _R:

        async with self.get_async_session(session) as ems_session:
            logger.debug(f'[{self.__class__.__name__}] 正在异步获取数据-{self.url()}')
            resp = await ems_session.post(url=self.url(), data=_data, timeout=RequestConfig.XTU_EMS_REQUEST_TIMEOUT,
                                          allow_redirects=False)
            if resp.status == 200:
                pdf = PDF(BytesIO(await resp.content.read()))
                return self._extra_info(pdf)

    def handler(self, session: Session, *args, **kwargs):
        with self.get_session(session) as ems_session:
            logger.debug(f'[{self.__class__.__name__}] 正在获取数据-{self.url()}')
            resp = ems_session.post(url=self.url(), data=_data, timeout=RequestConfig.XTU_EMS_REQUEST_TIMEOUT,
                                    allow_redirects=False)
            if resp.status_code == 200:
                pdf = PDF(BytesIO(resp.content))
                return self._extra_info(pdf)

    def url(self):
        return XTUEMSConfig.XTU_EMS_STUDENT_TRANSCRIPT_URL

    def _extra_info(self, pdf):
        page = pdf.pages[0]
        text: str = page.extract_text_lines()[1]['text']

        scoreboard = ScoreBoard()

        scoreboard.college = extract_field(text, '院系：', '专业：')

        scoreboard.major = extract_field(text, '专业：', '姓名：')

        scoreboard.name = extract_field(text, '姓名：', '学号：')

        scoreboard.student_id = extract_field(text, '学号：')

        table = page.extract_table()

        for row in table:

            if not isinstance(row[0], str) or row[0] == "课程名称":
                continue
            if not row[1]:
                self._parse_score(scoreboard, row[0])
                continue
            s = Score(name=pre_proc(row[0]),
                      type=row[1],
                      credit=pre_proc(row[2]),
                      score=pre_proc(row[3]),
                      term=int(pre_proc(row[4])))
            scoreboard.scores.append(s)
            if row[5]:
                s = Score(name=pre_proc(row[5]),
                          type=row[6],
                          credit=pre_proc(row[7]),
                          score=pre_proc(row[8]),
                          term=int(pre_proc(row[9])))
                scoreboard.scores.append(s)
        return scoreboard

    def _parse_score(self, scoreboard: ScoreBoard, detail: str):
        pieces = detail.split(' ')
        total = {}
        for i, piece in enumerate(pieces):
            chunks = piece.split('：')
            if len(chunks) == 1:
                pieces[i + 1] = chunks[0] + pieces[i + 1]
                continue
            if len(chunks) == 3:
                total[chunks[1]] = chunks[2]
                continue
            k, vs = chunks
            total[k.strip()] = vs.strip()
        for k, v in total.items():
            match k:
                case '总学分要求':
                    _, h = scoreboard.total_credit
                    scoreboard.total_credit = (v, h)
                case '已修总学分':
                    t, _ = scoreboard.total_credit
                    scoreboard.total_credit = (t, v)
                case '必修学分要求':
                    _, h = scoreboard.compulsory_credit
                    scoreboard.compulsory_credit = (v, h)
                case '已修必修学分':
                    t, _ = scoreboard.compulsory_credit
                    scoreboard.compulsory_credit = (t, v)
                case '选修学分要求':
                    _, h = scoreboard.elective_credit
                    scoreboard.elective_credit = (v, h)
                case '已修选修学分':
                    t, _ = scoreboard.elective_credit
                    scoreboard.elective_credit = (t, v)
                case '跨学科选修学分要求':
                    _, h = scoreboard.cross_course_credit
                    scoreboard.cross_course_credit = (v, h)
                case '跨学科选修学分':
                    t, _ = scoreboard.cross_course_credit
                    scoreboard.cross_course_credit = (t, v)
                case '平均学分绩点':
                    scoreboard.gpa = v
                case '平均成绩':
                    scoreboard.average_score = v
                case 'CET4':
                    scoreboard.cet4 = v
                case 'CET6':
                    scoreboard.cet6 = v
                case '辅修专业学士学位学分要求':
                    _, h = scoreboard.total_credit
                    scoreboard.total_credit = (v, h)


@cache
def get_all_terms(end_term, start_term=1998) -> list[str]:
    terms = []
    while True:
        term = f"{start_term}-{start_term + 1}-1"
        terms.append(term)
        if term == end_term:
            break
        term = f"{start_term}-{start_term + 1}-2"
        terms.append(term)
        if term == end_term:
            break
        start_term += 1
    return terms


class StudentRankGetter(EMSPoster[RankInfo]):
    def __init__(self, terms=None):
        """
        获取学生排名

        Args:
            terms: 学期列表，默认为获取所有学期
        """
        super().__init__()
        self._terms: list[str] = terms

    @property
    def terms(self):
        if self._terms is None:
            return get_all_terms(XTUEMSConfig.get_current_term())
        return self._terms

    def _data(self):
        return {
            'kksj': self.terms,
            'kclb': [1, 7],
            'zsb': 0
        }

    def url(self):
        return XTUEMSConfig.XTU_EMS_STUDENT_RANK_URL

    def _extra_info(self, soup: BeautifulSoup):
        trs = soup.find_all('tr')
        if len(trs) < 3:
            return []
        tr = trs[2]
        return self._extra_rank_info(tr)

    def _extra_rank_info(self, tr: BeautifulSoup):
        tds = tr.find_all('td')
        d = {
            'gpa': tds[0].text.strip(),
            'avg_score': tds[1].text.strip(),
            'class_rank': tds[2].text.strip(),
            'school_rank': tds[3].text.strip(),
        }
        return RankInfo(
            average_score=d['avg_score'],
            gpa=d['gpa'],
            class_rank=int(d['class_rank']),
            major_rank=int(d['school_rank']),
            terms=self._terms or ['*']
        )


class StudentRankGetterForCompulsory(StudentRankGetter):
    """
    获取学生排名，仅获取必修课
    """

    def _data(self):
        return {
            'kksj': self.terms,
            'kclb': [1],
            'zsb': 0
        }


class StudentTranscriptGetterForAcademicMinor(StudentTranscriptGetter):
    """
    获取学生成绩单，仅获取必修课
    """

    def _data(self):
        return {
            'sjlj': ['110', '120'],
            'sjcj': '2',
            'bblx': 'all'}

    def url(self):
        return XTUEMSConfig.XTU_EMS_STUDENT_TRANSCRIPT_MINOR_URL
