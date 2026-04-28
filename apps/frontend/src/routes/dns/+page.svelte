<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchDNSStats } from '$lib/api/client';
	import type { DNSStats } from '$lib/types';
	import { formatNumber } from '$lib/utils';

	let stats: DNSStats | null = $state(null);
	let loading = $state(true);
	let error = $state('');
	let errorType = $state<'not_configured' | 'connection_error'>('not_configured');
	let showSetup = $state(false);

	onMount(async () => {
		try {
			stats = await fetchDNSStats();
		} catch (e) {
			const msg = e instanceof Error ? e.message : 'DNS unavailable';
			error = msg;
			// If backend says "not available" → not configured; otherwise connection/API error
			errorType = msg.includes('not available') || msg.includes('unavailable')
				? 'not_configured'
				: 'connection_error';
		}
		loading = false;
	});
</script>

<svelte:head><title>DNS Gateway — Viyoga</title></svelte:head>

<div class="p-6 space-y-6 animate-fade-in">
	<div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
		<div>
			<h1 class="text-2xl font-bold text-base-content flex items-center gap-3">
				<span class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-info to-primary text-white shadow-lg shadow-info/20">
					<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
					</svg>
				</span>
				DNS Gateway
			</h1>
			<p class="text-base-content/60 text-sm mt-1 ml-[52px]">Network-level ad/tracker blocking via Technitium DNS</p>
		</div>
	</div>

	{#if error}
		<!-- Setup Guide -->
		<div class="card bg-base-200 border border-base-300/50 overflow-hidden animate-slide-up">
			<div class="card-body p-0">
				<!-- Gradient banner -->
				<div class="bg-gradient-to-r {errorType === 'connection_error' ? 'from-error/20 via-warning/10' : 'from-info/20 via-primary/10'} to-transparent px-6 py-5 border-b border-base-300/30">
					<div class="flex items-center gap-3">
						<div class="flex h-12 w-12 items-center justify-center rounded-2xl {errorType === 'connection_error' ? 'bg-error/20 text-error' : 'bg-info/20 text-info'}">
							<svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
								{#if errorType === 'connection_error'}
									<path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
								{:else}
									<path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
								{/if}
							</svg>
						</div>
						<div>
							<h2 class="text-lg font-semibold">{errorType === 'connection_error' ? 'DNS Connection Failed' : 'DNS Gateway Not Configured'}</h2>
							<p class="text-sm text-base-content/60">
								{#if errorType === 'connection_error'}
									Technitium DNS is configured but not reachable: <span class="text-error font-mono text-xs">{error}</span>
								{:else}
									Set up Technitium DNS Server to enable network-level ad blocking
								{/if}
							</p>
						</div>
					</div>
				</div>

				<div class="p-6 space-y-5">
					{#if errorType === 'connection_error'}
						<!-- Troubleshooting for connection error -->
						<div class="space-y-4">
							<h3 class="text-sm font-semibold uppercase tracking-wider text-base-content/50">Troubleshooting</h3>

							{#each [
								{ title: 'Check Technitium is running', cmd: 'sudo systemctl status technitium-dns-server', desc: 'Verify the DNS server service is active' },
								{ title: 'Verify the API URL', cmd: null, desc: 'Try opening http://localhost:5380 in your browser on the server' },
								{ title: 'Check your API token', cmd: null, desc: 'Go to Technitium Admin → Settings → API Token and regenerate if needed' },
								{ title: 'Restart Viyoga after config changes', cmd: 'sudo systemctl restart viyoga', desc: 'Apply changes to viyoga.yaml' }
							] as item, i}
								<div class="flex gap-4 p-4 rounded-xl bg-base-300/30 hover:bg-base-300/50 transition-colors">
									<div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-warning/15 text-warning text-sm font-bold">
										{i + 1}
									</div>
									<div class="flex-1 min-w-0">
										<p class="text-sm font-medium">{item.title}</p>
										<p class="text-xs text-base-content/50 mt-0.5">{item.desc}</p>
										{#if item.cmd}
											<code class="block mt-2 text-xs font-mono bg-base-100 rounded-lg px-3 py-2 text-primary overflow-x-auto">{item.cmd}</code>
										{/if}
									</div>
								</div>
							{/each}
						</div>
					{:else}
						<!-- Setup steps for not configured -->
						<div class="space-y-4">
							<h3 class="text-sm font-semibold uppercase tracking-wider text-base-content/50">Setup Steps</h3>

							{#each [
								{ step: 1, title: 'Install Technitium DNS Server', cmd: 'curl -sSL https://download.technitium.com/dns/install.sh | sudo bash', desc: 'Installs Technitium DNS on your server' },
								{ step: 2, title: 'Get API Token', cmd: null, desc: 'Open http://your-server:5380 → Settings → API Token → Copy' },
								{ step: 3, title: 'Configure Viyoga', cmd: 'sudo nano /etc/viyoga/viyoga.yaml', desc: 'Add the dns section with your API URL and token' },
								{ step: 4, title: 'Restart Viyoga', cmd: 'sudo systemctl restart viyoga', desc: 'Apply the configuration changes' }
							] as item}
								<div class="flex gap-4 p-4 rounded-xl bg-base-300/30 hover:bg-base-300/50 transition-colors">
									<div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-primary/15 text-primary text-sm font-bold">
										{item.step}
									</div>
									<div class="flex-1 min-w-0">
										<p class="text-sm font-medium">{item.title}</p>
										<p class="text-xs text-base-content/50 mt-0.5">{item.desc}</p>
										{#if item.cmd}
											<code class="block mt-2 text-xs font-mono bg-base-100 rounded-lg px-3 py-2 text-primary overflow-x-auto">{item.cmd}</code>
										{/if}
									</div>
								</div>
							{/each}
						</div>
					{/if}

					<!-- Config example -->
					<button class="btn btn-sm btn-ghost gap-2 text-base-content/60" onclick={() => showSetup = !showSetup}>
						<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 transition-transform {showSetup ? 'rotate-90' : ''}" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
							<path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
						</svg>
						View example config
					</button>

					{#if showSetup}
						<div class="bg-base-100 rounded-xl p-4 font-mono text-xs leading-relaxed animate-slide-up">
							<pre class="text-base-content/80"><span class="text-primary">dns:</span>
  <span class="text-info">engine:</span> <span class="text-warning">"technitium"</span>
  <span class="text-info">api_url:</span> <span class="text-warning">"http://localhost:5380"</span>
  <span class="text-info">api_token:</span> <span class="text-warning">"your-api-token-here"</span></pre>
						</div>
					{/if}

					<!-- Action -->
					<div class="flex items-center gap-3 pt-2">
						<a href="https://technitium.com/dns/" target="_blank" rel="noopener noreferrer" class="btn btn-primary btn-sm gap-2">
							<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
								<path stroke-linecap="round" stroke-linejoin="round" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
							</svg>
							Technitium Docs
						</a>
						<button class="btn btn-ghost btn-sm" onclick={() => window.location.reload()}>
							Retry Connection
						</button>
					</div>
				</div>
			</div>
		</div>

	{:else if loading}
		<div class="flex justify-center py-20"><span class="loading loading-spinner loading-lg text-primary"></span></div>

	{:else if stats}
		<div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
			<div class="card bg-base-200 border border-base-300/50 p-5 card-hover">
				<p class="text-xs text-base-content/60 uppercase tracking-wider">Total Queries</p>
				<p class="text-3xl font-bold text-primary mt-2">{formatNumber(stats.total_queries)}</p>
			</div>
			<div class="card bg-base-200 border border-base-300/50 p-5 card-hover">
				<p class="text-xs text-base-content/60 uppercase tracking-wider">Blocked</p>
				<p class="text-3xl font-bold text-error mt-2">{formatNumber(stats.total_blocked_queries)}</p>
			</div>
			<div class="card bg-base-200 border border-base-300/50 p-5 card-hover">
				<p class="text-xs text-base-content/60 uppercase tracking-wider">Block Rate</p>
				<p class="text-3xl font-bold text-warning mt-2">{stats.blocked_percent.toFixed(1)}%</p>
			</div>
			<div class="card bg-base-200 border border-base-300/50 p-5 card-hover">
				<p class="text-xs text-base-content/60 uppercase tracking-wider">Clients</p>
				<p class="text-3xl font-bold text-accent mt-2">{stats.total_clients}</p>
			</div>
		</div>
	{/if}
</div>
