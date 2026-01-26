# FCS Task Backend API (черновик)

Этот документ описывает минимальный REST API контракт, который нужен фронтенду.
Все эндпоинты возвращают JSON.

## Общие правила

- Базовый префикс: `/api`
- Заголовок авторизации: `Authorization: Bearer <token>`
- Формат ошибки:

```json
{
  "error": {
    "code": "...",
    "message": "..."
  }
}
```

## Пользователь

### GET `/api/me`

```json
{
  "username": "student",
  "initials": "ST",
  "role": "instance_admin"
}
```

## Курсы

### GET `/api/courses`

```json
[
  {
    "id": "algorithms",
    "name": "Algorithms 101",
    "status": "in_progress",
    "url": "/course/algorithms"
  }
]
```

### POST `/api/courses`

```json
{
  "name": "Advanced C++",
  "slug": "advanced-cpp",
  "status": "created",
  "startDate": "2024-10-01",
  "endDate": "2024-12-20",
  "repoTemplate": "git@gitlab.local/course-template.git",
  "description": "..."
}
```

### GET `/api/courses/:courseId`

```json
{
  "id": "algorithms",
  "name": "Algorithms 101",
  "status": "in_progress",
  "startDate": "2024-10-01",
  "endDate": "2024-12-20",
  "repoTemplate": "git@gitlab.local/course-template.git",
  "description": "..."
}
```

### PUT `/api/courses/:courseId`

Body same as POST `/api/courses`.

## Доска заданий

### GET `/api/courses/:courseId/board`

```json
{
  "courseName": "Algorithms 101",
  "courseStatus": "in_progress",
  "solvedScore": 126,
  "maxScore": 200,
  "solvedPercent": 63,
  "groups": [
    {
      "id": "week-1",
      "name": "Week 1: Warmup",
      "isSpecial": false,
      "startedAt": "2024-10-01T09:00:00Z",
      "endsAt": "2024-10-14T18:00:00Z",
      "deadlines": [
        {
          "id": "d1",
          "label": "Checkpoint",
          "percent": 0.6,
          "dueAt": "2024-10-10T18:00:00Z",
          "status": "active"
        }
      ],
      "tasks": [
        {
          "id": "t1",
          "name": "Arrays Sprint",
          "score": 20,
          "scoreEarned": 10,
          "stats": 0.64,
          "isBonus": false,
          "isSpecial": false,
          "url": "https://..."
        }
      ]
    }
  ]
}
```

## Все результаты

### GET `/api/courses/:courseId/scores`

```json
[
  { "id": 1, "student": "alex", "score": 192, "submitted": "2024-10-02" }
]
```

## Namespace

### GET `/api/namespaces`

```json
[
  {
    "id": "ns-01",
    "name": "Core CS",
    "slug": "core-cs",
    "description": "Foundational tracks",
    "gitlabGroupId": "22411",
    "coursesCount": 5,
    "usersCount": 180
  }
]
```

### GET `/api/namespaces/:namespaceId`

```json
{
  "namespace": {
    "id": "ns-01",
    "name": "Core CS",
    "slug": "core-cs",
    "description": "Foundational tracks",
    "gitlabGroupId": "22411"
  },
  "users": [
    { "id": "u-1", "username": "alex", "rmsId": "rms-210", "role": "namespace_admin" }
  ],
  "courses": [
    {
      "id": "c-101",
      "name": "Algorithms 101",
      "status": "running",
      "gitlabGroup": "algorithms-101",
      "owners": ["alex"],
      "url": "/course/algorithms"
    }
  ]
}
```

### POST `/api/namespaces/:namespaceId/users`

```json
{ "username": "new-user", "role": "student" }
```

### PUT `/api/namespaces/:namespaceId/users/:userId`

```json
{ "role": "program_manager" }
```

## Инстанс-админ

### GET `/api/instance/summary`

```json
{ "totalCourses": 21, "totalUsers": 920, "totalNamespaces": 6, "healthStatus": "ok" }
```

## Регистрация

### POST `/api/signup`

```json
{ "inviteCode": "...", "email": "...", "telegram": "...", "group": "..." }
```

### GET `/api/signup/status`

Опциональный эндпоинт для страницы завершения, если нужен.
