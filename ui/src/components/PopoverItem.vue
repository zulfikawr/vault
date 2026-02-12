<script setup lang="ts">
interface Props {
  icon?: any;
  variant?: 'default' | 'danger';
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default'
});

const emit = defineEmits<{
  click: [];
}>();
</script>

<template>
  <button
    @click="emit('click')"
    :class="[
      'w-full flex items-center gap-3 px-4 py-2.5 text-sm transition-colors text-left',
      variant === 'danger' 
        ? 'text-error hover:bg-error/10' 
        : 'text-text hover:bg-surface-dark'
    ]"
  >
    <component v-if="icon" :is="icon" class="w-4 h-4 flex-shrink-0" />
    <slot />
  </button>
</template>
