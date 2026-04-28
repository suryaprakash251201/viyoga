<script lang="ts">
	import type { NavItem } from '$lib/types';
	import { page } from '$app/state';

	let { collapsed = $bindable(false) } = $props();

	const navItems: NavItem[] = [
		{ label: 'Dashboard', href: '/', icon: '📊' },
		{ label: 'System', href: '/system', icon: '🖥️' },
		{ label: 'Hardware', href: '/hardware', icon: '🔧' },
		{ label: 'Containers', href: '/docker', icon: '🐳' },
		{ label: 'DNS Gateway', href: '/dns', icon: '🛡️' },
		{ label: 'Monitor', href: '/monitor', icon: '📡' },
		{ label: 'Terminal', href: '/terminal', icon: '💻' },
		{ label: 'Security', href: '/security', icon: '🔒' }
	];

	function isActive(href: string): boolean {
		if (href === '/') return page.url.pathname === '/';
		return page.url.pathname.startsWith(href);
	}
</script>

<aside
	class="flex h-screen flex-col border-r border-base-300/60 bg-base-200 transition-all duration-300 ease-in-out"
	class:w-64={!collapsed}
	class:w-16={collapsed}
>
	<!-- Logo -->
	<div class="flex h-16 items-center border-b border-base-300/60 px-4">
		{#if !collapsed}
			<div class="flex items-center gap-2.5 animate-fade-in">
				<div class="flex h-9 w-9 items-center justify-center rounded-xl bg-gradient-to-br from-primary to-accent text-primary-content font-bold text-sm shadow-lg shadow-primary/20">V</div>
				<span class="text-lg font-bold tracking-tight">
					<span class="bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent">Vi</span><span class="text-base-content">yoga</span>
				</span>
			</div>
		{:else}
			<div class="mx-auto flex h-9 w-9 items-center justify-center rounded-xl bg-gradient-to-br from-primary to-accent text-primary-content font-bold text-sm shadow-lg shadow-primary/20 animate-scale-in">V</div>
		{/if}
	</div>

	<!-- Navigation -->
	<nav class="flex-1 overflow-y-auto p-2.5">
		<ul class="menu gap-0.5">
			{#each navItems as item, i}
				<li style="animation-delay: {i * 40}ms">
					<a
						href={item.href}
						class="flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-all duration-200
							{isActive(item.href)
								? 'bg-primary/15 text-primary shadow-sm shadow-primary/10 nav-active'
								: 'text-base-content/70 hover:bg-base-300/60 hover:text-base-content'}
							{collapsed ? 'justify-center px-2' : ''}
							animate-slide-in-left"
						title={collapsed ? item.label : ''}
					>
						<span class="text-lg transition-transform duration-200 {isActive(item.href) ? 'scale-110' : 'group-hover:scale-105'}">{item.icon}</span>
						{#if !collapsed}
							<span class="flex-1">{item.label}</span>
							{#if item.badge}
								<span class="badge badge-sm badge-primary ml-auto">{item.badge}</span>
							{/if}
						{/if}
					</a>
				</li>
			{/each}
		</ul>
	</nav>

	<!-- Collapse toggle -->
	<div class="border-t border-base-300/60 p-2.5">
		<button
			onclick={() => (collapsed = !collapsed)}
			class="btn btn-ghost btn-sm w-full rounded-xl hover:bg-base-300/60"
			title={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
		>
			{#if collapsed}
				<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 transition-transform duration-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
					<path stroke-linecap="round" stroke-linejoin="round" d="M13 5l7 7-7 7M5 5l7 7-7 7" />
				</svg>
			{:else}
				<span class="flex items-center gap-2 text-xs text-base-content/50">
					<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 transition-transform duration-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M11 19l-7-7 7-7m8 14l-7-7 7-7" />
					</svg>
					Collapse
				</span>
			{/if}
		</button>
	</div>
</aside>
