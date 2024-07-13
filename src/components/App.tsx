import {
	Match,
	Switch,
	createContext,
	createEffect,
	createSignal,
	lazy,
	onCleanup,
	onMount,
} from "solid-js"
import type { EmailBase } from "../email"
import { Overview, OverviewProps } from "./overview/Overview"
import type { Accessor, Setter } from "solid-js"
import { EmailEventsWebsocket } from "../services/websocket"
import { fetch } from "../services/fetch"

const Show = lazy(() =>
	import("./show/Show").then((m) => ({ default: m.Show })),
)

interface SelectedEmailActions {
	delete: () => void
	select: (e: EmailBase) => void
	deSelect: () => void
}

type SelectedEmailContextType = [
	Accessor<EmailBase | undefined>,
	SelectedEmailActions,
]

export const SelectedEmailContext = createContext<SelectedEmailContextType>(
	[] as unknown as SelectedEmailContextType,
)

type SearchContextType = [Accessor<string>, { set: Setter<string> }]

export const SearchContext = createContext<SearchContextType>(
	[] as unknown as SearchContextType,
)

function throttleFn(fn: () => Promise<void>): () => Promise<void> {
	let running = false
	let reRun = false
	return async () => {
		if (running) {
			reRun = true
			return
		}
		running = true
		reRun = true
		while (reRun) {
			reRun = false
			try {
				await fn()
			} catch (e) {
				console.error(e)
			}
			await new Promise((res) => setTimeout(res, 500))
		}
		running = false
	}
}

export function App() {
	const [emails, setEmails] = createSignal<Array<EmailBase>>()
	const [selectedEmail, setSelectedEmail] = createSignal<EmailBase>()
	const [searchValue, setSearchValue] = createSignal<string>("")

	const fetchEmails = throttleFn(async () => {
		let url = "/api/emails"

		if (searchValue()?.length) {
			url += "?search=" + encodeURIComponent(searchValue())
		}

		const response = await fetch(url)
		const data: Array<EmailBase> = await response.json()
		setEmails(data)
	})

	const emailsWebsocket = new EmailEventsWebsocket(fetchEmails)

	createEffect(() => {
		searchValue()
		fetchEmails()
	})

	onMount(() => {
		emailsWebsocket.start()
	})

	onCleanup(() => {
		emailsWebsocket.close()
	})

	const deleteSelected = async () => {
		const id = selectedEmail()?.id
		if (!id) return

		setSelectedEmail(undefined)
		await fetch("/api/emails/" + id, {
			method: "DELETE",
		})
		await fetchEmails()
	}

	const selectedEmailContext: () => SelectedEmailContextType = () => [
		selectedEmail,
		{
			delete: deleteSelected,
			select: setSelectedEmail,
			deSelect: () => setSelectedEmail(undefined),
		},
	]

	const searchContext: () => SearchContextType = () => [
		searchValue,
		{ set: setSearchValue },
	]

	return (
		<SelectedEmailContext.Provider value={selectedEmailContext()}>
			<SearchContext.Provider value={searchContext()}>
				<div h-full w-full overflow-hidden>
					<Switch
						fallback={
							<Overview
								emails={() => emails() ?? []}
								loading={() => !emails()}
							/>
						}
					>
						<Match
							when={emails() !== undefined && selectedEmail() !== undefined}
						>
							<LeftOverviewRightEmail
								emails={emails}
								selectedEmail={() => selectedEmail()!}
							/>
						</Match>
					</Switch>
				</div>
			</SearchContext.Provider>
		</SelectedEmailContext.Provider>
	)
}

interface LeftOverviewRightEmailProps extends OverviewProps {
	selectedEmail: Accessor<EmailBase>
}

function LeftOverviewRightEmail({ emails }: LeftOverviewRightEmailProps) {
	return (
		<div flex items-stretch>
			<div
				h-screen
				self-stretch
				overflow-y-auto
				w="[400px]"
				border-0
				border-r
				border-r-solid
				border-zinc-700
			>
				<Overview emails={emails} />
			</div>
			<div h-screen overflow-y-auto self-stretch flex-1>
				<Show />
			</div>
		</div>
	)
}
