import { DotText } from "./DotText"

export function EmptyReader(props: { inboxEmpty: boolean }) {
	return (
		<main
			class="widget"
			flex-1
			min-w-0
			flex
			flex-col
			items-center
			justify-center
			gap-5
			p-8
			text-center
		>
			<div aria-hidden="true" text-muted mb-1>
				<DotText text="✉" dot={8} />
			</div>
			<p class="mono-label" m-0>
				{props.inboxEmpty ? "Waiting for emails" : "No email selected"}
			</p>
			<p m-0 text-sm text-fg2 max-w="[280px]">
				{props.inboxEmpty
					? "New emails show up here as soon as they arrive."
					: "Pick an email from the list to read it here."}
			</p>
		</main>
	)
}
