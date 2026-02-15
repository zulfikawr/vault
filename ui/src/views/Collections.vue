<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import Table from '../components/Table.vue';
import Popover from '../components/Popover.vue';
import PopoverItem from '../components/PopoverItem.vue';
import Checkbox from '../components/Checkbox.vue';
import ConfirmModal from '../components/ConfirmModal.vue';
import { FolderOpen, Filter, Plus, MoreHorizontal, Settings, Trash2 } from 'lucide-vue-next';

interface Field {
  name: string;
  type: string;
  required?: boolean;
}

interface Collection {
  name: string;
  type: 'base' | 'auth' | 'system';
  fields: Field[];
  created: string;
  recordCount?: number;
}

const router = useRouter();
const collections = ref<Collection[]>([]);
const recordCounts = ref<Record<string, number>>({});
const showSystemCollections = ref(false);
const filterTypes = ref({
  base: true,
  auth: true,
  system: true,
});
const showDeleteModal = ref(false);
const collectionToDelete = ref<Collection | null>(null);

// Sorting state
const sortKey = ref('name');
const sortOrder = ref<'asc' | 'desc'>('asc');

const filteredCollections = computed(() => {
  let result = collections.value.filter((col: Collection) => {
    // Filter by system collections
    if (!showSystemCollections.value && col.name.startsWith('_')) {
      return false;
    }
    // Filter by type
    return filterTypes.value[col.type];
  });

  // Apply sorting
  if (sortKey.value) {
    result.sort((a, b) => {
      let aValue, bValue;

      if (sortKey.value === 'records') {
        // Special handling for record count
        aValue = recordCounts.value[a.name] ?? 0;
        bValue = recordCounts.value[b.name] ?? 0;
      } else {
        const key = sortKey.value as keyof Collection;
        aValue = a[key];
        bValue = b[key];
      }

      // Handle null/undefined values
      if (aValue == null && bValue == null) return 0;
      if (aValue == null) return sortOrder.value === 'asc' ? 1 : -1;
      if (bValue == null) return sortOrder.value === 'asc' ? -1 : 1;

      // Handle dates
      if (
        typeof aValue === 'string' &&
        !isNaN(Date.parse(aValue)) &&
        typeof bValue === 'string' &&
        !isNaN(Date.parse(bValue))
      ) {
        const dateA = new Date(aValue).getTime();
        const dateB = new Date(bValue).getTime();
        return sortOrder.value === 'asc' ? dateA - dateB : dateB - dateA;
      }

      // Handle numbers
      if (typeof aValue === 'number' && typeof bValue === 'number') {
        return sortOrder.value === 'asc' ? aValue - bValue : bValue - aValue;
      }

      // Handle strings
      if (typeof aValue === 'string' && typeof bValue === 'string') {
        return sortOrder.value === 'asc'
          ? aValue.localeCompare(bValue)
          : bValue.localeCompare(aValue);
      }

      // Fallback comparison
      const strA = String(aValue);
      const strB = String(bValue);
      return sortOrder.value === 'asc' ? strA.localeCompare(strB) : strB.localeCompare(strA);
    });
  }

  return result;
});

const fetchCollections = async () => {
  try {
    const response = await axios.get('/api/admin/collections');
    collections.value = response.data.data;

    // Fetch record counts for each collection
    for (const collection of collections.value) {
      try {
        const recordResponse = await axios.get(
          `/api/collections/${collection.name}/records?perPage=1`
        );
        recordCounts.value[collection.name] = recordResponse.data.totalItems || 0;
      } catch (error) {
        console.error(`Failed to fetch record count for collection ${collection.name}:`, error);
        recordCounts.value[collection.name] = 0;
      }
    }
  } catch (error) {
    console.error('Failed to fetch collections', error);
  }
};

const handleDeleteClick = (collection: Record<string, unknown>) => {
  collectionToDelete.value = collection as unknown as Collection;
  showDeleteModal.value = true;
};

const confirmDelete = async () => {
  if (!collectionToDelete.value) return;

  try {
    await axios.delete(`/api/admin/collections/${collectionToDelete.value.name}`);
    // Refresh the collections list
    await fetchCollections();
    showDeleteModal.value = false;
    collectionToDelete.value = null;
  } catch (error) {
    console.error('Failed to delete collection', error);
    // Optionally show an error message to the user
  }
};

onMounted(() => {
  fetchCollections();
});
</script>

