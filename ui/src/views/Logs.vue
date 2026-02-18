<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import Table from '../components/Table.vue';
import ConfirmModal from '../components/ConfirmModal.vue';
import { Trash2, RefreshCw, Terminal, History, ChevronRight, Activity, Search, Filter } from 'lucide-vue-next';

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

const activeTab = ref<'system' | 'audit'>('system');
const tabs = [
  { id: 'system', label: 'System Logs', icon: Terminal },
  { id: 'audit', label: 'Audit Logs', icon: History },
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
      message="Are you sure you want to clear the logs? This action cannot be undone."
      confirm-text="Clear"
      cancel-text="Cancel"
      variant="danger"
      @confirm="handleClearLogs"
      @cancel="showClearModal = false"
    />

    <AppHeader>
      <template #breadcrumb>
        <div class="flex items-center gap-2 text-xs font-bold tracking-wide">
          <Activity class="w-4 h-4 text-primary" />
          <span class="text-primary">Logs</span>
        </div>
      </template>
    </AppHeader>

    <main class="flex-1 min-h-0 p-4 lg:p-8 max-w-7xl mx-auto w-full">
      <div class="space-y-8">
        <!-- Header -->
        <div class="flex flex-col md:flex-row md:items-end justify-between gap-6 border-b border-border/30 pb-6">
          <div>
            <h1 class="text-3xl font-bold text-text tracking-tight">Logs</h1>
            <p class="text-text-dim text-xs mt-1">Monitor system events and audit trails.</p>
          </div>
          
          <div class="flex bg-surface-dark/50 p-1 rounded-xl border border-border/50 backdrop-blur-sm">
            <button 
              v-for="tab in tabs"
              :key="tab.id"
              @click="activeTab = tab.id as any"
              :class="activeTab === tab.id ? 'bg-primary text-white shadow-lg' : 'text-text-dim hover:text-text'"
              class="px-4 py-2 rounded-lg text-xs font-bold transition-all flex items-center gap-2"
            >
              <component :is="tab.icon" class="w-3.5 h-3.5" />
              {{ tab.label }}
            </button>
          </div>
        </div>

        <!-- Toolbar -->
        <div class="flex flex-col lg:flex-row gap-4">
          <div class="relative flex-1">
            <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-text-dim" />
            <input 
              type="text" 
              placeholder="Search logs..." 
              class="w-full bg-surface-dark/50 border border-border/50 rounded-xl pl-10 pr-4 h-10 text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary/50 transition-all"
            />
          </div>
          <div class="flex items-center gap-3">
            <Button variant="secondary" size="md">
              <Filter class="w-4 h-4 mr-2" />
              Filter
            </Button>
            <div class="w-px h-6 bg-border/50 hidden sm:block"></div>
            <Button
              variant="secondary"
              size="md"
              :loading="systemLoading || auditLoading"
              @click="refresh"
            >
              <RefreshCw class="w-4 h-4 mr-2" />
              Refresh
            </Button>
            <Button
              v-if="activeTab === 'system'"
              variant="destructive"
              size="md"
              @click="showClearModal = true"
            >
              <Trash2 class="w-4 h-4 mr-2" />
              Clear
            </Button>
          </div>
        </div>

        <!-- Tables -->
        <div v-if="activeTab === 'system'" class="animate-in fade-in slide-in-from-bottom-4 duration-500">
          <Table
            :headers="systemHeaders"
            :items="parsedSystemLogs"
            :loading="systemLoading"
            empty-text="No system logs found"
            :enable-pagination="true"
            :default-page-size="20"
          >
            <template #cell(level)="{ item }">
              <span 
                :class="[
                  'px-2 py-0.5 rounded text-[10px] font-bold',
                  item.level === 'DEBUG' ? 'bg-info/10 text-info border border-info/20' :
                  item.level === 'INFO' ? 'bg-success/10 text-success border border-success/20' :
                  item.level === 'WARN' ? 'bg-warning/10 text-warning border border-warning/20' :
                  'bg-error/10 text-error border border-error/20'
                ]"
              >
                {{ item.level }}
              </span>
            </template>
            <template #cell(message)="{ item }">
              <span class="text-xs font-mono text-text leading-relaxed">{{ item.message }}</span>
            </template>
            <template #cell(time)="{ item }">
              <span class="text-[10px] font-bold text-text-dim">{{ item.time }}</span>
            </template>
          </Table>
        </div>

        <div v-else class="animate-in fade-in slide-in-from-bottom-4 duration-500">
          <Table
            :headers="auditHeaders"
            :items="parsedAuditLogs"
            :loading="auditLoading"
            empty-text="No audit logs found"
            :enable-pagination="true"
            :default-page-size="20"
          >
            <template #cell(action)="{ item }">
              <span class="px-2 py-0.5 rounded bg-white/5 border border-white/5 text-[10px] font-bold">
                {{ item.action }}
              </span>
            </template>
            <template #cell(details)="{ item }">
              <span class="text-xs font-mono text-text-muted truncate block max-w-xs">{{ item.details }}</span>
            </template>
            <template #cell(timestamp)="{ item }">
              <span class="text-[10px] font-bold text-text-dim">{{ item.timestamp }}</span>
            </template>
          </Table>
        </div>
      </div>
    </main>
  </AppLayout>
</template>
