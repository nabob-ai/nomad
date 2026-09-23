<template>
  <div class="nomad-slider">
    <div v-if="labelValue || showValue" class="nomad-slider-header">
      <label v-if="labelValue" :for="sliderId" class="nomad-slider-label">
        {{ labelValue }}
      </label>
      <span v-if="showValue" class="nomad-slider-value">{{ currentValue }}</span>
    </div>
    <input
      :id="sliderId"
      :value="currentValue"
      type="range"
      class="nomad-slider-input"
      :min="min"
      :max="max"
      :step="step"
      :disabled="disabledValue"
      @input="handleInput"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { DataValue, PropertyValue } from '@/types/ui-protocol';
import { useNomadDataBinding } from '@/composables/useNomadDataBinding';

/**
 * NomadSlider Component
 *
 * Range slider with min/max/step and two-way binding support.
 * Uses useNomadDataBinding for reactive data binding with the data model.
 */

interface Props {
  componentId?: string;
  surfaceId?: string;
  /** Value property - can be literal or path reference */
  value?: PropertyValue | number;
  /** Label property - can be literal or path reference */
  label?: PropertyValue | string;
  min?: number;
  max?: number;
  step?: number;
  /** Disabled property - can be literal or path reference */
  disabled?: PropertyValue | boolean;
  showValue?: boolean;
  /** Path for two-way binding (legacy support) */
  valuePath?: string;
}

const props = withDefaults(defineProps<Props>(), {
  min: 0,
  max: 100,
  step: 1,
  showValue: false,
});

const emit = defineEmits<{
  'update:value': [path: string, value: DataValue];
}>();

// Generate unique ID for label association
const sliderId = computed(() => `nomad-slider-${props.componentId ?? Math.random().toString(36).slice(2)}`);

// Convert props to PropertyValue format for binding
function toPropertyValue(value: PropertyValue | string | boolean | number | undefined): PropertyValue | undefined {
  if (value === undefined) return undefined;
  if (typeof value === 'string') return { literalString: value };
  if (typeof value === 'boolean') return { literalBoolean: value };
  if (typeof value === 'number') return { literalNumber: value };
  return value as PropertyValue;
}

// Use data binding for value
const { value: boundValue, updateValue, path: valuePath } = useNomadDataBinding({
  propertyValue: toPropertyValue(props.value),
  defaultValue: props.min,
});

// Use data binding for label
const { value: boundLabel } = useNomadDataBinding({
  propertyValue: toPropertyValue(props.label),
  defaultValue: '',
});

// Use data binding for disabled
const { value: boundDisabled } = useNomadDataBinding({
  propertyValue: toPropertyValue(props.disabled),
  defaultValue: false,
});

// Computed values for template
const currentValue = computed(() => Number(boundValue.value ?? props.min));
const labelValue = computed(() => boundLabel.value ? String(boundLabel.value) : undefined);
const disabledValue = computed(() => Boolean(boundDisabled.value));

function handleInput(event: Event) {
  const target = event.target as HTMLInputElement;
  const newValue = Number(target.value);

  // Update via data binding
  updateValue(newValue);

  // Also emit for legacy support
  const path = valuePath ?? props.valuePath;
  if (path) {
    emit('update:value', path, newValue);
  }
}
</script>

<style scoped>
.nomad-slider {
  @apply flex flex-col gap-2;
}

.nomad-slider-header {
  @apply flex items-center justify-between;
}

.nomad-slider-label {
  @apply text-sm font-medium text-gray-700 dark:text-gray-300;
}

.nomad-slider-value {
  @apply text-sm text-gray-500 dark:text-gray-400 tabular-nums;
}

.nomad-slider-input {
  @apply w-full h-2 bg-gray-200 dark:bg-gray-700 rounded-lg appearance-none cursor-pointer;
  @apply disabled:opacity-50 disabled:cursor-not-allowed;
}

.nomad-slider-input::-webkit-slider-thumb {
  @apply appearance-none w-4 h-4 bg-blue-600 rounded-full cursor-pointer;
}

.nomad-slider-input::-moz-range-thumb {
  @apply w-4 h-4 bg-blue-600 rounded-full cursor-pointer border-0;
}

.nomad-slider-input:focus {
  @apply outline-none;
}

.nomad-slider-input:focus::-webkit-slider-thumb {
  @apply ring-2 ring-blue-500 ring-offset-2;
}
</style>
