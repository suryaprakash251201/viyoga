<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchContainers, fetchImages, containerAction } from '$lib/api/client';
	import { formatBytes } from '$lib/utils';
	import type { DockerContainer, DockerImage } from '$lib/types';

	let containers: DockerContainer[] = $state([]);
	let images: DockerImage[] = $state([]);
	let tab = $state<'containers' | 'images'>('containers');
	let showAll = $state(true);
	let loading = $state(true);
	let actionLoading = $state('');
	let error = $state('');

	onMount(async () => {
		await Promise.all([loadContainers(), loadImages()]);
	});

	async function loadContainers() {
		loading = true;
		try {
			containers = await fetchContainers(showAll);
			error = '';
		} catch (e) {
			error = e instanceof Error ? e.message : 'Docker unavailable';
			containers = [];
		}
		loading = false;
	}

	async function loadImages() {
		try {
			images = await fetchImages();
		} catch {
			images = [];
		}
	}

	async function handleAction(id: string, action: string) {
		actionLoading = `${id}-${action}`;
		try {
			await containerAction(id, action);
			await loadContainers();
		} catch (e) {
			console.error(e);
		}
		actionLoading = '';
	}

	function stateColor(state: string) {
		switch (state) {
			case 'running': return 'badge-success';
			case 'exited': return 'badge-error';
			case 'paused': return 'badge-warning';
			case 'created': return 'badge-info';
			default: return 'badge-ghost';
		}
	}

	function stateIcon(state: string) {
		switch (state) {
			case 'running': return '🟢';
			case 'exited': return '🔴';
			case 'paused': return '🟡';
			default: return '⚪';
		}
	}
</script>

<svelte:head>
	<title>Docker Manager — Viyoga</title>
</svelte:head>

<div class="p-6 space-y-6">
	<div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
		<div>
			<h1 class="text-2xl font-bold text-base-content">Container Manager</h1>
			<p class="text-base-content/60 text-sm mt-1">
				{containers.length} containers • {images.length} images
			</p>
		</div>
		<button class="btn btn-primary btn-sm" onclick={() => { loadContainers(); loadImages(); }}>Refresh</button>
	</div>

	{#if error}
		<div class="alert alert-warning">
			<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z" /></svg>
			<span>{error}</span>
		</div>
	{/if}

	<!-- Tabs -->
	<div role="tablist" class="tabs tabs-bordered">
		<button role="tab" class="tab {tab === 'containers' ? 'tab-active' : ''}" onclick={() => tab = 'containers'}>
			🐳 Containers ({containers.length})
		</button>
		<button role="tab" class="tab {tab === 'images' ? 'tab-active' : ''}" onclick={() => tab = 'images'}>
			📦 Images ({images.length})
		</button>
	</div>

	{#if tab === 'containers'}
		<div class="flex items-center gap-2 mb-2">
			<label class="label cursor-pointer gap-2">
				<span class="text-xs">Show stopped</span>
				<input type="checkbox" class="toggle toggle-xs toggle-primary" bind:checked={showAll} onchange={() => loadContainers()} />
			</label>
		</div>

		{#if loading}
			<div class="flex justify-center py-20">
				<span class="loading loading-spinner loading-lg text-primary"></span>
			</div>
		{:else}
			<div class="grid gap-4">
				{#each containers as c (c.id)}
					<div class="card bg-base-200 shadow-sm border border-base-300/50 hover:border-primary/30 transition-colors">
						<div class="card-body p-4">
							<div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-3">
								<div class="flex items-center gap-3">
									<span class="text-lg">{stateIcon(c.state)}</span>
									<div>
										<h3 class="font-bold text-base-content">{c.name}</h3>
										<p class="text-xs text-base-content/50 font-mono">{c.id} • {c.image}</p>
									</div>
								</div>
								<div class="flex items-center gap-2">
									<span class="badge {stateColor(c.state)} badge-sm">{c.state}</span>
									<div class="join">
										{#if c.state !== 'running'}
											<button class="btn btn-xs btn-success join-item" disabled={!!actionLoading} onclick={() => handleAction(c.id, 'start')}>▶</button>
										{/if}
										{#if c.state === 'running'}
											<button class="btn btn-xs btn-error join-item" disabled={!!actionLoading} onclick={() => handleAction(c.id, 'stop')}>⏹</button>
										{/if}
										<button class="btn btn-xs btn-warning join-item" disabled={!!actionLoading} onclick={() => handleAction(c.id, 'restart')}>🔄</button>
										<button class="btn btn-xs btn-ghost join-item" disabled={!!actionLoading} onclick={() => handleAction(c.id, 'remove')}>🗑</button>
									</div>
								</div>
							</div>
							{#if c.ports && c.ports.length > 0}
								<div class="mt-2 flex flex-wrap gap-1">
									{#each c.ports as p}
										{#if p.public}
											<span class="badge badge-outline badge-xs font-mono">
												{p.ip || '0.0.0.0'}:{p.public} → {p.private}/{p.type}
											</span>
										{/if}
									{/each}
								</div>
							{/if}
							<p class="text-xs text-base-content/40 mt-1">{c.status}</p>
						</div>
					</div>
				{/each}
			</div>
		{/if}

	{:else if tab === 'images'}
		<div class="overflow-x-auto">
			<table class="table table-sm w-full">
				<thead>
					<tr class="text-base-content/70">
						<th>Image ID</th>
						<th>Tags</th>
						<th class="text-right">Size</th>
					</tr>
				</thead>
				<tbody>
					{#each images as img (img.id)}
						<tr class="hover:bg-base-300/30 transition-colors">
							<td class="font-mono text-xs text-primary">{img.id}</td>
							<td>
								<div class="flex flex-wrap gap-1">
									{#each img.tags as tag}
										<span class="badge badge-ghost badge-xs font-mono">{tag}</span>
									{/each}
									{#if !img.tags || img.tags.length === 0}
										<span class="text-base-content/40 text-xs">&lt;none&gt;</span>
									{/if}
								</div>
							</td>
							<td class="text-right text-sm">{formatBytes(img.size)}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>
