<script lang="ts">
	import { getUsageHex } from '$lib/utils';
	import { onMount } from 'svelte';

	let { percent = 0, size = 120, strokeWidth = 8, label = '', sublabel = '' } = $props<{
		percent: number;
		size?: number;
		strokeWidth?: number;
		label?: string;
		sublabel?: string;
	}>();

	const radius = $derived((size - strokeWidth) / 2);
	const circumference = $derived(2 * Math.PI * radius);

	let animatedPercent = $state(0);
	let mounted = $state(false);

	onMount(() => {
		mounted = true;
	});

	$effect(() => {
		if (mounted) {
			// Smooth animation toward target
			const target = Math.min(Math.max(percent, 0), 100);
			const step = () => {
				const diff = target - animatedPercent;
				if (Math.abs(diff) < 0.5) {
					animatedPercent = target;
				} else {
					animatedPercent += diff * 0.15;
					requestAnimationFrame(step);
				}
			};
			requestAnimationFrame(step);
		}
	});

	const offset = $derived(circumference - (animatedPercent / 100) * circumference);
	const color = $derived(getUsageHex(animatedPercent));
</script>

<div class="flex flex-col items-center gap-2">
	<div class="relative" style="width: {size}px; height: {size}px;">
		<svg width={size} height={size} class="-rotate-90">
			<!-- Background ring -->
			<circle
				cx={size / 2}
				cy={size / 2}
				r={radius}
				fill="none"
				stroke="oklch(0.22 0.02 270)"
				stroke-width={strokeWidth}
			/>
			<!-- Progress ring -->
			<circle
				cx={size / 2}
				cy={size / 2}
				r={radius}
				fill="none"
				stroke={color}
				stroke-width={strokeWidth}
				stroke-dasharray={circumference}
				stroke-dashoffset={offset}
				stroke-linecap="round"
				class="transition-all duration-500 ease-out"
				style="filter: drop-shadow(0 0 6px {color}40);"
			/>
		</svg>
		<!-- Center text -->
		<div class="absolute inset-0 flex flex-col items-center justify-center">
			<span class="text-xl font-bold tabular-nums" style="color: {color}">
				{animatedPercent.toFixed(1)}%
			</span>
		</div>
	</div>
	{#if label}
		<span class="text-xs font-medium text-base-content/60">{label}</span>
	{/if}
	{#if sublabel}
		<span class="text-xs text-base-content/30">{sublabel}</span>
	{/if}
</div>
