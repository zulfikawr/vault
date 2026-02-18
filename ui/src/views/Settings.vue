<script setup lang="ts">
import { ref, onMounted } from 'vue';
import axios from 'axios';
import { useUIStore, type Theme } from '../stores/ui';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import Input from '../components/Input.vue';
import Dropdown from '../components/Dropdown.vue';
import DropdownItem from '../components/DropdownItem.vue';
import { 
  Save, 
  Palette, 
  Settings as SettingsIcon, 
  Monitor, 
  Layout, 
  Check,
  ChevronRight,
  ChevronDown,
  Shield,
  Cpu
} from 'lucide-vue-next';

interface AppSettings {
  port: number;
  log_level: string;
  log_format: string;
  jwt_expiry: number;
  max_file_upload_size: number;
  cors_origins: string;
  rate_limit_per_min: number;
  tls_enabled: boolean;
}

const ui = useUIStore();
const settings = ref<AppSettings | null>(null);
const loading = ref(false);
const saving = ref(false);
const activeTab = ref<'appearance' | 'system'>('appearance');

const fetchSettings = async () => {
  loading.value = true;
  try {
    const response = await axios.get('/api/admin/settings');
    settings.value = response.data.data || response.data;
  } catch (error) {
    console.error('Failed to fetch settings', error);
  } finally {
    loading.value = false;
  }
};

const saveSettings = async () => {
  if (!settings.value) return;
  saving.value = true;
  try {
    await axios.patch('/api/admin/settings', settings.value);
  } catch (error) {
    console.error('Failed to save settings', error);
  } finally {
    saving.value = false;
  }
};

onMounted(fetchSettings);

const themes: { name: Theme; label: string; colors: string[] }[] = [
  { name: 'gruvbox', label: 'Gruvbox', colors: ['#282828', '#ebdbb2', '#ee822b'] },
  { name: 'monokai', label: 'Monokai', colors: ['#272822', '#f8f8f2', '#f92672'] },
  { name: 'nord', label: 'Nord', colors: ['#2e3440', '#eceff4', '#88c0d0'] },
  { name: 'dracula', label: 'Dracula', colors: ['#282a36', '#f8f8f2', '#bd93f9'] },
  { name: 'solarized', label: 'Solarized', colors: ['#002b36', '#839496', '#268bd2'] },
];
</script>

