package handler

import (
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func TodayHandler(c *gin.Context) {
	timestamps := getTodayTimestamps()

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
			timeLabels = append(timeLabels, time.Unix(ts, 0).Format("15:04"))
			values = append(values, val)
		} else {
			timeLabels = append(timeLabels, time.Unix(ts, 0).Format("15:04"))
			values = append(values, 0)
		}
	}

	// 修改渲染逻辑，使用 gin 的响应方式
	c.Writer.Header().Set("Content-Type", "image/png")
	renderChart(c.Writer, timeLabels, values, "最近24小时黄金价格走势")
}

// 生成最近24小时的所有10分钟时间戳
func getTodayTimestamps() []int64 {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)

	// 从当前时间往前推24小时
	startTime := now.Add(-24 * time.Hour)
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
	// 计算从24小时前到现在的所有10分钟时间戳
	for i := 0; i < 24*6+1; i++ { // 24小时 * 每小时6个10分钟 + 当前时间点
		t := startTime.Add(time.Duration(i*10) * time.Minute)
		// 如果超过当前时间，就停止
		if t.After(now) {
			break
		}
		timestamps = append(timestamps, getNormalizedTimestamp(t))
	}
	return timestamps
}
