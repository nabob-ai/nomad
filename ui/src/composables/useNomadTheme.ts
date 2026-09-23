import { ref, computed, watch, type Ref, type ComputedRef } from 'vue';

/**
 * Nomad UI Protocol Theme System
 *
 * Provides theming support using CSS custom properties.
 * Themes can be applied via:
 * 1. Protocol styles in BeginRenderingMessage
 * 2. Programmatic theme switching
 * 3. CSS class-based theming (dark mode)
 */

/**
 * Theme mode
 */
export type ThemeMode = 'light' | 'dark' | 'system';

/**
 * Theme variables that can be customized
 */
export interface NomadThemeVariables {
  // Primary Colors
  '--nomad-primary'?: string;
  '--nomad-primary-hover'?: string;
  '--nomad-primary-light'?: string;
  '--nomad-primary-dark'?: string;
  '--nomad-primary-contrast'?: string;

  // Secondary Colors
  '--nomad-secondary'?: string;
  '--nomad-secondary-hover'?: string;
  '--nomad-secondary-light'?: string;
  '--nomad-secondary-dark'?: string;
  '--nomad-secondary-contrast'?: string;

  // Surface Colors
  '--nomad-surface'?: string;
  '--nomad-surface-hover'?: string;
  '--nomad-surface-elevated'?: string;

  // Background Colors
  '--nomad-background'?: string;
  '--nomad-background-alt'?: string;

  // Border Colors
  '--nomad-border'?: string;
  '--nomad-border-focus'?: string;

  // Text Colors
  '--nomad-text'?: string;
  '--nomad-text-secondary'?: string;
  '--nomad-text-muted'?: string;
  '--nomad-text-inverse'?: string;

  // Status Colors
  '--nomad-success'?: string;
  '--nomad-success-light'?: string;
  '--nomad-warning'?: string;
  '--nomad-warning-light'?: string;
  '--nomad-error'?: string;
  '--nomad-error-light'?: string;
  '--nomad-info'?: string;
  '--nomad-info-light'?: string;

  // Typography
  '--nomad-font-family'?: string;
  '--nomad-font-family-mono'?: string;
  '--nomad-font-size-xs'?: string;
  '--nomad-font-size-sm'?: string;
  '--nomad-font-size-base'?: string;
  '--nomad-font-size-lg'?: string;
  '--nomad-font-size-xl'?: string;
  '--nomad-font-size-2xl'?: string;
  '--nomad-font-size-3xl'?: string;

  // Spacing
  '--nomad-spacing-xs'?: string;
  '--nomad-spacing-sm'?: string;
  '--nomad-spacing-md'?: string;
  '--nomad-spacing-lg'?: string;
  '--nomad-spacing-xl'?: string;
  '--nomad-spacing-2xl'?: string;

  // Border Radius
  '--nomad-radius-sm'?: string;
  '--nomad-radius-md'?: string;
  '--nomad-radius-lg'?: string;
  '--nomad-radius-xl'?: string;
  '--nomad-radius-full'?: string;

  // Shadows
  '--nomad-shadow-sm'?: string;
  '--nomad-shadow-md'?: string;
  '--nomad-shadow-lg'?: string;
  '--nomad-shadow-xl'?: string;

  // Transitions
  '--nomad-transition-fast'?: string;
  '--nomad-transition-normal'?: string;
  '--nomad-transition-slow'?: string;
  '--nomad-transition-easing'?: string;

  // Allow any custom CSS variable
  [key: `--${string}`]: string | undefined;
}

/**
 * Theme preset definition
 */
export interface NomadThemePreset {
  name: string;
  mode: ThemeMode;
  variables: NomadThemeVariables;
}

/**
 * Default light theme preset
 */
export const LIGHT_THEME: NomadThemePreset = {
  name: 'light',
  mode: 'light',
  variables: {},
};

/**
 * Default dark theme preset
 */
export const DARK_THEME: NomadThemePreset = {
  name: 'dark',
  mode: 'dark',
  variables: {},
};

/**
 * Global theme state
 */
const globalThemeMode = ref<ThemeMode>('system');
const globalCustomVariables = ref<NomadThemeVariables>({});

/**
 * Check if system prefers dark mode
 */
function getSystemPrefersDark(): boolean {
  if (typeof window === 'undefined') return false;
  return window.matchMedia('(prefers-color-scheme: dark)').matches;
}

/**
 * Get effective theme mode (resolves 'system' to actual mode)
 */
function getEffectiveMode(mode: ThemeMode): 'light' | 'dark' {
  if (mode === 'system') {
    return getSystemPrefersDark() ? 'dark' : 'light';
  }
  return mode;
}

/**
 * Apply theme mode to document
 */
function applyThemeModeToDocument(mode: ThemeMode): void {
  if (typeof document === 'undefined') return;

  const effectiveMode = getEffectiveMode(mode);
  const html = document.documentElement;

  if (effectiveMode === 'dark') {
    html.classList.add('dark');
    html.setAttribute('data-theme', 'dark');
  }
  else {
    html.classList.remove('dark');
    html.setAttribute('data-theme', 'light');
  }
}

/**
 * Apply CSS variables to an element
 */
function applyCSSVariables(
  element: HTMLElement,
  variables: Record<string, string>,
): void {
  for (const [key, value] of Object.entries(variables)) {
    if (key.startsWith('--') && value) {
      element.style.setProperty(key, value);
    }
  }
}

/**
 * Remove CSS variables from an element
 */
function removeCSSVariables(
  element: HTMLElement,
  variables: Record<string, string>,
): void {
  for (const key of Object.keys(variables)) {
    if (key.startsWith('--')) {
      element.style.removeProperty(key);
    }
  }
}

