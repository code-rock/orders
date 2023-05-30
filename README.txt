Запуск сервера 
nats-streaming-server --config cluster.conf

Запуск отправки сообщений
go run ./cmd/message-pablesher

Запуск
go run ./cmd/app