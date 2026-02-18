<script setup lang="ts">
import { ref, computed } from 'vue';
import Dropdown from './Dropdown.vue';
import DropdownItem from './DropdownItem.vue';
import Checkbox from './Checkbox.vue';
import {
  ChevronFirst,
  ChevronLast,
  ChevronLeft,
  ChevronRight,
  ArrowUp,
  ArrowDown,
} from 'lucide-vue-next';

interface Header {
  key: string;
  label: string;
  align?: 'left' | 'center' | 'right';
  sticky?: boolean;
  sortable?: boolean;
}

interface Props {
  headers: Header[];
  items: Record<string, unknown>[];
  loading?: boolean;
  emptyText?: string;
  rowClickable?: boolean;
  enablePagination?: boolean;
  defaultPageSize?: number;
  sortKey?: string;
  sortOrder?: 'asc' | 'desc';
  selectable?: boolean;
  selectionKey?: string;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  emptyText: 'No data found',
  rowClickable: false,
  enablePagination: false,
  defaultPageSize: 20,
  sortKey: undefined,
  sortOrder: 'asc',
  selectable: false,
  selectionKey: 'id',
});

const emit = defineEmits<{
  rowClick: [item: Record<string, unknown>, event: Event];
  sortChange: [sortKey: string, sortOrder: 'asc' | 'desc'];
  selectionChange: [selectedKeys: any[]];
}>();

const selectedKeys = ref<any[]>([]);

const toggleItem = (item: Record<string, unknown>) => {
  const key = item[props.selectionKey] as any;
  const index = selectedKeys.value.indexOf(key);
  if (index > -1) {
    selectedKeys.value.splice(index, 1);
  } else {
    selectedKeys.value.push(key);
  }
  emit('selectionChange', selectedKeys.value);
};

const isItemSelected = (item: Record<string, unknown>) => {
  return selectedKeys.value.includes(item[props.selectionKey] as any);
};

const isAllSelected = computed({
  get: () => {
    if (paginatedItems.value.length === 0) return false;
    return paginatedItems.value.every((item) =>
      selectedKeys.value.includes(item[props.selectionKey] as any)
    );
  },
  set: (value: boolean) => {
    if (value) {
      const currentPageKeys = paginatedItems.value.map((item) => item[props.selectionKey] as any);
      // Merge with already selected keys (if any from other pages)
      const otherPageKeys = selectedKeys.value.filter(
        (key) => !paginatedItems.value.some((item) => item[props.selectionKey] === key)
      );
      selectedKeys.value = [...otherPageKeys, ...currentPageKeys];
    } else {
      // Unselect only current page items
      selectedKeys.value = selectedKeys.value.filter(
        (key) => !paginatedItems.value.some((item) => item[props.selectionKey] === key)
      );
    }
    emit('selectionChange', selectedKeys.value);
  },
});

const handleSort = (headerKey: string) => {
  if (!props.headers.find((h) => h.key === headerKey)?.sortable) {
    return;
  }

  let newSortOrder: 'asc' | 'desc' = 'asc';
  if (props.sortKey === headerKey) {
    // Toggle order if clicking the same column
    newSortOrder = props.sortOrder === 'asc' ? 'desc' : 'asc';
  }

  emit('sortChange', headerKey, newSortOrder);
};

const currentPage = ref(1);
const pageSize = ref(props.defaultPageSize);

// Sort items based on sortKey and sortOrder
const sortedItems = computed(() => {
  if (!props.sortKey) {
    return props.items;
  }

  return [...props.items].sort((a, b) => {
    // Check if sortKey is defined before accessing
    if (!props.sortKey) return 0;

    const aValue = a[props.sortKey as keyof typeof a];
    const bValue = b[props.sortKey as keyof typeof b];

    // Handle null/undefined values
    if (aValue == null && bValue == null) return 0;
    if (aValue == null) return props.sortOrder === 'asc' ? 1 : -1;
    if (bValue == null) return props.sortOrder === 'asc' ? -1 : 1;

    // Handle dates
    if (
      typeof aValue === 'string' &&
      !isNaN(Date.parse(aValue)) &&
      typeof bValue === 'string' &&
      !isNaN(Date.parse(bValue))
    ) {
      const dateA = new Date(aValue).getTime();
      const dateB = new Date(bValue).getTime();
      return props.sortOrder === 'asc' ? dateA - dateB : dateB - dateA;
    }

    // Handle numbers
    if (typeof aValue === 'number' && typeof bValue === 'number') {
      return props.sortOrder === 'asc' ? aValue - bValue : bValue - aValue;
    }

    // Handle strings
    if (typeof aValue === 'string' && typeof bValue === 'string') {
      return props.sortOrder === 'asc'
        ? aValue.localeCompare(bValue)
        : bValue.localeCompare(aValue);
    }

    // Fallback comparison
    const strA = String(aValue);
    const strB = String(bValue);
    return props.sortOrder === 'asc' ? strA.localeCompare(strB) : strB.localeCompare(strA);
  });
});

