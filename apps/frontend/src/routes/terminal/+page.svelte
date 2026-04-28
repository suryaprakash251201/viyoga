<script lang="ts">
	import { onMount, tick } from 'svelte';

	interface TerminalLine {
		type: 'prompt' | 'output' | 'error' | 'system';
		text: string;
		timestamp: string;
	}

	let lines: TerminalLine[] = $state([
		{ type: 'system', text: '┌──────────────────────────────────────────────────────┐', timestamp: '' },
		{ type: 'system', text: '│  Viyoga Terminal v0.2.0 — Server Web Console         │', timestamp: '' },
		{ type: 'system', text: '│  Type "help" for available commands                   │', timestamp: '' },
		{ type: 'system', text: '└──────────────────────────────────────────────────────┘', timestamp: '' },
		{ type: 'output', text: '', timestamp: '' }
	]);
	let currentInput = $state('');
	let inputRef: HTMLInputElement;
	let scrollContainer: HTMLDivElement;
	let hostname = $state('viyoga');
	let user = $state('admin');
	let cwd = $state('~');
	let history: string[] = $state([]);
	let historyIdx = $state(-1);

	function timestamp() {
		return new Date().toLocaleTimeString('en-US', { hour12: false });
	}

	function addLine(type: TerminalLine['type'], text: string) {
		lines = [...lines, { type, text, timestamp: timestamp() }];
	}

	async function scrollToBottom() {
		await tick();
		if (scrollContainer) {
			scrollContainer.scrollTop = scrollContainer.scrollHeight;
		}
	}

	async function handleCommand() {
		const cmd = currentInput.trim();
		currentInput = '';
		historyIdx = -1;

		if (!cmd) {
			addLine('prompt', '');
			await scrollToBottom();
			return;
		}

		history = [cmd, ...history.slice(0, 99)];
		addLine('prompt', `${user}@${hostname}:${cwd}$ ${cmd}`);

		const parts = cmd.split(/\s+/);
		const command = parts[0].toLowerCase();
		const args = parts.slice(1);

		switch (command) {
			case 'help':
				addLine('output', 'Available commands:');
				addLine('output', '  help              Show this help message');
				addLine('output', '  clear             Clear terminal');
				addLine('output', '  whoami            Show current user');
				addLine('output', '  hostname          Show server hostname');
				addLine('output', '  uptime            Show system uptime');
				addLine('output', '  date              Show current date/time');
				addLine('output', '  uname             Show system information');
				addLine('output', '  df                Show disk usage');
				addLine('output', '  free              Show memory usage');
				addLine('output', '  ps                Show top processes');
				addLine('output', '  neofetch          Show system summary');
				addLine('output', '  echo [text]       Print text');
				addLine('output', '  history           Show command history');
				addLine('output', '  theme [dark|light] Toggle theme');
				addLine('output', '');
				addLine('system', '💡 Full PTY terminal requires Linux deployment with pty.go');
				break;

			case 'clear':
				lines = [];
				await scrollToBottom();
				return;

			case 'whoami':
				addLine('output', user);
				break;

			case 'hostname':
				addLine('output', hostname);
				break;

			case 'date':
				addLine('output', new Date().toString());
				break;

			case 'uptime': {
				try {
					const res = await fetch('/api/v1/system');
					const json = await res.json();
					if (json.data) {
						const secs = json.data.uptime_seconds;
						const d = Math.floor(secs / 86400);
						const h = Math.floor((secs % 86400) / 3600);
						const m = Math.floor((secs % 3600) / 60);
						addLine('output', ` ${new Date().toLocaleTimeString()} up ${d} days, ${h}:${m.toString().padStart(2, '0')}`);
					}
				} catch {
					addLine('error', 'Failed to fetch uptime');
				}
				break;
			}

			case 'uname':
				try {
					const res = await fetch('/api/v1/system');
					const json = await res.json();
					if (json.data) {
						const d = json.data;
						addLine('output', `${d.os} ${d.hostname} ${d.kernel_version} ${d.kernel_arch} ${d.platform}`);
					}
				} catch {
					addLine('error', 'Failed to fetch system info');
				}
				break;

			case 'df':
				try {
					const res = await fetch('/api/v1/metrics/disk');
					const json = await res.json();
					if (json.data?.partitions) {
						addLine('output', 'Filesystem      Size    Used   Avail  Use%  Mounted on');
						for (const p of json.data.partitions) {
							const size = (p.total_bytes / 1e9).toFixed(1) + 'G';
							const used = (p.used_bytes / 1e9).toFixed(1) + 'G';
							const avail = (p.free_bytes / 1e9).toFixed(1) + 'G';
							const pct = p.used_percent.toFixed(0) + '%';
							addLine('output', `${p.device.padEnd(16)}${size.padStart(6)}  ${used.padStart(6)}  ${avail.padStart(6)}  ${pct.padStart(4)}  ${p.mount_point}`);
						}
					}
				} catch {
					addLine('error', 'Failed to fetch disk info');
				}
				break;

			case 'free':
				try {
					const res = await fetch('/api/v1/metrics/memory');
					const json = await res.json();
					if (json.data) {
						const d = json.data;
						const gb = (b: number) => (b / 1e9).toFixed(1) + 'G';
						addLine('output', '              total       used       free     avail');
						addLine('output', `Mem:     ${gb(d.total_bytes).padStart(10)} ${gb(d.used_bytes).padStart(10)} ${gb(d.free_bytes).padStart(10)} ${gb(d.available_bytes).padStart(10)}`);
						addLine('output', `Swap:    ${gb(d.swap_total_bytes).padStart(10)} ${gb(d.swap_used_bytes).padStart(10)} ${gb(d.swap_free_bytes).padStart(10)}`);
					}
				} catch {
					addLine('error', 'Failed to fetch memory info');
				}
				break;

			case 'ps':
				try {
					const res = await fetch('/api/v1/hardware/processes');
					const json = await res.json();
					if (json.data) {
						addLine('output', '  PID  CPU%   MEM%  NAME');
						for (const p of json.data.slice(0, 15)) {
							addLine('output', `${String(p.pid).padStart(5)}  ${p.cpu_percent.toFixed(1).padStart(5)}  ${p.mem_percent.toFixed(1).padStart(5)}  ${p.name}`);
						}
						addLine('output', `... ${json.data.length} total processes`);
					}
				} catch {
					addLine('error', 'Failed to fetch processes');
				}
				break;

			case 'neofetch': {
				try {
					const res = await fetch('/api/v1/system');
					const json = await res.json();
					if (json.data) {
						const d = json.data;
						const up_d = Math.floor(d.uptime_seconds / 86400);
						const up_h = Math.floor((d.uptime_seconds % 86400) / 3600);
						addLine('system', '        ╭─────────────────────────╮');
						addLine('system', `   🖥️   │  ${d.hostname.padEnd(24)}│`);
						addLine('system', '        ╰─────────────────────────╯');
						addLine('output', `  OS      ${d.platform} ${d.platform_version}`);
						addLine('output', `  Kernel  ${d.kernel_version}`);
						addLine('output', `  Arch    ${d.kernel_arch}`);
						addLine('output', `  Uptime  ${up_d}d ${up_h}h`);
						addLine('output', `  Procs   ${d.procs}`);
						addLine('output', `  Go      ${d.go_version}`);
						addLine('output', `  Viyoga  v${d.viyoga_version}`);
					}
				} catch {
					addLine('error', 'Failed to fetch system info');
				}
				break;
			}

			case 'echo':
				addLine('output', args.join(' '));
				break;

			case 'history':
				history.forEach((cmd, i) => {
					addLine('output', `  ${(i + 1).toString().padStart(4)}  ${cmd}`);
				});
				break;

			case 'theme':
				if (args[0] === 'dark' || args[0] === 'light') {
					document.documentElement.setAttribute('data-theme', args[0] === 'dark' ? 'viyoga' : 'viyoga-light');
					addLine('output', `Theme switched to ${args[0]}`);
				} else {
					addLine('output', 'Usage: theme [dark|light]');
				}
				break;

			default:
				addLine('error', `command not found: ${command}`);
				addLine('output', "Type 'help' for available commands");
		}

		await scrollToBottom();
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'ArrowUp') {
			e.preventDefault();
			if (historyIdx < history.length - 1) {
				historyIdx++;
				currentInput = history[historyIdx];
			}
		} else if (e.key === 'ArrowDown') {
			e.preventDefault();
			if (historyIdx > 0) {
				historyIdx--;
				currentInput = history[historyIdx];
			} else {
				historyIdx = -1;
				currentInput = '';
			}
		}
	}

	function focusInput() {
		inputRef?.focus();
	}

	onMount(async () => {
		try {
			const res = await fetch('/api/v1/system');
			const json = await res.json();
			if (json.data) {
				hostname = json.data.hostname;
			}
		} catch {
			// use default
		}
		focusInput();
		await scrollToBottom();
	});
