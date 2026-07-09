# 1. Сборка фронтенда
FROM node:20-alpine AS frontend-build

WORKDIR /app/web

COPY web/package*.json ./
RUN npm install

COPY web/ ./
RUN npm run build


# 2. Сборка Go backend
FROM golang:1.26-alpine AS backend-build

WORKDIR /app

RUN apk update && apk add --no-cache git curl

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Берём собранный Vite build из web/dist
RUN rm -rf ./web
COPY --from=frontend-build /app/web/dist ./web

RUN go build -o /app/main .


# 3. Финальный контейнер
FROM golang:1.26-alpine

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=backend-build /app/main /app/main
COPY --from=backend-build /app/web /app/web
COPY --from=backend-build /app/configs /app/configs

EXPOSE 8080

CMD ["/app/main"]