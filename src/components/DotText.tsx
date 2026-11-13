import { Index, createMemo, createSignal, onCleanup, onMount } from "solid-js"

const ROWS = 7
const DIGIT_COLUMNS = 5

const GLYPHS: Record<string, Array<string>> = {
	I: ["#", "#", "#", "#", "#", "#", "#"],
	N: ["#...#", "#...#", "##..#", "#.#.#", "#..##", "#...#", "#...#"],
	B: ["####.", "#...#", "#...#", "####.", "#...#", "#...#", "####."],
	O: [".###.", "#...#", "#...#", "#...#", "#...#", "#...#", ".###."],
	X: ["#...#", "#...#", ".#.#.", "..#..", ".#.#.", "#...#", "#...#"],
	"0": [".###.", "#...#", "#...#", "#...#", "#...#", "#...#", ".###."],
	"1": [".#", "##", ".#", ".#", ".#", ".#", ".#"],
	"2": [".###.", "#...#", "....#", "...#.", "..#..", ".#...", "#####"],
	"3": [".###.", "#...#", "....#", "..##.", "....#", "#...#", ".###."],
	"4": ["...#.", "..##.", ".#.#.", "#..#.", "#####", "...#.", "...#."],
	"5": ["#####", "#....", "####.", "....#", "....#", "#...#", ".###."],
	"6": ["..##.", ".#...", "#....", "####.", "#...#", "#...#", ".###."],
	"7": ["#####", "....#", "...#.", "..#..", ".#...", ".#...", ".#..."],
	"8": [".###.", "#...#", "#...#", ".###.", "#...#", "#...#", ".###."],
	"9": [".###.", "#...#", "#...#", ".####", "....#", "...#.", ".##.."],
	":": [".", ".", "#", ".", "#", ".", "."],
	"✉": [
		"#########",
		"##.....##",
		"#.#...#.#",
		"#..#.#..#",
		"#...#...#",
		"#.......#",
		"#########",
	],
}

interface Raster {
	lit: Array<{ col: number; row: number }>
	columns: number
}

function rasterize(text: string, tabular = false): Raster {
	const lit: Raster["lit"] = []
	let col = 0
	for (const char of text) {
		const rows = GLYPHS[char]
		if (!rows) continue

		const glyphColumns = rows[0].length
		const columns =
			tabular && char >= "0" && char <= "9" ? DIGIT_COLUMNS : glyphColumns
		const left = col + Math.floor((columns - glyphColumns) / 2)
		rows.forEach((row, rowIdx) => {
			for (let c = 0; c < glyphColumns; c++) {
				if (row[c] === "#") lit.push({ col: left + c, row: rowIdx })
			}
		})
		col += columns + 1
	}
	return { lit, columns: Math.max(0, col - 1) }
}

function paint(
	fills: Array<string>,
	columns: number,
	raster: Raster,
	col0: number,
	row0: number,
	fill: string,
) {
	for (const { col, row } of raster.lit) {
		const x = col0 + col
		if (x >= 0 && x < columns) fills[(row0 + row) * columns + x] = fill
	}
}

const pitchFor = (dot: number) => dot + Math.max(1, Math.round(dot / 3))

interface DotsProps {
	fills: Array<string>
	columns: number
	pitchX: number
	pitchY: number
	dot: number
}

function Dots(props: DotsProps) {
	return (
		// Index keeps one circle per grid cell, so a cell's fill change can transition.
		<Index each={props.fills}>
			{(fill, i) => (
				<circle
					cx={(i % props.columns) * props.pitchX + props.dot / 2}
					cy={Math.floor(i / props.columns) * props.pitchY + props.dot / 2}
					r={props.dot / 2}
					style={{ fill: fill() }}
				/>
			)}
		</Index>
	)
}

export interface DotTextProps {
	/** Only characters in GLYPHS render; anything else is skipped. */
	text: string
	/** Dot diameter in px. */
	dot: number
	label?: string
	unlit?: string
	class?: string
	/** Every digit gets the same width, so a changing number keeps its size. */
	tabular?: boolean
}

/** Lit dots take the current text color. */
export function DotText(props: DotTextProps) {
	const layout = createMemo(() => {
		const dot = props.dot
		const pitch = pitchFor(dot)
		const raster = rasterize(props.text, props.tabular)
		const fills = Array<string>(raster.columns * ROWS).fill(
			props.unlit ?? "var(--unlit)",
		)
		paint(fills, raster.columns, raster, 0, 0, "currentColor")

		return {
			fills,
			columns: raster.columns,
			pitch,
			width: Math.max(0, (raster.columns - 1) * pitch + dot),
			height: (ROWS - 1) * pitch + dot,
		}
	})

	return (
		<svg
			class={`dot-matrix ${props.class ?? ""}`}
			role="img"
			aria-label={props.label ?? props.text}
			width={layout().width}
			height={layout().height}
			viewBox={`0 0 ${layout().width} ${layout().height}`}
			style={{ display: "block", "flex-shrink": 0 }}
		>
			<Dots
				fills={layout().fills}
				columns={layout().columns}
				pitchX={layout().pitch}
				pitchY={layout().pitch}
				dot={props.dot}
			/>
		</svg>
	)
}

const WORD_GAP = 2

export interface DotMatrixProps {
	/** Pinned to the first column. */
	start: string
	/** Pinned to the last column; drops to a band below `start` when both don't fit on one. */
	end: string
	/** Dot diameter in px. */
	dot: number
	label: string
	/**
	 * CSS color for `end`'s lit dots; defaults to the current text color.
	 * A color rather than a class: a dot that swaps class mid-fade fades toward the old currentColor.
	 */
	endColor?: string
}

/** A dot grid as wide as its container. */
export function DotMatrix(props: DotMatrixProps) {
	let svg!: SVGSVGElement
	const [width, setWidth] = createSignal(0)

	onMount(() => {
		const observer = new ResizeObserver(([entry]) =>
			setWidth(entry.contentRect.width),
		)
		observer.observe(svg)
		onCleanup(() => observer.disconnect())
	})

	const layout = createMemo(() => {
		const dot = props.dot
		const pitch = pitchFor(dot)
		const columns = Math.max(1, Math.floor((width() - dot) / pitch) + 1)
		// Stretched slightly so the last column ends flush with the container edge.
		const pitchX = columns > 1 ? (width() - dot) / (columns - 1) : pitch

		const start = rasterize(props.start)
		const end = rasterize(props.end)
		const wrapped = start.columns + WORD_GAP + end.columns > columns
		const rows = wrapped ? 2 * ROWS + 1 : ROWS

		const fills = Array<string>(columns * rows).fill("var(--unlit)")
		paint(fills, columns, start, 0, 0, "currentColor")
		paint(
			fills,
			columns,
			end,
			columns - end.columns,
			wrapped ? ROWS + 1 : 0,
			props.endColor ?? "currentColor",
		)

		return { fills, columns, pitch, pitchX, height: (rows - 1) * pitch + dot }
	})

	return (
		<svg
			ref={svg}
			class="dot-matrix"
			role="img"
			aria-label={props.label}
			width="100%"
			height={layout().height}
			style={{ display: "block" }}
		>
			<Dots
				fills={layout().fills}
				columns={layout().columns}
				pitchX={layout().pitchX}
				pitchY={layout().pitch}
				dot={props.dot}
			/>
		</svg>
	)
}
