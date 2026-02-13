<script setup lang="ts">
interface Props {
  value: string | number;
  icon?: any;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  click: [value: string | number];
  select: [value: string | number];
}>();

const handleClick = () => {
  emit('click', props.value);
  emit('select', props.value);
};
</script>

<template>
  <button
    @click="handleClick"
    type="button"
    class="w-full flex items-center gap-3 px-4 py-2.5 text-sm transition-colors text-left text-text hover:bg-surface-dark"
  >
    <component v-if="icon" :is="icon" class="w-4 h-4 flex-shrink-0" />
    <slot />
  </button>
</template>
