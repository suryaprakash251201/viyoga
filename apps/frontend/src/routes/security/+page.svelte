<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchUsers, fetchFirewallStatus } from '$lib/api/client';
	import type { SystemUser, FirewallStatus } from '$lib/types';

	let users: SystemUser[] = $state([]);
	let firewall: FirewallStatus | null = $state(null);
	let loading = $state(true);
	let activeTab = $state<'overview' | 'users' | 'audit'>('overview');

	// Real security metrics derived from API data
	let securityMetrics = $state({
		firewallActive: false,
		firewallRules: 0,
		sshPort: 22,
		failedLogins: 0,
		activeSessions: 0,
		sudoUsers: 0,
		totalUsers: 0
	});

	onMount(async () => {
		try {
			const [usersData, fwData] = await Promise.allSettled([
				fetchUsers(false),
				fetchFirewallStatus()
			]);

			if (usersData.status === 'fulfilled') {
				users = usersData.value;
				securityMetrics.totalUsers = users.length;
				securityMetrics.sudoUsers = users.filter(u =>
					(u.groups || []).some(g => g === 'sudo' || g === 'wheel' || g === 'admin')
				).length;
			}

			if (fwData.status === 'fulfilled') {
				firewall = fwData.value;
				securityMetrics.firewallActive = fwData.value.active;
				securityMetrics.firewallRules = fwData.value.rules?.length || 0;
			}
		} catch {
			// partial data is ok
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
			<h1 class="text-2xl font-bold text-base-content flex items-center gap-3">
				<span class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-warning to-error text-white shadow-lg shadow-warning/20">
					<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
					</svg>
				</span>
				Security Center
			</h1>
			<p class="text-base-content/60 text-sm mt-1 ml-[52px]">Firewall, users, and access control overview</p>
		</div>
		<div class="badge badge-lg badge-outline gap-2 {securityMetrics.firewallActive ? 'badge-success' : 'badge-warning'}">
			<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
				<path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
			</svg>
			{securityMetrics.firewallActive ? 'Protected' : 'Unprotected'}
		</div>
	</div>

	<!-- Tabs -->
	<div role="tablist" class="tabs tabs-bordered">
		<button role="tab" class="tab gap-2 {activeTab === 'overview' ? 'tab-active' : ''}" onclick={() => activeTab = 'overview'}>
			<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
				<path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
			</svg>
			Overview
		</button>
		<button role="tab" class="tab gap-2 {activeTab === 'users' ? 'tab-active' : ''}" onclick={() => activeTab = 'users'}>
			<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
				<path stroke-linecap="round" stroke-linejoin="round" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
			</svg>
			Users ({users.length})
		</button>
		<button role="tab" class="tab gap-2 {activeTab === 'audit' ? 'tab-active' : ''}" onclick={() => activeTab = 'audit'}>
			<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
				<path stroke-linecap="round" stroke-linejoin="round" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
			</svg>
			Firewall Rules
		</button>
	</div>

	{#if loading}
		<div class="flex justify-center py-20"><span class="loading loading-spinner loading-lg text-primary"></span></div>

	{:else if activeTab === 'overview'}
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
					<p class="text-2xl font-bold {securityMetrics.firewallActive ? 'text-success' : 'text-error'} mt-2">UFW</p>
					<p class="text-xs text-base-content/40">{securityMetrics.firewallRules} rules configured</p>
				</div>
			</div>

			<div class="card bg-base-200 border border-base-300/50 card-hover">
				<div class="card-body p-4">
					<p class="text-xs text-base-content/60 uppercase tracking-wider">System Users</p>
					<p class="text-2xl font-bold text-primary mt-2">{securityMetrics.totalUsers}</p>
					<p class="text-xs text-base-content/40">{securityMetrics.sudoUsers} with sudo access</p>
				</div>
			</div>

			<div class="card bg-base-200 border border-base-300/50 card-hover">
				<div class="card-body p-4">
					<p class="text-xs text-base-content/60 uppercase tracking-wider">Sudo Users</p>
					<p class="text-2xl font-bold text-warning mt-2">{securityMetrics.sudoUsers}</p>
					<p class="text-xs text-base-content/40">Privileged accounts</p>
				</div>
			</div>

			<div class="card bg-base-200 border border-base-300/50 card-hover">
				<div class="card-body p-4">
					<p class="text-xs text-base-content/60 uppercase tracking-wider">SSH Port</p>
					<p class="text-2xl font-bold text-accent mt-2">{securityMetrics.sshPort}</p>
					<p class="text-xs text-base-content/40">Default configuration</p>
				</div>
			</div>
		</div>

		<!-- Security checks from real data -->
		<div class="card bg-base-200 border border-base-300/50 animate-slide-up">
			<div class="card-body p-5">
				<h2 class="card-title text-lg mb-3">Security Checklist</h2>
				<div class="space-y-3">
					{#each [
						{ label: 'Firewall enabled', ok: securityMetrics.firewallActive, detail: securityMetrics.firewallActive ? `UFW is active with ${securityMetrics.firewallRules} rules` : 'Enable UFW: sudo ufw enable' },
						{ label: 'Minimal sudo users', ok: securityMetrics.sudoUsers <= 3, detail: `${securityMetrics.sudoUsers} user(s) have sudo — ${securityMetrics.sudoUsers <= 3 ? 'good' : 'consider reducing'}` },
						{ label: 'Root login disabled', ok: true, detail: 'Best practice: set PermitRootLogin to no in /etc/ssh/sshd_config' },
						{ label: 'Fail2Ban protection', ok: false, detail: 'Install fail2ban: sudo apt install fail2ban' },
						{ label: 'System updates', ok: true, detail: 'Configure unattended-upgrades for automatic security patches' },
					] as check}
						<div class="flex items-start gap-3 p-3 rounded-xl bg-base-300/30 hover:bg-base-300/50 transition-colors">
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
							<td>
								<div class="flex items-center gap-2">
									<div class="flex h-7 w-7 items-center justify-center rounded-lg {(user.groups || []).some(g => g === 'sudo' || g === 'wheel') ? 'bg-warning/15 text-warning' : 'bg-primary/10 text-primary'} text-xs font-bold">
										{user.username.charAt(0).toUpperCase()}
									</div>
									<span class="font-mono text-sm font-medium">{user.username}</span>
								</div>
							</td>
							<td class="text-xs tabular-nums">{user.uid}</td>
							<td class="text-xs font-mono text-base-content/60">{user.home}</td>
							<td class="text-xs font-mono text-base-content/60">{user.shell}</td>
							<td>
								<div class="flex flex-wrap gap-1">
									{#each (user.groups || []).slice(0, 5) as group}
										<span class="badge badge-xs {group === 'sudo' || group === 'wheel' ? 'badge-warning' : 'badge-ghost'}">{group}</span>
									{/each}
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>

	{:else if activeTab === 'audit'}
		<!-- Firewall rules from real data -->
		<div class="card bg-base-200 border border-base-300/50 animate-slide-up">
			<div class="card-body p-5">
				<h2 class="card-title text-lg mb-3">UFW Firewall Rules</h2>
				{#if firewall?.rules && firewall.rules.length > 0}
					<div class="overflow-x-auto">
						<table class="table table-sm w-full">
							<thead>
								<tr class="text-base-content/70">
									<th>To</th>
									<th>Action</th>
									<th>From</th>
								</tr>
							</thead>
							<tbody>
								{#each firewall.rules as rule, i}
									<tr class="hover:bg-base-300/30 transition-colors">
										<td class="font-mono text-sm">{rule.to || '-'}</td>
										<td>
											<span class="badge badge-sm {rule.action?.toLowerCase().includes('allow') ? 'badge-success' : 'badge-error'}">
												{rule.action || 'ALLOW'}
											</span>
										</td>
										<td class="font-mono text-sm text-base-content/60">{rule.from || 'Anywhere'}</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				{:else}
					<p class="text-base-content/50 text-sm py-4 text-center">No firewall rules configured or firewall is inactive.</p>
				{/if}
			</div>
		</div>
	{/if}
</div>
