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
      'w-full bg-surface border border-border rounded-lg text-text placeholder-text-muted/50 focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition-all resize-none disabled:opacity-50 disabled:cursor-not-allowed',
      props.size === 'sm'
        ? 'px-3 py-1.5 text-sm'
        : props.size === 'lg'
          ? 'px-5 py-3 text-base'
          : 'px-4 py-2.5 text-base',
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
      'w-full bg-surface border border-border rounded-lg text-text placeholder-text-muted/50 focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition-all disabled:opacity-50 disabled:cursor-not-allowed',
      props.size === 'sm'
        ? 'px-3 py-1.5 text-sm'
        : props.size === 'lg'
          ? 'px-5 py-3 text-base'
          : 'px-4 py-2.5 text-base',
    ]"
    @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
  />
</template>
