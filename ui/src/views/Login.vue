<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import { Lock, User, AlertCircle } from 'lucide-vue-next';

const identity = ref('');
const password = ref('');
const error = ref('');
const auth = useAuthStore();
const router = useRouter();

const handleLogin = async () => {
  error.value = '';
  const success = await auth.login(identity.value, password.value);
  if (success) {
    router.push({ name: 'Dashboard' });
  } else {
    error.value = 'Invalid identity or password';
  }
};
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-monokai-bg p-6">
    <div class="max-w-md w-full">
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold text-white tracking-tight">Vault</h1>
        <p class="text-monokai-gray mt-2">Administrative Interface</p>
      </div>

      <div class="bg-monokai-panel p-8 rounded-2xl border border-white/10 shadow-2xl">
        <form @submit.prevent="handleLogin" class="space-y-5">
          <div>
            <label class="block text-xs font-semibold text-monokai-gray uppercase tracking-wider mb-2">Identity</label>
            <div class="relative">
              <User class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-monokai-gray" />
              <input 
                v-model="identity" 
                type="text" 
                required 
                placeholder="Email or username"
                class="w-full bg-monokai-header border border-monokai-gray/30 rounded-xl py-3 pl-10 pr-4 text-white focus:border-monokai-blue focus:ring-1 focus:ring-monokai-blue/50 outline-none transition-all" 
              />
            </div>
          </div>
          
          <div>
            <label class="block text-xs font-semibold text-monokai-gray uppercase tracking-wider mb-2">Password</label>
            <div class="relative">
              <Lock class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-monokai-gray" />
              <input 
                v-model="password" 
                type="password" 
                required 
                placeholder="••••••••"
                class="w-full bg-monokai-header border border-monokai-gray/30 rounded-xl py-3 pl-10 pr-4 text-white focus:border-monokai-blue focus:ring-1 focus:ring-monokai-blue/50 outline-none transition-all" 
              />
            </div>
          </div>
          
          <div v-if="error" class="flex items-center space-x-2 bg-monokai-pink/10 border border-monokai-pink/20 text-monokai-pink p-3 rounded-xl text-sm">
            <AlertCircle class="w-4 h-4" />
            <span>{{ error }}</span>
          </div>
          
          <button 
            type="submit" 
            class="w-full py-3 bg-monokai-blue hover:bg-monokai-blue/90 text-monokai-bg font-bold rounded-xl transition-all shadow-lg shadow-monokai-blue/20 active:scale-[0.98]"
          >
            Sign In
          </button>
        </form>
      </div>
      
      <p class="text-center text-monokai-gray text-xs mt-8 uppercase tracking-widest opacity-50">
        Vault Framework v0.1.0
      </p>
    </div>
  </div>
</template>
