import {
	Accessor,
	createEffect,
	createResource,
	createSignal,
	For,
	JSX,
	lazy,
	Match,
	on,
	Show,
	Switch,
} from "solid-js"
import IconUnfold from "~icons/material-symbols/unfold-more"
import { EmailBase, EmailClient, EmailRemainder } from "../../email"
import { fetch, getUrl } from "../../services/fetch"

const Code = lazy(() => import("./Code").then((m) => ({ default: m.Code })))

// Outside the component so the choice carries over to the next email; undefined means the server's default.
const [chosenClient, setChosenClient] = createSignal<string>()

export interface BodyProps {
	email: Accessor<EmailBase>
	emailRemainder: Accessor<EmailRemainder | undefined>
	/** Rendered at the start of the view switcher row. */
	children?: JSX.Element
}

export function Body(props: BodyProps) {
	const { email, emailRemainder } = props
	const [view, setView] = createSignal<"preview" | "plain" | "html" | "raw">(
		"plain",
	)

	const isPreview = () => view() === "preview"
	const isPlain = () => view() === "plain"
	const isHtml = () => view() === "html"
	const isRaw = () => view() === "raw"

	const plain = () => emailRemainder()?.textBody ?? ""
	const html = () => emailRemainder()?.htmlBody ?? ""
	const raw = () => emailRemainder()?.raw ?? ""

	const [clients] = createResource<EmailClient[]>(async () => {
		const response = await fetch("/api/email-clients")
		return response.json()
	})
	const client = () => chosenClient() ?? clients()?.[0]?.id
	const selectedClient = () => clients()?.find((c) => c.id === client())
	const clientName = () => selectedClient()?.name
	const pageUrl = () =>
		`/api/emails/${email().id}/page?client=${encodeURIComponent(client() ?? "")}`

	createEffect(
		on(email, (newEmail) => {
			// When a email updates, check if the new content type is compatible with whats shown
			const currentView = view()
			if (newEmail.bodyType === "html") {
				if (currentView === "plain") setView("preview")
			} else {
				// Check if the current view is incompatible with the new email
				if (currentView === "preview" || currentView === "html")
					setView("plain")
			}
		}),
	)

	return (
		<>
			<div shrink-0 flex flex-wrap items-center gap="[10px]">
				{props.children}
				<span flex-1 />
				<Show when={isPreview() && clients()}>
					{(list) => (
						<label
							class="pill-button select-pill border-line text-fg2 hover:border-fg"
							relative
						>
							<span text-muted>Client</span>
							<span text-fg aria-hidden="true">
								{clientName()}
							</span>
							<IconUnfold aria-hidden="true" class="-mr-1" />
							{/* Transparent over the pill so the native picker handles input and accessibility. */}
							<select
								aria-label="Show as rendered by email client"
								value={client()}
								onChange={(e) => setChosenClient(e.currentTarget.value)}
								absolute
								inset-0
								w-full
								h-full
								opacity-0
								cursor-pointer
							>
								<For each={list()}>
									{(c) => <option value={c.id}>{c.name}</option>}
								</For>
							</select>
						</label>
					)}
				</Show>
				<div
					role="group"
					aria-label="Email view"
					inline-flex
					items-center
					ml-auto
					gap="[2px]"
					p="[3px]"
					rounded-full
					border
					border-solid
					border-line
				>
					<Show when={email().bodyType === "html"}>
						<Segment onclick={() => setView("preview")} selected={isPreview}>
							Preview
						</Segment>
						<Segment onclick={() => setView("html")} selected={isHtml}>
							HTML
						</Segment>
					</Show>
					<Segment onclick={() => setView("plain")} selected={isPlain}>
						Plain
					</Segment>
					<Segment onclick={() => setView("raw")} selected={isRaw}>
						Raw
					</Segment>
				</div>
			</div>
			<div
				class={isPreview() ? "bg-white" : "bg-card"}
				border
				border-solid
				border-line
				flex-1
				min-h="[320px]"
				rounded="[18px]"
				overflow-auto
			>
				<Switch>
					<Match when={isPreview() && client()}>
						<iframe
							title={`Email preview as shown by ${clientName() ?? "email client"}`}
							src={getUrl(pageUrl())}
							sandbox="allow-popups allow-popups-to-escape-sandbox"
							block
							w-full
							h-full
							border-0
							style={
								selectedClient()?.width
									? `max-width:${selectedClient()!.width}px;margin-inline:auto;box-shadow:0 0 0 1px var(--line)`
									: undefined
							}
						/>
					</Match>
					<Match when={isHtml()}>
						<Code lang="xml" code={html} />
					</Match>
					<Match when={isPlain()}>
						<Show
							when={plain()}
							fallback={
								<p class="mono-label" m-0 p-4>
									This email does not have a plain text fallback
								</p>
							}
						>
							<TextBlock text={plain()} />
						</Show>
					</Match>
					<Match when={isRaw()}>
						<TextBlock text={raw()} />
					</Match>
				</Switch>
			</div>
		</>
	)
}

function TextBlock(props: { text: string }) {
	return (
		<pre
			m-0
			p-4
			font-mono
			text="[13px]"
			leading-relaxed
			whitespace-pre-wrap
			break-words
			innerText={props.text}
		/>
	)
}

interface SegmentProps {
	onclick?: () => void
	selected: Accessor<boolean>
	children?: JSX.Element
}

function Segment({ onclick, selected, children }: SegmentProps) {
	return (
		<button
			type="button"
			onclick={onclick}
			aria-pressed={selected()}
			class={
				selected()
					? "segment bg-fg text-card"
					: "segment text-fg2 hover:text-fg hover:bg-card2"
			}
		>
			<Show when={selected()}>
				<span aria-hidden="true" w="[6px]" h="[6px]" rounded-full bg-signal />
			</Show>
			{children}
		</button>
	)
}
