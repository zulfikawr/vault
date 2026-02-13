<script setup lang="ts">
interface Props {
  icon?: unknown;
  variant?: 'default' | 'danger';
}

withDefaults(defineProps<Props>(), {
  icon: undefined,
  variant: 'default',
});

const emit = defineEmits<{
  click: [];
}>();
</script>

<template>
  <button
    :class="[
      'w-full flex items-center gap-3 px-4 py-2.5 text-sm transition-colors text-left',
      variant === 'danger' ? 'text-error hover:bg-error/10' : 'text-text hover:bg-surface-dark',
    ]"
    @click="emit('click')"
  >
    <component :is="icon" v-if="icon" class="w-4 h-4 flex-shrink-0" />
    <slot />
  </button>
</template>
