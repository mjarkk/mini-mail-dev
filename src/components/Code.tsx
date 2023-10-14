import hljs from "highlight.js/lib/core"
import xml from "highlight.js/lib/languages/xml"
import { Accessor, createMemo } from "solid-js"
import "highlight.js/styles/github-dark.css"
import "./Code.css"

export interface CodeProps {
	lang: "xml"
	code: Accessor<string>
}

export function Code({ lang, code }: CodeProps) {
	hljs.registerLanguage("xml", xml)

	const highlighted = createMemo(
		() =>
			`<code class="hljs">${
				hljs.highlight(code(), { language: lang }).value
			}</code>`,
	)

	return <pre innerHTML={highlighted()} />
}
