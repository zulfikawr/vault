<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useRouter, useRoute, RouterLink } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { useUIStore } from '../stores/ui';
import axios from 'axios';
import ConfirmModal from './ConfirmModal.vue';
import {
  LayoutDashboard,
  FolderOpen,
  Terminal,
  Settings,
  Cloud,
  LogOut,
  ChevronLeft,
  ChevronRight,
  Menu,
  Database,
  X,
  Search,
  Bell
} from 'lucide-vue-next';

interface Collection {
  name: string;
}

const auth = useAuthStore();
const ui = useUIStore();
const router = useRouter();
const route = useRoute();
const collections = ref<Collection[]>([]);
const showLogoutModal = ref(false);

const isActive = (path: string) => {
  if (path === '/') return route.path === '/';
  return route.path.startsWith(path);
};

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

const navItems = [
  { name: 'Dashboard', path: '/', icon: LayoutDashboard },
  { name: 'Collections', path: '/collections', icon: FolderOpen },
  { name: 'Storage', path: '/storage', icon: Cloud },
  { name: 'SQL Editor', path: '/sql-editor', icon: Database },
  { name: 'Logs', path: '/logs', icon: Terminal },
  { name: 'Settings', path: '/settings', icon: Settings },
];
</script>

<template>
  <div class="flex h-screen h-[100dvh] bg-background text-text overflow-hidden relative font-mono selection:bg-primary/30 selection:text-primary">
    <ConfirmModal
      :show="showLogoutModal"
      title="Terminate Session"
      message="Are you sure you want to log out? Active connections and background tasks may be interrupted."
      confirm-text="Logout"
      cancel-text="Cancel"
      variant="danger"
      @confirm="handleLogout"
      @cancel="showLogoutModal = false"
    />

    <!-- Mobile Sidebar Backdrop -->
    <Transition
      enter-active-class="transition-opacity duration-300 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition-opacity duration-200 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="ui.isMobileMenuOpen"
        class="fixed inset-0 bg-black/60 backdrop-blur-sm z-40 lg:hidden"
        @click="ui.closeMobileMenu"
      ></div>
    </Transition>

    <!-- Sidebar (Desktop & Mobile Drawer) -->
    <aside
      :class="[
        ui.sidebarCollapsed ? 'lg:w-[72px]' : 'lg:w-64',
        ui.isMobileMenuOpen ? 'translate-x-0 w-72' : '-translate-x-full lg:translate-x-0',
      ]"
      class="fixed lg:static inset-y-0 left-0 z-50 flex-shrink-0 border-r border-border bg-surface/80 backdrop-blur-md flex flex-col justify-between transition-all duration-300 ease-in-out shadow-2xl lg:shadow-none"
    >
      <div class="flex flex-col h-full overflow-hidden">
        <!-- Brand & Toggle -->
        <div
          class="h-16 flex items-center border-b border-border/50 px-4 shrink-0"
          :class="ui.sidebarCollapsed ? 'lg:justify-center lg:px-0' : 'justify-between'"
        >
          <div v-if="!ui.sidebarCollapsed || ui.isMobileMenuOpen" class="flex items-center gap-3 overflow-hidden">
            <div class="w-8 h-8 rounded-lg bg-primary/20 flex items-center justify-center border border-primary/30">
              <img src="/favicon-32x32.png" class="w-5 h-5 flex-shrink-0" alt="Vault logo" />
            </div>
            <span class="font-bold text-xl tracking-tight text-text group flex items-center gap-1">
              Vault <span class="text-[10px] bg-primary/20 text-primary px-1 rounded leading-none border border-primary/20 font-medium">v1</span>
            </span>
          </div>
          
          <div class="hidden lg:block">
            <button 
              class="w-8 h-8 flex items-center justify-center rounded-md hover:bg-white/5 transition-colors text-text-muted hover:text-primary"
              @click="ui.toggleSidebar"
            >
              <ChevronLeft v-if="!ui.sidebarCollapsed" class="w-5 h-5" />
              <ChevronRight v-else class="w-5 h-5" />
            </button>
          </div>

          <button 
            class="lg:hidden w-8 h-8 flex items-center justify-center rounded-md hover:bg-white/5 text-text-muted"
            @click="ui.closeMobileMenu"
          >
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Main Navigation -->
        <nav class="flex-1 overflow-y-auto py-4 px-3 space-y-1 custom-scrollbar">
          <div v-if="!ui.sidebarCollapsed || ui.isMobileMenuOpen" class="text-[10px] font-bold text-text-dim tracking-wider mb-3 px-3">
            System
          </div>
          
          <RouterLink
            v-for="item in navItems"
            :key="item.path"
            :to="item.path"
            :class="[
              isActive(item.path) 
                ? 'bg-primary/10 text-primary border-primary/20' 
                : 'text-text-muted hover:text-text hover:bg-white/5 border-transparent',
              ui.sidebarCollapsed && !ui.isMobileMenuOpen ? 'justify-center p-0 h-10 w-10 mx-auto' : 'px-3 py-2 gap-3',
            ]"
            class="flex items-center rounded-lg border transition-all duration-200 group relative"
            @click="ui.closeMobileMenu"
          >
            <component :is="item.icon" class="w-5 h-5 flex-shrink-0" />
            <span v-if="!ui.sidebarCollapsed || ui.isMobileMenuOpen" class="text-sm font-medium">{{ item.name }}</span>
            
            <!-- Tooltip for collapsed sidebar -->
            <div v-if="ui.sidebarCollapsed && !ui.isMobileMenuOpen" class="fixed left-[80px] px-2 py-1 bg-surface-dark border border-border rounded text-xs text-text opacity-0 group-hover:opacity-100 pointer-events-none transition-opacity z-[100] whitespace-nowrap shadow-xl">
              {{ item.name }}
            </div>
          </RouterLink>

          <!-- Collections Section -->
          <div class="pt-6">
            <div 
              v-if="!ui.sidebarCollapsed || ui.isMobileMenuOpen" 
              class="text-[10px] font-bold text-text-dim tracking-wider mb-3 px-3 flex items-center justify-between"
            >
              <span>Recent Collections</span>
              <button class="hover:text-primary transition-colors"><Search class="w-3 h-3" /></button>
            </div>
            
            <div class="space-y-0.5">
              <RouterLink
                v-for="col in collections.filter((c) => !c.name.startsWith('_')).slice(0, 5)"
                :key="col.name"
                :to="`/collections/${col.name}`"
                :class="[
                  route.params.name === col.name
                    ? 'bg-primary/10 text-primary border-primary/20'
                    : 'text-text-muted hover:text-text hover:bg-white/5 border-transparent',
                  ui.sidebarCollapsed && !ui.isMobileMenuOpen ? 'justify-center p-0 h-10 w-10 mx-auto' : 'px-3 py-1.5 gap-3',
                ]"
                class="flex items-center rounded-lg border transition-all duration-200 group relative"
              >
                <div v-if="ui.sidebarCollapsed && !ui.isMobileMenuOpen" class="w-5 h-5 flex items-center justify-center text-[10px] font-bold border border-current rounded">
                  {{ col.name.charAt(0).toUpperCase() }}
                </div>
                <div v-else class="w-1.5 h-1.5 rounded-full bg-primary/40 group-hover:bg-primary transition-colors"></div>
                <span v-if="!ui.sidebarCollapsed || ui.isMobileMenuOpen" class="text-xs truncate">{{ col.name }}</span>

                <div v-if="ui.sidebarCollapsed && !ui.isMobileMenuOpen" class="fixed left-[80px] px-2 py-1 bg-surface-dark border border-border rounded text-xs text-text opacity-0 group-hover:opacity-100 pointer-events-none transition-opacity z-[100] whitespace-nowrap shadow-xl">
                  {{ col.name }}
                </div>
              </RouterLink>
            </div>
          </div>
        </nav>

        <!-- User Profile -->
        <div class="p-3 border-t border-border/50 bg-surface-dark/30">
          <div
            v-if="!ui.sidebarCollapsed || ui.isMobileMenuOpen"
            class="flex items-center gap-3 w-full p-2 rounded-xl bg-white/5 border border-white/5 mb-2"
          >
            <div class="w-9 h-9 rounded-lg bg-primary/20 border border-primary/30 flex items-center justify-center text-primary font-bold shadow-inner">
              {{ auth.user?.data?.username?.charAt(0).toUpperCase() || 'A' }}
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-xs font-bold truncate text-text">
                {{ auth.user?.data?.username || 'Admin' }}
              </p>
              <div class="flex items-center gap-1.5">
                <div class="w-1.5 h-1.5 rounded-full bg-success"></div>
                <p class="text-[10px] text-text-dim truncate">Connected</p>
              </div>
            </div>
          </div>
          <div v-else class="flex justify-center mb-2">
            <div class="w-10 h-10 rounded-lg bg-primary/20 border border-primary/30 flex items-center justify-center text-primary font-bold shadow-inner cursor-pointer hover:bg-primary/30 transition-colors">
              {{ auth.user?.data?.username?.charAt(0).toUpperCase() || 'A' }}
            </div>
          </div>
          
          <button
            class="group w-full flex items-center justify-center gap-2 p-3 rounded-lg text-error hover:bg-error/10 border border-transparent hover:border-error/20 transition-all duration-200"
            :class="ui.sidebarCollapsed && !ui.isMobileMenuOpen ? 'p-0 h-10 w-10 mx-auto' : ''"
            @click="showLogoutModal = true"
          >
                      <LogOut class="w-4 h-4 flex-shrink-0" />
                      <span v-if="!ui.sidebarCollapsed || ui.isMobileMenuOpen" class="text-xs font-bold tracking-wide">Sign Out</span>
                    </button>
            
        </div>
      </div>
    </aside>

    <!-- Main Content Area -->
    <main class="flex-1 flex flex-col min-w-0 min-h-0 bg-background overflow-hidden relative">
      <!-- Mobile Top Header -->
      <div class="lg:hidden h-16 flex items-center justify-between px-4 border-b border-border bg-surface shrink-0 z-30">
        <button 
          class="w-10 h-10 flex items-center justify-center rounded-xl bg-surface-dark border border-border text-text hover:text-primary active:scale-95 transition-all"
          @click="ui.toggleMobileMenu"
        >
          <Menu class="h-5 w-5" />
        </button>
        
        <div class="flex items-center gap-2">
          <img src="/favicon-32x32.png" class="w-6 h-6" alt="Vault logo" />
          <span class="font-bold text-lg tracking-tighter text-text uppercase">Vault</span>
        </div>

        <button class="w-10 h-10 flex items-center justify-center rounded-xl bg-surface-dark border border-border text-text-muted hover:text-primary">
          <Bell class="h-5 w-5" />
        </button>
      </div>

      <!-- Content Slot -->
      <div class="flex-1 overflow-auto custom-scrollbar relative">
        <slot />
        
        <!-- Mobile Bottom Spacing for Nav -->
        <div class="h-20 lg:hidden"></div>
      </div>

      <!-- Mobile Bottom Navigation -->
      <nav class="lg:hidden fixed bottom-4 left-4 right-4 h-16 bg-surface/90 backdrop-blur-xl border border-white/10 rounded-2xl shadow-2xl z-40 flex items-center justify-around px-2">
        <RouterLink
          v-for="item in navItems.slice(0, 4)"
          :key="item.path"
          :to="item.path"
          :class="[
            isActive(item.path) ? 'text-primary' : 'text-text-muted',
          ]"
          class="flex flex-col items-center justify-center w-14 h-12 transition-colors relative"
        >
          <component :is="item.icon" class="w-5 h-5 mb-1" />
          <span class="text-[9px] font-bold uppercase tracking-tighter">{{ item.name.split(' ')[0] }}</span>
          <div v-if="isActive(item.path)" class="absolute -bottom-1 w-1 h-1 rounded-full bg-primary shadow-[0_0_8px_rgba(var(--primary-rgb),0.8)]"></div>
        </RouterLink>
        <button 
          class="flex flex-col items-center justify-center w-14 h-12 text-text-muted hover:text-primary transition-colors"
          @click="ui.toggleMobileMenu"
        >
          <Settings class="w-5 h-5 mb-1" />
          <span class="text-[9px] font-bold uppercase tracking-tighter">More</span>
        </button>
      </nav>
    </main>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: var(--border);
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: var(--primary);
}
</style>
