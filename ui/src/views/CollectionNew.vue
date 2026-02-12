<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import { 
  Plus, 
  Trash2, 
  Save, 
  X, 
  FolderPlus 
  FolderOpen, 
  Terminal, 
  Settings, 
  Cloud} from 'lucide-vue-next';

const router = useRouter();

const collectionFormData = ref({
  name: '',
  type: 'base',
  fields: [{ name: 'name', type: 'text', required: true }]
});

const addField = () => {
  collectionFormData.value.fields.push({ name: '', type: 'text', required: false });
};

const removeField = (index: number) => {
  collectionFormData.value.fields.splice(index, 1);
};

const saveCollection = async () => {
  try {
    await axios.post('/api/admin/collections', collectionFormData.value);
    router.push('/');
  } catch (error) {
    console.error('Collection creation failed', error);
  }
};

const handleLogout = () => {
  auth.logout();
  router.push({ name: 'Login' });
};
</script>

<template>
  <AppLayout>
    
      <!-- Header -->
      <header class="h-16 flex items-center justify-between px-8 border-b border-border bg-surface z-10">
        <div class="flex items-center text-sm text-text-muted">
          <span class="hover:text-text cursor-pointer" @click="router.push('/')">Vault</span>
          <span class="mx-2">/</span>
          <span class="hover:text-text cursor-pointer" @click="router.push('/')">Collections</span>
          <span class="mx-2">/</span>
          <span class="font-medium text-text">New</span>
        </div>
      </header>

      <!-- Form Content -->
      <div class="flex-1 overflow-auto p-8">
        <div class="space-y-6">
          <div>
            <h1 class="text-2xl font-bold text-text">Create New Collection</h1>
            <p class="text-sm text-text-muted mt-1">Define your database schema and fields</p>
          </div>

          <form @submit.prevent="saveCollection" class="space-y-6">
            <!-- Basic Info Card -->
            <div class="bg-surface-dark border border-border rounded-lg p-6">
              <h2 class="text-lg font-semibold text-text mb-4 flex items-center gap-2">
                <FolderPlus class="w-5 h-5 text-primary" />
                Basic Information
              </h2>
              
              <div class="grid grid-cols-2 gap-6">
                <div>
                  <label class="block text-sm font-medium text-text mb-2">Collection Name</label>
                  <input 
                    v-model="collectionFormData.name" 
                    type="text" 
                    required 
                    placeholder="e.g. products, users, posts"
                    class="w-full bg-surface border border-border rounded-lg px-4 py-2.5 text-text placeholder-text-muted focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition-all"
                  />
                  <p class="text-xs text-text-dim mt-1">Lowercase, no spaces (use underscores)</p>
                </div>
                
                <div>
                  <label class="block text-sm font-medium text-text mb-2">Collection Type</label>
                  <select 
                    v-model="collectionFormData.type" 
                    class="w-full bg-surface border border-border rounded-lg px-4 py-2.5 text-text focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition-all"
                  >
                    <option value="base">Base (Generic Data)</option>
                    <option value="auth">Auth (User Records)</option>
                  </select>
                  <p class="text-xs text-text-dim mt-1">Choose the collection purpose</p>
                </div>
              </div>
            </div>

            <!-- Fields Card -->
            <div class="bg-surface-dark border border-border rounded-lg p-6">
              <div class="flex items-center justify-between mb-4">
                <h2 class="text-lg font-semibold text-text flex items-center gap-2">
                  <Plus class="w-5 h-5 text-primary" />
                  Fields Definition
                </h2>
                <button 
                  type="button" 
                  @click="addField" 
                  class="px-4 py-2 bg-primary/10 text-primary rounded-lg text-sm font-medium hover:bg-primary/20 transition-colors flex items-center gap-2"
                >
                  <Plus class="w-4 h-4" />
                  Add Field
                </button>
              </div>

              <div class="space-y-3">
                <div 
                  v-for="(field, index) in collectionFormData.fields" 
                  :key="index" 
                  class="flex items-center gap-3 bg-surface p-4 rounded-lg border border-border"
                >
                  <div class="flex-1">
                    <input 
                      v-model="field.name" 
                      placeholder="field_name" 
                      class="w-full bg-surface-dark border border-border rounded-lg px-3 py-2 text-sm text-text placeholder-text-muted focus:outline-none focus:ring-1 focus:ring-primary/50 focus:border-primary transition-all"
                    />
                  </div>
                  
                  <div class="w-40">
                    <select 
                      v-model="field.type" 
                      class="w-full bg-surface-dark border border-border rounded-lg px-3 py-2 text-sm text-text focus:outline-none focus:ring-1 focus:ring-primary/50 focus:border-primary transition-all"
                    >
                      <option value="text">Text</option>
                      <option value="number">Number</option>
                      <option value="bool">Boolean</option>
                      <option value="json">JSON</option>
                      <option value="file">File</option>
                    </select>
                  </div>

                  <label class="flex items-center gap-2 text-sm text-text-muted cursor-pointer">
                    <input 
                      v-model="field.required" 
                      type="checkbox" 
                      class="w-4 h-4 text-primary bg-surface-dark border-border rounded focus:ring-primary"
                    />
                    Required
                  </label>
                  
                  <button 
                    type="button" 
                    @click="removeField(index)" 
                    class="p-2 text-error hover:bg-error/10 rounded-lg transition-colors"
                    :disabled="collectionFormData.fields.length === 1"
                  >
                    <Trash2 class="w-4 h-4" />
                  </button>
                </div>
              </div>
            </div>

            <!-- Actions -->
            <div class="flex items-center justify-end gap-3">
              <button 
                type="button" 
                @click="router.push('/')" 
                class="px-6 py-2.5 bg-surface border border-border text-text rounded-lg font-medium hover:bg-surface-dark transition-colors"
              >
                Cancel
              </button>
              <button 
                type="submit" 
                class="px-6 py-2.5 bg-primary hover:bg-primary-hover text-white rounded-lg font-medium transition-colors flex items-center gap-2"
              >
                <Save class="w-4 h-4" />
                Create Collection
              </button>
            </div>
          </form>
        </div>
      </div>
      </AppLayout>
</template>


