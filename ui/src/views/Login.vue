<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import Button from '../components/Button.vue';
import Input from '../components/Input.vue';
import Checkbox from '../components/Checkbox.vue';
import { Mail, KeyRound, AlertCircle, LockKeyhole } from 'lucide-vue-next';

const identity = ref('');
const password = ref('');
const rememberMe = ref(false);
const error = ref('');
const systemStatus = ref<'checking' | 'online' | 'offline'>('checking');
const backendPort = ref('');
const auth = useAuthStore();
const router = useRouter();

const handleLogin = async () => {
  error.value = '';
  const success = await auth.login(identity.value, password.value);
  if (success) {
    if (rememberMe.value) {
      const expiryDate = new Date();
      expiryDate.setDate(expiryDate.getDate() + 7);
      localStorage.setItem('rememberMe', 'true');
      localStorage.setItem('rememberMeExpiry', expiryDate.toISOString());
    }
    router.push({ name: 'Dashboard' });
  } else {
    error.value = 'Invalid identity or password';
  }
};

const checkSystemStatus = async () => {
  try {
    const response = await fetch('/api/health', { method: 'HEAD' });
    if (response.ok) {
      systemStatus.value = 'online';
      // Try to get port from response header or use window location
      backendPort.value =
        response.headers.get('X-Server-Port') || import.meta.env.VITE_API_PORT || '8090';
    } else {
      systemStatus.value = 'offline';
    }
  } catch {
    systemStatus.value = 'offline';
  }
};

onMounted(() => {
  checkSystemStatus();
  // Check if remember me is still valid
  const rememberMeExpiry = localStorage.getItem('rememberMeExpiry');
  if (rememberMeExpiry && new Date(rememberMeExpiry) > new Date()) {
    rememberMe.value = true;
  } else {
    localStorage.removeItem('rememberMe');
    localStorage.removeItem('rememberMeExpiry');
  }
});
</script>

<template>
  <div
    class="h-[100dvh] flex flex-col items-center justify-center bg-background p-4 sm:p-6 relative"
  >
    <!-- Background Pattern -->
    <div
      class="fixed inset-0 z-0 pointer-events-none opacity-10"
      style="
        background-image: radial-gradient(var(--color-border) 1px, transparent 1px);
        background-size: 24px 24px;
      "
    ></div>

    <main class="w-full max-w-md relative z-10">
      <!-- Brand Header -->
      <div class="text-center mb-8">
        <div
          class="inline-flex items-center justify-center w-16 h-16 rounded-xl bg-surface border border-border shadow-lg mb-4 group transition-transform duration-300 hover:scale-105"
        >
          <LockKeyhole class="text-primary w-8 h-8 group-hover:text-white transition-colors" />
        </div>
        <h1 class="text-3xl font-bold tracking-tight text-text mb-1">Vault Admin</h1>
        <p class="text-sm text-text-muted">Backend Management System</p>
      </div>

      <!-- Login Card -->
      <div class="bg-surface border border-border shadow-2xl rounded-xl overflow-hidden">
        <!-- Top Accent Line -->
        <div class="h-1 w-full bg-gradient-to-r from-primary via-accent to-secondary"></div>

        <div class="p-8">
          <form class="space-y-6" @submit.prevent="handleLogin">
            <!-- Email Input -->
            <div class="space-y-2">
              <label for="email" class="block text-sm font-medium text-text">Email Address</label>
              <div class="relative group">
                <Mail
                  class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-text-muted group-focus-within:text-primary transition-colors pointer-events-none"
                />
                <Input
                  id="email"
                  v-model="identity"
                  type="text"
                  required
                  placeholder="admin@company.com"
                  class="!pl-10"
                />
              </div>
            </div>

            <!-- Password Input -->
            <div class="space-y-2">
              <div class="flex items-center justify-between">
                <label for="password" class="block text-sm font-medium text-text">Password</label>
                <a
                  href="#"
                  class="text-xs font-medium text-primary hover:text-primary-hover hover:underline transition-colors"
                  >Forgot password?</a
                >
              </div>
              <div class="relative group">
                <KeyRound
                  class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-text-muted group-focus-within:text-primary transition-colors pointer-events-none"
                />
                <Input
                  id="password"
                  v-model="password"
                  type="password"
                  required
                  placeholder="••••••••"
                  class="!pl-10"
                />
              </div>
            </div>

            <!-- Remember Me -->
            <Checkbox v-model="rememberMe" label="Keep me signed in for 7 days" />

            <!-- Error Message -->
            <div
              v-if="error"
              class="flex items-center space-x-2 bg-error/10 border border-error/20 text-error p-3 rounded-lg text-sm"
            >
              <AlertCircle class="w-4 h-4" />
              <span>{{ error }}</span>
            </div>

            <!-- Submit Button -->
            <Button type="submit" class="w-full"> Authenticate Access </Button>
          </form>
        </div>

        <!-- Bottom Status -->
        <div
          class="px-8 py-4 bg-surface-dark border-t border-border flex justify-between items-center"
        >
          <div class="flex items-center space-x-2">
            <span v-if="systemStatus === 'checking'" class="relative flex h-2.5 w-2.5">
              <span
                class="animate-ping absolute inline-flex h-full w-full rounded-full bg-warning opacity-75"
              ></span>
              <span class="relative inline-flex rounded-full h-2.5 w-2.5 bg-warning"></span>
            </span>
            <span v-else-if="systemStatus === 'online'" class="relative flex h-2.5 w-2.5">
              <span
                class="animate-ping absolute inline-flex h-full w-full rounded-full bg-success opacity-75"
              ></span>
              <span class="relative inline-flex rounded-full h-2.5 w-2.5 bg-success"></span>
            </span>
            <span v-else class="relative flex h-2.5 w-2.5 bg-error rounded-full"></span>
            <span class="text-xs text-text-muted uppercase tracking-wide">
              {{
                systemStatus === 'checking'
                  ? 'Checking...'
                  : systemStatus === 'online'
                    ? `System Operational (Port: ${backendPort})`
                    : 'System Offline'
              }}
            </span>
          </div>
          <a href="#" class="text-xs text-text-muted hover:text-primary transition-colors">Help</a>
        </div>
      </div>

      <!-- Footer -->
      <div class="mt-8 text-center">
        <div class="flex justify-center space-x-4 text-xs text-text-dim">
          <a href="#" class="hover:text-primary transition-colors">Privacy Policy</a>
          <a href="#" class="hover:text-primary transition-colors">Terms of Service</a>
        </div>
      </div>
    </main>
  </div>
</template>
