import { Match, Switch, createEffect, createSignal } from "solid-js"
import type { EmailHint } from "../email"
import { EmailsList, EmailsListProps } from "./EmailRow"
import { Accessor } from "solid-js"
import { Email } from "./Email"

export function App() {
	const [emails, setEmails] = createSignal<Array<EmailHint>>()
	const [selectedEmail, setSelectedEmail] = createSignal<EmailHint>()

	createEffect(async () => {
		const response = await fetch("http://localhost:3000/api/emails")
		const data: Array<EmailHint> = await response.json()
		setEmails(data)
	})

	return (
		<div h-full w-full overflow-hidden>
			<Switch fallback={<div>Loading...</div>}>
				<Match when={emails() !== undefined && selectedEmail() !== undefined}>
					<LayoutWithEmail
						emails={emails}
						selectedEmail={() => selectedEmail()!}
						setSelectedEmail={setSelectedEmail}
					/>
				</Match>
				<Match when={emails() !== undefined}>
					<EmailsList emails={emails} setSelectedEmail={setSelectedEmail} />
				</Match>
			</Switch>
		</div>
	)
}

interface LayoutWithEmailProps extends EmailsListProps {
	selectedEmail: Accessor<EmailHint>
}

function LayoutWithEmail({
	emails,
	selectedEmail,
	setSelectedEmail,
}: LayoutWithEmailProps) {
	return (
		<div flex items-stretch>
			<div
				h-screen
				self-stretch
				overflow-y-auto
				w="[400px]"
				border-0
				border-r
				border-r-solid
				border-zinc-700
			>
				<EmailsList emails={emails} setSelectedEmail={setSelectedEmail} />
			</div>
			<div h-screen overflow-y-auto self-stretch flex-1>
				<Email email={selectedEmail} />
			</div>
		</div>
	)
}
