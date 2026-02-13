<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Popover from '../components/Popover.vue';
import PopoverItem from '../components/PopoverItem.vue';
import Checkbox from '../components/Checkbox.vue';
import { 
  FolderOpen,
  Filter,
  Plus,
  MoreHorizontal,
  Settings,
  Trash2
} from 'lucide-vue-next';

const router = useRouter();
const collections = ref([]);
const showSystemCollections = ref(true);
const filterTypes = ref({
  base: true,
  auth: true,
  system: true
});

const filteredCollections = computed(() => {
  return collections.value.filter((col: any) => {
    // Filter by system collections
    if (!showSystemCollections.value && col.name.startsWith('_')) {
      return false;
    }
    // Filter by type
    return filterTypes.value[col.type as keyof typeof filterTypes.value];
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
              <p class="mt-1 text-sm text-text-muted">Manage your database schemas and content types.</p>
            </div>
            <div class="flex items-center gap-3">
              <Popover align="right">
                <template #trigger>
                  <button class="flex-1 sm:flex-none px-4 py-2 bg-surface-dark border border-border rounded text-sm font-medium text-text hover:bg-surface transition-colors flex items-center justify-center gap-2">
                    <Filter class="w-4 h-4" />
                    Filter
                  </button>
                </template>
                <template #default>
                  <div class="p-3 min-w-[200px]">
                    <div class="text-xs font-semibold text-text-muted uppercase tracking-wider mb-2 px-1">
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
              <button @click="router.push('/collections/new')" class="flex-1 sm:flex-none px-4 py-2 bg-primary hover:bg-primary-hover text-white rounded text-sm font-medium shadow-sm hover:shadow transition-all flex items-center justify-center gap-2">
                <Plus class="w-4 h-4" />
                New Collection
              </button>
            </div>
          </div>

          <!-- Data Table -->
          <div class="bg-surface-dark rounded-lg border border-border shadow-sm overflow-hidden">
            <div class="overflow-x-auto">
              <table class="w-full text-left text-sm whitespace-nowrap">
                <thead class="bg-surface border-b border-border">
                  <tr>
                    <th class="px-6 py-3 text-xs font-medium text-text-muted uppercase tracking-wider">Name</th>
                    <th class="px-6 py-3 text-xs font-medium text-text-muted uppercase tracking-wider">Type</th>
                    <th class="px-6 py-3 text-xs font-medium text-text-muted uppercase tracking-wider">Fields</th>
                    <th class="px-6 py-3 text-xs font-medium text-text-muted uppercase tracking-wider">Created</th>
                    <th class="px-6 py-3 text-xs font-medium text-text-muted uppercase tracking-wider">Status</th>
                    <th class="sticky right-0 bg-surface px-6 py-3 text-center text-xs font-medium text-text-muted uppercase tracking-wider">Actions</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-border">
                  <tr v-for="col in filteredCollections" :key="col.name" class="hover:bg-background/50 transition-colors group">
                    <td @click="router.push(`/collections/${col.name}`)" class="px-6 py-4 cursor-pointer">
                      <div class="flex items-center gap-3">
                        <div class="p-1.5 rounded bg-primary/10 text-primary">
                          <FolderOpen class="w-4 h-4" />
                        </div>
                        <span class="font-medium text-text">{{ col.name }}</span>
                      </div>
                    </td>
                    <td @click="router.push(`/collections/${col.name}`)" class="px-6 py-4 cursor-pointer">
                      <span class="text-text-muted">{{ col.type }}</span>
                    </td>
                    <td @click="router.push(`/collections/${col.name}`)" class="px-6 py-4 cursor-pointer">
                      <span class="text-text-muted">{{ col.fields?.length || 0 }} fields</span>
                    </td>
                    <td @click="router.push(`/collections/${col.name}`)" class="px-6 py-4 cursor-pointer">
                      <span class="text-text-muted text-xs">{{ col.created ? new Date(col.created).toLocaleDateString() : '-' }}</span>
                    </td>
                    <td @click="router.push(`/collections/${col.name}`)" class="px-6 py-4 cursor-pointer">
                      <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-success/10 text-success">
                        Active
                      </span>
                    </td>
                    <td class="sticky right-0 bg-surface-dark px-6 py-4 group-hover:bg-background">
                      <div class="flex justify-center">
                      <Popover align="right">
                        <template #trigger>
                          <button class="text-text-muted hover:text-text hover:bg-surface-dark rounded p-1 transition-colors">
                            <MoreHorizontal class="w-4 h-4" />
                          </button>
                        </template>
                        <template #default="{ close }">
                          <PopoverItem 
                            :icon="Settings" 
                            @click="close(); router.push(`/collections/${col.name}/settings`)"
                          >
                            Settings
                          </PopoverItem>
                          <PopoverItem 
                            :icon="Trash2" 
                            variant="danger"
                            @click="close()"
                          >
                            Delete
                          </PopoverItem>
                        </template>
                      </Popover>
                      </div>
                    </td>
                  </tr>
                  <tr v-if="filteredCollections.length === 0">
                    <td colspan="6" class="px-6 py-12 text-center text-text-muted">
                      <FolderOpen class="w-12 h-12 mx-auto mb-3 opacity-30" />
                      <p class="text-sm mb-4">No collections found</p>
                      <button @click="router.push('/collections/new')" class="text-sm text-primary hover:text-primary-hover">Create your first collection</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
            
            <!-- Pagination -->
            <div class="bg-surface px-6 py-3 border-t border-border flex items-center justify-between">
              <div class="text-xs text-text-muted">
                Showing <span class="font-medium text-text">{{ filteredCollections.length }}</span> of <span class="font-medium text-text">{{ collections.length }}</span> results
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


