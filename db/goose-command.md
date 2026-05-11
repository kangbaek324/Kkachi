goose -dir db/migrations postgres "postgres://root:1234@localhost:5432/kkachi?sslmode=disable" up

goose -dir db/migrations create [name] sql
