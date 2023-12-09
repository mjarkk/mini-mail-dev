import { Attachment } from "../email"
import { AttachmentIcon } from "./AttachmentIcon"

export const isPdf = (t: string) =>
	t === "application/pdf" ||
	t === "application/x-pdf" ||
	t === "application/x-bzpdf"

export const isImage = (t: string) => t.startsWith("image/")

interface AttachmentButtonProps {
	onclick?: () => void
	attachment: () => Attachment
}

export function AttachmentButton({
	attachment,
	onclick,
}: AttachmentButtonProps) {
	return (
		<button onclick={onclick} p-4 border-0 w="180px">
			<AttachmentIcon contentType={() => attachment().contentType} />
			<div truncate>{attachment().filename}</div>
		</button>
	)
}
