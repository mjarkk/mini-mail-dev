import { getUrl } from "./fetch"

export class EmailEventsWebsocket {
	private closed = false
	private socket: WebSocket | null = null
	private newEmail: () => void
	private connectionChange: (connected: boolean) => void

	constructor(
		newEmail: () => void,
		connectionChange: (connected: boolean) => void,
	) {
		this.newEmail = newEmail
		this.connectionChange = connectionChange
	}

	async start() {
		if (this.socket) {
			throw "cannot start a Websocket twice"
		}

		let first = true
		while (true) {
			await new Promise((res) => {
				this.socket = new WebSocket(getUrl("/api/emails-events", true))
				this.socket.onopen = () => {
					this.connectionChange(true)
					// Emails that arrived while disconnected were never announced.
					if (!this.closed && !first) {
						this.newEmail()
					}
				}
				this.socket.onmessage = (ev) => {
					if (!this.closed && ev.data === "new-email") {
						this.newEmail()
					}
				}
				this.socket.onclose = () => {
					this.connectionChange(false)
					setTimeout(res, 5_000)
				}
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
