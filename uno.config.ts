import { presetAttributify, presetWind, defineConfig } from "unocss"
import transformerDirectives from "@unocss/transformer-directives"

export default defineConfig({
	presets: [presetWind(), presetAttributify()],
	transformers: [transformerDirectives()],
})
