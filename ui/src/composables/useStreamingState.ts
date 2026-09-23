/**
 * Nomad UI Protocol - Streaming State Preservation
 *
 * Composable for preserving UI state during incremental updates.
 * Maintains scroll position and input focus across component tree rebuilds.
 *
 * @module composables/useStreamingState
 */

import { ref, onMounted, onBeforeUnmount, nextTick, type Ref } from 'vue';

/**
 * Preserved state for a surface
 */
export interface PreservedState {
  /** Scroll positions by element selector or ID */
  scrollPositions: Map<string, { top: number; left: number }>;
  /** Currently focused element ID or selector */
  focusedElementId: string | null;
  /** Input values that may not have been synced yet */
  pendingInputValues: Map<string, string>;
}

/**
 * Create initial preserved state
 */
function createPreservedState(): PreservedState {
  return {
    scrollPositions: new Map(),
    focusedElementId: null,
    pendingInputValues: new Map(),
  };
}

/**
 * Global state store for all surfaces
 */
const surfaceStates = new Map<string, PreservedState>();

/**
 * Get or create preserved state for a surface
 */
export function getPreservedState(surfaceId: string): PreservedState {
  let state = surfaceStates.get(surfaceId);
  if (!state) {
    state = createPreservedState();
    surfaceStates.set(surfaceId, state);
  }
  return state;
}

/**
 * Clear preserved state for a surface
 */
export function clearPreservedState(surfaceId: string): void {
  surfaceStates.delete(surfaceId);
}

/**
 * Options for useStreamingState composable
 */
export interface UseStreamingStateOptions {
  /** Surface ID */
  surfaceId: string;
  /** Container element ref */
  containerRef: Ref<HTMLElement | null>;
  /** Whether to preserve scroll position */
  preserveScroll?: boolean;
  /** Whether to preserve input focus */
  preserveFocus?: boolean;
}

/**
 * Composable for preserving UI state during streaming updates
 *
 * @example
 * ```vue
 * <script setup>
 * const containerRef = ref<HTMLElement | null>(null);
 * const { saveState, restoreState } = useStreamingState({
 *   surfaceId: 'my-surface',
 *   containerRef,
 * });
 *
 * // Before update
 * saveState();
 *
 * // After update
 * await nextTick();
 * restoreState();
 * </script>
 * ```
 */
export function useStreamingState(options: UseStreamingStateOptions) {
  const {
    surfaceId,
    containerRef,
    preserveScroll = true,
    preserveFocus = true,
  } = options;

  const isRestoring = ref(false);

  /**
   * Save current scroll positions within the container
   */
  function saveScrollPositions(): void {
    if (!preserveScroll || !containerRef.value) return;

    const state = getPreservedState(surfaceId);
    state.scrollPositions.clear();

    // Save container scroll position
    state.scrollPositions.set('__container__', {
      top: containerRef.value.scrollTop,
      left: containerRef.value.scrollLeft,
    });

    // Save scroll positions of scrollable children
    const scrollableElements = containerRef.value.querySelectorAll('[data-scroll-preserve]');
    scrollableElements.forEach((el) => {
      const id = el.getAttribute('data-scroll-preserve') || el.id;
      if (id) {
        state.scrollPositions.set(id, {
          top: (el as HTMLElement).scrollTop,
          left: (el as HTMLElement).scrollLeft,
        });
      }
    });
  }

  /**
   * Restore saved scroll positions
   */
  function restoreScrollPositions(): void {
    if (!preserveScroll || !containerRef.value) return;

    const state = getPreservedState(surfaceId);

    // Restore container scroll position
    const containerScroll = state.scrollPositions.get('__container__');
    if (containerScroll) {
      containerRef.value.scrollTop = containerScroll.top;
      containerRef.value.scrollLeft = containerScroll.left;
    }

    // Restore scroll positions of scrollable children
    state.scrollPositions.forEach((pos, id) => {
      if (id === '__container__') return;

      const el = containerRef.value?.querySelector(`[data-scroll-preserve="${id}"]`)
        || containerRef.value?.querySelector(`#${id}`);
      if (el) {
        (el as HTMLElement).scrollTop = pos.top;
        (el as HTMLElement).scrollLeft = pos.left;
      }
    });
  }

  /**
   * Save currently focused element
   */
  function saveFocusState(): void {
    if (!preserveFocus || !containerRef.value) return;

    const state = getPreservedState(surfaceId);
    const activeElement = document.activeElement;

    if (activeElement && containerRef.value.contains(activeElement)) {
      // Try to get a unique identifier for the element
      const id = activeElement.id
        || activeElement.getAttribute('data-component-id')
        || activeElement.getAttribute('name');

      if (id) {
        state.focusedElementId = id;

        // Save pending input value if it's an input element
        if (activeElement instanceof HTMLInputElement || activeElement instanceof HTMLTextAreaElement) {
          state.pendingInputValues.set(id, activeElement.value);
        }
      }
    } else {
      state.focusedElementId = null;
    }
  }

  /**
   * Restore focus to previously focused element
   */
  function restoreFocusState(): void {
    if (!preserveFocus || !containerRef.value) return;

    const state = getPreservedState(surfaceId);

    if (state.focusedElementId) {
      const el = containerRef.value.querySelector(`#${state.focusedElementId}`)
        || containerRef.value.querySelector(`[data-component-id="${state.focusedElementId}"]`)
        || containerRef.value.querySelector(`[name="${state.focusedElementId}"]`);

      if (el && el instanceof HTMLElement) {
        // Use requestAnimationFrame to ensure DOM is ready
        requestAnimationFrame(() => {
          el.focus();

          // Restore cursor position for input elements
          if (el instanceof HTMLInputElement || el instanceof HTMLTextAreaElement) {
            const pendingValue = state.pendingInputValues.get(state.focusedElementId!);
            if (pendingValue !== undefined) {
              // Set selection to end of input
              const len = el.value.length;
              el.setSelectionRange(len, len);
            }
          }
        });
      }
    }
  }

  /**
   * Save all state before an update
   */
  function saveState(): void {
    saveScrollPositions();
    saveFocusState();
  }

  /**
   * Restore all state after an update
   */
  async function restoreState(): Promise<void> {
    if (isRestoring.value) return;

    isRestoring.value = true;

    try {
      await nextTick();
      restoreScrollPositions();
      restoreFocusState();
    } finally {
      isRestoring.value = false;
    }
  }

  /**
   * Clear all preserved state
   */
  function clearState(): void {
    clearPreservedState(surfaceId);
  }

  // Cleanup on unmount
  onBeforeUnmount(() => {
    clearState();
  });

  return {
    saveState,
    restoreState,
    clearState,
    isRestoring,
    // Expose individual functions for fine-grained control
    saveScrollPositions,
    restoreScrollPositions,
    saveFocusState,
    restoreFocusState,
  };
}

/**
 * Hook to mark an element as having preservable scroll position
 *
 * @example
 * ```vue
 * <div ref="scrollRef" v-scroll-preserve="'my-list'">
 *   <!-- scrollable content -->
 * </div>
 * ```
 */
export function useScrollPreserve(elementRef: Ref<HTMLElement | null>, id: string) {
  onMounted(() => {
    if (elementRef.value) {
      elementRef.value.setAttribute('data-scroll-preserve', id);
    }
  });
}
