package icalendar

type IcsLocation struct {
	name string
}

func (l *IcsLocation) ToIcs(_ *Timezone) string {
	return "LOCATION:" + l.name
}

func (l *IcsLocation) SetName(name string) {
	l.name = name
}
