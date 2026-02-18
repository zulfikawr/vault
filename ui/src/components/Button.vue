<script setup lang="ts">
import { computed } from 'vue';
import { Loader2 } from 'lucide-vue-next';

interface Props {
  variant?: 'primary' | 'secondary' | 'outline' | 'destructive' | 'ghost' | 'link';
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl';
  type?: 'button' | 'submit' | 'reset';
  disabled?: boolean;
  loading?: boolean;
  as?: string;
  to?: string;
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  type: 'button',
  disabled: false,
  loading: false,
  as: 'button',
  to: undefined,
});

const variantClasses = {
  primary: 'bg-primary text-white hover:bg-primary-hover shadow-[0_0_15px_rgba(var(--primary-rgb),0.2)] hover:shadow-[0_0_20px_rgba(var(--primary-rgb),0.4)] border border-primary/20',
  secondary: 'bg-surface-dark border border-border text-text hover:bg-surface hover:border-text-muted/30 shadow-sm',
  outline: 'bg-transparent border-2 border-primary text-primary hover:bg-primary/10 hover:shadow-[0_0_15px_rgba(var(--primary-rgb),0.2)]',
  destructive: 'bg-error/10 border border-error/30 text-error hover:bg-error hover:text-white shadow-sm',
  ghost: 'bg-transparent text-text-muted hover:bg-white/5 hover:text-text',
  link: 'bg-transparent text-primary hover:underline !p-0 !h-auto !ring-0',
};

const sizeClasses = {
  xs: 'px-2 h-7 text-[10px] font-bold uppercase tracking-wider',
  sm: 'px-3 h-8 text-xs font-bold uppercase tracking-wide',
  md: 'px-4 h-10 text-xs sm:text-sm font-bold uppercase tracking-wide',
  lg: 'px-6 h-12 text-sm sm:text-base font-bold uppercase tracking-wide',
  xl: 'px-8 h-14 text-base sm:text-lg font-bold uppercase tracking-wide',
};

const classes = computed(() => {
  return [
    'inline-flex items-center justify-center gap-2 transition-all duration-300 rounded-xl disabled:opacity-50 disabled:cursor-not-allowed active:scale-[0.96] outline-none focus:ring-2 focus:ring-primary/40 focus:ring-offset-2 focus:ring-offset-background',
    variantClasses[props.variant],
    props.variant !== 'link' ? sizeClasses[props.size] : '',
  ];
});
</script>

<template>
  <component
    :is="as"
    :type="as === 'button' ? type : undefined"
    :to="to"
    :class="classes"
    :disabled="disabled || loading"
  >
    <slot v-if="!loading" name="leftIcon" />
    <Loader2 v-if="loading" class="h-4 w-4 animate-spin shrink-0" />
    <slot />
    <slot v-if="!loading" name="rightIcon" />
  </component>
</template>
