<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick, watch } from 'vue';

interface Props {
  align?: 'left' | 'right';
}

const props = withDefaults(defineProps<Props>(), {
  align: 'right'
});

const isOpen = ref(false);
const triggerRef = ref<HTMLElement | null>(null);
const popoverRef = ref<HTMLElement | null>(null);
const coords = ref({ top: 0, left: 0 });

const updatePosition = async () => {
  if (!isOpen.value || !triggerRef.value) return;
  
  await nextTick();
  if (!popoverRef.value) return;

  const triggerRect = triggerRef.value.getBoundingClientRect();
  const popoverRect = popoverRef.value.getBoundingClientRect();
  const viewportWidth = window.innerWidth;
  const viewportHeight = window.innerHeight;

  let top = triggerRect.bottom + 8;
  let left = props.align === 'left' 
    ? triggerRect.left 
    : triggerRect.right - popoverRect.width;

  // Vertical check: if it goes below viewport, show it above the trigger
  if (top + popoverRect.height > viewportHeight - 12) {
    const spaceAbove = triggerRect.top;
    const spaceBelow = viewportHeight - triggerRect.bottom;
    
    if (spaceAbove > spaceBelow) {
      top = triggerRect.top - popoverRect.height - 8;
    }
  }

  // Horizontal check: ensure it stays within viewport
  if (left + popoverRect.width > viewportWidth - 12) {
    left = viewportWidth - popoverRect.width - 12;
  }
  if (left < 12) {
    left = 12;
  }

  coords.value = { top, left };
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
    popoverRef.value &&
    !triggerRef.value.contains(event.target as Node) &&
    !popoverRef.value.contains(event.target as Node)
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
            top: `${coords.top}px`,
            left: `${coords.left}px`,
            visibility: coords.top === 0 ? 'hidden' : 'visible'
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
