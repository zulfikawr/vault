<script setup lang="ts">
interface Header {
  key: string;
  label: string;
  align?: 'left' | 'center' | 'right';
  sticky?: boolean;
}

interface Props {
  headers: Header[];
  items: any[];
  loading?: boolean;
  emptyText?: string;
  rowClickable?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  emptyText: 'No data found',
  rowClickable: false
});

const emit = defineEmits<{
  rowClick: [item: any];
}>();

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
                header.sticky ? 'sticky right-0 bg-surface z-10' : ''
              ]"
            >
              {{ header.label }}
            </th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border">
          <tr 
            v-for="(item, index) in items" 
            :key="index"
            @click="rowClickable ? emit('rowClick', item) : null"
            :class="[
              'hover:bg-background/50 transition-colors group',
              rowClickable ? 'cursor-pointer' : ''
            ]"
          >
            <td 
              v-for="header in headers" 
              :key="header.key"
              :class="[
                'px-4 sm:px-6 py-4',
                getAlignClass(header.align),
                header.sticky ? 'sticky right-0 bg-surface-dark group-hover:bg-background z-10' : ''
              ]"
            >
              <slot :name="`cell(${header.key})`" :item="item">
                <span class="text-text">{{ item[header.key] ?? '-' }}</span>
              </slot>
            </td>
          </tr>
          
          <tr v-if="items.length === 0 && !loading">
            <td :colspan="headers.length" class="px-4 sm:px-6 py-12 text-center text-text-muted">
              <slot name="empty">
                <p class="text-sm">{{ emptyText }}</p>
              </slot>
            </td>
          </tr>

          <tr v-if="loading">
            <td :colspan="headers.length" class="px-4 sm:px-6 py-12 text-center text-text-muted">
              <div class="flex justify-center">
                <svg class="animate-spin h-6 w-6 text-primary" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    
    <!-- Footer Slot (e.g. for pagination) -->
    <slot name="footer" />
  </div>
</template>
