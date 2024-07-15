import { createEffect, createSignal, JSX, useContext } from "solid-js"
import type { EmailRemainder } from "../../email"
import { SelectedEmailContext } from "../App"
import { fetch } from "../../services/fetch"
import { Header } from "./Header"
import { Attachments } from "./Attachments"
import { Body } from "./Body"

export interface ShowProps {}

export function Show({}: ShowProps) {
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

	const hasAttachments = () => (emailRemainder()?.attachments?.length ?? 0) > 0

	const bodyStyle = (): JSX.CSSProperties => ({
		"grid-template-columns": "1fr",
		"grid-template-rows": hasAttachments()
			? "auto auto 60px 1fr"
			: "auto 60px 1fr",
	})

	return (
		<div
			flex-1
			self-stretch
			overflow-y-hidden
			h-screen
			max-h-screen
			grid
			style={bodyStyle()}
		>
			<Header email={mustEmail} emailRemainder={emailRemainder} />
			<Attachments
				email={mustEmail}
				emailRemainder={emailRemainder}
				hasAttachments={hasAttachments}
			/>
			<Body email={mustEmail} emailRemainder={emailRemainder} />
		</div>
	)
}