const totalPages = computed(() => Math.ceil(sortedItems.value.length / pageSize.value));

const paginatedItems = computed(() => {
  if (!props.enablePagination) {
    return sortedItems.value;
  }

  const startIndex = (currentPage.value - 1) * pageSize.value;
  const endIndex = startIndex + pageSize.value;
  return sortedItems.value.slice(startIndex, endIndex);
});

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
  }
};

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value = currentPage.value - 1;
  }
};

const firstPage = () => {
  currentPage.value = 1;
};

const lastPage = () => {
  currentPage.value = totalPages.value;
};

const getAlignClass = (align?: string) => {
  if (align === 'center') return 'text-center';
  if (align === 'right') return 'text-right';
  return 'text-left';
};

const getFlexAlignClass = (align?: string) => {
  if (align === 'center') return 'justify-center';
  if (align === 'right') return 'justify-end';
  return 'justify-start';
};
</script>

<template>
  <div class="bg-surface/30 backdrop-blur-md rounded-2xl border border-border/50 shadow-2xl overflow-hidden font-mono selection:bg-primary/20">
    <div class="overflow-x-auto custom-scrollbar">
      <table class="w-full text-left text-xs whitespace-nowrap border-collapse">
        <thead>
          <tr class="bg-surface-dark/50 border-b border-border/50">
            <th v-if="selectable" class="px-5 py-4 w-4">
              <Checkbox v-model="isAllSelected" size="sm" />
            </th>
            <th
              v-for="header in headers"
              :key="header.key"
              :class="[
                'px-6 py-4 text-[10px] font-black text-text-dim uppercase tracking-[0.2em]',
                getAlignClass(header.align),
                header.sticky ? 'sticky right-0 bg-surface-dark z-10 shadow-[-10px_0_15px_rgba(0,0,0,0.1)]' : '',
                header.sortable ? 'cursor-pointer select-none hover:text-primary transition-colors' : '',
              ]"
              @click="header.sortable ? handleSort(header.key) : null"
            >
              <div class="flex items-center gap-2" :class="getFlexAlignClass(header.align)">
                <span>{{ header.label }}</span>
                <div v-if="header.sortable" class="flex flex-col opacity-30 group-hover:opacity-100">
                  <ArrowUp
                    :class="[
                      'w-3 h-3 -mb-1',
                      sortKey === header.key && sortOrder === 'asc'
                        ? 'text-primary'
                        : 'text-text-muted',
                    ]"
                  />
                  <ArrowDown
                    :class="[
                      'w-3 h-3',
                      sortKey === header.key && sortOrder === 'desc'
                        ? 'text-primary'
                        : 'text-text-muted',
                    ]"
                  />
                </div>
              </div>
            </th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border/20">
          <tr
            v-for="(item, index) in paginatedItems"
            :key="enablePagination ? (currentPage - 1) * pageSize + index : index"
            :class="[
              'group transition-all duration-200',
              rowClickable ? 'cursor-pointer hover:bg-primary/5' : '',
            ]"
            @click="rowClickable ? emit('rowClick', item, $event) : null"
          >
            <td v-if="selectable" class="px-5 py-3 w-4">
              <Checkbox
                :model-value="isItemSelected(item)"
                size="sm"
                @update:model-value="toggleItem(item)"
                @click.stop
              />
            </td>
            <td
              v-for="header in headers"
              :key="header.key"
              :class="[
                'px-6 py-3.5',
                getAlignClass(header.align),
                header.sticky ? 'sticky right-0 bg-surface/80 backdrop-blur-md z-10 shadow-[-10px_0_15px_rgba(0,0,0,0.05)]' : '',
              ]"
            >
              <slot :name="`cell(${header.key})`" :item="item">
                <span class="text-text font-medium">{{ item[header.key] ?? '-' }}</span>
              </slot>
            </td>
          </tr>

          <tr v-if="paginatedItems.length === 0 && !loading">
            <td :colspan="headers.length + (selectable ? 1 : 0)" class="px-6 py-16 text-center">
              <slot name="empty">
                <p class="text-[10px] font-black uppercase tracking-widest text-text-dim">{{ emptyText }}</p>
              </slot>
            </td>
          </tr>

          <tr v-if="loading">
            <td :colspan="headers.length + (selectable ? 1 : 0)" class="px-6 py-16 text-center">
              <div class="flex justify-center">
                <div class="w-8 h-8 border-2 border-primary/20 border-t-primary rounded-full animate-spin"></div>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Default Pagination Footer -->
    <div
      v-if="enablePagination && items.length > 0"
      class="bg-surface-dark/30 px-6 py-4 border-t border-border/30 flex flex-col sm:flex-row items-center justify-between gap-4"
    >
      <div class="text-[10px] font-bold text-text-dim uppercase tracking-widest flex items-center gap-3">
        <div class="flex items-center whitespace-nowrap">
          <span class="mr-2">Density</span>
          <Dropdown align="left" size="xs" headless>
            <template #trigger>
              <button class="px-2 py-1 rounded bg-white/5 border border-white/5 hover:border-primary/30 transition-all flex items-center gap-1">
                {{ pageSize }}
                <ChevronDown class="w-3 h-3 text-text-dim" />
              </button>
            </template>
            <template #default>
              <DropdownItem :value="10" @click="pageSize = 10">10</DropdownItem>
              <DropdownItem :value="20" @click="pageSize = 20">20</DropdownItem>
              <DropdownItem :value="50" @click="pageSize = 50">50</DropdownItem>
              <DropdownItem :value="100" @click="pageSize = 100">100</DropdownItem>
            </template>
          </Dropdown>
          <span class="ml-3">Records: {{ items.length }}</span>
        </div>
      </div>
      
      <div class="flex items-center gap-2">
        <button
          :disabled="currentPage <= 1"
          class="p-2 rounded-lg bg-white/5 border border-white/5 hover:border-primary/30 disabled:opacity-30 disabled:cursor-not-allowed transition-all"
          @click="firstPage"
        >
          <ChevronFirst class="w-3.5 h-3.5 text-text" />
        </button>
        <button
          :disabled="currentPage <= 1"
          class="p-2 rounded-lg bg-white/5 border border-white/5 hover:border-primary/30 disabled:opacity-30 disabled:cursor-not-allowed transition-all"
          @click="prevPage"
        >
          <ChevronLeft class="w-3.5 h-3.5 text-text" />
        </button>
        
        <div class="flex items-center gap-1.5 mx-2">
          <span class="text-[10px] font-black text-primary">{{ currentPage }}</span>
          <span class="text-[10px] font-bold text-text-dim">/</span>
          <span class="text-[10px] font-bold text-text-dim">{{ totalPages }}</span>
        </div>

        <button
          :disabled="currentPage >= totalPages"
          class="p-2 rounded-lg bg-white/5 border border-white/5 hover:border-primary/30 disabled:opacity-30 disabled:cursor-not-allowed transition-all"
          @click="nextPage"
        >
          <ChevronRight class="w-3.5 h-3.5 text-text" />
        </button>
        <button
          :disabled="currentPage >= totalPages"
          class="p-2 rounded-lg bg-white/5 border border-white/5 hover:border-primary/30 disabled:opacity-30 disabled:cursor-not-allowed transition-all"
          @click="lastPage"
        >
          <ChevronLast class="w-3.5 h-3.5 text-text" />
        </button>
      </div>
    </div>

    <!-- Footer Slot (e.g. for custom pagination) -->
    <slot name="footer" />
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  height: 6px;
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: var(--border);
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: var(--primary);
}
</style>
