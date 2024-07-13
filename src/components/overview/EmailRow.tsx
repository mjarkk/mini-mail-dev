import { Match, Show, Switch, createMemo } from "solid-js"
import { Address, EmailBase } from "../../email"
import { Accessor } from "solid-js"
import { EmailAddr } from "../EmailAddress"

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
