Запуск сервера 
nats-streaming-server --config cluster.conf

Запуск отправки сообщений
go build ./cmd/message-pablesher
go run ./cmd/message-pablesher

Запуск
go build ./cmd/app
go run ./cmd/app