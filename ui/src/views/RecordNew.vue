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
  ChevronLeft,
  ChevronRight,
  Save,
  X
} from 'lucide-vue-next';

const auth = useAuthStore();
const router = useRouter();
const route = useRoute();
const sidebarCollapsed = ref(false);
const collections = ref([]);
const collection = ref(null);
const formData = ref({});

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
    // Initialize form data with empty values
    col?.fields?.forEach((field: any) => {
      formData.value[field.name] = '';
    });
  } catch (error) {
    console.error('Failed to fetch collection', error);
  }
};

const saveRecord = async () => {
  try {
    await axios.post(`/api/collections/${collectionName.value}/records`, formData.value);
    router.push(`/collections/${collectionName.value}`);
  } catch (error) {
    console.error('Save failed', error);
  }
};

const handleLogout = () => {
  auth.logout();
  router.push({ name: 'Login' });
};

onMounted(() => {
  fetchCollections();
  fetchCollection();
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
        <div class="flex items-center text-sm text-text-muted">
          <span class="hover:text-text cursor-pointer" @click="router.push('/')">Vault</span>
          <span class="mx-2">/</span>
          <span class="hover:text-text cursor-pointer" @click="router.push('/collections')">Collections</span>
          <span class="mx-2">/</span>
          <span class="hover:text-text cursor-pointer" @click="router.push(`/collections/${collectionName}`)">{{ collectionName }}</span>
          <span class="mx-2">/</span>
          <span class="font-medium text-text">New</span>
        </div>
      </header>

      <!-- Form Content -->
      <div class="flex-1 overflow-auto p-8">
        <div class="space-y-6">
          <div>
            <h1 class="text-2xl font-bold text-text">Create New Record</h1>
            <p class="text-sm text-text-muted mt-1">Add a new record to {{ collectionName }}</p>
          </div>

          <form @submit.prevent="saveRecord" class="space-y-6">
            <div class="bg-surface-dark border border-border rounded-lg p-6">
              <div class="space-y-6">
                <div v-for="field in collection?.fields" :key="field.name">
                  <label class="block text-sm font-medium text-text mb-2">
                    {{ field.name }}
                    <span v-if="field.required" class="text-error">*</span>
                  </label>
                  
                  <input 
                    v-if="field.type === 'text'"
                    v-model="formData[field.name]"
                    type="text"
                    :required="field.required"
                    class="w-full bg-surface border border-border rounded-lg px-4 py-2.5 text-text placeholder-text-muted focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition-all"
                  />
                  
                  <input 
                    v-else-if="field.type === 'number'"
                    v-model="formData[field.name]"
                    type="number"
                    :required="field.required"
                    class="w-full bg-surface border border-border rounded-lg px-4 py-2.5 text-text placeholder-text-muted focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition-all"
                  />
                  
                  <select 
                    v-else-if="field.type === 'bool'"
                    v-model="formData[field.name]"
                    :required="field.required"
                    class="w-full bg-surface border border-border rounded-lg px-4 py-2.5 text-text focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition-all"
                  >
                    <option value="">Select...</option>
                    <option :value="true">True</option>
                    <option :value="false">False</option>
                  </select>
                  
                  <textarea 
                    v-else-if="field.type === 'json'"
                    v-model="formData[field.name]"
                    :required="field.required"
                    rows="4"
                    placeholder='{"key": "value"}'
                    class="w-full bg-surface border border-border rounded-lg px-4 py-2.5 text-text placeholder-text-muted focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition-all font-mono text-sm"
                  ></textarea>
                  
                  <input 
                    v-else
                    v-model="formData[field.name]"
                    type="text"
                    :required="field.required"
                    class="w-full bg-surface border border-border rounded-lg px-4 py-2.5 text-text placeholder-text-muted focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition-all"
                  />
                  
                  <p class="text-xs text-text-dim mt-1">{{ field.type }}</p>
                </div>
              </div>
            </div>

            <div class="flex items-center justify-end gap-3">
              <button 
                type="button" 
                @click="router.push(`/collections/${collectionName}`)" 
                class="px-6 py-2.5 bg-surface border border-border text-text rounded-lg font-medium hover:bg-surface-dark transition-colors"
              >
                Cancel
              </button>
              <button 
                type="submit" 
                class="px-6 py-2.5 bg-primary hover:bg-primary-hover text-white rounded-lg font-medium transition-colors flex items-center gap-2"
              >
                <Save class="w-4 h-4" />
                Create Record
              </button>
            </div>
          </form>
        </div>
      </div>
    </main>
  </div>
</template>
