<script lang="ts">
	import type { Player } from '$lib/types/game';
	import { UIConfig } from '$lib/config/uiConfig';

	interface Props {
		player: Player;
		timeRemaining: number;
		isActive: boolean;
		label?: string;
		onTimeOut?: () => void;
	}

	let { player, isActive, onTimeOut, label = '', timeRemaining: propTimeRemaining }: Props = $props();
	let serverTimeBase = $state(0);
	let serverTimeTimestamp = $state(Date.now());
	let hasTriggeredTimeout = $state(false);
	let tick = $state(0);

	$effect(() => {
		isActive;
		serverTimeBase = propTimeRemaining;
		serverTimeTimestamp = Date.now();
		hasTriggeredTimeout = false;
	});

	const displayTime = $derived(() => {
		tick;
		if (!isActive) return Math.round(serverTimeBase);
		return Math.max(0, Math.round(serverTimeBase - (Date.now() - serverTimeTimestamp) / 1000));
	});

	$effect(() => {
		const current = displayTime();
		if (current <= 0 && isActive && !hasTriggeredTimeout) {
			hasTriggeredTimeout = true;
			onTimeOut?.();
		}
	});

	$effect(() => {
		if (!isActive) return;
		const id = setInterval(() => tick++, UIConfig.timerUpdateIntervalMs);
		return () => clearInterval(id);
	});

	function formatTime(seconds: number): string {
		seconds = Math.max(0, Math.round(seconds));
		return `${Math.floor(seconds / 60)}:${(seconds % 60).toString().padStart(2, '0')}`;
	}

	const isLowTime = $derived(displayTime() < UIConfig.lowTimeThresholdSeconds);
</script>

<div class="w-full max-w-[900px] mx-auto glass-panel rounded-2xl px-3 py-2 flex items-center gap-2 transition {isActive ? '' : 'opacity-70'}">
	<span class="h-3 w-3 rounded-full shrink-0 {player === 'red' ? 'bg-[#db6d63]' : 'bg-emerald-700'}"></span>
	<span class="text-sm font-extrabold {player === 'red' ? 'text-[#b84f48]' : 'text-emerald-800'}">{player === 'red' ? 'Đỏ O' : 'Xanh X'}</span>
	{#if label}<span class="truncate text-xs font-semibold text-emerald-950/60">{label}</span>{/if}
	<span class="ml-auto font-mono text-lg font-black {isLowTime && isActive ? 'text-rose-600 animate-pulse' : 'text-emerald-950'}">{formatTime(displayTime())}</span>
</div>