<template>
  <AppLayout>
    <!-- Header -->
    <AppHeader>
      <template #breadcrumb>
        <div class="flex items-center text-sm text-text-muted truncate gap-2">
          <span
            class="hover:text-text cursor-pointer font-medium text-text"
            @click="router.push('/')"
            >Collections</span
          >
        </div>
      </template>
    </AppHeader>

    <!-- Main Scrollable Area -->
    <div class="flex-1 overflow-auto min-h-0 p-4 sm:p-8 pb-24 sm:pb-8">
      <div class="max-w-7xl mx-auto space-y-6 sm:space-y-8">
        <!-- Page Title -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h1 class="text-2xl font-bold text-text tracking-tight">Data Collections</h1>
            <p class="mt-1 text-sm text-text-muted">
              Manage your database schemas and content types.
            </p>
          </div>
          <div class="flex items-center gap-3">
            <Popover align="right">
              <template #trigger>
                <Button variant="secondary" size="sm" class="flex-1 sm:flex-none">
                  <Filter class="w-4 h-4" />
                  Filter
                </Button>
              </template>
              <template #default>
                <div class="p-3 min-w-[200px]">
                  <div
                    class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-2 px-1"
                  >
                    Collection Type
                  </div>
                  <div class="space-y-2 mb-4">
                    <Checkbox v-model="filterTypes.base" label="Base" />
                    <Checkbox v-model="filterTypes.auth" label="Auth" />
                    <Checkbox v-model="filterTypes.system" label="System" />
                  </div>
                  <div class="border-t border-border pt-3">
                    <Checkbox v-model="showSystemCollections" label="Show system collections" />
                  </div>
                </div>
              </template>
            </Popover>
            <Button size="sm" class="flex-1 sm:flex-none" @click="router.push('/collections/new')">
              <Plus class="w-4 h-4" />
              New Collection
            </Button>
          </div>
        </div>

        <!-- Data Table -->
        <Table
          :headers="[
            { key: 'name', label: 'Name', sortable: true },
            { key: 'type', label: 'Type', sortable: true },
            { key: 'fields', label: 'Fields', sortable: true },
            { key: 'records', label: 'Records', sortable: true },
            { key: 'created', label: 'Created', sortable: true },
            { key: 'actions', label: 'Actions', align: 'center', sticky: true },
          ]"
          :items="filteredCollections"
          :enable-pagination="true"
          :default-page-size="10"
          :sort-key="sortKey"
          :sort-order="sortOrder"
          row-clickable
          @sort-change="
            (key, order) => {
              sortKey = key;
              sortOrder = order;
            }
          "
          @row-click="
            (item, event) => {
              // Only navigate if the click didn't originate from the actions column
              const target = event.target as HTMLElement;
              if (!target.closest('.actions-cell')) {
                router.push(`/collections/${(item as unknown as Collection).name}`);
              }
            }
          "
        >
          <template #cell(name)="{ item }">
            <div class="flex items-center gap-3">
              <div class="p-1.5 rounded bg-primary/10 text-primary">
                <FolderOpen class="w-4 h-4" />
              </div>
              <span class="font-medium text-text">{{ item.name }}</span>
            </div>
          </template>

          <template #cell(type)="{ item }">
            <span class="text-text-muted">{{ item.type }}</span>
          </template>

          <template #cell(fields)="{ item }">
            <span class="text-text-muted"
              >{{ (item as unknown as Collection).fields?.length ?? 0 }} fields</span
            >
          </template>

          <template #cell(records)="{ item }">
            <span class="text-text-muted"
              >{{ recordCounts[(item as unknown as Collection).name] ?? 0 }} records</span
            >
          </template>

          <template #cell(created)="{ item }">
            <span class="text-text-muted text-xs">{{
              item.created ? new Date(item.created as string).toLocaleDateString() : '-'
            }}</span>
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
                  <PopoverItem
                    :icon="Settings"
                    @click="
                      close();
                      router.push(`/collections/${item.name}/settings`);
                    "
                  >
                    Settings
                  </PopoverItem>
                  <PopoverItem
                    :icon="Trash2"
                    variant="danger"
                    @click="
                      close();
                      handleDeleteClick(item);
                    "
                  >
                    Delete
                  </PopoverItem>
                </template>
              </Popover>
            </div>
          </template>

          <template #empty>
            <div class="py-12 text-center text-text-muted">
              <FolderOpen class="w-12 h-12 mx-auto mb-3 opacity-30" />
              <p class="text-sm mb-4">No collections found</p>
              <Button variant="link" @click="router.push('/collections/new')"
                >Create your first collection</Button
              >
            </div>
          </template>
        </Table>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <ConfirmModal
      :show="showDeleteModal"
      title="Confirm Delete Collection"
      :message="`Are you sure you want to delete the collection '${collectionToDelete?.name}'? This action cannot be undone and will permanently remove all data in this collection.`"
      confirm-text="Delete"
      cancel-text="Cancel"
      variant="danger"
      @confirm="confirmDelete"
      @cancel="showDeleteModal = false"
    />
  </AppLayout>
</template>
