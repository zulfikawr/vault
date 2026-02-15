<script setup lang="ts">
interface Tab {
  id: string;
  label: string;
}

defineProps<{
  tabs: Tab[];
  modelValue: string;
}>();

defineEmits<{
  (e: 'update:modelValue', value: string): void;
}>();
</script>

<template>
  <div class="border-b border-border mb-6">
    <nav class="flex space-x-8" aria-label="Tabs">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        :class="[
          modelValue === tab.id
            ? 'border-primary text-primary'
            : 'border-transparent text-text-muted hover:text-text hover:border-border',
          'whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm transition-all duration-200',
        ]"
        @click="$emit('update:modelValue', tab.id)"
      >
        {{ tab.label }}
      </button>
    </nav>
  </div>
</template>
