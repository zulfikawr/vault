<script setup lang="ts">
import { ref, computed } from 'vue';
import Dropdown from './Dropdown.vue';
import DropdownItem from './DropdownItem.vue';
import { ChevronFirst, ChevronLast, ChevronLeft, ChevronRight, ArrowUp, ArrowDown } from 'lucide-vue-next';

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
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  emptyText: 'No data found',
  rowClickable: false,
  enablePagination: false,
  defaultPageSize: 20,
  sortKey: undefined,
  sortOrder: 'asc',
});

const emit = defineEmits<{
  rowClick: [item: Record<string, unknown>, event: Event];
  sortChange: [sortKey: string, sortOrder: 'asc' | 'desc'];
}>();

const handleSort = (headerKey: string) => {
  if (!props.headers.find(h => h.key === headerKey)?.sortable) {
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
    if (typeof aValue === 'string' && !isNaN(Date.parse(aValue)) && 
        typeof bValue === 'string' && !isNaN(Date.parse(bValue))) {
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
    return props.sortOrder === 'asc' 
      ? strA.localeCompare(strB) 
      : strB.localeCompare(strA);
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
  <div class="bg-surface-dark rounded-lg border border-border shadow-sm overflow-hidden">
    <div class="overflow-x-auto">
      <table class="w-full text-left text-sm whitespace-nowrap">
        <thead class="bg-surface border-b border-border">
          <tr>
            <th
              v-for="header in headers"
              :key="header.key"
              :class="[
                'px-4 sm:px-6 py-3 text-xs font-medium text-text-muted uppercase tracking-wider',
                getAlignClass(header.align),
                header.sticky ? 'sticky right-0 bg-surface z-10' : '',
                header.sortable ? 'cursor-pointer select-none hover:text-text' : ''
              ]"
              @click="header.sortable ? handleSort(header.key) : null"
            >
              <div class="flex items-center gap-1" :class="getFlexAlignClass(header.align)">
                <span>{{ header.label }}</span>
                <div v-if="header.sortable" class="flex flex-col">
                  <ArrowUp 
                    :class="[
                      'w-2 h-2',
                      sortKey === header.key && sortOrder === 'asc' ? 'text-text' : 'text-text-muted'
                    ]" 
                  />
                  <ArrowDown 
                    :class="[
                      'w-2 h-2',
                      sortKey === header.key && sortOrder === 'desc' ? 'text-text' : 'text-text-muted'
                    ]" 
                  />
                </div>
              </div>
            </th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border">
          <tr
            v-for="(item, index) in paginatedItems"
            :key="enablePagination ? (currentPage - 1) * pageSize + index : index"
            :class="[
              'hover:bg-background/50 transition-colors group',
              rowClickable ? 'cursor-pointer' : '',
            ]"
            @click="rowClickable ? emit('rowClick', item, $event) : null"
          >
            <td
              v-for="header in headers"
              :key="header.key"
              :class="[
                'px-4 sm:px-6 py-4',
                getAlignClass(header.align),
                header.sticky
                  ? 'sticky right-0 bg-transparent z-10'
                  : '',
              ]"
            >
              <slot :name="`cell(${header.key})`" :item="item">
                <span class="text-text">{{ item[header.key] ?? '-' }}</span>
              </slot>
            </td>
          </tr>

          <tr v-if="paginatedItems.length === 0 && !loading">
            <td :colspan="headers.length" class="px-4 sm:px-6 py-12 text-center text-text-muted">
              <slot name="empty">
                <p class="text-sm">{{ emptyText }}</p>
              </slot>
            </td>
          </tr>

          <tr v-if="loading">
            <td :colspan="headers.length" class="px-4 sm:px-6 py-12 text-center text-text-muted">
              <div class="flex justify-center">
                <svg
                  class="animate-spin h-6 w-6 text-primary"
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
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Default Pagination Footer -->
    <div
      v-if="enablePagination"
      class="bg-surface px-4 sm:px-6 py-3 border-t border-border flex items-center justify-between"
    >
      <div class="text-xs text-text-muted flex items-center gap-1">
        <div class="flex items-center whitespace-nowrap flex-shrink-0">
          <span class="mr-1">Showing</span>
          <Dropdown align="left">
            <template #trigger>
              <span class="pr-2">
                {{pageSize}}
              </span>
            </template>
            <template #default>
              <DropdownItem :value="10" @click="pageSize = 10">10</DropdownItem>
              <DropdownItem :value="20" @click="pageSize = 20">20</DropdownItem>
              <DropdownItem :value="50" @click="pageSize = 50">50</DropdownItem>
              <DropdownItem :value="100" @click="pageSize = 100">100</DropdownItem>
            </template>
          </Dropdown>
          <span class="ml-1">of {{ items.length }} results</span>
        </div>
      </div>
      <div class="flex items-center gap-1">
        <button
          :disabled="currentPage <= 1"
          @click="firstPage"
          class="p-1.5 rounded hover:bg-surface-dark disabled:opacity-50 disabled:cursor-not-allowed"
          :class="{ 'opacity-50 cursor-not-allowed': currentPage <= 1 }"
        >
          <ChevronFirst class="w-4 h-4 text-text" />
        </button>
        <button
          :disabled="currentPage <= 1"
          @click="prevPage"
          class="p-1.5 rounded hover:bg-surface-dark disabled:opacity-50 disabled:cursor-not-allowed"
          :class="{ 'opacity-50 cursor-not-allowed': currentPage <= 1 }"
        >
          <ChevronLeft class="w-4 h-4 text-text" />
        </button>
        <span class="px-2 text-sm text-text">
          {{ currentPage }} of {{ totalPages }}
        </span>
        <button
          :disabled="currentPage >= totalPages"
          @click="nextPage"
          class="p-1.5 rounded hover:bg-surface-dark disabled:opacity-50 disabled:cursor-not-allowed"
          :class="{ 'opacity-50 cursor-not-allowed': currentPage >= totalPages }"
        >
          <ChevronRight class="w-4 h-4 text-text" />
        </button>
        <button
          :disabled="currentPage >= totalPages"
          @click="lastPage"
          class="p-1.5 rounded hover:bg-surface-dark disabled:opacity-50 disabled:cursor-not-allowed"
          :class="{ 'opacity-50 cursor-not-allowed': currentPage >= totalPages }"
        >
          <ChevronLast class="w-4 h-4 text-text" />
        </button>
      </div>
    </div>

    <!-- Footer Slot (e.g. for custom pagination) -->
    <slot name="footer" />
  </div>
</template>
