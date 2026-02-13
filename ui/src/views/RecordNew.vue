<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import Input from '../components/Input.vue';
import Dropdown from '../components/Dropdown.vue';
import DropdownItem from '../components/DropdownItem.vue';
import { Save } from 'lucide-vue-next';

const router = useRouter();
const route = useRoute();
const collections = ref<Record<string, unknown>[]>([]);
const collection = ref<Record<string, unknown> | null>(null);
const formData = ref<Record<string, unknown>>({});

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
    const col = response.data.data.find(
      (c: Record<string, unknown>) => c.name === collectionName.value
    );
    collection.value = col;
    // Initialize form data with empty values
    col?.fields?.forEach((field: Record<string, unknown>) => {
      formData.value[field.name as string] = '';
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
        <div class="flex items-center text-sm text-text-muted overflow-hidden whitespace-nowrap">
          <span class="hover:text-text cursor-pointer shrink-0" @click="router.push('/')"
            >Vault</span
          >
          <span class="mx-2 shrink-0">/</span>
          <span
            class="hover:text-text cursor-pointer shrink-0 hidden sm:inline"
            @click="router.push('/collections')"
            >Collections</span
          >
          <span class="mx-2 shrink-0 hidden sm:inline">/</span>
          <span
            class="hover:text-text cursor-pointer truncate"
            @click="router.push(`/collections/${collectionName}`)"
            >{{ collectionName }}</span
          >
          <span class="mx-2 shrink-0">/</span>
          <span class="font-medium text-text shrink-0">New</span>
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
                  v-model="formData[field.name]"
                  type="text"
                  :required="field.required"
                />

                <Input
                  v-else-if="field.type === 'number'"
                  v-model="formData[field.name]"
                  type="number"
                  :required="field.required"
                />

                <Dropdown
                  v-else-if="field.type === 'bool'"
                  v-model="formData[field.name]"
                  align="left"
                >
                  <template #trigger>
                    {{
                      formData[field.name] === ''
                        ? 'Select...'
                        : formData[field.name]
                          ? 'True'
                          : 'False'
                    }}
                  </template>
                  <DropdownItem value="" @select="formData[field.name] = ''"
                    >Select...</DropdownItem
                  >
                  <DropdownItem :value="true" @select="formData[field.name] = true"
                    >True</DropdownItem
                  >
                  <DropdownItem :value="false" @select="formData[field.name] = false"
                    >False</DropdownItem
                  >
                </Dropdown>

                <Input
                  v-else-if="field.type === 'json'"
                  v-model="formData[field.name]"
                  type="textarea"
                  :required="field.required"
                  placeholder='{"key": "value"}'
                  :rows="4"
                />

                <Input
                  v-else
                  v-model="formData[field.name]"
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
