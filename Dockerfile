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

RUN apk update && apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN rm -rf ./web
COPY --from=frontend-build /app/web/dist ./web

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/main ./cmd/app


# 3. Финальный контейнер
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=backend-build /app/main /app/main
COPY --from=backend-build /app/web /app/web
COPY --from=backend-build /app/internal/config /app/internal/config

EXPOSE 8080

USER nobody

CMD ["/app/main"]

#для запуска
#docker build -t path-eater .
#docker run --rm -p 8080:8080 --name path-eater path-eater