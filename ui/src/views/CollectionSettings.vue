<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import Input from '../components/Input.vue';
import Checkbox from '../components/Checkbox.vue';
import Dropdown from '../components/Dropdown.vue';
import DropdownItem from '../components/DropdownItem.vue';
import { Plus, Trash2, Settings, Save } from 'lucide-vue-next';

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
    const col = response.data.data.find(
      (c: Record<string, unknown>) => c.name === collectionName.value
    );
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
      fields: fields.value,
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

        <form class="space-y-6" @submit.prevent="saveSettings">
          <div class="bg-surface-dark border border-border rounded-lg p-4 sm:p-6">
            <div class="flex items-center justify-between mb-4">
              <h2 class="text-lg font-semibold text-text flex items-center gap-2">
                <Settings class="w-5 h-5 text-primary" />
                Fields Definition
              </h2>
              <Button variant="secondary" size="sm" @click="addField">
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
                  <Input v-model="field.name" placeholder="field_name" type="text" />
                </div>

                <div class="w-full sm:w-40">
                  <Dropdown v-model="field.type" align="left">
                    <template #trigger>
                      {{ field.type.charAt(0).toUpperCase() + field.type.slice(1) }}
                    </template>
                    <DropdownItem value="text" @select="field.type = 'text'">Text</DropdownItem>
                    <DropdownItem value="number" @select="field.type = 'number'"
                      >Number</DropdownItem
                    >
                    <DropdownItem value="bool" @select="field.type = 'bool'">Boolean</DropdownItem>
                    <DropdownItem value="json" @select="field.type = 'json'">JSON</DropdownItem>
                    <DropdownItem value="file" @select="field.type = 'file'">File</DropdownItem>
                  </Dropdown>
                </div>

                <div class="flex items-center justify-between w-full sm:w-auto gap-4">
                  <Checkbox v-model="field.required" label="Required" />

                  <Button
                    variant="ghost"
                    size="xs"
                    class="!text-error hover:!bg-error/10"
                    :disabled="fields.length === 1"
                    @click="removeField(index)"
                  >
                    <Trash2 class="w-4 h-4" />
                  </Button>
                </div>
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
              Save Changes
            </Button>
          </div>
        </form>
      </div>
    </div>
  </AppLayout>
</template>
