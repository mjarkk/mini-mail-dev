import {
	Match,
	Show,
	Switch,
	createContext,
	createEffect,
	createSignal,
	lazy,
	onCleanup,
	onMount,
} from "solid-js"
import type { EmailBase, EmailPosition } from "../email"
import { Overview } from "./overview/Overview"
import { EmptyReader } from "./EmptyReader"
import type { Accessor, Setter } from "solid-js"
import { EmailEventsWebsocket } from "../services/websocket"
import { fetch } from "../services/fetch"

const ShowEmail = lazy(() =>
	import("./show/Show").then((m) => ({ default: m.ShowEmail })),
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

const calculateIsMobile = () => window.innerWidth < 840

export function App() {
	const [emails, setEmails] = createSignal<Array<EmailBase>>()
	const [selectedEmail, setSelectedEmail] = createSignal<EmailBase>()
	const [searchValue, setSearchValue] = createSignal<string>("")
	const [isMobile, setIsMobile] = createSignal(calculateIsMobile())
	const [connected, setConnected] = createSignal(false)

	const fetchEmails = throttleFn(async () => {
		let url = "/api/emails"

		if (searchValue()?.length) {
			url += "?search=" + encodeURIComponent(searchValue())
		}

		const response = await fetch(url)
		const data: Array<EmailBase> = await response.json()
		setEmails(data)
	})

	const emailsWebsocket = new EmailEventsWebsocket(fetchEmails, setConnected)

	createEffect(() => {
		searchValue()
		fetchEmails()
	})

	const onResize = () => {
		const newIsMobile = calculateIsMobile()
		if (newIsMobile !== isMobile()) setIsMobile(newIsMobile)
	}

	onMount(() => {
		emailsWebsocket.start()
		window.addEventListener("resize", onResize)
	})

	onCleanup(() => {
		emailsWebsocket.close()
		window.removeEventListener("resize", onResize)
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

	const select = async (email: EmailBase) => {
		setSelectedEmail(email)
		if (email.read) return

		setEmails((list) =>
			list?.map((e) => (e.id === email.id ? { ...e, read: true } : e)),
		)
		await fetch(`/api/emails/${email.id}/read`, { method: "POST" })
		// A list fetch that started before the POST landed would bring the dot back.
		await fetchEmails()
	}

	const selectedEmailContext: () => SelectedEmailContextType = () => [
		selectedEmail,
		{
			delete: deleteSelected,
			select,
			deSelect: () => setSelectedEmail(undefined),
		},
	]

	const searchContext: () => SearchContextType = () => [
		searchValue,
		{ set: setSearchValue },
	]

	const showEmail = () =>
		emails() !== undefined && selectedEmail() !== undefined

	const list = () => (
		<Overview
			emails={() => emails() ?? []}
			loading={() => !emails()}
			connected={connected}
		/>
	)

	const position = (): EmailPosition | undefined => {
		const list = emails() ?? []
		const id = selectedEmail()?.id
		const index = list.findIndex((e) => e.id === id)
		return index < 0 ? undefined : { index: index + 1, total: list.length }
	}

	return (
		<SelectedEmailContext.Provider value={selectedEmailContext()}>
			<SearchContext.Provider value={searchContext()}>
				<div
					h="[100dvh]"
					w-full
					box-border
					flex
					gap="[10px]"
					p-2
					p="sm:[10px]"
					overflow-hidden
				>
					<Switch
						fallback={
							<>
								<div w="[380px]" shrink-0 flex>
									{list()}
								</div>
								<Show
									when={showEmail()}
									fallback={<EmptyReader inboxEmpty={emails()?.length === 0} />}
								>
									<ShowEmail position={position} />
								</Show>
							</>
						}
					>
						<Match when={isMobile() && showEmail()}>
							<ShowEmail position={position} />
						</Match>
						<Match when={isMobile()}>{list()}</Match>
					</Switch>
				</div>
			</SearchContext.Provider>
		</SelectedEmailContext.Provider>
	)
}
