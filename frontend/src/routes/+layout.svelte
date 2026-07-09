<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import '../app.pcss';
	import { base } from '$app/paths';
	import { onMount } from 'svelte';
	import BackgroundPicker from '$lib/components/BackgroundPicker.svelte';
	import SoundToggle from '$lib/components/SoundToggle.svelte';

	let { children } = $props();
	let backgroundId = $state('workspace');

	const files: Record<string, string> = {
		'workspace': 'workspace.webp',
		'mint-family': 'mint-family.webp',
		'pastel-garden': 'pastel-garden.webp',
		'cat-tulips': 'cat-tulips.webp'
	};

	onMount(() => {
		const saved = localStorage.getItem('caro-pastel-background');
		if (saved && files[saved]) backgroundId = saved;
	});

	function selectBackground(id: string) {
		if (!files[id]) return;
		backgroundId = id;
		localStorage.setItem('caro-pastel-background', id);
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>Caro Pastel Online</title>
	<meta name="description" content="Cờ Caro pastel: chơi với AI, hai người cùng máy hoặc tạo phòng online." />
</svelte:head>

<div
	class="fixed inset-0 -z-20 bg-cover bg-center bg-no-repeat"
	style={`background-image: url('${base}/backgrounds/${files[backgroundId]}')`}
></div>
<div class="fixed inset-0 -z-10 bg-gradient-to-b from-emerald-950/20 via-emerald-900/5 to-emerald-950/20"></div>

<nav class="sticky top-0 z-40 border-b border-emerald-200/70 bg-white/82 backdrop-blur-xl">
	<div class="max-w-7xl mx-auto px-3 py-2 sm:px-5 sm:py-3 flex items-center justify-between gap-3">
		<a href={`${base}/`} class="font-extrabold tracking-tight text-emerald-900 text-lg sm:text-xl">
			Caro Pastel
		</a>
		<div class="flex items-center gap-2">
			<a href={`${base}/game?mode=ai`} class="hidden sm:inline-flex soft-button rounded-xl px-3 py-2 text-sm font-semibold">Chơi với AI</a>
			<a href={`${base}/game?mode=online`} class="hidden sm:inline-flex mint-button rounded-xl px-3 py-2 text-sm font-semibold">Chơi online</a>
			<BackgroundPicker selected={backgroundId} onSelect={selectBackground} />
			<SoundToggle />
		</div>
	</div>
</nav>

<main class="min-h-[calc(100vh-64px)]">
	{@render children()}
</main>
