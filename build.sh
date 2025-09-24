#!/bin/bash

# Postman Collection Tester 크로스 플랫폼 빌드 스크립트

set -e

echo "🔨 Postman Collection Tester 빌드 시작..."

# 빌드할 플랫폼들
PLATFORMS=(
    "windows/amd64"
    "windows/386" 
    "linux/amd64"
    "linux/386"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
)

# 빌드 디렉토리 생성
BUILD_DIR="builds"
if [ -d "$BUILD_DIR" ]; then
    rm -rf "$BUILD_DIR"
fi
mkdir -p "$BUILD_DIR"

# 각 플랫폼별로 빌드
for platform in "${PLATFORMS[@]}"; do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    
    output_name="postman-tester"
    
    if [ $GOOS = "windows" ]; then
        output_name+=".exe"
    fi
    
    output_path="$BUILD_DIR/${output_name}-${GOOS}-${GOARCH}"
    if [ $GOOS = "windows" ]; then
        output_path+=".exe"
    fi
    
    echo "🔄 빌드 중: $GOOS/$GOARCH -> $output_path"
    
    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o $output_path
    
    if [ $? -ne 0 ]; then
        echo "❌ 빌드 실패: $platform"
        exit 1
    fi
done

# README 파일 생성
cat > "$BUILD_DIR/README.md" << 'EOF'
# Postman Collection Tester

## 사용법

### Windows
```
postman-tester-windows-amd64.exe
postman-tester-windows-amd64.exe -dir ./collections
postman-tester-windows-amd64.exe -output report.html -format html
```

### macOS
```
# Intel Mac
./postman-tester-darwin-amd64

# Apple Silicon Mac
./postman-tester-darwin-arm64
```

### Linux
```
# 64bit
./postman-tester-linux-amd64

# 32bit  
./postman-tester-linux-386

# ARM64
./postman-tester-linux-arm64
```

## 옵션
- `-dir`: 컬렉션 파일 디렉토리 (기본: ./postman)
- `-output`: 결과 파일명 (선택사항)
- `-format`: 출력 형식 (text, json, html)
- `-parallel`: 병렬 실행 수 (기본: 1)
- `-verbose`: 상세 출력
- `-help`: 도움말

## 예시
```bash
# 기본 실행
./postman-tester

# HTML 리포트 생성
./postman-tester -output report.html -format html

# 3개 컬렉션 동시 실행
./postman-tester -parallel 3 -verbose
```
EOF

echo ""
echo "✅ 빌드 완료!"
echo ""
echo "📦 생성된 파일들:"
ls -la "$BUILD_DIR"

# 파일 크기 요약
echo ""
echo "📊 바이너리 크기:"
for file in "$BUILD_DIR"/postman-tester-*; do
    if [ -f "$file" ]; then
        size=$(du -h "$file" | cut -f1)
        echo "  $(basename "$file"): $size"
    fi
done

echo ""
echo "🎉 모든 플랫폼 빌드 완료! builds/ 디렉토리를 확인하세요."