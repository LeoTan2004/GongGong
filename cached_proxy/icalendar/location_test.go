package icalendar

import "testing"

func TestIcsLocation_ToIcs(t *testing.T) {
	type fields struct {
		name string
	}
	type args struct {
		in0 *Timezone
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Test IcsLocation ToIcs ",
			fields: fields{
				name: "name",
			},
			args: args{
				in0: nil,
			},
			want: "LOCATION:name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &IcsLocation{
				name: tt.fields.name,
			}
			if got := l.ToIcs(tt.args.in0); got != tt.want {
				t.Errorf("ToIcs() = %v, want %v", got, tt.want)
			}
		})
	}
}
