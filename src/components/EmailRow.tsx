import {
	For,
	Match,
	Show,
	Switch,
	createMemo,
	createSignal,
	useContext,
} from "solid-js"
import { Address, EmailHint } from "../email"
import { Accessor, onMount, onCleanup } from "solid-js"
import { EmailAddr } from "./EmailAddress"
import { SelectedEmailContext } from "./App"

export interface EmailsListProps {
	emails: Accessor<Array<EmailHint> | undefined>
}

export function EmailsList({ emails }: EmailsListProps) {
	const [selectedEmail, selectedEmailActions] = useContext(SelectedEmailContext)
	const [display, setDisplay] = createSignal<"lg" | "md" | "sm">()
	const getDisplayValue = (w: number) => {
		if (!w) return undefined
		else if (w > 900) return "lg"
		else if (w > 600) return "md"
		else return "sm"
	}

	let ref: HTMLDivElement | null = null

	const onResize = () => {
		const width = ref?.getBoundingClientRect().width
		if (width) setDisplay(getDisplayValue(width))
	}

	onMount(() => {
		window.addEventListener("resize", onResize)
		onResize()
	})

	onCleanup(() => {
		window.removeEventListener("resize", onResize)
	})

	return (
		<div ref={(el) => (ref = el)}>
			<h1 px-3>Emails</h1>
			<For
				each={emails()}
				fallback={
					<p p-3 text-zinc-400>
						No emails
					</p>
				}
			>
				{(email) => {
					const selected = () => email.id == selectedEmail()?.id
					return (
						<EmailRow
							display={() => display()}
							selected={selected}
							email={email}
							onClick={() =>
								selected()
									? selectedEmailActions.deSelect()
									: selectedEmailActions.select(email)
							}
						/>
					)
				}}
			</For>
		</div>
	)
}

export interface EmailRowProps {
	display: Accessor<"sm" | "md" | "lg" | undefined>
	email: EmailHint
	onClick?: () => void
	selected: Accessor<boolean>
}

export function EmailRow({ display, selected, email, onClick }: EmailRowProps) {
	const firstFrom = () => email.from[0]

	const date = createMemo(() => {
		const d = new Date(email.realDate)
		return d.toLocaleDateString() + " " + d.toLocaleTimeString()
	})

	return (
		<button
			class="block"
			onClick={onClick}
			cursor-pointer
			w-full
			text-left
			border-0
			border-t
			border-t-solid
			border-zinc-800
			p-3
			bg={selected() ? "zinc-800" : "zinc-950"}
			rounded-none
		>
			<Switch>
				<Match when={display() === "lg"}>
					<div flex items-center>
						<p m-0 truncate w="[300px]">
							<From from={firstFrom} />
						</p>
						<p my-0 mx-3 truncate flex-1>
							<ContentHint email={() => email} />
						</p>
						<MailDate date={date} />
					</div>
				</Match>
				<Match when={display() === "md"}>
					<div flex items-center>
						<p m-0 truncate flex-1>
							<From from={firstFrom} />
						</p>
						<MailDate date={date} />
					</div>
					<p m-0 mt-0 truncate>
						<ContentHint email={() => email} />
					</p>
				</Match>
				<Match when={display() === "sm"}>
					<p m-0 truncate>
						<From from={firstFrom} />
					</p>
					<p m-0 mb-1 truncate>
						<ContentHint email={() => email} />
					</p>
					<MailDate date={date} />
				</Match>
			</Switch>
		</button>
	)
}

function From({ from }: { from: Accessor<Address | undefined> }) {
	return (
		<Show when={from()} fallback={<span text-zinc-500>Unknown address</span>}>
			<EmailAddr address={() => from()!} />
		</Show>
	)
}

function ContentHint({ email }: { email: Accessor<EmailHint> }) {
	return (
		<>
			<span font-bold text-zinc-200>
				{email().subject}
			</span>
			<Show when={email().textBodyHint}>
				<span font-normal text-zinc-400>
					{" "}
					- {email().textBodyHint}
				</span>
			</Show>
		</>
	)
}

function MailDate({ date }: { date: Accessor<string> }) {
	return (
		<p m-0 text-zinc-500 text="13px" w="[130px]">
			{date()}
		</p>
	)
}
