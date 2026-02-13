<script setup lang="ts">
import { ref, onMounted } from 'vue';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import { Settings } from 'lucide-vue-next';

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

const settings = ref<AppSettings | null>(null);
const loading = ref(false);

const fetchSettings = async () => {
  loading.value = true;
  try {
    const response = await axios.get('/api/admin/settings');
    settings.value = response.data;
  } catch (error) {
    console.error('Failed to fetch settings', error);
  } finally {
    loading.value = false;
  }
};

onMounted(fetchSettings);
</script>

<template>
  <AppLayout>
    <AppHeader>
      <template #breadcrumb>
        <div class="flex items-center gap-2">
          <span class="text-sm text-text-muted">Settings</span>
        </div>
      </template>
    </AppHeader>

    <main class="flex-1 overflow-auto">
      <div class="p-6 max-w-4xl mx-auto">
        <div class="flex items-center gap-3 mb-6">
          <Settings class="w-8 h-8 text-primary" />
          <h1 class="text-2xl font-bold text-text">System Settings</h1>
        </div>

        <div v-if="loading" class="text-center text-text-muted">Loading settings...</div>

        <div v-else-if="settings" class="space-y-6">
          <div class="bg-surface border border-border rounded-lg p-6">
            <h2 class="text-lg font-semibold text-text mb-4">Server Configuration</h2>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <p class="text-sm text-text-muted">Port</p>
                <p class="text-lg font-medium text-text">{{ settings.port }}</p>
              </div>
              <div>
                <p class="text-sm text-text-muted">TLS Enabled</p>
                <p class="text-lg font-medium text-text">
                  {{ settings.tls_enabled ? 'Yes' : 'No' }}
                </p>
              </div>
            </div>
          </div>

          <div class="bg-surface border border-border rounded-lg p-6">
            <h2 class="text-lg font-semibold text-text mb-4">Logging</h2>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <p class="text-sm text-text-muted">Log Level</p>
                <p class="text-lg font-medium text-text">{{ settings.log_level }}</p>
              </div>
              <div>
                <p class="text-sm text-text-muted">Log Format</p>
                <p class="text-lg font-medium text-text">{{ settings.log_format }}</p>
              </div>
            </div>
          </div>

          <div class="bg-surface border border-border rounded-lg p-6">
            <h2 class="text-lg font-semibold text-text mb-4">Security & Limits</h2>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <p class="text-sm text-text-muted">JWT Expiry (hours)</p>
                <p class="text-lg font-medium text-text">{{ settings.jwt_expiry }}</p>
              </div>
              <div>
                <p class="text-sm text-text-muted">Max File Upload</p>
                <p class="text-lg font-medium text-text">
                  {{ (settings.max_file_upload_size / 1024 / 1024).toFixed(2) }} MB
                </p>
              </div>
              <div>
                <p class="text-sm text-text-muted">Rate Limit</p>
                <p class="text-lg font-medium text-text">
                  {{ settings.rate_limit_per_min }} req/min
                </p>
              </div>
              <div>
                <p class="text-sm text-text-muted">CORS Origins</p>
                <p class="text-lg font-medium text-text">{{ settings.cors_origins }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </AppLayout>
</template>
