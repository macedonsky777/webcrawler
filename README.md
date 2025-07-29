# Web Crawler 🕷️

Швидкий та ефективний веб-краулер, написаний на Go, який здатен паралельно сканувати веб-сайти та генерувати детальні звіти про внутрішні посилання.

## ✨ Особливості

- **🚀 Конкурентність**: Використовує goroutines для паралельного сканування сторінок
- **🎯 Контроль навантаження**: Налаштування максимальної кількості одночасних запитів
- **📊 Розумні обмеження**: Встановлення максимальної кількості сторінок для сканування
- **🔒 Thread-safe**: Безпечна робота з даними в багатопоточному середовищі
- **📈 Відсортовані звіти**: Красиві звіти з сортуванням за популярністю та алфавітом
- **🌐 Домен-специфічний**: Сканує тільки сторінки в межах одного домену
- **⚡ CLI інтерфейс**: Прості та гнучкі налаштування через командний рядок

## 🛠️ Встановлення

### Передумови
- Go 1.21+ 
- Доступ до інтернету

### Клонування репозиторію
```bash
git clone https://github.com/yourusername/web-crawler-bootdev.git
cd web-crawler-bootdev
```

### Встановлення залежностей
```bash
go mod tidy
```

### Збірка
```bash
go build -o crawler
```

## 🚀 Використання

### Базовий синтаксис
```bash
./crawler <URL> <maxConcurrency> <maxPages>
```

### Параметри
- **URL**: Стартова URL для сканування (обов'язковий)
- **maxConcurrency**: Максимальна кількість одночасних HTTP запитів (1-20)
- **maxPages**: Максимальна кількість сторінок для сканування

### Приклади використання

#### Базове сканування
```bash
# Сканувати максимум 10 сторінок з 3 одночасними запитами
./crawler "https://example.com" 3 10
```

#### Швидке сканування
```bash
# Високошвидкісне сканування з 10 одночасними запитами
./crawler "https://blog.example.com" 10 50
```

#### Обережне сканування
```bash
# Повільне сканування для невеликих серверів
./crawler "https://small-site.com" 1 5
```

### Альтернативний запуск
```bash
# Без збірки (прямий запуск)
go run . "https://example.com" 5 20
```

## 📋 Приклад виводу

```
starting crawl of: https://example.com
crawling: https://example.com
crawling: https://example.com/about
crawling: https://example.com/contact
crawling: https://example.com/blog
crawling: https://example.com/products

=============================
  REPORT for https://example.com
=============================
Found 45 internal links to https://example.com
Found 12 internal links to https://example.com/about
Found 8 internal links to https://example.com/blog
Found 3 internal links to https://example.com/contact
Found 2 internal links to https://example.com/products
```

## 🏗️ Архітектура

### Ключові компоненти

#### Структура `config`
```go
type config struct {
    pages              map[string]int    // Зберігає кількість посилань на кожну сторінку
    baseURL            *url.URL          // Базовий URL домену
    mu                 *sync.Mutex       // Мьютекс для thread-safe доступу
    concurrencyControl chan struct{}     // Канал для контролю concurrency
    wg                 *sync.WaitGroup   // WaitGroup для синхронізації goroutines
    maxPages           int               // Максимальна кількість сторінок
}
```

#### Основні функції
- `crawlPage()`: Рекурсивне сканування сторінок з goroutines
- `normalizeURL()`: Нормалізація URL для уникнення дублікатів
- `getHTML()`: HTTP запити з перевіркою content-type
- `getURLsFromHTML()`: Парсинг HTML та витягування посилань
- `printReport()`: Генерація відсортованих звітів

## 🔧 Технічні деталі

### Конкурентність
Краулер використовує буферизований канал для контролю кількості одночасних HTTP запитів:
```go
concurrencyControl: make(chan struct{}, maxConcurrency)
```

### Безпека потоків
Всі операції з спільною мапою `pages` захищені мьютексом для запобігання race conditions.

### Нормалізація URL
- Видалення trailing slashes
- Приведення до нижнього регістру хосту
- Очищення шляхів

### Обмеження
- Сканує тільки HTML сторінки (ігнорує XML, JSON тощо)
- Працює тільки в межах одного домену
- Респектує налаштування maxPages та maxConcurrency

## 📁 Структура проекту

```
web-crawler-bootdev/
├── main.go                 # Точка входу та CLI logic
├── normalize_url.go        # Основна логіка краулера
├── normalize_url_test.go   # Unit тести
├── go.mod                  # Go модуль та залежності
├── go.sum                  # Checksums залежностей
└── README.md              # Цей файл
```

## 🧪 Тестування

### Запуск unit тестів
```bash
go test -v
```

### Запуск конкретного тесту
```bash
go test -v -run TestNormalizeURL
```

## ⚠️ Обмеження та рекомендації

### Обмеження
- Не підтримує JavaScript-рендеринг
- Не обробляє редиректи автоматично
- Обмежений одним доменом

### Рекомендації
- Починайте з невеликих значень concurrency (1-3)
- Використовуйте розумні обмеження maxPages
- Будьте обережні з великими сайтами
- Респектуйте robots.txt (функціонал не реалізований)

## 🤝 Розробка

Цей проект створений як частина курсу [Boot.dev](https://boot.dev) для вивчення:
- Go programming language
- Конкурентності та goroutines
- HTTP клієнтів
- HTML парсингу
- CLI інструментів

## 📄 Ліцензія

MIT License - деталі в файлі LICENSE

## 🙏 Подяки

- [Boot.dev](https://boot.dev) за чудовий курс
- Go команда за потужну мову програмування
- Golang.org/x/net за HTML парсинг