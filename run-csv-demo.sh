#!/bin/bash

# Postman Collection Tester - CSV Demo Script
# 데모용 스크립트: 컬렉션을 실행하고 CSV로 결과 저장

echo "🚀 Postman Collection Tester - CSV Demo"
echo "======================================="

# 현재 시간으로 파일명 생성
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
OUTPUT_FILE="demo_results_${TIMESTAMP}.csv"

# 실행파일 존재 확인
if [ ! -f "./postman-tester" ]; then
    echo "❌ postman-tester 실행파일을 찾을 수 없습니다."
    echo "   먼저 'go build -o postman-tester' 명령을 실행해주세요."
    exit 1
fi

# 테스트 컬렉션 파일 존재 확인
if [ ! -f "./test-collection.json" ]; then
    echo "❌ test-collection.json 파일을 찾을 수 없습니다."
    exit 1
fi

echo "📁 테스트 파일: test-collection.json"
echo "📊 출력 파일: ${OUTPUT_FILE}"
echo ""

# 컬렉션 실행 (CSV 형식으로 저장)
echo "🔄 테스트 실행 중..."
./postman-tester -file test-collection.json -format csv -output "${OUTPUT_FILE}" -verbose

# 결과 확인
if [ $? -eq 0 ]; then
    echo ""
    echo "✅ 테스트 완료!"
    echo "📄 결과 파일: ${OUTPUT_FILE}"
    
    # CSV 파일의 헤더와 첫 몇 줄 미리보기
    echo ""
    echo "📋 CSV 파일 미리보기:"
    echo "--------------------"
    head -5 "${OUTPUT_FILE}"
    
    # 파일 크기 정보
    FILE_SIZE=$(ls -lh "${OUTPUT_FILE}" | awk '{print $5}')
    LINE_COUNT=$(wc -l < "${OUTPUT_FILE}")
    echo ""
    echo "📊 파일 정보:"
    echo "   크기: ${FILE_SIZE}"
    echo "   라인 수: ${LINE_COUNT}줄 (헤더 포함)"
    
    # Excel에서 열기 안내
    echo ""
    echo "💡 사용법:"
    echo "   - Excel에서 열기: Excel > 파일 > 열기 > ${OUTPUT_FILE}"
    echo "   - 터미널에서 보기: cat ${OUTPUT_FILE}"
    echo "   - 스프레드시트로 보기: open -a 'Numbers' ${OUTPUT_FILE}"
    
else
    echo "❌ 테스트 실행 중 오류가 발생했습니다."
    exit 1
fi