/**
 * Convert protocol styles to CSS style object
 * Validates that only CSS custom properties are applied
 */
export function convertProtocolStyles(
  styles: Record<string, string> | undefined,
): Record<string, string> {
  if (!styles) return {};

  const cssVariables: Record<string, string> = {};

  for (const [key, value] of Object.entries(styles)) {
    // Only allow CSS custom properties (starting with --)
    if (key.startsWith('--') && typeof value === 'string') {
      cssVariables[key] = value;
    }
  }

  return cssVariables;
}

/**
 * Composable options
 */
export interface UseNomadThemeOptions {
  /** Initial theme mode */
  initialMode?: ThemeMode;
  /** Initial custom variables */
  initialVariables?: NomadThemeVariables;
  /** Whether to sync with global theme state */
  syncGlobal?: boolean;
}

/**
 * Composable return type
 */
export interface UseNomadThemeReturn {
  /** Current theme mode */
  themeMode: Ref<ThemeMode>;
  /** Effective theme mode (resolved from 'system') */
  effectiveMode: ComputedRef<'light' | 'dark'>;
  /** Whether dark mode is active */
  isDark: ComputedRef<boolean>;
  /** Custom theme variables */
  customVariables: Ref<NomadThemeVariables>;
  /** Set theme mode */
  setThemeMode: (mode: ThemeMode) => void;
  /** Set custom variables */
  setCustomVariables: (variables: NomadThemeVariables) => void;
  /** Apply theme to an element */
  applyTheme: (element: HTMLElement, variables?: Record<string, string>) => void;
  /** Remove theme from an element */
  removeTheme: (element: HTMLElement, variables?: Record<string, string>) => void;
  /** Toggle between light and dark mode */
  toggleTheme: () => void;
  /** Convert protocol styles to CSS variables */
  convertStyles: typeof convertProtocolStyles;
}

/**
 * Nomad Theme Composable
 *
 * Provides theme management for Nomad UI components.
 *
 * @example
 * ```ts
 * const { themeMode, isDark, setThemeMode, applyTheme } = useNomadTheme();
 *
 * // Set theme mode
 * setThemeMode('dark');
 *
 * // Apply custom variables to an element
 * applyTheme(element, { '--nomad-primary': '#ff0000' });
 * ```
 */
export function useNomadTheme(options: UseNomadThemeOptions = {}): UseNomadThemeReturn {
  const {
    initialMode = 'system',
    initialVariables = {},
    syncGlobal = true,
  } = options;

  // Local state (synced with global if syncGlobal is true)
  const themeMode = syncGlobal ? globalThemeMode : ref<ThemeMode>(initialMode);
  const customVariables = syncGlobal ? globalCustomVariables : ref<NomadThemeVariables>(initialVariables);

  // Initialize if not syncing global
  if (!syncGlobal) {
    themeMode.value = initialMode;
    customVariables.value = initialVariables;
  }

  // Computed effective mode
  const effectiveMode = computed(() => getEffectiveMode(themeMode.value));

  // Computed isDark
  const isDark = computed(() => effectiveMode.value === 'dark');

  // Watch theme mode changes and apply to document
  watch(
    themeMode,
    (mode) => {
      applyThemeModeToDocument(mode);
    },
    { immediate: true },
  );

  // Listen for system theme changes
  if (typeof window !== 'undefined') {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    const handleChange = () => {
      if (themeMode.value === 'system') {
        applyThemeModeToDocument('system');
      }
    };
    mediaQuery.addEventListener('change', handleChange);
  }

  /**
   * Set theme mode
   */
  function setThemeMode(mode: ThemeMode): void {
    themeMode.value = mode;
  }

  /**
   * Set custom variables
   */
  function setCustomVariables(variables: NomadThemeVariables): void {
    customVariables.value = { ...customVariables.value, ...variables };
  }

  /**
   * Apply theme to an element
   */
  function applyTheme(element: HTMLElement, variables?: Record<string, string>): void {
    const varsToApply = variables ?? customVariables.value;
    applyCSSVariables(element, varsToApply as Record<string, string>);
  }

  /**
   * Remove theme from an element
   */
  function removeTheme(element: HTMLElement, variables?: Record<string, string>): void {
    const varsToRemove = variables ?? customVariables.value;
    removeCSSVariables(element, varsToRemove as Record<string, string>);
  }

  /**
   * Toggle between light and dark mode
   */
  function toggleTheme(): void {
    if (themeMode.value === 'system') {
      // If system, switch to opposite of current effective mode
      themeMode.value = effectiveMode.value === 'dark' ? 'light' : 'dark';
    }
    else {
      themeMode.value = themeMode.value === 'dark' ? 'light' : 'dark';
    }
  }

  return {
    themeMode,
    effectiveMode,
    isDark,
    customVariables,
    setThemeMode,
    setCustomVariables,
    applyTheme,
    removeTheme,
    toggleTheme,
    convertStyles: convertProtocolStyles,
  };
}

/**
 * Create a theme preset
 */
export function createThemePreset(
  name: string,
  mode: ThemeMode,
  variables: NomadThemeVariables,
): NomadThemePreset {
  return { name, mode, variables };
}

/**
 * Merge theme variables
 */
export function mergeThemeVariables(
  ...themes: (NomadThemeVariables | undefined)[]
): NomadThemeVariables {
  return themes.reduce<NomadThemeVariables>((acc, theme) => {
    if (theme) {
      return { ...acc, ...theme };
    }
    return acc;
  }, {});
}
