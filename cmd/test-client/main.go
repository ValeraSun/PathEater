package main

import (
    "flag"
    "log"
    "time"
    "github.com/gorilla/websocket"
)

func main() {
    var (
        url      = flag.String("url", "ws://localhost:8080/ws", "WebSocket URL")
        count    = flag.Int("count", 5, "Количество сообщений")
        delay    = flag.Duration("delay", 2*time.Second, "Задержка между сообщениями")
        jsonMode = flag.Bool("json", false, "Отправлять JSON сообщения")
    )
    flag.Parse()

    log.Printf("🔌 Подключение к %s", *url)
    
    ws, _, err := websocket.DefaultDialer.Dial(*url, nil)
    if err != nil {
        log.Fatal("❌ Ошибка подключения:", err)
    }
    defer ws.Close()

    log.Println("✅ Подключено успешно!")

    // Городина для чтения ответов
    go func() {
        for {
            _, msg, err := ws.ReadMessage()
            if err != nil {
                log.Println("⚠️ Ошибка чтения:", err)
                return
            }
            log.Printf("📥 Получено: %s", string(msg))
        }
    }()

    // Отправка сообщений
    for i := 1; i <= *count; i++ {
        var message []byte
        
        if *jsonMode {
            // JSON сообщение
            jsonMsg := map[string]interface{}{
                "type":    "test",
                "id":      i,
                "message": "Test message #" + string(rune('0'+i)),
                "timestamp": time.Now().Unix(),
            }
            jsonMsgBytes, _ := json.Marshal(jsonMsg)
            message = jsonMsgBytes
        } else {
            message = []byte("Test message #" + string(rune('0'+i)))
        }

        log.Printf("📤 Отправка [%d/%d]: %s", i, *count, string(message))
        
        err := ws.WriteMessage(websocket.TextMessage, message)
        if err != nil {
            log.Println("❌ Ошибка отправки:", err)
            break
        }

        time.Sleep(*delay)
    }

    log.Println("🏁 Тестирование завершено")
}