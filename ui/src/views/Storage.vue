<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import Table from '../components/Table.vue';
import ConfirmModal from '../components/ConfirmModal.vue';
import Input from '../components/Input.vue';
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
const fileToDelete = ref<FileInfo | null>(null);
const uploadForm = ref({ collection: '', recordID: '', file: null as File | null });
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

  try {
    uploadProgress.value = 50;
    await axios.post('/api/files', formData);
    uploadProgress.value = 100;
    showUploadModal.value = false;
    uploadForm.value = { collection: '', recordID: '', file: null };
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

async function deleteFile() {
  if (!fileToDelete.value) return;

  try {
    const data: { path: string; recursive?: boolean } = { path: fileToDelete.value.path };

    // If it's a directory, add recursive flag
    if (fileToDelete.value.isFolder) {
      data.recursive = true;
    }

    await axios.delete('/api/admin/storage', { data });
    showDeleteModal.value = false;
    fileToDelete.value = null;
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

    <!-- Upload Modal -->
    <div
      v-if="showUploadModal"
      class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4"
      @click.self="showUploadModal = false"
    >
      <div class="bg-surface rounded-lg w-full max-w-md border border-border">
        <div class="flex items-center justify-between p-6 border-b border-border">
          <h2 class="text-lg font-semibold text-text">Upload File</h2>
          <button class="text-text-muted hover:text-text" @click="showUploadModal = false">
            ×
          </button>
        </div>
        <div class="p-6 space-y-4">
          <div>
            <label class="block text-sm font-medium text-text-muted mb-2">Collection</label>
            <Input v-model="uploadForm.collection" placeholder="e.g., posts" />
          </div>
          <div>
            <label class="block text-sm font-medium text-text-muted mb-2">Record ID</label>
            <Input v-model="uploadForm.recordID" placeholder="e.g., abc123" />
          </div>
          <div>
            <label class="block text-sm font-medium text-text-muted mb-2">File</label>
            <input
              type="file"
              class="w-full text-sm text-text-muted file:mr-4 file:py-2 file:px-4 file:rounded file:border-0 file:text-sm file:font-semibold file:bg-primary/10 file:text-primary hover:file:bg-primary/20"
              @change="handleFileSelect"
            />
          </div>
          <div v-if="uploadProgress > 0" class="h-1 bg-surface-hover rounded overflow-hidden">
            <div
              class="h-full bg-primary transition-all"
              :style="{ width: uploadProgress + '%' }"
            ></div>
          </div>
        </div>
        <div class="flex justify-end gap-3 p-6 border-t border-border">
          <Button variant="secondary" @click="showUploadModal = false">Cancel</Button>
          <Button :disabled="!canUpload" @click="uploadFile">
            <template #leftIcon><Upload class="w-4 h-4" /></template>
            Upload
          </Button>
        </div>
      </div>
    </div>

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
          <Button size="sm" @click="showUploadModal = true">
            <template #leftIcon><Upload class="w-4 h-4" /></template>
            Upload File
          </Button>
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
