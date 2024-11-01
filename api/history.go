package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func HistoryHandler(w http.ResponseWriter, r *http.Request) {
	timestamps := getLast3DaysTimestamps()

	// 准备所有键
	keys := make([]string, len(timestamps))
	for i, ts := range timestamps {
		keys[i] = "price:" + strconv.FormatInt(ts, 10)
	}

	// 批量获取数据
	pipe := redisClient.Pipeline()
	cmds := make([]*redis.StringCmd, len(keys))
	for i, key := range keys {
		cmds[i] = pipe.Get(context.Background(), key)
	}
	pipe.Exec(context.Background())

	// 处理结果
	timeValues := make(map[int64]float64)
	var lastValue float64

	for i, cmd := range cmds {
		priceStr, err := cmd.Result()
		if err == nil {
			price, _ := strconv.ParseFloat(priceStr, 64)
			timeValues[timestamps[i]] = price
			lastValue = price
		} else if lastValue != 0 {
			timeValues[timestamps[i]] = lastValue
		}
	}

	// 准备图表数据
	var timeLabels []string
	var values []float64

	for _, ts := range timestamps {
		if val, exists := timeValues[ts]; exists {
			timeLabels = append(timeLabels, time.Unix(ts, 0).Format("01-02 15:04"))
			values = append(values, val)
		} else {
			timeLabels = append(timeLabels, time.Unix(ts, 0).Format("01-02 15:04"))
			values = append(values, 0)
		}
	}

	renderChart(w, timeLabels, values, "近3天黄金价格走势")
}

// 获取最近3天的时间戳
func getLast3DaysTimestamps() []int64 {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)

	// 从当前时间往前推3天
	startTime := now.AddDate(0, 0, -3)
	// 归一化开始时间到10分钟
	startTime = time.Date(
		startTime.Year(),
		startTime.Month(),
		startTime.Day(),
		startTime.Hour(),
		(startTime.Minute()/10)*10,
		0, 0, loc,
	)

	var timestamps []int64
	// 计算从3天前到现在的所有10分钟时间戳
	for i := 0; i < 3*24*6+1; i++ {
		t := startTime.Add(time.Duration(i*10) * time.Minute)
		// 如果超过当前时间，就停止
		if t.After(now) {
			break
		}
		timestamps = append(timestamps, getNormalizedTimestamp(t))
	}
	return timestamps
}
