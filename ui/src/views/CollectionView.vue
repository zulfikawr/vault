<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import { 
  FolderOpen, 
  Terminal, 
  Settings, 
  Cloud
  Search,
  Bell,
  Filter,
  Plus,
  MoreHorizontal
  Edit,
  Trash2
} from 'lucide-vue-next';

const router = useRouter();
const route = useRoute();
const collections = ref([]);
const collection = ref(null);
const records = ref([]);

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
  } catch (error) {
    console.error('Failed to fetch collection', error);
  }
};

const fetchRecords = async () => {
  try {
    const response = await axios.get(`/api/collections/${collectionName.value}/records`);
    records.value = response.data.items || [];
  } catch (error) {
    console.error('Failed to fetch records', error);
  }
};

const deleteRecord = async (id: string) => {
  if (!confirm('Are you sure you want to delete this record?')) return;
  try {
    await axios.delete(`/api/collections/${collectionName.value}/records/${id}`);
    fetchRecords();
  } catch (error) {
    console.error('Delete failed', error);
  }
};

const handleLogout = () => {
  auth.logout();
  router.push({ name: 'Login' });
};

onMounted(() => {
  fetchCollections();
  fetchCollection();
  fetchRecords();
});
</script>

<template>
  <AppLayout>
    
      <!-- Header -->
      <AppHeader>
      <template #breadcrumb>
        <div class="flex items-center text-sm text-text-muted">
          <span class="hover:text-text cursor-pointer" @click="router.push('/')">Vault</span>
          <span class="mx-2">/</span>
          <span class="hover:text-text cursor-pointer" @click="router.push('/collections')">Collections</span>
          <span class="mx-2">/</span>
          <span class="font-medium text-text">{{ collectionName }}</span>
        </div>
      </template>
    </AppHeader>

      <!-- Main Scrollable Area -->
      <div class="flex-1 overflow-auto p-8">
        <div class="max-w-7xl mx-auto space-y-8">
          <!-- Page Title -->
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div>
              <h1 class="text-2xl font-bold text-text tracking-tight">{{ collectionName }}</h1>
              <p class="mt-1 text-sm text-text-muted">{{ collection?.type }} collection • {{ collection?.fields?.length || 0 }} fields</p>
            </div>
            <div class="flex items-center gap-3">
              <button @click="router.push(`/collections/${collectionName}/settings`)" class="px-4 py-2 bg-surface-dark border border-border rounded text-sm font-medium text-text hover:bg-surface transition-colors flex items-center gap-2">
                <Settings class="w-4 h-4" />
                Settings
              </button>
              <button class="px-4 py-2 bg-surface-dark border border-border rounded text-sm font-medium text-text hover:bg-surface transition-colors flex items-center gap-2">
                <Filter class="w-4 h-4" />
                Filter
              </button>
              <button @click="router.push(`/collections/${collectionName}/new`)" class="px-4 py-2 bg-primary hover:bg-primary-hover text-white rounded text-sm font-medium shadow-sm hover:shadow transition-all flex items-center gap-2">
                <Plus class="w-4 h-4" />
                New Record
              </button>
            </div>
          </div>

          <!-- Data Table -->
          <div class="bg-surface-dark rounded-lg border border-border shadow-sm overflow-hidden">
            <div class="overflow-x-auto">
              <table class="w-full text-left text-sm whitespace-nowrap">
                <thead class="bg-surface border-b border-border">
                  <tr>
                    <th v-for="field in collection?.fields" :key="field.name" class="px-6 py-3 text-xs font-medium text-text-muted uppercase tracking-wider">
                      {{ field.name }}
                    </th>
                    <th class="px-6 py-3 text-right text-xs font-medium text-text-muted uppercase tracking-wider">Actions</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-border">
                  <tr v-for="record in records" :key="record.id" class="hover:bg-background/50 transition-colors group">
                    <td v-for="field in collection?.fields" :key="field.name" class="px-6 py-4">
                      <span class="text-text">{{ record[field.name] || '-' }}</span>
                    </td>
                    <td class="px-6 py-4 text-right">
                      <div class="flex justify-end items-center gap-2">
                        <button class="p-2 text-text-muted hover:text-primary hover:bg-primary/10 rounded transition-colors">
                          <Edit class="w-4 h-4" />
                        </button>
                        <button @click="deleteRecord(record.id)" class="p-2 text-text-muted hover:text-error hover:bg-error/10 rounded transition-colors">
                          <Trash2 class="w-4 h-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                  <tr v-if="records.length === 0">
                    <td :colspan="(collection?.fields?.length || 0) + 1" class="px-6 py-12 text-center text-text-muted">
                      <FolderOpen class="w-12 h-12 mx-auto mb-3 opacity-30" />
                      <p class="text-sm mb-4">No records found</p>
                      <button class="text-sm text-primary hover:text-primary-hover">Create your first record</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            
            <!-- Pagination -->
            <div class="bg-surface px-6 py-3 border-t border-border flex items-center justify-between">
              <div class="text-xs text-text-muted">
                Showing <span class="font-medium text-text">{{ records.length }}</span> of <span class="font-medium text-text">{{ records.length }}</span> results
              </div>
              <div class="flex gap-2">
                <button class="px-3 py-1 text-xs font-medium rounded border border-border bg-surface-dark text-text hover:bg-surface transition-colors disabled:opacity-50" disabled>Previous</button>
                <button class="px-3 py-1 text-xs font-medium rounded border border-border bg-surface-dark text-text hover:bg-surface transition-colors">Next</button>
              </div>
            </div>
          </div>
        </div>
      </div>
      </AppLayout>
</template>


