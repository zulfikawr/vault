<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import { 
  FolderOpen, 
  Terminal, 
  Settings, 
  Cloud
  Search,
  Bell,
  Filter,
  Plus,
  MoreHorizontal} from 'lucide-vue-next';

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

const handleLogout = () => {
  auth.logout();
  router.push({ name: 'Login' });
};

onMounted(fetchCollections);
</script>

<template>
  <AppLayout>
    
      <!-- Header -->
      <header class="h-16 flex items-center justify-between px-8 border-b border-border bg-surface z-10">
        <!-- Breadcrumbs -->
        <div class="flex items-center text-sm text-text-muted">
          <span class="hover:text-text cursor-pointer" @click="router.push('/')">Vault</span>
          <span class="mx-2">/</span>
          <span class="font-medium text-text">Collections</span>
        </div>
        
        <!-- Actions -->
        <div class="flex items-center gap-4">
          <!-- Environment Badge -->
          <span class="inline-flex items-center px-2.5 py-0.5 rounded text-xs font-medium bg-success/10 text-success border border-success/20">
            <span class="w-1.5 h-1.5 rounded-full bg-success mr-1.5"></span>
            Production
          </span>
          
          <!-- Search -->
          <div class="relative hidden sm:block group">
            <input 
              type="text" 
              placeholder="Search (Ctrl+K)" 
              class="w-64 bg-surface-dark border-border rounded-md py-1.5 pl-9 pr-3 text-sm focus:ring-1 focus:ring-primary focus:border-primary text-text placeholder-text-muted transition-all"
            />
            <Search class="absolute left-2.5 top-2 text-text-muted w-4 h-4" />
          </div>
          
          <!-- Notifications -->
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
              <h1 class="text-2xl font-bold text-text tracking-tight">Data Collections</h1>
              <p class="mt-1 text-sm text-text-muted">Manage your database schemas, content types, and API endpoints.</p>
            </div>
            <div class="flex items-center gap-3">
              <button class="px-4 py-2 bg-surface-dark border border-border rounded text-sm font-medium text-text hover:bg-surface transition-colors flex items-center gap-2">
                <Filter class="w-4 h-4" />
                Filter
              </button>
              <button @click="router.push('/collections/new')" class="px-4 py-2 bg-primary hover:bg-primary-hover text-white rounded text-sm font-medium shadow-sm hover:shadow transition-all flex items-center gap-2">
                <Plus class="w-4 h-4" />
                New Collection
              </button>
            </div>
          </div>

          <!-- Data Table -->
          <div class="bg-surface-dark rounded-lg border border-border shadow-sm overflow-hidden">
            <div class="overflow-x-auto">
              <table class="w-full text-left text-sm whitespace-nowrap">
                <thead class="bg-surface border-b border-border">
                  <tr>
                    <th class="px-6 py-3 text-xs font-medium text-text-muted uppercase tracking-wider">Name</th>
                    <th class="px-6 py-3 text-xs font-medium text-text-muted uppercase tracking-wider">Type</th>
                    <th class="px-6 py-3 text-xs font-medium text-text-muted uppercase tracking-wider">Fields</th>
                    <th class="px-6 py-3 text-xs font-medium text-text-muted uppercase tracking-wider">Status</th>
                    <th class="px-6 py-3 text-right text-xs font-medium text-text-muted uppercase tracking-wider">Actions</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-border">
                  <tr v-for="col in collections" :key="col.name" @click="router.push(`/collections/${col.name}`)" class="hover:bg-background/50 transition-colors group cursor-pointer">
                    <td class="px-6 py-4">
                      <div class="flex items-center gap-3">
                        <div class="p-1.5 rounded bg-primary/10 text-primary">
                          <FolderOpen class="w-4 h-4" />
                        </div>
                        <span class="font-medium text-text">{{ col.name }}</span>
                      </div>
                    </td>
                    <td class="px-6 py-4">
                      <span class="text-text-muted">{{ col.type }}</span>
                    </td>
                    <td class="px-6 py-4">
                      <span class="text-text-muted">{{ col.fields?.length || 0 }} fields</span>
                    </td>
                    <td class="px-6 py-4">
                      <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-success/10 text-success">
                        Active
                      </span>
                    </td>
                    <td class="px-6 py-4 text-right">
                      <button class="text-text-muted hover:text-primary transition-colors">
                        <MoreHorizontal class="w-4 h-4" />
                      </button>
                    </td>
                  </tr>
                  <tr v-if="collections.length === 0">
                    <td colspan="5" class="px-6 py-12 text-center text-text-muted">
                      <FolderOpen class="w-12 h-12 mx-auto mb-3 opacity-30" />
                      <p class="text-sm mb-4">No collections found</p>
                      <button @click="router.push('/collections/new')" class="text-sm text-primary hover:text-primary-hover">Create your first collection</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            
            <!-- Pagination -->
            <div class="bg-surface px-6 py-3 border-t border-border flex items-center justify-between">
              <div class="text-xs text-text-muted">
                Showing <span class="font-medium text-text">{{ collections.length }}</span> of <span class="font-medium text-text">{{ collections.length }}</span> results
              </div>
              <div class="flex gap-2">
                <button class="px-3 py-1 text-xs font-medium rounded border border-border bg-surface-dark text-text hover:bg-surface transition-colors disabled:opacity-50" disabled>Previous</button>
                <button class="px-3 py-1 text-xs font-medium rounded border border-border bg-surface-dark text-text hover:bg-surface transition-colors">Next</button>
              </div>
            </div>
          </div>
        </div>
      </div>
      </AppLayout>
</template>

