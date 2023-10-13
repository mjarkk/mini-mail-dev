import { For, Show, createMemo, useContext } from "solid-js"
import { EmailHint } from "../email"
import { Accessor } from "solid-js"
import { EmailAddr } from "./EmailAddress"
import { SelectedEmailContext } from "./App"

export interface EmailsListProps {
	emails: Accessor<Array<EmailHint> | undefined>
}

export function EmailsList({ emails }: EmailsListProps) {
	const [selectedEmail, selectedEmailActions] = useContext(SelectedEmailContext)

	return (
		<div>
			<h1 px-3>Emails</h1>
			<For
				each={emails()}
				fallback={
					<p p-3 text-zinc-400>
						No emails
					</p>
				}
			>
				{(email) => (
					<EmailRow
						selected={() => email.id == selectedEmail()?.id}
						email={email}
						onClick={() => selectedEmailActions.select(email)}
					/>
				)}
			</For>
		</div>
	)
}

export interface EmailRowProps {
	email: EmailHint
	onClick?: () => void
	selected: Accessor<boolean>
}

export function EmailRow({ selected, email, onClick }: EmailRowProps) {
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
			<p m-0 truncate>
				<Show
					when={firstFrom()}
					fallback={<span text-zinc-500>Unknown address</span>}
				>
					<EmailAddr address={firstFrom} />
				</Show>
			</p>
			<p m-0 mt-0 truncate>
				<span font-bold text-zinc-200>
					{email.subject}
				</span>
				<Show when={email.textBodyHint}>
					<span font-normal text-zinc-400>
						{" "}
						- {email.textBodyHint}
					</span>
				</Show>
			</p>
			<p m-0 mt-1 text-zinc-500>
				{date()}
			</p>
		</button>
	)
}
