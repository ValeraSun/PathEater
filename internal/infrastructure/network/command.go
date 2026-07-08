package network

import (
    "context"
    "encoding/json"
    "errors"
    "ws/internal/ws"
)

type Command interface {
    Name() string
    Execute(ctx context.Context, client *ws.Client, payload json.RawMessage) error
}