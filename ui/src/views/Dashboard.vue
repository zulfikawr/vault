<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useAuthStore } from '../stores/auth';
import { useRouter } from 'vue-router';
import axios from 'axios';
import { 
  LayoutGrid, 
  Database, 
  Plus, 
  LogOut, 
  Search, 
  Trash2, 
  Edit3, 
  X,
  ChevronRight,
  ShieldCheck,
  Zap,
  Save,
  ChevronLeft
} from 'lucide-vue-next';

const auth = useAuthStore();
const router = useRouter();
const collections = ref([]);
const selectedCollection = ref(null);
const records = ref([]);

// Sidebar Drawer States
const showRecordDrawer = ref(false);
const showCollectionDrawer = ref(false);

// Forms
const editingRecord = ref(null);
const recordFormData = ref({});
const collectionFormData = ref({
  name: '',
  type: 'base',
  fields: [{ name: 'name', type: 'text', required: true }]
});

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
    records.value = response.data.items || [];
  } catch (error) {
    console.error('Failed to fetch records', error);
  }
};

const handleLogout = () => {
  auth.logout();
  router.push({ name: 'Login' });
};

// Record Actions
const openRecordDrawer = (record = null) => {
  editingRecord.value = record;
  recordFormData.value = record ? { ...record } : {};
  showRecordDrawer.value = true;
};

const saveRecord = async () => {
  try {
    const url = editingRecord.value 
      ? `/api/collections/${selectedCollection.value.name}/records/${editingRecord.value.id}`
      : `/api/collections/${selectedCollection.value.name}/records`;
    const method = editingRecord.value ? 'patch' : 'post';
    await axios[method](url, recordFormData.value);
    showRecordDrawer.value = false;
    fetchRecords();
  } catch (error) {
    console.error('Save failed', error);
  }
};

const deleteRecord = async (id) => {
  if (!confirm('Are you sure you want to delete this record?')) return;
  try {
    await axios.delete(`/api/collections/${selectedCollection.value.name}/records/${id}`);
    fetchRecords();
  } catch (error) {
    console.error('Delete failed', error);
  }
};

// Collection Actions
const openCollectionDrawer = () => {
  collectionFormData.value = {
    name: '',
    type: 'base',
    fields: [{ name: 'name', type: 'text', required: true }]
  };
  showCollectionDrawer.value = true;
};

const saveCollection = async () => {
  try {
    await axios.post('/api/admin/collections', collectionFormData.value);
    showCollectionDrawer.value = false;
    fetchCollections();
  } catch (error) {
    console.error('Collection creation failed', error);
  }
};

onMounted(fetchCollections);
</script>

