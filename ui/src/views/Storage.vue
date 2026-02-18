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
// Sorting state
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

  // Apply sorting
  if (sortKey.value) {
    items.sort((a, b) => {
      let aValue, bValue;

      if (sortKey.value === 'size') {
        // Special handling for size
        aValue = a.size || 0;
        bValue = b.size || 0;
      } else if (sortKey.value === 'modified') {
        // Special handling for modified date
        aValue = a.modified || 0;
        bValue = b.modified || 0;
      } else if (sortKey.value === 'name') {
        // Special handling for name
        aValue = a.name || '';
        bValue = b.name || '';
      } else if (sortKey.value === 'type') {
        // Special handling for type
        aValue = a.type || '';
        bValue = b.type || '';
      } else {
        aValue = a[sortKey.value as keyof typeof a];
        bValue = b[sortKey.value as keyof typeof b];
      }

      // Handle null/undefined values
      if (aValue == null && bValue == null) return 0;
      if (aValue == null) return sortOrder.value === 'asc' ? 1 : -1;
      if (bValue == null) return sortOrder.value === 'asc' ? -1 : 1;

      // Handle numbers
      if (typeof aValue === 'number' && typeof bValue === 'number') {
        return sortOrder.value === 'asc' ? aValue - bValue : bValue - aValue;
      }

      // Handle strings
      if (typeof aValue === 'string' && typeof bValue === 'string') {
        return sortOrder.value === 'asc'
          ? aValue.localeCompare(bValue)
          : bValue.localeCompare(aValue);
      }

      // Fallback comparison
      const strA = String(aValue);
      const strB = String(bValue);
      return sortOrder.value === 'asc' ? strA.localeCompare(strB) : strB.localeCompare(strA);
    });
  }

  return items;
});

const tableHeaders = [
  { key: 'name', label: 'Name', align: 'left' as const, sortable: true },
  { key: 'size', label: 'Size', align: 'left' as const, sortable: true },
  { key: 'type', label: 'Type', align: 'left' as const, sortable: true },
  { key: 'modified', label: 'Modified', align: 'left' as const, sortable: true },
  { key: 'actions', label: 'Actions', align: 'center' as const, sticky: true },
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
    alert('Upload failed');
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
    alert(err.response?.data?.message || 'Rename failed');
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
    alert(err.response?.data?.message || 'Failed to create folder');
  }
}

