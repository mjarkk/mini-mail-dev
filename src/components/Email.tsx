import {
	Accessor,
	For,
	Show,
	createEffect,
	createMemo,
	createSignal,
} from "solid-js"
import type { Address, Email, EmailHint } from "../email"
import { EmailAddr } from "./EmailAddress"

export interface EmailProps {
	email: Accessor<EmailHint>
}

export function Email({ email }: EmailProps) {
	const [fullEmail, setFullEmail] = createSignal<Email>()

	const getFullMail = async (id: string) => {
		setFullEmail(undefined)
		const response = await fetch("http://localhost:3000/api/emails/" + id)
		const data = await response.json()
		setFullEmail(data)
	}

	createEffect(() => {
		getFullMail(email().id)
	})

	const fullOrPartialEmail = createMemo(() => fullEmail() ?? email())

	return (
		<div>
			<Header email={fullOrPartialEmail} />
			<Attachments />
			<Body email={fullOrPartialEmail} />
		</div>
	)
}

interface HeaderProps {
	email: Accessor<EmailHint | Email>
}

function Header({ email }: HeaderProps) {
	const fields = createMemo(() => {
		const e = email()

		const resp: Array<{
			name: string
			value: Array<any>
		}> = []

		const mightAddList = (name: string, list: Array<Address>) => {
			if (list && list.length > 0) {
				resp.push({
					name,
					value: list.map((addr) => <EmailAddr address={() => addr} />),
				})
			}
		}

		if (e.sender) {
			resp.push({
				name: "Sender",
				value: [<EmailAddr address={() => e.sender!} />],
			})
		}
		mightAddList("From", e.from)
		mightAddList("To", e.to)
		if ("cc" in e) {
			mightAddList("ReplyTo", e.replyTo)
			mightAddList("CC", e.cc)
			mightAddList("BCC", e.bcc)
		}

		return resp
	})

	return (
		<div bg-zinc-900>
			<div
				flex
				justify-between
				items-center
				p-4
				border-0
				border-b
				border-b-solid
				border-zinc-800
			>
				<h2 font-bold m-0>
					{email().subject}
				</h2>
				<div>
					<button>Delete</button>
				</div>
			</div>
			<div border-0 border-b border-b-solid border-zinc-800 p-4 flex flex-wrap>
				<For each={fields()}>
					{(field) => (
						<span mr-3>
							{field.name}: {field.value}
						</span>
					)}
				</For>
			</div>
		</div>
	)
}

interface AttachmentsProps {}

function Attachments({}: AttachmentsProps) {
	return <div></div>
}

interface BodyProps {
	email: Accessor<EmailHint | Email>
}

function Body({ email }: BodyProps) {
	const htmlBody = () => {
		const e = email()
		return "htmlBody" in e ? e.htmlBody || null : null
	}

	const fallbackBody = () => {
		const e = email()
		if ("textBody" in e && e.textBody) return e.textBody
		if ("textBodyHint" in e && e.textBodyHint) return e.textBodyHint
		return ""
	}

	// TODO: Fix remove everything scriptable from the email html body

	return (
		<div p-4>
			<Show
				when={htmlBody()}
				fallback={
					<>
						<pre innerText={fallbackBody()} />
						<ContentTypeHint kind="text/plain" />
					</>
				}
			>
				<div innerHTML={htmlBody()!} />
				<ContentTypeHint kind="text/html" />
			</Show>
			{/* <pre>{JSON.stringify(email(), null, 4)}</pre> */}
		</div>
	)
}

function ContentTypeHint({ kind }: { kind: "text/plain" | "text/html" }) {
	return (
		<p italic text-zinc-500>
			{kind}
		</p>
	)
}
