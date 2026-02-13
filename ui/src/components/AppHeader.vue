<script setup lang="ts">
import { ref, onMounted } from 'vue';
import Input from './Input.vue';
import { Search } from 'lucide-vue-next';

const systemStatus = ref<'checking' | 'online' | 'offline'>('checking');
const backendPort = ref('');

const checkSystemStatus = async () => {
  try {
    const response = await fetch('/api/health', { method: 'HEAD' });
    if (response.ok) {
      systemStatus.value = 'online';
      backendPort.value = response.headers.get('X-Server-Port') || import.meta.env.VITE_API_PORT || '8090';
    } else {
      systemStatus.value = 'offline';
    }
  } catch {
    systemStatus.value = 'offline';
  }
};

onMounted(checkSystemStatus);
</script>

<template>
  <header class="hidden sm:flex h-16 items-center justify-between px-4 sm:px-8 border-b border-border bg-surface z-10 shrink-0">
    <slot name="breadcrumb" />
    
    <div class="flex items-center gap-4">
      <span v-if="systemStatus === 'checking'" class="inline-flex items-center px-2.5 py-0.5 rounded text-xs font-medium bg-warning/10 text-warning border border-warning/20">
        <span class="relative flex h-2 w-2 mr-1.5">
          <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-warning opacity-75"></span>
          <span class="relative inline-flex rounded-full h-2 w-2 bg-warning"></span>
        </span>
        Checking...
      </span>
      <span v-else-if="systemStatus === 'online'" class="inline-flex items-center px-2.5 py-0.5 rounded text-xs font-medium bg-success/10 text-success border border-success/20">
        <span class="relative flex h-2 w-2 mr-1.5">
          <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-success opacity-75"></span>
          <span class="relative inline-flex rounded-full h-2 w-2 bg-success"></span>
        </span>
        System Operational (Port: {{ backendPort }})
      </span>
      <span v-else class="inline-flex items-center px-2.5 py-0.5 rounded text-xs font-medium bg-error/10 text-error border border-error/20">
        <span class="relative flex h-2 w-2 mr-1.5 bg-error rounded-full"></span>
        System Offline
      </span>
      
      <div class="relative hidden sm:block group">
        <Input 
          type="text" 
          placeholder="Search (Ctrl+K)"
          class="w-64 !bg-surface-dark"
        />
        <Search class="absolute left-2.5 top-1/2 -translate-y-1/2 text-text-muted w-4 h-4 pointer-events-none" />
      </div>
    </div>
  </header>
</template>
