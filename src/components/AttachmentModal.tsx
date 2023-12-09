import { Match, Portal, Switch } from "solid-js/web"
import { Attachment } from "../email"
import { isImage, isPdf } from "./AttachmentButton"

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
				bg="black/40"
				onClick={modalClick}
			>
				<div mx-auto bg-white rounded max-w="800px" my-10 text-black shadow-xl>
					<div overflow-hidden rounded>
						<div p-6 flex justify-between items-center>
							<p m-0>{attachment().filename}</p>
							<button onClick={downloadAttachment} bg="white hover:zinc-300">Download</button>
						</div>
						<ShowAttachment
							attachment={attachment}
							attachmentUrl={attachmentUrl}
						/>
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
	return (
		<Switch fallback={<p>Geen preview mogelijk</p>}>
			<Match when={isImage(attachment().contentType)}>
				<img src={attachmentUrl()} />
			</Match>
			<Match when={isPdf(attachment().contentType)}>
				<iframe src={attachmentUrl()} />
			</Match>
		</Switch>
	)
}
