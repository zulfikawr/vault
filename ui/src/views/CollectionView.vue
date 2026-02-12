<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import ConfirmModal from '../components/ConfirmModal.vue';
import Popover from '../components/Popover.vue';
import PopoverItem from '../components/PopoverItem.vue';
import Checkbox from '../components/Checkbox.vue';
import { 
  FolderOpen, 
  Filter,
  Plus,
  MoreHorizontal,
  Edit,
  Trash2,
  Settings
} from 'lucide-vue-next';

const router = useRouter();
const route = useRoute();
const collections = ref([]);
const collection = ref(null);
const records = ref([]);
const showDeleteModal = ref(false);
const recordToDelete = ref('');
const visibleFields = ref<Record<string, boolean>>({});

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
    
    // Initialize all fields as visible
    if (col?.fields) {
      visibleFields.value = col.fields.reduce((acc: Record<string, boolean>, field: any) => {
        acc[field.name] = true;
        return acc;
      }, {});
    }
  } catch (error) {
    console.error('Failed to fetch collection', error);
  }
};

const filteredFields = computed(() => {
  if (!collection.value?.fields) return [];
  return collection.value.fields.filter((f: any) => visibleFields.value[f.name]);
});

const fetchRecords = async () => {
  try {
    const response = await axios.get(`/api/collections/${collectionName.value}/records`);
    records.value = response.data.items || [];
  } catch (error) {
    console.error('Failed to fetch records', error);
  }
};

const confirmDelete = (id: string) => {
  recordToDelete.value = id;
  showDeleteModal.value = true;
};

const deleteRecord = async () => {
  try {
    await axios.delete(`/api/collections/${collectionName.value}/records/${recordToDelete.value}`);
    showDeleteModal.value = false;
    recordToDelete.value = '';
    fetchRecords();
  } catch (error) {
    console.error('Delete failed', error);
  }
};

onMounted(() => {
  fetchCollections();
  fetchCollection();
  fetchRecords();
});
</script>

<template>
  <AppLayout>
    <ConfirmModal
      :show="showDeleteModal"
      title="Delete Record"
      message="Are you sure you want to delete this record? This action cannot be undone."
      confirm-text="Delete"
      cancel-text="Cancel"
      variant="danger"
      @confirm="deleteRecord"
      @cancel="showDeleteModal = false"
    />
    
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
              <Popover align="right">
                <template #trigger>
                  <button class="px-4 py-2 bg-surface-dark border border-border rounded text-sm font-medium text-text hover:bg-surface transition-colors flex items-center gap-2">
                    <Filter class="w-4 h-4" />
                    Filter
                  </button>
                </template>
                <template #default>
                  <div class="p-3">
                    <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-2 px-1">
                      Visible Fields
                    </div>
                    <div class="space-y-2">
                      <Checkbox
                        v-for="field in collection?.fields"
                        :key="field.name"
                        v-model="visibleFields[field.name]"
                        :label="field.name"
                      />
                    </div>
                  </div>
                </template>
              </Popover>
              <button @click="router.push(`/collections/${collectionName}/new`)" class="px-4 py-2 bg-primary hover:bg-primary-hover text-white rounded text-sm font-medium shadow-sm hover:shadow transition-all flex items-center gap-2">
                <Plus class="w-4 h-4" />
                New Record
              </button>
            </div>
          </div>

          <!-- Data Table -->
          <div class="bg-surface-dark rounded-lg border border-border shadow-sm">
            <div class="overflow-x-auto">
              <table class="w-full text-left text-sm whitespace-nowrap">
                <thead class="bg-surface border-b border-border">
                  <tr>
                    <th v-for="field in filteredFields" :key="field.name" class="px-6 py-3 text-xs font-medium text-text-muted uppercase tracking-wider">
                      {{ field.name }}
                    </th>
                    <th class="sticky right-0 bg-surface px-6 py-3 text-center text-xs font-medium text-text-muted uppercase tracking-wider">Actions</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-border">
                  <tr v-for="record in records" :key="record.id" class="hover:bg-background/50 transition-colors group">
                    <td v-for="field in filteredFields" :key="field.name" class="px-6 py-4">
                      <span v-if="field.type === 'bool'" class="text-text">
                        {{ record.data?.[field.name] === 1 || record.data?.[field.name] === true ? 'true' : 'false' }}
                      </span>
                      <span v-else class="text-text">{{ record.data?.[field.name] ?? '-' }}</span>
                    </td>
                    <td class="sticky right-0 bg-surface-dark px-6 py-4 group-hover:bg-background">
                      <div class="flex justify-center">
                        <Popover align="right">
                          <template #trigger>
                            <button class="p-2 text-text-muted hover:text-text hover:bg-surface-dark rounded transition-colors">
                              <MoreHorizontal class="w-4 h-4" />
                            </button>
                          </template>
                          <template #default="{ close }">
                            <PopoverItem 
                              :icon="Edit" 
                              @click="close(); router.push(`/collections/${collectionName}/edit/${record.id}`)"
                            >
                              Edit
                            </PopoverItem>
                            <PopoverItem 
                              :icon="Trash2" 
                              variant="danger"
                              @click="close(); confirmDelete(record.id)"
                            >
                              Delete
                            </PopoverItem>
                          </template>
                        </Popover>
                      </div>
                    </td>
                  </tr>
                  <tr v-if="records.length === 0">
                    <td :colspan="filteredFields.length + 1" class="px-6 py-12 text-center text-text-muted">
                      <FolderOpen class="w-12 h-12 mx-auto mb-3 opacity-30" />
                      <p class="text-sm mb-4">No records found</p>
                      <button @click="router.push(`/collections/${collectionName}/new`)" class="text-sm text-primary hover:text-primary-hover">Create your first record</button>
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


