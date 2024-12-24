import datetime
import hashlib
from abc import abstractmethod
from dataclasses import dataclass
from typing import Optional, Union, Literal


class Component:
    """ics组件"""

    @abstractmethod
    def to_ical(self) -> str:
        """转化成ics格式中的组件"""
        pass


@dataclass
class BaseAlarm(Component):
    """基础提醒"""
    action: str = "DISPLAY"
    """提醒动作， ACTION"""
    trigger: Union[datetime.timedelta, None] = None
    """提醒时间， TRIGGER"""
    description: str = ""
    """提醒描述， DESCRIPTION"""

    def to_ical(self) -> str:
        """转化成ics格式中的提醒"""
        lines = ["BEGIN:VALARM", f"ACTION:{self.action}"]
        # DESCRIPTION
        if self.description:
            lines.append(f"DESCRIPTION:{self.description}")
        # TRIGGER
        if self.trigger:
            # 处理相对时间
            total_seconds = int(self.trigger.total_seconds())
            is_negative = total_seconds < 0  # 判断是否为负时间
            total_seconds = abs(total_seconds)

            # 按优先级选择最简洁的单位
            days = total_seconds // 86400
            hours = (total_seconds % 86400) // 3600
            minutes = (total_seconds % 3600) // 60

            if days > 0:
                trigger_value = f"-P{days}D" if is_negative else f"P{days}D"
            elif hours > 0:
                trigger_value = f"-PT{hours}H" if is_negative else f"PT{hours}H"
            else:
                trigger_value = f"-PT{minutes}M" if is_negative else f"PT{minutes}M"

            lines.append(f"TRIGGER:{trigger_value}")

        lines.append("END:VALARM")
        return "\n".join(lines)

    duration: Union[datetime.timedelta, None] = None


@dataclass
class BaseRepeatRule(Component):
    """基础重复规则"""
    freq: Literal['SECONDLY', 'MINUTELY', 'HOURLY', 'DAILY', 'WEEKLY', 'MONTHLY', 'YEARLY'] = "WEEKLY"
    """重复频率， FREQ"""
    interval: int = 1
    """重复间隔， INTERVAL"""
    count: int = 0
    """重复次数， COUNT"""
    until: Union[datetime.date, datetime.datetime, None] = None
    """结束时间， UNTIL"""

    def to_ical(self) -> str:
        """转化成ics格式中的重复规则"""
        if not self.freq:
            raise ValueError("The 'freq' field is required for a repeat rule.")

        rrule_parts = [f"FREQ={self.freq.upper()}"]

        # Add INTERVAL if it's greater than 1
        if self.interval > 1:
            rrule_parts.append(f"INTERVAL={self.interval}")

        # Add COUNT if specified and greater than 0
        if self.count > 0:
            rrule_parts.append(f"COUNT={self.count}")

        # Add UNTIL if specified
        if self.until:
            if isinstance(self.until, datetime.datetime):
                until_str = self.until.strftime("%Y%m%dT%H%M%S")
            elif isinstance(self.until, datetime.date):
                until_str = self.until.strftime("%Y%m%d")
            else:
                raise ValueError("The 'until' field must be a date or datetime object.")
            rrule_parts.append(f"UNTIL={until_str}")

        # Join all parts to form the RRULE
        return ";".join(rrule_parts)


@dataclass
class BaseEvent(Component):
    """基础事件"""
    summary: str = ""
    """事件标题， SUMMARY"""
    description: str = ""
    """描述， DESCRIPTION"""
    location: Optional[str] = None
    """地点， LOCATION"""
    start_time: Union[datetime.date, datetime.datetime, None] = None
    """开始时间， DTSTART"""
    end_time: Union[datetime.date, datetime.datetime, None] = None
    """结束时间， DTEND"""
    category: str = ""
    """分类， CATEGORIES"""
    rrule: Union[BaseRepeatRule, str] = ""
    """重复规则， RRULE"""
    dtstamp: datetime.datetime = datetime.datetime.now()
    """创建时间， DTSTAMP"""

    alarm: Union[BaseAlarm, list[BaseAlarm]] = ""
    """提醒， ALARM"""

    def to_ical(self) -> str:
        """转化成ics格式中的事件"""
        lines = ["BEGIN:VEVENT", f"DTSTAMP;TZID=Asia/Shanghai:{self.dtstamp.strftime('%Y%m%dT%H%M%S')}"]

        # SUMMARY
        if self.summary:
            lines.append(f"SUMMARY:{self.summary}")

        # DESCRIPTION
        if self.description:
            lines.append(f"DESCRIPTION:{self.description}")

        # LOCATION
        if self.location:
            lines.append(f"LOCATION:{self.location}")

        # DTSTART
        if self.start_time:
            if isinstance(self.start_time, datetime.date) and not isinstance(self.start_time, datetime.datetime):
                # All-day event
                lines.append(f"DTSTART;VALUE=DATE:{self.start_time.strftime('%Y%m%d')}")
            else:
                # Date-time event
                lines.append(f"DTSTART;TZID=Asia/Shanghai:{self.start_time.strftime('%Y%m%dT%H%M%S')}")

        # DTEND
        if self.end_time:
            if isinstance(self.end_time, datetime.date) and not isinstance(self.end_time, datetime.datetime):
                # All-day event
                lines.append(f"DTEND;VALUE=DATE:{self.end_time.strftime('%Y%m%d')}")
            else:
                # Date-time event
                lines.append(f"DTEND;TZID=Asia/Shanghai:{self.end_time.strftime('%Y%m%dT%H%M%S')}")

        # CATEGORIES
        if self.category:
            lines.append(f"CATEGORIES:{self.category}")

        # RRULE
        if self.rrule:
            if isinstance(self.rrule, BaseRepeatRule):
                lines.append(f"RRULE:{self.rrule.to_ical()}")
            else:
                lines.append(f"RRULE:{self.rrule}")

        # ALARM
        if self.alarm:
            if isinstance(self.alarm, BaseAlarm):
                lines.append(f"{self.alarm.to_ical()}")
            else:
                for alarm in self.alarm:
                    lines.append(f"{alarm.to_ical()}")

        # Add unique identifier and timestamp
        hash_value = hashlib.md5((self.summary + self.start_time.__str__()).encode()).hexdigest()
        uid = f"{hash_value}@gong.sky31.com"
        lines.append(f"UID:{uid}")

        lines.append("END:VEVENT")
        return "\n".join(lines)


class BaseCalendar:
    """基础日历"""
    events: list[BaseEvent] = []
    """事件列表， EVENTS"""
    prodid: str = "-//Sky31//Gong 3.0//CN"

    timezone = """\
BEGIN:VTIMEZONE
TZID:Asia/Shanghai
X-LIC-LOCATION:Asia/Shanghai
BEGIN:STANDARD
TZOFFSETFROM:+0800
TZOFFSETTO:+0800
TZNAME:CST
DTSTART:19700101T000000
END:STANDARD
END:VTIMEZONE"""

    def add_event(self, event: BaseEvent):
        """添加事件"""
        self.events.append(event)

    def to_ical(self) -> str:
        """转化成ics格式中的日历"""
        lines = ["BEGIN:VCALENDAR", "VERSION:2.0", f"PRODID:{self.prodid}", self.timezone]
        # Add events
        [lines.append(e.to_ical()) for e in self.events]

        lines.append("END:VCALENDAR")
        return "\n".join(lines)
