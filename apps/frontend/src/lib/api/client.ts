import type {
	APIResponse,
	SystemdService,
	ServiceDetail,
	LogEntry,
	FirewallStatus,
	SystemUser,
	CronJob,
	ProcessInfo,
	AlertRule,
	AlertEvent,
	DockerContainer,
	DockerImage,
	DNSStats,
	DNSQueryLog,
	DNSBlockList,
	MonitorTargetStatus,
	MonitorTarget
} from '$lib/types';

const API_BASE = '/api/v1';

async function fetchAPI<T>(endpoint: string): Promise<T> {
	const res = await fetch(`${API_BASE}${endpoint}`);
	if (!res.ok) {
		throw new Error(`API error: ${res.status} ${res.statusText}`);
	}
	const json: APIResponse<T> = await res.json();
	if (!json.success) {
		throw new Error(json.error || 'Unknown API error');
	}
	return json.data as T;
}

async function postAPI<T>(endpoint: string, body?: unknown): Promise<T> {
	const res = await fetch(`${API_BASE}${endpoint}`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: body ? JSON.stringify(body) : undefined
	});
	if (!res.ok) {
		throw new Error(`API error: ${res.status} ${res.statusText}`);
	}
	const json: APIResponse<T> = await res.json();
	if (!json.success) {
		throw new Error(json.error || 'Unknown API error');
	}
	return json.data as T;
}

async function deleteAPI<T>(endpoint: string, body?: unknown): Promise<T> {
	const res = await fetch(`${API_BASE}${endpoint}`, {
		method: 'DELETE',
		headers: { 'Content-Type': 'application/json' },
		body: body ? JSON.stringify(body) : undefined
	});
	if (!res.ok) {
		throw new Error(`API error: ${res.status} ${res.statusText}`);
	}
	const json: APIResponse<T> = await res.json();
	if (!json.success) {
		throw new Error(json.error || 'Unknown API error');
	}
	return json.data as T;
}

// ── Phase 1: Metrics ──
export async function fetchCurrentMetrics() {
	return fetchAPI('/metrics/current');
}
export async function fetchCPUMetrics() {
	return fetchAPI('/metrics/cpu');
}
export async function fetchMemoryMetrics() {
	return fetchAPI('/metrics/memory');
}
export async function fetchDiskMetrics() {
	return fetchAPI('/metrics/disk');
}
export async function fetchNetworkMetrics() {
	return fetchAPI('/metrics/network');
}
export async function fetchSystemInfo() {
	return fetchAPI('/system');
}
export async function fetchHealth() {
	return fetchAPI('/health');
}

// ── Phase 2: Linux Management ──
export async function fetchServices() {
	return fetchAPI<SystemdService[]>('/linux/services');
}
export async function fetchServiceDetail(name: string) {
	return fetchAPI<ServiceDetail>(`/linux/services/${name}`);
}
export async function serviceAction(name: string, action: string) {
	return postAPI(`/linux/services/${name}/${action}`);
}
export async function fetchLogs(params?: { unit?: string; priority?: string; since?: string; grep?: string }) {
	const query = new URLSearchParams(params as Record<string, string>).toString();
	return fetchAPI<LogEntry[]>(`/linux/logs?${query}`);
}
export async function fetchFirewallStatus() {
	return fetchAPI<FirewallStatus>('/linux/firewall');
}
export async function addFirewallRule(port: string, proto: string, action: string) {
	return postAPI('/linux/firewall/rules', { port, proto, action });
}
export async function deleteFirewallRule(ruleNumber: number) {
	return deleteAPI('/linux/firewall/rules', { rule_number: ruleNumber });
}
export async function fetchUsers(includeSystem = false) {
	return fetchAPI<SystemUser[]>(`/linux/users?system=${includeSystem}`);
}
export async function fetchCronJobs(user?: string) {
	const q = user ? `?user=${user}` : '';
	return fetchAPI<CronJob[]>(`/linux/cron${q}`);
}

// ── Phase 3: Hardware + Alerting ──
export async function fetchProcesses() {
	return fetchAPI<ProcessInfo[]>('/hardware/processes');
}
export async function fetchAlertRules() {
	return fetchAPI<AlertRule[]>('/hardware/alerts/rules');
}
export async function fetchAlertEvents() {
	return fetchAPI<AlertEvent[]>('/hardware/alerts/events');
}

// ── Phase 4: Docker ──
export async function fetchContainers(all = true) {
	return fetchAPI<DockerContainer[]>(`/docker/containers?all=${all}`);
}
export async function containerAction(id: string, action: string) {
	return postAPI(`/docker/containers/${id}/${action}`);
}
export async function fetchContainerLogs(id: string, lines = 100) {
	return fetchAPI<{ logs: string }>(`/docker/containers/${id}/logs?lines=${lines}`);
}
export async function fetchImages() {
	return fetchAPI<DockerImage[]>('/docker/images');
}
export async function removeImage(id: string) {
	return deleteAPI(`/docker/images/${id}`);
}
export async function pruneDocker(type: string) {
	return postAPI('/docker/prune', { type });
}

// ── Phase 5: DNS Gateway ──
export async function fetchDNSStats() {
	return fetchAPI<DNSStats>('/dns/stats');
}
export async function fetchDNSQueryLog() {
	return fetchAPI<DNSQueryLog[]>('/dns/querylog');
}
export async function fetchDNSBlockLists() {
	return fetchAPI<DNSBlockList[]>('/dns/blocklists');
}

// ── Phase 6: Web Monitor ──
export async function fetchMonitorStatus() {
	return fetchAPI<MonitorTargetStatus[]>('/monitor/status');
}
export async function fetchMonitorTargets() {
	return fetchAPI<MonitorTarget[]>('/monitor/targets');
}
export async function addMonitorTarget(target: Partial<MonitorTarget>) {
	return postAPI('/monitor/targets', target);
}
export async function removeMonitorTarget(id: number) {
	return deleteAPI(`/monitor/targets/${id}`);
}
