<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import { useAuthStore } from '../stores/auth';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import Table from '../components/Table.vue';
import { FolderOpen, Plus, Terminal, Settings, Database, TrendingUp, Cloud } from 'lucide-vue-next';

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
const totalRecords = ref(0);
const apiRequests = ref(0);
const storage = ref('-');

const fetchCollections = async () => {
  try {
    const response = await axios.get('/api/admin/collections');
    collections.value = response.data.data;
  } catch (error) {
    console.error('Failed to fetch collections', error);
  }
};

const fetchDashboardStats = async () => {
  try {
    // Fetch total records count from all collections
    let total = 0;
    for (const col of collections.value) {
      if (!col.name.startsWith('_')) {
        try {
          const response = await axios.get(`/api/collections/${col.name}/records?perPage=1`);
          total += response.data.totalItems || 0;
        } catch (error) {
          console.error(`Failed to fetch records for ${col.name}:`, error);
        }
      }
    }
    totalRecords.value = total;

    // Fetch storage stats
    try {
      const storageResponse = await axios.get('/api/admin/storage/stats');
      const stats = storageResponse.data.data;
      storage.value = formatBytes(stats.total_size || 0);
    } catch (error) {
      console.error('Failed to fetch storage stats:', error);
      storage.value = '-';
    }

    // Fetch logs to count API requests (approximate)
    try {
      const logsResponse = await axios.get('/api/admin/logs?limit=1');
      apiRequests.value = logsResponse.data.totalItems || 0;
    } catch (error) {
      console.error('Failed to fetch API request count:', error);
      apiRequests.value = 0;
    }
  } catch (error) {
    console.error('Failed to fetch dashboard stats', error);
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
  // Ensure user is authenticated
  if (!auth.isAuthenticated) {
    router.push('/login');
    return;
  }

  await fetchCollections();
  await fetchDashboardStats();
});
</script>

