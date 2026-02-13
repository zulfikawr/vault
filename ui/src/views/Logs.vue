<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import Table from '../components/Table.vue';
import ConfirmModal from '../components/ConfirmModal.vue';
import { Trash2, RefreshCw } from 'lucide-vue-next';

interface LogEntry {
  time: string;
  level: string;
  message: string;
  [key: string]: string;
}

const rawLogs = ref<string[]>([]);
const loading = ref(false);
const showClearModal = ref(false);

const parsedLogs = computed(() => {
  const logs = rawLogs.value.map((log) => {
    const timeMatch = log.match(/time=([^ ]+)/);
    const levelMatch = log.match(/level=([^ ]+)/);
    const msgMatch = log.match(/msg="([^"]*)"|msg=([^ ]+)/);

    const entry: LogEntry = {
      time: timeMatch ? new Date(timeMatch[1]!).toLocaleString() : 'N/A',
      level: levelMatch ? levelMatch[1]! : 'N/A',
      message: msgMatch ? (msgMatch[1]! || msgMatch[2]!) : 'N/A',
    };

    // Extract all key=value pairs
    const kvMatches = log.matchAll(/(\w+)=([^ ]+)/g);
    for (const match of kvMatches) {
      const key = match[1]!;
      if (!['time', 'level', 'msg'].includes(key)) {
        entry[key] = match[2]!;
      }
    }

    return entry;
  });

  // Sort by time descending (latest first)
  return logs.reverse();
});

const headers = computed(() => {
  const baseHeaders = [
    { key: 'time', label: 'Time' },
    { key: 'level', label: 'Level' },
    { key: 'message', label: 'Message' },
  ];

  // Get all unique keys from logs
  const allKeys = new Set<string>();
  parsedLogs.value.forEach((log) => {
    Object.keys(log).forEach((key) => {
      if (!['time', 'level', 'message'].includes(key)) {
        allKeys.add(key);
      }
    });
  });

  // Add dynamic columns
  Array.from(allKeys).forEach((key) => {
    baseHeaders.push({ key, label: key });
  });

  return baseHeaders;
});

const fetchLogs = async () => {
  loading.value = true;
  try {
    const response = await axios.get('/api/admin/logs?limit=500');
    rawLogs.value = response.data.data || [];
  } catch (error) {
    console.error('Failed to fetch logs', error);
  } finally {
    loading.value = false;
  }
};

const handleClearLogs = async () => {
  try {
    await axios.delete('/api/admin/logs');
    rawLogs.value = [];
    showClearModal.value = false;
  } catch (error) {
    console.error('Failed to clear logs', error);
  }
};

onMounted(fetchLogs);
</script>

<template>
  <AppLayout>
    <ConfirmModal
      :show="showClearModal"
      title="Clear Logs"
      message="Are you sure you want to clear all logs? This action cannot be undone."
      confirm-text="Clear"
      cancel-text="Cancel"
      variant="warning"
      @confirm="handleClearLogs"
      @cancel="showClearModal = false"
    />

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
            <Button variant="secondary" size="sm" :disabled="loading" @click="fetchLogs">
              <RefreshCw class="w-4 h-4" />
              Refresh
            </Button>
            <Button variant="outline" size="sm" class="!text-error" @click="showClearModal = true">
              <Trash2 class="w-4 h-4" />
              Clear
            </Button>
          </div>
        </div>

        <Table
          :headers="headers"
          :items="parsedLogs"
          :loading="loading"
          empty-text="No logs available"
        />
      </div>
    </main>
  </AppLayout>
</template>
