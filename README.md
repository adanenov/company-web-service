# Company Web Service

Проект: Контейнеризированный веб-сервис для демонстрации Docker, Nginx, CI/CD и автоматического перезапуска контейнера.

Описание:

Простой веб-сайт компании с маршрутами:

- / — главная страница "IT Company Website"
- /about — информация о компании
- /api — возвращает JSON {"status": "ok"}
- /crash — для демонстрации автоматического рестарта контейнера

Проект демонстрирует: контейнеризацию через Docker и docker-compose, Reverse Proxy через Nginx, автоматический restart контейнера (restart: always), управление конфигурацией через .env и CI/CD через GitHub Actions (сборка и push Docker-образа).

Стек технологий: Go (Golang), Docker + docker-compose, Nginx, GitHub Actions, .env

Быстрый старт:

1. Клонировать репозиторий:
   git clone https://github.com/adanenov/company-web-service.git
   cd company-web-service

2. Создать .env файл:
   echo "PORT=3000" > .env

3. Собрать и поднять контейнеры:
   docker-compose up -d

4. Проверить статус контейнеров:
   docker ps

Доступ к сервису:

- Главная страница: http://localhost
- О сайте: http://localhost/about
- API: http://localhost/api

Все запросы идут через Nginx Reverse Proxy. Открыты только порты 80 и 443.

Демонстрация автоматического рестарта:

1. Вызвать краш-контрол:
   curl http://localhost/crash

2. Проверить контейнер:
   docker ps
   - Контейнер автоматически поднимается → restart policy работает

3. Проверить логи контейнера:
   docker logs -f company-web
   В логах будет:
   Server running on port 3000
   Container crash triggered
   Server running on port 3000

Структура проекта:

.
├── main.go          # Go-сервер с /, /about, /api, /crash
├── Dockerfile       # Сборка контейнера Go
├── docker-compose.yml
├── nginx.conf       # Конфигурация Nginx Reverse Proxy
└── .env             # Порт веб-сервиса

CI/CD (GitHub Actions):

- Workflow: CI/CD Docker Build & Push
- Срабатывает при пуше в main
- Проверяет код, собирает Docker-образ и пушит в Docker Hub

Безопасность:

- Веб-контейнер недоступен напрямую снаружи
- Доступ только через Nginx (порты 80 и 443)
- Переменные окружения вынесены в .env

Логи и мониторинг:

- Просмотр логов веб-контейнера:
   docker logs -f company-web

- Просмотр логов Nginx:
   docker exec -it nginx-proxy sh
   cat /var/log/nginx/access.log

Итог:

- Контейнеры Web + Nginx работают
- Автоматический restart проверен
- Конфигурация через .env
- CI/CD workflow работает
- Все демонстрации можно показать преподавателю

Код Go с endpoint для краша (/crash):

package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func home(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "IT Company Website")
}

func about(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "About our company")
}

func api(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, `{"status": "ok"}`)
}

func crash(w http.ResponseWriter, r *http.Request) {
    log.Println("Container crash triggered")
    os.Exit(1)
}

func main() {
    http.HandleFunc("/", home)
    http.HandleFunc("/about", about)
    http.HandleFunc("/api", api)
    http.HandleFunc("/crash", crash)
    fmt.Println("Server running on port 3000")
    log.Fatal(http.ListenAndServe(":3000", nil))
}

docker-compose.yml:

version: "3.9"

services:
  web:
    build: .
    container_name: company-web
    restart: always
    environment:
      - PORT=${PORT}

  nginx:
    image: nginx:alpine
    container_name: nginx-proxy
    restart: always
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/conf.d/default.conf:ro

nginx.conf:

