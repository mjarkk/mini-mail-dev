import { Accessor, createMemo, For, useContext } from "solid-js"
import { Address, EmailBase, EmailRemainder } from "../../email"
import { SelectedEmailContext } from "../App"
import { EmailAddr } from "../EmailAddress"
import IconClose from "~icons/material-symbols/close-small"

export interface HeaderProps {
	email: Accessor<EmailBase>
	emailRemainder: Accessor<EmailRemainder | undefined>
}

export function Header({ email, emailRemainder }: HeaderProps) {
	const [_, selectedEmailActions] = useContext(SelectedEmailContext)

	const fields = createMemo(() => {
		const e = email()

		const resp: Array<{
			name: string
			value: Array<any>
		}> = []

		const mightAddList = (name: string, list: Array<Address>) => {
			if (list && list.length > 0) {
				resp.push({
					name,
					value: list.map((addr) => <EmailAddr address={() => addr} />),
				})
			}
		}

		if (e.sender) {
			resp.push({
				name: "Sender",
				value: [<EmailAddr address={() => e.sender!} />],
			})
		}
		mightAddList("From", e.from)
		mightAddList("To", e.to)

		const remainder = emailRemainder()
		if (remainder) {
			mightAddList("ReplyTo", remainder.replyTo)
			mightAddList("CC", remainder.cc)
			mightAddList("BCC", remainder.bcc)
		}

		return resp
	})

	return (
		<div bg-zinc-900>
			<div
				flex
				justify-between
				items-center
				p-4
				border-0
				border-b
				border-b-solid
				border-zinc-800
			>
				<h2 font-bold m-0 flex items-center>
					<button
						onclick={selectedEmailActions.deSelect}
						p-0
						text-zinc-400
						inline-flex
						mt-1
						mr-2
					>
						<IconClose />
					</button>
					{email().subject}
				</h2>
				<div>
					<button
						bg-zinc-700
						hover:bg-zinc-600
						onclick={selectedEmailActions.delete}
					>
						Delete
					</button>
				</div>
			</div>
			<div border-0 border-b border-b-solid border-zinc-800 p-4 flex flex-wrap>
				<For each={fields()}>
					{(field) => (
						<span mr-3>
							{field.name}: {field.value}
						</span>
					)}
				</For>
			</div>
		</div>
	)
}
