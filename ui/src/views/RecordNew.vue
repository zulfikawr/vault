<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import Input from '../components/Input.vue';
import { Save } from 'lucide-vue-next';

interface Field {
  name: string;
  type: string;
  required: boolean;
}

interface Collection {
  name: string;
  fields: Field[];
}

interface Record {
  [key: string]: string | number | boolean;
}

const router = useRouter();
const route = useRoute();
const collections = ref<Collection[]>([]);
const collection = ref<Collection | null>(null);
const formData = ref<Record>({});

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
    const col = response.data.data.find((c: Collection) => c.name === collectionName.value);
    collection.value = col;
    // Initialize form data with appropriate default values
    col?.fields?.forEach((field: Field) => {
      if (field.type === 'bool') {
        formData.value[field.name] = false;
      } else if (field.type === 'number') {
        formData.value[field.name] = 0;
      } else {
        formData.value[field.name] = '';
      }
    });
  } catch (error) {
    console.error('Failed to fetch collection', error);
  }
};

const saveRecord = async () => {
  try {
    await axios.post(`/api/collections/${collectionName.value}/records`, formData.value);
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
        <div class="flex items-center text-sm text-text-muted truncate gap-2">
          <span class="hover:text-text cursor-pointer font-medium text-text" @click="router.push('/collections')">Collections</span>
          <span class="text-text-muted flex-shrink-0">/</span>
          <span class="hover:text-text cursor-pointer text-text truncate" @click="router.push(`/collections/${collectionName}`)">{{ collectionName }}</span>
          <span class="text-text-muted flex-shrink-0">/</span>
          <span class="font-medium text-text flex-shrink-0">New</span>
        </div>
      </template>
    </AppHeader>

    <!-- Form Content -->
    <div class="flex-1 overflow-auto min-h-0 p-4 sm:p-8 pb-24 sm:pb-8">
      <div class="max-w-4xl mx-auto space-y-6">
        <div>
          <h1 class="text-2xl font-bold text-text">Create New Record</h1>
          <p class="text-sm text-text-muted mt-1">Add a new record to {{ collectionName }}</p>
        </div>

        <form class="space-y-6" @submit.prevent="saveRecord">
          <div class="bg-surface-dark border border-border rounded-lg p-4 sm:p-6">
            <div class="space-y-6">
              <div v-for="field in collection?.fields" :key="field.name">
                <label class="block text-sm font-medium text-text mb-2">
                  {{ field.name }}
                  <span v-if="field.required" class="text-error">*</span>
                </label>

                <Input
                  v-if="field.type === 'text'"
                  v-model="formData[field.name] as string"
                  type="text"
                  :required="field.required"
                />

                <Input
                  v-else-if="field.type === 'number'"
                  v-model="formData[field.name] as number"
                  type="number"
                  :required="field.required"
                />

                <input
                  v-else-if="field.type === 'bool'"
                  v-model="formData[field.name]"
                  type="checkbox"
                  class="w-4 h-4 bg-surface border border-border rounded focus:ring-2 focus:ring-primary/50 focus:border-primary"
                />

                <Input
                  v-else-if="field.type === 'json'"
                  v-model="formData[field.name] as string"
                  type="textarea"
                  :required="field.required"
                  placeholder='{"key": "value"}'
                  :rows="4"
                />

                <Input
                  v-else
                  v-model="formData[field.name] as string"
                  type="text"
                  :required="field.required"
                />

                <p class="text-xs text-text-dim mt-1">{{ field.type }}</p>
              </div>
            </div>
          </div>

          <div class="flex items-center justify-end gap-3">
            <Button
              variant="secondary"
              class="flex-1 sm:flex-none"
              @click="router.push(`/collections/${collectionName}`)"
            >
              Cancel
            </Button>
            <Button type="submit" class="flex-1 sm:flex-none">
              <Save class="w-4 h-4" />
              Create Record
            </Button>
          </div>
        </form>
      </div>
    </div>
  </AppLayout>
</template>
