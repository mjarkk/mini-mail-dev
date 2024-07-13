import { createEffect, createSignal, useContext } from "solid-js"
import type { EmailRemainder } from "../../email"
import { SelectedEmailContext } from "../App"
import { fetch } from "../../services/fetch"
import { Header } from "./Header"
import { Attachments } from "./Attachments"
import { Body } from "./Body"

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
