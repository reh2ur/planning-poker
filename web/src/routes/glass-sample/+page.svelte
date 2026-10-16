<script lang="ts">
	import { onMount } from 'svelte';

	const hand = ['0', '1', '2', '3', '5', '8', '13', '21', '?', 'coffee'];
	const table = [
		{ name: 'Alex', vote: '5' },
		{ name: 'Sam', vote: '8' },
		{ name: 'Jordan', vote: '5' },
		{ name: 'Robin', vote: '5' },
		{ name: 'Casey', vote: '3' },
		{ name: 'Morgan', vote: '?' }
	];
	const typed = '1';

	let dark = $state(false);

	onMount(() => {
		dark = matchMedia('(prefers-color-scheme: dark)').matches;
	});

	$effect(() => {
		document.documentElement.classList.toggle('dark', dark);
	});
</script>

<svelte:head>
	<title>Glass sample</title>
</svelte:head>

<div
	class="relative min-h-dvh overflow-hidden bg-linear-to-br from-sky-300 via-fuchsia-200 to-amber-200 p-4 text-zinc-950 sm:p-8 dark:from-indigo-950 dark:via-slate-900 dark:to-fuchsia-950 dark:text-white"
>
	<div class="absolute -top-24 -left-24 size-96 rounded-full bg-fuchsia-500/50 blur-3xl"></div>
	<div class="absolute -right-24 bottom-0 size-96 rounded-full bg-cyan-400/50 blur-3xl"></div>

	<main
		class="relative mx-auto flex max-w-xl flex-col gap-6 rounded-3xl border border-white/40 bg-white/30 p-5 shadow-2xl backdrop-blur-xl sm:p-8 dark:border-white/15 dark:bg-white/10"
	>
		<header class="flex items-center justify-between gap-4">
			<div>
				<p class="text-sm font-bold tracking-widest uppercase opacity-80">team-checkout</p>
				<h1 class="text-3xl font-black tracking-tight">Checkout redesign</h1>
			</div>
			<button
				type="button"
				onclick={() => (dark = !dark)}
				class="rounded-full border-2 border-zinc-950 bg-white px-4 py-2 text-sm font-bold text-zinc-950 dark:border-white dark:bg-zinc-950 dark:text-white"
			>
				{dark ? 'Light' : 'Dark'}
			</button>
		</header>

		<section aria-label="Table" class="flex flex-col gap-3">
			<div class="flex items-baseline justify-between">
				<h2 class="text-xl font-black">Revealed</h2>
				<p class="text-sm font-bold opacity-80">6 of 6 voted · 2 watching</p>
			</div>
			<ul class="grid grid-cols-3 gap-3">
				{#each table as seat (seat.name)}
					<li class="flex flex-col items-center gap-2">
						<div
							class="grid aspect-3/4 w-full place-items-center rounded-2xl border-4 border-zinc-950 bg-white text-5xl font-black text-zinc-950 dark:border-white dark:bg-zinc-950 dark:text-white"
						>
							{seat.vote}
						</div>
						<span
							class="rounded-lg bg-zinc-950 px-3 py-1 text-base font-bold text-white dark:bg-white dark:text-zinc-950"
						>
							{seat.name}
						</span>
					</li>
				{/each}
			</ul>
		</section>

		<section
			aria-label="Result"
			class="flex flex-wrap items-center justify-between gap-x-4 gap-y-1 rounded-2xl border-4 border-zinc-950 bg-white px-4 py-3 text-zinc-950 dark:border-white dark:bg-zinc-950 dark:text-white"
		>
			<p class="text-2xl font-black">3×5 · 1×8 · 1×3</p>
			<p class="text-base font-bold">Spread: Casey 3 ↔ Sam 8</p>
		</section>

		<section aria-label="Your hand" class="flex flex-col gap-3">
			<p class="text-sm font-bold opacity-80">
				Type a value, Enter commits: <kbd
					class="rounded border-2 border-zinc-950 bg-white px-1.5 font-mono text-zinc-950 dark:border-white dark:bg-zinc-950 dark:text-white"
					>{typed}</kbd
				>
			</p>
			<ul class="grid grid-cols-5 gap-2">
				{#each hand as card (card)}
					{@const highlighted = card.startsWith(typed)}
					<li
						class="grid aspect-3/4 place-items-center rounded-xl border-4 text-2xl font-black {highlighted
							? 'border-amber-500 bg-amber-300 text-zinc-950 ring-4 ring-zinc-950 dark:ring-white'
							: 'border-zinc-950 bg-white text-zinc-950 dark:border-white dark:bg-zinc-950 dark:text-white'}"
					>
						{card === 'coffee' ? '☕' : card}
					</li>
				{/each}
			</ul>
		</section>
	</main>
</div>
