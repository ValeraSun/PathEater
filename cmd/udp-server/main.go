package main

import (
    "fmt"
    "log"
    "net"
    "time"
)

func main() {
    addr, err := net.ResolveUDPAddr("udp", ":8081")
    if err != nil {
        log.Fatal("Ошибка резолва адреса:", err)
    }

    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        log.Fatal("Ошибка запуска UDP сервера:", err)
    }
    defer conn.Close()

    log.Println("UDP-сервер запущен на :8081")

    buffer := make([]byte, 1024)
    var messageCount int

    for {
        n, clientAddr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            log.Println("Ошибка чтения:", err)
            continue
        }

        messageCount++
        message := string(buffer[:n])
        timestamp := time.Now().Format("15:04:05.000")
        
        log.Printf("📨 [%d] %s от %s: %s", 
            messageCount, timestamp, clientAddr, message)

        // Отправляем подтверждение обратно
        response := fmt.Sprintf("OK: получено сообщение #%d", messageCount)
        _, err = conn.WriteToUDP([]byte(response), clientAddr)
        if err != nil {
            log.Println("Ошибка отправки ответа:", err)
        }
    }
}