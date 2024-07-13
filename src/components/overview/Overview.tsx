import {
	Accessor,
	createSignal,
	For,
	onCleanup,
	onMount,
	Show,
	useContext,
} from "solid-js"
import { Header } from "./Header"
import { SearchContext, SelectedEmailContext } from "../App"
import { EmailBase } from "../../email"
import { EmailRow } from "./EmailRow"

export interface OverviewProps {
	emails: Accessor<Array<EmailBase> | undefined>
	loading?: Accessor<boolean>
}

export function Overview({ emails, loading }: OverviewProps) {
	const [selectedEmail, selectedEmailActions] = useContext(SelectedEmailContext)
	const [search, searchActions] = useContext(SearchContext)
	const [display, setDisplay] = createSignal<"lg" | "md" | "sm">()
	const getDisplayValue = (w: number) => {
		if (!w) return undefined
		else if (w > 900) return "lg"
		else if (w > 600) return "md"
		else return "sm"
	}

	let ref: HTMLDivElement | null = null

	const onResize = () => {
		const width = ref?.getBoundingClientRect().width
		if (width) setDisplay(getDisplayValue(width))
	}

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
		window.addEventListener("resize", onResize)
		window.addEventListener("keydown", onKeydown)
		onResize()
	})

	onCleanup(() => {
		window.removeEventListener("resize", onResize)
		window.removeEventListener("keydown", onKeydown)
	})

	return (
		<div ref={(el) => (ref = el)}>
			<Header search={search} setSearch={searchActions.set} />

			<Show
				when={!loading || !loading()}
				fallback={
					<p p-3 text-zinc-400>
						Loading...
					</p>
				}
			>
				<For
					each={emails()}
					fallback={
						<p p-3 text-zinc-400>
							No emails
						</p>
					}
				>
					{(email) => {
						const selected = () => email.id == selectedEmail()?.id
						return (
							<EmailRow
								display={() => display()}
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
	)
}
