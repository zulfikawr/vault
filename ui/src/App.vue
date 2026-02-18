<script setup lang="ts">
import { onMounted } from 'vue';
import { useUIStore } from './stores/ui';

const uiStore = useUIStore();

onMounted(() => {
  uiStore.initTheme();
});
</script>

<template>
  <router-view />
</template>
