import { Match, Portal, Switch } from "solid-js/web"
import { Attachment } from "../email"
import { isImage, isPdf } from "./AttachmentButton"
import { Show, createEffect } from "solid-js"
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
				p-10
				bg="black/40"
				onClick={modalClick}
			>
				<div
					mx-auto
					bg-zinc-900
					rounded
					max-w="800px"
					h-full
					shadow-xl
					overflow-hidden
					flex
					flex-col
				>
					<div
						p-6
						flex
						justify-between
						items-center
						border-0
						border-zinc-800
						border-b
						border-b-solid
					>
						<p m-0>{attachment().filename}</p>
						<button onClick={downloadAttachment} border-zinc-500>
							Download
						</button>
					</div>
					<ShowAttachment
						attachment={attachment}
						attachmentUrl={attachmentUrl}
					/>
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
					<iframe src={fileUrl()} />
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
		<p flex flex-col justify-center items-center h-full text-zinc-400>
			<AttachmentIcon contentType={contentType} />
			<span mt-3>Content not displayable</span>
		</p>
	)
}
