import { Accessor, createEffect, createSignal, on, onCleanup } from "solid-js"

const DURATION_MS = 600

/** Eases through the integers to each new `target`; jumps straight there under reduced motion. */
export function createTween(target: Accessor<number>): Accessor<number> {
	const [shown, setShown] = createSignal(target())
	let frame = 0

	createEffect(
		on(
			target,
			(to) => {
				cancelAnimationFrame(frame)
				const from = shown()
				if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) {
					setShown(to)
					return
				}

				const begin = performance.now()
				const tick = (now: number) => {
					const t = Math.min(1, (now - begin) / DURATION_MS)
					const eased = 1 - (1 - t) ** 3
					setShown(Math.round(from + (to - from) * eased))
					if (t < 1) frame = requestAnimationFrame(tick)
				}
				frame = requestAnimationFrame(tick)
			},
			{ defer: true },
		),
	)
	onCleanup(() => cancelAnimationFrame(frame))

	return shown
}
