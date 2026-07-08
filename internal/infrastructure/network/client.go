package network

import (
    "context"
    "errors"
    "log"
    "time"
    "sync"
    "fmt"
    
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
    closed bool
}

func generateID() {
    id := uuid.New()
    fmt.Println(id.String())
}

//Создаёт нового клиента
func NewClient(ws *websocket.Conn) *Client{
    ctx, cancel := context.WithCancel(context.Background())
    return &Client{
        conn:   ws,
        ID:     generateID(), 
        msg:    make(chan []byte, 100),
        done:   make(chan struct{}),
        ctx:    ctx,
        cancel: cancel,
        state:  NewMenuState(),z
        closed: false,
    }
}

//Закрывает клиента
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

//Читает сообщения с клиента
func (c *Client) ReadMessages() {
    defer c.Close()
    for {
        select {
        case <-c.done:
            return
        default:
            c.mu.Lock()
            if c.conn == nil {
                c.mu.Unlock()
                return
            }
            conn := c.conn
            c.mu.Unlock()

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
                c.SendError("invalid_json")
                continue
            }
    
            //обработка команды
            if err := c.State.HandleCommand(c.ctx, c, req.Cmd, req.Payload); err != nil {
                c.SendError(err.Error())
            }
        }
    }
}

//Рассылает сообщения от сервера к клиентам
func (c *Client) WriteMessages() {
    defer c.Close()
    for {
        select {
        case <-c.done:
            return
        case msg := <-c.msg:
            c.mu.Lock()
            if c.conn == nil {
                c.mu.Unlock()
                return
            }
            conn := c.conn
            c.mu.Unlock()
            err := conn.WriteMessage(websocket.TextMessage, msg)
            if err != nil {
                log.Println("Ошибка записи:", err)
                return
            }
        }
    }
}

//Передаёт сообщение клиентам
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