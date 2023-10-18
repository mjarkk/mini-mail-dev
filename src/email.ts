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

export interface EmailBase {
	id: string
	realDate: Date
	subject: string
	bodyHint: string
	bodyType: "html" | "text"
	sender: Address | null
	from: Address[]
	to: Address[]
}

export interface EmailRemainder {
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
