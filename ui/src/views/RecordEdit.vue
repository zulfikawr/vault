<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import { 
  FolderOpen, 
  X,
  Save
} from 'lucide-vue-next';

const router = useRouter();
const route = useRoute();
const collection = ref(null);
const formData = ref({});
const loading = ref(true);

const collectionName = computed(() => route.params.name as string);
const recordId = computed(() => route.params.id as string);

const fetchCollection = async () => {
  try {
    const response = await axios.get(`/api/admin/collections`);
    const col = response.data.data.find((c: any) => c.name === collectionName.value);
    collection.value = col;
  } catch (error) {
    console.error('Failed to fetch collection', error);
  }
};

const fetchRecord = async () => {
  try {
    const response = await axios.get(`/api/collections/${collectionName.value}/records/${recordId.value}`);
    formData.value = response.data.data || {};
    loading.value = false;
  } catch (error) {
    console.error('Failed to fetch record', error);
    loading.value = false;
  }
};

const handleSubmit = async () => {
  try {
    await axios.patch(`/api/collections/${collectionName.value}/records/${recordId.value}`, formData.value);
    router.push(`/collections/${collectionName.value}`);
  } catch (error) {
    console.error('Update failed', error);
    alert('Failed to update record');
  }
};

onMounted(() => {
  fetchCollection();
  fetchRecord();
});
</script>

<template>
  <AppLayout>
    <AppHeader>
      <template #breadcrumb>
        <div class="flex items-center text-sm text-text-muted overflow-hidden whitespace-nowrap">
          <span class="hover:text-text cursor-pointer shrink-0" @click="router.push('/')">Vault</span>
          <span class="mx-2 shrink-0">/</span>
          <span class="hover:text-text cursor-pointer shrink-0 hidden sm:inline" @click="router.push('/collections')">Collections</span>
          <span class="mx-2 shrink-0 hidden sm:inline">/</span>
          <span class="hover:text-text cursor-pointer truncate" @click="router.push(`/collections/${collectionName}`)">{{ collectionName }}</span>
          <span class="mx-2 shrink-0">/</span>
          <span class="font-medium text-text shrink-0">Edit</span>
        </div>
      </template>
    </AppHeader>

    <div class="flex-1 overflow-auto min-h-0 p-4 sm:p-8 pb-24 sm:pb-8">
      <div class="max-w-4xl mx-auto space-y-6">
        <div>
          <h1 class="text-2xl font-bold text-text mb-2">Edit Record</h1>
          <p class="text-text-muted">Update record in {{ collectionName }} collection</p>
        </div>

        <form v-if="!loading && collection" @submit.prevent="handleSubmit" class="bg-surface-dark rounded-lg border border-border p-4 sm:p-6 space-y-6">
          <div class="space-y-4">
            <div v-for="field in collection.fields" :key="field.name">
              <label :for="field.name" class="block text-sm font-medium text-text mb-1.5">
                {{ field.name }}
                <span v-if="field.required" class="text-error">*</span>
              </label>
              
              <input
                v-if="field.type === 'text'"
                :id="field.name"
                v-model="formData[field.name]"
                type="text"
                :required="field.required"
                class="w-full px-3 py-2 bg-surface border border-border rounded-lg text-text placeholder-text-muted focus:ring-1 focus:ring-primary focus:border-primary"
              />
              
              <input
                v-else-if="field.type === 'number'"
                :id="field.name"
                v-model.number="formData[field.name]"
                type="number"
                :required="field.required"
                class="w-full px-3 py-2 bg-surface border border-border rounded-lg text-text placeholder-text-muted focus:ring-1 focus:ring-primary focus:border-primary"
              />
              
              <input
                v-else-if="field.type === 'email'"
                :id="field.name"
                v-model="formData[field.name]"
                type="email"
                :required="field.required"
                class="w-full px-3 py-2 bg-surface border border-border rounded-lg text-text placeholder-text-muted focus:ring-1 focus:ring-primary focus:border-primary"
              />
              
              <input
                v-else-if="field.type === 'date'"
                :id="field.name"
                v-model="formData[field.name]"
                type="date"
                :required="field.required"
                class="w-full px-3 py-2 bg-surface border border-border rounded-lg text-text placeholder-text-muted focus:ring-1 focus:ring-primary focus:border-primary"
              />
              
              <div v-else-if="field.type === 'bool'" class="flex items-center">
                <input
                  :id="field.name"
                  v-model="formData[field.name]"
                  type="checkbox"
                  class="w-4 h-4 text-primary bg-surface border-border rounded focus:ring-primary focus:ring-2"
                />
                <label :for="field.name" class="ml-2 text-sm text-text-muted">Enable</label>
              </div>
              
              <textarea
                v-else-if="field.type === 'json'"
                :id="field.name"
                v-model="formData[field.name]"
                rows="4"
                :required="field.required"
                class="w-full px-3 py-2 bg-surface border border-border rounded-lg text-text placeholder-text-muted focus:ring-1 focus:ring-primary focus:border-primary font-mono text-sm"
              ></textarea>
              
              <input
                v-else
                :id="field.name"
                v-model="formData[field.name]"
                type="text"
                :required="field.required"
                class="w-full px-3 py-2 bg-surface border border-border rounded-lg text-text placeholder-text-muted focus:ring-1 focus:ring-primary focus:border-primary"
              />
            </div>
          </div>

          <div class="flex flex-col sm:flex-row justify-end gap-3 mt-6 pt-6 border-t border-border">
            <button 
              type="button"
              @click="router.push(`/collections/${collectionName}`)"
              class="px-6 py-2.5 bg-surface border border-border rounded-lg font-medium text-text hover:bg-surface-dark transition-colors flex items-center justify-center gap-2"
            >
              <X class="w-4 h-4" />
              Cancel
            </button>
            <button 
              type="submit" 
              class="px-6 py-2.5 bg-primary hover:bg-primary-hover text-white rounded-lg font-medium transition-colors flex items-center justify-center gap-2"
            >
              <Save class="w-4 h-4" />
              Update Record
            </button>
          </div>
        </form>

        <div v-else class="bg-surface-dark rounded-lg border border-border p-12 text-center">
          <p class="text-text-muted">Loading...</p>
        </div>
      </div>
    </div>
  </AppLayout>
</template>
