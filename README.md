# Graph-Blog
Cервис для работы с постами и комментариями, реализованный с использованием GraphQL

### [Задание](./docs/task.docx)

## Описание

`GraphQL-Blog`— сервис для работы с постами и комментариями, реализованный с использованием GraphQL.

Посты могут иметь неограниченное количество вложенных комментариев. Автор поста может разрешать или запрещать добавление новых комментариев.

Для получения новых комментариев в реальном времени реализован GraphQL Subscription через WebSocket.

Сервис поддерживает два варианта хранения данных:

- `memory` — хранение данных в памяти;

- `postgres` — хранение данных в PostgreSQL.

Тип хранилища определяется при запуске через переменную окружения `REPOSITORY_TYPE`.

---

## Технологии
- **Go** — основной язык разработки
- **GraphQL** — API для взаимодействия с данными
- **PostgresSQL** — хранилище данных
- **Docker + Docker Compose** — контейнеризация и запуск окружения

---

## Функциональность

Посты:
- создание поста;
- получение списка постов;
- получение конкретного поста;
- включение/отключение комментариев автором поста.

Комментарии:
- создание комментариев;
- неограниченная вложенность комментариев;
- ответы на существующие комментарии;
- ограничение длины комментария — 2000 символов;
- пагинация комментариев;
- получение дочерних комментариев.

GraphQL Subscription:
- подписка на новые комментарии определённого поста;
- доставка комментариев через WebSocket без повторного выполнения GraphQL-запроса;
- несколько клиентов могут одновременно подписаться на один пост.

---

## Начало работы
### Установка
Клонирование репозитория
```sh
git clone https://github.com/AugustSerenity/GraphQL-Blog
```
### Запуск сервиса
Создаем .env файл с дефолтными настройками
```sh
make env
```
Настроки хранилища данных 
`REPOSITORY_TYPE`:
- =postgres для хранения данных в БД 
- =memory для хранения данных в inmemory

Eсли значение `REPOSITORY_TYPE` не задано, используется memory

Запускаем контейнер с помощью Makefile
```sh
make run
```

---

### Аутентификация

Для идентификации пользователя используется HTTP-заголовок:

`X-User-ID`

Например:

X-User-ID: user-1

Отдельной системы регистрации пользователей нет — переданный X-User-ID используется как идентификатор текущего пользователя.

## Пример использования

#### Создание поста
Запрос
```sh
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user-1" \
  -d '{
    "query": "mutation { createPost(input: { content: \"Создание поста\" }) { id content commentsEnabled author { id } } }"
  }'
```
Ответ 
```json
{"data":{"createPost":{"id":"843cb15a-2b69-4d68-acf4-847c299b825f","content":"Создание поста","commentsEnabled":true,"author":{"id":"user-1"}}}}
```

#### Получение списка постов
запрос
```sh
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user-1" \
  -d '{
    "query": "query { posts { id content commentsEnabled author { id } } }"
  }'
```

ответ
```json
{"data":{"posts":[{"id":"843cb15a-2b69-4d68-acf4-847c299b825f","content":"Создание поста","commentsEnabled":true,"author":{"id":"user-1"}},{"id":"196cf5b2-d62a-4cff-91da-7b42707bf8ee","content":"Post by user 2","commentsEnabled":true,"author":{"id":"user-2"}}]}}
```

#### Получение конкретного поста
запрос
```sh
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user-1" \
  -d '{
    "query": "query { post(id: \"843cb15a-2b69-4d68-acf4-847c299b825f\") { id content commentsEnabled author { id } } }"
  }'
```

ответ
```json
{"data":{"post":{"id":"843cb15a-2b69-4d68-acf4-847c299b825f","content":"Создание поста","commentsEnabled":true,"author":{"id":"user-1"}}}}
```

#### Создание комментария
запрос
```sh
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user-2" \
  -d '{
    "query": "mutation { createComment(input: { postId: \"843cb15a-2b69-4d68-acf4-847c299b825f\", content: \"Первый комментарий\" }) { id content author { id } parent { id } } }"
  }'
```

ответ
```json
{"data":{"createComment":{"id":"1cd64429-05d8-45c0-afa5-ff53e7056563","content":"Первый комментарий","author":{"id":"user-2"},"parent":null}}}
```

#### Создание вложенного комментария
запрос
```sh
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user-1" \
  -d '{
    "query": "mutation { createComment(input: { postId: \"843cb15a-2b69-4d68-acf4-847c299b825f\", parentId: \"1cd64429-05d8-45c0-afa5-ff53e7056563\", content: \"Ответ на первый комментарий\" }) { id content author { id } parent { id content } } }"
  }'
```

