<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import Table from '../components/Table.vue';
import ConfirmModal from '../components/ConfirmModal.vue';
import Tabs from '../components/Tabs.vue';
import { Trash2, RefreshCw } from 'lucide-vue-next';

interface LogEntry {
  time: string;
  level: string;
  message: string;
  [key: string]: string;
}

interface AuditLog {
  id: string;
  action: string;
  resource: string;
  admin_id: string;
  details: string;
  timestamp: string;
}

const activeTab = ref('system');
const tabs = [
  { id: 'system', label: 'System Logs' },
  { id: 'audit', label: 'Audit Logs' },
];

// System Logs
const rawLogs = ref<string[]>([]);
const systemLoading = ref(false);
const showClearModal = ref(false);

const parsedSystemLogs = computed(() => {
  const logs = rawLogs.value.map((log) => {
    const timeMatch = log.match(/time=([^ ]+)/);
    const levelMatch = log.match(/level=([^ ]+)/);
    const msgMatch = log.match(/msg="([^"]*)"|msg=([^ ]+)/);

    const entry: LogEntry = {
      time: timeMatch ? new Date(timeMatch[1]!).toLocaleString() : 'N/A',
      level: levelMatch ? levelMatch[1]! : 'N/A',
      message: msgMatch ? msgMatch[1]! || msgMatch[2]! : 'N/A',
    };

    const kvMatches = log.matchAll(/(\w+)=([^ ]+)/g);
    for (const match of kvMatches) {
      const key = match[1]!;
      if (!['time', 'level', 'msg'].includes(key)) {
        entry[key] = match[2]!;
      }
    }

    return entry;
  });

  return logs.reverse();
});

const systemHeaders = computed(() => {
  const baseHeaders = [
    { key: 'time', label: 'Time' },
    { key: 'level', label: 'Level' },
    { key: 'message', label: 'Message' },
  ];

  const allKeys = new Set<string>();
  parsedSystemLogs.value.forEach((log) => {
    Object.keys(log).forEach((key) => {
      if (!['time', 'level', 'message'].includes(key)) {
        allKeys.add(key);
      }
    });
  });

  Array.from(allKeys).forEach((key) => {
    baseHeaders.push({ key, label: key });
  });

  return baseHeaders;
});

const fetchSystemLogs = async () => {
  systemLoading.value = true;
  try {
    const response = await axios.get('/api/admin/logs?limit=500');
    rawLogs.value = response.data.data || [];
  } catch (error) {
    console.error('Failed to fetch logs', error);
  } finally {
    systemLoading.value = false;
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

// Audit Logs
const auditLogsData = ref<AuditLog[]>([]);
const auditLoading = ref(false);

const parsedAuditLogs = computed(() => {
  return auditLogsData.value.map((log) => ({
    ...log,
    timestamp: new Date(log.timestamp).toLocaleString(),
    details: typeof log.details === 'string' ? log.details : JSON.stringify(log.details),
  }));
});

const auditHeaders = [
  { key: 'timestamp', label: 'Timestamp' },
  { key: 'action', label: 'Action' },
  { key: 'resource', label: 'Resource' },
  { key: 'admin_id', label: 'Admin' },
  { key: 'details', label: 'Details' },
];

const fetchAuditLogs = async () => {
  auditLoading.value = true;
  try {
    const response = await axios.get(
      '/api/collections/_audit_logs/records?limit=100&sort=-timestamp'
    );
    auditLogsData.value = response.data.data || [];
  } catch (error) {
    console.error('Failed to fetch audit logs', error);
  } finally {
    auditLoading.value = false;
  }
};

const refresh = () => {
  if (activeTab.value === 'system') {
    fetchSystemLogs();
  } else {
    fetchAuditLogs();
  }
};

watch(activeTab, (newTab) => {
  if (newTab === 'system' && rawLogs.value.length === 0) {
    fetchSystemLogs();
  } else if (newTab === 'audit' && auditLogsData.value.length === 0) {
    fetchAuditLogs();
  }
});

onMounted(() => {
  fetchSystemLogs();
});
</script>

<template>
  <AppLayout>
    <ConfirmModal
      :show="showClearModal"
      title="Clear Logs"
      message="Are you sure you want to log out? You will need to sign in again to access the admin panel."
      confirm-text="Clear"
      cancel-text="Cancel"
      variant="warning"
      @confirm="handleClearLogs"
      @cancel="showClearModal = false"
    />

    <AppHeader>
      <template #breadcrumb>
        <div class="flex items-center text-sm text-text-muted">
          <span class="font-medium text-primary">Logs</span>
        </div>
      </template>
    </AppHeader>

    <main class="flex-1 overflow-auto min-h-0 p-4 sm:p-8 pb-24 sm:pb-8">
      <div class="max-w-7xl mx-auto space-y-6 sm:space-y-8">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h1 class="text-2xl font-bold text-text tracking-tight">Logs</h1>
            <p class="mt-1 text-sm text-text-muted">Monitor system events and audit trails.</p>
          </div>
          <div class="flex gap-2">
            <Button
              variant="secondary"
              size="sm"
              :disabled="systemLoading || auditLoading"
              @click="refresh"
            >
              <RefreshCw class="w-4 h-4" />
              Refresh
            </Button>
            <Button
              v-if="activeTab === 'system'"
              variant="destructive"
              size="sm"
              @click="showClearModal = true"
            >
              <Trash2 class="w-4 h-4" />
              Clear
            </Button>
          </div>
        </div>

        <Tabs v-model="activeTab" :tabs="tabs" />

        <div v-if="activeTab === 'system'">
          <Table
            :headers="systemHeaders"
            :items="parsedSystemLogs"
            :loading="systemLoading"
            empty-text="No system logs available"
            :enable-pagination="true"
            :default-page-size="20"
          />
        </div>
        <div v-else>
          <Table
            :headers="auditHeaders"
            :items="parsedAuditLogs"
            :loading="auditLoading"
            empty-text="No audit logs available"
            :enable-pagination="true"
            :default-page-size="20"
          />
        </div>
      </div>
    </main>
  </AppLayout>
</template>
