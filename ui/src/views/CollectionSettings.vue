<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import { 
  FolderOpen, 
  Plus,
  Trash2,
  Settings,
  Save
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

onMounted(() => {
  fetchCollections();
  fetchCollection();
});
</script>

<template>
  <AppLayout>
    
      <!-- Header -->
      <AppHeader>
        <template #breadcrumb>
          <div class="flex items-center text-sm text-text-muted overflow-hidden whitespace-nowrap">
            <span class="hover:text-text cursor-pointer shrink-0" @click="router.push('/')">Vault</span>
            <span class="mx-2 shrink-0">/</span>
            <span class="hover:text-text cursor-pointer shrink-0 hidden sm:inline" @click="router.push('/collections')">Collections</span>
            <span class="mx-2 shrink-0 hidden sm:inline">/</span>
            <span class="hover:text-text cursor-pointer truncate" @click="router.push(`/collections/${collectionName}`)">{{ collectionName }}</span>
            <span class="mx-2 shrink-0">/</span>
            <span class="font-medium text-text shrink-0">Settings</span>
          </div>
        </template>
      </AppHeader>

      <!-- Form Content -->
      <div class="flex-1 overflow-auto min-h-0 p-4 sm:p-8 pb-24 sm:pb-8">
        <div class="max-w-4xl mx-auto space-y-6">
          <div>
            <h1 class="text-2xl font-bold text-text">Collection Settings</h1>
            <p class="text-sm text-text-muted mt-1">Manage fields for {{ collectionName }}</p>
          </div>

          <form @submit.prevent="saveSettings" class="space-y-6">
            <div class="bg-surface-dark border border-border rounded-lg p-4 sm:p-6">
              <div class="flex items-center justify-between mb-4">
                <h2 class="text-lg font-semibold text-text flex items-center gap-2">
                  <Settings class="w-5 h-5 text-primary" />
                  Fields Definition
                </h2>
                <Button 
                  @click="addField" 
                  variant="secondary"
                  size="sm"
                >
                  <Plus class="w-4 h-4" />
                  <span class="hidden sm:inline">Add Field</span>
                  <span class="sm:hidden">Add</span>
                </Button>
              </div>

              <div class="space-y-3">
                <div 
                  v-for="(field, index) in fields" 
                  :key="index" 
                  class="flex flex-col sm:flex-row items-start sm:items-center gap-3 bg-surface p-4 rounded-lg border border-border"
                >
                  <div class="w-full sm:flex-1">
                    <input 
                      v-model="field.name" 
                      placeholder="field_name" 
                      class="w-full bg-surface-dark border border-border rounded-lg px-3 py-2 text-sm text-text placeholder-text-muted focus:outline-none focus:ring-1 focus:ring-primary/50 focus:border-primary transition-all"
                    />
                  </div>
                  
                  <div class="w-full sm:w-40">
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

                  <div class="flex items-center justify-between w-full sm:w-auto gap-4">
                    <label class="flex items-center gap-2 text-sm text-text-muted cursor-pointer">
                      <input 
                        v-model="field.required" 
                        type="checkbox" 
                        class="w-4 h-4 text-primary bg-surface-dark border-border rounded focus:ring-primary"
                      />
                      Required
                    </label>
                    
                    <Button 
                      @click="removeField(index)" 
                      variant="ghost"
                      size="xs"
                      class="!text-error hover:!bg-error/10"
                      :disabled="fields.length === 1"
                    >
                      <Trash2 class="w-4 h-4" />
                    </Button>
                  </div>
                </div>
              </div>
            </div>

            <div class="flex items-center justify-end gap-3">
              <Button 
                @click="router.push(`/collections/${collectionName}`)" 
                variant="secondary"
                class="flex-1 sm:flex-none"
              >
                Cancel
              </Button>
              <Button 
                type="submit" 
                class="flex-1 sm:flex-none"
              >
                <Save class="w-4 h-4" />
                Save Changes
              </Button>
            </div>
          </form>
        </div>
      </div>
      </AppLayout>
</template>