ответ
```json
{"data":{"createComment":{"id":"5dd5342a-8506-4cd1-9043-361c94b920a6","content":"Ответ на первый комментарий","author":{"id":"user-1"},"parent":{"id":"1cd64429-05d8-45c0-afa5-ff53e7056563","content":""}}}}
```

#### Получение комментариев
запрос
```sh
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user-1" \
  -d '{
    "query": "query { comments(postId: \"843cb15a-2b69-4d68-acf4-847c299b825f\") { items { id content author { id } parent { id } children { id content author { id } } } pageInfo { hasNextPage endCursor } } }"
  }'
```

ответ
```json
{"data":{"comments":{"items":[{"id":"1cd64429-05d8-45c0-afa5-ff53e7056563","content":"Первый комментарий","author":{"id":"user-2"},"parent":null,"children":[{"id":"5dd5342a-8506-4cd1-9043-361c94b920a6","content":"Ответ на первый комментарий","author":{"id":"user-1"}}]}],"pageInfo":{"hasNextPage":false,"endCursor":"1cd64429-05d8-45c0-afa5-ff53e7056563"}}}}
```

#### Проверка пагинации
Создал 6 комментариев и указал limit 2, в запросе указал вернуть 2 страницу,
то есть комментарий 3 и 4, так же аналогично с первой страницей и последней страницей 
```sh
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user-1" \
  -d '{
    "query": "query { comments(postId: \"843cb15a-2b69-4d68-acf4-847c299b825f\", pagination: { limit: 2, cursor: \"5546ccc9-e79e-474e-aca5-7dbf9d32fa50\" }) { items { id content author { id } children { id content author { id } } } pageInfo { hasNextPage endCursor } } }"
  }'
```
Ответ
```json
{"data":{"comments":{"items":[{"id":"c8b7721f-0890-438e-a12f-a22ad7befe02","content":"Комментарий для проверки пагинации 4","author":{"id":"user-2"},"children":[]},{"id":"d5ed7fae-aee4-429e-871a-77f60e3db04e","content":"Комментарий для проверки пагинации 3","author":{"id":"user-2"},"children":[]}],"pageInfo":{"hasNextPage":true,"endCursor":"d5ed7fae-aee4-429e-871a-77f60e3db04e"}}}}
```

#### Запрет комментариев
запрос
```sh
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user-1" \
  -d '{
    "query": "mutation { setCommentsEnabled(postId: \"843cb15a-2b69-4d68-acf4-847c299b825f\", enabled: false) { id commentsEnabled } }"
  }'
```

ответ
```json
{"data":{"setCommentsEnabled":{"id":"843cb15a-2b69-4d68-acf4-847c299b825f","commentsEnabled":false}}}
```
Автор `user-1` успешно отключил комментарии

#### Попытка отключить/включить комментарии другим пользователем
запрос
```sh
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user-2" \
  -d '{
    "query": "mutation { setCommentsEnabled(postId: \"843cb15a-2b69-4d68-acf4-847c299b825f\", enabled: true) { id commentsEnabled } }"
  }'
```

ответ
```json
{"errors":[{"message":"user is not post author","path":["setCommentsEnabled"],"locations":[{"line":1,"column":12}]}],"data":null}
```

#### Попытка создать комментарий после запрета
запрос
```sh
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user-2" \
  -d '{
    "query": "mutation { createComment(input: { postId: \"843cb15a-2b69-4d68-acf4-847c299b825f\", content: \"Комментарий после запрета\" }) { id content author { id } } }"
  }'
```

ответ
```json
{"errors":[{"message":"comments are disabled","path":["createComment"],"locations":[{"line":1,"column":12}]}],"data":null}
```

#### Снова разрешаем комментарии
Делаем это от имени владельца поста `user-1`:

запрос
```sh
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user-1" \
  -d '{
    "query": "mutation { setCommentsEnabled(postId: \"843cb15a-2b69-4d68-acf4-847c299b825f\", enabled: true) { id commentsEnabled } }"
  }'
```

ответ
```json
{"data":{"setCommentsEnabled":{"id":"843cb15a-2b69-4d68-acf4-847c299b825f","commentsEnabled":true}}}
```

#### Подписка на новые комментарии
Для проверки данного функционала необходимо подключиться к сервису через 
WebSocket. Не забываем заголовок 

