package main

import (
	"fmt"
	"image"
	"image/color"

	gdaldraw "gitee.com/wuxi_jiufang/gdal_draw"
	"github.com/batchatco/go-native-netcdf/netcdf"
	"github.com/sirupsen/logrus"
)

const NcPath = "./nc_files/bye.nc"

func main() {
	data, err := ReadDataFromNcFile(NcPath)
	if err != nil {
		logrus.Errorf("ReadDataFromNcFile: %v", err)
	}

	minWs := float64(-1)
	maxWs := float64(32)

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

	latVariable, err := group.GetVariable("lat_rho")
	if err != nil {
		return nil, fmt.Errorf("get variable lat_rho: %v", err)
	}
	latList := latVariable.Values.([][]float64)

	lonVariable, err := group.GetVariable("lon_rho")
	if err != nil {
		return nil, fmt.Errorf("get variable lon_rho: %v", err)
	}
	lonList := lonVariable.Values.([][]float64)

	tempVariable, err := group.GetVariable("temp")
	if err != nil {
		return nil, fmt.Errorf("get variable temp: %v", err)
	}
	tempList := tempVariable.Values.([][][][]float32)[0][0]

	var latArr, lonArr, valueArr []float64
	for i := 0; i < 502; i++ {
		for j := 0; j < 349; j++ {
			var lat = latList[i][j]
			var lon = lonList[i][j]
			var temp = tempList[i][j]

			if lat > 90 {
				continue
			}

			if lon > 360 {
				continue
			}

			if temp > 32 {
				continue
			}

			latArr = append(latArr, float64(lat))
			lonArr = append(lonArr, float64(lon))
			valueArr = append(valueArr, float64(temp))
		}
	}

	return &gdaldraw.IrregularData{
		LatList:     latArr,
		LonList:     lonArr,
		Accuracy:    0.1,
		ValueList:   valueArr,
		OutFilePath: "./bye_3857.png",
	}, nil
}
