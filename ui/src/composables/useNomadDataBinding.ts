/**
 * Nomad UI Protocol - Data Binding Composable
 *
 * Provides reactive data binding for Nomad UI components.
 * Implements two-way binding between UI components and the data model.
 *
 * @module composables/useNomadDataBinding
 */

import { ref, computed, watch, inject, type Ref, type ComputedRef } from 'vue';
import type { DataValue, DataMap, PropertyValue } from '@/types/ui-protocol';
import { isPathReference } from '@/types/ui-protocol';
import { getData, setData } from '@/protocol/path-resolver';
import type { MessageProcessor } from '@/protocol/message-processor';

/**
 * Data binding context provided by NomadSurface
 */
export interface DataBindingContext {
  processor: Ref<MessageProcessor | null>;
  surfaceId: string;
  dataModel: Ref<DataMap>;
}

/**
 * Injection key for data binding context
 */
export const DATA_BINDING_CONTEXT_KEY = Symbol('nomad-data-binding-context');

/**
 * Options for useNomadDataBinding
 */
export interface UseNomadDataBindingOptions {
  /** Property value (can be literal or path reference) */
  propertyValue?: PropertyValue;
  /** Default value if path doesn't exist */
  defaultValue?: DataValue;
  /** Debounce delay for updates (ms) */
  debounceMs?: number;
}

/**
 * Return type for useNomadDataBinding
 */
export interface UseNomadDataBindingReturn {
  /** Resolved value (reactive) */
  value: ComputedRef<DataValue | undefined>;
  /** Whether the property is a path reference */
  isPath: boolean;
  /** Path string if it's a path reference */
  path: string | null;
  /** Update the value (for two-way binding) */
  updateValue: (newValue: DataValue) => void;
}

/**
 * Composable for reactive data binding in Nomad UI components
 *
 * @param options - Binding options
 * @returns Reactive binding utilities
 *
 * @example
 * ```typescript
 * const { value, updateValue, isPath } = useNomadDataBinding({
 *   propertyValue: props.value,
 *   defaultValue: '',
 * });
 *
 * // Use in template
 * // <input :value="value" @input="updateValue($event.target.value)" />
 * ```
 */
export function useNomadDataBinding(
  options: UseNomadDataBindingOptions = {},
): UseNomadDataBindingReturn {
  const { propertyValue, defaultValue = null, debounceMs = 0 } = options;

  // Try to inject context from NomadSurface
  const context = inject<DataBindingContext | null>(DATA_BINDING_CONTEXT_KEY, null);

  // Determine if this is a path reference
  const isPath = propertyValue ? isPathReference(propertyValue) : false;
  const path = isPath && propertyValue ? (propertyValue as { path: string }).path : null;

  // Debounce timer
  let debounceTimer: ReturnType<typeof setTimeout> | null = null;

  // Computed value that resolves the property
  const value = computed<DataValue | undefined>(() => {
    if (!propertyValue) {
      return defaultValue ?? undefined;
    }

    // Literal values
    if ('literalString' in propertyValue) {
      return propertyValue.literalString;
    }
    if ('literalNumber' in propertyValue) {
      return propertyValue.literalNumber;
    }
    if ('literalBoolean' in propertyValue) {
      return propertyValue.literalBoolean;
    }

    // Path reference - get from data model
    if ('path' in propertyValue && context?.dataModel.value) {
      const result = getData(context.dataModel.value, propertyValue.path);
      return result ?? defaultValue ?? undefined;
    }

    return defaultValue ?? undefined;
  });

  /**
   * Update the value in the data model
   */
  function updateValue(newValue: DataValue): void {
    if (!isPath || !path || !context) {
      return;
    }

    const doUpdate = () => {
      if (context.processor.value) {
        context.processor.value.setData(context.surfaceId, path, newValue);
      }
      else if (context.dataModel.value) {
        setData(context.dataModel.value, path, newValue);
      }
    };

    if (debounceMs > 0) {
      if (debounceTimer) {
        clearTimeout(debounceTimer);
      }
      debounceTimer = setTimeout(doUpdate, debounceMs);
    }
    else {
      doUpdate();
    }
  }

  return {
    value,
    isPath,
    path,
    updateValue,
  };
}

/**
 * Create a data binding context for NomadSurface
 *
 * @param processor - MessageProcessor instance
 * @param surfaceId - Surface ID
 * @param dataModel - Reactive data model reference
 * @returns Data binding context
 */
export function createDataBindingContext(
  processor: Ref<MessageProcessor | null>,
  surfaceId: string,
  dataModel: Ref<DataMap>,
): DataBindingContext {
  return {
    processor,
    surfaceId,
    dataModel,
  };
}

/**
 * Composable for watching data model changes
 *
 * @param path - JSON Pointer path to watch
 * @param callback - Callback when value changes
 * @returns Stop watching function
 */
export function useWatchDataPath(
  path: string,
  callback: (newValue: DataValue | null, oldValue: DataValue | null) => void,
): () => void {
  const context = inject<DataBindingContext | null>(DATA_BINDING_CONTEXT_KEY, null);

  if (!context) {
    console.warn('useWatchDataPath: No data binding context found');
    return () => {};
  }

  const stopWatch = watch(
    () => getData(context.dataModel.value, path),
    (newValue, oldValue) => {
      callback(newValue, oldValue);
    },
    { deep: true },
  );

  return stopWatch;
}
