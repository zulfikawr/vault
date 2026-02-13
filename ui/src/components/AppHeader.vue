<script setup lang="ts">
import { ref, onMounted } from 'vue';
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
        <input 
          type="text" 
          placeholder="Search (Ctrl+K)" 
          class="w-64 bg-surface-dark border-border rounded-md py-1.5 pl-9 pr-3 text-sm focus:ring-1 focus:ring-primary focus:border-primary focus:outline-none text-text placeholder-text-muted transition-all"
        />
        <Search class="absolute left-2.5 top-2 text-text-muted w-4 h-4" />
      </div>
    </div>
  </header>
</template>
