<template>
  <div id="app" class="min-h-screen" :class="{ 'public-shell': isPublicRoute }">
    <router-view />
    <DialogHost />
  </div>
</template>

<script setup lang="ts">
import { computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import DialogHost from './components/DialogHost.vue'

const route = useRoute()
const isPublicRoute = computed(() => !route.path.startsWith('/ops') && !route.path.startsWith('/auth'))

watch(isPublicRoute, (isPublic) => {
  if (isPublic) document.documentElement.dataset.surface = 'public'
  else delete document.documentElement.dataset.surface
}, { immediate: true })
</script>
