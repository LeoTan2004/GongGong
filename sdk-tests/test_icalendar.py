import datetime
from unittest import TestCase

from xtu_ems.util.icalendar import BaseCalendar, BaseRepeatRule, BaseEvent, BaseAlarm


class TestBaseCalendar(TestCase):
    def test_to_ical(self):
        target = """\
BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Sky31//Gong 3.0//CN
BEGIN:VTIMEZONE
TZID:Asia/Shanghai
X-LIC-LOCATION:Asia/Shanghai
BEGIN:STANDARD
TZOFFSETFROM:+0800
TZOFFSETTO:+0800
TZNAME:CST
DTSTART:19700101T000000
END:STANDARD
END:VTIMEZONE
BEGIN:VEVENT
DTSTAMP;TZID=Asia/Shanghai:20240101T000000
SUMMARY:Christmas Celebration
DESCRIPTION:Celebrate Christmas all day.
LOCATION:Home
DTSTART;TZID=Asia/Shanghai:20241225T120000
DTEND;TZID=Asia/Shanghai:20241225T123000
CATEGORIES:Holiday
RRULE:FREQ=WEEKLY;COUNT=21
BEGIN:VALARM
ACTION:DISPLAY
DESCRIPTION:Alarm
TRIGGER:-PT22H
END:VALARM
UID:0ca8519470fb2c684ae26a8f010183dd@gong.sky31.com
END:VEVENT
END:VCALENDAR\
"""
        alarm = BaseAlarm(trigger=datetime.timedelta(days=-1, hours=2), description="Alarm")
        event = BaseEvent(
            summary="Christmas Celebration",
            description="Celebrate Christmas all day.",
            location="Home",
            start_time=datetime.datetime(2024, 12, 25, 12, 0, 0),
            end_time=datetime.datetime(2024, 12, 25, 12, 30, 0),
            category="Holiday",
            rrule=BaseRepeatRule(freq="WEEKLY", count=21),
            alarm=alarm,
            dtstamp=datetime.datetime(2024, 1, 1, 0, 0, 0)
        )
        calendar = BaseCalendar()
        calendar.add_event(event)
        ical_content = calendar.to_ical()
        self.assertEqual(target, ical_content)
        print(ical_content)
