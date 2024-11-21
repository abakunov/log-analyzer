# Log Analyzer

## Описание проекта

**Log Analyzer** – это приложение для анализа логов веб-серверов, предоставляющее удобный способ получения метрик и генерации отчётов. Оно позволяет анализировать логи с учётом фильтров, вычислять метрики, такие как RPS (Requests Per Second) и количество уникальных IP-адресов.

### Метрики:
- `FileNames`: Список обработанных файлов.
- `StartDate` и `EndDate`: Диапазон времени логов.
- `TotalRequests`: Общее количество запросов.
- `TotalRespSize`: Общий размер ответов.
- `AverageRespSize`: Средний размер ответа.
- `Percentile95`: 95-й перцентиль размера ответа.
- `Resources`: Частота запросов на ресурсы.
- `StatusCodes`: Частота кодов ответов.
- `UniqueIPs`: Количество уникальных IP-адресов (**дополнительные баллы**).
- `RPS`: Количество запросов в секунду (**дополнительные баллы**).

### Фильтрация логов (дополнительные баллы):
- `ip`: Фильтрация по IP-адресу.
- `timestamp`: Фильтрация по временным меткам.
- `method`: HTTP-метод (GET, POST и т.д.).
- `url`: URL-адрес запроса.
- `protocol`: Протокол запроса (HTTP/1.1 и т.д.).
- `status`: Код ответа HTTP.
- `response_size`: Размер ответа в байтах.
- `referer`: URL реферера.
- `agent`: User-Agent клиента.

**Особенности фильтрации**:
- Указание `*` в конце значения ищет совпадения по началу строки.
- Без `*` происходит точное сравнение.

---

## Установка и запуск

**Установка через `install.sh`**:
``` bash
sh install.sh
```
После установки программа доступна как analyzer.

**Запуск**:
``` bash
analyzer --path "logs/*" 
```

Параметры:

- `path`: Путь(и) к лог-файлам или паттерн (обязательный). Вводится в двойных кавычках. 
- `from`: Начальная дата (опционально).
- `to`: Конечная дата (опционально).
- `format`: Формат отчёта (markdown, adoc). Если не указан, выводится в консоль.
- `filter-field`: Поле для фильтрации (опционально). 
- `filter-value`: Значение для фильтрации (опционально). Вводится в двойных кавычках.

Приложение поддерживает фильтрацию логов по указанным полям. 
Значение для фильтрации может быть точным или содержать символ `*` в конце для поиска по началу строки. 
Если `*` отсутствует, производится поиск по точному совпадению.

**Приложение можно запустить без установки**

1. Сбилдить проект:
``` bash
go build -o analyzer cmd/run/main.go
```

2. Запустить:
``` bash
./analyzer --path "logs/*"
```

### Пример отчета

#### Report created: 21.11.2024 18:27:47

#### General Information

| Files | https://raw.githubusercontent.com/elastic/examples/master/Common%20Data%20Formats/nginx_logs/nginx_logs |
| --- | --- |
| Start Date | 17.05.2015 |
| End Date | 04.06.2015 |
| Total Requests | 51462 |
| Unique IPs Count | 2660 |
| RPS (Requests/sec) | 0.03 |
| Average Response Size | 659510b |
| 95th Percentile Size | 1768b |

#### Requested Resources

| Resource | Count |
| --- | --- |
| /downloads/product_1 | 30285 |
| /downloads/product_2 | 21104 |
| /downloads/product_3 | 73 |

#### Response Codes

| Code | Count |
| --- | --- |
| 404 | 33876 |
| 304 | 13330 |
| 200 | 4028 |
| 206 | 186 |
| 403 | 38 |
| 416 | 4 |



## Структура проекта

``` bash
/log-analyzer
├── cmd
│   └── run                 
│       ├── main.go          # Точка входа, настройка CLI и запуск анализа
│       ├── main_test.go     # Тесты для CLI
├── internal
│   ├── application          # Логика приложения
│   │   ├── log_analyzer.go  # Основной анализ логов
│   │   ├── log_parser.go    # Парсинг строк логов
│   │   ├── log_processor.go # Логика фильтрации
│   │   ├── metrics_updater.go # Обновление метрик
│   │   ├── file_processor.go  # Чтение логов из файлов и URL
│   │   ├── filter.go         # Проверка фильтров
│   │   ├── utils.go          # Вспомогательные функции
│   │   ├── *_test.go         # Тесты для всех модулей приложения
│   ├── domain               # Модели данных
│   │   └── models.go        # Структуры LogRecord и Metrics
│   ├── infrastructure       # Вспомогательные модули
│       ├── input_parsers.go # Парсинг времени и путей
│       ├── report_formatter.go # Форматирование отчётов
│       ├── report_output.go # Вывод отчётов
│       ├── *_test.go        # Тесты для инфраструктурных модулей
```

### Тестирование

Для запуска тестов используйте команду:
``` bash
make test
```
Тесты охватывают:
- Проверку логики фильтрации, обработки файлов и метрик.
- Корректность работы CLI.
