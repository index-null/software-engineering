package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Input struct {
	Prompt string `json:"prompt"`
}

type Parameters struct {
	//Style string `json:"style"`
	Size string `json:"size"`
	//N    int    `json:"n"`// 默认为4
}

type RequestBody struct {
	Model      string     `json:"model"`
	Input      Input      `json:"input"`
	Parameters Parameters `json:"parameters"`
}

type TaskResponse struct {
	Output struct {
		TaskStatus string `json:"task_status"`
		TaskID     string `json:"task_id"`
	} `json:"output"`
	RequestID string `json:"request_id"`
}

type TaskStatusResponse struct {
	RequestID string `json:"request_id"`
	Output    struct {
		TaskStatus string `json:"task_status"`
		TaskID     string `json:"task_id"`
		Code       string `json:"code"`
		Message    string `json:"message"`
		Result     []struct {
			URL string `json:"url"`
		} `json:"results"`
	} `json:"output"`
}

func main() {
	// 创建 HTTP 客户端
	client := &http.Client{}

	// 构建请求体
	requestBody := RequestBody{
		Model: "wanx-v1",
		Input: Input{
			Prompt: "奔跑小猫",
		},
		Parameters: Parameters{
			//Style: "<auto>",
			Size: "1024*1024",
			//N:    1,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal(err)
	}

	// 创建 POST 请求来创建任务
	req, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/api/v1/services/aigc/text2image/image-synthesis", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}

	// 设置请求头
	apiKey := "sk-6e79f5171c934d8fbbbdb0f4cd42d669"
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DashScope-Async", "enable")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("请求失败%v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("读取响应失败%v", err)
	}

	// 解析响应以获取任务ID
	var taskResponse TaskResponse
	err = json.Unmarshal(bodyText, &taskResponse)
	if err != nil {
		log.Fatalf("解析响应失败%v", err)
	}

	taskID := taskResponse.Output.TaskID
	if taskID == "" {
		fmt.Printf("任务ID为空，请检查请求是否成功%v", taskResponse)
		return
	}
	fmt.Printf("任务已创建，任务ID: %s\n", taskID)
	fmt.Printf(": %v\n", taskResponse)

	// 轮询任务状态
	for {
		time.Sleep(40 * time.Second) // 每分钟轮询一次

		// 创建 GET 请求来查询任务状态
		statusReq, err := http.NewRequest("GET", fmt.Sprintf("https://dashscope.aliyuncs.com/api/v1/tasks/%s", taskID), nil)
		if err != nil {
			log.Fatal(err)
		}

		// 设置请求头
		statusReq.Header.Set("Authorization", "Bearer "+apiKey)

		// 发送请求
		statusResp, err := client.Do(statusReq)
		if err != nil {
			log.Fatalf("请求失败%v", err)
		}
		defer statusResp.Body.Close()

		// 读取响应体
		statusBodyText, err := io.ReadAll(statusResp.Body)
		if err != nil {
			log.Fatalf("读取响应失败%v", err)
		}

		// 解析响应以获取任务状态
		var taskStatusResponse TaskStatusResponse
		err = json.Unmarshal(statusBodyText, &taskStatusResponse)
		if err != nil {
			log.Fatalf("解析响应失败%v", err)
		}

		fmt.Printf("任务状态: %s\n", taskStatusResponse.Output.TaskStatus)

		if taskStatusResponse.Output.TaskStatus == "SUCCEEDED" {
			fmt.Printf("任务已完成，结果: %s\n", taskStatusResponse.Output.Result)
			break
		} else if taskStatusResponse.Output.TaskStatus == "FAILED" {
			fmt.Printf("任务失败%v,message:%v", taskStatusResponse.Output.Code, taskStatusResponse.Output.Message)
			break
		}
	}
}