<template>
  <AppLayout>
    <AppHeader>
      <template #breadcrumb>
        <div class="flex items-center text-sm text-text-muted">
          <span class="font-medium text-text">Dashboard</span>
        </div>
      </template>
    </AppHeader>

    <div class="flex-1 overflow-auto min-h-0 p-4 sm:p-8 pb-24 sm:pb-8">
      <div class="max-w-7xl mx-auto space-y-8">
        <div>
          <h1 class="text-2xl font-bold text-text tracking-tight">Dashboard</h1>
          <p class="mt-1 text-sm text-text-muted">Overview of your Vault instance</p>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6">
          <div class="bg-surface-dark p-5 rounded-lg border border-border shadow-sm">
            <div class="flex justify-between items-start">
              <div>
                <p class="text-xs font-medium text-text-muted uppercase tracking-wider">
                  Collections
                </p>
                <h3 class="mt-1 text-2xl font-bold text-text">{{ collections.length }}</h3>
              </div>
              <span class="p-2 rounded bg-primary/10 text-primary">
                <FolderOpen class="w-5 h-5" />
              </span>
            </div>
          </div>

          <div class="bg-surface-dark p-5 rounded-lg border border-border shadow-sm">
            <div class="flex justify-between items-start">
              <div>
                <p class="text-xs font-medium text-text-muted uppercase tracking-wider">
                  Total Records
                </p>
                <h3 class="mt-1 text-2xl font-bold text-text">{{ totalRecords }}</h3>
              </div>
              <span class="p-2 rounded bg-primary/10 text-primary">
                <Database class="w-5 h-5" />
              </span>
            </div>
          </div>

          <div class="bg-surface-dark p-5 rounded-lg border border-border shadow-sm">
            <div class="flex justify-between items-start">
              <div>
                <p class="text-xs font-medium text-text-muted uppercase tracking-wider">
                  API Requests
                </p>
                <h3 class="mt-1 text-2xl font-bold text-text">{{ apiRequests }}</h3>
              </div>
              <span class="p-2 rounded bg-primary/10 text-primary">
                <TrendingUp class="w-5 h-5" />
              </span>
            </div>
          </div>

          <div class="bg-surface-dark p-5 rounded-lg border border-border shadow-sm">
            <div class="flex justify-between items-start">
              <div>
                <p class="text-xs font-medium text-text-muted uppercase tracking-wider">Storage</p>
                <h3 class="mt-1 text-2xl font-bold text-text">{{ storage }}</h3>
              </div>
              <span class="p-2 rounded bg-primary/10 text-primary">
                <Cloud class="w-5 h-5" />
              </span>
            </div>
          </div>
        </div>

        <div>
          <h2 class="text-lg font-semibold text-text mb-4">Quick Actions</h2>
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            <Button
              variant="secondary"
              class="!flex !items-center !gap-4 !p-4 !justify-start !h-auto"
              @click="router.push('/collections/new')"
            >
              <div class="p-3 rounded-lg bg-primary/10 text-primary">
                <Plus class="w-6 h-6" />
              </div>
              <div>
                <h3 class="font-medium text-text text-left">New Collection</h3>
                <p class="text-xs text-text-muted mt-1 text-left">Create a new data schema</p>
              </div>
            </Button>

            <Button
              variant="secondary"
              class="!flex !items-center !gap-4 !p-4 !justify-start !h-auto"
            >
              <div class="p-3 rounded-lg bg-primary/10 text-primary">
                <Terminal class="w-6 h-6" />
              </div>
              <div>
                <h3 class="font-medium text-text text-left">View Logs</h3>
                <p class="text-xs text-text-muted mt-1 text-left">Check system activity</p>
              </div>
            </Button>

            <Button
              variant="secondary"
              class="!flex !items-center !gap-4 !p-4 !justify-start !h-auto"
            >
              <div class="p-3 rounded-lg bg-primary/10 text-primary">
                <Settings class="w-6 h-6" />
              </div>
              <div>
                <h3 class="font-medium text-text text-left">Settings</h3>
                <p class="text-xs text-text-muted mt-1 text-left">Configure your instance</p>
              </div>
            </Button>
          </div>
        </div>

        <div>
          <div class="py-4 flex items-center justify-between">
            <h2 class="text-lg font-semibold text-text">Recent Collections</h2>
            <Button variant="link" class="text-xs sm:text-sm" @click="router.push('/collections')"
              >View All</Button
            >
          </div>

          <Table
            :headers="[
              { key: 'name', label: 'Name' },
              { key: 'type', label: 'Type' },
              { key: 'fields', label: 'Fields', align: 'right' },
              { key: 'created', label: 'Created', align: 'right' },
            ]"
            :items="
              collections.filter((c) => !c.name.startsWith('_')).slice(0, 5) as Record<
                string,
                unknown
              >[]
            "
            row-clickable
            @row-click="(col: Record<string, unknown>) => router.push(`/collections/${col.name}`)"
            :enable-pagination="false"
          >
            <template #cell(name)="{ item }">
              <div class="flex items-center gap-3">
                <div class="p-2 rounded bg-primary/10 text-primary">
                  <FolderOpen class="w-4 h-4" />
                </div>
                <span class="font-medium text-text">{{ item.name }}</span>
              </div>
            </template>
            <template #cell(type)="{ item }">
              <span class="text-text-muted truncate block max-w-[150px] sm:max-w-none"
                >{{ item.type }} collection</span
              >
            </template>
            <template #cell(fields)="{ item }">
              <span class="text-text-muted shrink-0"
                >{{ (item as unknown as Collection).fields?.length || 0 }} fields</span
              >
            </template>
            <template #cell(created)="{ item }">
              <span class="text-text-muted text-xs">
                {{
                  item.created ? new Date(item.created as string).toLocaleDateString() : '-'
                }}
              </span>
            </template>
            <template #empty>
              <div class="py-6 text-center text-text-muted">
                <FolderOpen class="w-12 h-12 mx-auto mb-3 opacity-30" />
                <p class="text-sm">No collections yet</p>
                <Button variant="link" class="mt-4" @click="router.push('/collections/new')"
                  >Create your first collection</Button
                >
              </div>
            </template>
          </Table>
        </div>
      </div>
    </div>
  </AppLayout>
</template>
