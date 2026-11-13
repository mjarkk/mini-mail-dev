import { Accessor, Show } from "solid-js"
import { Address } from "../email"

export interface EmailAddrProps {
	address: Accessor<Address>
}

export function EmailAddr({ address }: EmailAddrProps) {
	const name = () => address().name
	const addr = () => address().address

	return (
		<Show when={name()} fallback={<span block truncate text-fg2>{addr()}</span>}>
			<span block truncate font-medium>
				{name()}
			</span>
			<Show when={addr()}>
				<span block truncate text-xs text-muted>
					{addr()}
				</span>
			</Show>
		</Show>
	)
}
