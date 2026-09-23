<template>
  <div class="nomad-text-field">
    <label v-if="labelValue" :for="inputId" class="nomad-text-field-label">
      {{ labelValue }}
    </label>
    <textarea
      v-if="multiline"
      :id="inputId"
      :value="currentValue"
      class="nomad-text-field-input nomad-text-field-textarea"
      :placeholder="placeholderValue"
      :disabled="disabledValue"
      :maxlength="maxLength"
      @input="handleInput"
      @keydown.enter.ctrl="handleSubmit"
    />
    <input
      v-else
      :id="inputId"
      :value="currentValue"
      class="nomad-text-field-input"
      :type="inputType"
      :placeholder="placeholderValue"
      :disabled="disabledValue"
      :maxlength="maxLength"
      @input="handleInput"
      @keydown.enter="handleSubmit"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import type { DataValue, PropertyValue, UIActionEvent } from '@/types/ui-protocol';
import { createUIActionEvent } from '@/types/ui-protocol';
import { useNomadDataBinding } from '@/composables/useNomadDataBinding';
import { useUIAction } from '@/composables/useUIAction';

/**
 * NomadTextField Component
 *
 * Text input with label, placeholder, and two-way binding support.
 * Uses useNomadDataBinding for reactive data binding with the data model.
 */

interface Props {
  componentId?: string;
  surfaceId?: string;
  /** Value property - can be literal or path reference */
  value?: PropertyValue | string;
  /** Label property - can be literal or path reference */
  label?: PropertyValue | string;
  /** Placeholder property - can be literal or path reference */
  placeholder?: PropertyValue | string;
  multiline?: boolean;
  /** Disabled property - can be literal or path reference */
  disabled?: PropertyValue | boolean;
  maxLength?: number;
  inputType?: 'text' | 'password' | 'email' | 'number' | 'tel' | 'url';
  /** Path for two-way binding (legacy support) */
  valuePath?: string;
}

const props = withDefaults(defineProps<Props>(), {
  multiline: false,
  inputType: 'text',
});

const emit = defineEmits<{
  'update:value': [path: string, value: DataValue];
  action: [event: UIActionEvent];
}>();

// Generate unique ID for label association
const inputId = computed(() => `nomad-text-field-${props.componentId ?? Math.random().toString(36).slice(2)}`);

// Setup UI action emitter
const { emitFullAction } = useUIAction({
  surfaceId: props.surfaceId,
  componentId: props.componentId,
});

// Convert props to PropertyValue format for binding
function toPropertyValue(value: PropertyValue | string | boolean | undefined): PropertyValue | undefined {
  if (value === undefined) return undefined;
  if (typeof value === 'string') return { literalString: value };
  if (typeof value === 'boolean') return { literalBoolean: value };
  return value as PropertyValue;
}

// Use data binding for value
const { value: boundValue, updateValue, path: valuePath } = useNomadDataBinding({
  propertyValue: toPropertyValue(props.value),
  defaultValue: '',
});

// Use data binding for label
const { value: boundLabel } = useNomadDataBinding({
  propertyValue: toPropertyValue(props.label),
  defaultValue: '',
});

// Use data binding for placeholder
const { value: boundPlaceholder } = useNomadDataBinding({
  propertyValue: toPropertyValue(props.placeholder),
  defaultValue: '',
});

// Use data binding for disabled
const { value: boundDisabled } = useNomadDataBinding({
  propertyValue: toPropertyValue(props.disabled),
  defaultValue: false,
});

// Computed values for template
const currentValue = computed(() => String(boundValue.value ?? ''));
const labelValue = computed(() => boundLabel.value ? String(boundLabel.value) : undefined);
const placeholderValue = computed(() => String(boundPlaceholder.value ?? ''));
const disabledValue = computed(() => Boolean(boundDisabled.value));

function handleInput(event: Event) {
  const target = event.target as HTMLInputElement | HTMLTextAreaElement;
  const newValue = target.value;

  // Update via data binding
  updateValue(newValue);

  // Also emit for legacy support
  const path = valuePath ?? props.valuePath;
  if (path) {
    emit('update:value', path, newValue);
  }
}

function handleSubmit(event: KeyboardEvent) {
  // Emit submit action when Enter is pressed (Ctrl+Enter for multiline)
  if (props.componentId && props.surfaceId) {
    const actionEvent = createUIActionEvent(
      props.surfaceId,
      props.componentId,
      'submit',
      { value: currentValue.value },
    );
    emit('action', actionEvent);
    emitFullAction(actionEvent);
  }
}
</script>

<style scoped>
.nomad-text-field {
  @apply flex flex-col;
  gap: var(--nomad-spacing-xs, 0.25rem);
}

.nomad-text-field-label {
  @apply font-medium;
  font-size: var(--nomad-font-size-sm, 0.875rem);
  color: var(--nomad-text-secondary, #6b7280);
}

.dark .nomad-text-field-label {
  color: var(--nomad-text-secondary, #cbd5e1);
}

.nomad-text-field-input {
  @apply w-full focus:outline-none focus:ring-2 focus:border-transparent;
  padding: var(--nomad-spacing-sm, 0.5rem) var(--nomad-spacing-md, 1rem);
  font-size: var(--nomad-font-size-sm, 0.875rem);
  border: 1px solid var(--nomad-border, #e5e7eb);
  border-radius: var(--nomad-radius-lg, 0.5rem);
  background-color: var(--nomad-surface, #ffffff);
  color: var(--nomad-text, #111827);
  --tw-ring-color: var(--nomad-border-focus, #3b82f6);
}

.dark .nomad-text-field-input {
  border-color: var(--nomad-border, #475569);
  background-color: var(--nomad-surface, #334155);
  color: var(--nomad-text, #f1f5f9);
}

.nomad-text-field-input::placeholder {
  color: var(--nomad-text-muted, #9ca3af);
}

.dark .nomad-text-field-input::placeholder {
  color: var(--nomad-text-muted, #94a3b8);
}

.nomad-text-field-input:disabled {
  @apply cursor-not-allowed opacity-50;
  background-color: var(--nomad-background-alt, #f3f4f6);
}

.dark .nomad-text-field-input:disabled {
  background-color: var(--nomad-background-alt, #0f172a);
}

.nomad-text-field-textarea {
  @apply min-h-[100px] resize-y;
}
</style>
