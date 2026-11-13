import { Accessor, createEffect, createSignal, useContext } from "solid-js"
import type { EmailPosition, EmailRemainder } from "../../email"
import { SelectedEmailContext } from "../App"
import { fetch } from "../../services/fetch"
import { Header } from "./Header"
import { Attachments } from "./Attachments"
import { Body } from "./Body"

export interface ShowEmailProps {
	position: Accessor<EmailPosition | undefined>
}

export function ShowEmail({ position }: ShowEmailProps) {
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
		<main
			class="widget"
			flex-1
			min-w-0
			flex
			flex-col
			gap="[14px]"
			px-4
			px="sm:[22px]"
			pt-4
			sm:pt-5
			pb-4
			pb="sm:[22px]"
			overflow-y-auto
		>
			<Header
				email={mustEmail}
				emailRemainder={emailRemainder}
				position={position}
			/>
			<Body email={mustEmail} emailRemainder={emailRemainder}>
				<Attachments email={mustEmail} emailRemainder={emailRemainder} />
			</Body>
		</main>
	)
}
