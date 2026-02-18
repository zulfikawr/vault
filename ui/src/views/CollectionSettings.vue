<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import Input from '../components/Input.vue';
import Checkbox from '../components/Checkbox.vue';
import ConfirmModal from '../components/ConfirmModal.vue';
import Dropdown from '../components/Dropdown.vue';
import DropdownItem from '../components/DropdownItem.vue';
import { Plus, Trash2, Settings, Save } from 'lucide-vue-next';

interface Field {
  name: string;
  type: string;
  required: boolean;
}

interface Collection {
  id: string;
  name: string;
  fields: Field[];
}

const router = useRouter();
const route = useRoute();
const collections = ref<Collection[]>([]);
const collection = ref<Collection | null>(null);
const fields = ref<Field[]>([]);
const rules = ref({
  list_rule: '',
  view_rule: '',
  create_rule: '',
  update_rule: '',
  delete_rule: '',
});
const showDeleteModal = ref(false);

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
    collection.value = col || null;
    fields.value = JSON.parse(JSON.stringify(col?.fields || []));
    rules.value = {
      list_rule: col?.list_rule || '',
      view_rule: col?.view_rule || '',
      create_rule: col?.create_rule || '',
      update_rule: col?.update_rule || '',
      delete_rule: col?.delete_rule || '',
    };
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
  if (!collection.value) return;

  try {
    await axios.patch(`/api/admin/collections/${collection.value.id}`, {
      ...collection.value,
      fields: fields.value,
      ...rules.value,
    });
    router.push(`/collections/${collectionName.value}`);
  } catch (error) {
    console.error('Save failed', error);
  }
};

const deleteCollection = async () => {
  if (!collection.value) return;

  try {
    await axios.delete(`/api/admin/collections/${collection.value.id}`);
    router.push('/collections'); // Redirect to collections list after deletion
  } catch (error) {
    console.error('Failed to delete collection', error);
    // TODO: Show error notification to user
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
          <span class="font-medium text-primary flex-shrink-0">Settings</span>
        </div>
      </template>
    </AppHeader>

    <!-- Form Content -->
    <div class="flex-1 overflow-auto min-h-0 p-4 sm:p-8 pb-24 sm:pb-8">
      <div class="max-w-7xl mx-auto space-y-6 sm:space-y-8">
        <!-- Page Title and Actions -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h1 class="text-xl font-semibold text-text">Collection Settings</h1>
            <p class="mt-1 text-sm text-text-muted">Manage fields for {{ collectionName }}</p>
          </div>
          <div class="flex items-center gap-2">
            <Button
              variant="secondary"
              size="sm"
              class="px-3 py-1.5 text-sm"
              @click="router.push(`/collections/${collectionName}`)"
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              size="sm"
              class="px-3 py-1.5 text-sm"
              @click="showDeleteModal = true"
            >
              <Trash2 class="w-4 h-4" />
              Delete
            </Button>
            <Button size="sm" class="px-3 py-1.5 text-sm" @click="saveSettings">
              <Save class="w-4 h-4" />
              Save
            </Button>
          </div>
        </div>

        <form id="collection-settings-form" class="space-y-4" @submit.prevent="saveSettings">
          <div class="flex items-center justify-between mb-3">
            <h2 class="text-base font-medium text-text flex items-center gap-2">
              <Settings class="w-4 h-4 text-primary" />
              Fields Definition
            </h2>
            <Button variant="secondary" size="sm" class="text-xs px-2.5 py-1" @click="addField">
              <Plus class="w-4 h-4" />
              <span class="hidden sm:inline">Add Field</span>
            </Button>
          </div>

          <div class="space-y-2">
            <div
              v-for="(field, index) in fields"
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
                  <DropdownItem value="number" @select="field.type = 'number'">Number</DropdownItem>
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
                  :disabled="fields.length === 1"
                  @click="removeField(index)"
                >
                  <Trash2 class="w-4 h-4" />
                </Button>
              </div>
            </div>
          </div>

          <div class="pt-4 border-t border-border">
            <h2 class="text-base font-medium text-text flex items-center gap-2 mb-3">
              <Settings class="w-4 h-4 text-primary" />
              API Rules
            </h2>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="space-y-2">
                <label class="text-xs font-semibold text-text-muted uppercase tracking-wider"
                  >List Rule</label
                >
                <Input
                  v-model="rules.list_rule"
                  placeholder="e.g. id = @request.auth.id"
                  size="sm"
                />
                <p class="text-[10px] text-text-dim">Leave empty for public access</p>
              </div>
              <div class="space-y-2">
                <label class="text-xs font-semibold text-text-muted uppercase tracking-wider"
                  >View Rule</label
                >
                <Input
                  v-model="rules.view_rule"
                  placeholder="e.g. id = @request.auth.id"
                  size="sm"
                />
              </div>
              <div class="space-y-2">
                <label class="text-xs font-semibold text-text-muted uppercase tracking-wider"
                  >Create Rule</label
                >
                <Input
                  v-model="rules.create_rule"
                  placeholder="e.g. @request.auth.id != ''"
                  size="sm"
                />
              </div>
              <div class="space-y-2">
                <label class="text-xs font-semibold text-text-muted uppercase tracking-wider"
                  >Update Rule</label
                >
                <Input
                  v-model="rules.update_rule"
                  placeholder="e.g. id = @request.auth.id"
                  size="sm"
                />
              </div>
              <div class="space-y-2">
                <label class="text-xs font-semibold text-text-muted uppercase tracking-wider"
                  >Delete Rule</label
                >
                <Input
                  v-model="rules.delete_rule"
                  placeholder="e.g. @request.auth.id = 'admin-id'"
                  size="sm"
                />
              </div>
            </div>
          </div>
        </form>

        <!-- Delete Confirmation Modal -->
        <ConfirmModal
          :show="showDeleteModal"
          title="Confirm Delete Collection"
          :message="`Are you sure you want to delete the collection '${collectionName}'? This action cannot be undone and will permanently remove all data in this collection.`"
          confirm-text="Delete"
          cancel-text="Cancel"
          variant="danger"
          @confirm="deleteCollection"
          @cancel="showDeleteModal = false"
        />
      </div>
    </div>
  </AppLayout>
</template>
