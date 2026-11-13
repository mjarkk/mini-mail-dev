export interface Address {
	name: string
	address: string
}

export interface Attachment {
	filename: string
	contentType: string
}

export interface EmbeddedFile {
	cid: string
	contentType: string
}

/** An email client whose rendering the server can simulate. */
export interface EmailClient {
	id: string
	name: string
	/** Viewport width in CSS pixels of a mobile client. */
	width?: number
}

export interface EmailBase {
	id: string
	realDate: Date
	subject: string
	bodyHint: string
	bodyType: "html" | "text"
	read: boolean
	sender: Address | null
	from: Address[]
	to: Address[]
}

/** `index` is 1-based within the current, possibly search-filtered, list. */
export interface EmailPosition {
	index: number
	total: number
}

export interface EmailRemainder {
	raw: string
	header: Record<string, Array<string>>

	replyTo: Array<Address>
	cc: Array<Address>
	bcc: Array<Address>
	date: Date
	messageId: string
	inReplyTo: string[] | null
	references: string[] | null

	resentFrom: Array<Address>
	resentSender: Address | null
	resentTo: Array<Address>
	resentDate: Date
	resentCc: Array<Address>
	resentBcc: Array<Address>
	resentMessageId: string

	contentType: string
	// content: Readable // TODO

	htmlBody: string
	textBody: string

	attachments: Attachment[] | null
	embeddedFiles: EmbeddedFile[] | null
}
