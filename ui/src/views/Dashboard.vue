<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import { 
  FolderOpen,
  Plus,
  Terminal,
  Settings,
  Database,
  TrendingUp,
  Cloud
} from 'lucide-vue-next';

const router = useRouter();
const collections = ref([]);

const fetchCollections = async () => {
  try {
    const response = await axios.get('/api/admin/collections');
    collections.value = response.data.data;
  } catch (error) {
    console.error('Failed to fetch collections', error);
  }
};

onMounted(fetchCollections);
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
                <p class="text-xs font-medium text-text-muted uppercase tracking-wider">Collections</p>
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
                <p class="text-xs font-medium text-text-muted uppercase tracking-wider">Total Records</p>
                <h3 class="mt-1 text-2xl font-bold text-text">-</h3>
              </div>
              <span class="p-2 rounded bg-primary/10 text-primary">
                <Database class="w-5 h-5" />
              </span>
            </div>
          </div>

          <div class="bg-surface-dark p-5 rounded-lg border border-border shadow-sm">
            <div class="flex justify-between items-start">
              <div>
                <p class="text-xs font-medium text-text-muted uppercase tracking-wider">API Requests</p>
                <h3 class="mt-1 text-2xl font-bold text-text">-</h3>
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
                <h3 class="mt-1 text-2xl font-bold text-text">-</h3>
              </div>
              <span class="p-2 rounded bg-primary/10 text-primary">
                <Cloud class="w-5 h-5" />
              </span>
            </div>
          </div>
        </div>

        <div class="bg-surface-dark border border-border rounded-lg p-4 sm:p-6">
          <h2 class="text-lg font-semibold text-text mb-4">Quick Actions</h2>
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            <Button @click="router.push('/collections/new')" variant="secondary" class="!flex !items-center !gap-4 !p-4 !justify-start !h-auto">
              <div class="p-3 rounded-lg bg-primary/10 text-primary">
                <Plus class="w-6 h-6" />
              </div>
              <div>
                <h3 class="font-medium text-text text-left">New Collection</h3>
                <p class="text-xs text-text-muted mt-1 text-left">Create a new data schema</p>
              </div>
            </Button>

            <Button variant="secondary" class="!flex !items-center !gap-4 !p-4 !justify-start !h-auto">
              <div class="p-3 rounded-lg bg-primary/10 text-primary">
                <Terminal class="w-6 h-6" />
              </div>
              <div>
                <h3 class="font-medium text-text text-left">View Logs</h3>
                <p class="text-xs text-text-muted mt-1 text-left">Check system activity</p>
              </div>
            </Button>

            <Button variant="secondary" class="!flex !items-center !gap-4 !p-4 !justify-start !h-auto">
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

        <div class="bg-surface-dark border border-border rounded-lg overflow-hidden">
          <div class="px-4 sm:px-6 py-4 border-b border-border flex items-center justify-between">
            <h2 class="text-lg font-semibold text-text">Recent Collections</h2>
            <Button @click="router.push('/collections')" variant="link">View All</Button>
          </div>
          <div class="divide-y divide-border">
            <div v-for="col in collections.filter(c => !c.name.startsWith('_')).slice(0, 5)" :key="col.name" @click="router.push(`/collections/${col.name}`)" class="px-4 sm:px-6 py-4 hover:bg-surface transition-colors cursor-pointer">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <div class="p-2 rounded bg-primary/10 text-primary">
                    <FolderOpen class="w-4 h-4" />
                  </div>
                  <div>
                    <h3 class="font-medium text-text">{{ col.name }}</h3>
                    <p class="text-xs text-text-muted truncate max-w-[150px] sm:max-w-none">{{ col.type }} collection</p>
                  </div>
                </div>
                <span class="text-xs text-text-muted shrink-0">{{ col.fields?.length || 0 }} fields</span>
              </div>
            </div>
            <div v-if="collections.filter(c => !c.name.startsWith('_')).length === 0" class="px-6 py-12 text-center text-text-muted">
              <FolderOpen class="w-12 h-12 mx-auto mb-3 opacity-30" />
              <p class="text-sm">No collections yet</p>
              <Button @click="router.push('/collections/new')" variant="link" class="mt-4">Create your first collection</Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>
