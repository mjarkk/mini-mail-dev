import { Attachment } from "../email"
import IconImage from "~icons/material-symbols/image-outline"
import IconFilePresent from "~icons/material-symbols/file-present-outline"
import IconPictureAsPdf from "~icons/material-symbols/picture-as-pdf-outline"
import { Match, Switch } from "solid-js"

export const isPdf = (t: string) =>
	t === "application/pdf" ||
	t === "application/x-pdf" ||
	t === "application/x-bzpdf"

export const isImage = (t: string) => t.startsWith("image/")

interface AttachmentButtonProps extends Attachment {
	onclick?: () => void
}

export function AttachmentButton({
	contentType,
	filename,
	onclick,
}: AttachmentButtonProps) {
	const iconStyle = { height: "30px", width: "30px" }
	return (
		<button onclick={onclick} p-4 border-0 w="180px">
			<Switch fallback={<IconFilePresent style={iconStyle} />}>
				<Match when={isImage(contentType)}>
					<IconImage style={iconStyle} />
				</Match>
				<Match when={isPdf(contentType)}>
					<IconPictureAsPdf style={iconStyle} />
				</Match>
			</Switch>
			<div truncate>{filename}</div>
		</button>
	)
}
