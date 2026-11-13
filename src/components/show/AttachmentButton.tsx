import { Attachment } from "../../email"
import { AttachmentIcon } from "./AttachmentIcon"

export const isPdf = (t: string) =>
	t === "application/pdf" ||
	t === "application/x-pdf" ||
	t === "application/x-bzpdf"

export const isImage = (t: string) => t.startsWith("image/")

function typeLabel({ contentType, filename }: Attachment): string {
	if (isPdf(contentType)) return "PDF"
	if (isImage(contentType)) {
		return contentType.slice("image/".length).split(/[+;]/)[0]
	}
	const extension = filename.includes(".") ? filename.split(".").pop() : ""
	return extension && extension.length <= 5 ? extension : "File"
}

interface AttachmentButtonProps {
	onclick?: () => void
	attachment: () => Attachment
}

export function AttachmentButton({
	attachment,
	onclick,
}: AttachmentButtonProps) {
	return (
		<button
			type="button"
			onclick={onclick}
			title={attachment().filename}
			inline-flex
			items-center
			gap="[10px]"
			h-10
			max-w-full
			pl="[5px]"
			pr-4
			rounded-full
			bg-card2
			hover:bg-line
			transition-colors
		>
			<span
				shrink-0
				inline-flex
				items-center
				justify-center
				w="[30px]"
				h="[30px]"
				rounded-full
				bg-card
				text-fg
			>
				<AttachmentIcon
					contentType={() => attachment().contentType}
					size={16}
				/>
			</span>
			<span truncate max-w="[220px]" font-mono text-xs text-fg>
				{attachment().filename}
			</span>
			<span class="mono-label" shrink-0>
				{typeLabel(attachment())}
			</span>
		</button>
	)
}
