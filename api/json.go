package handler

//package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GetNormalizedTimeStamp 返回归一化为10分钟的北京时间的 Unix 时间戳
func GetNormalizedTimeStamp() int64 {
	// 设置上海时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}

	// 获取当前上海时间
	now := time.Now().In(loc)

	// 计算归一化的分钟
	normalizedMinutes := (now.Minute() / 10) * 10

	// 创建新的时间，设置归一化的分钟
	normalizedTime := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		normalizedMinutes,
		0, 0,
		loc,
	)

	return normalizedTime.Unix()
}

func callCaibai() (*RowData, error) {
	// 定义目标 URL
	url := "http://111.198.86.222/BAP/OpenApi"

	// 构建 JSON 数据
	data := map[string]interface{}{
		"Context": map[string]string{
			"token":     "",
			"version":   "",
			"from":      "2",
			"mchid":     "",
			"appid":     "",
			"timestamp": "",
		},
		"SQLBuilderItem": []map[string]interface{}{
			{
				"SQLBuilderID": "{005A5001-B9AD-41CB-8409-8F7675D19143}",
				"TableName":    "BS_POS_GP_MA",
				"Caption":      "每日金价",
				"Select": map[string]string{
					"FMID":          "{4F054C98-16B8-8A9E-3112-F8AFC1BC77E9}",
					"FPID":          "{4F054C98-16B8-8A9E-3112-F8AFC1BC77E9}",
					"FTID":          "",
					"FUID":          "",
					"FOID":          "{7D77D027-9824-4156-A25E-12FC59527DDE}",
					"FWID":          "",
					"FORG_STORE_ID": "",
				},
			},
		},
	}

	// 将数据编码为 JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// 创建 HTTP POST 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求并获取响应
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 解析响应 JSON
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	// 检查是否有有效的 Row 数据
	if len(response.JsonData) > 0 && len(response.JsonData[0].Row) > 0 {
		// 获取 RowData 中的最后一个元素
		lastRow := response.JsonData[0].Row[len(response.JsonData[0].Row)-1]
		return &lastRow, nil
	} else {
		return nil, errors.New("no data available")
	}
}

func insertToRedis(rowData *RowData) error {
	key := fmt.Sprintf("price:%d", GetNormalizedTimeStamp())
	if price, err := strconv.ParseFloat(strings.ReplaceAll(rowData.FPriceBase, " 元/克", ""), 64); err == nil {
		log.Println(redisClient.Set(context.Background(), key, price, 0).Result())
	}
	return nil
}

// Handler 处理 HTTP 请求
func JsonHandler(c *gin.Context) {
	if lastRow, err := callCaibai(); err == nil {
		insertToRedis(lastRow)
		c.JSON(http.StatusOK, lastRow)
	} else {
		// 获取上海时间
		loc, _ := time.LoadLocation("Asia/Shanghai")
		currentTime := time.Now().In(loc).Format("2006-01-02 15:04:05")

		defaultRow := RowData{
			FKindName:  "Gold Price",
			FPriceBase: "0 元/克",
			FNewTime:   currentTime,
			FTopRemark: "No data available",
			FRemark:    "Generated due to error",
		}
		c.JSON(http.StatusInternalServerError, defaultRow)
	}
}