async function deleteFile() {
  if (selectedPaths.value.length === 0 && !fileToDelete.value) return;

  try {
    const pathsToDelete = selectedPaths.value.length > 0 ? selectedPaths.value : [fileToDelete.value!.path];
    const recursive = true; // For simplicity in batch delete

    await axios.delete('/api/admin/storage', {
      data: {
        paths: pathsToDelete,
        recursive,
      },
    });

    showDeleteModal.value = false;
    fileToDelete.value = null;
    selectedPaths.value = [];
    loadStats();
    loadFiles(currentPath.value);
  } catch (err) {
    console.error('Delete error:', err);
    alert('Delete failed');
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
          ? `Are you sure you want to delete the folder '${fileToDelete?.name}' and all its contents? This action cannot be undone.`
          : `Are you sure you want to delete '${fileToDelete?.name}'? This action cannot be undone.`
      "
      confirm-text="Delete"
      cancel-text="Cancel"
      variant="danger"
      @confirm="deleteFile"
      @cancel="showDeleteModal = false"
    />

    <!-- New Folder Modal -->
    <Modal :show="showNewFolderModal" title="Create Folder" maxWidth="sm" @close="showNewFolderModal = false">
      <div class="space-y-4">
        <div class="space-y-2">
          <label class="block text-sm font-semibold text-text">Folder Name</label>
          <Input v-model="newFolderName" placeholder="e.g. documents" class="w-full" @keyup.enter="createFolder" />
          <p class="text-[11px] text-text-muted italic">Create a new directory in current path.</p>
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

    <!-- Rename Modal -->
    <Modal :show="showRenameModal" title="Rename" maxWidth="sm" @close="showRenameModal = false">
      <div class="space-y-4">
        <div class="space-y-2">
          <label class="block text-sm font-semibold text-text">New Name</label>
          <Input v-model="newName" placeholder="Enter new name" class="w-full" @keyup.enter="renameFile" />
          <p class="text-[11px] text-text-muted italic">Renaming '{{ fileToRename?.name }}'</p>
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

    <!-- Upload Modal -->
    <Modal :show="showUploadModal" title="Upload File" maxWidth="2xl" @close="showUploadModal = false">
      <div class="space-y-6">
        <div class="space-y-2">
          <label class="block text-sm font-semibold text-text">Target Collection</label>
          <Input v-model="uploadForm.collection" placeholder="e.g., posts" class="w-full" />
          <p class="text-[11px] text-text-muted italic">Specify which collection this file belongs to.</p>
        </div>
        <div class="space-y-2">
          <label class="block text-sm font-semibold text-text">Record ID</label>
          <Input v-model="uploadForm.recordID" placeholder="e.g., abc123" class="w-full" />
          <p class="text-[11px] text-text-muted italic">The specific record this file is attached to.</p>
        </div>
        <div class="space-y-2">
          <label class="block text-sm font-semibold text-text">Choose File</label>
          <div class="relative">
            <input
              type="file"
              class="w-full text-sm text-text-muted file:mr-4 file:py-2.5 file:px-4 file:rounded-lg file:border-0 file:text-xs file:font-bold file:bg-primary/10 file:text-primary hover:file:bg-primary/20 file:transition-colors file:cursor-pointer border border-border rounded-lg p-1 bg-surface-dark/30"
              @change="handleFileSelect"
            />
          </div>
        </div>
        <div class="pt-2">
          <Checkbox v-model="uploadForm.preserveName" label="Keep original filename" />
          <p class="text-[11px] text-text-muted mt-1 ml-7 italic">If unchecked, a random unique name will be generated.</p>
        </div>
        <div v-if="uploadProgress > 0" class="space-y-2">
          <div class="flex justify-between text-xs font-medium text-primary">
            <span>Uploading...</span>
            <span>{{ uploadProgress }}%</span>
          </div>
          <div class="h-1.5 bg-surface-hover rounded-full overflow-hidden border border-border/50">
            <div
              class="h-full bg-primary transition-all duration-300 ease-out shadow-[0_0_10px_rgba(var(--color-primary),0.5)]"
              :style="{ width: uploadProgress + '%' }"
            ></div>
          </div>
        </div>
      </div>

      <template #footer>
        <Button variant="secondary" size="sm" @click="showUploadModal = false">Cancel</Button>
        <Button :disabled="!canUpload || uploadProgress > 0" size="sm" @click="uploadFile">
          <template #leftIcon><Upload class="w-4 h-4" /></template>
          Start Upload
        </Button>
      </template>
    </Modal>

    <!-- Header -->
    <AppHeader>
      <template #breadcrumb>
        <div class="flex items-center text-sm text-text-muted truncate gap-2">
          <span
            :class="[
              !currentPath
                ? 'font-medium text-primary'
                : 'hover:text-primary cursor-pointer transition-colors duration-200',
            ]"
            @click="navigateTo('')"
            >Storage</span
          >
          <template v-for="(part, index) in pathParts" :key="index">
            <span class="text-text-muted flex-shrink-0">/</span>
            <span
              :class="[
                index === pathParts.length - 1
                  ? 'font-medium text-primary'
                  : 'hover:text-primary cursor-pointer transition-colors duration-200',
              ]"
              class="truncate"
              @click="navigateTo(pathParts.slice(0, index + 1).join('/'))"
            >
              {{ part }}
            </span>
          </template>
        </div>
      </template>
    </AppHeader>

    <!-- Main Content -->
    <div class="flex-1 overflow-auto min-h-0 p-4 sm:p-8 pb-24 sm:pb-8">
      <div class="max-w-7xl mx-auto space-y-6 sm:space-y-8">
        <!-- Page Title -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h1 class="text-2xl font-bold text-text tracking-tight">Storage Bucket</h1>
            <p class="mt-1 text-sm text-text-muted">Manage uploaded files and media assets.</p>
          </div>
          <div class="flex gap-3">
            <Button
              v-if="selectedPaths.length > 0"
              variant="destructive"
              size="sm"
              @click="showDeleteModal = true"
            >
              <template #leftIcon><Trash2 class="w-4 h-4" /></template>
              Delete Selected ({{ selectedPaths.length }})
            </Button>
            <Button variant="secondary" size="sm" @click="showNewFolderModal = true">
              <template #leftIcon><FolderPlus class="w-4 h-4" /></template>
              New Folder
            </Button>
            <Button size="sm" @click="showUploadModal = true">
              <template #leftIcon><Upload class="w-4 h-4" /></template>
              Upload File
            </Button>
          </div>
        </div>

        <!-- Stats Cards -->
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div class="bg-surface border border-border rounded-lg p-6">
            <div class="text-3xl font-bold text-primary">{{ stats.total_files }}</div>
            <div class="text-sm text-text-muted mt-1">Total Files</div>
          </div>
          <div class="bg-surface border border-border rounded-lg p-6">
            <div class="text-3xl font-bold text-primary">{{ formatSize(stats.total_size) }}</div>
            <div class="text-sm text-text-muted mt-1">Total Size</div>
          </div>
          <div class="bg-surface border border-border rounded-lg p-6">
            <div class="text-3xl font-bold text-primary">{{ stats.total_collections }}</div>
            <div class="text-sm text-text-muted mt-1">Collections</div>
          </div>
        </div>

        <!-- File List -->
        <Table
          :headers="tableHeaders"
          :items="tableItems"
          :loading="loading"
          empty-text="No files or folders"
          :enable-pagination="true"
          :default-page-size="15"
          :sort-key="sortKey"
          :sort-order="sortOrder"
          row-clickable
          selectable
          selection-key="path"
          @selection-change="selectedPaths = $event"
          @sort-change="
            (key, order) => {
              sortKey = key;
              sortOrder = order;
            }
          "
          @row-click="
            (item, event) => {
              const target = event.target as HTMLElement;
              if (!target.closest('.actions-cell')) {
                handleRowClick(item as Record<string, unknown>);
              }
            }
          "
        >
          <template #cell(name)="{ item }">
            <div
              class="flex items-center"
              :class="item.isFolder ? 'cursor-pointer hover:text-primary' : ''"
              @click="item.isFolder && navigateTo(item.path as string)"
            >
              <FolderOpen v-if="item.isFolder" class="w-5 h-5 text-primary mr-3" />
              <File v-else class="w-5 h-5 text-text-muted mr-3" />
              <span class="text-sm font-medium text-text">{{ item.name }}</span>
            </div>
          </template>

          <template #cell(size)="{ item }">
            <span class="text-sm text-text-muted">
              {{ item.isFolder ? '-' : formatSize(item.size as number) }}
            </span>
          </template>

          <template #cell(type)="{ item }">
            <span class="text-sm text-text-muted">{{ item.type }}</span>
          </template>

          <template #cell(modified)="{ item }">
            <span class="text-sm text-text-muted">{{ formatDate(item.modified as number) }}</span>
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
                  <PopoverItem
                    v-if="!item.isFolder"
                    :icon="Download"
                    @click="
                      close();
                      downloadFile(item as unknown as FileInfo);
                    "
                  >
                    Download
                  </PopoverItem>
                  <PopoverItem v-if="item.isFolder" :icon="Settings" @click="close()">
                    Settings
                  </PopoverItem>
                  <PopoverItem
                    :icon="Pencil"
                    @click="
                      close();
                      confirmRename(item as unknown as FileInfo);
                    "
                  >
                    Rename
                  </PopoverItem>
                  <PopoverItem
                    v-if="!item.isFolder"
                    :icon="Trash2"
                    variant="danger"
                    @click="
                      close();
                      confirmDelete(item as unknown as FileInfo);
                    "
                  >
                    Delete
                  </PopoverItem>
                  <PopoverItem
                    v-if="item.isFolder"
                    :icon="Trash2"
                    variant="danger"
                    @click="
                      close();
                      confirmDelete(item as unknown as FileInfo);
                    "
                  >
                    Delete
                  </PopoverItem>
                </template>
              </Popover>
            </div>
          </template>
        </Table>
      </div>
    </div>
  </AppLayout>
</template>
