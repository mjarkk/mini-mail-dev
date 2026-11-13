import { Accessor, createEffect, createSignal, For, Show } from "solid-js"
import { EmailBase, EmailRemainder } from "../../email"
import { AttachmentModal } from "./AttachmentModal"
import { getUrl } from "../../services/fetch"
import { AttachmentButton } from "./AttachmentButton"

export interface AttachmentsProps {
	email: Accessor<EmailBase>
	emailRemainder: Accessor<EmailRemainder | undefined>
}

export function Attachments({ emailRemainder, email }: AttachmentsProps) {
	const [showAttachment, setShowAttachment] = createSignal<number>()
	const attachments = () => emailRemainder()?.attachments ?? []

	createEffect(() => {
		// When the email changes hide the show attachment
		email()
		setShowAttachment(undefined)
	})

	return (
		<>
			<Show when={attachments().length > 0}>
				<div flex flex-wrap gap-2 min-w-0>
					<For each={attachments()}>
						{(attachment, index) => (
							<AttachmentButton
								attachment={() => attachment}
								onclick={() => setShowAttachment(index)}
							/>
						)}
					</For>
				</div>
			</Show>
			<Show when={showAttachment() !== undefined && attachments().length > 0}>
				<AttachmentModal
					onClose={() => setShowAttachment(undefined)}
					attachmentUrl={() =>
						getUrl(
							"/api/emails/" + email().id + "/attachments/" + showAttachment(),
						)
					}
					attachment={() => attachments()[showAttachment()!]}
				/>
			</Show>
		</>
	)
}
