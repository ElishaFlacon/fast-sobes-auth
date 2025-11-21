# Gosling

CLI-утилита для автоматической генерации компонентов микросервисной архитектуры (handler, usecase, repository).

## Использование

### Создание сервиса

# Дефолтный сервис
`make gosling-make SERVICE=name`

# С кастомными методами
`make gosling-make SERVICE=name METHODS="method1 method2 method3"`

# С пропуском слоев
`make gosling-make SERVICE=name NO_HANDLER=true NO_REPOSITORY=true`

### Удаление сервиса

# Полное удаление
`make gosling-remove SERVICE=name`

# С пропуском слоев
`make gosling-remove SERVICE=name NO_USECASE=true`

## Флаги

- `SERVICE` - имя сервиса (обязательно)
- `METHODS` - список методов через пробел
- `NO_HANDLER` - пропустить handler
- `NO_USECASE` - пропустить usecase
- `NO_REPOSITORY` - пропустить repository
- `NO_APP` - пропустить обновление provider

## После генерации

Зарегистрируйте handler в `internal/app/app.go`:
