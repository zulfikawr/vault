<script setup lang="ts">
interface Props {
  modelValue?: string | number | boolean;
  type?: string;
  placeholder?: string;
  required?: boolean;
  disabled?: boolean;
  readonly?: boolean;
  rows?: number;
  size?: 'sm' | 'md' | 'lg';
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  type: 'text',
  placeholder: '',
  required: false,
  disabled: false,
  readonly: false,
  rows: 4,
  size: 'md',
});

const emit = defineEmits<{
  'update:modelValue': [value: string | number];
}>();

const isTextarea = () => props.type === 'textarea';

const baseClasses = 'w-full bg-surface-dark/50 backdrop-blur-sm border border-border/50 rounded-xl text-text placeholder:text-text-dim/40 focus:outline-none focus:ring-4 focus:ring-primary/10 focus:border-primary/50 transition-all duration-300 font-mono disabled:opacity-30 disabled:cursor-not-allowed selection:bg-primary/20 shadow-inner';
const sizeClasses = {
  sm: 'px-3 h-8 text-xs',
  md: 'px-4 h-10 text-sm',
  lg: 'px-6 h-12 text-base',
};
</script>

<template>
  <textarea
    v-if="isTextarea()"
    :value="modelValue"
    :placeholder="placeholder"
    :required="required"
    :disabled="disabled"
    :readonly="readonly"
    :rows="rows || 4"
    :class="[
      'w-full bg-surface-dark/50 backdrop-blur-sm border border-border/50 rounded-xl text-text placeholder:text-text-dim/40 focus:outline-none focus:ring-4 focus:ring-primary/10 focus:border-primary/50 transition-all duration-300 font-mono disabled:opacity-30 disabled:cursor-not-allowed selection:bg-primary/20 shadow-inner',
      props.size === 'sm' ? 'px-3 py-1.5 text-xs' : props.size === 'lg' ? 'px-6 py-4 text-base' : 'px-4 py-2.5 text-sm',
      'resize-none'
    ]"
    @input="emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
  />
  <input
    v-else
    :value="modelValue"
    :type="type"
    :placeholder="placeholder"
    :required="required"
    :disabled="disabled"
    :readonly="readonly"
    :class="[
      baseClasses,
      sizeClasses[props.size],
    ]"
    @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
  />
</template>
