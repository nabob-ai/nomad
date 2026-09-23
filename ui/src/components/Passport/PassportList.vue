<template>
  <div class="passport-list">
    <!-- Header -->
    <div class="list-header">
      <div class="header-left">
        <h2 class="list-title">全球护照排名</h2>
        <span class="passport-count">{{ passports.length }} 个</span>
      </div>
    </div>

    <!-- Search -->
    <div class="list-search">
      <svg class="search-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
      </svg>
      <input v-model="searchQuery" type="text" placeholder="搜索 Passport..." class="search-input" />
    </div>

    <!-- Loading -->
    <LoadingSpinner v-if="loading" class="my-8" />

    <!-- Passport List -->
    <div v-else class="passport-grid">
      <div v-for="passport in filteredPassports" :key="passport.position" class="passport-card" @click="$emit('select', passport)">
        <div class="passport-header">
          <div class="passport-icon">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"
              ></path>
            </svg>
          </div>

          <div class="passport-info">
            <h3 class="passport-name">国家：{{ passport.passport }}</h3>
          </div>
        </div>

        <div class="passport-info">
          <h3 class="passport-name">排名：{{ passport.rank }}</h3>
        </div>

        <div class="passport-info">
          <h3 class="passport-name">国家代码：{{ passport.code }}</h3>
        </div>

        <div class="passport-info">
          <h3 class="passport-name">区域：{{ passport.region }}</h3>
        </div>

        <div class="passport-info">
          <h3 class="passport-name">得分：{{ passport.score }}</h3>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, computed, type PropType } from "vue";
import LoadingSpinner from "../Common/LoadingSpinner.vue";
import EmptyState from "../Common/EmptyState.vue";
import type { Passport } from "@/types";

export default defineComponent({
  name: "PassportList",
  components: {
    LoadingSpinner,
    EmptyState,
  },
  props: {
    passports: {
      type: Array as PropType<Passport[]>,
      required: true,
    },
    loading: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const searchQuery = ref("");
    const filteredPassports = computed(() => {
      if (!searchQuery.value) return props.passports;

      const query = searchQuery.value.toLowerCase();
      return props.passports.filter((passport) => passport.passport.toLowerCase().includes(query));
    });

    return {
      searchQuery,
      filteredPassports,
    };
  },
});
</script>

<style scoped>
.passport-list {
  @apply space-y-4;
}

.list-header {
  @apply flex items-center justify-between;
}

.header-left {
  @apply flex items-center gap-3;
}

.list-title {
  @apply text-2xl font-bold text-text dark:text-text-dark;
}

.passport-count {
  @apply px-2 py-1 bg-primary/10 dark:bg-primary/20 text-primary dark:text-primary-light text-sm font-medium rounded-full;
}

.btn-create {
  @apply flex items-center gap-2 px-4 py-2 bg-primary hover:bg-primary-hover text-white rounded-lg font-medium transition-colors;
}

.list-search {
  @apply relative;
}

.search-icon {
  @apply absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-secondary dark:text-secondary-dark;
}

.search-input {
  @apply w-full pl-10 pr-4 py-2 bg-surface dark:bg-surface-dark border border-border dark:border-border-dark rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20;
}

.passport-grid {
  @apply grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4;
}

.passport-card {
  @apply bg-surface dark:bg-surface-dark border border-border dark:border-border-dark rounded-lg p-4 hover:shadow-md transition-all cursor-pointer;
}

.passport-header {
  @apply flex items-start gap-3 mb-3;
}

.passport-icon {
  @apply w-12 h-12 rounded-full bg-primary/10 dark:bg-primary/20 flex items-center justify-center text-primary dark:text-primary-light flex-shrink-0;
}

.passport-info {
  @apply flex-1 min-w-0;
}

.passport-name {
  @apply text-base font-semibold text-text dark:text-text-dark truncate;
}

.passport-members {
  @apply text-sm text-secondary dark:text-secondary-dark mt-1;
}

.passport-description {
  @apply text-sm text-secondary dark:text-secondary-dark mb-3 line-clamp-2;
}

.passport-footer {
  @apply flex items-center justify-between pt-3 border-t border-border dark:border-border-dark;
}

.passport-time {
  @apply text-xs text-secondary dark:text-secondary-dark;
}

.btn-delete {
  @apply p-1.5 hover:bg-red-50 dark:hover:bg-red-900/20 text-red-600 dark:text-red-400 rounded transition-colors;
}
</style>
