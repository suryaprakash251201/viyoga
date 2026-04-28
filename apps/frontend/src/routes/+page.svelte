<script lang="ts">
	import StatCard from '$lib/components/metrics/StatCard.svelte';
	import GaugeChart from '$lib/components/metrics/GaugeChart.svelte';
	import AreaChart from '$lib/components/charts/AreaChart.svelte';
	import {
		cpuMetrics,
		memoryMetrics,
		diskMetrics,
		networkMetrics,
		systemInfo,
		wsConnected,
		cpuHistory,
		memoryHistory,
		networkSentHistory,
		networkRecvHistory,
		networkSendRate,
		networkRecvRate
	} from '$lib/stores/metrics';
	import { formatBytes, formatPercent, formatUptime, formatNumber } from '$lib/utils';

	// Derived chart data
	const cpuChartSeries = $derived([
		{ name: 'CPU %', data: $cpuHistory.map((h) => parseFloat(h.value.toFixed(1))) }
	]);
	const cpuChartCategories = $derived($cpuHistory.map((h) => h.time));

	const memChartSeries = $derived([
		{ name: 'RAM %', data: $memoryHistory.map((h) => parseFloat(h.value.toFixed(1))) }
	]);
	const memChartCategories = $derived($memoryHistory.map((h) => h.time));

	const netChartSeries = $derived([
		{ name: '↑ Send', data: $networkSentHistory.map((h) => parseFloat(h.value.toFixed(0))) },
		{ name: '↓ Recv', data: $networkRecvHistory.map((h) => parseFloat(h.value.toFixed(0))) }
	]);
	const netChartCategories = $derived($networkSentHistory.map((h) => h.time));

	// Primary disk (first partition or the largest)
	const primaryDisk = $derived(
		$diskMetrics?.partitions?.length
			? $diskMetrics.partitions.reduce((a, b) => (b.total_bytes > a.total_bytes ? b : a))
			: null
	);
</script>

