<script setup lang="ts">
import { X } from 'lucide-vue-next';
import { onMounted, onUnmounted, watch } from 'vue';
import Button from './Button.vue';

interface Props {
  show: boolean;
  title: string;
  maxWidth?: 'sm' | 'md' | 'lg' | 'xl' | '2xl' | '4xl';
}

const props = withDefaults(defineProps<Props>(), {
  maxWidth: 'md',
});

const emit = defineEmits<{
  close: [];
}>();

const handleEscape = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && props.show) {
    emit('close');
  }
};

onMounted(() => {
  window.addEventListener('keydown', handleEscape);
});

onUnmounted(() => {
  window.removeEventListener('keydown', handleEscape);
});

watch(() => props.show, (newVal) => {
  if (newVal) {
    document.body.style.overflow = 'hidden';
  } else {
    document.body.style.overflow = '';
  }
});
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="show" class="fixed inset-0 z-[100] flex items-center justify-center p-4 sm:p-6 overflow-y-auto">
        <!-- Backdrop -->
        <div class="fixed inset-0 bg-background/80 backdrop-blur-md transition-opacity duration-300" @click="emit('close')"></div>

        <!-- Modal Content -->
        <div
          class="relative bg-surface/90 border border-border shadow-[0_25px_50px_-12px_rgba(0,0,0,0.5)] rounded-3xl w-full transition-all duration-300 font-mono"
          :class="[
            maxWidth === 'sm' ? 'max-w-sm' : '',
            maxWidth === 'md' ? 'max-w-md' : '',
            maxWidth === 'lg' ? 'max-w-lg' : '',
            maxWidth === 'xl' ? 'max-w-xl' : '',
            maxWidth === '2xl' ? 'max-w-2xl' : '',
            maxWidth === '4xl' ? 'max-w-4xl' : '',
          ]"
        >
          <!-- Terminal Header -->
          <div class="flex items-center justify-between px-6 py-4 border-b border-border/50">
            <div class="flex items-center gap-2">
              <div class="flex gap-1.5 mr-3 hidden sm:flex">
                <div class="w-3 h-3 rounded-full bg-error/30 border border-error/50"></div>
                <div class="w-3 h-3 rounded-full bg-warning/30 border border-warning/50"></div>
                <div class="w-3 h-3 rounded-full bg-success/30 border border-success/50"></div>
              </div>
              <h3 class="text-xs font-black uppercase tracking-widest text-text leading-none">{{ title }}</h3>
            </div>
            
            <button 
              class="w-10 h-10 flex items-center justify-center rounded-xl bg-surface-dark border border-border text-text-muted hover:text-error hover:border-error/30 transition-all active:scale-95"
              @click="emit('close')"
              aria-label="Close modal"
            >
              <X class="w-5 h-5" />
            </button>
          </div>

          <!-- Body -->
          <div class="p-6 sm:p-8 custom-scrollbar max-h-[70vh] overflow-y-auto">
            <slot />
          </div>

          <!-- Footer -->
          <div v-if="$slots.footer" class="flex flex-col sm:flex-row justify-end gap-3 px-6 py-5 border-t border-border/50 bg-surface-dark/40 rounded-b-3xl">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}

.modal-enter-active .relative,
.modal-leave-active .relative {
  transition: transform 0.5s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.4s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .relative,
.modal-leave-to .relative {
  transform: scale(0.92) translateY(20px);
  opacity: 0;
}

.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: var(--border);
  border-radius: 10px;
}
</style>
