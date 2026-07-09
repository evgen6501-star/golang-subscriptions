# 📦 Subscription Service API

Сервис для управления подписками пользователей. Реализован на **Go** с использованием **PostgreSQL**, **Docker** и **Swagger**.  
Проект выполнен в рамках тестового задания.

---

##  Функциональность

- операции над подписками:
  - Создание (`POST /subscriptions`)
  - Получение по ID (`GET /subscriptions/{id}`)
  - Обновление (`PUT /subscriptions/{id}`)
  - Удаление (`DELETE /subscriptions/{id}`)
  - Список подписок с пагинацией (`GET /subscriptions?limit=10&offset=0`)
  - Подсчёт суммы подписок за период с фильтрацией по пользователю и сервису (`GET /subscriptions/total`)

---

## Запуск

### Локально (без Docker)

```bash
make run
```
### Через Docker Compose

```bash
make up
```
### Остановка контейнера

```bash
make down
```
### Полная очистка

```bash
make clean
```