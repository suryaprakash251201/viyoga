<script lang="ts">
	import './layout.css';
	import Sidebar from '$lib/components/layout/Sidebar.svelte';
	import Header from '$lib/components/layout/Header.svelte';
	import { connectWebSocket, disconnectWebSocket } from '$lib/stores/metrics';
	import { theme } from '$lib/stores/theme';
	import { onMount, onDestroy } from 'svelte';
	import { browser } from '$app/environment';

	let { children } = $props();
	let sidebarCollapsed = $state(false);

	onMount(() => {
		if (browser) {
			connectWebSocket();
			// Ensure theme is applied on mount
			const stored = localStorage.getItem('viyoga-theme');
			if (stored) {
				document.documentElement.setAttribute('data-theme', stored);
			}
		}
	});

	onDestroy(() => {
		if (browser) {
			disconnectWebSocket();
		}
	});
</script>

<svelte:head>
	<title>Viyoga — Server Dashboard</title>
	<meta name="description" content="Self-hosted Ubuntu Server Dashboard" />
</svelte:head>

<div class="flex h-screen overflow-hidden bg-base-100">
	<!-- Sidebar -->
	<Sidebar bind:collapsed={sidebarCollapsed} />

	<!-- Main content area -->
	<div class="flex flex-1 flex-col overflow-hidden">
		<!-- Header -->
		<Header />

		<!-- Page content -->
		<main class="flex-1 overflow-y-auto">
			{@render children()}
		</main>
	</div>
</div>
