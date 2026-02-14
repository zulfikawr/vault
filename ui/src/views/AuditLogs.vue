<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import Table from '../components/Table.vue';
import { RefreshCw } from 'lucide-vue-next';

interface AuditLog {
  id: string;
  action: string;
  resource: string;
  admin_id: string;
  details: string;
  timestamp: string;
}

const auditLogs = ref<AuditLog[]>([]);
const loading = ref(false);

const parsedLogs = computed(() => {
  return auditLogs.value.map((log) => ({
    ...log,
    timestamp: new Date(log.timestamp).toLocaleString(),
    details: typeof log.details === 'string' ? log.details : JSON.stringify(log.details),
  }));
});

const headers = [
  { key: 'timestamp', label: 'Timestamp' },
  { key: 'action', label: 'Action' },
  { key: 'resource', label: 'Resource' },
  { key: 'admin_id', label: 'Admin' },
  { key: 'details', label: 'Details' },
];

const fetchAuditLogs = async () => {
  loading.value = true;
  try {
    const response = await axios.get('/api/collections/_audit_logs/records?limit=100&sort=-timestamp');
    auditLogs.value = response.data.data || [];
  } catch (error) {
    console.error('Failed to fetch audit logs', error);
  } finally {
    loading.value = false;
  }
};

onMounted(fetchAuditLogs);
</script>

<template>
  <AppLayout>
    <AppHeader>
      <template #breadcrumb>
        <div class="flex items-center text-sm text-text-muted">
          <span class="hover:text-text cursor-pointer" @click="$router.push('/')">Vault</span>
          <span class="mx-2">/</span>
          <span class="font-medium text-text">Audit Logs</span>
        </div>
      </template>
    </AppHeader>

    <main class="flex-1 overflow-auto min-h-0 p-4 sm:p-8 pb-24 sm:pb-8">
      <div class="max-w-7xl mx-auto space-y-6 sm:space-y-8">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h1 class="text-2xl font-bold text-text tracking-tight">Audit Logs</h1>
            <p class="mt-1 text-sm text-text-muted">Track all collection management operations.</p>
          </div>
          <Button variant="secondary" size="sm" :disabled="loading" @click="fetchAuditLogs">
            <RefreshCw class="w-4 h-4" />
            Refresh
          </Button>
        </div>

        <Table
          :headers="headers"
          :items="parsedLogs"
          :loading="loading"
          empty-text="No audit logs available"
          :enable-pagination="true"
          :default-page-size="20"
        />
      </div>
    </main>
  </AppLayout>
</template>
