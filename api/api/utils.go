package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// 获取指定时间的10分钟归一化时间戳
func getNormalizedTimestamp(t time.Time) int64 {
	// 确保时间是上海时区
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t = t.In(loc)

	minutes := t.Minute()
	normalizedMinutes := (minutes / 10) * 10
	normalizedTime := time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		normalizedMinutes,
		0, 0,
		loc,
	)
	return normalizedTime.Unix()
}

func renderChart(w http.ResponseWriter, timeLabels []string, values []float64, title string) {
	// 计算最新价格和Y轴范围
	var latestPrice float64
	minPrice, maxPrice := values[0], values[0]
	for _, v := range values {
		if v != 0 {
			if v < minPrice || minPrice == 0 {
				minPrice = v
			}
			if v > maxPrice {
				maxPrice = v
			}
			latestPrice = v
		}
	}

	// 计算Y轴范围，留出5%的边距
	yAxisMin := minPrice - (maxPrice-minPrice)*0.05
	yAxisMax := maxPrice + (maxPrice-minPrice)*0.05

	// 更新标题，加入最新价格并换行
	titleWithPrice := fmt.Sprintf("%s\n(最新价格: %.2f元/克)", title, latestPrice)

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme:  "white",
			Width:  "100%",
			Height: "600px",
		}),
		charts.WithTitleOpts(opts.Title{
			Title: titleWithPrice,
			Left:  "center",
			Top:   "20px",
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      0,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:    true,
			Trigger: "axis",
			AxisPointer: &opts.AxisPointer{
				Type: "line",
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "价格 (元/克)",
			Min:  yAxisMin,
			Max:  yAxisMax,
			AxisLabel: &opts.AxisLabel{
				Show:      true,
				Formatter: "{value}",
			},
			SplitLine: &opts.SplitLine{
				Show: true,
				LineStyle: &opts.LineStyle{
					Color: "#f0f0f0",
					Type:  "dashed",
				},
			},
			SplitNumber: 10,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "时间",
			Show: true,
			AxisLabel: &opts.AxisLabel{
				Show:   true,
				Rotate: 45,
				Margin: 15,
			},
			SplitLine: &opts.SplitLine{
				Show: true,
				LineStyle: &opts.LineStyle{
					Color: "#f0f0f0",
					Type:  "dashed",
					Width: 1,
				},
			},
			AxisTick: &opts.AxisTick{
				Show: true,
			},
		}),
	)

	line.SetXAxis(timeLabels).
		AddSeries("黄金价格", generateLineItems(values)).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{
				Smooth:     false,
				ShowSymbol: true,
				Symbol:     "circle",
				SymbolSize: 6,
			}),
			charts.WithLabelOpts(opts.Label{
				Show: false,
			}),
			charts.WithItemStyleOpts(opts.ItemStyle{
				Color: "#ff4d4f",
			}),
			charts.WithAreaStyleOpts(opts.AreaStyle{
				Opacity: 0.1,
				Color:   "#ff4d4f",
			}),
			charts.WithMarkPointNameTypeItemOpts(
				opts.MarkPointNameTypeItem{
					Name:      "最高价",
					Type:      "max",
					ItemStyle: &opts.ItemStyle{Color: "#ff4d4f"},
				},
				opts.MarkPointNameTypeItem{
					Name:      "最低价",
					Type:      "min",
					ItemStyle: &opts.ItemStyle{Color: "#52c41a"},
				},
			),
		)

	w.Header().Set("Content-Type", "text/html")
	line.Render(w)
}

// 生成K线数据
func generateKLineData(timeLabels []string, values []float64) []opts.KlineData {
	items := make([]opts.KlineData, 0)
	windowSize := 6 // 每小时的数据点数（10分钟一个点）

	for i := 0; i < len(values); i += windowSize {
		end := i + windowSize
		if end > len(values) {
			end = len(values)
		}

		var open, close, high, low float64
		validData := false

		// 获取这个时间窗口的数据
		for j := i; j < end; j++ {
			if values[j] != 0 {
				if !validData {
					open = values[j]
					high = values[j]
					low = values[j]
					validData = true
				}
				if values[j] > high {
					high = values[j]
				}
				if values[j] < low {
					low = values[j]
				}
				close = values[j]
			}
		}

		if validData {
			items = append(items, opts.KlineData{
				Value: [4]float64{open, close, low, high},
			})
		}
	}
	return items
}

// 添加 generateLineItems 函数
func generateLineItems(values []float64) []opts.LineData {
	items := make([]opts.LineData, len(values))
	for i := 0; i < len(values); i++ {
		items[i] = opts.LineData{Value: values[i]}
	}
	return items
}
