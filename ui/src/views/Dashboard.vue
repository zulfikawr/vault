<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import { useAuthStore } from '../stores/auth';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import Table from '../components/Table.vue';
import { 
  FolderOpen, 
  Plus, 
  Terminal, 
  Settings, 
  Database, 
  TrendingUp, 
  Cloud,
  ChevronRight,
  Activity,
  Zap,
  Clock,
  ExternalLink,
  BookOpen
} from 'lucide-vue-next';

interface Field {
  name: string;
  type: string;
  required?: boolean;
}

interface Collection {
  name: string;
  type: string;
  fields: Field[];
  created: string;
}

const router = useRouter();
const auth = useAuthStore();
const collections = ref<Collection[]>([]);
const recordCounts = ref<Record<string, number>>({});
const totalRecords = ref(0);
const apiRequests = ref(0);
const storage = ref('-');
const isLoading = ref(true);

const fetchCollections = async () => {
  try {
    const response = await axios.get('/api/admin/collections');
    collections.value = response.data.data;
  } catch (error) {
    console.error('Failed to fetch collections', error);
  }
};

const fetchDashboardStats = async () => {
  isLoading.value = true;
  try {
    let total = 0;
    const countPromises = collections.value
      .filter(col => !col.name.startsWith('_'))
      .map(async (col) => {
        try {
          const response = await axios.get(`/api/collections/${col.name}/records?perPage=1`);
          const count = response.data.totalItems || 0;
          recordCounts.value[col.name] = count;
          return count;
        } catch (error) {
          recordCounts.value[col.name] = 0;
          return 0;
        }
      });
    
    const counts = await Promise.all(countPromises);
    total = counts.reduce((acc, curr) => acc + curr, 0);
    totalRecords.value = total;

    try {
      const storageResponse = await axios.get('/api/admin/storage/stats');
      const stats = storageResponse.data.data;
      storage.value = formatBytes(stats.total_size || 0);
    } catch (error) {
      storage.value = '0 B';
    }

    try {
      const logsResponse = await axios.get('/api/admin/logs?limit=1');
      apiRequests.value = logsResponse.data.totalItems || 0;
    } catch (error) {
      apiRequests.value = 0;
    }
  } catch (error) {
    console.error('Failed to fetch dashboard stats', error);
  } finally {
    isLoading.value = false;
  }
};

const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

onMounted(async () => {
  if (!auth.isAuthenticated) {
    router.push('/login');
    return;
  }
  await fetchCollections();
  await fetchDashboardStats();
});

const stats = [
  { label: 'Collections', value: () => collections.value.length, icon: FolderOpen, color: 'text-primary' },
  { label: 'Total Records', value: () => totalRecords.value, icon: Database, color: 'text-secondary' },
  { label: 'API Requests', value: () => apiRequests.value, icon: TrendingUp, color: 'text-accent' },
  { label: 'Storage', value: () => storage.value, icon: Cloud, color: 'text-info' },
];
</script>

