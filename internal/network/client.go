package network

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn   *websocket.Conn
	ID     string
	msg    chan []byte
	done   chan struct{}
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
	state  State
	room   *GameRoom
	closed bool
}

// Генерация ID
func generateID() string {
	return uuid.New().String()
}

// Создаёт нового клиента
func NewClient(ws *websocket.Conn) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		conn:   ws,
		ID:     generateID(),
		msg:    make(chan []byte, 100),
		done:   make(chan struct{}),
		ctx:    ctx,
		cancel: cancel,
		state:  MainMenuState(),
		room:   nil,
		closed: false,
	}
}

func (c *Client) SetState(s State) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.state = s
}

// Закрывает клиента
func (c *Client) Close() error {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return nil
	}
	c.closed = true

	msgCh := c.msg
	doneCh := c.done
	conn := c.conn

	c.msg = nil
	c.done = nil
	c.conn = nil
	c.mu.Unlock()

	c.cancel()

	if msgCh != nil {
		close(msgCh)
	}
	if doneCh != nil {
		close(doneCh)
	}
	if conn != nil {
		return conn.Close()
	}
	return nil
}

// Читает сообщения с клиента
func (c *Client) ReadMessages() {
	log.Printf("Начал слушать клиента %v\n", c.ID)
	defer c.Close()
	for {
		c.mu.Lock()
		if c.closed || c.done == nil || c.conn == nil || c.msg == nil {
			c.mu.Unlock()
			return
		}
		done := c.done
		conn := c.conn
		c.mu.Unlock()

		select {
		case <-done:
			return
		default:
			//чтение сообщения
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("Ошибка чтения:", err)
				return
			}

			//превращение данных из json в структуру команды
			var req struct {
				Cmd     string          `json:"cmd"`
				Payload json.RawMessage `json:"payload"`
			}
			if err := json.Unmarshal(msg, &req); err != nil {
				c.SendError(err)
				continue
			}

			//обработка команды
			log.Printf("Принял команду %+v\n от клиента %v", req, c.ID)
			if err := c.state.HandleCommand(c, req.Cmd, req.Payload); err != nil {
				c.SendError(err)
			}
		}
	}
}

// Рассылает сообщения от сервера к клиенту
func (c *Client) WriteMessages() {
	defer c.Close()
	for {
		c.mu.Lock()
		if c.closed || c.done == nil || c.conn == nil || c.msg == nil {
			c.mu.Unlock()
			return
		}
		done := c.done
		conn := c.conn
		msgCh := c.msg
		c.mu.Unlock()

		select {
		case <-done:
			return
		case msg := <-msgCh:
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("Ошибка записи:", err)
				return
			}
		}
	}
}

// Передаёт сообщение клиенту
func (c *Client) Send(msg []byte) error {
	c.mu.Lock()
	if c.closed || c.msg == nil {
		c.mu.Unlock()
		return errors.New("клиент закрыт")
	}
	msgCh := c.msg
	c.mu.Unlock()

	select {
	case msgCh <- msg:
		return nil
	default:
		return errors.New("канал сообщений переполнен")
	}
}

// Передаёт данные клиенту
func (c *Client) SendMessage(typeMsg string, data interface{}) error {
	c.mu.Lock()

	if c.conn == nil {
		c.mu.Unlock()
		return errors.New("нет соединения")
	}

	c.mu.Unlock()

	message := struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}{
		Type: typeMsg,
		Data: data,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("ошибка создания json: %w", err)
	}

	if err = c.Send(jsonData); err != nil {
		return fmt.Errorf("ошибка отправки: %w", err)
	}
	fmt.Printf("Сервер отправил %+v клиенту %v\n", message, c.ID)
	return nil
}

// передаёт ошибку клиенту
func (c *Client) SendError(err error) error {
	return c.SendMessage("error", map[string]interface{}{
		"message": err.Error(),
	})
}

// передаёт клиенту текстовое сообщение
func (c *Client) SendText(TypeMsg string, text string) error {
	return c.SendMessage(TypeMsg, map[string]interface{}{
		"message": text,
	})
}
