import { Accessor, Show } from "solid-js"
import IconClose from "~icons/material-symbols/close-small"

export interface HeaderProps {
	setSearch: (v: string) => void
	search: Accessor<string>
}

export function Header({ setSearch, search }: HeaderProps) {
	return (
		<div flex items-center flex-wrap gap-4 justify-between mb-2 px-2 py-4>
			<h1 m-0>Emails</h1>
			<div flex justify-center items-center gap-2>
				<Show when={search()}>
					<button
						outline-none
						rounded-lg
						p-2
						bg-transparent
						border-gray-500
						hover:border-gray-400
						onclick={() => setSearch("")}
					>
						<IconClose />
					</button>
				</Show>
				<input
					outline-none
					rounded-lg
					font-normal
					border-gray-500
					focus:border-gray-400
					placeholder="Search"
					w-full
					sm:w-auto
					px-4
					py-2
					value={search()}
					onInput={(e) => setSearch(e.target.value)}
				/>
			</div>
		</div>
	)
}
