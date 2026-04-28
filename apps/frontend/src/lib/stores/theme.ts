import { writable } from 'svelte/store';
import { browser } from '$app/environment';

type Theme = 'viyoga' | 'viyoga-light';

function createThemeStore() {
	const stored = browser ? (localStorage.getItem('viyoga-theme') as Theme) : null;
	const initial: Theme = stored || 'viyoga';

	const { subscribe, set, update } = writable<Theme>(initial);

	function applyTheme(theme: Theme) {
		if (browser) {
			document.documentElement.setAttribute('data-theme', theme);
			localStorage.setItem('viyoga-theme', theme);
		}
	}

	// Apply initial theme
	if (browser) {
		applyTheme(initial);
	}

	return {
		subscribe,
		toggle: () => {
			update((current) => {
				const next: Theme = current === 'viyoga' ? 'viyoga-light' : 'viyoga';
				applyTheme(next);
				return next;
			});
		},
		set: (theme: Theme) => {
			applyTheme(theme);
			set(theme);
		}
	};
}

export const theme = createThemeStore();
export const isDark = {
	subscribe(fn: (value: boolean) => void) {
		return theme.subscribe((t) => fn(t === 'viyoga'));
	}
};
