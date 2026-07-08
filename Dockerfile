# 1. Сборка фронтенда
FROM node:20-alpine AS frontend-build

WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm install

COPY frontend/ ./
RUN npm run build


# 2. Сборка Go backend
FROM golang:1.23-alpine AS backend-build

WORKDIR /app

RUN apk update && apk add --no-cache git curl

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY --from=frontend-build /app/web ./web

RUN go build -o /app/main .


# 3. Финальный контейнер
FROM golang:1.23-alpine

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=backend-build /app/main /app/main
COPY --from=backend-build /app/web /app/web
COPY --from=backend-build /app/configs /app/configs

EXPOSE 8080

CMD ["/app/main"]