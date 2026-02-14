<script setup lang="ts">
import { ref, computed } from 'vue';
import Dropdown from './Dropdown.vue';
import DropdownItem from './DropdownItem.vue';
import { ChevronFirst, ChevronLast, ChevronLeft, ChevronRight } from 'lucide-vue-next';

interface Header {
  key: string;
  label: string;
  align?: 'left' | 'center' | 'right';
  sticky?: boolean;
}

interface Props {
  headers: Header[];
  items: Record<string, unknown>[];
  loading?: boolean;
  emptyText?: string;
  rowClickable?: boolean;
  enablePagination?: boolean;
  defaultPageSize?: number;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  emptyText: 'No data found',
  rowClickable: false,
  enablePagination: false,
  defaultPageSize: 20,
});

const emit = defineEmits<{
  rowClick: [item: Record<string, unknown>];
}>();

const currentPage = ref(1);
const pageSize = ref(props.defaultPageSize);

const totalPages = computed(() => Math.ceil(props.items.length / pageSize.value));

const paginatedItems = computed(() => {
  if (!props.enablePagination) {
    return props.items;
  }
  
  const startIndex = (currentPage.value - 1) * pageSize.value;
  const endIndex = startIndex + pageSize.value;
  return props.items.slice(startIndex, endIndex);
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
              ]"
            >
              {{ header.label }}
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
            @click="rowClickable ? emit('rowClick', item) : null"
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