<template>
  <AppLayout>
    <AppHeader>
      <template #breadcrumb>
        <div class="flex items-center gap-2 text-xs font-bold tracking-wide">
          <Activity class="w-4 h-4 text-primary" />
          <span class="text-text-dim">Dashboard</span>
          <ChevronRight class="w-3 h-3 text-text-dim" />
          <span class="text-primary">Overview</span>
        </div>
      </template>
    </AppHeader>

    <div class="flex-1 min-h-0 p-4 lg:p-8 space-y-8 max-w-7xl mx-auto w-full">
      <!-- Welcome Header -->
      <div class="flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
          <h1 class="text-3xl font-bold text-text tracking-tight">
            Dashboard <span class="text-primary text-2xl ml-2 opacity-50 font-normal">Overview</span>
          </h1>
          <p class="text-text-dim text-xs mt-1 flex items-center gap-2">
            <Clock class="w-3 h-3" />
            Last system update: {{ new Date().toLocaleTimeString() }}
          </p>
        </div>
        <div class="flex items-center gap-2">
          <div class="px-3 py-1.5 rounded-lg bg-surface-dark border border-border flex items-center gap-2">
            <div class="w-2 h-2 rounded-full bg-success"></div>
            <span class="text-[10px] font-bold tracking-tight text-text-muted">System Active</span>
          </div>
        </div>
      </div>

      <!-- Stats Grid -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div 
          v-for="stat in stats" 
          :key="stat.label"
          class="group bg-surface/40 backdrop-blur-sm border border-border/50 p-6 rounded-2xl transition-all duration-300"
        >
          <div class="flex justify-between items-start">
            <div>
              <p class="text-[10px] font-bold text-text-dim tracking-widest">{{ stat.label }}</p>
              <h3 class="text-2xl font-bold text-text mt-1 tracking-tight">
                {{ typeof stat.value === 'function' ? stat.value() : stat.value }}
              </h3>
            </div>
            <div :class="['p-2 rounded-xl bg-surface-dark border border-border transition-transform', stat.color]">
              <component :is="stat.icon" class="w-5 h-5" />
            </div>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Left Column: Recent Collections -->
        <div class="lg:col-span-2 space-y-8">
          <section>
            <div class="flex items-center justify-between mb-4 px-2">
              <div class="flex items-center gap-3">
                <Database class="w-4 h-4 text-secondary" />
                <h2 class="text-sm font-bold tracking-widest text-text">Recent Collections</h2>
              </div>
              <Button variant="ghost" size="xs" class="!text-[10px] font-bold tracking-wider hover:text-primary" @click="router.push('/collections')">
                View All
              </Button>
            </div>

            <div class="bg-surface/30 backdrop-blur-sm border border-border/50 rounded-2xl overflow-hidden shadow-sm">
              <Table
                :headers="[
                  { key: 'name', label: 'Name' },
                  { key: 'type', label: 'Type' },
                  { key: 'fields', label: 'Fields', align: 'right' },
                  { key: 'records', label: 'Records', align: 'right' },
                  { key: 'created', label: 'Created', align: 'right' },
                ]"
                :items="collections.filter((c) => !c.name.startsWith('_')).slice(0, 5) as Record<string, unknown>[]"
                row-clickable
                :enable-pagination="false"
                @row-click="(col: Record<string, unknown>) => router.push(`/collections/${col.name}`)"
              >
                <template #cell(name)="{ item }">
                  <div class="flex items-center gap-3 py-1">
                    <div class="w-8 h-8 rounded-lg bg-surface-dark border border-border flex items-center justify-center text-primary font-bold text-[10px]">
                      {{ item.name.toString().substring(0, 2).toUpperCase() }}
                    </div>
                    <span class="font-bold text-xs text-text">{{ item.name }}</span>
                  </div>
                </template>
                <template #cell(type)="{ item }">
                  <span class="px-2 py-0.5 rounded bg-white/5 border border-white/5 text-[9px] font-bold text-text-muted">
                    {{ item.type }}
                  </span>
                </template>
                <template #cell(fields)="{ item }">
                  <span class="text-text-muted text-xs">{{ (item as any).fields?.length || 0 }} fields</span>
                </template>
                <template #cell(records)="{ item }">
                  <span class="text-text-muted text-xs">{{ recordCounts[item.name as string] ?? 0 }} records</span>
                </template>
                <template #cell(created)="{ item }">
                  <span class="text-text-muted text-xs">{{ item.created ? new Date(item.created as string).toLocaleDateString() : '-' }}</span>
                </template>
                <template #empty>
                  <div class="py-12 text-center">
                    <Database class="w-12 h-12 mx-auto mb-4 opacity-10" />
                    <p class="text-xs font-bold tracking-widest text-text-dim">No collections found</p>
                    <Button variant="link" class="mt-4 text-[10px] font-bold tracking-widest" @click="router.push('/collections/new')">
                      Create First Collection
                    </Button>
                  </div>
                </template>
              </Table>
            </div>
          </section>
        </div>

        <!-- Right Column: Quick Actions & Resources -->
        <div class="space-y-8">
          <section class="space-y-4">
            <div class="flex items-center gap-3 px-2">
              <Zap class="w-4 h-4 text-accent" />
              <h2 class="text-sm font-bold tracking-widest text-text">Quick Actions</h2>
            </div>
            
            <div class="space-y-2">
              <button 
                @click="router.push('/collections/new')"
                class="w-full group flex items-center p-4 rounded-2xl bg-surface-dark/50 border border-border hover:border-primary/50 transition-all text-left"
              >
                <div class="p-2.5 rounded-xl bg-primary/10 text-primary mr-4 group-hover:scale-110 transition-transform">
                  <Plus class="w-5 h-5" />
                </div>
                <div class="flex-1">
                  <h3 class="font-bold text-xs tracking-wider text-text">New Collection</h3>
                  <p class="text-[10px] text-text-dim mt-0.5">Define a new data structure</p>
                </div>
                <ChevronRight class="w-4 h-4 text-text-dim group-hover:text-primary transition-colors" />
              </button>

              <button 
                @click="router.push('/logs')"
                class="w-full group flex items-center p-4 rounded-2xl bg-surface-dark/50 border border-border hover:border-secondary/50 transition-all text-left"
              >
                <div class="p-2.5 rounded-xl bg-secondary/10 text-secondary mr-4 group-hover:scale-110 transition-transform">
                  <Terminal class="w-5 h-5" />
                </div>
                <div class="flex-1">
                  <h3 class="font-bold text-xs tracking-wider text-text">View Logs</h3>
                  <p class="text-[10px] text-text-dim mt-0.5">Monitor system activity</p>
                </div>
                <ChevronRight class="w-4 h-4 text-text-dim group-hover:text-secondary transition-colors" />
              </button>

              <button 
                @click="router.push('/settings')"
                class="w-full group flex items-center p-4 rounded-2xl bg-surface-dark/50 border border-border hover:border-accent/50 transition-all text-left"
              >
                <div class="p-2.5 rounded-xl bg-accent/10 text-accent mr-4 group-hover:scale-110 transition-transform">
                  <Settings class="w-5 h-5" />
                </div>
                <div class="flex-1">
                  <h3 class="font-bold text-xs tracking-wider text-text">Settings</h3>
                  <p class="text-[10px] text-text-dim mt-0.5">Manage configuration</p>
                </div>
                <ChevronRight class="w-4 h-4 text-text-dim group-hover:text-accent transition-colors" />
              </button>
            </div>
          </section>

          <section class="p-6 rounded-3xl border border-border bg-surface-dark/30 space-y-4">
            <div class="flex items-center gap-3">
              <BookOpen class="w-4 h-4 text-info" />
              <h2 class="text-sm font-bold tracking-widest text-text">Resources</h2>
            </div>
            <p class="text-[10px] text-text-dim leading-relaxed">
              Access the documentation to learn more about the API, collection rules, and storage configuration.
            </p>
            <Button variant="outline" size="sm" class="w-full !text-[10px] font-bold tracking-widest border-border/50">
              <ExternalLink class="w-3 h-3 mr-2" />
              Documentation
            </Button>
          </section>
        </div>
      </div>
    </div>
  </AppLayout>
</template>
