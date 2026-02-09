<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';
import axios from 'axios';

const auth = useAuthStore();
const router = useRouter();
const collections = ref([]);
const selectedCollection = ref(null);
const records = ref([]);
const showEditor = ref(false);
const editingRecord = ref(null);
const formData = ref({});

const fetchCollections = async () => {
  try {
    const response = await axios.get('/api/admin/collections');
    collections.value = response.data.data;
  } catch (error) {
    console.error('Failed to fetch collections', error);
  }
};

const selectCollection = async (col) => {
  selectedCollection.value = col;
  fetchRecords();
};

const fetchRecords = async () => {
  if (!selectedCollection.value) return;
  try {
    const response = await axios.get(`/api/collections/${selectedCollection.value.name}/records`);
    records.value = response.data.items;
  } catch (error) {
    console.error('Failed to fetch records', error);
  }
};

const handleLogout = () => {
  auth.logout();
  router.push({ name: 'Login' });
};

const openEditor = (record = null) => {
  editingRecord.value = record;
  formData.value = record ? { ...record } : {};
  showEditor.value = true;
};

const saveRecord = async () => {
  try {
    const url = editingRecord.value 
      ? `/api/collections/${selectedCollection.value.name}/records/${editingRecord.value.id}`
      : `/api/collections/${selectedCollection.value.name}/records`;
    
    const method = editingRecord.value ? 'patch' : 'post';
    
    await axios[method](url, formData.value);
    showEditor.value = false;
    fetchRecords();
  } catch (error) {
    console.error('Failed to save record', error);
  }
};

const deleteRecord = async (id) => {
  if (!confirm('Are you sure?')) return;
  try {
    await axios.delete(`/api/collections/${selectedCollection.value.name}/records/${id}`);
    fetchRecords();
  } catch (error) {
    console.error('Failed to delete record', error);
  }
};

onMounted(fetchCollections);
</script>

<template>
  <div class="flex h-screen bg-gray-100">
    <!-- Sidebar -->
    <div class="w-64 bg-slate-800 text-white flex flex-col">
      <div class="p-4 text-xl font-bold border-b border-slate-700 flex justify-between items-center">
        <span>Vault</span>
        <button @click="handleLogout" class="text-sm text-slate-400 hover:text-white">Logout</button>
      </div>
      <div class="flex-1 overflow-y-auto">
        <div class="p-4 text-xs font-semibold text-slate-400 uppercase tracking-wider">Collections</div>
        <nav class="space-y-1 px-2">
          <button
            v-for="col in collections"
            :key="col.name"
            @click="selectCollection(col)"
            class="w-full text-left flex items-center px-2 py-2 text-sm font-medium rounded-md hover:bg-slate-700 transition-colors"
            :class="selectedCollection?.name === col.name ? 'bg-slate-700 text-white' : 'text-slate-300'"
          >
            {{ col.name }}
          </button>
        </nav>
      </div>
    </div>

    <!-- Main Content -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <header class="bg-white shadow-sm z-10 flex justify-between items-center px-6 py-4">
        <h1 class="text-lg font-semibold leading-6 text-gray-900">
          {{ selectedCollection ? selectedCollection.name : 'Select a collection' }}
        </h1>
        <button v-if="selectedCollection" @click="openEditor()" class="bg-blue-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-blue-700">
          New Record
        </button>
      </header>

      <main class="flex-1 overflow-y-auto p-6">
        <div v-if="selectedCollection">
          <div class="bg-white shadow rounded-lg overflow-hidden border border-gray-200">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th v-for="field in selectedCollection.fields" :key="field.name" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    {{ field.name }}
                  </th>
                  <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr v-for="record in records" :key="record.id">
                  <td v-for="field in selectedCollection.fields" :key="field.name" class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 overflow-hidden max-w-xs truncate">
                    {{ record[field.name] }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium space-x-2">
                    <button @click="openEditor(record)" class="text-blue-600 hover:text-blue-900">Edit</button>
                    <button @click="deleteRecord(record.id)" class="text-red-600 hover:text-red-900">Delete</button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        <div v-else class="flex flex-col items-center justify-center h-full text-gray-500">
          <p>Please select a collection from the sidebar to view data.</p>
        </div>
      </main>
    </div>

    <!-- Record Editor Modal -->
    <div v-if="showEditor" class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" @click="showEditor = false"></div>
        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <form @submit.prevent="saveRecord">
            <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
              <h3 class="text-lg leading-6 font-medium text-gray-900 mb-4" id="modal-title">
                {{ editingRecord ? 'Edit Record' : 'New Record' }}
              </h3>
              <div class="space-y-4">
                <div v-for="field in selectedCollection.fields" :key="field.name">
                  <label class="block text-sm font-medium text-gray-700">{{ field.name }}</label>
                  <input
                    v-if="field.type === 'text' || field.type === 'number'"
                    v-model="formData[field.name]"
                    :type="field.type === 'number' ? 'number' : 'text'"
                    class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 focus:ring-blue-500 focus:border-blue-500"
                  />
                  <select
                    v-else-if="field.type === 'bool'"
                    v-model="formData[field.name]"
                    class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm p-2 focus:ring-blue-500 focus:border-blue-500"
                  >
                    <option :value="true">True</option>
                    <option :value="false">False</option>
                  </select>
                </div>
              </div>
            </div>
            <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
              <button type="submit" class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:ml-3 sm:w-auto sm:text-sm">
                Save
              </button>
              <button type="button" @click="showEditor = false" class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm">
                Cancel
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>