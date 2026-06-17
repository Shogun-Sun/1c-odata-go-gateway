# Этап сборки
FROM alt:p11 AS builder

# Устанавливаем Go
RUN apt-get update && apt-get install -y golang && apt-get clean

WORKDIR /app

# Копируем файлы зависимостей и скачиваем их (для кэширования слоев)
COPY go.mod ./
RUN go mod download

# Копируем код
COPY . .

# Собираем бинарник. 
# CGO_ENABLED=0 делает бинарник статическим (не зависит от системных библиотек)
RUN CGO_ENABLED=0 GOOS=linux go build -o /academic-api ./cmd/api/main.go

# Финальный образ
FROM alt:p11

# Копируем только готовый файл из этапа builder
COPY --from=builder /academic-api /academic-api

# Указываем порт, который слушает твое приложение
EXPOSE 4000

# Запуск приложения
ENTRYPOINT ["/academic-api"]