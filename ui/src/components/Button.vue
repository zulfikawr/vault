<script setup lang="ts">
import { computed } from 'vue';
import { Loader2 } from 'lucide-vue-next';

interface Props {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'link';
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl';
  type?: 'button' | 'submit' | 'reset';
  disabled?: boolean;
  loading?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  type: 'button',
  disabled: false,
  loading: false
});

const variantClasses = {
  primary: 'bg-primary text-white hover:bg-primary-hover shadow-sm',
  secondary: 'bg-surface border border-border text-text hover:bg-surface-dark shadow-sm',
  outline: 'bg-transparent border border-primary text-primary hover:bg-primary/10',
  ghost: 'bg-transparent text-text-muted hover:bg-surface-dark hover:text-text',
  link: 'bg-transparent text-primary hover:underline !p-0 !h-auto'
};

const sizeClasses = {
  xs: 'px-2 py-1 text-xs',
  sm: 'px-3 py-1.5 text-sm',
  md: 'px-4 py-2 text-sm sm:text-base',
  lg: 'px-6 py-2.5 text-base sm:text-lg',
  xl: 'px-8 py-3.5 text-lg sm:text-xl'
};

const classes = computed(() => {
  return [
    'inline-flex items-center justify-center gap-2 font-medium transition-all duration-200 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed active:scale-[0.98] outline-none focus:ring-2 focus:ring-primary/20',
    variantClasses[props.variant],
    props.variant !== 'link' ? sizeClasses[props.size] : ''
  ];
});
</script>

<template>
  <button :type="type" :class="classes" :disabled="disabled || loading">
    <slot v-if="!loading" name="leftIcon" />
    <Loader2 v-if="loading" class="h-4 w-4 animate-spin shrink-0" />
    <slot />
    <slot v-if="!loading" name="rightIcon" />
  </button>
</template>
