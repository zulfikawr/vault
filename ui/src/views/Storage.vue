<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import Table from '../components/Table.vue';
import ConfirmModal from '../components/ConfirmModal.vue';
import Modal from '../components/Modal.vue';
import Input from '../components/Input.vue';
import Checkbox from '../components/Checkbox.vue';
import Popover from '../components/Popover.vue';
import PopoverItem from '../components/PopoverItem.vue';
import {
  Upload,
  Download,
  Trash2,
  FolderOpen,
  File,
  MoreHorizontal,
  Settings,
  Pencil,
  FolderPlus,
  Cloud,
  ChevronRight,
  Search,
  LayoutGrid,
  List
} from 'lucide-vue-next';

interface FileInfo {
  name: string;
  path: string;
  size: number;
  is_dir: boolean;
  modified: number;
  mime_type?: string;
}

interface StorageStats {
  total_files: number;
  total_size: number;
  total_collections: number;
}

const currentPath = ref('');
const folders = ref<FileInfo[]>([]);
const files = ref<FileInfo[]>([]);
const sortKey = ref('name');
const sortOrder = ref<'asc' | 'desc'>('asc');
const stats = ref<StorageStats>({ total_files: 0, total_size: 0, total_collections: 0 });
const loading = ref(false);
const showUploadModal = ref(false);
const showDeleteModal = ref(false);
const showRenameModal = ref(false);
const showNewFolderModal = ref(false);
const fileToDelete = ref<FileInfo | null>(null);
const selectedPaths = ref<string[]>([]);
const fileToRename = ref<FileInfo | null>(null);
const newName = ref('');
const newFolderName = ref('');
const uploadForm = ref({ collection: '', recordID: '', file: null as File | null, preserveName: true });
const uploadProgress = ref(0);

const pathParts = computed(() => {
  return currentPath.value ? currentPath.value.split('/').filter(Boolean) : [];
});

const canUpload = computed(() => {
  return uploadForm.value.collection && uploadForm.value.recordID && uploadForm.value.file;
});

const tableItems = computed(() => {
  let items = [
    ...folders.value.map((f) => ({ ...f, type: 'folder', isFolder: true })),
    ...files.value.map((f) => ({ ...f, type: getFileType(f.mime_type || ''), isFolder: false })),
  ];

  if (sortKey.value) {
    items.sort((a, b) => {
      let aValue, bValue;
      if (sortKey.value === 'size') {
        aValue = a.size || 0;
        bValue = b.size || 0;
      } else if (sortKey.value === 'modified') {
        aValue = a.modified || 0;
        bValue = b.modified || 0;
      } else if (sortKey.value === 'name') {
        aValue = a.name || '';
        bValue = b.name || '';
      } else if (sortKey.value === 'type') {
        aValue = a.type || '';
        bValue = b.type || '';
      } else {
        aValue = a[sortKey.value as keyof typeof a];
        bValue = b[sortKey.value as keyof typeof b];
      }
      if (aValue == null && bValue == null) return 0;
      if (aValue == null) return sortOrder.value === 'asc' ? 1 : -1;
      if (bValue == null) return sortOrder.value === 'asc' ? -1 : 1;
      if (typeof aValue === 'number' && typeof bValue === 'number') {
        return sortOrder.value === 'asc' ? aValue - bValue : bValue - aValue;
      }
      if (typeof aValue === 'string' && typeof bValue === 'string') {
        return sortOrder.value === 'asc' ? aValue.localeCompare(bValue) : bValue.localeCompare(aValue);
      }
      return sortOrder.value === 'asc' ? String(aValue).localeCompare(String(bValue)) : String(bValue).localeCompare(String(aValue));
    });
  }
  return items;
});

const tableHeaders = [
  { key: 'name', label: 'Name', align: 'left' as const, sortable: true },
  { key: 'size', label: 'Size', align: 'left' as const, sortable: true },
  { key: 'type', label: 'Type', align: 'left' as const, sortable: true },
  { key: 'modified', label: 'Modified', align: 'left' as const, sortable: true },
  { key: 'actions', label: '', align: 'center' as const, sticky: true },
];

onMounted(() => {
  loadStats();
  loadFiles('');
});

async function loadStats() {
  try {
    const response = await axios.get('/api/admin/storage/stats');
    stats.value = response.data.data;
  } catch (err) {
    console.error('Failed to load stats:', err);
  }
}

async function loadFiles(path: string) {
  loading.value = true;
  try {
    const url = `/api/admin/storage?path=${encodeURIComponent(path)}`;
    const response = await axios.get(url);
    folders.value = response.data.data.folders || [];
    files.value = response.data.data.files || [];
  } catch (err) {
    console.error('Failed to load files:', err);
  } finally {
    loading.value = false;
  }
}

