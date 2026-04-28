<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchDNSStats } from '$lib/api/client';
	import type { DNSStats } from '$lib/types';
	import { formatNumber } from '$lib/utils';

	let stats: DNSStats | null = $state(null);
	let loading = $state(true);
	let error = $state('');

	onMount(async () => {
		try {
			stats = await fetchDNSStats();
		} catch (e) {
			error = e instanceof Error ? e.message : 'DNS unavailable';
		}
		loading = false;
	});
</script>

<svelte:head><title>DNS Gateway — Viyoga</title></svelte:head>

<div class="p-6 space-y-6">
	<h1 class="text-2xl font-bold">DNS Gateway</h1>
	<p class="text-base-content/60 text-sm">Network-level ad/tracker blocking</p>

	{#if error}
		<div class="alert alert-warning">
			<span>DNS Gateway Unavailable — configure Technitium in viyoga.yaml</span>
		</div>
	{:else if loading}
		<div class="flex justify-center py-20"><span class="loading loading-spinner loading-lg text-primary"></span></div>
	{:else if stats}
		<div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
			<div class="card bg-base-200 border border-base-300/50 p-4">
				<p class="text-xs text-base-content/60 uppercase">Total Queries</p>
				<p class="text-3xl font-bold text-primary">{formatNumber(stats.total_queries)}</p>
			</div>
			<div class="card bg-base-200 border border-base-300/50 p-4">
				<p class="text-xs text-base-content/60 uppercase">Blocked</p>
				<p class="text-3xl font-bold text-error">{formatNumber(stats.total_blocked_queries)}</p>
			</div>
			<div class="card bg-base-200 border border-base-300/50 p-4">
				<p class="text-xs text-base-content/60 uppercase">Block Rate</p>
				<p class="text-3xl font-bold text-warning">{stats.blocked_percent.toFixed(1)}%</p>
			</div>
			<div class="card bg-base-200 border border-base-300/50 p-4">
				<p class="text-xs text-base-content/60 uppercase">Clients</p>
				<p class="text-3xl font-bold text-accent">{stats.total_clients}</p>
			</div>
		</div>
	{/if}
</div>
