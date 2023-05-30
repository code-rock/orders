Запуск сервера 
nats-streaming-server --config cluster.conf

Запуск отправки сообщений
go run ./cmd/message-pablesher

Запуск
go run ./cmd/app

Пример env файла
POSTGRE_SQL_HOST = "localhost"
POSTGRE_SQL_PORT = 5432
POSTGRE_SQL_USER = "admin"
POSTGRE_SQL_PASSWORD = "12345"
POSTGRE_SQL_NAME = "orders"