<script setup lang="ts">
import Input from './Input.vue';
import { Search } from 'lucide-vue-next';
</script>

<template>
  <header
    class="flex h-16 items-center justify-between px-4 sm:px-8 border-b border-border bg-surface z-10 shrink-0"
  >
    <div class="flex-1 min-w-0">
      <slot name="breadcrumb" />
    </div>

    <div class="relative hidden sm:block group ml-4">
      <Input
        type="text"
        placeholder="Search (Ctrl+K)"
        class="w-64 !bg-surface-dark pl-9 h-9 text-sm"
      />
      <Search
        class="absolute left-2.5 top-1/2 -translate-y-1/2 text-text-muted w-4 h-4 pointer-events-none"
      />
    </div>
  </header>
</template>
