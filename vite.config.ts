import { defineConfig } from "vite"
import solid from "vite-plugin-solid"
import unocss from "unocss/vite"
import Icons from "unplugin-icons/vite"

export default defineConfig({
	plugins: [Icons({ compiler: "solid" }), unocss(), solid()],
	server: {
		port: 3001,
	},
})
