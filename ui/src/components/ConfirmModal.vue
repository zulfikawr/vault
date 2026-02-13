<script setup lang="ts">
import { X, AlertTriangle } from 'lucide-vue-next';
import Button from './Button.vue';

interface Props {
  show: boolean;
  title: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  variant?: 'danger' | 'warning' | 'info';
}

const props = withDefaults(defineProps<Props>(), {
  confirmText: 'Confirm',
  cancelText: 'Cancel',
  variant: 'info'
});

const emit = defineEmits<{
  confirm: [];
  cancel: [];
}>();
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <!-- Backdrop -->
        <div 
          class="absolute inset-0 bg-black/70"
          @click="emit('cancel')"
        ></div>
        
        <!-- Modal -->
        <div class="relative bg-surface border border-border rounded-lg shadow-2xl max-w-lg w-full">
          <!-- Header -->
          <div class="flex items-start justify-between p-6 border-b border-border">
            <div class="flex items-center gap-3">
              <div 
                :class="[
                  'p-2 rounded-lg',
                  variant === 'danger' ? 'bg-error/10 text-error' : '',
                  variant === 'warning' ? 'bg-warning/10 text-warning' : '',
                  variant === 'info' ? 'bg-primary/10 text-primary' : ''
                ]"
              >
                <AlertTriangle class="w-5 h-5" />
              </div>
              <h3 class="text-lg font-semibold text-text">{{ title }}</h3>
            </div>
            <Button 
              @click="emit('cancel')"
              variant="ghost"
              size="xs"
              class="!p-1"
            >
              <X class="w-5 h-5" />
            </Button>
          </div>
          
          <!-- Body -->
          <div class="p-6">
            <p class="text-text-muted">{{ message }}</p>
          </div>
          
          <!-- Footer -->
          <div class="flex justify-end gap-3 p-6 border-t border-border bg-surface-dark/50">
            <Button 
              @click="emit('cancel')"
              variant="secondary"
            >
              {{ cancelText }}
            </Button>
            <Button 
              @click="emit('confirm')"
              :variant="variant === 'danger' ? 'primary' : 'primary'"
              :class="variant === 'danger' ? '!bg-error hover:!bg-error/90' : ''"
            >
              {{ confirmText }}
            </Button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-active .relative,
.modal-leave-active .relative {
  transition: transform 0.2s ease, opacity 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .relative,
.modal-leave-to .relative {
  transform: scale(0.95);
  opacity: 0;
}
</style>
