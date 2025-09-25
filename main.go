package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	directory = flag.String("dir", "./postman", "Postman 컬렉션 파일들이 있는 디렉토리")
	file      = flag.String("file", "", "단일 Postman 컬렉션 파일 (이 옵션 사용시 -dir 무시)")
	output    = flag.String("output", "", "결과를 저장할 파일 (선택사항, 기본값: 콘솔 출력)")
	format    = flag.String("format", "text", "출력 형식 (text, json, html, csv)")
	parallel  = flag.Int("parallel", 1, "병렬 실행할 컬렉션 수 (기본값: 1)")
	timeout   = flag.Int("timeout", 30, "요청 타임아웃 (초, 기본값: 30)")
	verbose   = flag.Bool("verbose", false, "상세 출력")
	help      = flag.Bool("help", false, "도움말 표시")
)

func main() {
	flag.Parse()

	if *help {
		printUsage()
		return
	}

	var files []string

	// 단일 파일이 지정된 경우
	if *file != "" {
		if _, err := os.Stat(*file); os.IsNotExist(err) {
			log.Fatalf("파일을 찾을 수 없습니다: %s", *file)
		}
		files = []string{*file}
	} else {
		// 디렉토리 확인
		if _, err := os.Stat(*directory); os.IsNotExist(err) {
			log.Fatalf("디렉토리를 찾을 수 없습니다: %s", *directory)
		}

		// Postman 컬렉션 파일 찾기
		var err error
		files, err = findCollectionFiles(*directory)
		if err != nil {
			log.Fatalf("컬렉션 파일 검색 실패: %v", err)
		}

		if len(files) == 0 {
			log.Fatalf("디렉토리에서 Postman 컬렉션 파일을 찾을 수 없습니다: %s", *directory)
		}
	}

	fmt.Printf("🚀 %d개의 Postman 컬렉션을 테스트합니다...\n\n", len(files))

	// 모든 컬렉션 실행 (병렬 처리 지원)
	allResults := make([]*TestSummary, 0, len(files))
	
	if *parallel <= 1 {
		// 순차 실행
		runner := NewRunner()
		for i, file := range files {
			result := runSingleCollection(runner, file, i+1, len(files), *verbose)
			if result != nil {
				allResults = append(allResults, result)
			}
		}
	} else {
		// 병렬 실행
		allResults = runCollectionsInParallel(files, *parallel, *verbose)
	}

	// 최종 결과 출력
	reporter := NewReporter(*format)
	if *output != "" {
		err := reporter.SaveToFile(allResults, *output)
		if err != nil {
			log.Fatalf("결과 저장 실패: %v", err)
		}
		fmt.Printf("📊 결과가 저장되었습니다: %s\n", *output)
	} else {
		reporter.Print(allResults)
	}

	// 전체 요약
	printOverallSummary(allResults)
}

func findCollectionFiles(dir string) ([]string, error) {
	var files []string
	
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".json") {
			files = append(files, path)
		}
		
		return nil
	})
	
	return files, err
}

func printOverallSummary(results []*TestSummary) {
	totalCollections := len(results)
	totalTests := 0
	totalPassed := 0
	totalFailed := 0
	successfulCollections := 0

	for _, result := range results {
		totalTests += result.TotalTests
		totalPassed += result.PassedTests
		totalFailed += result.FailedTests
		if result.FailedTests == 0 {
			successfulCollections++
		}
	}

	fmt.Println("=" + strings.Repeat("=", 50))
	fmt.Println("📋 전체 테스트 요약")
	fmt.Println("=" + strings.Repeat("=", 50))
	fmt.Printf("컬렉션: %d개 (성공: %d개)\n", totalCollections, successfulCollections)
	fmt.Printf("테스트: %d개 (성공: %d개, 실패: %d개)\n", totalTests, totalPassed, totalFailed)
	
	if totalFailed > 0 {
		fmt.Printf("🔴 전체 성공률: %.1f%%\n", float64(totalPassed)/float64(totalTests)*100)
		os.Exit(1)
	} else {
		fmt.Println("🟢 모든 테스트가 성공했습니다!")
	}
}

