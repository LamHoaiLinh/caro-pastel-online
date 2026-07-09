<script lang="ts">
	import { onMount } from 'svelte';
	import type { Cell } from '$lib/types/game';
	import CellComponent from './Cell.svelte';
	import WinningLine from './WinningLine.svelte';
	import { calculateGhostStonePosition, isValidCell, computeCellSize } from '$lib/utils/boardUtils';
	import { vibrateOnValidMove, vibrateOnInvalidMove } from '$lib/utils/haptics';
	import { GameConfig } from '$lib/config/gameConfig';

	interface Props {
		board: Cell[];
		onMove: (x: number, y: number) => void;
		winningLine?: Array<{ x: number; y: number }>;
		lastMove?: { x: number; y: number } | null;
		interactive?: boolean;
	}

	let {
		board,
		onMove,
		winningLine = [],
		lastMove = null,
		interactive = true
	}: Props = $props();

	let ghostPosition = $state<{ x: number; y: number } | null>(null);
	let cellSize = $state(computeCellSize(typeof window !== 'undefined' ? window.innerWidth : 1024));
	let boardEl: HTMLDivElement | undefined = $state();

	const labelSize = $derived(Math.max(cellSize * 0.52, 12));
	const labelFont = $derived(`${Math.max(labelSize * 0.68, 9)}px`);
	const cols = $derived(Array.from({ length: GameConfig.boardSize }, (_, i) => String.fromCharCode(65 + i)));
	const rows = $derived(Array.from({ length: GameConfig.boardSize }, (_, i) => i + 1));

	function handleCellClick(x: number, y: number) {
		if (!interactive) {
			vibrateOnInvalidMove();
			return;
		}
		const cell = board[y * GameConfig.boardSize + x];
		if (!cell || cell.player !== 'none') {
			vibrateOnInvalidMove();
			return;
		}
		vibrateOnValidMove();
		onMove(x, y);
	}

	function handleTouchMove(event: TouchEvent) {
		if (!interactive) return;
		const touch = event.touches[0];
		const element = document.elementFromPoint(touch.clientX, touch.clientY);
		if (element instanceof HTMLElement) {
			const x = parseInt(element.dataset.x ?? '-1');
			const y = parseInt(element.dataset.y ?? '-1');
			if (isValidCell(x, y)) {
				const rect = element.getBoundingClientRect();
				ghostPosition = calculateGhostStonePosition(rect.left + rect.width / 2, rect.top + rect.height / 2, cellSize * 0.78);
			}
		}
	}

	onMount(() => {
		const observer = new ResizeObserver((entries) => {
			for (const entry of entries) {
				const width = entry.contentRect.width;
				if (width > 0) cellSize = computeCellSize(width);
			}
		});
		if (boardEl) observer.observe(boardEl);
		return () => observer.disconnect();
	});
</script>

<div class="w-full max-w-[900px] mx-auto overflow-x-auto pb-1" bind:this={boardEl} ontouchmove={handleTouchMove}>
	<div class="relative inline-block min-w-max left-1/2 -translate-x-1/2">
		<div
			class="grid gap-0 touch-none select-none rounded-2xl overflow-hidden border border-emerald-300/70 shadow-[0_16px_38px_rgba(30,101,70,0.18)]"
			style="display: grid; grid-template-columns: {labelSize}px repeat({GameConfig.boardSize}, {cellSize}px) {labelSize}px; grid-template-rows: {labelSize}px repeat({GameConfig.boardSize}, {cellSize}px) {labelSize}px;"
		>
			<div class="label-cell"></div>
			{#each cols as col}
				<div class="label-cell" style="font-size: {labelFont};">{col}</div>
			{/each}
			<div class="label-cell"></div>

			{#each rows as row, y}
				<div class="label-cell" style="font-size: {labelFont};">{row}</div>
				{#each cols as _, x}
					{@const cell = board[y * GameConfig.boardSize + x]}
					<CellComponent
						x={x}
						y={y}
						player={cell.player}
						isLastMove={lastMove !== null && x === lastMove.x && y === lastMove.y}
						{cellSize}
						{interactive}
						onclick={() => handleCellClick(x, y)}
						onkeydown={(e) => e.key === 'Enter' && handleCellClick(x, y)}
					/>
				{/each}
				<div class="label-cell" style="font-size: {labelFont};">{row}</div>
			{/each}

			<div class="label-cell"></div>
			{#each cols as col}
				<div class="label-cell" style="font-size: {labelFont};">{col}</div>
			{/each}
			<div class="label-cell"></div>
		</div>

		<WinningLine winningLine={winningLine} boardSize={GameConfig.boardSize} {cellSize} {labelSize} />

		{#if ghostPosition}
			<div class="fixed pointer-events-none rounded-full border-2 border-dashed border-emerald-700 opacity-60" style="width: {cellSize}px; height: {cellSize}px; left: {ghostPosition.x - cellSize / 2}px; top: {ghostPosition.y - cellSize / 2}px;"></div>
		{/if}
	</div>
</div>

<style>
	.label-cell {
		display: flex;
		align-items: center;
		justify-content: center;
		background: rgba(217, 242, 226, 0.96);
		color: #347358;
		font-weight: 800;
		font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
	}
</style>
