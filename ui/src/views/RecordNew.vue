<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import Input from '../components/Input.vue';
import Checkbox from '../components/Checkbox.vue';
import {
  Save,
  Type,
  Hash,
  ToggleLeft,
  Calendar,
  Code,
  FileText,
  AlertCircle,
  X,
} from 'lucide-vue-next';

interface Field {
  name: string;
  type: string;
  required: boolean;
}

interface Collection {
  name: string;
  fields: Field[];
}

interface RecordData {
  [key: string]: string | number | boolean;
}

const router = useRouter();
const route = useRoute();
const collections = ref<Collection[]>([]);
const collection = ref<Collection | null>(null);
const formData = ref<RecordData>({});
const errors = ref<Record<string, string>>({});

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

const validate = () => {
  errors.value = {};
  let isValid = true;

  collection.value?.fields.forEach((field) => {
    if (field.required && !formData.value[field.name]) {
      if (field.type !== 'bool') {
        // bool is never 'empty' as it's false
        errors.value[field.name] = `${field.name} is required`;
        isValid = false;
      }
    }

    if (field.type === 'json' && formData.value[field.name]) {
      try {
        JSON.parse(String(formData.value[field.name]));
      } catch (e) {
        errors.value[field.name] = 'Invalid JSON format';
        isValid = false;
      }
    }
  });

  return isValid;
};

const saveRecord = async () => {
  if (!validate()) return;

  try {
    await axios.post(`/api/collections/${collectionName.value}/records`, formData.value);
    router.push(`/collections/${collectionName.value}`);
  } catch (error) {
    console.error('Save failed', error);
    alert('Failed to save record. Check console for details.');
  }
};

const getFieldIcon = (type: string) => {
  switch (type) {
    case 'number':
      return Hash;
    case 'bool':
      return ToggleLeft;
    case 'date':
      return Calendar;
    case 'json':
      return Code;
    case 'file':
      return FileText;
    default:
      return Type;
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
          <span
            class="hover:text-primary cursor-pointer transition-colors duration-200"
            @click="router.push('/collections')"
            >Collections</span
          >
          <span class="text-text-muted flex-shrink-0">/</span>
          <span
            class="hover:text-primary cursor-pointer truncate transition-colors duration-200"
            @click="router.push(`/collections/${collectionName}`)"
            >{{ collectionName }}</span
          >
          <span class="text-text-muted flex-shrink-0">/</span>
          <span class="font-medium text-primary flex-shrink-0">New</span>
        </div>
      </template>
    </AppHeader>

    <!-- Form Content -->
    <div class="flex-1 overflow-auto min-h-0 p-4 sm:p-8 pb-24 sm:pb-8">
      <div class="max-w-4xl mx-auto space-y-6">
        <!-- Page Title -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-border pb-6">
          <div>
            <h1 class="text-2xl font-bold text-text tracking-tight">Create New Record</h1>
            <p class="mt-1 text-sm text-text-muted">Populate fields for your {{ collectionName }} collection</p>
          </div>
          <div class="flex items-center gap-3">
            <Button
              variant="secondary"
              size="sm"
              @click="router.push(`/collections/${collectionName}`)"
            >
              <X class="w-4 h-4" />
              Cancel
            </Button>
            <Button size="sm" @click="saveRecord">
              <Save class="w-4 h-4" />
              Create
            </Button>
          </div>
        </div>

        <form class="space-y-6" @submit.prevent="saveRecord">
          <div class="grid grid-cols-1 gap-6">
            <div v-for="field in collection?.fields" :key="field.name" class="space-y-2">
              <div class="flex items-center justify-between">
                <label class="flex items-center gap-2 text-sm font-semibold text-text">
                  <span>{{ field.name }}</span>
                  <span v-if="field.required" class="text-error" title="Required">*</span>
                </label>
                <div 
                  class="flex items-center gap-1.5 px-2 py-0.5 rounded-full bg-surface-dark border border-border text-[10px] font-bold uppercase tracking-widest text-text-muted"
                >
                  <component :is="getFieldIcon(field.type)" class="w-3 h-3" />
                  {{ field.type }}
                </div>
              </div>

              <div class="relative">
                <Input
                  v-if="field.type === 'text' || field.type === 'email' || field.type === 'date' || field.type === 'url'"
                  v-model="formData[field.name] as string"
                  :type="field.type === 'date' ? 'date' : field.type"
                  size="md"
                  :placeholder="`Enter ${field.name}...`"
                  :class="{'!border-error !ring-error/20': errors[field.name]}"
                />

                <Input
                  v-else-if="field.type === 'number'"
                  :model-value="String(formData[field.name])"
                  type="number"
                  size="md"
                  :placeholder="`Enter ${field.name}...`"
                  :class="{'!border-error !ring-error/20': errors[field.name]}"
                  @update:model-value="formData[field.name] = Number($event)"
                />

                <div v-else-if="field.type === 'bool'" class="p-3 bg-surface border border-border rounded-lg">
                  <Checkbox
                    :model-value="Boolean(formData[field.name])"
                    :label="Boolean(formData[field.name]) ? 'Enabled' : 'Disabled'"
                    @update:model-value="formData[field.name] = $event"
                  />
                </div>

                <Input
                  v-else-if="field.type === 'json'"
                  v-model="formData[field.name] as string"
                  type="textarea"
                  size="md"
                  :placeholder='`{"key": "value"}`'
                  :rows="6"
                  class="font-mono text-sm"
                  :class="{'!border-error !ring-error/20': errors[field.name]}"
                />

                <div v-else-if="field.type === 'file'" class="p-8 border-2 border-dashed border-border rounded-lg bg-surface-dark/30 flex flex-col items-center justify-center text-text-muted">
                   <FileText class="w-10 h-10 mb-2 opacity-20" />
                   <p class="text-sm font-medium">File uploading via record forms coming soon</p>
                   <p class="text-[11px] mt-1">Use Storage tab to manage media</p>
                </div>

                <Input
                  v-else
                  v-model="formData[field.name] as string"
                  type="text"
                  size="md"
                  :class="{'!border-error !ring-error/20': errors[field.name]}"
                />
              </div>

              <!-- Error Message -->
              <div v-if="errors[field.name]" class="flex items-center gap-1.5 text-error text-[11px] font-medium">
                <AlertCircle class="w-3.5 h-3.5" />
                {{ errors[field.name] }}
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>
  </AppLayout>
</template>
