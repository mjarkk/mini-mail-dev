import {
	Accessor,
	createEffect,
	createSignal,
	lazy,
	Match,
	on,
	Show,
	Switch,
} from "solid-js"
import { EmailBase, EmailRemainder } from "../../email"
import { getUrl } from "../../services/fetch"

const Code = lazy(() => import("./Code").then((m) => ({ default: m.Code })))

export interface BodyProps {
	email: Accessor<EmailBase>
	emailRemainder: Accessor<EmailRemainder | undefined>
}

export function Body({ email, emailRemainder }: BodyProps) {
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
			<div flex gap-1 p-4>
				<Show when={email().bodyType === "html"}>
					<HtmlBodyViewButton
						onclick={() => setView("preview")}
						selected={isPreview}
					>
						Preview
					</HtmlBodyViewButton>
					<HtmlBodyViewButton onclick={() => setView("html")} selected={isHtml}>
						HTML
					</HtmlBodyViewButton>
				</Show>
				<HtmlBodyViewButton onclick={() => setView("plain")} selected={isPlain}>
					Plain text
				</HtmlBodyViewButton>
				<HtmlBodyViewButton onclick={() => setView("raw")} selected={isRaw}>
					Raw
				</HtmlBodyViewButton>
			</div>
			<div self-stretch overflow-y-auto>
				<Switch>
					<Match when={isPreview()}>
						<iframe
							flex-1
							src={getUrl(`/api/emails/${email().id}/page`)}
							w-full
							h-full
							bg-white
						/>
					</Match>
					<Match when={isHtml()}>
						<div p-4 pt-0>
							<Code lang="xml" code={html} />
						</div>
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
							<pre p-4 pt-0 innerText={plain()} />
						</Show>
					</Match>
					<Match when={isRaw()}>
						<pre p-4 pt-0 innerText={raw()} />
					</Match>
				</Switch>
			</div>
		</>
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
