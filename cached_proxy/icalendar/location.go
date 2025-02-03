package icalendar

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IcsLocation struct {
	name string
}

func (l *IcsLocation) ToIcs(_ *Timezone) string {
	return "LOCATION:" + l.name
}

func (l *IcsLocation) SetName(name string) {
	l.name = name
}

type GeoLocation struct {
	IcsLocation
	refreshed bool
	latitude  float64
	longitude float64
}

func (g *GeoLocation) SetName(name string) {
	g.IcsLocation.SetName(name)
	g.refreshed = false
}

func (g *GeoLocation) refresh() {
	name := g.name
	resp, err := http.Get("http://api.map.baidu.com/geocoder?address=" + name + "&output=json")
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	//	 json 解析
	decoder := json.NewDecoder(resp.Body)
	var data map[string]interface{}
	err = decoder.Decode(&data)
	if err != nil {
		return
	}
	lat := data["result"].(map[string]interface{})["location"].(map[string]interface{})["lat"]
	lng := data["result"].(map[string]interface{})["location"].(map[string]interface{})["lng"]
	g.latitude = lat.(float64)
	g.longitude = lng.(float64)
	g.refreshed = true
}

func (g *GeoLocation) ToIcs(t *Timezone) string {
	if !g.refreshed {
		g.refresh()
	}
	return g.IcsLocation.ToIcs(t) + fmt.Sprintf("\nGEO:%f;%f", g.latitude, g.longitude)
}
