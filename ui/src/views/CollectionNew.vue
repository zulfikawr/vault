<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import Input from '../components/Input.vue';
import Checkbox from '../components/Checkbox.vue';
import Dropdown from '../components/Dropdown.vue';
import DropdownItem from '../components/DropdownItem.vue';
import { Plus, Trash2, Save, FolderPlus } from 'lucide-vue-next';

const router = useRouter();

const collectionFormData = ref({
  name: '',
  type: 'base',
  fields: [{ name: 'name', type: 'text', required: true }],
});

const addField = () => {
  collectionFormData.value.fields.push({ name: '', type: 'text', required: false });
};

const removeField = (index: number) => {
  collectionFormData.value.fields.splice(index, 1);
};

const validate = () => {
  if (!collectionFormData.value.name) {
    alert('Collection name is required');
    return false;
  }
  if (collectionFormData.value.fields.some(f => !f.name)) {
    alert('All fields must have a name');
    return false;
  }
  return true;
};

const saveCollection = async () => {
  if (!validate()) return;

  try {
    await axios.post('/api/admin/collections', collectionFormData.value);
    router.push('/collections');
  } catch (error: unknown) {
    console.error('Collection creation failed', error);
    let message = 'Failed to create collection';
    if (axios.isAxiosError(error) && error.response?.data?.message) {
      message = error.response.data.message;
    }
    alert(message);
  }
};
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
          <span class="font-medium text-primary flex-shrink-0">New</span>
        </div>
      </template>
    </AppHeader>

    <!-- Form Content -->
    <div class="flex-1 overflow-auto min-h-0 p-4 sm:p-8 pb-24 sm:pb-8">
      <div class="max-w-7xl mx-auto space-y-6 sm:space-y-8">
        <!-- Page Title -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h1 class="text-xl font-semibold text-text">Create New Collection</h1>
            <p class="mt-1 text-sm text-text-muted">Define your database schema and fields</p>
          </div>
          <div class="flex items-center gap-2">
            <Button
              variant="secondary"
              size="sm"
              class="px-3 py-1.5 text-sm"
              @click="router.push('/collections')"
            >
              Cancel
            </Button>
            <Button size="sm" class="px-3 py-1.5 text-sm" @click="saveCollection">
              <Save class="w-4 h-4" />
              Create
            </Button>
          </div>
        </div>

        <form id="collection-form" class="space-y-4" @submit.prevent="saveCollection">
          <!-- Basic Info Card -->
          <div class="bg-surface-dark border border-border rounded-lg p-4">
            <h2 class="text-base font-medium text-text mb-3 flex items-center gap-2">
              <FolderPlus class="w-4 h-4 text-primary" />
              Basic Information
            </h2>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
              <div>
                <label class="block text-sm font-medium text-text mb-2">Collection Name</label>
                <Input
                  v-model="collectionFormData.name"
                  type="text"
                  size="sm"
                  required
                  placeholder="e.g. products, users, posts"
                />
                <p class="text-xs text-text-dim mt-1">Lowercase, no spaces (use underscores)</p>
              </div>

              <div>
                <label class="block text-sm font-medium text-text mb-2">Collection Type</label>
                <Dropdown v-model="collectionFormData.type" align="left" size="sm">
                  <template #trigger>
                    {{
                      collectionFormData.type === 'base'
                        ? 'Base (Generic Data)'
                        : 'Auth (User Records)'
                    }}
                  </template>
                  <DropdownItem value="base" @select="collectionFormData.type = 'base'"
                    >Base (Generic Data)</DropdownItem
                  >
                  <DropdownItem value="auth" @select="collectionFormData.type = 'auth'"
                    >Auth (User Records)</DropdownItem
                  >
                </Dropdown>
                <p class="text-xs text-text-dim mt-1">Choose the collection purpose</p>
              </div>
            </div>
          </div>

          <!-- Fields Card -->
          <div class="bg-surface-dark border border-border rounded-lg p-4">
            <div class="flex items-center justify-between mb-3">
              <h2 class="text-base font-medium text-text flex items-center gap-2">
                <Plus class="w-4 h-4 text-primary" />
                Fields Definition
              </h2>
              <Button variant="secondary" size="sm" class="text-xs px-2.5 py-1" @click="addField">
                <Plus class="w-4 h-4" />
                <span class="hidden sm:inline">Add Field</span>
              </Button>
            </div>

            <div class="space-y-2">
              <div
                v-for="(field, index) in collectionFormData.fields"
                :key="index"
                class="flex flex-col sm:flex-row items-start sm:items-center gap-2 bg-surface p-3 rounded border border-border"
              >
                <div class="w-full sm:flex-1">
                  <Input v-model="field.name" placeholder="field_name" type="text" size="sm" />
                </div>

                <div class="w-full sm:w-32">
                  <Dropdown v-model="field.type" align="left" size="sm">
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

                <div class="flex items-center justify-between w-full sm:w-auto gap-3">
                  <Checkbox v-model="field.required" label="Required" size="sm" />

                  <Button
                    variant="ghost"
                    size="xs"
                    class="!text-error hover:!bg-error/10"
                    :disabled="collectionFormData.fields.length === 1"
                    @click="removeField(index)"
                  >
                    <Trash2 class="w-4 h-4" />
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>
  </AppLayout>
</template>
