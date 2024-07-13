import { Accessor, createSignal, For, Show } from "solid-js"
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

	const hasAttachments = () => (emailRemainder()?.attachments?.length ?? 0) > 0

	return (
		<>
			<Show when={hasAttachments()}>
				<div border-0 border-b border-b-solid border-zinc-800 p-4>
					<p m-0 text-zinc-500>
						Attachments:
					</p>
					<div flex flex-wrap gap-2 mt-1>
						<For each={emailRemainder()!.attachments}>
							{(attachment, index) => (
								<AttachmentButton
									attachment={() => attachment}
									onclick={() => setShowAttachment(index)}
								/>
							)}
						</For>
					</div>
				</div>
			</Show>
			<Show when={showAttachment() !== undefined}>
				<AttachmentModal
					onClose={() => setShowAttachment(undefined)}
					attachmentUrl={() =>
						getUrl(
							"/api/emails/" + email().id + "/attachments/" + showAttachment(),
						)
					}
					attachment={() => emailRemainder()!.attachments![showAttachment()!]}
				/>
			</Show>
		</>
	)
}
