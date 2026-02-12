<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import axios from 'axios';
import { 
  Plus, 
  Trash2, 
  Save, 
  X, 
  FolderPlus,
  LayoutDashboard, 
  FolderOpen, 
  Terminal, 
  Settings, 
  Cloud, 
  LogOut,
  ChevronLeft,
  ChevronRight
} from 'lucide-vue-next';

const router = useRouter();
const auth = useAuthStore();
const sidebarCollapsed = ref(false);

const collectionFormData = ref({
  name: '',
  type: 'base',
  fields: [{ name: 'name', type: 'text', required: true }]
});

const addField = () => {
  collectionFormData.value.fields.push({ name: '', type: 'text', required: false });
};

const removeField = (index: number) => {
  collectionFormData.value.fields.splice(index, 1);
};

const saveCollection = async () => {
  try {
    await axios.post('/api/admin/collections', collectionFormData.value);
    router.push('/');
  } catch (error) {
    console.error('Collection creation failed', error);
  }
};

const handleLogout = () => {
  auth.logout();
  router.push({ name: 'Login' });
};
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
          <a href="/" :class="sidebarCollapsed ? 'justify-center px-0' : 'gap-3 px-3'" class="flex items-center py-2 text-sm font-medium rounded-lg bg-surface-dark text-primary border-l-2 border-primary transition-colors">
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
          <span class="hover:text-text cursor-pointer" @click="router.push('/')">Collections</span>
          <span class="mx-2">/</span>
          <span class="font-medium text-text">New</span>
        </div>
      </header>

      <!-- Form Content -->
      <div class="flex-1 overflow-auto p-8">
        <div class="space-y-6">
          <div>
            <h1 class="text-2xl font-bold text-text">Create New Collection</h1>
            <p class="text-sm text-text-muted mt-1">Define your database schema and fields</p>
          </div>

          <form @submit.prevent="saveCollection" class="space-y-6">
            <!-- Basic Info Card -->
            <div class="bg-surface-dark border border-border rounded-lg p-6">
              <h2 class="text-lg font-semibold text-text mb-4 flex items-center gap-2">
                <FolderPlus class="w-5 h-5 text-primary" />
                Basic Information
              </h2>
              
              <div class="grid grid-cols-2 gap-6">
                <div>
                  <label class="block text-sm font-medium text-text mb-2">Collection Name</label>
                  <input 
                    v-model="collectionFormData.name" 
                    type="text" 
                    required 
                    placeholder="e.g. products, users, posts"
                    class="w-full bg-surface border border-border rounded-lg px-4 py-2.5 text-text placeholder-text-muted focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition-all"
                  />
                  <p class="text-xs text-text-dim mt-1">Lowercase, no spaces (use underscores)</p>
                </div>
                
                <div>
                  <label class="block text-sm font-medium text-text mb-2">Collection Type</label>
                  <select 
                    v-model="collectionFormData.type" 
                    class="w-full bg-surface border border-border rounded-lg px-4 py-2.5 text-text focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition-all"
                  >
                    <option value="base">Base (Generic Data)</option>
                    <option value="auth">Auth (User Records)</option>
                  </select>
                  <p class="text-xs text-text-dim mt-1">Choose the collection purpose</p>
                </div>
              </div>
            </div>

            <!-- Fields Card -->
            <div class="bg-surface-dark border border-border rounded-lg p-6">
              <div class="flex items-center justify-between mb-4">
                <h2 class="text-lg font-semibold text-text flex items-center gap-2">
                  <Plus class="w-5 h-5 text-primary" />
                  Fields Definition
                </h2>
                <button 
                  type="button" 
                  @click="addField" 
                  class="px-4 py-2 bg-primary/10 text-primary rounded-lg text-sm font-medium hover:bg-primary/20 transition-colors flex items-center gap-2"
                >
                  <Plus class="w-4 h-4" />
                  Add Field
                </button>
              </div>

              <div class="space-y-3">
                <div 
                  v-for="(field, index) in collectionFormData.fields" 
                  :key="index" 
                  class="flex items-center gap-3 bg-surface p-4 rounded-lg border border-border"
                >
                  <div class="flex-1">
                    <input 
                      v-model="field.name" 
                      placeholder="field_name" 
                      class="w-full bg-surface-dark border border-border rounded-lg px-3 py-2 text-sm text-text placeholder-text-muted focus:outline-none focus:ring-1 focus:ring-primary/50 focus:border-primary transition-all"
                    />
                  </div>
                  
                  <div class="w-40">
                    <select 
                      v-model="field.type" 
                      class="w-full bg-surface-dark border border-border rounded-lg px-3 py-2 text-sm text-text focus:outline-none focus:ring-1 focus:ring-primary/50 focus:border-primary transition-all"
                    >
                      <option value="text">Text</option>
                      <option value="number">Number</option>
                      <option value="bool">Boolean</option>
                      <option value="json">JSON</option>
                      <option value="file">File</option>
                    </select>
                  </div>

                  <label class="flex items-center gap-2 text-sm text-text-muted cursor-pointer">
                    <input 
                      v-model="field.required" 
                      type="checkbox" 
                      class="w-4 h-4 text-primary bg-surface-dark border-border rounded focus:ring-primary"
                    />
                    Required
                  </label>
                  
                  <button 
                    type="button" 
                    @click="removeField(index)" 
                    class="p-2 text-error hover:bg-error/10 rounded-lg transition-colors"
                    :disabled="collectionFormData.fields.length === 1"
                  >
                    <Trash2 class="w-4 h-4" />
                  </button>
                </div>
              </div>
            </div>

            <!-- Actions -->
            <div class="flex items-center justify-end gap-3">
              <button 
                type="button" 
                @click="router.push('/')" 
                class="px-6 py-2.5 bg-surface border border-border text-text rounded-lg font-medium hover:bg-surface-dark transition-colors"
              >
                Cancel
              </button>
              <button 
                type="submit" 
                class="px-6 py-2.5 bg-primary hover:bg-primary-hover text-white rounded-lg font-medium transition-colors flex items-center gap-2"
              >
                <Save class="w-4 h-4" />
                Create Collection
              </button>
            </div>
          </form>
        </div>
      </div>
    </main>
  </div>
</template>
