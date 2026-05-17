# KKachi 🏦

Go 기반의 다중통화 지갑 서비스. **Gin REST API 백엔드**와 **Bubble Tea TUI 프론트엔드**로 구성된 Go 모노레포입니다.

## Monorepo Layout

```
.
├── go.work                  # 두 모듈을 묶는 워크스페이스
├── apps/
│   ├── backend/             # Gin + PostgreSQL REST API
│   └── frontend/            # Bubble Tea 기반 TUI 클라이언트
├── CLAUDE.md
└── README.md
```

- `apps/backend` 모듈: `github.com/kangbaek324/kkachi/apps/backend`
- `apps/frontend` 모듈: `github.com/kangbaek324/kkachi/apps/frontend`

## Stack

| 분류      | Backend                 | Frontend   |
| --------- | ----------------------- | ---------- |
| Language  | Go 1.26+                | Go 1.26+   |
| Framework | Gin                     | Bubble Tea |
| Styling   | –                       | Lipgloss   |
| Database  | PostgreSQL              | –          |
| DB Driver | pgx/v5                  | –          |
| Query Gen | sqlc                    | –          |
| Migration | goose                   | –          |
| Auth      | JWT (golang-jwt/jwt v5) | –          |
| Config    | godotenv                | os.Getenv  |

## Feature

### Auth (backend)

- 회원가입 / 로그인 (JWT)

### Wallet (backend)

- 지갑(계좌) 개설
- 계좌 목록 및 계좌 내역 조회
- 거래 내역 조회
- 송금 / 환전

### Currency (backend)

- 서비스 기준 환율 조회

## Getting Started

### 사전 요구사항

- Go 1.22+
- PostgreSQL
- [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html) (백엔드 SQL 변경 시)
- [goose](https://github.com/pressly/goose#install) (마이그레이션)

### 워크스페이스 셋업

```bash
# 모듈 다운로드
go mod download -C apps/backend
go mod download -C apps/frontend
```

`go.work`가 루트에 있어 두 모듈이 같은 빌드 그래프 안에서 동작합니다.

### 백엔드 환경 변수

```bash
cp apps/backend/.env.example apps/backend/.env
```

`apps/backend/.env` 내용:

```env
DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable
PORT=8080
JWT_SECRET=your_secret_here
GIN_MODE=debug              # 프로덕션은 release
API_KEY=your_koreaexim_apikey_here
```

### 백엔드 실행

```bash
# DB 마이그레이션
goose -dir apps/backend/db/migrations postgres "$DATABASE_URL" up

# 서버 실행 (저장소 루트에서)
go run ./apps/backend/cmd
```

### 프론트엔드(TUI) 실행

```bash
# 백엔드가 다른 호스트/포트에 있다면 지정
export KKACHI_API_URL=http://localhost:8080

go run ./apps/frontend/cmd
```

종료는 `q`, `esc`, 또는 `ctrl+c`.

### 빌드

```bash
go build -o bin/server ./apps/backend/cmd
go build -o bin/tui    ./apps/frontend/cmd
```

## Folder Structure

```
.
├── go.work
├── CLAUDE.md
├── README.md
└── apps/
    ├── backend/
    │   ├── cmd/
    │   │   └── main.go                 # 진입점
    │   ├── db/
    │   │   ├── postgres.go             # DB 커넥션 풀
    │   │   ├── schema.sql              # 스키마 참조용
    │   │   ├── migrations/             # goose 마이그레이션 (sqlc 스키마 소스)
    │   │   ├── queries/                # sqlc 입력 SQL
    │   │   │   ├── user.sql
    │   │   │   ├── wallet.sql
    │   │   │   └── currency.sql
    │   │   └── sqlc/                   # sqlc 생성 코드 (도메인 공유)
    │   ├── internal/
    │   │   ├── common/                 # 공통 Response, error 등
    │   │   ├── config/                 # 환경변수 로딩
    │   │   ├── middleware/             # JWT 미들웨어
    │   │   └── domain/
    │   │       ├── user/
    │   │       ├── wallet/
    │   │       └── currency/
    │   ├── routes/
    │   │   └── routes.go               # 라우트 등록, 도메인 위임
    │   ├── .env.example
    │   ├── .air.toml
    │   ├── sqlc.yaml
    │   └── go.mod
    └── frontend/
        ├── cmd/
        │   └── main.go                 # 진입점
        ├── internal/
        │   ├── app/                    # tea.Model / Update / View
        │   └── config/                 # 환경변수 로딩
        └── go.mod
```

## DB Schema

```
users          — 사용자 계정 (username, password, role)
wallets        — 계좌 (account_number, nickname)
balances       — 통화별 잔액 (account_id × currency_id)
currencies     — 지원 통화 목록 (code, name, unit)
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
