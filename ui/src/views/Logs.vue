<script setup lang="ts">
import { ref, onMounted } from 'vue';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import { Trash2, RefreshCw } from 'lucide-vue-next';

const logs = ref<string[]>([]);
const loading = ref(false);

const fetchLogs = async () => {
  loading.value = true;
  try {
    const response = await axios.get('/api/admin/logs?limit=500');
    logs.value = response.data.data || [];
  } catch (error) {
    console.error('Failed to fetch logs', error);
  } finally {
    loading.value = false;
  }
};

const clearLogs = async () => {
  if (!confirm('Are you sure you want to clear all logs?')) return;

  try {
    await axios.delete('/api/admin/logs');
    logs.value = [];
  } catch (error) {
    console.error('Failed to clear logs', error);
  }
};

onMounted(fetchLogs);
</script>

<template>
  <AppLayout>
    <AppHeader>
      <template #breadcrumb>
        <div class="flex items-center gap-2">
          <span class="text-sm text-text-muted">Logs</span>
        </div>
      </template>
    </AppHeader>

    <main class="flex-1 overflow-auto">
      <div class="p-6 max-w-7xl mx-auto">
        <div class="flex items-center justify-between mb-6">
          <h1 class="text-2xl font-bold text-text">System Logs</h1>
          <div class="flex gap-2">
            <Button
              variant="secondary"
              size="sm"
              @click="fetchLogs"
              :disabled="loading"
            >
              <RefreshCw class="w-4 h-4" />
              Refresh
            </Button>
            <Button
              variant="outline"
              size="sm"
              @click="clearLogs"
              class="!text-error"
            >
              <Trash2 class="w-4 h-4" />
              Clear
            </Button>
          </div>
        </div>

        <div class="bg-surface border border-border rounded-lg overflow-hidden">
          <div
            v-if="logs.length === 0"
            class="p-8 text-center text-text-muted"
          >
            <p>No logs available</p>
          </div>
          <div v-else class="max-h-[600px] overflow-y-auto">
            <div class="font-mono text-xs bg-surface-dark p-4 space-y-1">
              <div
                v-for="(log, index) in logs"
                :key="index"
                class="text-text-muted hover:text-text transition-colors"
              >
                {{ log }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </AppLayout>
</template>
