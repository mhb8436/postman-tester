@echo off
REM Postman Collection Tester 크로스 플랫폼 빌드 스크립트 (Windows용)

echo 🔨 Postman Collection Tester 빌드 시작...

REM 빌드 디렉토리 정리
if exist builds rmdir /s /q builds
mkdir builds

echo 🔄 Windows 64bit 빌드 중...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o builds\postman-tester-windows-amd64.exe
if errorlevel 1 goto :error

echo 🔄 Windows 32bit 빌드 중...
set GOOS=windows
set GOARCH=386
go build -ldflags="-s -w" -o builds\postman-tester-windows-386.exe
if errorlevel 1 goto :error

echo 🔄 Linux 64bit 빌드 중...
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w" -o builds\postman-tester-linux-amd64
if errorlevel 1 goto :error

echo 🔄 Linux ARM64 빌드 중...
set GOOS=linux
set GOARCH=arm64
go build -ldflags="-s -w" -o builds\postman-tester-linux-arm64
if errorlevel 1 goto :error

echo 🔄 macOS Intel 빌드 중...
set GOOS=darwin
set GOARCH=amd64
go build -ldflags="-s -w" -o builds\postman-tester-darwin-amd64
if errorlevel 1 goto :error

echo 🔄 macOS Apple Silicon 빌드 중...
set GOOS=darwin
set GOARCH=arm64
go build -ldflags="-s -w" -o builds\postman-tester-darwin-arm64
if errorlevel 1 goto :error

REM README 파일 생성
(
echo # Postman Collection Tester
echo.
echo ## Windows에서 실행
echo ```
echo postman-tester-windows-amd64.exe
echo postman-tester-windows-amd64.exe -dir .\collections
echo postman-tester-windows-amd64.exe -output report.html -format html
echo ```
echo.
echo ## 옵션
echo - `-dir`: 컬렉션 파일 디렉토리 (기본: .\postman^)
echo - `-output`: 결과 파일명 (선택사항^)
echo - `-format`: 출력 형식 (text, json, html^)
echo - `-parallel`: 병렬 실행 수 (기본: 1^)
echo - `-verbose`: 상세 출력
echo - `-help`: 도움말
) > builds\README.md

echo.
echo ✅ 빌드 완료!
echo.
echo 📦 생성된 파일들:
dir /b builds

echo.
echo 🎉 모든 플랫폼 빌드 완료! builds\ 디렉토리를 확인하세요.
goto :end

:error
echo ❌ 빌드 실패!
exit /b 1

:end