server {
    listen 80;
    server_name localhost;

    location / {
        proxy_pass http://web:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

Dockerfile:

FROM golang:1.25.7-alpine
WORKDIR /app
COPY . .
RUN go build -o server main.go
EXPOSE 3000
CMD ["./server"]



# CI/CD Конвейер для веб-сервиса «Company Web Service»
Общая идея

CI/CD (Continuous Integration / Continuous Deployment) позволяет автоматически собирать, тестировать и деплоить код при изменениях. В нашем проекте конвейер обеспечивает:

Автоматический билд Docker-образа при пуше в GitHub.

Проверку кода Go (go fmt, go vet) перед сборкой.

Автоматический push образа в Docker Hub.

Автоматический запуск контейнера на локальном сервере через docker-compose.

Авто-рестарт контейнера при падении (restart: always).

Таким образом, любой код, который ты пушишь в ветку main, сразу становится рабочим на сервере или локальной машине.

[VS Code / Локальная машина]
        │
        │ 1. Разработка кода
        │    - main.go, Dockerfile, nginx.conf, docker-compose.yml
        │
        ▼
[GitHub Repository]
        │
        │ 2. Push в ветку main
        │    - GitHub Actions срабатывает автоматически
        │
        ▼
[GitHub Actions Workflow]
        │
        │ 3a. Checkout кода
        │     - actions/checkout@v3
        │
        │ 3b. Проверка кода Go (опционально, но рекомендуется)
        │     - go fmt ./...
        │     - go vet ./...
        │
        │ 3c. Сборка Docker-образа
        │     - docker build -t company-web .
        │
        │ 3d. Push Docker-образа на Docker Hub
        │     - docker tag company-web <DOCKER_USERNAME>/company-web:latest
        │     - docker push <DOCKER_USERNAME>/company-web:latest
        │
        ▼
[Docker Hub]
        │
        │ 4. Хранение образа
        │    - Образ доступен для развёртывания на сервере
        │
        ▼
[Локальный сервер / Dev Server]
        │
        │ 5. Запуск контейнера через docker-compose
        │    - docker-compose pull (новый образ)
        │    - docker-compose up -d
        │
        │ 6. Nginx проксирует запросы к контейнеру web
        │    - Порты: 80 и 443 открыты наружу
        │
        │ 7. Авто-рестарт контейнера
        │    - restart: always
        │    - Если контейнер падает (например, /crash), Docker автоматически его поднимает
        │
        ▼
[Пользователь / браузер]
        │
        │ 8. Доступ к сервису
        │    - http://localhost
        │    - http://localhost/about
        │    - http://localhost/api

Подробные объяснения шагов

Разработка кода в VS Code / локальной машине

Файлы проекта: main.go, Dockerfile, docker-compose.yml, nginx.conf.

Код Go обрабатывает маршруты /, /about, /api, /crash.

.env содержит порт для web-контейнера.

Push в ветку main на GitHub

Любое изменение кода коммитится и пушится.

Это запускает workflow GitHub Actions.

GitHub Actions Workflow

Checkout: копирует текущую версию кода из репозитория в рабочую среду GitHub Actions.

Проверка кода: go fmt форматирует код, go vet проверяет на ошибки. Если есть проблемы — workflow останавливается.

Сборка Docker-образа: создаётся образ Go-приложения, готовый к запуску.

Push Docker-образа в Docker Hub: образ становится доступен для развёртывания на любой машине.

Docker Hub

Служит хранилищем готового Docker-образа.

Образ можно подтянуть на сервер или локально командой docker pull.

Развёртывание на локальном сервере / Dev сервере

docker-compose pull && docker-compose up -d подтягивает новый образ и запускает контейнеры.

Контейнеры web и nginx поднимаются и начинают работать.

Nginx

Прокси для web-контейнера.

Запросы к портам 80/443 направляются в контейнер web на порт 3000.

Закрыт прямой доступ к web-контейнеру.

Авто-рестарт контейнера

Настроено через restart: always в docker-compose.yml.

Контейнер перезапускается автоматически, если падает, например, при вызове /crash.

Доступ пользователю

Пользователь открывает браузер и получает контент через Nginx.

Все изменения кода автоматически попадают на сервер после workflow → демонстрация CI/CD.


 # UML Deployment Diagram

             +-----------------------+
             |      Пользователь      |
             |  (браузер / HTTP)     |
             +-----------+-----------+
                         |
                         v
             +-----------------------+
             |        Nginx          |
             |  Reverse Proxy        |
             |  Порты: 80, 443      |
             +-----------+-----------+
                         |
        -----------------+-----------------
        |                                 |
        v                                 v
+-------------------+           +-------------------+
|   Web Container   |           | (возможная БД)    |
|   Go-приложение   |           | (не реализовано)  |
|   Порт: 3000      |           +-------------------+
|   restart: always |
+-------------------+

Пользователь / браузер

Отправляет HTTP/HTTPS запросы на сервис.

Nginx

Принимает внешние запросы на порты 80/443.

Проксирует все запросы в web-контейнер.

Web Container (Go-приложение)

Слушает порт 3000.

restart policy: always → контейнер автоматически перезапускается при падении.

База данных (опционально)

В твоём проекте пока нет, но UML показывает возможное расширение системы.