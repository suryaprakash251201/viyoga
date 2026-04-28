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
	<svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 shrink-0 transition-all duration-300 {active ? 'text-primary drop-shadow-[0_0_8px_var(--color-primary)]' : 'text-base-content/60'}" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
		{#if name === 'dashboard'}
			<path stroke-linecap="round" stroke-linejoin="round" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
		{:else if name === 'system'}
			<path stroke-linecap="round" stroke-linejoin="round" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01" />
		{:else if name === 'hardware'}
			<path stroke-linecap="round" stroke-linejoin="round" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
		{:else if name === 'docker'}
			<path stroke-linecap="round" stroke-linejoin="round" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
		{:else if name === 'dns'}
			<path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
		{:else if name === 'monitor'}
			<path stroke-linecap="round" stroke-linejoin="round" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
		{:else if name === 'terminal'}
			<path stroke-linecap="round" stroke-linejoin="round" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
		{:else if name === 'security'}
			<path stroke-linecap="round" stroke-linejoin="round" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
		{/if}
	</svg>
{/snippet}

<aside
	class="flex h-screen flex-col border-r border-base-300/50 bg-base-200/80 backdrop-blur-xl transition-all duration-300 ease-in-out shadow-2xl shadow-black/20"
	class:w-64={!collapsed}
	class:w-16={collapsed}
>
	<!-- Logo -->
	<div class="flex h-16 items-center border-b border-base-300/50 px-4">
		{#if !collapsed}
			<div class="flex items-center gap-2.5 animate-fade-in">
				<div class="flex h-9 w-9 items-center justify-center rounded-xl bg-gradient-to-br from-primary to-accent text-primary-content font-bold text-sm shadow-lg shadow-primary/20">V</div>
				<span class="text-xl font-bold tracking-tight">
					<span class="bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent drop-shadow-sm">Vi</span><span class="text-base-content">yoga</span>
				</span>
			</div>
		{:else}
			<div class="mx-auto flex h-9 w-9 items-center justify-center rounded-xl bg-gradient-to-br from-primary to-accent text-primary-content font-bold text-sm shadow-lg shadow-primary/20 animate-scale-in">V</div>
		{/if}
	</div>

	<!-- Navigation -->
	<nav class="flex-1 overflow-y-auto scrollbar-hide p-3">
		<ul class="menu gap-1.5 p-0">
			{#each navItems as item, i}
				<li style="animation-delay: {i * 30}ms">
					<a
						href={item.href}
						class="flex items-center gap-3.5 rounded-xl px-3.5 py-3 text-[14px] font-medium transition-all duration-300 group
							{isActive(item.href)
								? 'bg-gradient-to-r from-primary/10 to-transparent text-primary shadow-sm shadow-primary/5 border-l-2 border-primary'
								: 'text-base-content/60 hover:bg-base-300/40 hover:text-base-content border-l-2 border-transparent'}
							{collapsed ? 'justify-center px-2 border-l-0' : ''}
							animate-slide-in-left"
						title={collapsed ? item.label : ''}
					>
						{@render iconSvg(item.icon, isActive(item.href))}
						{#if !collapsed}
							<span class="flex-1 tracking-wide">{item.label}</span>
							{#if item.badge}
								<span class="badge badge-sm badge-primary ml-auto shadow-sm shadow-primary/20">{item.badge}</span>
							{/if}
						{/if}
					</a>
				</li>
			{/each}
		</ul>
	</nav>

	<!-- Collapse toggle -->
	<div class="border-t border-base-300/50 p-3">
		<button
			onclick={() => (collapsed = !collapsed)}
			class="btn btn-ghost btn-sm w-full rounded-xl hover:bg-base-300/50 text-base-content/50 hover:text-base-content transition-colors"
			title={collapsed ? 'Expand sidebar' : 'Collapse sidebar'}
		>
			{#if collapsed}
				<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 transition-transform duration-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
					<path stroke-linecap="round" stroke-linejoin="round" d="M13 5l7 7-7 7M5 5l7 7-7 7" />
				</svg>
			{:else}
				<span class="flex items-center gap-2 text-xs font-medium">
					<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 transition-transform duration-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
						<path stroke-linecap="round" stroke-linejoin="round" d="M11 19l-7-7 7-7m8 14l-7-7 7-7" />
					</svg>
					Collapse
				</span>
			{/if}
		</button>
	</div>
</aside>
