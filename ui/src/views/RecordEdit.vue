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
  X,
  Save,
  Type,
  Hash,
  ToggleLeft,
  Calendar,
  Code,
  FileText,
  AlertCircle,
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
const collection = ref<Collection | null>(null);
const formData = ref<RecordData>({});
const errors = ref<Record<string, string>>({});

const collectionName = computed(() => route.params.name as string);
const recordId = computed(() => route.params.id as string);

const fetchCollection = async () => {
  try {
    const response = await axios.get(`/api/admin/collections`);
    const col = response.data.data.find((c: Collection) => c.name === collectionName.value);
    collection.value = col;
  } catch (error) {
    console.error('Failed to fetch collection', error);
  }
};

const fetchRecord = async () => {
  try {
    const response = await axios.get(
      `/api/collections/${collectionName.value}/records/${recordId.value}`
    );
    // The actual record fields are in response.data.data.data
    formData.value = response.data.data?.data || {};
  } catch (error) {
    console.error('Failed to fetch record', error);
  }
};

const validate = () => {
  errors.value = {};
  let isValid = true;

  collection.value?.fields.forEach((field) => {
    const value = formData.value[field.name];
    
    if (field.required && (value === undefined || value === null || value === '')) {
      if (field.type !== 'bool') {
        errors.value[field.name] = `${field.name} is required`;
        isValid = false;
      }
    }

    if (field.type === 'json' && value) {
      try {
        if (typeof value === 'string') {
          JSON.parse(value);
        }
      } catch (e) {
        errors.value[field.name] = 'Invalid JSON format';
        isValid = false;
      }
    }
  });

  return isValid;
};

const handleSubmit = async () => {
  if (!validate()) {
    alert('Please fix the validation errors before saving.');
    return;
  }

  try {
    await axios.patch(
      `/api/collections/${collectionName.value}/records/${recordId.value}`,
      formData.value
    );
    router.push(`/collections/${collectionName.value}`);
  } catch (error: unknown) {
    console.error('Update failed', error);
    let message = 'Failed to update record';
    if (axios.isAxiosError(error) && error.response?.data?.message) {
      message = error.response.data.message;
    }
    alert(message);
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
  fetchCollection();
  fetchRecord();
});
</script>

<template>
  <AppLayout>
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
          <span class="font-medium text-primary flex-shrink-0">Edit</span>
        </div>
      </template>
    </AppHeader>

    <div class="flex-1 overflow-auto min-h-0 p-4 sm:p-8 pb-24 sm:pb-8">
      <div class="max-w-7xl mx-auto space-y-6">
        <!-- Page Title -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-border pb-6">
          <div>
            <h1 class="text-2xl font-bold text-text tracking-tight">Edit Record</h1>
            <p class="mt-1 text-sm text-text-muted">
              Updating record <span class="font-mono text-primary">{{ recordId }}</span> in {{ collectionName }}
            </p>
          </div>
          <div class="flex items-center gap-2">
            <Button
              variant="secondary"
              size="sm"
              @click="router.push(`/collections/${collectionName}`)"
            >
              <X class="w-4 h-4" />
              Cancel
            </Button>
            <Button size="sm" @click="handleSubmit">
              <Save class="w-4 h-4" />
              Update
            </Button>
          </div>
        </div>

        <form
          v-if="collection"
          id="record-edit-form"
          class="space-y-6"
          @submit.prevent="handleSubmit"
        >
          <div class="grid grid-cols-1 gap-6">
            <div v-for="field in collection.fields" :key="field.name" class="space-y-2">
              <div class="flex items-center justify-between">
                <label :for="field.name" class="flex items-center gap-2 text-sm font-semibold text-text">
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
                  :model-value="String(formData[field.name] || '')"
                  :type="field.type === 'date' ? 'date' : field.type"
                  size="md"
                  @update:model-value="formData[field.name] = $event"
                  :class="{'!border-error !ring-error/20': errors[field.name]}"
                />

                <Input
                  v-else-if="field.type === 'number'"
                  :model-value="String(formData[field.name] || 0)"
                  type="number"
                  size="md"
                  @update:model-value="formData[field.name] = Number($event)"
                  :class="{'!border-error !ring-error/20': errors[field.name]}"
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
                  :model-value="typeof formData[field.name] === 'object' ? JSON.stringify(formData[field.name], null, 2) : String(formData[field.name] || '')"
                  type="textarea"
                  size="md"
                  :rows="6"
                  class="font-mono text-sm"
                  @update:model-value="formData[field.name] = $event"
                  :class="{'!border-error !ring-error/20': errors[field.name]}"
                />

                <div v-else-if="field.type === 'file'" class="p-8 border-2 border-dashed border-border rounded-lg bg-surface-dark/30 flex flex-col items-center justify-center text-text-muted">
                   <FileText class="w-10 h-10 mb-2 opacity-20" />
                   <p class="text-sm font-medium">File management via record forms coming soon</p>
                   <p class="text-[11px] mt-1">Use Storage tab to manage media</p>
                </div>

                <Input
                  v-else
                  :model-value="String(formData[field.name] || '')"
                  type="text"
                  size="md"
                  @update:model-value="formData[field.name] = $event"
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
