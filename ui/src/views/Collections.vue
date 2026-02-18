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
import { 
  FolderOpen, 
  Filter, 
  Plus, 
  MoreHorizontal, 
  Settings, 
  Trash2,
  Database,
  ChevronRight,
  Activity,
  Search,
  LayoutGrid,
  List
} from 'lucide-vue-next';

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
const sortKey = ref('name');
const sortOrder = ref<'asc' | 'desc'>('asc');

const filteredCollections = computed(() => {
  let result = collections.value.filter((col: Collection) => {
    if (!showSystemCollections.value && col.name.startsWith('_')) return false;
    return filterTypes.value[col.type];
  });

  if (sortKey.value) {
    result.sort((a, b) => {
      let aValue, bValue;
      if (sortKey.value === 'records') {
        aValue = recordCounts.value[a.name] ?? 0;
        bValue = recordCounts.value[b.name] ?? 0;
      } else {
        const key = sortKey.value as keyof Collection;
        aValue = a[key];
        bValue = b[key];
      }
      if (aValue == null && bValue == null) return 0;
      if (aValue == null) return sortOrder.value === 'asc' ? 1 : -1;
      if (bValue == null) return sortOrder.value === 'asc' ? -1 : 1;
      if (typeof aValue === 'number' && typeof bValue === 'number') {
        return sortOrder.value === 'asc' ? aValue - bValue : bValue - aValue;
      }
      if (typeof aValue === 'string' && typeof bValue === 'string') {
        return sortOrder.value === 'asc' ? aValue.localeCompare(bValue) : bValue.localeCompare(aValue);
      }
      return sortOrder.value === 'asc' ? String(aValue).localeCompare(String(bValue)) : String(bValue).localeCompare(String(aValue));
    });
  }
  return result;
});

const fetchCollections = async () => {
  try {
    const response = await axios.get('/api/admin/collections');
    collections.value = response.data.data;
    for (const collection of collections.value) {
      try {
        const recordResponse = await axios.get(`/api/collections/${collection.name}/records?perPage=1`);
        recordCounts.value[collection.name] = recordResponse.data.totalItems || 0;
      } catch (error) {
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
    await fetchCollections();
    showDeleteModal.value = false;
    collectionToDelete.value = null;
  } catch (error) {
    console.error('Failed to delete collection', error);
  }
};

onMounted(fetchCollections);
</script>

<template>
  <AppLayout>
    <AppHeader>
      <template #breadcrumb>
        <div class="flex items-center gap-2 text-xs font-bold tracking-wide">
          <Database class="w-4 h-4 text-primary" />
          <span class="text-primary">Collections</span>
        </div>
      </template>
    </AppHeader>

    <main class="flex-1 min-h-0 p-4 lg:p-8 max-w-7xl mx-auto w-full">
      <div class="space-y-8">
        <!-- Header -->
        <div class="flex flex-col md:flex-row md:items-end justify-between gap-6 border-b border-border/30 pb-6">
          <div>
            <h1 class="text-3xl font-bold text-text tracking-tight">Data Collections</h1>
            <p class="text-text-dim text-xs mt-1">Manage database schemas and content structures.</p>
          </div>
          
          <div class="flex gap-3">
            <Popover align="right">
              <template #trigger>
                <Button variant="secondary" size="sm">
                  <Filter class="w-4 h-4 mr-2" />
                  Filter
                </Button>
              </template>
              <template #default>
                <div class="p-4 min-w-[240px] space-y-4">
                  <div class="text-[10px] font-bold text-text-dim tracking-wider">Collection Type</div>
                  <div class="space-y-2">
                    <Checkbox v-model="filterTypes.base" label="Base" />
                    <Checkbox v-model="filterTypes.auth" label="Auth" />
                    <Checkbox v-model="filterTypes.system" label="System" />
                  </div>
                  <div class="border-t border-border/50 pt-3">
                    <Checkbox v-model="showSystemCollections" label="Show system collections" />
                  </div>
                </div>
              </template>
            </Popover>
            <Button size="sm" @click="router.push('/collections/new')">
              <Plus class="w-4 h-4 mr-2" />
              New Collection
            </Button>
          </div>
        </div>

        <!-- Data Table -->
        <div class="animate-in fade-in slide-in-from-bottom-4 duration-500">
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
            @sort-change="(key, order) => { sortKey = key; sortOrder = order; }"
            @row-click="(item, event) => { if (!(event.target as HTMLElement).closest('.actions-cell')) router.push(`/collections/${item.name}`); }"
          >
            <template #cell(name)="{ item }">
              <div class="flex items-center py-1 group/item">
                <div class="w-8 h-8 rounded-lg bg-surface-dark border border-border flex items-center justify-center mr-3 group-hover/item:border-primary/30 transition-colors">
                  <FolderOpen class="w-4 h-4 text-primary" />
                </div>
                <span class="font-bold text-xs tracking-tight text-text transition-colors">{{ item.name }}</span>
              </div>
            </template>

            <template #cell(type)="{ item }">
              <span class="px-2 py-0.5 rounded bg-white/5 border border-white/5 text-[9px] font-bold text-text-muted">
                {{ item.type }}
              </span>
            </template>

            <template #cell(fields)="{ item }">
              <span class="text-[10px] font-bold text-text-dim tracking-wider">{{ item.fields?.length ?? 0 }} fields</span>
            </template>

            <template #cell(records)="{ item }">
              <span class="text-xs font-mono text-text-muted">{{ recordCounts[item.name as string] ?? 0 }} records</span>
            </template>

            <template #cell(created)="{ item }">
              <span class="text-[10px] font-medium text-text-dim">{{ item.created ? new Date(item.created as string).toLocaleDateString() : '-' }}</span>
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
                    <PopoverItem :icon="Settings" @click="close(); router.push(`/collections/${item.name}/settings`);">Settings</PopoverItem>
                    <PopoverItem :icon="Trash2" variant="danger" @click="close(); handleDeleteClick(item);">Delete</PopoverItem>
                  </template>
                </Popover>
              </div>
            </template>

            <template #empty>
              <div class="py-16 text-center">
                <Database class="w-12 h-12 mx-auto mb-4 opacity-10" />
                <p class="text-xs font-bold tracking-widest text-text-dim">No collections found</p>
                <Button variant="link" class="mt-4 text-[10px] font-bold tracking-widest" @click="router.push('/collections/new')">Create your first collection</Button>
              </div>
            </template>
          </Table>
        </div>
      </div>
    </main>

    <ConfirmModal
      :show="showDeleteModal"
      title="Confirm Delete"
      :message="`Are you sure you want to delete the collection '${collectionToDelete?.name}'? This action cannot be undone.`"
      confirm-text="Delete"
      cancel-text="Cancel"
      variant="danger"
      @confirm="confirmDelete"
      @cancel="showDeleteModal = false"
    />
  </AppLayout>
</template>
