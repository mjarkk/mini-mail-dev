import { presetAttributify, presetWind3, defineConfig } from "unocss"
import transformerDirectives from "@unocss/transformer-directives"

export default defineConfig({
	presets: [presetWind3(), presetAttributify()],
	transformers: [transformerDirectives()],
	theme: {
		colors: {
			ground: "var(--ground)",
			card: "var(--card)",
			card2: "var(--card2)",
			line: "var(--line)",
			fg: "var(--fg)",
			fg2: "var(--fg2)",
			muted: "var(--muted)",
			signal: "var(--signal)",
		},
	},
	shortcuts: {
		widget: "bg-card rounded-[22px]",
		"mono-label":
			"font-mono text-[11px] font-medium tracking-[0.14em] uppercase text-muted",
		"pill-button":
			"inline-flex items-center justify-center gap-2 h-9 px-4 rounded-full border border-solid font-mono text-[11px] font-medium tracking-[0.14em] uppercase transition-colors",
		"outline-button":
			"pill-button border-fg text-fg hover:bg-fg hover:text-card",
		segment:
			"inline-flex items-center justify-center gap-2 h-[28px] px-[14px] rounded-full font-mono text-[11px] font-medium tracking-[0.14em] uppercase transition-colors",
	},
})
