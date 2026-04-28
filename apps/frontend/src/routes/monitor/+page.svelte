<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchMonitorStatus } from '$lib/api/client';
	import type { MonitorTargetStatus } from '$lib/types';

	let statuses: MonitorTargetStatus[] = $state([]);
	let loading = $state(true);
	let error = $state('');

	onMount(async () => {
		try { statuses = await fetchMonitorStatus(); }
		catch (e) { error = e instanceof Error ? e.message : 'Monitor unavailable'; }
		loading = false;
	});
</script>

<svelte:head><title>Web Monitor — Viyoga</title></svelte:head>

<div class="p-6 space-y-6">
	<h1 class="text-2xl font-bold">Web Service Monitor</h1>
	<p class="text-base-content/60 text-sm">HTTP endpoint health checks</p>

	{#if error}
		<div class="alert alert-info"><span>{error}</span></div>
	{:else if loading}
		<div class="flex justify-center py-20"><span class="loading loading-spinner loading-lg text-primary"></span></div>
	{:else if statuses.length === 0}
		<div class="card bg-base-200 border border-base-300/50 p-8 text-center">
			<p class="text-base-content/60">No monitoring targets configured.</p>
			<p class="text-xs text-base-content/40 mt-2">Add targets via the API: POST /api/v1/monitor/targets</p>
		</div>
	{:else}
		<div class="grid gap-4">
			{#each statuses as s (s.target.id)}
				<div class="card bg-base-200 border border-base-300/50 hover:border-primary/30 transition-colors">
					<div class="card-body p-4">
						<div class="flex justify-between items-center">
							<div class="flex items-center gap-3">
								<span class="text-2xl">{s.last_result?.is_up ? '🟢' : '🔴'}</span>
								<div>
									<h3 class="font-bold">{s.target.name}</h3>
									<p class="text-xs font-mono text-base-content/50">{s.target.url}</p>
								</div>
							</div>
							<div class="text-right">
								<p class="text-lg font-bold {s.uptime_percent >= 99 ? 'text-success' : s.uptime_percent >= 95 ? 'text-warning' : 'text-error'}">
									{s.uptime_percent.toFixed(1)}%
								</p>
								<p class="text-xs text-base-content/50">{s.avg_response_ms}ms avg</p>
							</div>
						</div>
						{#if s.history && s.history.length > 0}
							<div class="flex gap-0.5 mt-3">
								{#each s.history.slice(-30) as r}
									<div class="flex-1 h-6 rounded-sm {r.is_up ? 'bg-success/40' : 'bg-error/40'}" title="{r.response_time_ms}ms"></div>
								{/each}
							</div>
						{/if}
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>
