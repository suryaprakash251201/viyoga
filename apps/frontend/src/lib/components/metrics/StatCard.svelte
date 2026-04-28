<script lang="ts">
	import { getUsageColor } from '$lib/utils';

	let { title, value, subtitle = '', icon = '', percent = -1, trend = '' } = $props<{
		title: string;
		value: string;
		subtitle?: string;
		icon?: string;
		percent?: number;
		trend?: string;
	}>();

	const iconColors: Record<string, string> = {
		'cpu': 'from-cyan-500 to-blue-600',
		'memory': 'from-violet-500 to-purple-600',
		'network-up': 'from-emerald-500 to-green-600',
		'network-down': 'from-blue-500 to-indigo-600',
	};

	function getIconGradient(ic: string): string {
		if (ic.includes('⚡') || ic === 'cpu') return iconColors['cpu'];
		if (ic.includes('🧠') || ic === 'memory') return iconColors['memory'];
		if (ic.includes('📤') || ic.includes('↑')) return iconColors['network-up'];
		if (ic.includes('📥') || ic.includes('↓')) return iconColors['network-down'];
		return 'from-primary to-accent';
	}
</script>

<div class="card bg-base-200 border border-base-300/60 shadow-lg hover:border-primary/30 transition-all duration-300 hover:shadow-primary/5 card-hover group">
	<div class="card-body p-4">
		<div class="flex items-start justify-between">
			<div class="flex-1">
				<p class="text-xs font-medium uppercase tracking-wider text-base-content/40">{title}</p>
				<p class="mt-1 text-2xl font-bold tabular-nums {percent >= 0 ? getUsageColor(percent) : 'text-base-content'}">
					{value}
				</p>
				{#if subtitle}
					<p class="mt-0.5 text-xs text-base-content/50">{subtitle}</p>
				{/if}
			</div>
			{#if icon}
				<div class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br {getIconGradient(icon)} text-white shadow-lg opacity-80 group-hover:opacity-100 group-hover:scale-110 transition-all duration-300">
					{#if icon.includes('⚡') || icon === 'cpu'}
						<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
							<path stroke-linecap="round" stroke-linejoin="round" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
						</svg>
					{:else if icon.includes('🧠') || icon === 'memory'}
						<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
							<path stroke-linecap="round" stroke-linejoin="round" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4" />
						</svg>
					{:else if icon.includes('↑') || icon.includes('📤')}
						<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
							<path stroke-linecap="round" stroke-linejoin="round" d="M7 11l5-5m0 0l5 5m-5-5v12" />
						</svg>
					{:else if icon.includes('↓') || icon.includes('📥')}
						<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
							<path stroke-linecap="round" stroke-linejoin="round" d="M17 13l-5 5m0 0l-5-5m5 5V6" />
						</svg>
					{:else}
						<span class="text-lg">{icon}</span>
					{/if}
				</div>
			{/if}
		</div>

		{#if percent >= 0}
			<div class="mt-3">
				<div class="flex justify-between text-xs text-base-content/40 mb-1">
					<span>Usage</span>
					<span class="tabular-nums">{percent.toFixed(1)}%</span>
				</div>
				<div class="h-1.5 w-full rounded-full bg-base-300/60 overflow-hidden">
					<div
						class="h-full rounded-full transition-all duration-500 ease-out
							{percent < 50 ? 'bg-gradient-to-r from-success to-success/80' :
							 percent < 70 ? 'bg-gradient-to-r from-info to-info/80' :
							 percent < 90 ? 'bg-gradient-to-r from-warning to-warning/80' :
							 'bg-gradient-to-r from-error to-error/80'}"
						style="width: {Math.min(percent, 100)}%"
					></div>
				</div>
			</div>
		{/if}

		{#if trend}
			<p class="mt-1 text-xs text-base-content/30">{trend}</p>
		{/if}
	</div>
</div>
