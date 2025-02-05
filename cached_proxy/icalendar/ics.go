package icalendar

import (
	"fmt"
	"strings"
	"time"
)

type Timezone interface {
	Component
	// SetID 设置时区的 ID
	SetID(id string)
	// SetOffsetFrom 设置时区的偏移
	SetOffsetFrom(offset time.Duration)
	// SetOffsetTo 设置时区的偏移
	SetOffsetTo(offset time.Duration)
	// SetName 设置时区的名称
	SetName(name string)
	// SetStart 设置时区的开始时间
	SetStart(start time.Time)
	// GetID 获取时区的 ID
	GetID() string
}

// Component 是 ICS 文件中的组件。
type Component interface {
	// ToIcs 将组件转换为 ICS 格式
	ToIcs(timezone *Timezone) string
}
type Action string

const (
	DISPLAY = "DISPLAY"
	AUDIO   = "AUDIO"
)

// Alarm 是 ICS 文件中的提醒。
type Alarm interface {
	Component
	// SetAction 设置提醒的动作
	SetAction(action Action)
	// SetTrigger 设置提醒的触发时间
	SetTrigger(trigger time.Duration)
	// SetDescription 设置提醒的描述
	SetDescription(description string)
}

// RepeatRule 是 ICS 文件中的重复规则。
type RepeatRule interface {
	Component
	// SetFrequency 设置重复规则的频率
	SetFrequency(frequency string)
	// SetInterval 设置重复规则的间隔
	SetInterval(interval int)
	// SetCount 设置重复规则的次数
	SetCount(count int)
	// SetUntil 设置重复规则的结束时间
	SetUntil(until time.Time)
}

type Location interface {
	Component
	// SetName 设置地点的名称
	SetName(name string)
}

// Event 是 ICS 文件中的事件。
type Event interface {
	Component
	// SetSummary 设置事件的摘要
	SetSummary(summary string)
	// SetDescription 设置事件的描述
	SetDescription(description string)
	// SetLocation 设置事件的地点
	SetLocation(location Location)
	// SetStart 设置事件的开始时间
	SetStart(start time.Time)
	// SetEnd 设置事件的结束时间
	SetEnd(end time.Time)
	// AddAlarm 添加一个提醒
	AddAlarm(alarm Alarm)
	// SetRepeatRule 设置事件的重复规则
	SetRepeatRule(rule RepeatRule)
}

type Calendar interface {
	Component
	// AddEvent 添加一个事件
	AddEvent(event Event)
	// SetProductID 设置日历的产品 ID
	SetProductID(productID string)
	// SetTimezone 设置日历的时区
	SetTimezone(timezone Timezone)
}

func TimeToIcs(t time.Time, timezone *Timezone, sep string) string {
	ts := t.Format("20060102T150405")
	if timezone == nil {
		return fmt.Sprintf("%s%sZ", sep, ts)
	}
	return fmt.Sprintf(";TZID=%s%s%s", (*timezone).GetID(), sep, ts)
}

func DurationToIcs(d time.Duration) string {
	if d == 0 {
		return "PT0S"
	}
	result := strings.Builder{}
	if d < 0 {
		result.WriteRune('-')
		d = d.Abs()
	}
	result.WriteString("P")
	days := int(d / (24 * time.Hour))
	hours := int((d % (24 * time.Hour)) / time.Hour)
	minutes := int((d % time.Hour) / time.Minute)
	seconds := int((d % time.Minute) / time.Second)
	if days > 0 {
		result.WriteString(fmt.Sprintf("%dD", days))
	}
	if hours > 0 || minutes > 0 || seconds > 0 {
		result.WriteString("T")
		if hours > 0 {
			result.WriteString(fmt.Sprintf("%dH", hours))
		}
		if minutes > 0 {
			result.WriteString(fmt.Sprintf("%dM", minutes))
		}
		if seconds > 0 {
			result.WriteString(fmt.Sprintf("%dS", seconds))
		}
	}
	return result.String()
}

func OffsetToIcs(d time.Duration) string {
	if d == 0 {
		return "Z"
	}
	return fmt.Sprintf("+%02d00", int(d.Hours()))
}
