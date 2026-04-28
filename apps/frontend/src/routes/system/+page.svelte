<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchServices, serviceAction } from '$lib/api/client';
	import type { SystemdService } from '$lib/types';

	let services: SystemdService[] = $state([]);
	let filtered: SystemdService[] = $state([]);
	let search = $state('');
	let filter = $state('all'); // all, active, inactive, failed
	let loading = $state(true);
	let actionLoading = $state('');

	onMount(async () => {
		await loadServices();
	});

	async function loadServices() {
		loading = true;
		try {
			services = await fetchServices();
		} catch (e) {
			services = [];
		}
		loading = false;
	}

	$effect(() => {
		filtered = services.filter((s) => {
			const matchesSearch =
				s.name.toLowerCase().includes(search.toLowerCase()) ||
				s.description.toLowerCase().includes(search.toLowerCase());
			const matchesFilter =
				filter === 'all' ||
				(filter === 'active' && s.active_state === 'active') ||
				(filter === 'inactive' && s.active_state === 'inactive') ||
				(filter === 'failed' && s.active_state === 'failed');
			return matchesSearch && matchesFilter;
		});
	});

	async function handleAction(name: string, action: string) {
		actionLoading = `${name}-${action}`;
		try {
			await serviceAction(name, action);
			await loadServices();
		} catch (e) {
			console.error(e);
		}
		actionLoading = '';
	}

	function stateColor(state: string) {
		switch (state) {
			case 'active': return 'badge-success';
			case 'inactive': return 'badge-ghost';
			case 'failed': return 'badge-error';
			default: return 'badge-warning';
		}
	}

	function subStateIcon(sub: string) {
		switch (sub) {
			case 'running': return '▶';
			case 'exited': return '⏹';
			case 'dead': return '💀';
			case 'waiting': return '⏳';
			default: return '○';
		}
	}
</script>

<svelte:head>
	<title>System Services — Viyoga</title>
</svelte:head>

<div class="p-6 space-y-6">
	<!-- Header -->
	<div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
		<div>
			<h1 class="text-2xl font-bold text-base-content">System Services</h1>
			<p class="text-base-content/60 text-sm mt-1">Manage systemd services</p>
		</div>
		<button class="btn btn-primary btn-sm" onclick={() => loadServices()}>
			<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" /></svg>
			Refresh
		</button>
	</div>

	<!-- Filters -->
	<div class="flex flex-col sm:flex-row gap-3">
		<input type="text" placeholder="Search services..." class="input input-bordered input-sm flex-1" bind:value={search} />
		<div class="flex gap-1">
			{#each ['all', 'active', 'inactive', 'failed'] as f}
				<button class="btn btn-xs {filter === f ? 'btn-primary' : 'btn-ghost'}" onclick={() => filter = f}>
					{f.charAt(0).toUpperCase() + f.slice(1)}
				</button>
			{/each}
		</div>
	</div>

	<p class="text-xs text-base-content/50">{filtered.length} of {services.length} services</p>

	<!-- Services Table -->
	{#if loading}
		<div class="flex justify-center py-20">
			<span class="loading loading-spinner loading-lg text-primary"></span>
		</div>
	{:else}
		<div class="overflow-x-auto">
			<table class="table table-sm table-zebra w-full">
				<thead>
					<tr class="text-base-content/70">
						<th>Service</th>
						<th>State</th>
						<th>Sub</th>
						<th>Description</th>
						<th class="text-right">Actions</th>
					</tr>
				</thead>
				<tbody>
					{#each filtered as svc (svc.name)}
						<tr class="hover:bg-base-300/30 transition-colors">
							<td class="font-mono text-sm font-medium text-primary">{svc.name}</td>
							<td>
								<span class="badge badge-sm {stateColor(svc.active_state)}">{svc.active_state}</span>
							</td>
							<td class="text-sm">
								<span class="mr-1">{subStateIcon(svc.sub_state)}</span>{svc.sub_state}
							</td>
							<td class="text-sm text-base-content/70 max-w-xs truncate">{svc.description}</td>
							<td class="text-right">
								<div class="join">
									{#if svc.active_state !== 'active'}
										<button
											class="btn btn-xs btn-success join-item"
											disabled={actionLoading === `${svc.name}-start`}
											onclick={() => handleAction(svc.name, 'start')}
										>
											{actionLoading === `${svc.name}-start` ? '...' : '▶'}
										</button>
									{/if}
									{#if svc.active_state === 'active'}
										<button
											class="btn btn-xs btn-error join-item"
											disabled={actionLoading === `${svc.name}-stop`}
											onclick={() => handleAction(svc.name, 'stop')}
										>
											{actionLoading === `${svc.name}-stop` ? '...' : '⏹'}
										</button>
									{/if}
									<button
										class="btn btn-xs btn-warning join-item"
										disabled={actionLoading === `${svc.name}-restart`}
										onclick={() => handleAction(svc.name, 'restart')}
									>
										{actionLoading === `${svc.name}-restart` ? '...' : '🔄'}
									</button>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>
