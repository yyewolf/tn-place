import { putUint32 } from "./utils";

const WS_URL = import.meta.env.VITE_WS_URL;

export class Gateway {
    #socket;
    #listeners = new Map();

    constructor(listeners) {
        this.#socket = null;
        listeners = listeners || [];
        for (let i = 0; i < listeners.length; i++) {
            const [name, listener] = listeners[i];
            this.#listeners.set(name, listener);
        }
    }

    initConnection() {
        this.#connect(WS_URL);
    }

    #connect(url) {
		this.#socket = new WebSocket(url);

		const socketMessage = async (event) => {
			let b = await event.data.arrayBuffer();
            for (const listener of this.#listeners.values()) {
                listener(b);
            }
		};

		const socketClose = (event) => {
			this.#socket = null;
		};

		const socketError = (event) => {
			console.error("Error making WebSocket connection.");
			alert("Failed to connect.");
			this.#socket.close();
			window.location.href = "/";
		};

		this.#socket.addEventListener("message", socketMessage);
		this.#socket.addEventListener("close", socketClose);
		this.#socket.addEventListener("error", socketError);
	}

    addListener(name, listener) {
        this.#listeners.set(name, listener);
    }

    removeListener(name) {
        this.#listeners.delete(name);
    }

    setPixel(x, y, color) {
		if (this.#socket != null && this.#socket.readyState == 1) {
			let b = new Uint8Array(11);
			putUint32(b.buffer, 0, x);
			putUint32(b.buffer, 4, y);
			for (let i = 0; i < 3; i++) {
				b[8 + i] = color[i];
			}
			this.#socket.send(b);
		} else {
			alert("Disconnected.");
			console.error("Disconnected.");
		}
    }
}