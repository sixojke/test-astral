На реализацию кеша меня не хватило. Слишком объёмное тз.

## Запуск 

Cоздайте файл .env, запишите эти переменные
```env
POSTGRES_DB=documents
POSTGRES_USER=user
POSTGRES_PASSWORD=password

AUTH_ADMIN_TOKEN=sfuqwejqjoiu93e29
AUTH_SIGNING_KEY=fnweosiupfhjpioe

HASHER_SALT=43kolpcqjrq3v4rpr
```

Чтобы запустить приложение пропишите make up

## Миграции

Миграции лежат в папке schema/postgres. Накатываются сами

## Swagger

Перейдите по этой ссылке, чтобы открыть документацию 
```link
http://localhost:8080/api/swagger/index.html#/docs/get_docs
```