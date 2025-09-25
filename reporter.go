package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"
)


type Reporter struct {
	format string
}

func NewReporter(format string) *Reporter {
	return &Reporter{format: format}
}

func (r *Reporter) Print(summaries []*TestSummary) {
	switch r.format {
	case "json":
		r.printJSON(summaries)
	case "html":
		r.printHTML(summaries)
	case "csv":
		r.printCSV(summaries)
	default:
		r.printText(summaries)
	}
}

func (r *Reporter) SaveToFile(summaries []*TestSummary, filename string) error {
	var content string
	var err error

	switch r.format {
	case "json":
		content, err = r.generateJSON(summaries)
	case "html":
		content, err = r.generateHTML(summaries)
	case "csv":
		content, err = r.generateCSV(summaries)
	default:
		content, err = r.generateText(summaries)
	}

	if err != nil {
		return err
	}

	return os.WriteFile(filename, []byte(content), 0644)
}

func (r *Reporter) printText(summaries []*TestSummary) {
	content, _ := r.generateText(summaries)
	fmt.Print(content)
}

func (r *Reporter) generateText(summaries []*TestSummary) (string, error) {
	var sb strings.Builder

	sb.WriteString("📊 상세 테스트 결과\n")
	sb.WriteString("=" + strings.Repeat("=", 50) + "\n\n")

	for i, summary := range summaries {
		sb.WriteString(fmt.Sprintf("[%d] %s\n", i+1, summary.CollectionName))
		sb.WriteString(fmt.Sprintf("파일: %s\n", summary.FilePath))
		sb.WriteString(fmt.Sprintf("실행시간: %.2fs\n", summary.TotalTime.Seconds()))
		sb.WriteString(fmt.Sprintf("결과: %d개 성공, %d개 실패\n", summary.PassedTests, summary.FailedTests))
		sb.WriteString("-" + strings.Repeat("-", 30) + "\n")

		for j, result := range summary.Results {
			status := "✅"
			if !result.Success {
				status = "❌"
			}

			sb.WriteString(fmt.Sprintf("  [%d.%d] %s %s\n", i+1, j+1, status, result.Name))
			sb.WriteString(fmt.Sprintf("        %s %s\n", result.Method, result.URL))
			sb.WriteString(fmt.Sprintf("        응답: HTTP %d (%.2fs)\n", result.StatusCode, result.ResponseTime.Seconds()))

			if !result.Success {
				sb.WriteString(fmt.Sprintf("        오류: %s\n", result.ErrorMessage))
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func (r *Reporter) printCSV(summaries []*TestSummary) {
	content, _ := r.generateCSV(summaries)
	fmt.Print(content)
}

func (r *Reporter) generateCSV(summaries []*TestSummary) (string, error) {
	var sb strings.Builder
	
	// UTF-8 BOM 추가 (Excel에서 한글 제대로 표시하기 위함)
	sb.WriteString("\uFEFF")
	
	// CSV 헤더
	sb.WriteString("Collection,FilePath,TestName,Method,URL,StatusCode,Success,ResponseTime,ErrorMessage\n")
	
	// 각 컬렉션의 테스트 결과를 CSV 행으로 변환
	for _, summary := range summaries {
		for _, result := range summary.Results {
			// CSV 필드 값들을 이스케이프 처리
			collection := escapeCSV(summary.CollectionName)
			filePath := escapeCSV(summary.FilePath)
			testName := escapeCSV(result.Name)
			method := escapeCSV(result.Method)
			url := escapeCSV(result.URL)
			statusCode := fmt.Sprintf("%d", result.StatusCode)
			success := "true"
			if !result.Success {
				success = "false"
			}
			responseTime := fmt.Sprintf("%.3f", result.ResponseTime.Seconds())
			errorMessage := escapeCSV(result.ErrorMessage)
			
			sb.WriteString(fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s\n",
				collection, filePath, testName, method, url, statusCode, success, responseTime, errorMessage))
		}
	}
	
	return sb.String(), nil
}

// CSV 필드 값에 대한 이스케이프 처리
func escapeCSV(value string) string {
	// 빈 값 처리
	if value == "" {
		return ""
	}
	
	// 쌍따옴표, 쉼표, 줄바꿈이 포함된 경우 쌍따옴표로 감싸고 내부 쌍따옴표는 두 번 반복
	if strings.Contains(value, "\"") || strings.Contains(value, ",") || strings.Contains(value, "\n") {
		value = strings.ReplaceAll(value, "\"", "\"\"")
		return "\"" + value + "\""
	}
	
	return value
}

func (r *Reporter) printJSON(summaries []*TestSummary) {
	content, _ := r.generateJSON(summaries)
	fmt.Print(content)
}

func (r *Reporter) generateJSON(summaries []*TestSummary) (string, error) {
	data, err := json.MarshalIndent(summaries, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (r *Reporter) printHTML(summaries []*TestSummary) {
	content, _ := r.generateHTML(summaries)
	fmt.Print(content)
}

func (r *Reporter) generateHTML(summaries []*TestSummary) (string, error) {
	tmpl := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Postman Collection Test Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background: #f5f5f5; padding: 20px; border-radius: 5px; margin-bottom: 20px; }
        .collection { border: 1px solid #ddd; margin-bottom: 20px; border-radius: 5px; }
        .collection-header { background: #f8f9fa; padding: 15px; border-bottom: 1px solid #ddd; }
        .collection-name { font-size: 1.2em; font-weight: bold; margin: 0; }
        .collection-stats { color: #666; margin-top: 5px; }
        .test-item { padding: 10px 15px; border-bottom: 1px solid #eee; }
        .test-item:last-child { border-bottom: none; }
        .test-success { border-left: 4px solid #28a745; }
        .test-failed { border-left: 4px solid #dc3545; }
        .test-name { font-weight: bold; }
        .test-details { color: #666; margin-top: 5px; }
        .error-message { color: #dc3545; font-style: italic; margin-top: 5px; }
        .summary { background: #e9ecef; padding: 15px; border-radius: 5px; }
        .success-rate { font-size: 1.1em; font-weight: bold; }
    </style>
</head>
<body>
    <div class="header">
        <h1>🧪 Postman Collection Test Report</h1>
        <p>생성 시간: {{.Timestamp}}</p>
    </div>

    {{range $index, $summary := .Summaries}}
    <div class="collection">
        <div class="collection-header">
            <h2 class="collection-name">{{$summary.CollectionName}}</h2>
            <div class="collection-stats">
                파일: {{$summary.FilePath}}<br>
                실행시간: {{printf "%.2f" $summary.TotalTime.Seconds}}초 | 
                총 {{$summary.TotalTests}}개 테스트 | 
                성공: {{$summary.PassedTests}}개 | 
                실패: {{$summary.FailedTests}}개
            </div>
        </div>
        
        {{range $summary.Results}}
        <div class="test-item {{if .Success}}test-success{{else}}test-failed{{end}}">
            <div class="test-name">
                {{if .Success}}✅{{else}}❌{{end}} {{.Name}}
            </div>
            <div class="test-details">
                {{.Method}} {{.URL}}<br>
                응답: HTTP {{.StatusCode}} ({{printf "%.2f" .ResponseTime.Seconds}}초)
            </div>
            {{if not .Success}}
            <div class="error-message">오류: {{.ErrorMessage}}</div>
            {{end}}
        </div>
        {{end}}
    </div>
    {{end}}

    <div class="summary">
        <h3>📋 전체 요약</h3>
        <p>총 {{.TotalCollections}}개 컬렉션, {{.TotalTests}}개 테스트</p>
        <p class="success-rate">
            성공률: {{printf "%.1f" .SuccessRate}}% 
            ({{.TotalPassed}}개 성공 / {{.TotalFailed}}개 실패)
        </p>
    </div>
</body>
</html>`

	// 통계 계산
	data := struct {
		Summaries        []*TestSummary
		Timestamp        string
		TotalCollections int
		TotalTests       int
		TotalPassed      int
		TotalFailed      int
		SuccessRate      float64
	}{
		Summaries:        summaries,
		Timestamp:        time.Now().Format("2006-01-02 15:04:05"),
		TotalCollections: len(summaries),
	}

	for _, summary := range summaries {
		data.TotalTests += summary.TotalTests
		data.TotalPassed += summary.PassedTests
		data.TotalFailed += summary.FailedTests
	}

	if data.TotalTests > 0 {
		data.SuccessRate = float64(data.TotalPassed) / float64(data.TotalTests) * 100
	}

	t, err := template.New("report").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	err = t.Execute(&sb, data)
	if err != nil {
		return "", err
	}

	return sb.String(), nil
}