package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// 启动 HTTP 服务器，监听端口
	http.HandleFunc("/send-request", Handler)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}

	// 创建 HTTP POST 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求并获取响应
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// 打印响应状态码和响应体
	fmt.Println("Response status:", resp.Status)
	fmt.Println("Response body:", string(body))

	// 将响应体返回给客户端
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
