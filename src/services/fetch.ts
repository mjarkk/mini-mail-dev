export function getUrl(path: string, isWebsocket = false) {
	if (import.meta.env.DEV) {
		return `${isWebsocket ? "ws" : "http"}://localhost:3000${path}`
	}

	if (isWebsocket) {
		const protocol = location.protocol === "https:" ? "wss://" : "ws://"
		return protocol + location.host + path
	}

	return path
}

export function fetch(path: string, options?: RequestInit) {
	return window.fetch(getUrl(path), options)
}
