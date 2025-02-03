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

func TestGeoLocation_SetName(t *testing.T) {
	type fields struct {
		IcsLocation IcsLocation
		latitude    float64
		longitude   float64
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test GeoLocation SetName",
			fields: fields{
				IcsLocation: IcsLocation{
					name: "name",
				},
			},
			args: args{
				name: "北京",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GeoLocation{
				IcsLocation: tt.fields.IcsLocation,
				latitude:    tt.fields.latitude,
				longitude:   tt.fields.longitude,
			}
			g.SetName(tt.args.name)
		})
	}
}

func TestGeoLocation_ToIcs(t *testing.T) {
	type fields struct {
		IcsLocation IcsLocation
		latitude    float64
		longitude   float64
		refreshed   bool
	}
	type args struct {
		t *Timezone
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Test GeoLocation ToIcs",
			fields: fields{
				IcsLocation: IcsLocation{
					name: "湘潭大学碧泉书院",
				},
				refreshed: false,
			},
			args: args{
				t: nil,
			},
			want: "LOCATION:湘潭大学碧泉书院\nGEO:39.883451;116.196375",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GeoLocation{
				IcsLocation: tt.fields.IcsLocation,
				latitude:    tt.fields.latitude,
				longitude:   tt.fields.longitude,
				refreshed:   tt.fields.refreshed,
			}
			if got := g.ToIcs(tt.args.t); got != tt.want {
				t.Errorf("ToIcs() = %v, want %v", got, tt.want)
			}
		})
	}
}
