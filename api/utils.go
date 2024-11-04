package handler

import (
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
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Theme:  "white",
			Width:  "100%",
			Height: "600px",
		}),
		charts.WithTitleOpts(opts.Title{
			Title: title,
			Left:  "center",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:    true,
			Trigger: "axis",
			AxisPointer: &opts.AxisPointer{
				Type: "line",
			},
		}),
		charts.WithLegendOpts(opts.Legend{
			Show:  true,
			Right: "10%",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "时间",
			AxisLabel: &opts.AxisLabel{
				Rotate: 45,
			},
			SplitLine: &opts.SplitLine{
				Show: true,
				LineStyle: &opts.LineStyle{
					Color: "#f0f0f0",
					Type:  "dashed",
				},
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "价格 (元/克)",
			AxisLabel: &opts.AxisLabel{
				Formatter: "{value} 元",
			},
			SplitLine: &opts.SplitLine{
				Show: true,
				LineStyle: &opts.LineStyle{
					Color: "#f0f0f0",
					Type:  "dashed",
				},
			},
		}),
	)

	line.SetXAxis(timeLabels).
		AddSeries("黄金价格", generateLineItems(values)).
		SetSeriesOptions(
			charts.WithLineChartOpts(opts.LineChart{
				Smooth:     true,
				ShowSymbol: true,
				Symbol:     "circle",
				SymbolSize: 8,
			}),
			charts.WithLabelOpts(opts.Label{
				Show:      true,
				Formatter: "{c} 元",
			}),
			charts.WithItemStyleOpts(opts.ItemStyle{
				Color: "#ff4d4f", // 红色线条
			}),
			charts.WithAreaStyleOpts(opts.AreaStyle{
				Opacity: 0.2,
				Color:   "#ff4d4f", // 阴影区域颜色
			}),
		)

	w.Header().Set("Content-Type", "text/html")
	line.Render(w)
}

func generateLineItems(values []float64) []opts.LineData {
	items := make([]opts.LineData, 0, len(values))
	for _, v := range values {
		items = append(items, opts.LineData{Value: v})
	}
	return items
}
