import {
	Match,
	Switch,
	createContext,
	createEffect,
	createSignal,
	lazy,
	onCleanup,
} from "solid-js"
import type { EmailBase } from "../email"
import { EmailsList, EmailsListProps } from "./EmailRow"
import { Accessor } from "solid-js"
import { EmailEventsWebsocket } from "../services/websocket"
import { fetch } from "../services/fetch"

const Email = lazy(() => import("./Email").then((m) => ({ default: m.Email })))

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
		}
		await new Promise((res) => setTimeout(res, 500))
		running = false
	}
}

export function App() {
	const [emails, setEmails] = createSignal<Array<EmailBase>>()
	const [selectedEmail, setSelectedEmail] = createSignal<EmailBase>()

	const fetchEmails = throttleFn(async () => {
		const response = await fetch("/api/emails")
		const data: Array<EmailBase> = await response.json()
		setEmails(data)
	})

	const emailsWebsocket = new EmailEventsWebsocket(fetchEmails)

	createEffect(() => {
		fetchEmails()
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

	return (
		<SelectedEmailContext.Provider value={selectedEmailContext()}>
			<div h-full w-full overflow-hidden>
				<Switch
					fallback={
						<EmailsList
							emails={() => emails() ?? []}
							loading={() => !emails()}
						/>
					}
				>
					<Match when={emails() !== undefined && selectedEmail() !== undefined}>
						<LayoutWithEmail
							emails={emails}
							selectedEmail={() => selectedEmail()!}
						/>
					</Match>
				</Switch>
			</div>
		</SelectedEmailContext.Provider>
	)
}

interface LayoutWithEmailProps extends EmailsListProps {
	selectedEmail: Accessor<EmailBase>
}

function LayoutWithEmail({ emails }: LayoutWithEmailProps) {
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
				<EmailsList emails={emails} />
			</div>
			<div h-screen overflow-y-auto self-stretch flex-1>
				<Email />
			</div>
		</div>
	)
}