<div class="space-y-6">
	<!-- Section: Overview Gauges -->
	<section>
		<h2 class="mb-4 text-sm font-semibold uppercase tracking-wider text-base-content/40">System Overview</h2>
		<div class="grid grid-cols-2 gap-4 sm:gap-6 lg:grid-cols-4">
			<div class="card bg-base-200 border border-base-300 p-6 flex items-center justify-center">
				<GaugeChart
					percent={$cpuMetrics?.usage_percent ?? 0}
					label="CPU"
					sublabel={$cpuMetrics?.model_name?.split(' ').slice(0, 3).join(' ') ?? ''}
				/>
			</div>
			<div class="card bg-base-200 border border-base-300 p-6 flex items-center justify-center">
				<GaugeChart
					percent={$memoryMetrics?.usage_percent ?? 0}
					label="Memory"
					sublabel={$memoryMetrics ? `${formatBytes($memoryMetrics.used_bytes)} / ${formatBytes($memoryMetrics.total_bytes)}` : ''}
				/>
			</div>
			<div class="card bg-base-200 border border-base-300 p-6 flex items-center justify-center">
				<GaugeChart
					percent={primaryDisk?.used_percent ?? 0}
					label="Disk"
					sublabel={primaryDisk ? `${formatBytes(primaryDisk.used_bytes)} / ${formatBytes(primaryDisk.total_bytes)}` : ''}
				/>
			</div>
			<div class="card bg-base-200 border border-base-300 p-6 flex items-center justify-center">
				<GaugeChart
					percent={$memoryMetrics?.swap_usage_percent ?? 0}
					label="Swap"
					sublabel={$memoryMetrics ? `${formatBytes($memoryMetrics.swap_used_bytes)} / ${formatBytes($memoryMetrics.swap_total_bytes)}` : ''}
				/>
			</div>
		</div>
	</section>

	<!-- Section: Quick Stats -->
	<section>
		<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
			<StatCard
				title="CPU Usage"
				value={$cpuMetrics ? formatPercent($cpuMetrics.usage_percent) : '--'}
				subtitle="{$cpuMetrics?.logical_count ?? '--'} cores • {$cpuMetrics?.frequency_mhz?.toFixed(0) ?? '--'} MHz"
				icon="⚡"
				percent={$cpuMetrics?.usage_percent ?? -1}
			/>
			<StatCard
				title="Memory"
				value={$memoryMetrics ? formatBytes($memoryMetrics.used_bytes) : '--'}
				subtitle="of {$memoryMetrics ? formatBytes($memoryMetrics.total_bytes) : '--'} total"
				icon="🧠"
				percent={$memoryMetrics?.usage_percent ?? -1}
			/>
			<StatCard
				title="Network ↑"
				value={formatBytes($networkSendRate) + '/s'}
				subtitle="Total: {$networkMetrics ? formatBytes($networkMetrics.total_bytes_sent) : '--'}"
				icon="📤"
			/>
			<StatCard
				title="Network ↓"
				value={formatBytes($networkRecvRate) + '/s'}
				subtitle="Total: {$networkMetrics ? formatBytes($networkMetrics.total_bytes_recv) : '--'}"
				icon="📥"
			/>
		</div>
	</section>

	<!-- Section: Real-time Charts -->
	<section>
		<h2 class="mb-4 text-sm font-semibold uppercase tracking-wider text-base-content/40">Real-time Metrics</h2>
		<div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
			<!-- CPU Chart -->
			<div class="card bg-base-200 border border-base-300 p-4">
				<AreaChart
					title="CPU Usage"
					series={cpuChartSeries}
					categories={cpuChartCategories}
					colors={['#00d4ff']}
					yAxisMax={100}
					yAxisFormatter={(v) => `${v.toFixed(0)}%`}
					height={220}
				/>
			</div>

			<!-- Memory Chart -->
			<div class="card bg-base-200 border border-base-300 p-4">
				<AreaChart
					title="Memory Usage"
					series={memChartSeries}
					categories={memChartCategories}
					colors={['#7c3aed']}
					yAxisMax={100}
					yAxisFormatter={(v) => `${v.toFixed(0)}%`}
					height={220}
				/>
			</div>

			<!-- Network Chart -->
			<div class="card bg-base-200 border border-base-300 p-4 lg:col-span-2">
				<AreaChart
					title="Network Throughput"
					series={netChartSeries}
					categories={netChartCategories}
					colors={['#22c55e', '#3b82f6']}
					yAxisFormatter={(v) => formatBytes(v) + '/s'}
					height={220}
				/>
			</div>
		</div>
	</section>

	<!-- Section: Disk Partitions -->
	{#if $diskMetrics?.partitions?.length}
		<section>
			<h2 class="mb-4 text-sm font-semibold uppercase tracking-wider text-base-content/40">Disk Partitions</h2>
			<div class="overflow-x-auto">
				<table class="table table-sm">
					<thead>
						<tr class="text-base-content/40">
							<th>Device</th>
							<th>Mount</th>
							<th>Type</th>
							<th>Total</th>
							<th>Used</th>
							<th>Free</th>
							<th>Usage</th>
						</tr>
					</thead>
					<tbody>
						{#each $diskMetrics.partitions as part}
							<tr class="hover:bg-base-300/30">
								<td class="font-mono text-xs">{part.device}</td>
								<td class="font-mono text-xs">{part.mount_point}</td>
								<td class="text-xs text-base-content/50">{part.fs_type}</td>
								<td class="text-xs tabular-nums">{formatBytes(part.total_bytes)}</td>
								<td class="text-xs tabular-nums">{formatBytes(part.used_bytes)}</td>
								<td class="text-xs tabular-nums">{formatBytes(part.free_bytes)}</td>
								<td>
									<div class="flex items-center gap-2">
										<progress
											class="progress h-1.5 w-16"
											class:progress-success={part.used_percent < 50}
											class:progress-info={part.used_percent >= 50 && part.used_percent < 70}
											class:progress-warning={part.used_percent >= 70 && part.used_percent < 90}
											class:progress-error={part.used_percent >= 90}
											value={part.used_percent}
											max="100"
										></progress>
										<span class="text-xs tabular-nums">{part.used_percent.toFixed(1)}%</span>
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</section>
	{/if}

	<!-- Section: Network Interfaces -->
	{#if $networkMetrics?.interfaces?.length}
		<section>
			<h2 class="mb-4 text-sm font-semibold uppercase tracking-wider text-base-content/40">Network Interfaces</h2>
			<div class="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
				{#each $networkMetrics.interfaces.filter((i) => i.bytes_sent > 0 || i.bytes_recv > 0) as iface}
					<div class="card bg-base-200 border border-base-300 p-4">
						<h3 class="font-mono text-sm font-semibold text-primary">{iface.name}</h3>
						<div class="mt-2 grid grid-cols-2 gap-2 text-xs">
							<div>
								<span class="text-base-content/40">Sent</span>
								<p class="font-semibold tabular-nums text-success">{formatBytes(iface.bytes_sent)}</p>
							</div>
							<div>
								<span class="text-base-content/40">Received</span>
								<p class="font-semibold tabular-nums text-info">{formatBytes(iface.bytes_recv)}</p>
							</div>
							<div>
								<span class="text-base-content/40">Packets ↑</span>
								<p class="tabular-nums">{formatNumber(iface.packets_sent)}</p>
							</div>
							<div>
								<span class="text-base-content/40">Packets ↓</span>
								<p class="tabular-nums">{formatNumber(iface.packets_recv)}</p>
							</div>
						</div>
					</div>
				{/each}
			</div>
		</section>
	{/if}

	<!-- Section: System Info -->
	{#if $systemInfo}
		<section>
			<h2 class="mb-4 text-sm font-semibold uppercase tracking-wider text-base-content/40">System Information</h2>
			<div class="card bg-base-200 border border-base-300">
				<div class="card-body p-4">
					<div class="grid grid-cols-2 gap-x-8 gap-y-3 text-sm sm:grid-cols-3 lg:grid-cols-4">
						<div>
							<span class="text-xs text-base-content/40">Hostname</span>
							<p class="font-mono font-semibold">{$systemInfo.hostname}</p>
						</div>
						<div>
							<span class="text-xs text-base-content/40">OS</span>
							<p class="font-semibold">{$systemInfo.platform} {$systemInfo.platform_version}</p>
						</div>
						<div>
							<span class="text-xs text-base-content/40">Kernel</span>
							<p class="font-mono text-xs">{$systemInfo.kernel_version}</p>
						</div>
						<div>
							<span class="text-xs text-base-content/40">Architecture</span>
							<p class="font-mono">{$systemInfo.kernel_arch}</p>
						</div>
						<div>
							<span class="text-xs text-base-content/40">Uptime</span>
							<p class="font-semibold">{formatUptime($systemInfo.uptime_seconds)}</p>
						</div>
						<div>
							<span class="text-xs text-base-content/40">Processes</span>
							<p class="font-semibold tabular-nums">{$systemInfo.procs}</p>
						</div>
						<div>
							<span class="text-xs text-base-content/40">Go Runtime</span>
							<p class="font-mono text-xs">{$systemInfo.go_version}</p>
						</div>
						<div>
							<span class="text-xs text-base-content/40">Viyoga</span>
							<p class="font-mono text-xs text-primary">v{$systemInfo.viyoga_version}</p>
						</div>
					</div>
				</div>
			</div>
		</section>
	{/if}

	<!-- Connection status banner -->
	{#if !$wsConnected}
		<div class="fixed bottom-4 right-4 z-50">
			<div class="alert alert-warning shadow-lg max-w-sm">
				<span class="text-sm">⚠️ Disconnected from server. Reconnecting...</span>
			</div>
		</div>
	{/if}
</div>
