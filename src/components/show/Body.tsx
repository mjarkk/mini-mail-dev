import { Accessor, createSignal, lazy, Match, Show, Switch } from "solid-js"
import { EmailBase, EmailRemainder } from "../../email"
import { getUrl } from "../../services/fetch"

const Code = lazy(() => import("./Code").then((m) => ({ default: m.Code })))

export interface BodyProps {
	email: Accessor<EmailBase>
	emailRemainder: Accessor<EmailRemainder | undefined>
}

export function Body({ email, emailRemainder }: BodyProps) {
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
