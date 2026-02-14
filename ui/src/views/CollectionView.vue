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
import { FolderOpen, Filter, Plus, MoreHorizontal, Edit, Trash2, Settings } from 'lucide-vue-next';

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
    const col = response.data.data.find((c: Collection) => c.name === collectionName.value);
    collection.value = col;

    // Initialize all fields as visible
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
        <div class="flex items-center text-sm text-text-muted truncate gap-2">
          <span class="hover:text-text cursor-pointer font-medium text-text" @click="router.push('/collections')">Collections</span>
          <span class="text-text-muted flex-shrink-0">/</span>
          <span class="font-medium text-text truncate">{{ collectionName }}</span>
        </div>
      </template>
    </AppHeader>

    <!-- Main Scrollable Area -->
    <div class="flex-1 overflow-auto min-h-0 p-4 sm:p-8 pb-24 sm:pb-8">
      <div class="max-w-7xl mx-auto space-y-6 sm:space-y-8">
        <!-- Page Title -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h1 class="text-2xl font-bold text-text tracking-tight">{{ collectionName }}</h1>
            <p class="mt-1 text-sm text-text-muted">
              {{ collection?.type }} collection • {{ collection?.fields?.length || 0 }} fields
            </p>
          </div>
          <div class="flex items-center gap-3">
            <Button
              variant="secondary"
              size="sm"
              class="flex-1 sm:flex-none"
              @click="router.push(`/collections/${collectionName}/settings`)"
            >
              <Settings class="w-4 h-4" />
              <span class="hidden sm:inline">Settings</span>
            </Button>
            <Popover align="right">
              <template #trigger>
                <Button variant="secondary" size="sm" class="flex-1 sm:flex-none">
                  <Filter class="w-4 h-4" />
                  <span class="hidden sm:inline">Filter</span>
                </Button>
              </template>
              <template #default>
                <div class="p-3 min-w-[200px]">
                  <div
                    class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-2 px-1"
                  >
                    Visible Fields
                  </div>
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
              class="flex-1 sm:flex-none"
              @click="router.push(`/collections/${collectionName}/new`)"
            >
              <Plus class="w-4 h-4" />
              <span class="whitespace-nowrap">New Record</span>
            </Button>
          </div>
        </div>

        <!-- Data Table -->
        <Table
          :headers="[
            ...filteredFields.map((f) => ({ key: f.name, label: f.name })),
            { key: 'actions', label: 'Actions', align: 'center', sticky: true },
          ]"
          :items="records"
          :enable-pagination="true"
          :default-page-size="15"
        >
          <template
            v-for="field in filteredFields"
            :key="field.name"
            #[`cell(${field.name})`]="{ item }"
          >
            <span v-if="field.type === 'bool'" class="text-text">
              {{
                (item.data as any)?.[field.name] === 1 || (item.data as any)?.[field.name] === true
                  ? 'true'
                  : 'false'
              }}
            </span>
            <span v-else class="text-text">{{ (item.data as any)?.[field.name] ?? '-' }}</span>
          </template>

          <template #cell(actions)="{ item }">
            <Popover align="right">
              <template #trigger>
                <Button variant="ghost" size="xs">
                  <MoreHorizontal class="w-4 h-4" />
                </Button>
              </template>
              <template #default="{ close }">
                <PopoverItem
                  :icon="Edit"
                  @click="
                    close();
                    router.push(`/collections/${collectionName}/edit/${item.id as string}`);
                  "
                >
                  Edit
                </PopoverItem>
                <PopoverItem
                  :icon="Trash2"
                  variant="danger"
                  @click="
                    close();
                    confirmDelete(item.id as string);
                  "
                >
                  Delete
                </PopoverItem>
              </template>
            </Popover>
          </template>

          <template #empty>
            <div class="py-12 text-center text-text-muted">
              <FolderOpen class="w-12 h-12 mx-auto mb-3 opacity-30" />
              <p class="text-sm mb-4">No records found</p>
              <Button variant="link" @click="router.push(`/collections/${collectionName}/new`)"
                >Create your first record</Button
              >
            </div>
          </template>

        </Table>
      </div>
    </div>
  </AppLayout>
</template>
