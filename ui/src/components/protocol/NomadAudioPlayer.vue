<template>
  <div class="nomad-audio-player">
    <div v-if="title" class="nomad-audio-title">{{ title }}</div>
    <audio
      v-if="isValidSrc"
      ref="audioRef"
      class="nomad-audio"
      :src="sanitizedSrc"
      :autoplay="autoplay"
      :loop="loop"
      controls
    >
      Your browser does not support the audio element.
    </audio>
    <div v-else class="nomad-audio-placeholder">
      <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3" />
      </svg>
      <span class="text-sm text-gray-500 mt-1">Invalid audio source</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { validateUrl, sanitizeUrl } from '@/protocol/security';

/**
 * NomadAudioPlayer Component
 *
 * Audio player with URL validation and playback controls.
 */

interface Props {
  componentId?: string;
  surfaceId?: string;
  src?: string;
  title?: string;
  autoplay?: boolean;
  loop?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  src: '',
  title: '',
  autoplay: false,
  loop: false,
});

const audioRef = ref<HTMLAudioElement | null>(null);

// Validate and sanitize URL
const isValidSrc = computed(() => validateUrl(props.src));
const sanitizedSrc = computed(() => sanitizeUrl(props.src));

// Expose audio element for external control
defineExpose({
  audioRef,
});
</script>

<style scoped>
.nomad-audio-player {
  @apply flex flex-col gap-2;
}

.nomad-audio-title {
  @apply text-sm font-medium text-gray-700 dark:text-gray-300;
}

.nomad-audio {
  @apply w-full max-w-md;
}

.nomad-audio-placeholder {
  @apply flex flex-col items-center justify-center p-4 bg-gray-100 dark:bg-gray-700 rounded-lg;
}
</style>
