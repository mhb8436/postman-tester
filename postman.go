package main

import (
	"time"
)

// Postman Collection 구조체 정의
type Collection struct {
	Info CollectionInfo `json:"info"`
	Item []Item         `json:"item"`
}

type CollectionInfo struct {
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

type Item struct {
	Name    string            `json:"name"`
	Item    []Item            `json:"item,omitempty"` // 중첩된 폴더 구조
	Request *Request          `json:"request,omitempty"`
	Event   []Event           `json:"event,omitempty"`
}

type Request struct {
	Method string      `json:"method"`
	Header []Header    `json:"header"`
	Body   *Body       `json:"body,omitempty"`
	URL    interface{} `json:"url"` // string 또는 URL 객체
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Body struct {
	Mode    string      `json:"mode"`
	Raw     string      `json:"raw,omitempty"`
	Options *BodyOptions `json:"options,omitempty"`
}

type BodyOptions struct {
	Raw *RawOptions `json:"raw,omitempty"`
}

type RawOptions struct {
	Language string `json:"language"`
}

type URL struct {
	Raw      string   `json:"raw"`
	Protocol string   `json:"protocol"`
	Host     []string `json:"host"`
	Port     string   `json:"port,omitempty"`
	Path     []string `json:"path"`
	Query    []Query  `json:"query,omitempty"`
}

type Query struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Event struct {
	Listen string       `json:"listen"`
	Script EventScript  `json:"script"`
}

type EventScript struct {
	Type string   `json:"type"`
	Exec []string `json:"exec"`
}

// 테스트 결과 구조체
type TestResult struct {
	Name           string        `json:"name"`
	Method         string        `json:"method"`
	URL            string        `json:"url"`
	StatusCode     int           `json:"status_code"`
	ResponseTime   time.Duration `json:"response_time"`
	Success        bool          `json:"success"`
	ErrorMessage   string        `json:"error_message,omitempty"`
	ResponseBody   string        `json:"response_body,omitempty"`
	RequestHeaders map[string]string `json:"request_headers"`
	Timestamp      time.Time     `json:"timestamp"`
}

type TestSummary struct {
	CollectionName string        `json:"collection_name"`
	FilePath       string        `json:"file_path"`
	TotalTests     int           `json:"total_tests"`
	PassedTests    int           `json:"passed_tests"`
	FailedTests    int           `json:"failed_tests"`
	TotalTime      time.Duration `json:"total_time"`
	Results        []TestResult  `json:"results"`
	StartTime      time.Time     `json:"start_time"`
	EndTime        time.Time     `json:"end_time"`
}