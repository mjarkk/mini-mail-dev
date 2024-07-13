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

	onMount(() => {
		window.addEventListener("resize", onResize)
		onResize()
	})

	onCleanup(() => {
		window.removeEventListener("resize", onResize)
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
