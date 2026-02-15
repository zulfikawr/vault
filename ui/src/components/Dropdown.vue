<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue';
import { ChevronDown } from 'lucide-vue-next';

interface Props {
  align?: 'left' | 'right';
  modelValue?: string | number;
  size?: 'xs' | 'sm' | 'md' | 'lg';
}

const props = withDefaults(defineProps<Props>(), {
  align: 'left',
  modelValue: '',
  size: 'md',
});

defineEmits<{
  'update:modelValue': [value: string | number];
}>();

const isOpen = ref(false);
const triggerRef = ref<HTMLElement | null>(null);
const dropdownRef = ref<HTMLElement | null>(null);
const coords = ref({ top: 0, left: 0, width: 0 });

const updatePosition = async () => {
  if (!isOpen.value || !triggerRef.value) return;

  await nextTick();
  if (!dropdownRef.value) return;

  const triggerRect = triggerRef.value.getBoundingClientRect();
  const dropdownRect = dropdownRef.value.getBoundingClientRect();
  const viewportWidth = window.innerWidth;
  const viewportHeight = window.innerHeight;

  let top = triggerRect.bottom + 8;
  let left = props.align === 'left' ? triggerRect.left : triggerRect.right - triggerRect.width;
  let width = triggerRect.width;

  // Vertical check: if it goes below viewport, show it above the trigger
  if (top + dropdownRect.height > viewportHeight - 12) {
    const spaceAbove = triggerRect.top;
    const spaceBelow = viewportHeight - triggerRect.bottom;

    if (spaceAbove > spaceBelow) {
      top = triggerRect.top - dropdownRect.height - 8;
    }
  }

  // Horizontal check: ensure it stays within viewport
  if (left + dropdownRect.width > viewportWidth - 12) {
    left = viewportWidth - dropdownRect.width - 12;
  }
  if (left < 12) {
    left = 12;
  }

  coords.value = { top, left, width };
};

const toggle = (event: Event) => {
  event.stopPropagation();
  isOpen.value = !isOpen.value;
};

const close = () => {
  isOpen.value = false;
};

const handleClickOutside = (event: MouseEvent) => {
  if (
    triggerRef.value &&
    dropdownRef.value &&
    !triggerRef.value.contains(event.target as Node) &&
    !dropdownRef.value.contains(event.target as Node)
  ) {
    close();
  }
};

watch(isOpen, (val) => {
  if (val) {
    updatePosition();
    window.addEventListener('scroll', updatePosition, true);
    window.addEventListener('resize', updatePosition);
  } else {
    window.removeEventListener('scroll', updatePosition, true);
    window.removeEventListener('resize', updatePosition);
  }
});

onMounted(() => {
  document.addEventListener('click', handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
  window.removeEventListener('scroll', updatePosition, true);
  window.removeEventListener('resize', updatePosition);
});

defineExpose({ close });
</script>

<template>
  <div class="relative inline-block w-full">
    <button
      ref="triggerRef"
      type="button"
      :class="[
        'w-full bg-surface border border-border rounded-lg text-text focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition-all flex items-center justify-between',
        props.size === 'xs'
          ? 'px-2 py-1 text-xs'
          : props.size === 'sm'
            ? 'px-3 py-1.5 text-sm'
            : props.size === 'lg'
              ? 'px-5 py-3 text-base'
              : 'px-4 py-2.5 text-base',
      ]"
      @click="toggle"
    >
      <slot name="trigger" />
      <ChevronDown
        class="w-4 h-4 text-text-muted transition-transform duration-200"
        :style="{ transform: isOpen ? 'rotate(180deg)' : 'rotate(0)' }"
      />
    </button>

    <Teleport to="body">
      <Transition name="dropdown">
        <div
          v-if="isOpen"
          ref="dropdownRef"
          class="fixed z-[100] bg-surface border border-border rounded-lg shadow-xl overflow-hidden"
          :style="{
            top: `${coords.top}px`,
            left: `${coords.left}px`,
            width: `${coords.width}px`,
            visibility: coords.top === 0 ? 'hidden' : 'visible',
          }"
        >
          <slot :close="close" :model-value="modelValue" />
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
.dropdown-enter-active,
.dropdown-leave-active {
  transition:
    opacity 0.15s ease,
    transform 0.15s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