<template>
  <AppLayout>
    <AppHeader>
      <template #breadcrumb>
        <div class="flex items-center gap-2 text-xs font-bold tracking-wide">
          <SettingsIcon class="w-4 h-4 text-primary" />
          <span class="text-text-dim">Settings</span>
          <ChevronRight class="w-3 h-3 text-text-dim" />
          <span class="text-primary">{{ activeTab === 'appearance' ? 'Appearance' : 'System' }}</span>
        </div>
      </template>
    </AppHeader>

    <main class="flex-1 min-h-0 p-4 lg:p-8 max-w-7xl mx-auto w-full">
      <div class="space-y-8">
        <!-- Header -->
        <div class="flex flex-col md:flex-row md:items-end justify-between gap-6 border-b border-border/30 pb-6">
          <div>
            <h1 class="text-3xl font-bold text-text tracking-tight">Settings</h1>
            <p class="text-text-dim text-xs mt-1">Manage your application preferences and server configuration.</p>
          </div>
          
          <div class="flex bg-surface-dark/50 p-1 rounded-xl border border-border/50 backdrop-blur-sm">
            <button 
              @click="activeTab = 'appearance'"
              :class="activeTab === 'appearance' ? 'bg-primary text-white shadow-lg' : 'text-text-dim hover:text-text'"
              class="px-4 py-2 rounded-lg text-xs font-bold transition-all flex items-center gap-2"
            >
              <Palette class="w-3.5 h-3.5" />
              Appearance
            </button>
            <button 
              @click="activeTab = 'system'"
              :class="activeTab === 'system' ? 'bg-primary text-white shadow-lg' : 'text-text-dim hover:text-text'"
              class="px-4 py-2 rounded-lg text-xs font-bold transition-all flex items-center gap-2"
            >
              <Cpu class="w-3.5 h-3.5" />
              System
            </button>
          </div>
        </div>

        <!-- Appearance Tab -->
        <div v-if="activeTab === 'appearance'" class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
          <section class="space-y-6">
            <div class="flex items-center gap-3">
              <Monitor class="w-4 h-4 text-primary" />
              <h2 class="text-sm font-bold text-text">Themes</h2>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
              <button 
                v-for="theme in themes" 
                :key="theme.name"
                @click="ui.setTheme(theme.name)"
                :class="ui.theme === theme.name ? 'border-primary ring-4 ring-primary/10' : 'border-border/50 hover:border-primary/30'"
                class="group relative flex flex-col p-4 rounded-2xl bg-surface/30 border-2 transition-all text-left overflow-hidden"
              >
                <div class="flex items-center justify-between mb-4">
                  <span class="text-xs font-bold">{{ theme.label }}</span>
                  <div v-if="ui.theme === theme.name" class="w-5 h-5 rounded-full bg-primary flex items-center justify-center">
                    <Check class="w-3 h-3 text-white" />
                  </div>
                </div>
                <div class="flex gap-2">
                  <div v-for="color in theme.colors" :key="color" :style="{ backgroundColor: color }" class="w-8 h-8 rounded-lg border border-white/10 shadow-inner"></div>
                </div>
                <div class="mt-4 p-2 rounded-lg bg-black/20 border border-white/5 space-y-1.5">
                  <div class="h-1 w-2/3 rounded-full" :style="{ backgroundColor: theme.colors[2] }"></div>
                  <div class="h-1 w-full rounded-full opacity-30" :style="{ backgroundColor: theme.colors[1] }"></div>
                </div>
              </button>
            </div>
          </section>

          <section class="space-y-6">
            <div class="flex items-center gap-3">
              <Layout class="w-4 h-4 text-secondary" />
              <h2 class="text-sm font-bold text-text">Layout Settings</h2>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="p-6 rounded-2xl bg-surface-dark/50 border border-border flex items-center justify-between">
                <div>
                  <h3 class="text-xs font-bold text-text">Collapse Sidebar</h3>
                  <p class="text-[10px] text-text-dim mt-1">Use a minimized navigation bar</p>
                </div>
                <button 
                  @click="ui.toggleSidebar"
                  :class="ui.sidebarCollapsed ? 'bg-primary border-primary' : 'bg-surface border-border'"
                  class="w-12 h-6 rounded-full border-2 relative transition-all duration-300"
                >
                  <div 
                    :class="ui.sidebarCollapsed ? 'translate-x-6' : 'translate-x-0'"
                    class="absolute top-0.5 left-0.5 w-4 h-4 rounded-full bg-white transition-transform duration-300 shadow-md"
                  ></div>
                </button>
              </div>

              <div class="p-6 rounded-2xl bg-surface-dark/50 border border-border flex items-center justify-between">
                <div>
                  <h3 class="text-xs font-bold text-text">Compact Display</h3>
                  <p class="text-[10px] text-text-dim mt-1">Reduce interface spacing and padding</p>
                </div>
                <button 
                  @click="ui.toggleCompactMode"
                  :class="ui.compactMode ? 'bg-primary border-primary' : 'bg-surface border-border'"
                  class="w-12 h-6 rounded-full border-2 relative transition-all duration-300"
                >
                  <div 
                    :class="ui.compactMode ? 'translate-x-6' : 'translate-x-0'"
                    class="absolute top-0.5 left-0.5 w-4 h-4 rounded-full bg-white transition-transform duration-300 shadow-md"
                  ></div>
                </button>
              </div>
            </div>
          </section>
        </div>

        <!-- System Tab -->
        <div v-if="activeTab === 'system'" class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
          <div v-if="loading" class="py-20 text-center">
            <div class="w-10 h-10 border-2 border-primary/20 border-t-primary rounded-full animate-spin mx-auto mb-4"></div>
            <p class="text-xs font-bold text-text-dim">Loading server configuration...</p>
          </div>

          <template v-else-if="settings">
            <section class="space-y-6">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <SettingsIcon class="w-4 h-4 text-primary" />
                  <h2 class="text-sm font-bold text-text">Server Configuration</h2>
                </div>
                <Button variant="primary" size="sm" :loading="saving" @click="saveSettings">
                  <Save class="w-4 h-4 mr-2" />
                  Save Changes
                </Button>
              </div>

              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="space-y-2">
                  <label class="text-[10px] font-bold text-text-dim ml-1">Server Port</label>
                  <Input
                    :model-value="String(settings.port)"
                    type="number"
                    @update:model-value="settings!.port = Number($event)"
                  />
                </div>
                
                <div class="space-y-2">
                  <label class="text-[10px] font-bold text-text-dim ml-1">TLS / HTTPS</label>
                  <Dropdown
                    :model-value="settings.tls_enabled ? 'true' : 'false'"
                    class="w-full"
                    headless
                    @update:model-value="settings!.tls_enabled = $event === 'true'"
                  >
                    <template #trigger>
                      <div class="w-full text-left px-4 h-10 rounded-xl bg-surface-dark border border-border/50 text-sm font-bold flex justify-between items-center hover:border-primary/30 transition-all cursor-pointer">
                        <span>{{ settings.tls_enabled ? 'Enabled' : 'Disabled' }}</span>
                        <ChevronDown class="w-4 h-4 text-text-dim" />
                      </div>
                    </template>
                    <template #default="{ close }">
                      <DropdownItem value="true" @click="settings!.tls_enabled = true; close()">Enabled</DropdownItem>
                      <DropdownItem value="false" @click="settings!.tls_enabled = false; close()">Disabled</DropdownItem>
                    </template>
                  </Dropdown>
                </div>
              </div>
            </section>

            <section class="space-y-6">
              <div class="flex items-center gap-3">
                <Shield class="w-4 h-4 text-success" />
                <h2 class="text-sm font-bold text-text">Security Settings</h2>
              </div>

              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="space-y-2">
                  <label class="text-[10px] font-bold text-text-dim ml-1">JWT Expiration (Hours)</label>
                  <Input
                    :model-value="String(settings.jwt_expiry)"
                    type="number"
                    @update:model-value="settings!.jwt_expiry = Number($event)"
                  />
                </div>
                <div class="space-y-2">
                  <label class="text-[10px] font-bold text-text-dim ml-1">Rate Limit (Requests/Min)</label>
                  <Input
                    :model-value="String(settings.rate_limit_per_min)"
                    type="number"
                    @update:model-value="settings!.rate_limit_per_min = Number($event)"
                  />
                </div>
                <div class="md:col-span-2 space-y-2">
                  <label class="text-[10px] font-bold text-text-dim ml-1">Allowed CORS Origins</label>
                  <Input v-model="settings.cors_origins" placeholder="e.g. *, http://localhost:3000" />
                </div>
              </div>
            </section>
          </template>
        </div>
      </div>
    </main>
  </AppLayout>
</template>
