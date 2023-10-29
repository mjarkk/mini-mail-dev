/* @refresh reload */
import "@unocss/reset/sanitize/sanitize.css"
import "@unocss/reset/sanitize/assets.css"
import "@unocss/reset/sanitize/forms.css"
import "@unocss/reset/sanitize/system-ui.css"
import "@unocss/reset/sanitize/ui-monospace.css"
import "@unocss/reset/sanitize/typography.css"
import "./global.css"

import "virtual:uno.css"

import { render } from "solid-js/web"

import { App } from "./components/App"

if ((window as any).trustedTypes && (window as any).trustedTypes.createPolicy) {
	;(window as any).trustedTypes.createPolicy("default", {
		createHTML: (string: any, sink: any) => {
			if (sink === "TrustedHTML") return string
			throw new Error("HTML is not trusted")
		},
	})
}

const root = document.getElementById("root")

render(() => <App />, root!)
