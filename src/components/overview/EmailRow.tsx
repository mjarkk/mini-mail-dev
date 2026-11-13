import { Show } from "solid-js"
import { EmailBase } from "../../email"
import { Accessor } from "solid-js"

export interface EmailRowProps {
	email: EmailBase
	onClick?: () => void
	selected: Accessor<boolean>
}

export function EmailRow({ selected, email, onClick }: EmailRowProps) {
	const date = new Date(email.realDate)
	const sender = email.from[0]?.name || email.from[0]?.address

	return (
		<button
			type="button"
			onClick={onClick}
			aria-current={selected() ? "true" : undefined}
			class={selected() ? "bg-card2" : "hover:bg-card2"}
			w-full
			grid
			grid-cols="[12px_1fr]"
			gap-x-2
			text-left
			rounded-2xl
			py-3
			pl="[10px]"
			pr="[14px]"
			transition-colors
		>
			<span
				class={email.read ? "" : "bg-signal"}
				aria-hidden="true"
				w="[7px]"
				h="[7px]"
				mt="[7px]"
				rounded-full
			/>
			<span flex flex-col gap-px min-w-0>
				<span flex items-baseline gap-2>
					<span flex-1 min-w-0 truncate font-medium text-fg>
						<Show when={!email.read}>
							<span sr-only>Unread: </span>
						</Show>
						{sender || "Unknown sender"}
					</span>
					<time
						class="mono-label"
						shrink-0
						dateTime={date.toISOString()}
						title={date.toLocaleString()}
					>
						{shortDate(date)}
					</time>
				</span>
				<span truncate text-fg2>
					{email.subject || "(no subject)"}
				</span>
				<Show when={email.bodyHint}>
					<span truncate text-muted text-sm>
						{email.bodyHint}
					</span>
				</Show>
			</span>
		</button>
	)
}

function shortDate(date: Date): string {
	const now = new Date()
	if (date.toDateString() === now.toDateString()) {
		return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })
	}
	return date.toLocaleDateString([], {
		day: "numeric",
		month: "short",
		year: date.getFullYear() === now.getFullYear() ? undefined : "numeric",
	})
}
