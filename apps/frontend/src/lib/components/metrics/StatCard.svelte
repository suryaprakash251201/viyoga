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
</script>

<div class="card bg-base-200 border border-base-300 shadow-lg hover:border-primary/30 transition-all duration-300 hover:shadow-primary/5">
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
				<div class="flex h-10 w-10 items-center justify-center rounded-lg bg-base-300/50 text-xl">
					{icon}
				</div>
			{/if}
		</div>

		{#if percent >= 0}
			<div class="mt-3">
				<div class="flex justify-between text-xs text-base-content/40 mb-1">
					<span>Usage</span>
					<span>{percent.toFixed(1)}%</span>
				</div>
				<progress
					class="progress h-1.5 w-full"
					class:progress-success={percent < 50}
					class:progress-info={percent >= 50 && percent < 70}
					class:progress-warning={percent >= 70 && percent < 90}
					class:progress-error={percent >= 90}
					value={percent}
					max="100"
				></progress>
			</div>
		{/if}

		{#if trend}
			<p class="mt-1 text-xs text-base-content/30">{trend}</p>
		{/if}
	</div>
</div>
