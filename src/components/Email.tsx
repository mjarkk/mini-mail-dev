import {
	Accessor,
	For,
	Show,
	createEffect,
	createMemo,
	createSignal,
	useContext,
} from "solid-js"
import type { Address, EmailBase, EmailRemainder } from "../email"
import { EmailAddr } from "./EmailAddress"
import { SelectedEmailContext } from "./App"
import { AttachmentButton } from "./AttachmentButton"
import { fetch } from "../services/fetch"

export interface EmailProps {}

export function Email({}: EmailProps) {
	const [email] = useContext(SelectedEmailContext)

	const [emailRemainder, setEmailRemainder] = createSignal<EmailRemainder>()

	const getFullMail = async (id: string) => {
		setEmailRemainder(undefined)
		const response = await fetch(`/api/emails/${id}/remainder`)
		const data = await response.json()
		if (response.status !== 200) throw data.error
		setEmailRemainder(data)
	}

	createEffect(() => {
		const e = email()
		if (e) getFullMail(e.id)
	})

	return (
		<div>
			<Header email={() => email()!} emailRemainder={emailRemainder} />
			<Attachments email={() => email()!} emailRemainder={emailRemainder} />
			<Body emailRemainder={emailRemainder} />
		</div>
	)
}

interface HeaderProps {
	email: Accessor<EmailBase>
	emailRemainder: Accessor<EmailRemainder | undefined>
}

function Header({ email, emailRemainder }: HeaderProps) {
	const [_, selectedEmailActions] = useContext(SelectedEmailContext)

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

		const remainder = emailRemainder()
		if (remainder) {
			mightAddList("ReplyTo", remainder.replyTo)
			mightAddList("CC", remainder.cc)
			mightAddList("BCC", remainder.bcc)
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
					<button border border-zinc-400 onclick={selectedEmailActions.delete}>
						Delete
					</button>
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

interface AttachmentsProps {
	email: Accessor<EmailBase>
	emailRemainder: Accessor<EmailRemainder | undefined>
}

function Attachments({ emailRemainder, email }: AttachmentsProps) {
	const downloadAttachment = (filename: string, index: number) => {
		const link = document.createElement("a")
		link.download = filename
		link.href =
			"http://localhost:3000/api/emails/" + email().id + "/attachments/" + index
		link.click()
	}

	const hasAttachments = () => (emailRemainder()?.attachments?.length ?? 0) > 0

	return (
		<Show when={hasAttachments()}>
			<div border-0 border-b border-b-solid border-zinc-800 p-4>
				<p m-0 text-zinc-500>
					Attachments:
				</p>
				<div flex flex-wrap gap-2 mt-1>
					<For each={emailRemainder()!.attachments}>
						{(attachment, index) => (
							<AttachmentButton
								{...attachment}
								onclick={() => downloadAttachment(attachment.filename, index())}
							/>
						)}
					</For>
				</div>
			</div>
		</Show>
	)
}

interface BodyProps {
	emailRemainder: Accessor<EmailRemainder | undefined>
}

function Body({ emailRemainder }: BodyProps) {
	const htmlBody = () => emailRemainder()?.htmlBody || null
	const fallbackBody = () => emailRemainder()?.textBody ?? ""

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