function navigateTo(path: string) {
  currentPath.value = path;
  loadFiles(path);
}

function handleRowClick(item: Record<string, unknown>) {
  if (item.isFolder) {
    navigateTo(item.path as string);
  }
}

function handleFileSelect(event: Event) {
  const target = event.target as HTMLInputElement;
  if (target.files && target.files[0]) {
    uploadForm.value.file = target.files[0];
  }
}

async function uploadFile() {
  if (!canUpload.value) return;
  const formData = new FormData();
  formData.append('file', uploadForm.value.file!);
  formData.append('collection', uploadForm.value.collection);
  formData.append('recordID', uploadForm.value.recordID);
  formData.append('preserve_name', String(uploadForm.value.preserveName));
  try {
    uploadProgress.value = 50;
    await axios.post('/api/files', formData);
    uploadProgress.value = 100;
    showUploadModal.value = false;
    uploadForm.value = { collection: '', recordID: '', file: null, preserveName: true };
    uploadProgress.value = 0;
    loadStats();
    loadFiles(currentPath.value);
  } catch (err) {
    console.error('Upload error:', err);
  }
}

function downloadFile(file: FileInfo) {
  const url = `/api/files/${file.path}`;
  window.open(url, '_blank');
}

function confirmDelete(file: FileInfo) {
  fileToDelete.value = file;
  showDeleteModal.value = true;
}

function confirmRename(file: FileInfo) {
  fileToRename.value = file;
  newName.value = file.name;
  showRenameModal.value = true;
}

async function renameFile() {
  if (!fileToRename.value || !newName.value || newName.value === fileToRename.value.name) {
    showRenameModal.value = false;
    return;
  }
  try {
    await axios.post('/api/admin/storage/rename', {
      old_path: fileToRename.value.path,
      new_name: newName.value,
    });
    showRenameModal.value = false;
    fileToRename.value = null;
    loadFiles(currentPath.value);
  } catch (err: any) {
    console.error('Rename error:', err);
  }
}

async function createFolder() {
  if (!newFolderName.value) return;
  try {
    await axios.post('/api/admin/storage/mkdir', {
      path: currentPath.value,
      name: newFolderName.value,
    });
    showNewFolderModal.value = false;
    newFolderName.value = '';
    loadFiles(currentPath.value);
  } catch (err: any) {
    console.error('Create folder error:', err);
  }
}

async function deleteFile() {
  if (selectedPaths.value.length === 0 && !fileToDelete.value) return;
  try {
    const pathsToDelete = selectedPaths.value.length > 0 ? selectedPaths.value : [fileToDelete.value!.path];
    await axios.delete('/api/admin/storage', {
      data: { paths: pathsToDelete, recursive: true },
    });
    showDeleteModal.value = false;
    fileToDelete.value = null;
    selectedPaths.value = [];
    loadStats();
    loadFiles(currentPath.value);
  } catch (err) {
    console.error('Delete error:', err);
  }
}

function formatSize(bytes: number): string {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i];
}

function formatDate(timestamp: number): string {
  const date = new Date(timestamp * 1000);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const hours = Math.floor(diff / 3600000);
  const days = Math.floor(diff / 86400000);
  if (hours < 1) return 'Just now';
  if (hours < 24) return `${hours}h ago`;
  if (days < 7) return `${days}d ago`;
  return date.toLocaleDateString();
}

function getFileType(mimeType: string): string {
  if (mimeType.startsWith('image/')) return 'image';
  if (mimeType.startsWith('video/')) return 'video';
  if (mimeType.startsWith('audio/')) return 'audio';
  if (mimeType === 'application/pdf') return 'pdf';
  return 'file';
}
</script>

