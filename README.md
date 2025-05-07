# Biathlon Event Processor

Утилита командной строки для обработки логов биатлонных соревнований.
Она читает JSON-конфигурацию и файлы событий, реконструирует ход гонки для каждого участника (круги, стрельбу, штрафы, финиш) и выводит:

- Хронологический лог
- Итоговый отчёт по каждому спортсмену

## Требования
- Go 1.21+
- GNU Make (или аналог на Windows)
- Git (для клонирования)

## Makefile
Все основные задачи автоматизированы через Makefile:

```bash
make	fmt, vet, test, build
make fmt	авто-форматирование (go fmt ./...)
make vet	статический анализ (go vet ./...)
make lint	запуск golangci-lint (при установленном)
make test	запуск всех тестов (go test ./...)
make build	сборка бинарника biathlon
make run	сборка и запуск бинарника (читает config.json и events/)
make clean	удаление бинарника
```

Сборка и запуск
Клонировать репозиторий:

```bash
git clone https://github.com/yourorg/YadroTest.git
cd YadroTest
```
Собрать и запустить:

```bash
make run
```
Утилита автоматически найдёт config.json и папку events/ в текущей директории.
Если нужна только сборка:

```bash
make build
```
Результат: файл biathlon (или biathlon.exe на Windows).
