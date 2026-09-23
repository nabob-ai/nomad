<template>
  <video
    v-if="isValidSrc"
    ref="videoRef"
    class="nomad-video"
    :src="sanitizedSrc"
    :poster="sanitizedPoster"
    :autoplay="autoplay"
    :controls="controls"
    :loop="loop"
    :muted="muted"
    playsinline
  >
    Your browser does not support the video tag.
  </video>
  <div v-else class="nomad-video-placeholder">
    <svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
    </svg>
    <span class="text-sm text-gray-500 mt-2">Invalid video source</span>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { validateUrl, sanitizeUrl } from '@/protocol/security';

/**
 * NomadVideo Component
 *
 * Video player with URL validation and playback controls.
 */

interface Props {
  componentId?: string;
  surfaceId?: string;
  src?: string;
  poster?: string;
  autoplay?: boolean;
  controls?: boolean;
  loop?: boolean;
  muted?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  src: '',
  poster: '',
  autoplay: false,
  controls: true,
  loop: false,
  muted: false,
});

const videoRef = ref<HTMLVideoElement | null>(null);

// Validate and sanitize URLs
const isValidSrc = computed(() => validateUrl(props.src));
const sanitizedSrc = computed(() => sanitizeUrl(props.src));
const sanitizedPoster = computed(() => props.poster ? sanitizeUrl(props.poster) : undefined);

// Expose video element for external control
defineExpose({
  videoRef,
});
</script>

<style scoped>
.nomad-video {
  @apply w-full max-w-2xl rounded-lg;
}

.nomad-video-placeholder {
  @apply flex flex-col items-center justify-center w-full h-48 bg-gray-100 dark:bg-gray-700 rounded-lg;
}
</style>
