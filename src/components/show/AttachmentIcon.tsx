import IconImage from "~icons/material-symbols/image-outline"
import IconFilePresent from "~icons/material-symbols/file-present-outline"
import IconPictureAsPdf from "~icons/material-symbols/picture-as-pdf-outline"
import { Match, Switch } from "solid-js"
import { isImage, isPdf } from "./AttachmentButton"

export interface AttachmentIconProps {
	contentType: () => string
}

export function AttachmentIcon({ contentType }: AttachmentIconProps) {
	const iconStyle = { height: "30px", width: "30px" }

	return (
		<Switch fallback={<IconFilePresent style={iconStyle} />}>
			<Match when={isImage(contentType())}>
				<IconImage style={iconStyle} />
			</Match>
			<Match when={isPdf(contentType())}>
				<IconPictureAsPdf style={iconStyle} />
			</Match>
		</Switch>
	)
}
