<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchMonitorStatus, addMonitorTarget } from '$lib/api/client';
	import type { MonitorTargetStatus } from '$lib/types';

	let statuses: MonitorTargetStatus[] = $state([]);
	let loading = $state(true);
	let error = $state('');
	let showAddForm = $state(false);
	let newName = $state('');
	let newUrl = $state('');
	let adding = $state(false);

	async function loadData() {
		try {
			statuses = await fetchMonitorStatus();
			error = '';
		} catch (e) {
			error = e instanceof Error ? e.message : 'Monitor unavailable';
			statuses = [];
		}
		loading = false;
	}

	async function handleAddTarget() {
		if (!newName.trim() || !newUrl.trim()) return;
		adding = true;
		try {
			await addMonitorTarget({ name: newName.trim(), url: newUrl.trim(), interval_seconds: 60 });
			newName = '';
			newUrl = '';
			showAddForm = false;
			await loadData();
		} catch (e) {
			error = e instanceof Error ? e.message : 'Failed to add target';
		}
		adding = false;
	}

	onMount(loadData);
</script>

<svelte:head><title>Web Monitor — Viyoga</title></svelte:head>

<div class="p-6 space-y-6 animate-fade-in">
	<div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
		<div>
			<h1 class="text-2xl font-bold text-base-content flex items-center gap-3">
				<span class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-success to-accent text-white shadow-lg shadow-success/20">
					<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
					</svg>
				</span>
				Web Service Monitor
			</h1>
			<p class="text-base-content/60 text-sm mt-1 ml-[52px]">HTTP endpoint health checks & uptime tracking</p>
		</div>

		<button class="btn btn-primary btn-sm gap-2" onclick={() => showAddForm = !showAddForm}>
			<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
				<path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
			</svg>
			Add Target
		</button>
	</div>

	<!-- Add target form -->
	{#if showAddForm}
		<div class="card bg-base-200 border border-primary/20 animate-slide-up">
			<div class="card-body p-5">
				<h3 class="text-sm font-semibold mb-3">Add Monitoring Target</h3>
				<div class="flex flex-col sm:flex-row gap-3">
					<input type="text" bind:value={newName} placeholder="Service name (e.g. API Server)" class="input input-bordered input-sm flex-1" />
					<input type="url" bind:value={newUrl} placeholder="https://example.com/health" class="input input-bordered input-sm flex-[2]" />
					<button class="btn btn-primary btn-sm gap-2" onclick={handleAddTarget} disabled={adding || !newName.trim() || !newUrl.trim()}>
						{#if adding}
							<span class="loading loading-spinner loading-xs"></span>
						{/if}
						Add
					</button>
					<button class="btn btn-ghost btn-sm" onclick={() => showAddForm = false}>Cancel</button>
				</div>
			</div>
		</div>
	{/if}

	{#if error}
		<div class="alert alert-info"><span>{error}</span></div>
	{/if}

	{#if loading}
		<div class="flex justify-center py-20"><span class="loading loading-spinner loading-lg text-primary"></span></div>
	{:else if statuses.length === 0 && !error}
		<!-- Empty state -->
		<div class="card bg-base-200 border border-base-300/50 animate-slide-up">
			<div class="card-body items-center text-center py-16">
				<div class="flex h-16 w-16 items-center justify-center rounded-2xl bg-primary/10 text-primary mb-4">
					<svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
						<path stroke-linecap="round" stroke-linejoin="round" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
					</svg>
				</div>
				<h3 class="text-lg font-semibold">No Monitoring Targets</h3>
				<p class="text-base-content/50 text-sm max-w-md mt-1">
					Add HTTP endpoints to monitor their health, response time, and uptime. Click "Add Target" above to get started.
				</p>
				<button class="btn btn-primary btn-sm mt-4 gap-2" onclick={() => showAddForm = true}>
					<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4" />
					</svg>
					Add Your First Target
				</button>
			</div>
		</div>
	{:else}
		<!-- Stats overview -->
		<div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
			<div class="card bg-base-200 border border-base-300/50 p-4 card-hover">
				<p class="text-xs text-base-content/60 uppercase tracking-wider">Total Targets</p>
				<p class="text-2xl font-bold text-primary mt-1">{statuses.length}</p>
			</div>
			<div class="card bg-base-200 border border-base-300/50 p-4 card-hover">
				<p class="text-xs text-base-content/60 uppercase tracking-wider">Healthy</p>
				<p class="text-2xl font-bold text-success mt-1">{statuses.filter(s => s.last_result?.is_up).length}</p>
			</div>
			<div class="card bg-base-200 border border-base-300/50 p-4 card-hover">
				<p class="text-xs text-base-content/60 uppercase tracking-wider">Down</p>
				<p class="text-2xl font-bold text-error mt-1">{statuses.filter(s => !s.last_result?.is_up).length}</p>
			</div>
			<div class="card bg-base-200 border border-base-300/50 p-4 card-hover">
				<p class="text-xs text-base-content/60 uppercase tracking-wider">Avg Response</p>
				<p class="text-2xl font-bold text-info mt-1">{statuses.length > 0 ? Math.round(statuses.reduce((a, s) => a + (s.avg_response_ms || 0), 0) / statuses.length) : 0}ms</p>
			</div>
		</div>

		<!-- Target list -->
		<div class="grid gap-4">
			{#each statuses as s (s.target.id)}
				<div class="card bg-base-200 border border-base-300/50 hover:border-primary/30 transition-all card-hover">
					<div class="card-body p-4">
						<div class="flex justify-between items-center">
							<div class="flex items-center gap-3">
								<div class="flex h-10 w-10 items-center justify-center rounded-xl {s.last_result?.is_up ? 'bg-success/15 text-success' : 'bg-error/15 text-error'}">
									{#if s.last_result?.is_up}
										<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
											<path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
										</svg>
									{:else}
										<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
											<path stroke-linecap="round" stroke-linejoin="round" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
										</svg>
									{/if}
								</div>
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
