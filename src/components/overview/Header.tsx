import { Accessor, Show } from "solid-js"
import IconClose from "~icons/material-symbols/close-small"
import IconSearch from "~icons/material-symbols/search"
import { DotMatrix } from "../DotText"
import { createTween } from "../../tween"

export interface HeaderProps {
	setSearch: (v: string) => void
	search: Accessor<string>
	count: Accessor<number>
	connected: Accessor<boolean>
}

export function Header({ setSearch, search, count, connected }: HeaderProps) {
	const shownCount = createTween(count)

	return (
		<>
			<header
				class="widget"
				shrink-0
				flex
				flex-col
				gap="[18px]"
				px-5
				pt="[18px]"
				pb-5
			>
				<div class="mono-label" flex items-center gap-2>
					<span flex-1>Mini Mail Dev</span>
					<span
						class={connected() ? "bg-signal" : "bg-muted"}
						w="[7px]"
						h="[7px]"
						rounded-full
					/>
					<span text-fg>{connected() ? "Live" : "Offline"}</span>
				</div>
				<h1 m-0 text-fg>
					<DotMatrix
						start="INBOX"
						end={String(shownCount()).padStart(2, "0")}
						endColor="var(--signal)"
						label={`Inbox, ${count()} emails`}
						dot={6}
					/>
				</h1>
			</header>

			<label
				class="search-pill"
				shrink-0
				flex
				items-center
				gap="[10px]"
				h-12
				pl="[18px]"
				pr-2
				rounded-full
				bg-card
				text-muted
			>
				<IconSearch aria-hidden="true" />
				<input
					type="search"
					aria-label="Search emails"
					placeholder="SEARCH"
					flex-1
					min-w-0
					h-full
					border-0
					bg-transparent
					outline-none
					text-fg
					font-mono
					class="tracking-[0.08em]"
					text="[12px]"
					value={search()}
					onInput={(e) => setSearch(e.target.value)}
				/>
				<Show when={search()}>
					<button
						type="button"
						aria-label="Clear search"
						inline-flex
						items-center
						justify-center
						w-8
						h-8
						rounded-full
						text-fg
						hover:bg-card2
						onClick={() => setSearch("")}
					>
						<IconClose aria-hidden="true" />
					</button>
				</Show>
			</label>
		</>
	)
}
