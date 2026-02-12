<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';
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
  Database,
  TrendingUp,
  CheckCircle,
  Users,
  FileText,
  Lock,
  ShoppingCart,
  Image,
  MoreHorizontal,
  ChevronLeft,
  ChevronRight
} from 'lucide-vue-next';

const auth = useAuthStore();
const router = useRouter();
const collections = ref([]);
const selectedCollection = ref(null);
const sidebarCollapsed = ref(false);

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
          <a href="/" :class="sidebarCollapsed ? 'justify-center px-0' : 'gap-3 px-3'" class="flex items-center py-2 text-sm font-medium rounded-lg bg-surface-dark text-primary border-l-2 border-primary transition-colors">
            <LayoutDashboard class="w-5 h-5 flex-shrink-0" />
            <span v-if="!sidebarCollapsed">Dashboard</span>
          </a>
          <a href="/collections" :class="sidebarCollapsed ? 'justify-center px-0' : 'gap-3 px-3'" class="flex items-center py-2 text-sm font-medium rounded-lg text-text-muted hover:bg-surface-dark hover:text-text transition-colors">
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
              <a href="#" class="flex items-center justify-between px-3 py-1.5 text-sm text-text-muted hover:text-text group transition-colors">
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
          <span class="font-medium text-text">Dashboard</span>
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
          <div>
            <h1 class="text-2xl font-bold text-text tracking-tight">Dashboard</h1>
            <p class="mt-1 text-sm text-text-muted">Overview of your Vault instance</p>
          </div>

          <!-- Stats Cards -->
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

          <!-- Quick Actions -->
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

          <!-- Recent Collections -->
          <div class="bg-surface-dark border border-border rounded-lg overflow-hidden">
            <div class="px-6 py-4 border-b border-border flex items-center justify-between">
              <h2 class="text-lg font-semibold text-text">Recent Collections</h2>
              <button @click="router.push('/collections')" class="text-sm text-primary hover:text-primary-hover transition-colors">View All</button>
            </div>
            <div class="divide-y divide-border">
              <div v-for="col in collections.slice(0, 5)" :key="col.name" class="px-6 py-4 hover:bg-surface transition-colors cursor-pointer">
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
              <div v-if="collections.length === 0" class="px-6 py-12 text-center text-text-muted">
                <FolderOpen class="w-12 h-12 mx-auto mb-3 opacity-30" />
                <p class="text-sm">No collections yet</p>
                <button @click="router.push('/collections/new')" class="mt-4 text-sm text-primary hover:text-primary-hover">Create your first collection</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>
