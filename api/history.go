package handler

import (
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func HistoryHandler(c *gin.Context) {
	timestamps := getLast7DaysTimestamps()

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

	loc, _ := time.LoadLocation("Asia/Shanghai")
	for _, ts := range timestamps {
		if val, exists := timeValues[ts]; exists {
			// 将时间戳转换为上海时间后再格式化
			t := time.Unix(ts, 0).In(loc)
			timeLabels = append(timeLabels, t.Format("01-02 15:04"))
			values = append(values, val)
		} else {
			t := time.Unix(ts, 0).In(loc)
			timeLabels = append(timeLabels, t.Format("01-02 15:04"))
			values = append(values, 0)
		}
	}

	// 对timeLabels和values做处理， 对比连续的数据， 如果前面的值和后面相等， 那就只记录前面的。 最后一个有值的数据也保留
	var filteredTimeLabels []string
	var filteredValues []float64

	for i := 0; i < len(values); i++ {
		if i == len(values)-1 || values[i] != values[i+1] {
			filteredTimeLabels = append(filteredTimeLabels, timeLabels[i])
			filteredValues = append(filteredValues, values[i])
		}
	}
	// 保留最后一个有值的数据
	if len(values) > 0 && (len(filteredValues) == 0 || filteredValues[len(filteredValues)-1] != values[len(values)-1]) {
		filteredTimeLabels = append(filteredTimeLabels, timeLabels[len(timeLabels)-1])
		filteredValues = append(filteredValues, values[len(values)-1])
	}

	timeLabels = filteredTimeLabels
	values = filteredValues

	// 修改渲染逻辑，使用 gin 的响应方式
	c.Writer.Header().Set("Content-Type", "image/png")
	renderChart(c.Writer, timeLabels, values, "近3天黄金价格走势")
}

// 获取最近3天的时间戳
func getLast3DaysTimestamps() []int64 {
	// 设置上海时区
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
		0, 0,
		loc,
	)

	var timestamps []int64
	// 计算从3天前到现在的所有10分钟时间戳
	for i := 0; i < 3*24*6+1; i++ {
		t := startTime.Add(time.Duration(i*10) * time.Minute)
		if t.After(now) {
			break
		}
		timestamps = append(timestamps, getNormalizedTimestamp(t))
	}
	return timestamps
}


// 获取最近7天的时间戳
func getLast7DaysTimestamps() []int64 {
	// 设置上海时区
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)

	// 从当前时间往前推7天
	startTime := now.AddDate(0, 0, -7)
	// 归一化开始时间到10分钟
	startTime = time.Date(
		startTime.Year(),
		startTime.Month(),
		startTime.Day(),
		startTime.Hour(),
		(startTime.Minute()/10)*10,
		0, 0,
		loc,
	)

	var timestamps []int64
	// 计算从7天前到现在的所有10分钟时间戳
	for i := 0; i < 7*24*6+1; i++ {
		t := startTime.Add(time.Duration(i*10) * time.Minute)
		if t.After(now) {
			break
		}
		timestamps = append(timestamps, getNormalizedTimestamp(t))
	}
	return timestamps
}
