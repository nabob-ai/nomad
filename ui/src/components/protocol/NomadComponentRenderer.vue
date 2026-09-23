<template>
  <component
    :is="componentType"
    v-bind="componentProps"
    @action="handleAction"
    @update:value="handleValueUpdate"
  >
    <template v-if="node.children && node.children.length > 0">
      <NomadComponentRenderer
        v-for="child in node.children"
        :key="child.id"
        :node="child"
        :surface-id="surfaceId"
        @action="handleAction"
        @update:value="handleValueUpdate"
      />
    </template>
  </component>
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent } from 'vue';
import type { AnyComponentNode, UIActionEvent, DataValue } from '@/types/ui-protocol';

/**
 * NomadComponentRenderer
 *
 * Recursively renders component nodes from the component tree.
 * Maps component types to their Vue implementations.
 */

interface Props {
  /** Component node to render */
  node: AnyComponentNode;
  /** Surface ID for event context */
  surfaceId: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  action: [event: UIActionEvent];
  'update:value': [path: string, value: DataValue];
}>();

// Component type mapping
const componentMap: Record<string, ReturnType<typeof defineAsyncComponent>> = {
  // Layout components
  Row: defineAsyncComponent(() => import('./NomadRow.vue')),
  Column: defineAsyncComponent(() => import('./NomadColumn.vue')),
  Card: defineAsyncComponent(() => import('./NomadCard.vue')),
  List: defineAsyncComponent(() => import('./NomadList.vue')),
  Tabs: defineAsyncComponent(() => import('./NomadTabs.vue')),
  Modal: defineAsyncComponent(() => import('./NomadModal.vue')),
  Divider: defineAsyncComponent(() => import('./NomadDivider.vue')),

  // Content components
  Text: defineAsyncComponent(() => import('./NomadText.vue')),
  Image: defineAsyncComponent(() => import('./NomadImage.vue')),
  Icon: defineAsyncComponent(() => import('./NomadIcon.vue')),
  Video: defineAsyncComponent(() => import('./NomadVideo.vue')),
  AudioPlayer: defineAsyncComponent(() => import('./NomadAudioPlayer.vue')),

  // Input components
  Button: defineAsyncComponent(() => import('./NomadButton.vue')),
  TextField: defineAsyncComponent(() => import('./NomadTextField.vue')),
  Checkbox: defineAsyncComponent(() => import('./NomadCheckbox.vue')),
  Select: defineAsyncComponent(() => import('./NomadSelect.vue')),
  DateTimeInput: defineAsyncComponent(() => import('./NomadDateTimeInput.vue')),
  Slider: defineAsyncComponent(() => import('./NomadSlider.vue')),
  MultipleChoice: defineAsyncComponent(() => import('./NomadMultipleChoice.vue')),
};

// Fallback component for unknown types
const FallbackComponent = defineAsyncComponent(() => import('./NomadFallback.vue'));

// Get the component to render
const componentType = computed(() => {
  return componentMap[props.node.type] ?? FallbackComponent;
});

// Get component props (excluding children which are handled separately)
const componentProps = computed(() => {
  const { children, ...rest } = props.node.props;
  return {
    ...rest,
    componentId: props.node.id,
    surfaceId: props.surfaceId,
  };
});

/**
 * Handle action events and add component context
 */
function handleAction(event: UIActionEvent) {
  emit('action', event);
}

/**
 * Handle value updates from input components
 */
function handleValueUpdate(path: string, value: DataValue) {
  emit('update:value', path, value);
}
</script>
