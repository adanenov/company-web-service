==============================
Демонстрация проекта: Company Web Service
==============================

Цель: показать работу контейнеризированного веб-сервиса с автоматическим перезапуском, Nginx, CI/CD и управлением конфигурацией через .env.

---

1️⃣ Подготовка проекта

1. Перейти в папку проекта:
    cd ~/projects/company-web-service

2. Проверить файл .env:
    cat .env
    Пример:
        PORT=3000

3. Пересобрать контейнеры и поднять их:
    docker-compose down
    docker-compose build --no-cache
    docker-compose up -d

4. Проверить, что контейнеры запущены:
    docker ps
    Должны быть:
        - company-web
        - nginx-proxy

---

2️⃣ Проверка веб-сервиса

1. В браузере открыть:
    http://localhost        → "IT Company Website"
    http://localhost/about  → "About our company"
    http://localhost/api    → {"status": "ok"}

---

3️⃣ Демонстрация Nginx Reverse Proxy

- Nginx проксирует запросы на web-контейнер.
- Проверка: http://localhost → работает через Nginx
- Порты 80 и 443 открыты наружу.

---

4️⃣ Демонстрация автоматического рестарта контейнера

1. Поднять контейнеры:

    docker-compose up -d

2. Проверить, что сайт работает:

    http://localhost

3. Вызвать краш-контрол для демонстрации рестарта:

    curl http://localhost/crash

4. Через секунду проверить статус контейнера:

    docker ps

- Контейнер снова в статусе Up, uptime обновлён → restart policy сработала

5. Проверить логи контейнера:

    docker logs -f company-web

- В логах будет видно:

    Server running on port 3000
    Container crash triggered
    Server running on port 3000

> Иными словами: контейнер реально падает → Docker автоматически его поднимает → restart policy работает

5️⃣ CI/CD workflow (GitHub Actions)

1. Сделать небольшое изменение в коде, например main.go:
    fmt.Println("Test CI/CD")

2. Commit & push:
    git add .
    git commit -m "Test CI/CD"
    git push

3. На GitHub открыть вкладку Actions:
    - Workflow сработал:
        - Checkout
        - Build Docker image
        - Push (если настроено)

---

6️⃣ Управление конфигурацией через .env

1. Изменить порт в .env:
    PORT=4000

2. Перезапустить контейнеры:
    docker-compose up -d

3. Проверить в браузере:
    http://localhost:4000
    - Сайт работает без изменения кода

---

7️⃣ Проверка безопасности

- Открыты только порты 80 и 443 для Nginx
- Порт приложения (3000) лучше закрыть наружу

Проверить:
    docker ps

---

8️⃣ Логи контейнера

- Проверяем работу приложения и перезапуск после краша:
    docker logs -f company-web

---

9️⃣ Код Go с endpoint для краша (/crash)

main.go:

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

---

10️⃣ docker-compose.yml

version: "3.9"

services:
  web:
    build: .
    container_name: company-web
    restart: always
    ports:
      - "${PORT}:3000"
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

---

11️⃣ nginx.conf

server {
    listen 80;
    server_name localhost;

    location / {
        proxy_pass http://web:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

---

12️⃣ Dockerfile

FROM golang:1.25.7-alpine
WORKDIR /app
COPY . .
RUN go build -o server main.go
EXPOSE 3000
CMD ["./server"]

---

13️⃣ Итог демонстрации

- Все сервисы работают: Web + Nginx  
- Автоматический restart проверен через /crash  
- CI/CD workflow проверен  
- Конфигурация через .env работает  
- Логи и безопасность показаны  
- Можно показать изменения кода → push → автоматическая сборка