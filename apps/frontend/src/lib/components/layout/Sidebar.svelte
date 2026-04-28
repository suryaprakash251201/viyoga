<script lang="ts">
	import type { NavItem } from '$lib/types';
	import { page } from '$app/state';

	let { collapsed = $bindable(false) } = $props();

	const navItems: NavItem[] = [
		{ label: 'Dashboard', href: '/', icon: 'dashboard' },
		{ label: 'System', href: '/system', icon: 'system' },
		{ label: 'Hardware', href: '/hardware', icon: 'hardware' },
		{ label: 'Containers', href: '/docker', icon: 'docker' },
		{ label: 'DNS Gateway', href: '/dns', icon: 'dns' },
		{ label: 'Monitor', href: '/monitor', icon: 'monitor' },
		{ label: 'Terminal', href: '/terminal', icon: 'terminal' },
		{ label: 'Security', href: '/security', icon: 'security' }
	];

	function isActive(href: string): boolean {
		if (href === '/') return page.url.pathname === '/';
		return page.url.pathname.startsWith(href);
	}
</script>

<!-- SVG icon map -->
{#snippet iconSvg(name: string, active: boolean)}
	<svg xmlns="http://www.w3.org/2000/svg" class="h-[18px] w-[18px] shrink-0 transition-all duration-200 {active ? 'drop-shadow-[0_0_6px_var(--color-primary)]' : ''}" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
		{#if name === 'dashboard'}
			<path stroke-linecap="round" stroke-linejoin="round" d="M4 5a1 1 0 011-1h4a1 1 0 011 1v5a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM14 5a1 1 0 011-1h4a1 1 0 011 1v2a1 1 0 01-1 1h-4a1 1 0 01-1-1V5zM4 15a1 1 0 011-1h4a1 1 0 011 1v2a1 1 0 01-1 1H5a1 1 0 01-1-1v-2zM14 12a1 1 0 011-1h4a1 1 0 011 1v5a1 1 0 01-1 1h-4a1 1 0 01-1-1v-5z" />
		{:else if name === 'system'}
			<path stroke-linecap="round" stroke-linejoin="round" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
		{:else if name === 'hardware'}
			<path stroke-linecap="round" stroke-linejoin="round" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
		{:else if name === 'docker'}
			<path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
		{:else if name === 'dns'}
			<path stroke-linecap="round" stroke-linejoin="round" d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
		{:else if name === 'monitor'}
			<path stroke-linecap="round" stroke-linejoin="round" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
		{:else if name === 'terminal'}
			<path stroke-linecap="round" stroke-linejoin="round" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
		{:else if name === 'security'}
			<path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
		{/if}
	</svg>
{/snippet}

<aside
	class="flex h-screen flex-col border-r border-base-300/40 bg-base-200/95 backdrop-blur-sm transition-all duration-300 ease-in-out"
	class:w-64={!collapsed}
	class:w-16={collapsed}
>
	<!-- Logo -->
	<div class="flex h-16 items-center border-b border-base-300/40 px-4">
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
		<ul class="menu gap-1">
			{#each navItems as item, i}
				<li style="animation-delay: {i * 40}ms">
					<a
						href={item.href}
						class="flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-all duration-200
							{isActive(item.href)
								? 'bg-primary/15 text-primary shadow-sm shadow-primary/10 nav-active'
								: 'text-base-content/60 hover:bg-base-300/50 hover:text-base-content'}
							{collapsed ? 'justify-center px-2' : ''}
							animate-slide-in-left"
						title={collapsed ? item.label : ''}
					>
						{@render iconSvg(item.icon, isActive(item.href))}
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
	<div class="border-t border-base-300/40 p-2.5">
		<button
			onclick={() => (collapsed = !collapsed)}
			class="btn btn-ghost btn-sm w-full rounded-xl hover:bg-base-300/50"
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
