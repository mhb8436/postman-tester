package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Runner struct {
	client *http.Client
}

func NewRunner() *Runner {
	return &Runner{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Postman 컬렉션 파일을 로드하고 파싱
func (r *Runner) LoadCollection(filepath string) (*Collection, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("파일을 읽을 수 없습니다: %v", err)
	}

	var collection Collection
	if err := json.Unmarshal(data, &collection); err != nil {
		return nil, fmt.Errorf("JSON 파싱 실패: %v", err)
	}

	return &collection, nil
}

// 컬렉션의 모든 요청 실행
func (r *Runner) RunCollection(collection *Collection) *TestSummary {
	summary := &TestSummary{
		StartTime: time.Now(),
		Results:   make([]TestResult, 0),
	}

	// 모든 아이템을 재귀적으로 실행
	r.executeItems(collection.Item, summary)

	summary.EndTime = time.Now()
	summary.TotalTime = summary.EndTime.Sub(summary.StartTime)
	summary.TotalTests = len(summary.Results)
	
	for _, result := range summary.Results {
		if result.Success {
			summary.PassedTests++
		} else {
			summary.FailedTests++
		}
	}

	return summary
}

// 아이템들을 재귀적으로 실행 (폴더 구조 지원)
func (r *Runner) executeItems(items []Item, summary *TestSummary) {
	for _, item := range items {
		if item.Request != nil {
			// 요청이 있는 아이템 실행
			result := r.executeRequest(item)
			summary.Results = append(summary.Results, result)
		} else if len(item.Item) > 0 {
			// 중첩된 아이템들 재귀 실행
			r.executeItems(item.Item, summary)
		}
	}
}

// 개별 요청 실행
func (r *Runner) executeRequest(item Item) TestResult {
	result := TestResult{
		Name:           item.Name,
		Timestamp:      time.Now(),
		RequestHeaders: make(map[string]string),
	}

	startTime := time.Now()

	// URL 파싱
	url := r.parseURL(item.Request.URL)
	result.URL = url
	result.Method = item.Request.Method

	// HTTP 요청 생성
	var body io.Reader
	if item.Request.Body != nil && item.Request.Body.Raw != "" {
		body = strings.NewReader(item.Request.Body.Raw)
	}

	req, err := http.NewRequest(item.Request.Method, url, body)
	if err != nil {
		result.Success = false
		result.ErrorMessage = fmt.Sprintf("요청 생성 실패: %v", err)
		return result
	}

	// 헤더 설정
	for _, header := range item.Request.Header {
		if header.Key != "" && header.Value != "" {
			req.Header.Set(header.Key, header.Value)
			result.RequestHeaders[header.Key] = header.Value
		}
	}

	// Content-Type 설정 (JSON body가 있는 경우)
	if item.Request.Body != nil && item.Request.Body.Mode == "raw" {
		if item.Request.Body.Options != nil && item.Request.Body.Options.Raw != nil {
			if item.Request.Body.Options.Raw.Language == "json" {
				req.Header.Set("Content-Type", "application/json")
			}
		}
	}

	// 요청 실행
	resp, err := r.client.Do(req)
	if err != nil {
		result.Success = false
		result.ErrorMessage = fmt.Sprintf("요청 실행 실패: %v", err)
		result.ResponseTime = time.Since(startTime)
		return result
	}
	defer resp.Body.Close()

	result.ResponseTime = time.Since(startTime)
	result.StatusCode = resp.StatusCode

	// 응답 본문 읽기
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Success = false
		result.ErrorMessage = fmt.Sprintf("응답 읽기 실패: %v", err)
		return result
	}
	result.ResponseBody = string(bodyBytes)

	// 성공 여부 판단 (2xx 상태코드)
	result.Success = resp.StatusCode >= 200 && resp.StatusCode < 300

	if !result.Success {
		result.ErrorMessage = fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	return result
}

// URL 파싱 (string 또는 URL 객체 모두 지원)
func (r *Runner) parseURL(urlInterface interface{}) string {
	switch v := urlInterface.(type) {
	case string:
		return v
	case map[string]interface{}:
		// URL 객체인 경우
		if raw, ok := v["raw"].(string); ok {
			return raw
		}
		// 수동으로 URL 구성
		return r.buildURL(v)
	default:
		return ""
	}
}

// URL 객체로부터 URL 문자열 구성
func (r *Runner) buildURL(urlObj map[string]interface{}) string {
	var buf bytes.Buffer

	// 프로토콜
	if protocol, ok := urlObj["protocol"].(string); ok {
		buf.WriteString(protocol)
		buf.WriteString("://")
	}

	// 호스트
	if hostInterface, ok := urlObj["host"]; ok {
		if hostArray, ok := hostInterface.([]interface{}); ok {
			hostParts := make([]string, len(hostArray))
			for i, part := range hostArray {
				if str, ok := part.(string); ok {
					hostParts[i] = str
				}
			}
			buf.WriteString(strings.Join(hostParts, "."))
		}
	}

	// 포트
	if port, ok := urlObj["port"].(string); ok && port != "" {
		buf.WriteString(":")
		buf.WriteString(port)
	}

	// 경로
	if pathInterface, ok := urlObj["path"]; ok {
		if pathArray, ok := pathInterface.([]interface{}); ok {
			for _, part := range pathArray {
				if str, ok := part.(string); ok {
					buf.WriteString("/")
					buf.WriteString(str)
				}
			}
		}
	}

	return buf.String()
}