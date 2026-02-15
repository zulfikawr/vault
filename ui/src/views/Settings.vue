<script setup lang="ts">
import { ref, onMounted } from 'vue';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import Input from '../components/Input.vue';
import Dropdown from '../components/Dropdown.vue';
import DropdownItem from '../components/DropdownItem.vue';
import { Save } from 'lucide-vue-next';

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
const saving = ref(false);

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
    alert('Settings saved successfully');
  } catch (error) {
    console.error('Failed to save settings', error);
    alert('Failed to save settings');
  } finally {
    saving.value = false;
  }
};

onMounted(fetchSettings);
</script>

<template>
  <AppLayout>
    <AppHeader>
      <template #breadcrumb>
        <div class="flex items-center text-sm text-text-muted">
          <span class="font-medium text-text">Settings</span>
        </div>
      </template>
    </AppHeader>

    <main class="flex-1 overflow-auto min-h-0 p-4 sm:p-8 pb-24 sm:pb-8">
      <div class="max-w-7xl mx-auto space-y-6 sm:space-y-8">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h1 class="text-xl font-semibold text-text">System Settings</h1>
            <p class="mt-1 text-sm text-text-muted">Configure your Vault instance.</p>
          </div>
          <Button
            v-if="settings"
            variant="primary"
            size="sm"
            class="px-3 py-1.5 text-sm"
            :disabled="saving"
            @click="saveSettings"
          >
            <Save class="w-4 h-4" />
            Save
          </Button>
        </div>

        <div v-if="loading" class="text-center text-text-muted">Loading settings...</div>

        <div v-else-if="settings" class="space-y-4">
          <div class="bg-surface border border-border rounded-lg p-6">
            <h2 class="text-lg font-semibold text-text mb-4">Server Configuration</h2>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm text-text-muted mb-2">Port</label>
                <Input
                  :model-value="String(settings!.port)"
                  type="number"
                  size="sm"
                  @update:model-value="settings!.port = Number($event)"
                />
              </div>
              <div>
                <label class="block text-sm text-text-muted mb-2">TLS Enabled</label>
                <Dropdown
                  :model-value="settings!.tls_enabled ? 'true' : 'false'"
                  size="sm"
                  @update:model-value="settings!.tls_enabled = $event === 'true'"
                >
                  <template #trigger>
                    {{ settings!.tls_enabled ? 'Yes' : 'No' }}
                  </template>
                  <template #default="{ close }">
                    <DropdownItem
                      value="true"
                      @click="
                        () => {
                          settings!.tls_enabled = true;
                          close();
                        }
                      "
                    >
                      Yes
                    </DropdownItem>
                    <DropdownItem
                      value="false"
                      @click="
                        () => {
                          settings!.tls_enabled = false;
                          close();
                        }
                      "
                    >
                      No
                    </DropdownItem>
                  </template>
                </Dropdown>
              </div>
            </div>
          </div>

          <div class="bg-surface border border-border rounded-lg p-6">
            <h2 class="text-lg font-semibold text-text mb-4">Logging</h2>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm text-text-muted mb-2">Log Level</label>
                <Dropdown v-model="settings!.log_level" size="sm">
                  <template #trigger>
                    {{ settings!.log_level }}
                  </template>
                  <template #default="{ close }">
                    <DropdownItem
                      value="DEBUG"
                      @click="
                        () => {
                          settings!.log_level = 'DEBUG';
                          close();
                        }
                      "
                    >
                      DEBUG
                    </DropdownItem>
                    <DropdownItem
                      value="INFO"
                      @click="
                        () => {
                          settings!.log_level = 'INFO';
                          close();
                        }
                      "
                    >
                      INFO
                    </DropdownItem>
                    <DropdownItem
                      value="WARN"
                      @click="
                        () => {
                          settings!.log_level = 'WARN';
                          close();
                        }
                      "
                    >
                      WARN
                    </DropdownItem>
                    <DropdownItem
                      value="ERROR"
                      @click="
                        () => {
                          settings!.log_level = 'ERROR';
                          close();
                        }
                      "
                    >
                      ERROR
                    </DropdownItem>
                  </template>
                </Dropdown>
              </div>
              <div>
                <label class="block text-sm text-text-muted mb-2">Log Format</label>
                <Dropdown v-model="settings!.log_format" size="sm">
                  <template #trigger>
                    {{ settings!.log_format }}
                  </template>
                  <template #default="{ close }">
                    <DropdownItem
                      value="text"
                      @click="
                        () => {
                          settings!.log_format = 'text';
                          close();
                        }
                      "
                    >
                      text
                    </DropdownItem>
                    <DropdownItem
                      value="json"
                      @click="
                        () => {
                          settings!.log_format = 'json';
                          close();
                        }
                      "
                    >
                      json
                    </DropdownItem>
                  </template>
                </Dropdown>
              </div>
            </div>
          </div>

          <div class="bg-surface border border-border rounded-lg p-6">
            <h2 class="text-lg font-semibold text-text mb-4">Security & Limits</h2>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm text-text-muted mb-2">JWT Expiry (hours)</label>
                <Input
                  :model-value="String(settings!.jwt_expiry)"
                  type="number"
                  size="sm"
                  @update:model-value="settings!.jwt_expiry = Number($event)"
                />
              </div>
              <div>
                <label class="block text-sm text-text-muted mb-2">Max File Upload (MB)</label>
                <Input
                  :model-value="String(settings!.max_file_upload_size / 1024 / 1024)"
                  type="number"
                  size="sm"
                  @update:model-value="
                    settings!.max_file_upload_size = Number($event) * 1024 * 1024
                  "
                />
              </div>
              <div>
                <label class="block text-sm text-text-muted mb-2">Rate Limit (req/min)</label>
                <Input
                  :model-value="String(settings!.rate_limit_per_min)"
                  type="number"
                  size="sm"
                  @update:model-value="settings!.rate_limit_per_min = Number($event)"
                />
              </div>
              <div>
                <label class="block text-sm text-text-muted mb-2">CORS Origins</label>
                <Input v-model="settings!.cors_origins" type="text" size="sm" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>
  </AppLayout>
</template>
