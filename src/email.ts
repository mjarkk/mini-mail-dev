export interface Address {
	name: string
	address: string
}

export interface Attachment {
	filename: string
	contentType: string
	// data: ; // TODO
}

export interface EmbeddedFile {
	cid: string
	contentType: string
	// data: ; // TODO
}

export interface EmailBase {
	id: string
	realDate: Date
	subject: string
	sender: Address | null
	from: Address[]
	to: Address[]
	attachments: Attachment[]
}

export interface Email extends EmailBase {
	// Recived form the client
	header: Record<string, Array<string>>

	replyTo: Array<Address>
	cc: Array<Address>
	bcc: Array<Address>
	date: Date
	messageId: string
	inReplyTo: string[]
	references: string[]

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

	attachments: Attachment[]
	embeddedFiles: EmbeddedFile[]
}

export interface EmailHint extends EmailBase {
	textBodyHint: string
}
