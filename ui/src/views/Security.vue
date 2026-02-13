<script setup lang="ts">
import { ref, onMounted } from 'vue';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import { RefreshCw, CheckCircle, AlertCircle } from 'lucide-vue-next';

interface HealthStatus {
  collections: number;
  status: string;
}

const health = ref<HealthStatus | null>(null);
const loading = ref(false);

const securityFeatures = [
  {
    name: 'Admin Authentication',
    description: 'All collection operations require admin credentials',
    status: 'enabled',
  },
  {
    name: 'Audit Logging',
    description: 'All collection changes are logged with admin ID and timestamp',
    status: 'enabled',
  },
  {
    name: 'Rate Limiting',
    description: 'Collection operations limited to 10 requests per minute',
    status: 'enabled',
  },
  {
    name: 'CSRF Protection',
    description: 'State-changing operations require CSRF token validation',
    status: 'enabled',
  },
  {
    name: 'CORS Validation',
    description: 'Cross-origin requests validated against whitelist',
    status: 'enabled',
  },
  {
    name: 'Transaction Support',
    description: 'Collection operations wrapped in database transactions',
    status: 'enabled',
  },
];

const fetchHealth = async () => {
  loading.value = true;
  try {
    const response = await axios.get('/api/health/collections');
    health.value = response.data.data;
  } catch (error) {
    console.error('Failed to fetch health status', error);
  } finally {
    loading.value = false;
  }
};

onMounted(fetchHealth);
</script>

<template>
  <AppLayout>
    <AppHeader>
      <template #breadcrumb>
        <div class="flex items-center text-sm text-text-muted">
          <span class="hover:text-text cursor-pointer" @click="$router.push('/')">Vault</span>
          <span class="mx-2">/</span>
          <span class="font-medium text-text">Security & Health</span>
        </div>
      </template>
    </AppHeader>

    <main class="flex-1 overflow-auto min-h-0 p-4 sm:p-8 pb-24 sm:pb-8">
      <div class="max-w-7xl mx-auto space-y-6 sm:space-y-8">
        <!-- Header -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h1 class="text-2xl font-bold text-text tracking-tight">Security & Health</h1>
            <p class="mt-1 text-sm text-text-muted">Monitor system health and security features.</p>
          </div>
          <Button variant="secondary" size="sm" :disabled="loading" @click="fetchHealth">
            <RefreshCw class="w-4 h-4" />
            Refresh
          </Button>
        </div>

        <!-- Health Status -->
        <div class="bg-bg-secondary border border-border rounded-lg p-6">
          <h2 class="text-lg font-semibold text-text mb-4">System Health</h2>
          <div v-if="health" class="space-y-3">
            <div class="flex items-center justify-between">
              <span class="text-text-muted">Collections</span>
              <span class="font-semibold text-text">{{ health.collections }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-text-muted">Status</span>
              <div class="flex items-center gap-2">
                <CheckCircle class="w-5 h-5 text-success" />
                <span class="font-semibold text-success capitalize">{{ health.status }}</span>
              </div>
            </div>
          </div>
          <div v-else class="text-text-muted">Loading...</div>
        </div>

        <!-- Security Features -->
        <div class="bg-bg-secondary border border-border rounded-lg p-6">
          <h2 class="text-lg font-semibold text-text mb-4">Security Features</h2>
          <div class="space-y-4">
            <div
              v-for="feature in securityFeatures"
              :key="feature.name"
              class="flex items-start gap-4 pb-4 border-b border-border last:border-b-0 last:pb-0"
            >
              <CheckCircle class="w-5 h-5 text-success flex-shrink-0 mt-0.5" />
              <div class="flex-1">
                <h3 class="font-semibold text-text">{{ feature.name }}</h3>
                <p class="text-sm text-text-muted mt-1">{{ feature.description }}</p>
              </div>
              <span class="text-xs font-semibold text-success uppercase">{{ feature.status }}</span>
            </div>
          </div>
        </div>

        <!-- Information -->
        <div class="bg-bg-secondary border border-border rounded-lg p-6">
          <h2 class="text-lg font-semibold text-text mb-4">Information</h2>
          <div class="space-y-3 text-sm text-text-muted">
            <p>
              <strong class="text-text">Audit Logs:</strong> All collection management operations are logged in the
              <router-link to="/audit-logs" class="text-primary hover:underline">Audit Logs</router-link> page.
            </p>
            <p>
              <strong class="text-text">Rate Limiting:</strong> Collection operations are limited to 10 requests per minute per IP
              address.
            </p>
            <p>
              <strong class="text-text">Authentication:</strong> All collection operations require admin authentication via email and
              password.
            </p>
            <p>
              <strong class="text-text">Data Integrity:</strong> Collection operations are wrapped in database transactions to ensure
              consistency.
            </p>
          </div>
        </div>
      </div>
    </main>
  </AppLayout>
</template>
