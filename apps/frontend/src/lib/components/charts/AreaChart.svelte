<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { browser } from '$app/environment';

	let {
		series = [],
		categories = [],
		title = '',
		height = 200,
		colors = ['#00d4ff'],
		yAxisMax = undefined,
		yAxisFormatter = (val: number) => val.toFixed(1)
	} = $props<{
		series: { name: string; data: number[] }[];
		categories: string[];
		title?: string;
		height?: number;
		colors?: string[];
		yAxisMax?: number | undefined;
		yAxisFormatter?: (val: number) => string;
	}>();

	let chartEl: HTMLDivElement | undefined = $state();
	let chart: any = null;

	onMount(async () => {
		if (!browser || !chartEl) return;

		const ApexCharts = (await import('apexcharts')).default;

		const options = {
			chart: {
				type: 'area' as const,
				height,
				sparkline: { enabled: false },
				toolbar: { show: false },
				zoom: { enabled: false },
				background: 'transparent',
				fontFamily: 'Inter, system-ui, sans-serif',
				animations: {
					enabled: true,
					easing: 'smooth',
					dynamicAnimation: { speed: 500 }
				}
			},
			series,
			xaxis: {
				categories,
				labels: {
					show: true,
					style: { colors: 'oklch(0.5 0.01 270)', fontSize: '10px' },
					rotate: 0,
					hideOverlappingLabels: true
				},
				axisBorder: { show: false },
				axisTicks: { show: false }
			},
			yaxis: {
				max: yAxisMax,
				labels: {
					show: true,
					style: { colors: 'oklch(0.5 0.01 270)', fontSize: '10px' },
					formatter: yAxisFormatter
				}
			},
			grid: {
				show: true,
				borderColor: 'oklch(0.22 0.02 270)',
				strokeDashArray: 3,
				padding: { left: 8, right: 8 }
			},
			stroke: {
				curve: 'smooth' as const,
				width: 2
			},
			colors,
			fill: {
				type: 'gradient',
				gradient: {
					shadeIntensity: 1,
					opacityFrom: 0.4,
					opacityTo: 0.05,
					stops: [0, 100]
				}
			},
			dataLabels: { enabled: false },
			tooltip: {
				theme: 'dark',
				x: { show: true },
				y: { formatter: yAxisFormatter }
			},
			legend: {
				show: series.length > 1,
				labels: { colors: 'oklch(0.7 0.01 270)' },
				position: 'top' as const
			},
			title: title
				? {
						text: title,
						style: { color: 'oklch(0.7 0.01 270)', fontSize: '12px', fontWeight: 500 }
					}
				: undefined
		};

		chart = new ApexCharts(chartEl, options);
		chart.render();
	});

	$effect(() => {
		if (chart && series.length > 0 && categories.length > 0) {
			chart.updateOptions(
				{
					series,
					xaxis: { categories }
				},
				false,
				true
			);
		}
	});

	onDestroy(() => {
		if (chart) {
			chart.destroy();
			chart = null;
		}
	});
</script>

<div bind:this={chartEl} class="w-full"></div>
