package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	// 启动 HTTP 服务器，监听端口
	http.HandleFunc("/send-request", Handler)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Response 结构定义
type Response struct {
	JsonResult  string      `json:"JsonResult"`
	JsonMessage JsonMessage `json:"JsonMessage"`
	JsonData    []JsonData  `json:"JsonData"`
}

type JsonMessage struct {
	MessageIndex string `json:"MessageIndex"`
	Remark       string `json:"Remark"`
	MessageInfo  string `json:"MessageInfo"`
}

type JsonData struct {
	SQLBuilderID string    `json:"SQLBuilderID"`
	Field        []Field   `json:"FIELD"`
	Row          []RowData `json:"ROW"`
}

type Field struct {
	AttrName  string `json:"attrname"`
	FieldType string `json:"fieldtype"`
	Width     string `json:"WIDTH"`
}

type RowData struct {
	FKindName  string `json:"FKIND_NAME"`
	FPriceBase string `json:"FPRICE_BASE"`
	FNewTime   string `json:"FNEWTIME"`
	FTopRemark string `json:"FTOP_REMARK"`
	FRemark    string `json:"FREMARK"`
}

// Handler 处理 HTTP 请求
func Handler(w http.ResponseWriter, r *http.Request) {
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
		respondWithError(w)
		return
	}

	// 创建 HTTP POST 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		respondWithError(w)
		return
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求并获取响应
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		respondWithError(w)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		respondWithError(w)
		return
	}

	// 解析响应 JSON
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		respondWithError(w)
		return
	}

	// 检查是否有有效的 Row 数据
	if len(response.JsonData) > 0 && len(response.JsonData[0].Row) > 0 {
		// 获取 RowData 中的最后一个元素
		lastRow := response.JsonData[0].Row[len(response.JsonData[0].Row)-1]

		// 返回最后一个 RowData 作为 JSON 响应
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(lastRow)
	} else {
		respondWithError(w)
	}
}

// respondWithError 返回默认的 RowData 错误响应
func respondWithError(w http.ResponseWriter) {
	// 获取当前时间并格式化
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// 构建默认的 RowData 响应
	defaultRow := RowData{
		FKindName:  "Gold Price",
		FPriceBase: "0 元/克",
		FNewTime:   currentTime,
		FTopRemark: "No data available",
		FRemark:    "Generated due to error",
	}

	// 设置响应头并返回默认 RowData
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(defaultRow)
}
