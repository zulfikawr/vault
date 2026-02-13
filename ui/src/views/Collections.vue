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
}

const router = useRouter();
const collections = ref<Collection[]>([]);
const showSystemCollections = ref(true);
const filterTypes = ref({
  base: true,
  auth: true,
  system: true,
});

const filteredCollections = computed(() => {
  return collections.value.filter((col: Collection) => {
    // Filter by system collections
    if (!showSystemCollections.value && col.name.startsWith('_')) {
      return false;
    }
    // Filter by type
    return filterTypes.value[col.type];
  });
});

const fetchCollections = async () => {
  try {
    const response = await axios.get('/api/admin/collections');
    collections.value = response.data.data;
  } catch (error) {
    console.error('Failed to fetch collections', error);
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
        <div class="flex items-center text-sm text-text-muted">
          <span class="hover:text-text cursor-pointer" @click="router.push('/')">Vault</span>
          <span class="mx-2">/</span>
          <span class="font-medium text-text">Collections</span>
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
                <Button variant="secondary" class="flex-1 sm:flex-none">
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
            <Button class="flex-1 sm:flex-none" @click="router.push('/collections/new')">
              <Plus class="w-4 h-4" />
              New Collection
            </Button>
          </div>
        </div>

        <!-- Data Table -->
        <Table
          :headers="[
            { key: 'name', label: 'Name' },
            { key: 'type', label: 'Type' },
            { key: 'fields', label: 'Fields' },
            { key: 'created', label: 'Created' },
            { key: 'status', label: 'Status' },
            { key: 'actions', label: 'Actions', align: 'center', sticky: true },
          ]"
          :items="filteredCollections"
        >
          <template #cell(name)="{ item }">
            <div
              class="flex items-center gap-3 cursor-pointer"
              @click="router.push(`/collections/${item.name}`)"
            >
              <div class="p-1.5 rounded bg-primary/10 text-primary">
                <FolderOpen class="w-4 h-4" />
              </div>
              <span class="font-medium text-text">{{ item.name }}</span>
            </div>
          </template>

          <template #cell(type)="{ item }">
            <span
              class="text-text-muted cursor-pointer"
              @click="router.push(`/collections/${item.name}`)"
              >{{ item.type }}</span
            >
          </template>

          <template #cell(fields)="{ item }">
            <span
              class="text-text-muted cursor-pointer"
              @click="router.push(`/collections/${(item as unknown as Collection).name}`)"
              >{{ (item as unknown as Collection).fields?.length ?? 0 }} fields</span
            >
          </template>

          <template #cell(created)="{ item }">
            <span
              class="text-text-muted text-xs cursor-pointer"
              @click="router.push(`/collections/${item.name}`)"
              >{{
                item.created ? new Date(item.created as string).toLocaleDateString() : '-'
              }}</span
            >
          </template>

          <template #cell(status)="{ item }">
            <span
              class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-success/10 text-success cursor-pointer"
              @click="router.push(`/collections/${item.name}`)"
            >
              Active
            </span>
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
                  :icon="Settings"
                  @click="
                    close();
                    router.push(`/collections/${item.name}/settings`);
                  "
                >
                  Settings
                </PopoverItem>
                <PopoverItem :icon="Trash2" variant="danger" @click="close()"> Delete </PopoverItem>
              </template>
            </Popover>
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

          <template #footer>
            <div
              class="bg-surface px-4 sm:px-6 py-3 border-t border-border flex items-center justify-between"
            >
              <div class="text-xs text-text-muted">
                Showing
                <span class="font-medium text-text">{{ filteredCollections.length }}</span> of
                <span class="font-medium text-text">{{ collections.length }}</span> results
              </div>
              <div class="flex gap-2">
                <Button variant="secondary" size="xs" disabled>Previous</Button>
                <Button variant="secondary" size="xs">Next</Button>
              </div>
            </div>
          </template>
        </Table>
      </div>
    </div>
  </AppLayout>
</template>
