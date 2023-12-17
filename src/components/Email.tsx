import {
	Accessor,
	For,
	Match,
	Show,
	Switch,
	createEffect,
	createMemo,
	createSignal,
	lazy,
	useContext,
} from "solid-js"
import type { Address, EmailBase, EmailRemainder } from "../email"
import { EmailAddr } from "./EmailAddress"
import { SelectedEmailContext } from "./App"
import { AttachmentButton } from "./AttachmentButton"
import { fetch, getUrl } from "../services/fetch"
import { AttachmentModal } from "./AttachmentModal"

const Code = lazy(() => import("./Code").then((m) => ({ default: m.Code })))

export interface EmailProps {}

export function Email({}: EmailProps) {
	const [email] = useContext(SelectedEmailContext)
	const mustEmail = () => email()!

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
		<div flex flex-col overflow-y-hidden h-screen>
			<Header email={mustEmail} emailRemainder={emailRemainder} />
			<Attachments email={mustEmail} emailRemainder={emailRemainder} />
			<Body email={mustEmail} emailRemainder={emailRemainder} />
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
	const [showAttachment, setShowAttachment] = createSignal<number>()

	const hasAttachments = () => (emailRemainder()?.attachments?.length ?? 0) > 0

	return (
		<>
			<Show when={hasAttachments()}>
				<div border-0 border-b border-b-solid border-zinc-800 p-4>
					<p m-0 text-zinc-500>
						Attachments:
					</p>
					<div flex flex-wrap gap-2 mt-1>
						<For each={emailRemainder()!.attachments}>
							{(attachment, index) => (
								<AttachmentButton
									attachment={() => attachment}
									onclick={() => setShowAttachment(index)}
								/>
							)}
						</For>
					</div>
				</div>
			</Show>
			<Show when={showAttachment() !== undefined}>
				<AttachmentModal
					onClose={() => setShowAttachment(undefined)}
					attachmentUrl={() =>
						getUrl(
							"/api/emails/" + email().id + "/attachments/" + showAttachment(),
						)
					}
					attachment={() => emailRemainder()!.attachments![showAttachment()!]}
				/>
			</Show>
		</>
	)
}

interface BodyProps {
	email: Accessor<EmailBase>
	emailRemainder: Accessor<EmailRemainder | undefined>
}

function Body({ email, emailRemainder }: BodyProps) {
	const htmlBody = () => emailRemainder()?.htmlBody ?? ""
	const fallbackBody = () => emailRemainder()?.textBody ?? ""

	return (
		<Switch>
			<Match when={email().bodyType === "text"}>
				<div p-4 flex-1 overflow-y-auto>
					<pre innerText={fallbackBody()} />
					<ContentTypeHint kind="text/plain" />
				</div>
			</Match>
			<Match when={email().bodyType === "html"}>
				<div flex-1>
					<HtmlBody email={email} html={htmlBody} plain={fallbackBody} />
					{/* <pre>{JSON.stringify(email(), null, 4)}</pre> */}
				</div>
			</Match>
		</Switch>
	)
}

interface HtmlBodyViewButtonProps {
	onclick?: () => void
	selected: Accessor<boolean>
	children?: any
}

function HtmlBodyViewButton({
	onclick,
	selected,
	children,
}: HtmlBodyViewButtonProps) {
	return (
		<button
			onclick={onclick}
			text={selected() ? "indigo-400" : "zinc-200"}
			border-0
		>
			{children}
		</button>
	)
}

interface HtmlBodyProps {
	email: Accessor<EmailBase>
	html: Accessor<string>
	plain: Accessor<string>
}

function HtmlBody({ email, html, plain }: HtmlBodyProps) {
	const [view, setView] = createSignal<"preview" | "plain" | "html">("preview")

	const isPreview = () => view() === "preview"
	const isPlain = () => view() === "plain"
	const isHtml = () => view() === "html"

	return (
		<div flex flex-col h-full>
			<div flex gap-1 p-4>
				<HtmlBodyViewButton
					onclick={() => setView("preview")}
					selected={isPreview}
				>
					Preview
				</HtmlBodyViewButton>
				<HtmlBodyViewButton onclick={() => setView("plain")} selected={isPlain}>
					Plain text
				</HtmlBodyViewButton>
				<HtmlBodyViewButton onclick={() => setView("html")} selected={isHtml}>
					HTML
				</HtmlBodyViewButton>
			</div>
			<Switch>
				<Match when={isPreview()}>
					<iframe
						flex-1
						src={getUrl(`/api/emails/${email().id}/page`)}
						w-full
						bg-white
						h-100
					/>
				</Match>
				<Match when={isPlain()}>
					<Show
						when={plain()}
						fallback={
							<p p-4 pt-0 text-zinc-500>
								This email does not have a plain text fallback 😞
							</p>
						}
					>
						<pre p-4 pt-0 overflow-y-auto innerText={plain()} />
					</Show>
				</Match>
				<Match when={isHtml()}>
					<div p-4 pt-0 overflow-y-auto>
						<Code lang="xml" code={html} />
					</div>
				</Match>
			</Switch>
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
