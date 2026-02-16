<script setup lang="ts">
import { ref, computed } from 'vue';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Table from '../components/Table.vue';
import Button from '../components/Button.vue';
import {
  Play,
  Database,
  FileJson,
  Table as TableIcon,
  AlertCircle,
  Copy,
  Check,
} from 'lucide-vue-next';

interface QueryResult {
  columns: string[];
  rows: Record<string, unknown>[];
}

const query = ref('SELECT * FROM users LIMIT 10;');
const loading = ref(false);
const result = ref<QueryResult | null>(null);
const error = ref<string | null>(null);
const activeTab = ref<'table' | 'json'>('table');
const copied = ref(false);

const executeQuery = async () => {
  if (!query.value.trim()) return;

  loading.value = true;
  error.value = null;
  result.value = null;

  try {
    const response = await axios.post('/api/admin/query', { query: query.value });
    result.value = response.data.data;
  } catch (err: unknown) {
    if (axios.isAxiosError(err)) {
      error.value =
        err.response?.data?.message || err.message || 'An error occurred while executing the query';
    } else {
      error.value = (err as Error).message || 'An error occurred while executing the query';
    }
    console.error('Query execution failed', err);
  } finally {
    loading.value = false;
  }
};

const tableHeaders = computed(() => {
  if (!result.value) return [];
  return result.value.columns.map((col) => ({
    key: col,
    label: col,
    sortable: true,
  }));
});

const jsonResult = computed(() => {
  if (!result.value) return '';
  return JSON.stringify(result.value.rows, null, 2);
});

const copyJson = () => {
  navigator.clipboard.writeText(jsonResult.value);
  copied.value = true;
  setTimeout(() => (copied.value = false), 2000);
};
</script>

<template>
  <AppLayout>
    <!-- Header -->
    <AppHeader>
      <template #breadcrumb>
        <div class="flex items-center text-sm text-text-muted truncate gap-2">
          <span class="font-medium text-primary">SQL Editor</span>
        </div>
      </template>
    </AppHeader>

    <div class="flex flex-col h-full bg-background overflow-hidden">
      <!-- Editor & Results (2 Pane) -->
      <div class="flex-1 flex flex-col min-h-0 overflow-hidden">
        <!-- Top Pane: Editor -->
        <div class="h-1/2 border-b border-border flex flex-col min-h-0 relative">
          <div class="absolute top-4 right-4 z-20">
            <Button
              variant="primary"
              size="sm"
              :loading="loading"
              class="gap-2 shadow-lg"
              @click="executeQuery"
            >
              <Play v-if="!loading" class="w-4 h-4" />
              Execute
            </Button>
          </div>
          <div class="flex-1 relative bg-surface-dark">
            <textarea
              v-model="query"
              class="w-full h-full p-6 pr-32 font-mono text-sm bg-transparent border-none focus:ring-0 resize-none text-text-muted focus:text-text transition-colors outline-none"
              placeholder="Enter your SQL query here..."
              spellcheck="false"
            ></textarea>
          </div>
        </div>

        <!-- Bottom Pane: Results -->
        <div class="h-1/2 flex flex-col min-h-0 bg-background">
          <!-- Result Tabs -->
          <div
            class="h-12 border-b border-border px-6 flex items-center justify-between bg-surface shrink-0"
          >
            <div class="flex gap-1 h-full">
              <button
                :class="[
                  'px-4 h-full text-sm font-medium border-b-2 transition-all',
                  activeTab === 'table'
                    ? 'border-primary text-primary'
                    : 'border-transparent text-text-muted hover:text-text',
                ]"
                @click="activeTab = 'table'"
              >
                <div class="flex items-center gap-2">
                  <TableIcon class="w-4 h-4" />
                  Table
                </div>
              </button>
              <button
                :class="[
                  'px-4 h-full text-sm font-medium border-b-2 transition-all',
                  activeTab === 'json'
                    ? 'border-primary text-primary'
                    : 'border-transparent text-text-muted hover:text-text',
                ]"
                @click="activeTab = 'json'"
              >
                <div class="flex items-center gap-2">
                  <FileJson class="w-4 h-4" />
                  JSON
                </div>
              </button>
            </div>

            <div v-if="result && activeTab === 'json'" class="flex items-center">
              <Button variant="ghost" size="xs" class="gap-1.5" @click="copyJson">
                <Check v-if="copied" class="w-3.5 h-3.5 text-success" />
                <Copy v-else class="w-3.5 h-3.5" />
                {{ copied ? 'Copied' : 'Copy JSON' }}
              </Button>
            </div>

            <div v-if="result" class="text-xs text-text-muted">
              {{ result.rows.length }} rows returned
            </div>
          </div>

          <!-- Content Area -->
          <div class="flex-1 overflow-auto p-6">
            <div
              v-if="error"
              class="bg-error/10 border border-error/20 rounded-lg p-4 flex gap-3 text-error"
            >
              <AlertCircle class="w-5 h-5 shrink-0" />
              <div>
                <p class="font-bold text-sm">Query Error</p>
                <p class="text-sm opacity-90">{{ error }}</p>
              </div>
            </div>

            <div
              v-else-if="loading"
              class="flex flex-col items-center justify-center h-full text-text-muted gap-4"
            >
              <svg
                class="animate-spin h-8 w-8 text-primary"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  class="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  stroke-width="4"
                ></circle>
                <path
                  class="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
              <p class="text-sm font-medium animate-pulse">Executing query...</p>
            </div>

            <div v-else-if="result" class="h-full">
              <div v-if="activeTab === 'table'" class="h-full">
                <Table
                  :headers="tableHeaders"
                  :items="result.rows"
                  class="!shadow-none !border-none"
                  enable-pagination
                  :default-page-size="50"
                />
              </div>
              <div v-else class="h-full">
                <div class="relative group h-full">
                  <pre
                    class="bg-surface-dark p-6 rounded-xl border border-border text-xs font-mono overflow-auto h-full max-h-[600px] text-text-muted scrollbar-thin scrollbar-thumb-border hover:scrollbar-thumb-text-muted/20"
                  ><code>{{ jsonResult }}</code></pre>
                </div>
              </div>
            </div>

            <div
              v-else
              class="flex flex-col items-center justify-center h-full text-text-muted/40 gap-4 border-2 border-dashed border-border rounded-2xl"
            >
              <Database class="w-12 h-12 stroke-[1]" />
              <div class="text-center">
                <p class="text-base font-medium text-text-muted">No results to display</p>
                <p class="text-sm">Run a query to see the data</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>

<style scoped>
textarea {
  tab-size: 2;
}

/* Custom Scrollbar for better look */
.scrollbar-thin::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

.scrollbar-thin::-webkit-scrollbar-track {
  background: transparent;
}

.scrollbar-thin::-webkit-scrollbar-thumb {
  background: var(--color-border);
  border-radius: 10px;
}

.scrollbar-thin::-webkit-scrollbar-thumb:hover {
  background: var(--color-text-muted);
}
</style>
