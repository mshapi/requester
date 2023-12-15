# Тестовое задание

Необходимо написать программу(сервис), реализующую функционал отправки запросов
на сторонний сервис(реализация webhook).

Пример конфигурации

```yaml
url: https://some.domain
requests:
    amount: 1000
    per_second: 10
```

Отправляем **POST** запросом на url запросы с **body** = `{ "iteration": ${number} }`,
где `number` порядковый номер итерации из `requests.amount`.
Реализация должна учитывать количество отправляемых запросов `request.per_second`.
Написать тест для предложенной реализации.

Используемый инструментарий
- go 1.21.x
- https://github.com/uber-go/zap
- https://github.com/spf13/cobra
- https://github.com/uber-go/fx
- https://github.com/go-resty/resty