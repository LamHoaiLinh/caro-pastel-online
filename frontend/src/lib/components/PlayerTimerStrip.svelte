<script lang="ts">
	import type { Player } from '$lib/types/game';
	import { UIConfig } from '$lib/config/uiConfig';

	interface Props {
		player: Player;
		timeRemaining: number;
		turnTimeRemaining?: number;
		moveTimeLimit?: number;
		isActive: boolean;
		label?: string;
		onTimeOut?: () => void;
	}

	let {
		player,
		isActive,
		onTimeOut,
		label = '',
		timeRemaining: propTimeRemaining,
		turnTimeRemaining: propTurnTimeRemaining = 0,
		moveTimeLimit = 0
	}: Props = $props();
	let serverTimeBase = $state(0);
	let turnTimeBase = $state(0);
	let serverTimeTimestamp = $state(Date.now());
	let hasTriggeredTimeout = $state(false);
	let tick = $state(0);

	$effect(() => {
		isActive;
		serverTimeBase = propTimeRemaining;
		turnTimeBase = propTurnTimeRemaining;
		serverTimeTimestamp = Date.now();
		hasTriggeredTimeout = false;
	});

	const displayTime = $derived(() => {
		tick;
		if (!isActive) return Math.round(serverTimeBase);
		return Math.max(0, Math.round(serverTimeBase - (Date.now() - serverTimeTimestamp) / 1000));
	});

	const displayTurnTime = $derived(() => {
		tick;
		if (moveTimeLimit <= 0) return 0;
		if (!isActive) return Math.max(0, Math.round(moveTimeLimit));
		return Math.max(0, Math.round(turnTimeBase - (Date.now() - serverTimeTimestamp) / 1000));
	});

	$effect(() => {
		const totalLeft = displayTime();
		const turnLeft = displayTurnTime();
		const turnExpired = moveTimeLimit > 0 && turnLeft <= 0;
		if ((totalLeft <= 0 || turnExpired) && isActive && !hasTriggeredTimeout) {
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
	const isLowTurnTime = $derived(moveTimeLimit > 0 && displayTurnTime() <= Math.min(10, Math.ceil(moveTimeLimit / 3)));
</script>

<div class="w-full max-w-[900px] mx-auto glass-panel rounded-2xl px-3 py-2 flex items-center gap-2 transition {isActive ? '' : 'opacity-70'}">
	<span class="h-3 w-3 rounded-full shrink-0 {player === 'red' ? 'bg-[#db6d63]' : 'bg-emerald-700'}"></span>
	<div class="min-w-0">
		<span class="text-sm font-extrabold {player === 'red' ? 'text-[#b84f48]' : 'text-emerald-800'}">{player === 'red' ? 'Đỏ O' : 'Xanh X'}</span>
		{#if label}<span class="ml-1.5 text-xs font-semibold text-emerald-950/60">{label}</span>{/if}
	</div>
	<div class="ml-auto flex items-center gap-2 sm:gap-3 text-right shrink-0">
		{#if moveTimeLimit > 0}
			<div class="rounded-lg border border-amber-200 bg-amber-50/90 px-2 py-1 leading-none">
				<span class="block text-[9px] sm:text-[10px] font-bold uppercase tracking-wide text-amber-700">{isActive ? 'Nước này' : 'Tối đa/lượt'}</span>
				<span class="mt-1 block font-mono text-sm sm:text-base font-black {isLowTurnTime && isActive ? 'text-rose-600 animate-pulse' : 'text-amber-950'}">{formatTime(displayTurnTime())}</span>
			</div>
		{/if}
		<div class="leading-none">
			<span class="block text-[9px] sm:text-[10px] font-bold uppercase tracking-wide text-emerald-700">Tổng giờ</span>
			<span class="mt-1 block font-mono text-lg font-black {isLowTime && isActive ? 'text-rose-600 animate-pulse' : 'text-emerald-950'}">{formatTime(displayTime())}</span>
		</div>
	</div>
</div>
