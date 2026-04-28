/**
 * Format bytes to human-readable string.
 */
export function formatBytes(bytes: number, decimals = 1): string {
	if (bytes === 0) return '0 B';
	const k = 1024;
	const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'];
	const i = Math.floor(Math.log(bytes) / Math.log(k));
	return `${parseFloat((bytes / Math.pow(k, i)).toFixed(decimals))} ${sizes[i]}`;
}

/**
 * Format seconds to human-readable uptime string.
 */
export function formatUptime(seconds: number): string {
	const days = Math.floor(seconds / 86400);
	const hours = Math.floor((seconds % 86400) / 3600);
	const mins = Math.floor((seconds % 3600) / 60);

	const parts: string[] = [];
	if (days > 0) parts.push(`${days}d`);
	if (hours > 0) parts.push(`${hours}h`);
	parts.push(`${mins}m`);

	return parts.join(' ');
}

/**
 * Format percentage with optional decimal places.
 */
export function formatPercent(value: number, decimals = 1): string {
	return `${value.toFixed(decimals)}%`;
}

/**
 * Format a number with SI suffixes (K, M, G).
 */
export function formatNumber(num: number): string {
	if (num >= 1e9) return `${(num / 1e9).toFixed(1)}G`;
	if (num >= 1e6) return `${(num / 1e6).toFixed(1)}M`;
	if (num >= 1e3) return `${(num / 1e3).toFixed(1)}K`;
	return num.toString();
}

/**
 * Get color class based on usage percentage thresholds.
 */
export function getUsageColor(percent: number): string {
	if (percent >= 90) return 'text-error';
	if (percent >= 70) return 'text-warning';
	if (percent >= 50) return 'text-info';
	return 'text-success';
}

/**
 * Get color hex based on usage percentage thresholds.
 */
export function getUsageHex(percent: number): string {
	if (percent >= 90) return '#ef4444';
	if (percent >= 70) return '#f59e0b';
	if (percent >= 50) return '#3b82f6';
	return '#22c55e';
}

/**
 * Clamp a number between min and max.
 */
export function clamp(value: number, min: number, max: number): number {
	return Math.min(Math.max(value, min), max);
}
