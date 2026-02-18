<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import Table from '../components/Table.vue';
import ConfirmModal from '../components/ConfirmModal.vue';
import Popover from '../components/Popover.vue';
import PopoverItem from '../components/PopoverItem.vue';
import Checkbox from '../components/Checkbox.vue';
import { FolderOpen, Filter, Plus, MoreHorizontal, Edit, Trash2, Settings, Database, ChevronRight, Activity, Search, RefreshCw } from 'lucide-vue-next';

interface Field {
  name: string;
  type: string;
}

interface Collection {
  name: string;
  type: string;
  fields: Field[];
}

interface RecordData {
  id: string;
  data?: Record<string, string | number | boolean>;
  [key: string]: string | number | boolean | Record<string, string | number | boolean> | undefined;
}

const router = useRouter();
const route = useRoute();
const collections = ref<Collection[]>([]);
const collection = ref<Collection | null>(null);
const records = ref<RecordData[]>([]);
const showDeleteModal = ref(false);
const recordToDelete = ref('');
const selectedRecords = ref<string[]>([]);
const visibleFields = ref<Record<string, boolean>>({});
const loading = ref(false);

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
    const col = response.data.data.find((c: Collection) => c.name === collectionName.value);
    collection.value = col;
    if (col?.fields) {
      visibleFields.value = col.fields.reduce((acc: Record<string, boolean>, field: Field) => {
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
  return collection.value.fields.filter((f: Field) => visibleFields.value[f.name]);
});

const fetchRecords = async () => {
  loading.value = true;
  try {
    const response = await axios.get(`/api/collections/${collectionName.value}/records`);
    records.value = response.data.items || [];
  } catch (error) {
    console.error('Failed to fetch records', error);
  } finally {
    loading.value = false;
  }
};

const confirmDelete = (id: string) => {
  recordToDelete.value = id;
  showDeleteModal.value = true;
};

const deleteRecord = async () => {
  try {
    if (selectedRecords.value.length > 0) {
      await axios.delete(`/api/collections/${collectionName.value}/records`, {
        data: { ids: selectedRecords.value },
      });
    } else {
      await axios.delete(`/api/collections/${collectionName.value}/records/${recordToDelete.value}`);
    }
    showDeleteModal.value = false;
    recordToDelete.value = '';
    selectedRecords.value = [];
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

    <AppHeader>
      <template #breadcrumb>
        <div class="flex items-center gap-2 text-xs font-bold tracking-wide">
          <Database class="w-4 h-4 text-primary" />
          <span class="text-text-dim hover:text-primary transition-colors cursor-pointer" @click="router.push('/collections')">Collections</span>
          <ChevronRight class="w-3 h-3 text-text-dim" />
          <span class="text-primary truncate">{{ collectionName }}</span>
        </div>
      </template>
    </AppHeader>

    <main class="flex-1 min-h-0 p-4 lg:p-8 max-w-7xl mx-auto w-full">
      <div class="space-y-8">
        <!-- Header -->
        <div class="flex flex-col md:flex-row md:items-end justify-between gap-6 border-b border-border/30 pb-6">
          <div>
            <h1 class="text-3xl font-bold text-text tracking-tight">{{ collectionName }}</h1>
            <p class="text-text-dim text-xs mt-1">
              {{ collection?.type }} collection • {{ collection?.fields?.length || 0 }} fields
            </p>
          </div>
          
          <div class="flex items-center gap-3">
            <Button
              v-if="selectedRecords.length > 0"
              variant="destructive"
              size="sm"
              @click="showDeleteModal = true"
            >
              <Trash2 class="w-4 h-4 mr-2" />
              Delete Selected ({{ selectedRecords.length }})
            </Button>
            <Button
              variant="secondary"
              size="sm"
              @click="router.push(`/collections/${collectionName}/settings`)"
            >
              <Settings class="w-4 h-4 mr-2" />
              Settings
            </Button>
            <Popover align="right">
              <template #trigger>
                <Button variant="secondary" size="sm">
                  <Filter class="w-4 h-4 mr-2" />
                  Filter
                </Button>
              </template>
              <template #default>
                <div class="p-4 min-w-[240px] space-y-4">
                  <div class="text-[10px] font-bold text-text-dim tracking-wider">Visible Fields</div>
                  <div class="space-y-2">
                    <Checkbox
                      v-for="field in collection?.fields"
                      :key="field.name"
                      :model-value="visibleFields[field.name] || false"
                      :label="field.name"
                      @update:model-value="visibleFields[field.name] = $event"
                    />
                  </div>
                </div>
              </template>
            </Popover>
            <Button
              size="sm"
              @click="router.push(`/collections/${collectionName}/new`)"
            >
              <Plus class="w-4 h-4 mr-2" />
              New Record
            </Button>
          </div>
        </div>

        <!-- Data Table -->
        <div class="animate-in fade-in slide-in-from-bottom-4 duration-500">
          <Table
            :headers="[
              ...filteredFields.map((f) => ({ key: f.name, label: f.name })),
              { key: 'actions', label: 'Actions', align: 'center', sticky: true },
            ]"
            :items="records"
            :loading="loading"
            :enable-pagination="true"
            :default-page-size="15"
            selectable
            selection-key="id"
            @selection-change="selectedRecords = $event"
          >
            <template
              v-for="field in filteredFields"
              :key="field.name"
              #[`cell(${field.name})`]="{ item }"
            >
              <span v-if="field.type === 'bool'" class="text-xs font-medium text-text-muted">
                {{ (item.data as any)?.[field.name] ? 'true' : 'false' }}
              </span>
              <span v-else class="text-xs font-medium text-text-muted">{{ (item.data as any)?.[field.name] ?? '-' }}</span>
            </template>

            <template #cell(actions)="{ item }">
              <div class="actions-cell">
                <Popover align="right">
                  <template #trigger>
                    <Button variant="ghost" size="xs">
                      <MoreHorizontal class="w-4 h-4" />
                    </Button>
                  </template>
                  <template #default="{ close }">
                    <PopoverItem :icon="Edit" @click="close(); router.push(`/collections/${collectionName}/edit/${item.id as string}`);">Edit</PopoverItem>
                    <PopoverItem :icon="Trash2" variant="danger" @click="close(); confirmDelete(item.id as string);">Delete</PopoverItem>
                  </template>
                </Popover>
              </div>
            </template>

            <template #empty>
              <div class="py-16 text-center">
                <Activity class="w-12 h-12 mx-auto mb-4 opacity-10" />
                <p class="text-xs font-bold tracking-widest text-text-dim">No records found</p>
                <Button variant="link" class="mt-4 text-[10px] font-bold tracking-widest" @click="router.push(`/collections/${collectionName}/new`)">Create your first record</Button>
              </div>
            </template>
          </Table>
        </div>
      </div>
    </main>
  </AppLayout>
</template>
