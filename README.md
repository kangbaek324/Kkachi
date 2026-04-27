# KKachi 🏦

Go 기반의 간단한 다중통화 지갑 REST API 서버입니다. 사용자가 계좌를 개설하고, 다양한 통화로 잔액을 관리하며 송금할 수 있습니다.

## Stack

| 분류 | 기술 |
|------|------|
| Language | Go 1.26+ |
| Framework | Gin |
| Database | PostgreSQL |
| DB Driver | pgx/v5 |
| Query Gen | sqlc |
| Migration | goose |
| Auth | JWT (golang-jwt/jwt v5) |
| Config | godotenv |

## Feature

### Auth

- 회원가입 / 로그인 (JWT)

### Wallet

- 지갑(계좌) 개설
- 계좌 목록 및 계좌 내역 조회
- 거래 내역 조회
- 송금

### Currency

- 서비스 기준 환율 조회

## Getting Started

### 사전 요구사항

- Go 1.22+
- PostgreSQL
- [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html)
- [goose](https://github.com/pressly/goose#install)

### 환경 변수 설정

프로젝트 루트에 `.env` 파일을 생성합니다.

```bash
cp .env.example .env
```

`.env` 파일 내용:

```env
DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
PORT=8080
JWT_SECRET=your_secret_here
GIN_MODE=debug   # 프로덕션 환경에서는 release
```

### 실행

```bash
# 의존성 설치
go mod tidy

# DB 마이그레이션
goose -dir db/migrations postgres $DATABASE_URL up

# 서버 실행
go run cmd/main.go
```

### 빌드

```bash
go build -o bin/server cmd/main.go
./bin/server
```

## Folder Structure

```
.
├── cmd/
│   └── main.go                     # 진입점
├── db/
│   ├── postgres.go                 # DB 커넥션 풀
│   ├── schema.sql                  # 스키마 참조용
│   ├── migrations/                 # goose 마이그레이션 파일 (sqlc 스키마 소스)
│   ├── queries/                    # sqlc용 SQL 쿼리 파일
│   │   ├── user.sql
│   │   ├── wallet.sql
│   │   └── currency.sql
│   └── sqlc/                       # sqlc 생성 코드 (전 도메인 공유)
│       ├── db.go
│       ├── models.go
│       ├── querier.go
│       └── *.sql.go
├── internal/
│   ├── common/
│   │   └── response.go             # 공통 Response 구조체
│   ├── config/
│   │   └── config.go               # 환경변수 로딩
│   ├── middleware/
│   │   └── auth.go                 # JWT 인증 미들웨어
│   └── domain/
│       ├── user/                   # 사용자 도메인
│       │   ├── handler.go
│       │   ├── service.go
│       │   ├── routes.go
│       │   └── dto.go              # 요청/응답 타입
│       ├── wallet/                 # 지갑/잔액 도메인
│       │   ├── handler.go
│       │   ├── service.go
│       │   └── routes.go
│       └── currency/               # 통화/환율 도메인
│           ├── handler.go
│           ├── service.go
│           └── routes.go
├── routes/
│   └── routes.go                   # 라우트 등록, 도메인별 위임
├── .env.example
├── sqlc.yaml
└── go.mod
```

## DB Schema

```
users         — 사용자 계정 (username, password, role)
wallets       — 계좌 (account_number, nickname)
balances      — 통화별 잔액 (account_id × currency_id)
currencies    — 지원 통화 목록 (code, name, unit)
exchange_rates — 환율 (currency_id, rate)
```

## API Response Format

모든 응답은 아래 형식을 따릅니다.

```json
// 성공
{ "code": 200, "success": true, "message": "ok", "data": { ... } }

// 실패
{ "code": 400, "success": false, "message": "invalid request" }
```
