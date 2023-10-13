import { Accessor, Show } from "solid-js"
import { Address } from "../email"

export interface EmailAddrProps {
	address: Accessor<Address>
}

export function EmailAddr({ address }: EmailAddrProps) {
	const name = () => address().name
	const addr = () => address().address

	return (
		<span font-bold text-zinc-200>
			{name()}
			<Show when={addr()}>
				<span font-normal text-zinc-500>{` <${addr()}>`}</span>
			</Show>
		</span>
	)
}