<template>
  <AppLayout>
    <ConfirmModal
      :show="showDeleteModal"
      :title="fileToDelete?.isFolder ? 'Delete Folder' : 'Delete File'"
      :message="
        fileToDelete?.isFolder
          ? `Are you sure you want to delete '${fileToDelete?.name}' and all its contents? This action cannot be undone.`
          : `Are you sure you want to delete '${fileToDelete?.name}'? This action cannot be undone.`
      "
      confirm-text="Delete"
      cancel-text="Cancel"
      variant="danger"
      @confirm="deleteFile"
      @cancel="showDeleteModal = false"
    />

    <Modal :show="showNewFolderModal" title="New Folder" maxWidth="sm" @close="showNewFolderModal = false">
      <div class="space-y-4">
        <div class="space-y-2">
          <label class="text-[10px] font-bold text-text-dim ml-1">Folder Name</label>
          <Input v-model="newFolderName" placeholder="e.g. documents" class="w-full" @keyup.enter="createFolder" />
        </div>
      </div>
      <template #footer>
        <Button variant="secondary" size="sm" @click="showNewFolderModal = false">Cancel</Button>
        <Button :disabled="!newFolderName" size="sm" @click="createFolder">
          <template #leftIcon><FolderPlus class="w-4 h-4" /></template>
          Create
        </Button>
      </template>
    </Modal>

    <Modal :show="showRenameModal" title="Rename" maxWidth="sm" @close="showRenameModal = false">
      <div class="space-y-4">
        <div class="space-y-2">
          <label class="text-[10px] font-bold text-text-dim ml-1">New Name</label>
          <Input v-model="newName" placeholder="Enter new name" class="w-full" @keyup.enter="renameFile" />
        </div>
      </div>
      <template #footer>
        <Button variant="secondary" size="sm" @click="showRenameModal = false">Cancel</Button>
        <Button :disabled="!newName || newName === fileToRename?.name" size="sm" @click="renameFile">
          <template #leftIcon><Pencil class="w-4 h-4" /></template>
          Rename
        </Button>
      </template>
    </Modal>

    <Modal :show="showUploadModal" title="Upload File" maxWidth="lg" @close="showUploadModal = false">
      <div class="space-y-6">
        <div class="grid grid-cols-2 gap-4">
          <div class="space-y-2">
            <label class="text-[10px] font-bold text-text-dim ml-1">Target Collection</label>
            <Input v-model="uploadForm.collection" placeholder="e.g. posts" />
          </div>
          <div class="space-y-2">
            <label class="text-[10px] font-bold text-text-dim ml-1">Record ID</label>
            <Input v-model="uploadForm.recordID" placeholder="e.g. abc123" />
          </div>
        </div>
        <div class="space-y-2">
          <label class="text-[10px] font-bold text-text-dim ml-1">Choose File</label>
          <div class="relative group">
            <input
              type="file"
              class="w-full text-xs text-text-muted file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:text-[10px] file:font-bold file:bg-primary/10 file:text-primary hover:file:bg-primary/20 transition-all cursor-pointer border border-border/50 rounded-xl p-1 bg-surface-dark/30"
              @change="handleFileSelect"
            />
          </div>
        </div>
        <div class="flex items-center gap-2">
          <Checkbox v-model="uploadForm.preserveName" label="Keep original filename" />
        </div>
        <div v-if="uploadProgress > 0" class="space-y-2">
          <div class="flex justify-between text-[10px] font-bold text-primary">
            <span>Uploading...</span>
            <span>{{ uploadProgress }}%</span>
          </div>
          <div class="h-1 bg-surface-dark rounded-full overflow-hidden">
            <div class="h-full bg-primary transition-all duration-300" :style="{ width: uploadProgress + '%' }"></div>
          </div>
        </div>
      </div>
      <template #footer>
        <Button variant="secondary" size="sm" @click="showUploadModal = false">Cancel</Button>
        <Button :disabled="!canUpload || uploadProgress > 0" size="sm" @click="uploadFile">
          <template #leftIcon><Upload class="w-4 h-4" /></template>
          Upload
        </Button>
      </template>
    </Modal>

    <AppHeader>
      <template #breadcrumb>
        <div class="flex items-center gap-2 text-xs font-bold tracking-wide">
          <Cloud class="w-4 h-4 text-primary" />
          <span class="text-text-dim hover:text-primary transition-colors cursor-pointer" @click="navigateTo('')">Storage</span>
          <template v-for="(part, index) in pathParts" :key="index">
            <ChevronRight class="w-3 h-3 text-text-dim" />
            <span
              :class="index === pathParts.length - 1 ? 'text-primary' : 'text-text-dim hover:text-primary cursor-pointer transition-colors'"
              @click="navigateTo(pathParts.slice(0, index + 1).join('/'))"
            >{{ part }}</span>
          </template>
        </div>
      </template>
    </AppHeader>

    <main class="flex-1 min-h-0 p-4 lg:p-8 max-w-7xl mx-auto w-full">
      <div class="space-y-8">
        <!-- Header -->
        <div class="flex flex-col md:flex-row md:items-end justify-between gap-6 border-b border-border/30 pb-6">
          <div>
            <h1 class="text-3xl font-bold text-text tracking-tight">Storage Bucket</h1>
            <p class="text-text-dim text-xs mt-1">Manage uploaded files and media assets.</p>
          </div>
          <div class="flex gap-3">
            <Button v-if="selectedPaths.length > 0" variant="destructive" size="sm" @click="showDeleteModal = true">
              <Trash2 class="w-4 h-4 mr-2" />
              Delete Selected ({{ selectedPaths.length }})
            </Button>
            <Button variant="secondary" size="sm" @click="showNewFolderModal = true">
              <FolderPlus class="w-4 h-4 mr-2" />
              New Folder
            </Button>
            <Button size="sm" @click="showUploadModal = true">
              <Upload class="w-4 h-4 mr-2" />
              Upload File
            </Button>
          </div>
        </div>

        <!-- Stats Grid -->
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div class="bg-surface/40 backdrop-blur-sm border border-border/50 p-5 rounded-2xl">
            <p class="text-[10px] font-bold text-text-dim tracking-wider mb-1">Total Files</p>
            <div class="flex items-baseline gap-2">
              <span class="text-2xl font-bold text-text">{{ stats.total_files }}</span>
            </div>
          </div>
          <div class="bg-surface/40 backdrop-blur-sm border border-border/50 p-5 rounded-2xl">
            <p class="text-[10px] font-bold text-text-dim tracking-wider mb-1">Total Size</p>
            <div class="flex items-baseline gap-2">
              <span class="text-2xl font-bold text-text">{{ formatSize(stats.total_size) }}</span>
            </div>
          </div>
          <div class="bg-surface/40 backdrop-blur-sm border border-border/50 p-5 rounded-2xl">
            <p class="text-[10px] font-bold text-text-dim tracking-wider mb-1">Collections</p>
            <div class="flex items-baseline gap-2">
              <span class="text-2xl font-bold text-text">{{ stats.total_collections }}</span>
            </div>
          </div>
        </div>

        <!-- Toolbar -->
        <div class="flex flex-col sm:flex-row gap-4">
          <div class="relative flex-1">
            <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-text-dim" />
            <input 
              type="text" 
              placeholder="Search files..." 
              class="w-full bg-surface-dark/50 border border-border/50 rounded-xl pl-10 pr-4 h-10 text-sm focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary/50 transition-all"
            />
          </div>
          <div class="flex bg-surface-dark/50 p-1 rounded-xl border border-border/50 backdrop-blur-sm">
            <button class="p-1.5 rounded-lg bg-primary text-white shadow-sm"><List class="w-4 h-4" /></button>
            <button class="p-1.5 rounded-lg text-text-dim hover:text-text transition-colors"><LayoutGrid class="w-4 h-4" /></button>
          </div>
        </div>

        <!-- File Table -->
        <div class="animate-in fade-in slide-in-from-bottom-4 duration-500">
          <Table
            :headers="tableHeaders"
            :items="tableItems"
            :loading="loading"
            empty-text="No files or folders found"
            :enable-pagination="true"
            :default-page-size="15"
            :sort-key="sortKey"
            :sort-order="sortOrder"
            row-clickable
            selectable
            selection-key="path"
            @selection-change="selectedPaths = $event"
            @sort-change="(key, order) => { sortKey = key; sortOrder = order; }"
            @row-click="(item, event) => { if (!(event.target as HTMLElement).closest('.actions-cell')) handleRowClick(item as any); }"
          >
            <template #cell(name)="{ item }">
              <div class="flex items-center py-1 group/item">
                <div class="w-8 h-8 rounded-lg bg-surface-dark border border-border flex items-center justify-center mr-3 group-hover/item:border-primary/30 transition-colors">
                  <FolderOpen v-if="item.isFolder" class="w-4 h-4 text-primary" />
                  <File v-else class="w-4 h-4 text-text-muted" />
                </div>
                <span :class="item.isFolder ? 'font-bold text-text cursor-pointer hover:text-primary' : 'text-text-muted'" class="text-xs transition-colors">
                  {{ item.name }}
                </span>
              </div>
            </template>
            <template #cell(size)="{ item }">
              <span class="text-xs text-text-dim">{{ item.isFolder ? '-' : formatSize(item.size as number) }}</span>
            </template>
            <template #cell(type)="{ item }">
              <span class="px-2 py-0.5 rounded bg-white/5 border border-white/5 text-[9px] font-bold text-text-muted">
                {{ item.isFolder ? 'folder' : item.type }}
              </span>
            </template>
            <template #cell(modified)="{ item }">
              <span class="text-[10px] font-medium text-text-dim">{{ formatDate(item.modified as number) }}</span>
            </template>
            <template #cell(actions)="{ item }">
              <div class="actions-cell">
                <Popover align="right">
                  <template #trigger>
                    <Button variant="ghost" size="xs">
                      <MoreHorizontal class="w-4 h-4" />
                    </Button>
                  </template>
                  <template #default="{ close }">
                    <PopoverItem v-if="!item.isFolder" :icon="Download" @click="close(); downloadFile(item as any);">Download</PopoverItem>
                    <PopoverItem :icon="Pencil" @click="close(); confirmRename(item as any);">Rename</PopoverItem>
                    <PopoverItem :icon="Trash2" variant="danger" @click="close(); confirmDelete(item as any);">Delete</PopoverItem>
                  </template>
                </Popover>
              </div>
            </template>
          </Table>
        </div>
      </div>
    </main>
  </AppLayout>
</template>
