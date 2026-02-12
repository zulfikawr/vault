<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import axios from 'axios';
import { 
  LayoutDashboard, 
  FolderOpen, 
  Terminal, 
  Settings, 
  Cloud, 
  LogOut,
  Search,
  Bell,
  Filter,
  Plus,
  MoreHorizontal,
  ChevronLeft,
  ChevronRight,
  Edit,
  Trash2
} from 'lucide-vue-next';

const auth = useAuthStore();
const router = useRouter();
const route = useRoute();
const sidebarCollapsed = ref(false);
const collections = ref([]);
const collection = ref(null);
const records = ref([]);

const collectionName = computed(() => route.params.name as string);

const fetchCollections = async () => {
  try {
    const response = await axios.get('/api/admin/collections');
    collections.value = response.data.data;
  } catch (error) {
    console.error('Failed to fetch collections', error);
  }
};

const fetchCollection = async () => {
  try {
    const response = await axios.get(`/api/admin/collections`);
    const col = response.data.data.find((c: any) => c.name === collectionName.value);
    collection.value = col;
  } catch (error) {
    console.error('Failed to fetch collection', error);
  }
};

const fetchRecords = async () => {
  try {
    const response = await axios.get(`/api/collections/${collectionName.value}/records`);
    records.value = response.data.items || [];
  } catch (error) {
    console.error('Failed to fetch records', error);
  }
};

const deleteRecord = async (id: string) => {
  if (!confirm('Are you sure you want to delete this record?')) return;
  try {
    await axios.delete(`/api/collections/${collectionName.value}/records/${id}`);
    fetchRecords();
  } catch (error) {
    console.error('Delete failed', error);
  }
};

const handleLogout = () => {
  auth.logout();
  router.push({ name: 'Login' });
};

onMounted(() => {
  fetchCollections();
  fetchCollection();
  fetchRecords();
});
</script>

