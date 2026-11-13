import { Match, Portal, Switch } from "solid-js/web"
import { Attachment } from "../../email"
import { isImage, isPdf } from "./AttachmentButton"
import { Show, createEffect, onCleanup, onMount } from "solid-js"
import { createSignal } from "solid-js"
import { AttachmentIcon } from "./AttachmentIcon"

export interface AttachmentModalProps {
	onClose: () => void
	attachmentUrl: () => string
	attachment: () => Attachment
}

export function AttachmentModal({
	onClose,
	attachmentUrl,
	attachment,
}: AttachmentModalProps) {
	let modalRef: HTMLDivElement | undefined
	const modalClick = (e: MouseEvent) => {
		if (
			!modalRef ||
			modalRef === e.target ||
			modalRef === (e.target as unknown as HTMLElement).children?.[0]
		) {
			onClose()
		}
	}

	const downloadAttachment = () => {
		const link = document.createElement("a")
		link.download = attachment().filename
		link.href = attachmentUrl()
		link.click()
	}

	const keyboardEvent = (e: KeyboardEvent) => {
		if (e.key === "Escape") {
			onClose()
		}
	}

	onMount(() => {
		document.addEventListener("keydown", keyboardEvent)
	})
	onCleanup(() => {
		document.removeEventListener("keydown", keyboardEvent)
	})

	return (
		<Portal mount={document.body}>
			<div
				ref={modalRef}
				fixed
				z-100
				top-0
				left-0
				h-screen
				w-screen
				box-border
				p-4
				sm:p-10
				bg="black/60"
				onClick={modalClick}
			>
				<div
					role="dialog"
					aria-modal="true"
					aria-label={attachment().filename}
					class="widget"
					mx-auto
					max-w="800px"
					h-full
					overflow-hidden
					flex
					flex-col
				>
					<div px="[22px]" py="[18px]" flex items-center gap-3>
						<span class="mono-label">Attachment</span>
						<p m-0 truncate font-mono text-sm text-fg>
							{attachment().filename}
						</p>
					</div>
					<div flex-1 min-h-0 mx="[10px]" rounded="[18px]" bg-card2 overflow-hidden flex>
						<ShowAttachment
							attachment={attachment}
							attachmentUrl={attachmentUrl}
						/>
					</div>
					<div px="[22px]" py-4 flex justify-end gap-2 items-center>
						<button type="button" class="outline-button" onClick={onClose}>
							Close
						</button>
						<button
							type="button"
							class="outline-button"
							onClick={downloadAttachment}
						>
							Download
						</button>
					</div>
				</div>
			</div>
		</Portal>
	)
}

interface ShowAttachmentProps {
	attachment: () => Attachment
	attachmentUrl: () => string
}

function ShowAttachment({ attachment, attachmentUrl }: ShowAttachmentProps) {
	const [fileUrl, setFileUrl] = createSignal<string | undefined>(undefined)

	const getData = async () => {
		const response = await fetch(attachmentUrl())
		const blob = await response.blob()
		const url = URL.createObjectURL(blob)
		setFileUrl(url)
	}

	createEffect(() => {
		getData()
	})

	return (
		<Show when={fileUrl()}>
			<Switch
				fallback={
					<ContentNotDisplayable contentType={() => attachment().contentType} />
				}
			>
				<Match when={isImage(attachment().contentType)}>
					<div
						h-full
						w-full
						style={{ "background-image": `url(${fileUrl()})` }}
						bg-center
						bg-contain
						bg-no-repeat
					/>
				</Match>
				<Match when={isPdf(attachment().contentType)}>
					<iframe
						title={attachment().filename}
						h-full
						w-full
						border-0
						src={fileUrl()}
					/>
				</Match>
			</Switch>
		</Show>
	)
}

interface ContentNotDisplayableProps {
	contentType: () => string
}

function ContentNotDisplayable({ contentType }: ContentNotDisplayableProps) {
	return (
		<p
			class="mono-label"
			m-0
			flex
			flex-1
			flex-col
			justify-center
			items-center
		>
			<AttachmentIcon contentType={contentType} />
			<span mt-3>Content not displayable</span>
		</p>
	)
}
