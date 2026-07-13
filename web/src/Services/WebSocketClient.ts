type Callback = (data: any) => void;

class EventEmitter {
    private listeners: Record<string, Callback[]> = {};

    public on(eventName: string, callback: Callback): void {
        if (!this.listeners[eventName]) {
            this.listeners[eventName] = [];
        }
        this.listeners[eventName].push(callback);
    }

    public emit(eventName: string, data: any): void {
        const eventCallbacks = this.listeners[eventName];
        
        if (eventCallbacks) {
            for (const callback of eventCallbacks) {
                callback(data);
            }
        }
    }
}

export class WebSocketClient extends EventEmitter {
    private socket: WebSocket | null = null;

    private isCorrectJSON(message: any): boolean {
        return (message && typeof message === 'object' && 'cmd' in message && 'payload' in message);
    }

    public connect(url: string) {
        this.socket = new WebSocket(url);

        this.socket.onmessage = (event: MessageEvent) => {
            try {
                const message = JSON.parse(event.data);
                if (this.isCorrectJSON(message)) {
                    this.emit(message.cmd, message.payload);
                } else {
                    console.warn("Получено некорректное сообщение с сервера", message);
                }
                

            } catch(error) {
                console.error("Ошибка при прасинге JSON: ", error)
            }
           
        };
    }

    public send(cmd: string, payload: any) {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify({ cmd, payload }));
        }
    }

     
}