<template>
  <div class="flex h-screen bg-monokai-bg font-sans text-monokai-fg overflow-hidden relative">
    <!-- Sidebar -->
    <div class="w-64 bg-monokai-panel flex flex-col shrink-0 border-r border-white/5">
      <div class="p-6 flex items-center">
        <span class="text-xl font-bold tracking-tight text-white">Vault</span>
      </div>
      
      <div class="flex-1 overflow-y-auto px-3">
        <div class="flex items-center justify-between px-3 py-4">
          <span class="text-[10px] font-bold text-monokai-gray uppercase tracking-widest">Collections</span>
          <button @click="openCollectionDrawer" class="p-1 hover:bg-monokai-header rounded-md text-monokai-green transition-colors cursor-pointer">
            <Plus class="w-4 h-4" />
          </button>
        </div>
        
        <nav class="space-y-1">
          <button
            v-for="col in collections"
            :key="col.name"
            @click="selectCollection(col)"
            class="w-full flex items-center justify-between px-3 py-2 rounded-xl text-sm font-medium transition-all group cursor-pointer"
            :class="selectedCollection?.name === col.name 
              ? 'bg-monokai-header text-white shadow-sm' 
              : 'text-monokai-gray hover:text-white hover:bg-white/5'"
          >
            <div class="flex items-center space-x-3 truncate">
              <Database class="w-4 h-4 shrink-0" :class="selectedCollection?.name === col.name ? 'text-monokai-blue' : 'text-monokai-gray/50'" />
              <span class="truncate">{{ col.name }}</span>
            </div>
            <ChevronRight class="w-3 h-3 opacity-0 group-hover:opacity-100 transition-opacity" />
          </button>
        </nav>
      </div>

      <div class="p-4 border-t border-white/5">
        <div class="flex items-center space-x-3 px-3 py-2 mb-2">
          <div class="w-8 h-8 bg-monokai-blue/20 rounded-full flex items-center justify-center text-monokai-blue">
            <ShieldCheck class="w-4 h-4" />
          </div>
          <div class="truncate">
            <p class="text-xs font-bold text-white truncate">{{ auth.user?.username || 'Admin' }}</p>
            <p class="text-[10px] text-monokai-gray truncate">Superuser</p>
          </div>
        </div>
        <button 
          @click="handleLogout" 
          class="w-full flex items-center space-x-3 px-3 py-2 rounded-xl text-xs font-bold text-monokai-pink hover:bg-monokai-pink/10 transition-colors cursor-pointer"
        >
          <LogOut class="w-4 h-4" />
          <span>Sign Out</span>
        </button>
      </div>
    </div>

    <!-- Main Content -->
    <div class="flex-1 flex flex-col min-w-0 bg-[#2d2e27] relative">
      <header v-if="selectedCollection" class="h-20 bg-monokai-bg border-b border-white/5 px-8 flex items-center justify-between shrink-0">
        <div>
          <h1 class="text-xl font-bold text-white tracking-tight flex items-center space-x-3">
            <span class="text-monokai-blue">/</span>
            <span>{{ selectedCollection.name }}</span>
          </h1>
        </div>
        
        <div class="flex items-center space-x-4">
          <div class="hidden lg:flex relative">
            <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-monokai-gray" />
            <input type="text" placeholder="Search records..." class="bg-monokai-panel border border-white/10 rounded-full py-2 pl-10 pr-4 text-xs text-white focus:border-monokai-blue outline-none transition-all w-64" />
          </div>
          <button 
            @click="openRecordDrawer()" 
            class="bg-monokai-green hover:bg-monokai-green/90 text-monokai-bg px-5 py-2 rounded-full text-xs font-bold transition-all shadow-lg shadow-monokai-green/10 flex items-center space-x-2 cursor-pointer"
          >
            <Plus class="w-4 h-4" />
            <span>New Record</span>
          </button>
        </div>
      </header>

      <main class="flex-1 overflow-auto p-8">
        <div v-if="selectedCollection">
          <div class="bg-monokai-panel rounded-2xl border border-white/5 shadow-xl overflow-hidden">
            <table class="w-full text-left border-collapse">
              <thead>
                <tr class="bg-white/5 border-b border-white/5">
                  <th v-for="field in selectedCollection.fields" :key="field.name" class="px-6 py-4 text-[10px] font-black text-monokai-gray uppercase tracking-widest">
                    {{ field.name }}
                  </th>
                  <th class="px-6 py-4 text-right"></th>
                </tr>
              </thead>
              <tbody class="divide-y divide-white/5">
                <tr v-for="record in records" :key="record.id" class="hover:bg-white/[0.02] transition-colors group">
                  <td v-for="field in selectedCollection.fields" :key="field.name" class="px-6 py-4 text-sm font-medium text-white/80 max-w-xs truncate font-mono">
                    {{ record[field.name] }}
                  </td>
                  <td class="px-6 py-4 text-right">
                    <div class="flex justify-end items-center space-x-2 opacity-0 group-hover:opacity-100 transition-opacity">
                      <button @click="openRecordDrawer(record)" class="p-2 hover:bg-monokai-blue/20 text-monokai-blue rounded-lg transition-colors cursor-pointer">
                        <Edit3 class="w-4 h-4" />
                      </button>
                      <button @click="deleteRecord(record.id)" class="p-2 hover:bg-monokai-pink/20 text-monokai-pink rounded-lg transition-colors cursor-pointer">
                        <Trash2 class="w-4 h-4" />
                      </button>
                    </div>
                  </td>
                </tr>
                <tr v-if="records.length === 0">
                  <td :colspan="selectedCollection.fields.length + 1" class="py-20 text-center">
                    <div class="flex flex-col items-center justify-center opacity-30">
                      <LayoutGrid class="w-12 h-12 mb-4" />
                      <p class="text-sm font-medium">No records found</p>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
        
        <div v-else class="h-full flex items-center justify-center">
          <div class="max-w-md w-full text-center p-12 bg-monokai-panel rounded-[2rem] border border-white/5 shadow-2xl">
            <div class="w-20 h-20 bg-monokai-orange/10 text-monokai-orange rounded-3xl flex items-center justify-center mx-auto mb-8">
              <Zap class="w-10 h-10" />
            </div>
            <h2 class="text-2xl font-bold text-white mb-3">System Ready</h2>
            <p class="text-monokai-gray text-sm mb-8 leading-relaxed">
              Select a collection on the left to manage data or create a new one to extend your schema.
            </p>
            <button @click="openCollectionDrawer" class="bg-monokai-blue text-monokai-bg px-8 py-3 rounded-2xl font-bold text-sm hover:scale-[1.02] transition-transform shadow-xl shadow-monokai-blue/10 cursor-pointer">
              Create New Collection
            </button>
          </div>
        </div>
      </main>
    </div>

    <!-- Backdrop -->
    <div v-if="showRecordDrawer || showCollectionDrawer" 
         @click="showRecordDrawer = false; showCollectionDrawer = false"
         class="fixed inset-0 bg-monokai-bg/60 backdrop-blur-sm z-40 transition-opacity">
    </div>

    <!-- Record Sidebar (Drawer) -->
    <div class="fixed right-0 top-0 h-full w-full max-w-md bg-monokai-panel border-l border-white/10 z-50 transform transition-transform duration-300 ease-in-out shadow-2xl overflow-hidden flex flex-col"
         :class="showRecordDrawer ? 'translate-x-0' : 'translate-x-full'">
      <div class="px-8 py-6 border-b border-white/5 flex items-center justify-between bg-white/5">
        <div class="flex items-center space-x-3">
          <Edit3 class="w-5 h-5 text-monokai-blue" />
          <h3 class="text-lg font-bold text-white">{{ editingRecord ? 'Edit Record' : 'New Record' }}</h3>
        </div>
        <button @click="showRecordDrawer = false" class="p-2 hover:bg-white/10 rounded-full text-monokai-gray hover:text-white transition-all cursor-pointer">
          <X class="w-5 h-5" />
        </button>
      </div>
      
      <form @submit.prevent="saveRecord" class="flex-1 flex flex-col overflow-hidden">
        <div class="flex-1 overflow-y-auto p-8 space-y-6">
          <div v-for="field in selectedCollection?.fields" :key="field.name">
            <label class="block text-[10px] font-black text-monokai-gray uppercase tracking-widest mb-2">{{ field.name }} ({{ field.type }})</label>
            <input
              v-if="field.type === 'text' || field.type === 'number'"
              v-model="recordFormData[field.name]"
              :type="field.type === 'number' ? 'number' : 'text'"
              :placeholder="'Enter ' + field.name + '...'"
              class="w-full bg-monokai-header border border-monokai-gray/30 rounded-xl p-3 text-white focus:border-monokai-blue outline-none transition-all placeholder:text-monokai-gray/50"
            />
            <select v-else-if="field.type === 'bool'" v-model="recordFormData[field.name]" class="w-full bg-monokai-header border border-monokai-gray/30 rounded-xl p-3 text-white outline-none cursor-pointer">
              <option :value="undefined" disabled>Select value...</option>
              <option :value="true">True</option>
              <option :value="false">False</option>
            </select>
          </div>
        </div>
        
        <div class="p-8 bg-white/5 border-t border-white/5 flex space-x-3">
          <button type="submit" class="flex-1 py-3 bg-monokai-blue text-monokai-bg font-bold rounded-xl hover:opacity-90 transition-all flex items-center justify-center space-x-2 cursor-pointer shadow-lg shadow-monokai-blue/10">
            <Save class="w-4 h-4" />
            <span>{{ editingRecord ? 'Update Record' : 'Create Record' }}</span>
          </button>
          <button type="button" @click="showRecordDrawer = false" class="px-6 py-3 bg-monokai-header text-white font-bold rounded-xl hover:bg-monokai-gray/20 transition-all cursor-pointer">
            Cancel
          </button>
        </div>
      </form>
    </div>

    <!-- Collection Sidebar (Drawer) -->
    <div class="fixed right-0 top-0 h-full w-full max-w-xl bg-monokai-panel border-l border-white/10 z-50 transform transition-transform duration-300 ease-in-out shadow-2xl overflow-hidden flex flex-col"
         :class="showCollectionDrawer ? 'translate-x-0' : 'translate-x-full'">
      <div class="px-8 py-6 border-b border-white/5 flex items-center justify-between bg-white/5">
        <div class="flex items-center space-x-3">
          <Plus class="w-5 h-5 text-monokai-green" />
          <h3 class="text-lg font-bold text-white">Create Collection</h3>
        </div>
        <button @click="showCollectionDrawer = false" class="p-2 hover:bg-white/10 rounded-full text-monokai-gray hover:text-white transition-all cursor-pointer">
          <X class="w-5 h-5" />
        </button>
      </div>
      
      <form @submit.prevent="saveCollection" class="flex-1 flex flex-col overflow-hidden">
        <div class="flex-1 overflow-y-auto p-8 space-y-8">
          <div class="grid grid-cols-2 gap-6">
            <div>
              <label class="block text-[10px] font-black text-monokai-gray uppercase tracking-widest mb-2">Internal Name</label>
              <input v-model="collectionFormData.name" type="text" required placeholder="e.g. products" class="w-full bg-monokai-header border border-monokai-gray/30 rounded-xl p-3 text-white focus:border-monokai-blue outline-none transition-all placeholder:text-monokai-gray/50" />
            </div>
            <div>
              <label class="block text-[10px] font-black text-monokai-gray uppercase tracking-widest mb-2">Collection Type</label>
              <select v-model="collectionFormData.type" class="w-full bg-monokai-header border border-monokai-gray/30 rounded-xl p-3 text-white outline-none cursor-pointer">
                <option value="base">Base (Generic Data)</option>
                <option value="auth">Auth (User Records)</option>
              </select>
            </div>
          </div>

          <div>
            <div class="flex justify-between items-center mb-4">
              <label class="text-[10px] font-black text-monokai-orange uppercase tracking-widest">Fields Definition</label>
              <button type="button" @click="collectionFormData.fields.push({ name: '', type: 'text' })" class="text-xs font-bold text-monokai-blue hover:underline flex items-center space-x-1 cursor-pointer">
                <Plus class="w-3 h-3" />
                <span>Add Field</span>
              </button>
            </div>
            <div class="space-y-3">
              <div v-for="(field, index) in collectionFormData.fields" :key="index" class="flex items-center space-x-3 bg-white/5 p-3 rounded-2xl border border-white/5">
                <input v-model="field.name" placeholder="field_name" class="flex-1 bg-transparent border-b border-monokai-gray/20 p-1 text-sm text-white focus:border-monokai-blue outline-none" />
                <select v-model="field.type" class="bg-monokai-header border border-monokai-gray/20 rounded-lg p-2 text-xs text-white outline-none cursor-pointer">
                  <option value="text">Text</option>
                  <option value="number">Number</option>
                  <option value="bool">Bool</option>
                  <option value="json">JSON</option>
                  <option value="file">File</option>
                </select>
                <button type="button" @click="collectionFormData.fields.splice(index, 1)" class="text-monokai-pink p-2 hover:bg-monokai-pink/10 rounded-lg transition-colors cursor-pointer">
                  <Trash2 class="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        </div>

        <div class="p-8 bg-white/5 border-t border-white/5 flex space-x-3">
          <button type="submit" class="flex-1 py-3 bg-monokai-green text-monokai-bg font-bold rounded-xl hover:opacity-90 transition-all flex items-center justify-center space-x-2 cursor-pointer shadow-lg shadow-monokai-green/10">
            <Zap class="w-4 h-4" />
            <span>Bootstrap Collection</span>
          </button>
          <button type="button" @click="showCollectionDrawer = false" class="px-6 py-3 bg-monokai-header text-white font-bold rounded-xl hover:bg-monokai-gray/20 transition-all cursor-pointer">
            Dismiss
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<style>
/* Smooth Sidebar Transition */
.translate-x-0 {
  transform: translateX(0);
}
.translate-x-full {
  transform: translateX(100%);
}
</style>
