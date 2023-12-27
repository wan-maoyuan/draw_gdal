package main

import (
	"fmt"
	"image"
	"image/color"

	gdaldraw "gitee.com/wuxi_jiufang/gdal_draw"
	"github.com/batchatco/go-native-netcdf/netcdf"
	"github.com/sirupsen/logrus"
)

const NcPath = "./nc_files/H2B.nc"

func main() {
	data, err := ReadDataFromNcFile(NcPath)
	if err != nil {
		logrus.Errorf("ReadDataFromNcFile: %v", err)
	}

	minWs := float64(0)
	maxWs := float64(25)

	gdaldraw.DrawIrregular3857(data, func(img *image.RGBA, x, y int, value float64) {
		var colorValue uint8
		if value == 0 {
			colorValue = 0
		} else {
			colorValue = uint8((value - minWs) / (maxWs - minWs) * 255)
		}

		img.SetRGBA(x, y, color.RGBA{
			colorValue, colorValue, colorValue, 255,
		})
	})
}

func ReadDataFromNcFile(path string) (*gdaldraw.IrregularData, error) {
	group, err := netcdf.Open(path)
	if err != nil {
		return nil, fmt.Errorf("read nc file: %v", err)
	}
	defer group.Close()

	latVariable, err := group.GetVariable("wvc_lat")
	if err != nil {
		return nil, fmt.Errorf("get variable wvc_lat: %v", err)
	}
	latList := latVariable.Values.([][]float32)

	lonVariable, err := group.GetVariable("wvc_lon")
	if err != nil {
		return nil, fmt.Errorf("get variable wvc_lon: %v", err)
	}
	lonList := lonVariable.Values.([][]float32)

	wsVariable, err := group.GetVariable("wind_speed_selection")
	if err != nil {
		return nil, fmt.Errorf("get variable wind_speed_selection: %v", err)
	}
	wsList := wsVariable.Values.([][]int16)
	fillValue, has := wsVariable.Attributes.Get("fill_value")
	if !has {
		return nil, fmt.Errorf("get attribute fill_value error")
	}
	addOffset, has := wsVariable.Attributes.Get("add_offset")
	if !has {
		return nil, fmt.Errorf("get attribute add_offset error")
	}
	scaleFactor, has := wsVariable.Attributes.Get("scale_factor")
	if !has {
		return nil, fmt.Errorf("get attribute scale_factor error")
	}

	var latArr, lonArr, valueArr []float64
	for i := 0; i < 1624; i++ {
		for j := 0; j < 76; j++ {
			var lat = latList[i][j]
			var lon = lonList[i][j]
			var ws = wsList[i][j]

			if lat > 90 {
				continue
			}

			if lon > 360 {
				continue
			}

			if ws == fillValue {
				continue
			}

			latArr = append(latArr, float64(lat))
			lonArr = append(lonArr, float64(lon))
			valueArr = append(valueArr, float64(convertInt16ToFloat(ws, scaleFactor.(float32), addOffset.(float32))))
		}
	}

	return &gdaldraw.IrregularData{
		LatList:     latArr,
		LonList:     lonArr,
		Accuracy:    0.1,
		ValueList:   valueArr,
		OutFilePath: "./h2b.png",
	}, nil
}

func convertInt16ToFloat(real int16, scaleFactor, addOffset float32) float32 {
	return float32(real)*scaleFactor + addOffset
}
