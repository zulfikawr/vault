<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter, useRoute, RouterLink } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import axios from 'axios';
import ConfirmModal from './ConfirmModal.vue';
import Button from './Button.vue';
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
} from 'lucide-vue-next';

interface Collection {
  name: string;
}

const auth = useAuthStore();
const router = useRouter();
const route = useRoute();
const sidebarCollapsed = ref(localStorage.getItem('sidebarCollapsed') === 'true');
const isMobileMenuOpen = ref(false);
const collections = ref<Collection[]>([]);
const showLogoutModal = ref(false);

const isActive = (path: string) => {
  if (path === '/') return route.path === '/';
  return route.path.startsWith(path);
};

const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value;
  localStorage.setItem('sidebarCollapsed', String(sidebarCollapsed.value));
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
</script>

<template>
  <div class="flex h-screen h-[100dvh] bg-background text-text overflow-hidden relative">
    <ConfirmModal
      :show="showLogoutModal"
      title="Confirm Logout"
      message="Are you sure you want to log out? You will need to sign in again to access the admin panel."
      confirm-text="Logout"
      cancel-text="Cancel"
      variant="warning"
      @confirm="handleLogout"
      @cancel="showLogoutModal = false"
    />

    <!-- Mobile Backdrop -->
    <div
      v-if="isMobileMenuOpen"
      class="fixed inset-0 bg-black/50 z-40 lg:hidden"
      @click="isMobileMenuOpen = false"
    ></div>

    <!-- Sidebar -->
    <aside
      :class="[
        sidebarCollapsed ? 'lg:w-16' : 'lg:w-64',
        isMobileMenuOpen ? 'translate-x-0 w-64' : '-translate-x-full lg:translate-x-0',
      ]"
      class="fixed lg:static inset-y-0 left-0 z-50 flex-shrink-0 border-r border-border bg-surface flex flex-col justify-between transition-all duration-300"
    >
      <div>
        <!-- Brand -->
        <div
          class="h-16 flex items-center border-b border-border px-6"
          :class="sidebarCollapsed ? 'lg:justify-center lg:px-0' : 'justify-between'"
        >
          <span
            v-if="!sidebarCollapsed || isMobileMenuOpen"
            class="font-bold text-lg tracking-tight text-primary"
            >Vault</span
          >
          <div class="hidden lg:block">
            <Button variant="ghost" size="sm" class="!p-2" @click="toggleSidebar">
              <ChevronLeft v-if="!sidebarCollapsed" class="w-5 h-5 text-text-muted" />
              <ChevronRight v-else class="w-5 h-5 text-text-muted" />
            </Button>
          </div>
        </div>

        <!-- Navigation -->
        <nav
          :class="[
            'p-4 space-y-1',
            sidebarCollapsed ? 'lg:px-0 lg:flex lg:flex-col lg:items-center' : '',
          ]"
        >
          <Button
            as="RouterLink"
            to="/"
            variant="ghost"
            size="sm"
            :class="[
              'h-9 transition-all duration-300',
              sidebarCollapsed ? 'lg:w-9 lg:!p-0' : 'w-full !justify-start gap-3 px-3',
              isActive('/') ? '!bg-primary/10 !text-primary' : '',
            ]"
          >
            <template #leftIcon>
              <LayoutDashboard class="w-5 h-5 flex-shrink-0" />
            </template>
            <span :class="sidebarCollapsed ? 'lg:hidden' : ''">Dashboard</span>
          </Button>
          <Button
            as="RouterLink"
            to="/collections"
            variant="ghost"
            size="sm"
            :class="[
              'h-9 transition-all duration-300',
              sidebarCollapsed ? 'lg:w-9 lg:!p-0' : 'w-full !justify-start gap-3 px-3',
              isActive('/collections') ? '!bg-primary/10 !text-primary' : '',
            ]"
          >
            <template #leftIcon>
              <FolderOpen class="w-5 h-5 flex-shrink-0" />
            </template>
            <span :class="sidebarCollapsed ? 'lg:hidden' : ''">Collections</span>
          </Button>
          <Button
            as="RouterLink"
            to="/storage"
            variant="ghost"
            size="sm"
            :class="[
              'h-9 transition-all duration-300',
              sidebarCollapsed ? 'lg:w-9 lg:!p-0' : 'w-full !justify-start gap-3 px-3',
              isActive('/storage') ? '!bg-primary/10 !text-primary' : '',
            ]"
          >
            <template #leftIcon>
              <Cloud class="w-5 h-5 flex-shrink-0" />
            </template>
            <span :class="sidebarCollapsed ? 'lg:hidden' : ''">Storage</span>
          </Button>
          <Button
            as="RouterLink"
            to="/logs"
            variant="ghost"
            size="sm"
            :class="[
              'h-9 transition-all duration-300',
              sidebarCollapsed ? 'lg:w-9 lg:!p-0' : 'w-full !justify-start gap-3 px-3',
              isActive('/logs') ? '!bg-primary/10 !text-primary' : '',
            ]"
          >
            <template #leftIcon>
              <Terminal class="w-5 h-5 flex-shrink-0" />
            </template>
            <span :class="sidebarCollapsed ? 'lg:hidden' : ''">Logs</span>
          </Button>
          <Button
            as="RouterLink"
            to="/settings"
            variant="ghost"
            size="sm"
            :class="[
              'h-9 transition-all duration-300',
              sidebarCollapsed ? 'lg:w-9 lg:!p-0' : 'w-full !justify-start gap-3 px-3',
              isActive('/settings') ? '!bg-primary/10 !text-primary' : '',
            ]"
          >
            <template #leftIcon>
              <Settings class="w-5 h-5 flex-shrink-0" />
            </template>
            <span :class="sidebarCollapsed ? 'lg:hidden' : ''">Settings</span>
          </Button>
        </nav>

        <!-- Collections Quick List -->
        <div v-if="!sidebarCollapsed || isMobileMenuOpen" class="px-4 mt-6">
          <div class="text-xs font-semibold text-text-dim uppercase tracking-wider mb-2 px-3">
            Recent Collections
          </div>
          <ul class="space-y-1">
            <li
              v-for="col in collections.filter((c) => !c.name.startsWith('_')).slice(0, 3)"
              :key="col.name"
            >
              <RouterLink
                :to="`/collections/${col.name}`"
                :class="[
                  route.params.name === col.name
                    ? 'bg-primary/10 text-primary'
                    : 'text-text-muted hover:text-text',
                ]"
                class="flex items-center justify-between px-3 py-1.5 text-sm group transition-colors"
              >
                <span class="text-xs">{{ col.name }}</span>
                <span
                  class="w-1.5 h-1.5 rounded-full bg-success group-hover:scale-110 transition-transform"
                ></span>
              </RouterLink>
            </li>
          </ul>
        </div>
      </div>

      <!-- User Profile -->
      <div
        :class="[
          'p-4 border-t border-border',
          sidebarCollapsed ? 'lg:px-0 lg:flex lg:flex-col lg:items-center' : '',
        ]"
      >
        <div
          v-if="!sidebarCollapsed || isMobileMenuOpen"
          class="flex items-center gap-3 w-full p-2 rounded-lg mb-2"
        >
          <div
            class="w-8 h-8 rounded-full bg-primary flex items-center justify-center text-white font-bold text-sm"
          >
            {{ auth.user?.data?.username?.charAt(0).toUpperCase() || 'A' }}
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium truncate text-text">
              {{ auth.user?.data?.username || 'Admin' }}
            </p>
            <p class="text-xs text-text-muted truncate">
              {{ auth.user?.data?.email || 'admin@vault.local' }}
            </p>
          </div>
        </div>
        <div v-else class="flex justify-center mb-2">
          <div
            class="w-8 h-8 rounded-full bg-primary flex items-center justify-center text-white font-bold text-sm"
          >
            {{ auth.user?.data?.username?.charAt(0).toUpperCase() || 'A' }}
          </div>
        </div>
        <Button
          variant="ghost"
          size="sm"
          class="!text-error hover:!bg-error/10 h-9 transition-all duration-300"
          :class="[
            sidebarCollapsed && !isMobileMenuOpen
              ? 'lg:w-9 lg:!p-0'
              : 'w-full !justify-start gap-3 px-3',
          ]"
          @click="showLogoutModal = true"
        >
          <LogOut class="w-5 h-5 flex-shrink-0" />
          <span v-if="!sidebarCollapsed || isMobileMenuOpen" class="text-xs font-bold"
            >Sign Out</span
          >
        </Button>
      </div>
    </aside>

    <!-- Main Content Slot -->
    <main class="flex-1 flex flex-col min-w-0 min-h-0 bg-background overflow-hidden">
      <!-- Mobile Header Toggle (Visible only on mobile) -->
      <div class="lg:hidden h-16 flex items-center px-4 border-b border-border bg-surface shrink-0">
        <Button variant="ghost" size="xs" class="-ml-2" @click="isMobileMenuOpen = true">
          <Menu class="h-6 w-6 text-text" />
        </Button>
        <span class="ml-4 font-bold text-lg tracking-tight text-primary">vault</span>
      </div>
      <slot />
    </main>
  </div>
</template>
