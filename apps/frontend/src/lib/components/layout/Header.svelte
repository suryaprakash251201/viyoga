<script lang="ts">
	import { wsConnected } from '$lib/stores/metrics';
	import { systemInfo } from '$lib/stores/metrics';
	import { theme, isDark } from '$lib/stores/theme';
	import { formatUptime } from '$lib/utils';

	let now = $state(new Date());
	$effect(() => {
		const interval = setInterval(() => (now = new Date()), 1000);
		return () => clearInterval(interval);
	});
</script>

<header class="flex h-16 items-center justify-between border-b border-base-300/40 bg-base-200/90 px-6 backdrop-blur-md">
	<!-- Left: Server info -->
	<div class="flex items-center gap-4">
		<div>
			<h1 class="text-sm font-semibold text-base-content flex items-center gap-2">
				{$systemInfo?.hostname ?? 'Connecting...'}
				{#if $wsConnected}
					<span class="relative flex h-2 w-2">
						<span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-success opacity-75"></span>
						<span class="relative inline-flex h-2 w-2 rounded-full bg-success"></span>
					</span>
				{/if}
			</h1>
			<p class="text-xs text-base-content/50">
				{#if $systemInfo}
					{$systemInfo.platform} {$systemInfo.platform_version} • Kernel {$systemInfo.kernel_version}
				{:else}
					Waiting for connection...
				{/if}
			</p>
		</div>
	</div>

	<!-- Right: Status indicators -->
	<div class="flex items-center gap-2.5">
		<!-- Uptime -->
		{#if $systemInfo}
			<div class="hidden items-center gap-1.5 rounded-lg bg-base-300/40 px-3 py-1.5 text-xs text-base-content/60 sm:flex">
				<svg xmlns="http://www.w3.org/2000/svg" class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
					<path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
				</svg>
				<span>Uptime: {formatUptime($systemInfo.uptime_seconds)}</span>
			</div>
		{/if}

		<!-- Connection status -->
		<div class="flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-xs font-medium {$wsConnected ? 'bg-success/10 text-success' : 'bg-error/10 text-error'}">
			{#if $wsConnected}
				<span class="relative flex h-1.5 w-1.5">
					<span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-success opacity-75"></span>
					<span class="relative inline-flex h-1.5 w-1.5 rounded-full bg-success"></span>
				</span>
				Live
			{:else}
				<span class="h-1.5 w-1.5 rounded-full bg-error"></span>
				Offline
			{/if}
		</div>

		<!-- Theme toggle -->
		<button
			onclick={() => theme.toggle()}
			class="btn btn-ghost btn-sm btn-circle relative overflow-hidden hover:bg-base-300/50"
			title={$isDark ? 'Switch to Light mode' : 'Switch to Dark mode'}
		>
			<div class="relative w-5 h-5">
				{#if $isDark}
					<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-warning animate-scale-in" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
					</svg>
				{:else}
					<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-primary animate-scale-in" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
					</svg>
				{/if}
			</div>
		</button>

		<!-- Time -->
		<div class="hidden rounded-lg bg-base-300/40 px-3 py-1.5 text-xs tabular-nums text-base-content/50 lg:block font-mono">
			{now.toLocaleTimeString('en-US', { hour12: false })}
		</div>
	</div>
</header>
