# Balance manager
Мкиросервис для выполнения транзакций с балансом пользователя.

## Запуск 
Для запуска нужно склонировать репозиторий:
```
git clone https://github.com/Kulallador/balance_manager.git
```
И выполнить команду:
```
docker compose up
```

## Запросы
1. GET: /balance. Запрос для получения баланса пользователя по user_id. Принимает {user_id int}, возвращает {user_id int, balance float}.
2. POST: /balance/inc. Запрос для пополнения баланса пользователя. Принимает {user_id int, money float}, возвращает код "200" в случае успешной транзакции и "400" в случае ошибки. 
3. POST: /balance/dec. Запрос для извлечения средств с баланса пользователя. Принимает {user_id int, money float}, возвращает код "200" в случае успешной транзакции и "400" в случае ошибки. 
4. POST: /balance/translate. Запрос для перевода средств от одного пользователя (from_id) другому (to_id). Принимает {from_id int, to_id int, money float}, возвращает код "200" в случае успешной транзакции и "400" в случае ошибки. 
5. GET: /reserve. Запрос для получения баланса резерва. Принимает {user_id int, service_id int, order_id int}, возвращает {user_id int, service_id int, order_id int, money float}.
6. POST: /reserve/inc. Запрос для резервирования средств. Принимает {user_id int, service_id int, order_id int, money float}, возвращает код "200" в случае успешной транзакции и "400" в случае ошибки. 
7. POST: /reserve/dec. Запрос для извлечения средств из резерва. Принимает {user_id int, service_id int, order_id int, money float}, возвращает код "200" в случае успешной транзакции и "400" в случае ошибки.
