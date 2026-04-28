<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchUsers } from '$lib/api/client';
	import type { SystemUser } from '$lib/types';

	let users: SystemUser[] = $state([]);
	let loading = $state(true);
	let activeTab = $state<'overview' | 'users' | 'audit'>('overview');

	// Security metrics (mocked for demo, can be wired to real backend)
	let securityMetrics = $state({
		firewallActive: true,
		sshPort: 22,
		failedLogins: 3,
		activeSessions: 1,
		lastUpdate: new Date().toISOString(),
		sslExpiry: '2027-01-15',
		authMethod: 'Password',
		sudoUsers: 2,
	});

	onMount(async () => {
		try {
			users = await fetchUsers(false);
		} catch {
			users = [];
		}
		loading = false;
	});
</script>

<svelte:head>
	<title>Security — Viyoga</title>
</svelte:head>

<div class="p-6 space-y-6 animate-fade-in">
	<div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
		<div>
			<h1 class="text-2xl font-bold text-base-content">Security Center</h1>
			<p class="text-base-content/60 text-sm mt-1">Authentication, audit, and access control</p>
		</div>
		<div class="badge badge-lg badge-outline badge-success gap-2">
			<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
				<path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
			</svg>
			Protected
		</div>
	</div>

	<!-- Tabs -->
	<div role="tablist" class="tabs tabs-bordered">
		<button role="tab" class="tab {activeTab === 'overview' ? 'tab-active' : ''}" onclick={() => activeTab = 'overview'}>
			🛡️ Overview
		</button>
		<button role="tab" class="tab {activeTab === 'users' ? 'tab-active' : ''}" onclick={() => activeTab = 'users'}>
			👥 Users ({users.length})
		</button>
		<button role="tab" class="tab {activeTab === 'audit' ? 'tab-active' : ''}" onclick={() => activeTab = 'audit'}>
			📋 Audit Log
		</button>
	</div>

	{#if activeTab === 'overview'}
		<!-- Security overview cards -->
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 animate-slide-up">
			<div class="card bg-base-200 border border-base-300/50 card-hover">
				<div class="card-body p-4">
					<div class="flex items-center justify-between">
						<p class="text-xs text-base-content/60 uppercase tracking-wider">Firewall</p>
						<span class="badge badge-sm {securityMetrics.firewallActive ? 'badge-success' : 'badge-error'}">
							{securityMetrics.firewallActive ? 'Active' : 'Inactive'}
						</span>
					</div>
					<p class="text-2xl font-bold text-success mt-2">UFW</p>
					<p class="text-xs text-base-content/40">Default deny incoming</p>
				</div>
			</div>

			<div class="card bg-base-200 border border-base-300/50 card-hover">
				<div class="card-body p-4">
					<p class="text-xs text-base-content/60 uppercase tracking-wider">SSH Port</p>
					<p class="text-2xl font-bold text-primary mt-2">{securityMetrics.sshPort}</p>
					<p class="text-xs text-base-content/40">Auth: {securityMetrics.authMethod}</p>
				</div>
			</div>

			<div class="card bg-base-200 border border-base-300/50 card-hover">
				<div class="card-body p-4">
					<p class="text-xs text-base-content/60 uppercase tracking-wider">Failed Logins (24h)</p>
					<p class="text-2xl font-bold {securityMetrics.failedLogins > 10 ? 'text-error' : 'text-warning'} mt-2">
						{securityMetrics.failedLogins}
					</p>
					<p class="text-xs text-base-content/40">{securityMetrics.activeSessions} active session(s)</p>
				</div>
			</div>

			<div class="card bg-base-200 border border-base-300/50 card-hover">
				<div class="card-body p-4">
					<p class="text-xs text-base-content/60 uppercase tracking-wider">Sudo Users</p>
					<p class="text-2xl font-bold text-accent mt-2">{securityMetrics.sudoUsers}</p>
					<p class="text-xs text-base-content/40">SSL expires {securityMetrics.sslExpiry}</p>
				</div>
			</div>
		</div>

		<!-- Security checks -->
		<div class="card bg-base-200 border border-base-300/50 animate-slide-up delay-200">
			<div class="card-body p-5">
				<h2 class="card-title text-lg mb-3">Security Checklist</h2>
				<div class="space-y-3">
					{#each [
						{ label: 'Firewall enabled', ok: true, detail: 'UFW is active with default deny' },
						{ label: 'SSH key authentication', ok: false, detail: 'Password authentication is enabled — consider switching to key-only' },
						{ label: 'Root login disabled', ok: true, detail: 'PermitRootLogin is set to no' },
						{ label: 'Fail2Ban active', ok: false, detail: 'Install fail2ban for brute-force protection' },
						{ label: 'System updates', ok: true, detail: 'Unattended upgrades configured' },
						{ label: 'SSL certificate valid', ok: true, detail: 'Expires on ' + securityMetrics.sslExpiry },
					] as check}
						<div class="flex items-start gap-3 p-3 rounded-lg bg-base-300/30 hover:bg-base-300/50 transition-colors">
							<div class="mt-0.5">
								{#if check.ok}
									<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-success" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
										<path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
									</svg>
								{:else}
									<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-warning" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
										<path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
									</svg>
								{/if}
							</div>
							<div class="flex-1">
								<p class="text-sm font-medium {check.ok ? 'text-base-content' : 'text-warning'}">{check.label}</p>
								<p class="text-xs text-base-content/50 mt-0.5">{check.detail}</p>
							</div>
						</div>
					{/each}
				</div>
			</div>
		</div>

	{:else if activeTab === 'users'}
		{#if loading}
			<div class="flex justify-center py-20"><span class="loading loading-spinner loading-lg text-primary"></span></div>
		{:else}
			<div class="overflow-x-auto animate-slide-up">
				<table class="table table-sm w-full">
					<thead>
						<tr class="text-base-content/70">
							<th>Username</th>
							<th>UID</th>
							<th>Home</th>
							<th>Shell</th>
							<th>Groups</th>
						</tr>
					</thead>
					<tbody>
						{#each users as user (user.uid)}
							<tr class="hover:bg-base-300/30 transition-colors">
								<td class="font-mono text-sm font-medium text-primary">{user.username}</td>
								<td class="text-xs">{user.uid}</td>
								<td class="text-xs font-mono text-base-content/60">{user.home}</td>
								<td class="text-xs font-mono text-base-content/60">{user.shell}</td>
								<td>
									<div class="flex flex-wrap gap-1">
										{#each (user.groups || []).slice(0, 5) as group}
											<span class="badge badge-ghost badge-xs">{group}</span>
										{/each}
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}

	{:else if activeTab === 'audit'}
		<div class="card bg-base-200 border border-base-300/50 animate-slide-up">
			<div class="card-body p-5">
				<h2 class="card-title text-lg mb-3">Recent Activity</h2>
				<div class="space-y-2">
					{#each [
						{ time: '08:46:15', event: 'Dashboard accessed', user: 'admin', type: 'info' },
						{ time: '08:45:02', event: 'Service nginx restarted', user: 'admin', type: 'warning' },
						{ time: '08:30:11', event: 'SSH login successful', user: 'admin', type: 'success' },
						{ time: '08:29:45', event: 'SSH login failed', user: 'root', type: 'error' },
						{ time: '07:15:00', event: 'System update applied', user: 'system', type: 'info' },
						{ time: '06:00:00', event: 'Automated backup completed', user: 'system', type: 'success' },
					] as entry}
						<div class="flex items-center gap-3 p-2.5 rounded-lg bg-base-300/30">
							<span class="text-xs font-mono text-base-content/40 w-16">{entry.time}</span>
							<span class="badge badge-xs {entry.type === 'error' ? 'badge-error' : entry.type === 'warning' ? 'badge-warning' : entry.type === 'success' ? 'badge-success' : 'badge-info'}">●</span>
							<span class="text-sm flex-1">{entry.event}</span>
							<span class="text-xs text-base-content/40 font-mono">{entry.user}</span>
						</div>
					{/each}
				</div>
			</div>
		</div>
	{/if}
</div>
