<script setup lang="ts">
interface Props {
  modelValue: boolean;
  label?: string;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  'update:modelValue': [value: boolean];
}>();

const toggle = () => {
  emit('update:modelValue', !props.modelValue);
};
</script>

<template>
  <label class="flex items-center gap-2 cursor-pointer group">
    <div class="relative">
      <input
        type="checkbox"
        :checked="modelValue"
        @change="toggle"
        class="sr-only"
      />
      <div
        :class="[
          'w-4 h-4 rounded border-2 transition-all',
          modelValue
            ? 'bg-primary border-primary'
            : 'bg-surface-dark border-border group-hover:border-text-muted'
        ]"
      >
        <svg
          v-if="modelValue"
          class="w-full h-full text-white"
          viewBox="0 0 16 16"
          fill="none"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            d="M13 4L6 11L3 8"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
      </div>
    </div>
    <span v-if="label" class="text-sm text-text select-none">{{ label }}</span>
    <slot v-else />
  </label>
</template>

<style scoped>
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border-width: 0;
}
</style>
