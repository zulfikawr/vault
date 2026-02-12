<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue';

interface Props {
  align?: 'left' | 'right';
}

const props = withDefaults(defineProps<Props>(), {
  align: 'right'
});

const isOpen = ref(false);
const triggerRef = ref<HTMLElement | null>(null);
const popoverRef = ref<HTMLElement | null>(null);
const position = ref({ top: 0, left: 0, right: 0 });

const toggle = (event: Event) => {
  event.stopPropagation();
  isOpen.value = !isOpen.value;
  
  if (isOpen.value && triggerRef.value) {
    const rect = triggerRef.value.getBoundingClientRect();
    position.value = {
      top: rect.bottom + window.scrollY + 4,
      left: rect.left + window.scrollX,
      right: window.innerWidth - rect.right - window.scrollX
    };
  }
};

const close = () => {
  isOpen.value = false;
};

const handleClickOutside = (event: MouseEvent) => {
  if (
    triggerRef.value &&
    popoverRef.value &&
    !triggerRef.value.contains(event.target as Node) &&
    !popoverRef.value.contains(event.target as Node)
  ) {
    close();
  }
};

onMounted(() => {
  document.addEventListener('click', handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
});

defineExpose({ close });
</script>

<template>
  <div class="relative inline-block">
    <div ref="triggerRef" @click="toggle">
      <slot name="trigger" />
    </div>
    
    <Teleport to="body">
      <Transition name="popover">
        <div
          v-if="isOpen"
          ref="popoverRef"
          class="fixed z-[100] min-w-[180px] bg-surface border border-border rounded-lg shadow-xl overflow-hidden"
          :style="{
            top: `${position.top}px`,
            [align === 'right' ? 'right' : 'left']: align === 'right' ? `${position.right}px` : `${position.left}px`
          }"
        >
          <slot :close="close" />
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style scoped>
.popover-enter-active,
.popover-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.popover-enter-from,
.popover-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