Например 
```sh
websocat -v \
  --protocol graphql-transport-ws \
  -H="X-User-ID: user-1" \
  ws://localhost:8080/query
```
Инициализируем подключение к graphQL
```sh
{"type":"connection_init"}
```
получаем ответ
```sh
{"type":"connection_ack"}
```
Делаем в терминале с webSocket подписку пост
```json
{"id":"1","type":"subscribe","payload":{"query":"subscription { commentAdded(postId: \"843cb15a-2b69-4d68-acf4-847c299b825f\") { id content author { id } post { id } } }"}}
```
Теперь создаем комментарий после подписка на пост
```sh
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user-2" \
  -d '{
    "query": "mutation { createComment(input: { postId: \"843cb15a-2b69-4d68-acf4-847c299b825f\", content: \"Комментарий через Subscription\" }) { id content author { id } } }"
  }'
```
И после публикации комментария получаем в терминале с webSocket событие 
```json
{"payload":{"data":{"commentAdded":{"id":"1d4a58f5-222b-42b3-b66c-870601822ae4","content":"Комментарий через Subscription","author":{"id":"user-2"},"post":{"id":"843cb15a-2b69-4d68-acf4-847c299b825f"}}}},"id":"1","type":"next"}
```
#### Проверим несколько подписчиков
Создадим еще одно подключение user-2 по webSocket, подпишимся на тот же пост
Те же самые дейстивия по подключению произведем, только с заголовком `user-2`

Далее оформим подписку на пост
```sh
{"id":"2","type":"subscribe","payload":{"query":"subscription { commentAdded(postId: \"843cb15a-2b69-4d68-acf4-847c299b825f\") { id content author { id } post { id } } }"}}
```

Создаю еще комментарий
```sh
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user-2" \
  -d '{
    "query": "mutation { createComment(input: { postId: \"843cb15a-2b69-4d68-acf4-847c299b825f\", content: \"Комментарий для двух подписчиков\" }) { id content } }"
  }'
```

И получаю ответ для каждого подписчика

первый
```json
{"payload":{"data":{"commentAdded":{"id":"13037646-dc1a-4759-9a13-7b7f31d7a39e","content":"Комментарий для двух подписчиков","author":{"id":"user-2"},"post":{"id":"843cb15a-2b69-4d68-acf4-847c299b825f"}}}},"id":"1","type":"next"}
```

второй 
```json
{"payload":{"data":{"commentAdded":{"id":"13037646-dc1a-4759-9a13-7b7f31d7a39e","content":"Комментарий для двух подписчиков","author":{"id":"user-2"},"post":{"id":"843cb15a-2b69-4d68-acf4-847c299b825f"}}}},"id":"2","type":"next"}
```

#### Проверим ограничение длины комментария в 2000 символов
создадим строку в терминале длиной 2001 символ
```sh
CONTENT=$(printf 'a%.0s' {1..2001})
```
Запрос 
```sh
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user-2" \
  -d "{\"query\":\"mutation { createComment(input: { postId: \\\"843cb15a-2b69-4d68-acf4-847c299b825f\\\", content: \\\"$CONTENT\\\" }) { id content } }\"}"
```
Ответ
```json
{"errors":[{"message":"content is invalid","path":["createComment"],"locations":[{"line":1,"column":12}]}],"data":null}
```

#### Проверим ровно 2000 символов
создадим строку в терминале
```sh
CONTENT=$(printf 'a%.0s' {1..2000})
```
Запрос 
```sh
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user-2" \
  -d "{\"query\":\"mutation { createComment(input: { postId: \\\"843cb15a-2b69-4d68-acf4-847c299b825f\\\", content: \\\"$CONTENT\\\" }) { id content } }\"}"
```
Ответ
```json
"data":{"createComment":{"id":"91e8414c-1a31-4561-adc2-5486484a895f","content":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
```

#### Проверка создания комментария к не существующему посту и несуществующему комментарию

Запрос 
```sh
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user-2" \
  -d '{
    "query": "mutation { createComment(input: { postId: \"00000000-0000-0000-0000-000000000000\", content: \"Комментарий к несуществующему посту\" }) { id content } }"
  }'
```
Ответ
```json
{"errors":[{"message":"not found","path":["createComment"],"locations":[{"line":1,"column":12}]}],"data":null}
```

и к не существующему комментарию parenID

Запрос 
```sh
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "X-User-ID: user-2" \
  -d '{
    "query": "mutation { createComment(input: { postId: \"843cb15a-2b69-4d68-acf4-847c299b825f\", parentId: \"00000000-0000-0000-0000-000000000000\", content: \"Ответ на несуществующий комментарий\" }) { id content parent { id } } }"
  }'
```
  
Ответ
```json
{"errors":[{"message":"not found","path":["createComment"],"locations":[{"line":1,"column":12}]}],"data":null}
```