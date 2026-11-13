import { Accessor, createMemo, For, Show, useContext } from "solid-js"
import { Address, EmailBase, EmailPosition, EmailRemainder } from "../../email"
import { SelectedEmailContext } from "../App"
import { EmailAddr } from "../EmailAddress"
import { DotText } from "../DotText"
import IconClose from "~icons/material-symbols/close-small"
import { createTween } from "../../tween"

export interface HeaderProps {
	email: Accessor<EmailBase>
	emailRemainder: Accessor<EmailRemainder | undefined>
	position: Accessor<EmailPosition | undefined>
}

export function Header({ email, emailRemainder, position }: HeaderProps) {
	const [_, selectedEmailActions] = useContext(SelectedEmailContext)

	const fields = createMemo(() => {
		const e = email()
		const remainder = emailRemainder()

		const candidates: Array<[string, Array<Address> | undefined]> = [
			["Sender", e.sender ? [e.sender] : undefined],
			["From", e.from],
			["To", e.to],
			["Reply-To", remainder?.replyTo],
			["CC", remainder?.cc],
			["BCC", remainder?.bcc],
		]

		return candidates
			.filter(([, addresses]) => addresses && addresses.length > 0)
			.map(([name, addresses]) => ({ name, addresses: addresses! }))
	})

	const received = () => new Date(email().realDate)
	const pad = (n: number) => String(n).padStart(2, "0")
	const shownMinutes = createTween(
		() => received().getHours() * 60 + received().getMinutes(),
	)

	return (
		<>
			<div shrink-0 flex items-center gap-2>
				<span class="mono-label" flex-1>
					Message
					<Show when={position()}>
						{(p) => ` ${pad(p().index)}/${pad(p().total)}`}
					</Show>
				</span>
				<button
					type="button"
					class="outline-button"
					onClick={selectedEmailActions.delete}
				>
					Delete
				</button>
				<button
					type="button"
					aria-label="Close email"
					inline-flex
					items-center
					justify-center
					w-9
					h-9
					rounded-full
					border
					border-solid
					border-fg
					text-fg
					hover:bg-fg
					hover:text-card
					transition-colors
					onClick={selectedEmailActions.deSelect}
				>
					<IconClose aria-hidden="true" />
				</button>
			</div>

			<div shrink-0 flex flex-wrap items-stretch gap="[10px]">
				<h2
					flex-1
					min-w="[min(100%,240px)]"
					self-center
					m-0
					class="tracking-[-0.03em] leading-[1.15]"
					text="[22px] sm:[28px]"
					font-medium
					break-words
				>
					{email().subject || "(no subject)"}
				</h2>
				<div
					flex
					flex-col
					justify-center
					gap-2
					px-4
					py-3
					rounded="[18px]"
					bg-card2
					title={received().toLocaleString()}
				>
					<DotText
						text={`${pad(Math.floor(shownMinutes() / 60))}:${pad(shownMinutes() % 60)}`}
						label={`Received ${received().toLocaleString()}`}
						dot={4}
						unlit="var(--unlit-on-card2)"
						tabular
					/>
					<span class="mono-label" text="[10px]">
						{received().toLocaleDateString([], {
							day: "numeric",
							month: "short",
							year: "numeric",
						})}
					</span>
				</div>
			</div>

			<dl
				shrink-0
				m-0
				grid
				gap-2
				grid-cols="[repeat(auto-fit,minmax(180px,1fr))]"
			>
				<For each={fields()}>
					{(field) => (
						<div min-w-0 px="[14px]" py="[10px]" rounded-2xl bg-card2>
							<dt class="mono-label" text="[10px]">
								{field.name}
							</dt>
							<For each={field.addresses}>
								{(address) => (
									<dd m-0 mt="[2px]" min-w-0>
										<EmailAddr address={() => address} />
									</dd>
								)}
							</For>
						</div>
					)}
				</For>
			</dl>
		</>
	)
}
