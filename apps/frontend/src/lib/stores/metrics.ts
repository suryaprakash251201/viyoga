import { writable, derived } from 'svelte/store';
import type { MetricsSnapshot, CPUMetrics, MemoryMetrics, DiskMetrics, NetworkMetrics, SystemInfo } from '$lib/types';

// Connection state
export const wsConnected = writable(false);
export const wsError = writable<string | null>(null);

// Latest metrics snapshot
export const metricsSnapshot = writable<MetricsSnapshot | null>(null);

// Metrics history (rolling window for charts)
const MAX_HISTORY = 60;
export const cpuHistory = writable<{ time: string; value: number }[]>([]);
export const memoryHistory = writable<{ time: string; value: number }[]>([]);
export const networkSentHistory = writable<{ time: string; value: number }[]>([]);
export const networkRecvHistory = writable<{ time: string; value: number }[]>([]);

// Derived stores for individual metric types
export const cpuMetrics = derived(metricsSnapshot, ($s) => $s?.cpu as CPUMetrics | null);
export const memoryMetrics = derived(metricsSnapshot, ($s) => $s?.memory as MemoryMetrics | null);
export const diskMetrics = derived(metricsSnapshot, ($s) => $s?.disk as DiskMetrics | null);
export const networkMetrics = derived(metricsSnapshot, ($s) => $s?.network as NetworkMetrics | null);
export const systemInfo = derived(metricsSnapshot, ($s) => $s?.system_info as SystemInfo | null);

// Previous network snapshot for calculating rates
let prevNetworkSent = 0;
let prevNetworkRecv = 0;
let prevTimestamp = 0;

// Network rate stores (bytes/sec)
export const networkSendRate = writable(0);
export const networkRecvRate = writable(0);

function pushHistory(store: typeof cpuHistory, time: string, value: number) {
	store.update((h) => {
		const next = [...h, { time, value }];
		if (next.length > MAX_HISTORY) next.shift();
		return next;
	});
}

function processSnapshot(snap: MetricsSnapshot) {
	metricsSnapshot.set(snap);

	const now = new Date(snap.timestamp);
	const timeStr = now.toLocaleTimeString('en-US', { hour12: false, hour: '2-digit', minute: '2-digit', second: '2-digit' });

	// CPU history
	if (snap.cpu) {
		pushHistory(cpuHistory, timeStr, snap.cpu.usage_percent);
	}

	// Memory history
	if (snap.memory) {
		pushHistory(memoryHistory, timeStr, snap.memory.usage_percent);
	}

	// Network rate calculation
	if (snap.network) {
		const currentTime = now.getTime();
		if (prevTimestamp > 0) {
			const elapsed = (currentTime - prevTimestamp) / 1000;
			if (elapsed > 0) {
				const sentRate = (snap.network.total_bytes_sent - prevNetworkSent) / elapsed;
				const recvRate = (snap.network.total_bytes_recv - prevNetworkRecv) / elapsed;
				networkSendRate.set(Math.max(0, sentRate));
				networkRecvRate.set(Math.max(0, recvRate));
				pushHistory(networkSentHistory, timeStr, Math.max(0, sentRate));
				pushHistory(networkRecvHistory, timeStr, Math.max(0, recvRate));
			}
		}
		prevNetworkSent = snap.network.total_bytes_sent;
		prevNetworkRecv = snap.network.total_bytes_recv;
		prevTimestamp = currentTime;
	}
}

// WebSocket manager
let ws: WebSocket | null = null;
let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
let reconnectAttempts = 0;
const MAX_RECONNECT_DELAY = 30000;

export function connectWebSocket() {
	if (ws?.readyState === WebSocket.OPEN) return;

	const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
	const wsUrl = `${protocol}//${window.location.host}/ws/metrics`;

	try {
		ws = new WebSocket(wsUrl);

		ws.onopen = () => {
			wsConnected.set(true);
			wsError.set(null);
			reconnectAttempts = 0;
		};

		ws.onmessage = (event) => {
			try {
				const snap: MetricsSnapshot = JSON.parse(event.data);
				processSnapshot(snap);
			} catch {
				// Ignore parse errors
			}
		};

		ws.onclose = () => {
			wsConnected.set(false);
			ws = null;
			scheduleReconnect();
		};

		ws.onerror = () => {
			wsError.set('Connection failed');
			wsConnected.set(false);
		};
	} catch {
		wsError.set('Failed to create WebSocket');
		scheduleReconnect();
	}
}

function scheduleReconnect() {
	if (reconnectTimer) return;
	const delay = Math.min(1000 * Math.pow(2, reconnectAttempts), MAX_RECONNECT_DELAY);
	reconnectAttempts++;

	reconnectTimer = setTimeout(() => {
		reconnectTimer = null;
		connectWebSocket();
	}, delay);
}

export function disconnectWebSocket() {
	if (reconnectTimer) {
		clearTimeout(reconnectTimer);
		reconnectTimer = null;
	}
	if (ws) {
		ws.close();
		ws = null;
	}
	wsConnected.set(false);
}
