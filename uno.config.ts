import { presetAttributify, presetWind3, defineConfig } from "unocss"
import transformerDirectives from "@unocss/transformer-directives"

export default defineConfig({
	presets: [presetWind3(), presetAttributify()],
	transformers: [transformerDirectives()],
})
