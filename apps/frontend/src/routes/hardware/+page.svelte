<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchProcesses } from '$lib/api/client';
	import { formatBytes } from '$lib/utils';
	import type { ProcessInfo } from '$lib/types';

	let processes: ProcessInfo[] = $state([]);
	let search = $state('');
	let sortBy = $state<'cpu' | 'mem' | 'pid' | 'name'>('cpu');
	let loading = $state(true);
	let autoRefresh = $state(true);
	let refreshInterval: ReturnType<typeof setInterval>;

	onMount(() => {
		loadProcesses();
		refreshInterval = setInterval(() => {
			if (autoRefresh) loadProcesses();
		}, 5000);
		return () => clearInterval(refreshInterval);
	});

	async function loadProcesses() {
		try {
			processes = await fetchProcesses();
		} catch (e) {
			processes = [];
		}
		loading = false;
	}

	let sorted = $derived.by(() => {
		let list = processes.filter(
			(p) =>
				p.name.toLowerCase().includes(search.toLowerCase()) ||
				p.cmdline.toLowerCase().includes(search.toLowerCase()) ||
				p.username.toLowerCase().includes(search.toLowerCase())
		);
		switch (sortBy) {
			case 'cpu': return list.sort((a, b) => b.cpu_percent - a.cpu_percent);
			case 'mem': return list.sort((a, b) => b.mem_percent - a.mem_percent);
			case 'pid': return list.sort((a, b) => a.pid - b.pid);
			case 'name': return list.sort((a, b) => a.name.localeCompare(b.name));
			default: return list;
		}
	});

	function cpuColor(pct: number) {
		if (pct >= 80) return 'text-error';
		if (pct >= 50) return 'text-warning';
		if (pct >= 20) return 'text-info';
		return 'text-success';
	}
</script>

<svelte:head>
	<title>Process Monitor — Viyoga</title>
</svelte:head>

<div class="p-6 space-y-6">
	<div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
		<div>
			<h1 class="text-2xl font-bold text-base-content">Process Monitor</h1>
			<p class="text-base-content/60 text-sm mt-1">Top {sorted.length} processes by resource usage</p>
		</div>
		<div class="flex items-center gap-3">
			<label class="label cursor-pointer gap-2">
				<span class="text-xs">Auto-refresh</span>
				<input type="checkbox" class="toggle toggle-primary toggle-xs" bind:checked={autoRefresh} />
			</label>
			<button class="btn btn-primary btn-sm" onclick={() => loadProcesses()}>Refresh</button>
		</div>
	</div>

	<!-- Controls -->
	<div class="flex flex-col sm:flex-row gap-3">
		<input type="text" placeholder="Search processes..." class="input input-bordered input-sm flex-1" bind:value={search} />
		<div class="flex gap-1">
			{#each [['cpu', 'CPU'], ['mem', 'Memory'], ['pid', 'PID'], ['name', 'Name']] as [key, label]}
				<button class="btn btn-xs {sortBy === key ? 'btn-primary' : 'btn-ghost'}" onclick={() => sortBy = key as typeof sortBy}>
					{label}
				</button>
			{/each}
		</div>
	</div>

	{#if loading}
		<div class="flex justify-center py-20">
			<span class="loading loading-spinner loading-lg text-primary"></span>
		</div>
	{:else}
		<div class="overflow-x-auto">
			<table class="table table-xs w-full">
				<thead>
					<tr class="text-base-content/70">
						<th>PID</th>
						<th>Name</th>
						<th>User</th>
						<th class="text-right">CPU %</th>
						<th class="text-right">Mem %</th>
						<th class="text-right">RSS</th>
						<th>Threads</th>
						<th>Status</th>
					</tr>
				</thead>
				<tbody>
					{#each sorted as proc (proc.pid)}
						<tr class="hover:bg-base-300/30 transition-colors">
							<td class="font-mono text-xs">{proc.pid}</td>
							<td class="font-mono text-sm font-medium text-primary truncate max-w-40" title={proc.cmdline}>
								{proc.name}
							</td>
							<td class="text-xs text-base-content/60">{proc.username}</td>
							<td class="text-right font-mono text-sm {cpuColor(proc.cpu_percent)}">
								{proc.cpu_percent.toFixed(1)}%
							</td>
							<td class="text-right font-mono text-sm {cpuColor(proc.mem_percent)}">
								{proc.mem_percent.toFixed(1)}%
							</td>
							<td class="text-right text-xs">{formatBytes(proc.mem_rss)}</td>
							<td class="text-center text-xs">{proc.num_threads}</td>
							<td>
								<span class="badge badge-xs {proc.status === 'R' ? 'badge-success' : proc.status === 'S' ? 'badge-info' : 'badge-ghost'}">
									{proc.status}
								</span>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>
