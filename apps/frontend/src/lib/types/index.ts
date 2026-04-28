// System metrics types

export interface CPUMetrics {
	usage_percent: number;
	per_core: number[];
	core_count: number;
	logical_count: number;
	model_name: string;
	frequency_mhz: number;
	load_avg_1: number;
	load_avg_5: number;
	load_avg_15: number;
}

export interface MemoryMetrics {
	total_bytes: number;
	used_bytes: number;
	available_bytes: number;
	free_bytes: number;
	usage_percent: number;
	cached_bytes: number;
	buffers_bytes: number;
	swap_total_bytes: number;
	swap_used_bytes: number;
	swap_free_bytes: number;
	swap_usage_percent: number;
}

export interface DiskPartition {
	device: string;
	mount_point: string;
	fs_type: string;
	total_bytes: number;
	used_bytes: number;
	free_bytes: number;
	used_percent: number;
}

export interface DiskIOStats {
	device: string;
	read_bytes: number;
	write_bytes: number;
	read_count: number;
	write_count: number;
	read_time_ms: number;
	write_time_ms: number;
}

export interface DiskMetrics {
	partitions: DiskPartition[];
	io: DiskIOStats[];
}

export interface NetworkInterface {
	name: string;
	bytes_sent: number;
	bytes_recv: number;
	packets_sent: number;
	packets_recv: number;
	err_in: number;
	err_out: number;
	drop_in: number;
	drop_out: number;
}

export interface ConnectionSummary {
	tcp: number;
	udp: number;
	listening: number;
	established: number;
}

export interface NetworkMetrics {
	interfaces: NetworkInterface[];
	connections: ConnectionSummary;
	total_bytes_sent: number;
	total_bytes_recv: number;
}

export interface SystemInfo {
	hostname: string;
	os: string;
	platform: string;
	platform_version: string;
	platform_family: string;
	kernel_version: string;
	kernel_arch: string;
	uptime_seconds: number;
	boot_time: number;
	procs: number;
	go_version: string;
	viyoga_version: string;
}

export interface MetricsSnapshot {
	timestamp: string;
	cpu: CPUMetrics | null;
	memory: MemoryMetrics | null;
	disk: DiskMetrics | null;
	network: NetworkMetrics | null;
	system_info: SystemInfo | null;
}

export interface APIResponse<T> {
	success: boolean;
	data?: T;
	error?: string;
	meta?: {
		timestamp: string;
		ws_clients: number;
	};
}

// Navigation types
export interface NavItem {
	label: string;
	href: string;
	icon: string;
	badge?: string;
	disabled?: boolean;
}

// ── Phase 2: Linux Management ──

export interface SystemdService {
	name: string;
	description: string;
	load_state: string;
	active_state: string;
	sub_state: string;
	enabled: string;
}

export interface ServiceDetail extends SystemdService {
	main_pid: number;
	memory_current: number;
	cpu_usage: string;
	fragment_path: string;
	restart: string;
	type: string;
}

export interface LogEntry {
	timestamp: string;
	hostname: string;
	unit: string;
	message: string;
	priority: string;
	pid: string;
}

export interface FirewallStatus {
	active: boolean;
	default_policy: string;
	rules: FirewallRule[];
}

export interface FirewallRule {
	number: number;
	to: string;
	action: string;
	from: string;
	direction: string;
	v6: boolean;
}

export interface SystemUser {
	username: string;
	uid: number;
	gid: number;
	comment: string;
	home: string;
	shell: string;
	groups: string[];
	is_system: boolean;
}

export interface CronJob {
	schedule: string;
	command: string;
	user: string;
	raw: string;
}

// ── Phase 3: Hardware Monitor + Alerting ──

export interface ProcessInfo {
	pid: number;
	name: string;
	username: string;
	cpu_percent: number;
	mem_percent: number;
	mem_rss: number;
	status: string;
	create_time: number;
	cmdline: string;
	ppid: number;
	num_threads: number;
}

export interface AlertRule {
	id: number;
	name: string;
	metric_type: string;
	condition: string;
	threshold: number;
	notify_channel: string;
	notify_target: string;
	enabled: boolean;
	cooldown_mins: number;
}

export interface AlertEvent {
	id: number;
	rule_id: number;
	rule_name: string;
	metric_type: string;
	value: number;
	threshold: number;
	triggered_at: string;
	resolved_at?: string;
	acknowledged: boolean;
}

// ── Phase 4: Docker ──

export interface DockerContainer {
	id: string;
	name: string;
	image: string;
	state: string;
	status: string;
	created: number;
	ports: PortMapping[];
	labels: Record<string, string>;
}

export interface PortMapping {
	private: number;
	public: number;
	type: string;
	ip: string;
}

export interface DockerImage {
	id: string;
	tags: string[];
	size: number;
	created: number;
}

// ── Phase 5: DNS Gateway ──

export interface DNSStats {
	total_queries: number;
	total_blocked_queries: number;
	blocked_percent: number;
	total_clients: number;
	top_domains: { domain: string; hits: number }[];
	top_blocked_domains: { domain: string; hits: number }[];
}

export interface DNSQueryLog {
	timestamp: string;
	client_ip: string;
	domain: string;
	type: string;
	response_code: string;
	blocked: boolean;
}

export interface DNSBlockList {
	name: string;
	url: string;
	enabled: boolean;
	entries: number;
}

// ── Phase 6: Web Monitor ──

export interface MonitorTarget {
	id: number;
	name: string;
	url: string;
	method: string;
	interval_seconds: number;
	timeout_seconds: number;
	expected_status: number;
	enabled: boolean;
}

export interface MonitorResult {
	target_id: number;
	target_name: string;
	url: string;
	status_code: number;
	response_time_ms: number;
	is_up: boolean;
	error?: string;
	checked_at: string;
	cert_expiry?: string;
}

export interface MonitorTargetStatus {
	target: MonitorTarget;
	last_result?: MonitorResult;
	uptime_percent: number;
	avg_response_ms: number;
	history: MonitorResult[];
}

