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
import { 
  Plus, 
  Trash2, 
  Save, 
  X, 
  FolderPlus,
  FolderOpen
} from 'lucide-vue-next';

const router = useRouter();

const collectionFormData = ref({
  name: '',
  type: 'base',
  fields: [{ name: 'name', type: 'text', required: true }]
});

const addField = () => {
  collectionFormData.value.fields.push({ name: '', type: 'text', required: false });
};

const removeField = (index: number) => {
  collectionFormData.value.fields.splice(index, 1);
};

const saveCollection = async () => {
  try {
    await axios.post('/api/admin/collections', collectionFormData.value);
    router.push('/');
  } catch (error) {
    console.error('Collection creation failed', error);
  }
};
</script>

<template>
  <AppLayout>
    
      <!-- Header -->
      <AppHeader>
        <template #breadcrumb>
          <div class="flex items-center text-sm text-text-muted">
            <span class="hover:text-text cursor-pointer" @click="router.push('/')">Vault</span>
            <span class="mx-2">/</span>
            <span class="hover:text-text cursor-pointer" @click="router.push('/collections')">Collections</span>
            <span class="mx-2">/</span>
            <span class="font-medium text-text">New</span>
          </div>
        </template>
      </AppHeader>

      <!-- Form Content -->
      <div class="flex-1 overflow-auto min-h-0 p-4 sm:p-8 pb-24 sm:pb-8">
        <div class="max-w-4xl mx-auto space-y-6">
          <div>
            <h1 class="text-2xl font-bold text-text">Create New Collection</h1>
            <p class="text-sm text-text-muted mt-1">Define your database schema and fields</p>
          </div>

          <form @submit.prevent="saveCollection" class="space-y-6">
            <!-- Basic Info Card -->
            <div class="bg-surface-dark border border-border rounded-lg p-4 sm:p-6">
              <h2 class="text-lg font-semibold text-text mb-4 flex items-center gap-2">
                <FolderPlus class="w-5 h-5 text-primary" />
                Basic Information
              </h2>
              
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-6">
                <div>
                  <label class="block text-sm font-medium text-text mb-2">Collection Name</label>
                  <Input 
                    v-model="collectionFormData.name" 
                    type="text" 
                    required 
                    placeholder="e.g. products, users, posts"
                  />
                  <p class="text-xs text-text-dim mt-1">Lowercase, no spaces (use underscores)</p>
                </div>
                
                <div>
                  <label class="block text-sm font-medium text-text mb-2">Collection Type</label>
                  <Dropdown v-model="collectionFormData.type" align="left">
                    <template #trigger>
                      {{ collectionFormData.type === 'base' ? 'Base (Generic Data)' : 'Auth (User Records)' }}
                    </template>
                    <DropdownItem value="base" @select="collectionFormData.type = 'base'">Base (Generic Data)</DropdownItem>
                    <DropdownItem value="auth" @select="collectionFormData.type = 'auth'">Auth (User Records)</DropdownItem>
                  </Dropdown>
                  <p class="text-xs text-text-dim mt-1">Choose the collection purpose</p>
                </div>
              </div>
            </div>

            <!-- Fields Card -->
            <div class="bg-surface-dark border border-border rounded-lg p-4 sm:p-6">
              <div class="flex items-center justify-between mb-4">
                <h2 class="text-lg font-semibold text-text flex items-center gap-2">
                  <Plus class="w-5 h-5 text-primary" />
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
                  v-for="(field, index) in collectionFormData.fields" 
                  :key="index" 
                  class="flex flex-col sm:flex-row items-start sm:items-center gap-3 bg-surface p-4 rounded-lg border border-border"
                >
                  <div class="w-full sm:flex-1">
                    <Input 
                      v-model="field.name" 
                      placeholder="field_name"
                      type="text"
                    />
                  </div>
                  
                  <div class="w-full sm:w-40">
                    <Dropdown v-model="field.type" align="left">
                      <template #trigger>
                        {{ field.type.charAt(0).toUpperCase() + field.type.slice(1) }}
                      </template>
                      <DropdownItem value="text" @select="field.type = 'text'">Text</DropdownItem>
                      <DropdownItem value="number" @select="field.type = 'number'">Number</DropdownItem>
                      <DropdownItem value="bool" @select="field.type = 'bool'">Boolean</DropdownItem>
                      <DropdownItem value="json" @select="field.type = 'json'">JSON</DropdownItem>
                      <DropdownItem value="file" @select="field.type = 'file'">File</DropdownItem>
                    </Dropdown>
                  </div>

                  <div class="flex items-center justify-between w-full sm:w-auto gap-4">
                    <Checkbox v-model="field.required" label="Required" />
                    
                    <Button 
                      @click="removeField(index)" 
                      variant="ghost"
                      size="xs"
                      class="!text-error hover:!bg-error/10"
                      :disabled="collectionFormData.fields.length === 1"
                    >
                      <Trash2 class="w-4 h-4" />
                    </Button>
                  </div>
                </div>
              </div>
            </div>

            <!-- Actions -->
            <div class="flex items-center justify-end gap-3">
              <Button 
                @click="router.push('/')" 
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
                Create Collection
              </Button>
            </div>
          </form>
        </div>
      </div>
      </AppLayout>
</template>


