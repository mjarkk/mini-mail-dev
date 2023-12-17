import {
	For,
	Match,
	Show,
	Switch,
	createMemo,
	createSignal,
	useContext,
} from "solid-js"
import { Address, EmailBase } from "../email"
import { Accessor, onMount, onCleanup } from "solid-js"
import { EmailAddr } from "./EmailAddress"
import { SearchContext, SelectedEmailContext } from "./App"
import IconClose from "~icons/material-symbols/close-small"

export interface EmailsListProps {
	emails: Accessor<Array<EmailBase> | undefined>
	loading?: Accessor<boolean>
}

export function EmailsList({ emails, loading }: EmailsListProps) {
	const [selectedEmail, selectedEmailActions] = useContext(SelectedEmailContext)
	const [search, searchActions] = useContext(SearchContext)
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
			<div flex items-center flex-wrap gap-4 justify-between mb-2 px-2 py-4>
				<h1 m-0>Emails</h1>
				<div flex justify-center items-center gap-2>
					<Show when={search()}>
						<button
							outline-none
							rounded-lg
							p-2
							bg-transparent
							border-gray-500
							hover:border-gray-400
							onclick={() => searchActions.set("")}
						>
							<IconClose />
						</button>
					</Show>
					<input
						outline-none
						rounded-lg
						font-normal
						border-gray-500
						focus:border-gray-400
						placeholder="Search"
						w-full
						sm:w-auto
						px-4
						py-2
						value={search()}
						onInput={(e) => searchActions.set(e.target.value)}
					/>
				</div>
			</div>
			<Show
				when={!loading || !loading()}
				fallback={
					<p p-3 text-zinc-400>
						Loading...
					</p>
				}
			>
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
			</Show>
		</div>
	)
}

export interface EmailRowProps {
	display: Accessor<"sm" | "md" | "lg" | undefined>
	email: EmailBase
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
			text-zinc-500
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
					<p m-0 truncate>
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
			{(from) => <EmailAddr address={from} />}
		</Show>
	)
}

function ContentHint({ email }: { email: Accessor<EmailBase> }) {
	return (
		<>
			<span font-bold text-zinc-200>
				{email().subject}
			</span>
			<Show when={email().bodyHint}>
				<span font-normal text-zinc-400>
					{" "}
					- {email().bodyHint}
				</span>
			</Show>
		</>
	)
}

interface MailDateProps {
	date: Accessor<string>
	fullWidth?: boolean
}

function MailDate({ date, fullWidth }: MailDateProps) {
	return (
		<p m-0 text-zinc-500 text="13px" w={fullWidth ? "130px" : "auto"}>
			{date()}
		</p>
	)
}
