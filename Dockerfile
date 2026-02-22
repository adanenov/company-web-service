# Используем официальный Go образ
FROM golang:1.25.7-alpine

# Создаём рабочую папку
WORKDIR /app

# Копируем весь проект
COPY . .

# Собираем бинарник
RUN go build -o server main.go

# Указываем порт
EXPOSE 3000

# Запускаем сервер
CMD ["./server"]