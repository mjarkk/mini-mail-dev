export class EmailEventsWebsocket {
	private closed = false
	private socket: WebSocket | null = null
	private newEmail: () => void

	constructor(newEmail: () => void) {
		this.newEmail = newEmail
	}

	async start() {
		if (this.socket) {
			throw "cannot start a Websocket twice"
		}

		let first = true
		while (true) {
			await new Promise((res) => {
				this.socket = new WebSocket("ws://localhost:3000/api/emails-events")
				this.socket.onopen = () => {
					if (this.closed && !first) {
						this.newEmail()
					}
				}
				this.socket.onmessage = (ev) => {
					if (!this.closed && ev.data === "new-email") {
						this.newEmail()
					}
				}
				this.socket.onclose = () => setTimeout(res, 5_000)
			})
			if (this.closed) {
				break
			}
			first = false
		}
	}

	close() {
		this.closed = true
		this.socket?.close()
	}
}
