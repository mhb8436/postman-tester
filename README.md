# Postman Collection Tester

Postman 컬렉션을 일괄 테스트할 수 있는 크로스 플랫폼 CLI 도구입니다.

## 🚀 주요 기능

- **크로스 플랫폼**: Windows, macOS, Linux 지원
- **다양한 입력 방식**: 단일 파일 또는 디렉토리 전체 테스트
- **여러 출력 형식**: 텍스트, JSON, HTML 리포트
- **병렬 실행**: 여러 컬렉션 동시 처리
- **상세한 결과**: 응답 시간, 상태 코드, 오류 메시지 포함

## 📦 설치 및 빌드

### Go 소스에서 빌드
```bash
go build -o postman-tester
```

### 크로스 플랫폼 빌드
```bash
# Linux/macOS에서
./build.sh

# Windows에서
build.bat
```

### Windows 전용 빌드 (맥/리눅스에서)
```bash
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o postman-tester-windows.exe
```

## 💻 Windows 11 사용법

### 준비사항
1. `postman-tester-windows.exe` (실행 파일)
2. `test-collection.json` (테스트용 컬렉션)
3. 본인의 Postman JSON 컬렉션 파일들

### 기본 사용법

#### 1. 파일 준비
```
C:\api-test\
├── postman-tester-windows.exe
├── test-collection.json
└── 본인파일.json
```

#### 2. 터미널 실행
- Windows 키 + R → `cmd` 입력
- 또는 PowerShell 실행

#### 3. 폴더 이동
```cmd
cd C:\Users\사용자명\Downloads\api-test
```

#### 4. 테스트 실행

**기본 테스트:**
```cmd
postman-tester-windows.exe -file test-collection.json
```

**상세 출력:**
```cmd
postman-tester-windows.exe -file test-collection.json -verbose
```

**본인 파일 테스트:**
```cmd
postman-tester-windows.exe -file "내파일.json"
```

**HTML 리포트 생성:**
```cmd
postman-tester-windows.exe -file test-collection.json -output report.html -format html
```

**JSON 리포트 생성:**
```cmd
postman-tester-windows.exe -file test-collection.json -output report.json -format json
```

**디렉토리 전체 테스트:**
```cmd
postman-tester-windows.exe -dir postman
```

**병렬 실행:**
```cmd
postman-tester-windows.exe -dir postman -parallel 3 -verbose
```

**도움말:**
```cmd
postman-tester-windows.exe -help
```

## 🔧 명령줄 옵션

| 옵션 | 설명 | 기본값 |
|------|------|--------|
| `-file` | 단일 Postman 컬렉션 파일 | - |
| `-dir` | 컬렉션 파일 디렉토리 | `./postman` |
| `-output` | 결과 저장 파일명 | 콘솔 출력 |
| `-format` | 출력 형식 (text, json, html) | `text` |
| `-parallel` | 병렬 실행 수 | `1` |
| `-timeout` | 요청 타임아웃(초) | `30` |
| `-verbose` | 상세 출력 | `false` |
| `-help` | 도움말 표시 | `false` |

## 📊 출력 예시

### 성공적인 실행
```
🚀 1개의 Postman 컬렉션을 테스트합니다...

[1/1] test-collection.json 실행 중...
  📄 컬렉션: JSONPlaceholder API 테스트
  ✅ 10개 모두 성공 (2.34s)

📋 전체 테스트 요약
컬렉션: 1개 (성공: 1개)
테스트: 10개 (성공: 10개, 실패: 0개)
🟢 모든 테스트가 성공했습니다!
```

### 상세 출력 (-verbose)
```
  [1.1] ✅ 모든 게시물 조회
        GET https://jsonplaceholder.typicode.com/posts
        응답: HTTP 200 (0.22s)

  [1.2] ✅ 특정 게시물 조회
        GET https://jsonplaceholder.typicode.com/posts/1
        응답: HTTP 200 (0.07s)
```

## 🧪 테스트용 컬렉션

`test-collection.json`은 JSONPlaceholder API를 사용한 테스트 컬렉션입니다:
- GET, POST, PUT, DELETE 요청 포함
- 10개의 다양한 API 테스트
- 공개 API로 즉시 테스트 가능

## ❓ 문제 해결

### "파일을 찾을 수 없습니다"
- 파일 경로 확인
- 파일명에 공백이 있으면 따옴표 사용: `"my collection.json"`

### "액세스가 거부되었습니다"
- 관리자 권한으로 터미널 실행
- 다른 폴더에서 실행 시도

### 네트워크 오류
- 방화벽/프록시 설정 확인
- VPN 연결 해제 후 재시도

## 🏗️ 프로젝트 구조

```
api-test/
├── main.go              # CLI 인터페이스
├── postman.go           # Postman 구조체 정의
├── runner.go            # HTTP 요청 실행 엔진
├── reporter.go          # 리포트 생성
├── build.sh             # Unix 빌드 스크립트
├── build.bat            # Windows 빌드 스크립트
├── test-collection.json # 테스트용 컬렉션
├── postman/            # 원본 컬렉션들
└── CLAUDE.md           # 개발 가이드
```

## 📝 라이선스

MIT License

## 🤝 기여

이슈 및 PR을 환영합니다!

---

**마지막 업데이트**: 2024-09-24