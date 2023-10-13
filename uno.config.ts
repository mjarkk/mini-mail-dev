import { presetAttributify, presetWind, defineConfig } from "unocss"

export default defineConfig({
	presets: [presetWind(), presetAttributify()],
})
