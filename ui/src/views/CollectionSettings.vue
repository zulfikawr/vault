<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import { 
  FolderOpen, 
  Terminal, 
  Settings, 
  CloudSave,
  Plus,
  Trash2
} from 'lucide-vue-next';

const router = useRouter();
const route = useRoute();
const collections = ref([]);
const collection = ref(null);
const fields = ref([]);

const collectionName = computed(() => route.params.name as string);

const fetchCollections = async () => {
  try {
    const response = await axios.get('/api/admin/collections');
    collections.value = response.data.data;
  } catch (error) {
    console.error('Failed to fetch collections', error);
  }
};

const fetchCollection = async () => {
  try {
    const response = await axios.get(`/api/admin/collections`);
    const col = response.data.data.find((c: any) => c.name === collectionName.value);
    collection.value = col;
    fields.value = JSON.parse(JSON.stringify(col?.fields || []));
  } catch (error) {
    console.error('Failed to fetch collection', error);
  }
};

const addField = () => {
  fields.value.push({ name: '', type: 'text', required: false });
};

const removeField = (index: number) => {
  fields.value.splice(index, 1);
};

const saveSettings = async () => {
  try {
    await axios.patch(`/api/admin/collections/${collection.value.id}`, {
      fields: fields.value
    });
    router.push(`/collections/${collectionName.value}`);
  } catch (error) {
    console.error('Save failed', error);
  }
};

const handleLogout = () => {
  auth.logout();
  router.push({ name: 'Login' });
};

onMounted(() => {
  fetchCollections();
  fetchCollection();
});
</script>

<template>
  <AppLayout>
    
      <!-- Header -->
      <header class="h-16 flex items-center justify-between px-8 border-b border-border bg-surface z-10">
        <div class="flex items-center text-sm text-text-muted">
          <span class="hover:text-text cursor-pointer" @click="router.push('/')">Vault</span>
          <span class="mx-2">/</span>
          <span class="hover:text-text cursor-pointer" @click="router.push('/collections')">Collections</span>
          <span class="mx-2">/</span>
          <span class="hover:text-text cursor-pointer" @click="router.push(`/collections/${collectionName}`)">{{ collectionName }}</span>
          <span class="mx-2">/</span>
          <span class="font-medium text-text">Settings</span>
        </div>
      </header>

      <!-- Form Content -->
      <div class="flex-1 overflow-auto p-8">
        <div class="space-y-6">
          <div>
            <h1 class="text-2xl font-bold text-text">Collection Settings</h1>
            <p class="text-sm text-text-muted mt-1">Manage fields for {{ collectionName }}</p>
          </div>

          <form @submit.prevent="saveSettings" class="space-y-6">
            <div class="bg-surface-dark border border-border rounded-lg p-6">
              <div class="flex items-center justify-between mb-4">
                <h2 class="text-lg font-semibold text-text flex items-center gap-2">
                  <Settings class="w-5 h-5 text-primary" />
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
                  v-for="(field, index) in fields" 
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
                    :disabled="fields.length === 1"
                  >
                    <Trash2 class="w-4 h-4" />
                  </button>
                </div>
              </div>
            </div>

            <div class="flex items-center justify-end gap-3">
              <button 
                type="button" 
                @click="router.push(`/collections/${collectionName}`)" 
                class="px-6 py-2.5 bg-surface border border-border text-text rounded-lg font-medium hover:bg-surface-dark transition-colors"
              >
                Cancel
              </button>
              <button 
                type="submit" 
                class="px-6 py-2.5 bg-primary hover:bg-primary-hover text-white rounded-lg font-medium transition-colors flex items-center gap-2"
              >
                <Save class="w-4 h-4" />
                Save Changes
              </button>
            </div>
          </form>
        </div>
      </div>
      </AppLayout>
</template>

