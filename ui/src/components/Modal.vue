<script setup lang="ts">
import { X } from 'lucide-vue-next';
import Button from './Button.vue';

interface Props {
  show: boolean;
  title: string;
  maxWidth?: 'sm' | 'md' | 'lg' | 'xl' | '2xl';
}

withDefaults(defineProps<Props>(), {
  maxWidth: 'md',
});

const emit = defineEmits<{
  close: [];
}>();
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center p-4 overflow-y-auto">
        <!-- Backdrop -->
        <div class="fixed inset-0 bg-black/70 transition-opacity" @click="emit('close')"></div>

        <!-- Modal Content -->
        <div
          class="relative bg-surface border border-border rounded-lg shadow-2xl w-full transition-all"
          :class="[
            maxWidth === 'sm' ? 'max-w-sm' : '',
            maxWidth === 'md' ? 'max-w-md' : '',
            maxWidth === 'lg' ? 'max-w-lg' : '',
            maxWidth === 'xl' ? 'max-w-xl' : '',
            maxWidth === '2xl' ? 'max-w-2xl' : '',
          ]"
        >
          <!-- Header -->
          <div class="flex items-center justify-between p-6 border-b border-border">
            <h3 class="text-lg font-semibold text-text tracking-tight">{{ title }}</h3>
            <Button variant="ghost" size="xs" class="!p-1 text-text-muted hover:text-text" @click="emit('close')">
              <X class="w-5 h-5" />
            </Button>
          </div>

          <!-- Body -->
          <div class="p-6">
            <slot />
          </div>

          <!-- Footer -->
          <div v-if="$slots.footer" class="flex justify-end gap-3 p-6 border-t border-border bg-surface-dark/30">
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
  transition: opacity 0.3s ease;
}

.modal-enter-active .relative,
.modal-leave-active .relative {
  transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1), opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .relative,
.modal-leave-to .relative {
  transform: scale(0.95) translateY(10px);
  opacity: 0;
}
</style>
