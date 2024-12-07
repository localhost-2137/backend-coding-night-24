package main

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"os"
	"sort"
	"strconv"
)

type nasaWhetherT struct {
	LastUTC  string `json:"Last_UTC"`
	FirstUTC string `json:"First_UTC"`
	AT       struct {
		Av float64 `json:"av"`
		Mn float64 `json:"mn"`
		Mx float64 `json:"mx"`
	} `json:"AT"`
	PRE struct {
		Av float64 `json:"av"`
		Mn float64 `json:"mn"`
		Mx float64 `json:"mx"`
	} `json:"PRE"`
	HWS struct {
		Av float64 `json:"av"`
		Mn float64 `json:"mn"`
		Mx float64 `json:"mx"`
	} `json:"HWS"`
}

type chartDataT struct {
	Date     string  `json:"date"`
	TempAvg  float64 `json:"temp_avg"`
	Pressure float64 `json:"pressure"`
	Wind     float64 `json:"wind"`
}

var nasaWhetherData map[string]nasaWhetherT
var nasaWhetherChartData []chartDataT

func initNasaData() error {
	//whetherDataRaw, err := http.Get("https://api.nasa.gov/insight_weather/?api_key=DEMO_KEY&feedtype=json&ver=1.0")
	//if err != nil {
	//	return err
	//}

	// open form file:
	whetherDataRaw, err := os.Open("nasa_whether.json")
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(whetherDataRaw); err != nil {
		return err
	}

	if err := whetherDataRaw.Close(); err != nil {
		return err
	}

	if err := json.Unmarshal(buf.Bytes(), &nasaWhetherData); err != nil {
		return err
	}

	dataKeys := make([]string, 0, len(nasaWhetherData))
	for key := range nasaWhetherData {
		dataKeys = append(dataKeys, key)
	}
	sort.Strings(dataKeys)

	nasaWhetherChartData = make([]chartDataT, 0, len(dataKeys))
	for _, key := range dataKeys {
		day := nasaWhetherData[key]
		nasaWhetherChartData = append(nasaWhetherChartData, chartDataT{
			Date:     day.LastUTC,
			TempAvg:  day.AT.Av,
			Pressure: day.PRE.Av,
			Wind:     day.HWS.Av,
		})
	}

	return nil
}

func statsEndpoint(ctx *fiber.Ctx) error {
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		limit = len(nasaWhetherChartData)
	}

	if limit > len(nasaWhetherChartData) {
		limit = len(nasaWhetherChartData)
	}
	
	data := nasaWhetherChartData[len(nasaWhetherChartData)-limit:]
	return ctx.JSON(data)
}
