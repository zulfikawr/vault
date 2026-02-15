<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import Input from '../components/Input.vue';
import Checkbox from '../components/Checkbox.vue';
import { X, Save } from 'lucide-vue-next';

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
const collection = ref<Collection | null>(null);
const formData = ref<Record>({});

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

const handleSubmit = async () => {
  try {
    await axios.patch(
      `/api/collections/${collectionName.value}/records/${recordId.value}`,
      formData.value
    );
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
        <div class="flex items-center text-sm text-text-muted truncate gap-2">
          <span
            class="hover:text-text cursor-pointer font-medium text-text"
            @click="router.push('/collections')"
            >Collections</span
          >
          <span class="text-text-muted flex-shrink-0">/</span>
          <span
            class="hover:text-text cursor-pointer text-text truncate"
            @click="router.push(`/collections/${collectionName}`)"
            >{{ collectionName }}</span
          >
          <span class="text-text-muted flex-shrink-0">/</span>
          <span class="font-medium text-text flex-shrink-0">Edit</span>
        </div>
      </template>
    </AppHeader>

    <div class="flex-1 overflow-auto min-h-0 p-4 sm:p-8 pb-24 sm:pb-8">
      <div class="max-w-7xl mx-auto space-y-6 sm:space-y-8">
        <!-- Page Title -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h1 class="text-xl font-semibold text-text">Edit Record</h1>
            <p class="mt-1 text-sm text-text-muted">
              Update record in {{ collectionName }} collection
            </p>
          </div>
          <div class="flex items-center gap-2">
            <Button
              variant="secondary"
              size="sm"
              class="px-3 py-1.5 text-sm"
              @click="router.push(`/collections/${collectionName}`)"
            >
              <X class="w-4 h-4" />
              Cancel
            </Button>
            <Button type="submit" size="sm" class="px-3 py-1.5 text-sm">
              <Save class="w-4 h-4" />
              Update
            </Button>
          </div>
        </div>

        <form
          v-if="collection"
          class="bg-surface-dark rounded-lg border border-border p-4 space-y-4"
          @submit.prevent="handleSubmit"
        >
          <div class="space-y-4">
            <div v-for="field in collection.fields" :key="field.name">
              <label :for="field.name" class="block text-sm font-medium text-text mb-1.5">
                {{ field.name }}
                <span v-if="field.required" class="text-error">*</span>
              </label>

              <Input
                v-if="field.type === 'text'"
                :model-value="String(formData[field.name] || '')"
                type="text"
                size="sm"
                :required="field.required"
                @update:model-value="formData[field.name] = $event"
              />

              <Input
                v-else-if="field.type === 'number'"
                :model-value="String(formData[field.name] || '')"
                type="number"
                size="sm"
                :required="field.required"
                @update:model-value="formData[field.name] = $event"
              />

              <Input
                v-else-if="field.type === 'email'"
                :model-value="String(formData[field.name] || '')"
                type="email"
                size="sm"
                :required="field.required"
                @update:model-value="formData[field.name] = $event"
              />

              <Input
                v-else-if="field.type === 'date'"
                :model-value="String(formData[field.name] || '')"
                type="date"
                size="sm"
                :required="field.required"
                @update:model-value="formData[field.name] = $event"
              />

              <div v-else-if="field.type === 'bool'" class="flex items-center">
                <Checkbox
                  :model-value="Boolean(formData[field.name])"
                  label="Enable"
                  @update:model-value="formData[field.name] = $event"
                />
              </div>

              <Input
                v-else-if="field.type === 'json'"
                :model-value="String(formData[field.name] || '')"
                type="textarea"
                size="sm"
                :required="field.required"
                :rows="4"
                @update:model-value="formData[field.name] = $event"
              />

              <Input
                v-else
                :model-value="String(formData[field.name] || '')"
                type="text"
                size="sm"
                :required="field.required"
                @update:model-value="formData[field.name] = $event"
              />
            </div>
          </div>
        </form>
      </div>
    </div>
  </AppLayout>
</template>
