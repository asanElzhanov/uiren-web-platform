

# 📘 API Endpoints

## 🔐 Basic Authentication

### `GET /api/sign-in`

**Request Body (JSON):**

```json
{
  "identificator": "asanelzhanov2@gmail.com",
  "password": "Pass@123456"
}
```

---

### `GET /api/register`

**Request Body (JSON):**

```json
{
  "username": "seab",
  "email": "seab@seab.ru",
  "password": "Pass@123456"
}
```

---

### `GET /api/verify/:username/:code`

**Path Parameters:**

* `:username` — имя пользователя
* `:code` — код подтверждения

---

### `GET /api/refresh-token?refresh_token=...`

**Query Parameters:**

* `refresh_token` — токен обновления

---

## 👤 Users (Admin Only)

### `GET /api/users/:id`

Получить информацию о пользователе по ID.

---

### `POST /api/users`

Создать нового пользователя.
**Request Body (JSON):**

```json
{
  "username": "seab",
  "email": "seab@seab.ru",
  "password": "Pass@123456"
}
```

---

### `PATCH /api/users/:id`

Обновить данные пользователя.
**Request Body (JSON):**

```json
{
  "first_name": "first",
  "last_name": "last",
  "phone": "+77474568595",
  "phone_region": "KZ"
}
```

---


---

## 📦 Modules (Admin Only)

### `GET /api/modules/:code`

Получить информацию о модуле по его коду.

---

### `POST /api/modules`

Создать новый модуль.
**Request Body (JSON):**

```json
{
  "code": "new",
  "title": "Основы программирования",
  "description": "Этот модуль охватывает базовые концепции программирования.",
  "goal": "Научиться писать простые программы.",
  "difficulty": "beginner",
  "unlock_requirements": {
    "previous_module": "",
    "min_xp": 0
  },
  "reward": {
    "xp": 100,
    "badge": "Beginner Coder"
  },
  "lessons": [
    "lesson_001",
    "lesson_002",
    "lesson_003"
  ]
}
```

---

### `PATCH /api/modules/:code`

Обновить данные модуля (частично или полностью).
**Request Body (JSON):**

```json
{
  "title": "Basics of Kazakh",
  "description": "Introduction to Kazakh language, alphabet, and greetings.",
  "goal": "Learn basic phrases and pronunciation.",
  "difficulty": "beginner",
  "unlock_requirements": {
    "previous_module": "intro-001",
    "min_xp": 50.0
  },
  "reward": {
    "xp": 100.0,
    "badge": "kazakh_starter"
  }
}
```

---

### `DELETE /api/modules/:code`

Удалить модуль по его коду.

---

### `POST /api/modules/:code/lessons-list/:lessonCode`

Добавить урок с кодом `lessonCode` в модуль `code`.

---

### `DELETE /api/modules/:code/lessons-list/:lessonCode`

Удалить урок с кодом `lessonCode` из модуля `code`.

---


---

## 📚 Lessons (Admin Only)

### `GET /api/lessons/:code`

Получить информацию об уроке по его коду.

---

### `POST /api/lessons`

Создать новый урок.
**Request Body (JSON):**

```json
{
  "code": "lesson_012",
  "title": "Основы грамматики",
  "description": "Изучение базовых грамматических правил",
  "created_at": "2025-03-02T12:00:00Z",
  "deleted_at": null
}
```

---

### `PATCH /api/lessons/:code`

Обновить данные урока (можно частично).
**Request Body (JSON):**

```json
{
  "title": "new title",
  "description": "desc"
}
```

---

### `DELETE /api/lessons/:code`

Удалить урок по его коду.

---

### `POST /api/lessons/:code/exercises-list/:exerciseCode`

Добавить упражнение с кодом `exerciseCode` в урок `code`.

---

### `DELETE /api/lessons/:code/exercises-list/:exerciseCode`

Удалить упражнение `exerciseCode` из урока `code`.

---

Вот красиво оформленный раздел **Exercises (admin only)** в том же стиле, что и ранее:

---

## 🧩 Exercises (Admin Only)

### `GET /api/exercises/:code`

Получить упражнение по его коду.

---

### `POST /api/exercises`

Создать новое упражнение. В зависимости от `type`, структура тела может меняться.

#### 📌 Тип: `match_pairs`

**Request Body (JSON):**

```json
{
  "code": "match_pairs_5",
  "type": "match_pairs",
  "question": "Соедини слова с их переводами. 4 with delete",
  "pairs": [
    { "term": "Салам", "match": "Привет" },
    { "term": "Рақмет", "match": "Спасибо" },
    { "term": "fdwfw", "match": "f489" }
  ],
  "explanation": "Сопоставь казахские слова с их русскими значениями. 2",
  "hints": [
    "'Салам' похоже на арабское приветствие.",
    "fdafdsafsadfasd"
  ],
  "correct_order": ["fdafd", "fdsfas"]
}
```

---

#### 📌 Тип: `order_words`

**Request Body (JSON):**

```json
{
  "code": "order_words_73",
  "type": "order_words",
  "question": "Состаравильное предложение: Менің атым Алихан",
  "options": [
    "Мен",
    "Алихан",
    "зовут",
    "меня",
    "привет",
    "как",
    "дела"
  ],
  "correct_order": [
    "меня",
    "зовут",
    "Алихан"
  ],
  "correct_answer": "fdafa",
  "explanation": "В казахском языке порядок слов отличается от русского.",
  "hints": [
    "Сначала подлежащее, затем сказуемое."
  ]
}
```

---

#### 📌 Тип: `manual_typing`

**Request Body (JSON):**

```json
{
  "code": "translate1_1",
  "type": "manual_typing",
  "question": "Переведи на казахский: 'Спасибо'",
  "correct_answer": "Рақмет",
  "explanation": "",
  "hints": [
    "Начинается на 'Р'.",
    "Есть похожее слово в тюркских языках."
  ]
}
```

---

### `PATCH /api/exercises/:code`

Обновить упражнение (можно частично). Нужно чтобы в теле запроса были все поля которые из всех возможных относятся именно к определенному типу упражнения

**Request Body (JSON):**

```json
{
  "type": "match_pairs",
  "correct_answer": "niga",
  "pairs": [
    { "term": "san", "match": "bazan" }
  ]
}
```

---

### `DELETE /api/exercises/:code`

Удалить упражнение по его коду.

---

Вот красиво оформленный раздел **Friends** в том же стиле:

---

## 👥 Friends

### `POST /api/friends/send-request`

Отправить запрос на добавление в друзья.
**Request Body (JSON):**

```json
{
  "requester_username": "s4ab",
  "recipient_username": "asan2"
}
```

---

### `POST /api/friends/handle-request`

Обработать входящий запрос в друзья (принять или отклонить).
**Request Body (JSON):**

```json
{
  "requester_username": "asan2",
  "recipient_username": "s4ab",
  "status": "declined"  // или "accepted"
}
```

---

### `GET /api/friends/friend-list?username=seab`

Получить список друзей пользователя.

---

### `GET /api/friends/request-list?username=seab`

Получить список входящих запросов в друзья.

---

Вот красиво оформленный раздел **Data** для обычных пользователей:

---

## 📊 Data (Для обычных пользователей)

### `GET /api/data/modules`

Получить список всех доступных модулей для прохождения.

---

### `GET /api/data/users?username=seab&withProgress=true`

Получить информацию о пользователе `seab`.
Если параметр `withProgress=true`, то дополнительно возвращается прогресс пользователя по модулям и урокам.

---