</script>

<svelte:head>
	<title>Terminal — Viyoga</title>
</svelte:head>

<!-- svelte-ignore a11y_click_events_have_key_events a11y_no_static_element_interactions -->
<div class="h-full flex flex-col animate-fade-in" onclick={focusInput}>
	<!-- Terminal header bar -->
	<div class="flex items-center gap-2 px-4 py-2 bg-base-300/60 border-b border-base-300 rounded-t-xl">
		<div class="flex gap-1.5">
			<div class="w-3 h-3 rounded-full bg-error/70"></div>
			<div class="w-3 h-3 rounded-full bg-warning/70"></div>
			<div class="w-3 h-3 rounded-full bg-success/70"></div>
		</div>
		<span class="flex-1 text-center text-xs font-mono text-base-content/40">
			{user}@{hostname} — viyoga terminal
		</span>
		<button class="btn btn-ghost btn-xs" onclick={() => { lines = []; }}>Clear</button>
	</div>

	<!-- Terminal body -->
	<div
		bind:this={scrollContainer}
		class="flex-1 overflow-y-auto p-4 terminal-container rounded-b-xl font-mono text-sm leading-relaxed"
	>
		{#each lines as line, i}
			<div class="terminal-line" style="animation-delay: {Math.min(i * 20, 200)}ms">
				{#if line.type === 'prompt'}
					<span class="terminal-prompt">{line.text}</span>
				{:else if line.type === 'error'}
					<span class="terminal-error">{line.text}</span>
				{:else if line.type === 'system'}
					<span class="text-primary/80">{line.text}</span>
				{:else}
					<span class="terminal-output">{line.text}</span>
				{/if}
			</div>
		{/each}

		<!-- Input line -->
		<div class="flex items-center">
			<span class="terminal-prompt">{user}@{hostname}:{cwd}$&nbsp;</span>
			<input
				bind:this={inputRef}
				bind:value={currentInput}
				onkeydown={(e) => {
					if (e.key === 'Enter') handleCommand();
					else handleKeydown(e);
				}}
				class="terminal-input flex-1 text-sm"
				spellcheck="false"
				autocomplete="off"
				autocorrect="off"
				autocapitalize="off"
			/>
			<span class="w-2 h-4 bg-primary/80 animate-blink-cursor ml-0.5"></span>
		</div>
	</div>
</div>