// 단일 컬렉션 실행 함수
func runSingleCollection(runner *Runner, file string, index, total int, verbose bool) *TestSummary {
	fmt.Printf("[%d/%d] %s 실행 중...\n", index, total, filepath.Base(file))
	
	collection, err := runner.LoadCollection(file)
	if err != nil {
		log.Printf("❌ 컬렉션 로드 실패: %s - %v", file, err)
		return nil
	}

	if verbose {
		fmt.Printf("  📄 컬렉션: %s\n", collection.Info.Name)
	}

	summary := runner.RunCollection(collection)
	summary.CollectionName = collection.Info.Name
	summary.FilePath = file

	// 간단한 결과 출력
	if summary.FailedTests > 0 {
		fmt.Printf("  ❌ %d개 실패 / %d개 총 테스트 (%.2fs)\n", 
			summary.FailedTests, summary.TotalTests, summary.TotalTime.Seconds())
	} else {
		fmt.Printf("  ✅ %d개 모두 성공 (%.2fs)\n", 
			summary.TotalTests, summary.TotalTime.Seconds())
	}
	fmt.Println()

	return summary
}

// 병렬 컬렉션 실행 함수
func runCollectionsInParallel(files []string, maxParallel int, verbose bool) []*TestSummary {
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]*TestSummary, 0, len(files))
	
	// 작업 채널과 워커 풀 생성
	jobs := make(chan string, len(files))
	
	// 워커 시작
	for w := 0; w < maxParallel; w++ {
		go func() {
			runner := NewRunner()
			for file := range jobs {
				result := processCollectionFile(runner, file, verbose)
				if result != nil {
					mu.Lock()
					results = append(results, result)
					mu.Unlock()
				}
				wg.Done()
			}
		}()
	}

	// 작업 큐에 파일들 추가
	for _, file := range files {
		wg.Add(1)
		jobs <- file
	}
	close(jobs)

	// 모든 작업 완료 대기
	wg.Wait()
	
	fmt.Printf("✅ %d개 컬렉션 병렬 실행 완료\n\n", len(results))
	return results
}

// 개별 컬렉션 파일 처리 (병렬용)
func processCollectionFile(runner *Runner, file string, verbose bool) *TestSummary {
	collection, err := runner.LoadCollection(file)
	if err != nil {
		log.Printf("❌ 컬렉션 로드 실패: %s - %v", file, err)
		return nil
	}

	if verbose {
		fmt.Printf("🔄 처리 중: %s\n", collection.Info.Name)
	}

	summary := runner.RunCollection(collection)
	summary.CollectionName = collection.Info.Name
	summary.FilePath = file

	status := "✅"
	if summary.FailedTests > 0 {
		status = "❌"
	}

	fmt.Printf("%s %s: %d/%d 성공 (%.2fs)\n", 
		status, filepath.Base(file), summary.PassedTests, summary.TotalTests, summary.TotalTime.Seconds())

	return summary
}

func printUsage() {
	fmt.Println("Postman Collection Tester")
	fmt.Println("")
	fmt.Println("사용법:")
	fmt.Printf("  %s [옵션]\n", os.Args[0])
	fmt.Println("")
	fmt.Println("옵션:")
	flag.PrintDefaults()
	fmt.Println("")
	fmt.Println("예시:")
	fmt.Printf("  %s                                    # ./postman 디렉토리의 모든 컬렉션 실행\n", os.Args[0])
	fmt.Printf("  %s -dir ./collections                 # 특정 디렉토리의 컬렉션 실행\n", os.Args[0])
	fmt.Printf("  %s -file test.json                    # 단일 파일 실행\n", os.Args[0])
	fmt.Printf("  %s -output report.html -format html   # HTML 리포트 생성\n", os.Args[0])
	fmt.Printf("  %s -parallel 3                        # 3개 컬렉션 동시 실행\n", os.Args[0])
	fmt.Printf("  %s -verbose                           # 상세 출력\n", os.Args[0])
}