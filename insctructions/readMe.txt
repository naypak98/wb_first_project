1. Запустить Docker Desktop.
2. В корне проекта выполнить:
   docker-compose up -d
3. Создать таблицы в БД:
   - Подключиться к PostgreSQL (логин/пароль: testuser/testpass, база orders_db)
   - Выполнить SQL-скрипт schema.sql
4. Запустить приложение (Go-сервис):
   go run cmd/main.go
5. Открыть Kafka UI в браузере: http://localhost:8082/
6. Создать топик "orders"
7. Отправить сообщение в Kafka (тело сообщения можно взять из model.json)
8. Открыть интерфейс orders_front.html
9. Ввести order_uid из JSON и нажать кнопку "Get order"
10. Должен отобразиться полный JSON заказа из кеша/БД