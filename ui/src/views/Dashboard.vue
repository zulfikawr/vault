<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import { 
  FolderOpen,
  Search,
  Bell,
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
    <header class="h-16 flex items-center justify-between px-8 border-b border-border bg-surface z-10">
      <div class="flex items-center text-sm text-text-muted">
        <span class="font-medium text-text">Dashboard</span>
      </div>
      
      <div class="flex items-center gap-4">
        <span class="inline-flex items-center px-2.5 py-0.5 rounded text-xs font-medium bg-success/10 text-success border border-success/20">
          <span class="w-1.5 h-1.5 rounded-full bg-success mr-1.5"></span>
          Production
        </span>
        
        <div class="relative hidden sm:block group">
          <input 
            type="text" 
            placeholder="Search (Ctrl+K)" 
            class="w-64 bg-surface-dark border-border rounded-md py-1.5 pl-9 pr-3 text-sm focus:ring-1 focus:ring-primary focus:border-primary text-text placeholder-text-muted transition-all"
          />
          <Search class="absolute left-2.5 top-2 text-text-muted w-4 h-4" />
        </div>
        
        <button class="relative p-1.5 text-text-muted hover:text-text rounded-md hover:bg-surface-dark transition-colors">
          <Bell class="w-5 h-5" />
          <span class="absolute top-1.5 right-1.5 w-2 h-2 rounded-full bg-primary border-2 border-surface"></span>
        </button>
      </div>
    </header>

    <div class="flex-1 overflow-auto p-8">
      <div class="max-w-7xl mx-auto space-y-8">
        <div>
          <h1 class="text-2xl font-bold text-text tracking-tight">Dashboard</h1>
          <p class="mt-1 text-sm text-text-muted">Overview of your Vault instance</p>
        </div>

        <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
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

        <div class="bg-surface-dark border border-border rounded-lg p-6">
          <h2 class="text-lg font-semibold text-text mb-4">Quick Actions</h2>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <button @click="router.push('/collections/new')" class="flex items-center gap-4 p-4 bg-surface border border-border rounded-lg hover:border-primary transition-colors text-left">
              <div class="p-3 rounded-lg bg-primary/10 text-primary">
                <Plus class="w-6 h-6" />
              </div>
              <div>
                <h3 class="font-medium text-text">New Collection</h3>
                <p class="text-xs text-text-muted mt-1">Create a new data schema</p>
              </div>
            </button>

            <button class="flex items-center gap-4 p-4 bg-surface border border-border rounded-lg hover:border-primary transition-colors text-left">
              <div class="p-3 rounded-lg bg-primary/10 text-primary">
                <Terminal class="w-6 h-6" />
              </div>
              <div>
                <h3 class="font-medium text-text">View Logs</h3>
                <p class="text-xs text-text-muted mt-1">Check system activity</p>
              </div>
            </button>

            <button class="flex items-center gap-4 p-4 bg-surface border border-border rounded-lg hover:border-primary transition-colors text-left">
              <div class="p-3 rounded-lg bg-primary/10 text-primary">
                <Settings class="w-6 h-6" />
              </div>
              <div>
                <h3 class="font-medium text-text">Settings</h3>
                <p class="text-xs text-text-muted mt-1">Configure your instance</p>
              </div>
            </button>
          </div>
        </div>

        <div class="bg-surface-dark border border-border rounded-lg overflow-hidden">
          <div class="px-6 py-4 border-b border-border flex items-center justify-between">
            <h2 class="text-lg font-semibold text-text">Recent Collections</h2>
            <button @click="router.push('/collections')" class="text-sm text-primary hover:text-primary-hover transition-colors">View All</button>
          </div>
          <div class="divide-y divide-border">
            <div v-for="col in collections.filter(c => !c.name.startsWith('_')).slice(0, 5)" :key="col.name" @click="router.push(`/collections/${col.name}`)" class="px-6 py-4 hover:bg-surface transition-colors cursor-pointer">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <div class="p-2 rounded bg-primary/10 text-primary">
                    <FolderOpen class="w-4 h-4" />
                  </div>
                  <div>
                    <h3 class="font-medium text-text">{{ col.name }}</h3>
                    <p class="text-xs text-text-muted">{{ col.type }} collection</p>
                  </div>
                </div>
                <span class="text-xs text-text-muted">{{ col.fields?.length || 0 }} fields</span>
              </div>
            </div>
            <div v-if="collections.filter(c => !c.name.startsWith('_')).length === 0" class="px-6 py-12 text-center text-text-muted">
              <FolderOpen class="w-12 h-12 mx-auto mb-3 opacity-30" />
              <p class="text-sm">No collections yet</p>
              <button @click="router.push('/collections/new')" class="mt-4 text-sm text-primary hover:text-primary-hover">Create your first collection</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>
