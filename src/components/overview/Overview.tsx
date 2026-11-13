import { Accessor, For, onCleanup, onMount, Show, useContext } from "solid-js"
import { Header } from "./Header"
import { SearchContext, SelectedEmailContext } from "../App"
import { EmailBase } from "../../email"
import { EmailRow } from "./EmailRow"

export interface OverviewProps {
	emails: Accessor<Array<EmailBase> | undefined>
	loading?: Accessor<boolean>
	connected: Accessor<boolean>
}

export function Overview({ emails, loading, connected }: OverviewProps) {
	const [selectedEmail, selectedEmailActions] = useContext(SelectedEmailContext)
	const [search, searchActions] = useContext(SearchContext)

	const onKeydown = ({ target, key }: KeyboardEvent) => {
		const currentSelectedEmail = selectedEmail()
		const currentEmails = emails()
		if (!currentSelectedEmail || !currentEmails || currentEmails.length < 2) {
			return
		}

		const direction: 1 | 0 | -1 =
			key === "ArrowUp" ? -1 : key === "ArrowDown" ? 1 : 0

		if (direction === 0) {
			return
		}

		// Check if we are not focused on a input
		if (target) {
			const typedTarget = target as HTMLElement
			const stopTagNames = new Set(["input", "textarea", "area", "select"])
			if (stopTagNames.has(typedTarget.tagName.toLowerCase())) {
				return
			}
		}

		for (let idx = 0; idx < currentEmails.length; idx++) {
			const email = currentEmails[idx]
			if (email.id === currentSelectedEmail.id) {
				const newSelectedEmail = currentEmails[idx + direction]
				if (newSelectedEmail) {
					selectedEmailActions.select(newSelectedEmail)
				}
				break
			}
		}
	}

	onMount(() => {
		window.addEventListener("keydown", onKeydown)
	})

	onCleanup(() => {
		window.removeEventListener("keydown", onKeydown)
	})

	return (
		<div flex flex-col gap="[10px]" flex-1 min-w-0 min-h-0>
			<Header
				search={search}
				setSearch={searchActions.set}
				count={() => emails()?.length ?? 0}
				connected={connected}
			/>

			<div class="widget" flex-1 min-h-0 overflow-y-auto p="[6px]">
				<Show
					when={!loading || !loading()}
					fallback={
						<p class="mono-label" m-0 p-4>
							Loading
						</p>
					}
				>
					<For
						each={emails()}
						fallback={
							<p class="mono-label" m-0 p-4>
								{search() ? "No matches" : "No emails"}
							</p>
						}
					>
						{(email) => {
							const selected = () => email.id == selectedEmail()?.id
							return (
								<EmailRow
									selected={selected}
									email={email}
									onClick={() =>
										selected()
											? selectedEmailActions.deSelect()
											: selectedEmailActions.select(email)
									}
								/>
							)
						}}
					</For>
				</Show>
			</div>
		</div>
	)
}