<template>
  <div class="flex h-screen bg-background text-text overflow-hidden">
    <!-- Sidebar -->
    <aside :class="sidebarCollapsed ? 'w-16' : 'w-64'" class="flex-shrink-0 border-r border-border bg-surface flex flex-col justify-between transition-all duration-300">
      <div>
        <!-- Brand -->
        <div class="h-16 flex items-center border-b border-border" :class="sidebarCollapsed ? 'justify-center' : 'justify-between px-6'">
          <span v-if="!sidebarCollapsed" class="font-bold text-lg tracking-tight text-primary">vault</span>
          <button @click="sidebarCollapsed = !sidebarCollapsed" class="p-1 hover:bg-surface-dark rounded transition-colors">
            <ChevronLeft v-if="!sidebarCollapsed" class="w-5 h-5 text-text-muted" />
            <ChevronRight v-else class="w-5 h-5 text-text-muted" />
          </button>
        </div>
        
        <!-- Navigation -->
        <nav class="p-4 space-y-1">
          <a href="/" :class="sidebarCollapsed ? 'justify-center px-0' : 'gap-3 px-3'" class="flex items-center py-2 text-sm font-medium rounded-lg text-text-muted hover:bg-surface-dark hover:text-text transition-colors">
            <LayoutDashboard class="w-5 h-5 flex-shrink-0" />
            <span v-if="!sidebarCollapsed">Dashboard</span>
          </a>
          <a href="/collections" :class="sidebarCollapsed ? 'justify-center px-0' : 'gap-3 px-3'" class="flex items-center py-2 text-sm font-medium rounded-lg bg-surface-dark text-primary border-l-2 border-primary transition-colors">
            <FolderOpen class="w-5 h-5 flex-shrink-0" />
            <span v-if="!sidebarCollapsed">Collections</span>
          </a>
          <a href="#" :class="sidebarCollapsed ? 'justify-center px-0' : 'gap-3 px-3'" class="flex items-center py-2 text-sm font-medium rounded-lg text-text-muted hover:bg-surface-dark hover:text-text transition-colors">
            <Terminal class="w-5 h-5 flex-shrink-0" />
            <span v-if="!sidebarCollapsed">Logs</span>
          </a>
          <a href="#" :class="sidebarCollapsed ? 'justify-center px-0' : 'gap-3 px-3'" class="flex items-center py-2 text-sm font-medium rounded-lg text-text-muted hover:bg-surface-dark hover:text-text transition-colors">
            <Settings class="w-5 h-5 flex-shrink-0" />
            <span v-if="!sidebarCollapsed">Settings</span>
          </a>
          <a href="#" :class="sidebarCollapsed ? 'justify-center px-0' : 'gap-3 px-3'" class="flex items-center py-2 text-sm font-medium rounded-lg text-text-muted hover:bg-surface-dark hover:text-text transition-colors">
            <Cloud class="w-5 h-5 flex-shrink-0" />
            <span v-if="!sidebarCollapsed">Storage</span>
          </a>
        </nav>
        
        <!-- Collections Quick List -->
        <div v-if="!sidebarCollapsed" class="px-4 mt-6">
          <div class="text-xs font-semibold text-text-dim uppercase tracking-wider mb-2 px-3">Recent Collections</div>
          <ul class="space-y-1">
            <li v-for="col in collections.filter(c => !c.name.startsWith('_')).slice(0, 3)" :key="col.name">
              <a :href="`/collections/${col.name}`" class="flex items-center justify-between px-3 py-1.5 text-sm text-text-muted hover:text-text group transition-colors">
                <span class="text-xs">{{ col.name }}</span>
                <span class="w-1.5 h-1.5 rounded-full bg-success group-hover:scale-110 transition-transform"></span>
              </a>
            </li>
          </ul>
        </div>
      </div>
      
      <!-- User Profile -->
      <div class="p-4 border-t border-border">
        <div v-if="!sidebarCollapsed" class="flex items-center gap-3 w-full p-2 rounded-lg mb-2">
          <div class="w-8 h-8 rounded-full bg-primary flex items-center justify-center text-white font-bold text-sm">
            {{ auth.user?.username?.charAt(0).toUpperCase() || 'A' }}
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium truncate text-text">{{ auth.user?.username || 'Admin' }}</p>
            <p class="text-xs text-text-muted truncate">Admin Workspace</p>
          </div>
        </div>
        <div v-else class="flex justify-center mb-2">
          <div class="w-8 h-8 rounded-full bg-primary flex items-center justify-center text-white font-bold text-sm">
            {{ auth.user?.username?.charAt(0).toUpperCase() || 'A' }}
          </div>
        </div>
        <button 
          @click="handleLogout" 
          :class="sidebarCollapsed ? 'justify-center px-0' : 'gap-3 px-3'"
          class="w-full flex items-center py-2 rounded-lg text-xs font-bold text-error hover:bg-error/10 transition-colors"
        >
          <LogOut class="w-4 h-4 flex-shrink-0" />
          <span v-if="!sidebarCollapsed">Sign Out</span>
        </button>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 flex flex-col min-w-0 bg-background overflow-hidden">
      <!-- Header -->
      <header class="h-16 flex items-center justify-between px-8 border-b border-border bg-surface z-10">
        <!-- Breadcrumbs -->
        <div class="flex items-center text-sm text-text-muted">
          <span class="hover:text-text cursor-pointer" @click="router.push('/')">Vault</span>
          <span class="mx-2">/</span>
          <span class="hover:text-text cursor-pointer" @click="router.push('/collections')">Collections</span>
          <span class="mx-2">/</span>
          <span class="font-medium text-text">{{ collectionName }}</span>
        </div>
        
        <!-- Actions -->
        <div class="flex items-center gap-4">
          <span class="inline-flex items-center px-2.5 py-0.5 rounded text-xs font-medium bg-success/10 text-success border border-success/20">
            <span class="w-1.5 h-1.5 rounded-full bg-success mr-1.5"></span>
            Production
          </span>
          
          <div class="relative hidden sm:block group">
            <input 
              type="text" 
              placeholder="Search records..." 
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

      <!-- Main Scrollable Area -->
      <div class="flex-1 overflow-auto p-8">
        <div class="max-w-7xl mx-auto space-y-8">
          <!-- Page Title -->
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div>
              <h1 class="text-2xl font-bold text-text tracking-tight">{{ collectionName }}</h1>
              <p class="mt-1 text-sm text-text-muted">{{ collection?.type }} collection • {{ collection?.fields?.length || 0 }} fields</p>
            </div>
            <div class="flex items-center gap-3">
              <button @click="router.push(`/collections/${collectionName}/settings`)" class="px-4 py-2 bg-surface-dark border border-border rounded text-sm font-medium text-text hover:bg-surface transition-colors flex items-center gap-2">
                <Settings class="w-4 h-4" />
                Settings
              </button>
              <button class="px-4 py-2 bg-surface-dark border border-border rounded text-sm font-medium text-text hover:bg-surface transition-colors flex items-center gap-2">
                <Filter class="w-4 h-4" />
                Filter
              </button>
              <button @click="router.push(`/collections/${collectionName}/new`)" class="px-4 py-2 bg-primary hover:bg-primary-hover text-white rounded text-sm font-medium shadow-sm hover:shadow transition-all flex items-center gap-2">
                <Plus class="w-4 h-4" />
                New Record
              </button>
            </div>
          </div>

          <!-- Data Table -->
          <div class="bg-surface-dark rounded-lg border border-border shadow-sm overflow-hidden">
            <div class="overflow-x-auto">
              <table class="w-full text-left text-sm whitespace-nowrap">
                <thead class="bg-surface border-b border-border">
                  <tr>
                    <th v-for="field in collection?.fields" :key="field.name" class="px-6 py-3 text-xs font-medium text-text-muted uppercase tracking-wider">
                      {{ field.name }}
                    </th>
                    <th class="px-6 py-3 text-right text-xs font-medium text-text-muted uppercase tracking-wider">Actions</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-border">
                  <tr v-for="record in records" :key="record.id" class="hover:bg-background/50 transition-colors group">
                    <td v-for="field in collection?.fields" :key="field.name" class="px-6 py-4">
                      <span class="text-text">{{ record[field.name] || '-' }}</span>
                    </td>
                    <td class="px-6 py-4 text-right">
                      <div class="flex justify-end items-center gap-2">
                        <button class="p-2 text-text-muted hover:text-primary hover:bg-primary/10 rounded transition-colors">
                          <Edit class="w-4 h-4" />
                        </button>
                        <button @click="deleteRecord(record.id)" class="p-2 text-text-muted hover:text-error hover:bg-error/10 rounded transition-colors">
                          <Trash2 class="w-4 h-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                  <tr v-if="records.length === 0">
                    <td :colspan="(collection?.fields?.length || 0) + 1" class="px-6 py-12 text-center text-text-muted">
                      <FolderOpen class="w-12 h-12 mx-auto mb-3 opacity-30" />
                      <p class="text-sm mb-4">No records found</p>
                      <button class="text-sm text-primary hover:text-primary-hover">Create your first record</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            
            <!-- Pagination -->
            <div class="bg-surface px-6 py-3 border-t border-border flex items-center justify-between">
              <div class="text-xs text-text-muted">
                Showing <span class="font-medium text-text">{{ records.length }}</span> of <span class="font-medium text-text">{{ records.length }}</span> results
              </div>
              <div class="flex gap-2">
                <button class="px-3 py-1 text-xs font-medium rounded border border-border bg-surface-dark text-text hover:bg-surface transition-colors disabled:opacity-50" disabled>Previous</button>
                <button class="px-3 py-1 text-xs font-medium rounded border border-border bg-surface-dark text-text hover:bg-surface transition-colors">Next</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>
