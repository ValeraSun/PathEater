type Callback = (data: any) => void;

class EventEmitter {
    private listeners: Record<string, Callback[]> = {};

    public on(eventName: string, callback: Callback): void {
        if (!this.listeners[eventName]) {
            this.listeners[eventName] = [];
        }
        this.listeners[eventName].push(callback);
    }

    public off(eventName: string, callback: Callback): void {
        const callbacks = this.listeners[eventName];
        if (!callbacks) return;

        this.listeners[eventName] = callbacks.filter(
            cb => cb !== callback
        );
    }

    public once(eventName: string, callback: Callback): void {
        const wrapper = (data: any) => {
            this.off(eventName, wrapper);
            callback(data);
        };
        this.on(eventName, wrapper);
    }

    public emit(eventName: string, data: any): void {
        const eventCallbacks = this.listeners[eventName];

        if (eventCallbacks) {
            for (const callback of [...eventCallbacks]) {
                callback(data);
            }
        }
    }
}

export class WebSocketClient extends EventEmitter {
    private socket: WebSocket | null = null;

    public get isConnected(): boolean {
        return this.socket?.readyState === WebSocket.OPEN;
    }

    private isCorrectJSON(message: any): boolean {
        return (
            message &&
            typeof message === "object" &&
            "type" in message &&
            "data" in message
        );
    }

    public connect(url: string): Promise<void> {
        return new Promise((resolve, reject) => {
            this.socket = new WebSocket(url);

            this.socket.onopen = () => {
                this.emit("connected", null);
                resolve();
            };

            this.socket.onerror = () => {
                this.emit("connectionError", null);
                reject(new Error("Не удалось подключиться к серверу"));
            };

            this.socket.onclose = () => {
                this.emit("disconnected", null);
            };

           this.socket.onmessage = (event: MessageEvent) => {
                try {
                    const message = JSON.parse(event.data);

                    if (this.isCorrectJSON(message)) {
                        let processedData = message.data;

                        if (typeof message.data === "string") {
                            try {
                                const decoded = atob(message.data);
                                processedData = JSON.parse(decoded); 
                            } catch (e) {
                            }
                        }

                        this.emit(message.type, processedData);
                    } else {
                        console.warn("Получено некорректное сообщение с сервера", message);
                    }
                } catch (error) {
                    console.error("Ошибка при парсинге JSON: ", error);
                }
            };
        });
    }

    public send(cmd: string, payload: any): void {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(JSON.stringify({ cmd, payload }));
        } else {
            console.warn(`Попытка отправить "${cmd}" без активного соединения`);
        }
    }

    public disconnect(): void {
        this.socket?.close();
        this.socket = null;
    }
}