import { For, Show, Switch } from "solid-js"
import { EmailHint } from "../email"
import { Accessor } from "solid-js"
import { EmailAddr } from "./EmailAddress"

export interface EmailsListProps {
	emails: Accessor<Array<EmailHint> | undefined>
	setSelectedEmail: (email: EmailHint) => void
}

export function EmailsList({ emails, setSelectedEmail }: EmailsListProps) {
	return (
		<div>
			<h1 px-3>Emails</h1>
			<For each={emails()}>
				{(email) => (
					<EmailRow email={email} onClick={() => setSelectedEmail(email)} />
				)}
			</For>
		</div>
	)
}

export interface EmailRowProps {
	email: EmailHint
	onClick?: () => void
}

export function EmailRow({ email, onClick }: EmailRowProps) {
	const firstFrom = () => email.from[0]
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
		</button>
	)
}
