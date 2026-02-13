<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import AppLayout from '../components/AppLayout.vue';
import AppHeader from '../components/AppHeader.vue';
import Button from '../components/Button.vue';
import Table from '../components/Table.vue';
import ConfirmModal from '../components/ConfirmModal.vue';
import Input from '../components/Input.vue';
import { Upload, Download, Trash2, FolderOpen, File } from 'lucide-vue-next';

const router = useRouter();

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

const tableData = computed(() => {
  const items = [
    ...folders.value.map(f => ({ ...f, type: 'folder' })),
    ...files.value.map(f => ({ ...f, type: getFileType(f.mime_type || '') }))
  ];
  return items;
});

const headers = [
  { key: 'name', label: 'Name' },
  { key: 'size', label: 'Size' },
  { key: 'type', label: 'Type' },
  { key: 'modified', label: 'Modified' },
  { key: 'actions', label: 'Actions' },
];

onMounted(() => {
  loadStats();
  loadFiles('');
});

async function loadStats() {
  try {
    const response = await axios.get('/api/admin/storage/stats');
    stats.value = response.data;
  } catch (err) {
    console.error('Failed to load stats:', err);
  }
}

async function loadFiles(path: string) {
  loading.value = true;
  try {
    const url = `/api/admin/storage?path=${encodeURIComponent(path)}`;
    const response = await axios.get(url);
    folders.value = response.data.folders || [];
    files.value = response.data.files || [];
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
    await axios.delete('/api/admin/storage', {
      data: { path: fileToDelete.value.path },
    });
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
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
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

function getFileIcon(mimeType: string): string {
  if (mimeType.startsWith('image/')) return '🖼️';
  if (mimeType.startsWith('video/')) return '🎥';
  if (mimeType.startsWith('audio/')) return '🎵';
  if (mimeType === 'application/pdf') return '📄';
  if (mimeType.includes('zip')) return '📦';
  return '📄';
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
      title="Delete File"
      :message="`Are you sure you want to delete ${fileToDelete?.name}? This action cannot be undone.`"
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
          <button
            @click="showUploadModal = false"
            class="text-text-muted hover:text-text"
          >
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
              @change="handleFileSelect"
              class="w-full text-sm text-text-muted file:mr-4 file:py-2 file:px-4 file:rounded file:border-0 file:text-sm file:font-semibold file:bg-primary/10 file:text-primary hover:file:bg-primary/20"
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
        <div class="flex items-center text-sm text-text-muted">
          <span class="hover:text-text cursor-pointer" @click="router.push('/')">Vault</span>
          <span class="mx-2">/</span>
          <span class="font-medium text-text">Storage</span>
        </div>
      </template>
    </AppHeader>

    <!-- Main Content -->
    <div class="flex-1 overflow-auto min-h-0 p-4 sm:p-8 pb-24 sm:pb-8">
      <div class="max-w-7xl mx-auto space-y-6 sm:space-y-8">
        <!-- Page Title -->
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h1 class="text-2xl font-bold text-text tracking-tight">Storage Browser</h1>
            <p class="mt-1 text-sm text-text-muted">
              Manage uploaded files and media assets.
            </p>
          </div>
          <Button @click="showUploadModal = true">
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

        <!-- Breadcrumb -->
        <div class="flex items-center text-sm text-text-muted">
          <span @click="navigateTo('')" class="hover:text-text cursor-pointer">Storage</span>
          <template v-for="(part, index) in pathParts" :key="index">
            <span class="mx-2">/</span>
            <span
              @click="navigateTo(pathParts.slice(0, index + 1).join('/'))"
              class="hover:text-text cursor-pointer"
            >
              {{ part }}
            </span>
          </template>
        </div>

        <!-- File List -->
        <div class="bg-surface border border-border rounded-lg overflow-hidden">
          <div v-if="loading" class="p-12 text-center text-text-muted">
            Loading...
          </div>

          <div
            v-else-if="folders.length === 0 && files.length === 0"
            class="p-12 text-center text-text-muted"
          >
            <p>No files or folders</p>
          </div>

          <div v-else class="overflow-x-auto">
            <table class="w-full">
              <thead class="bg-surface-hover border-b border-border">
                <tr>
                  <th class="text-left px-6 py-3 text-xs font-medium text-text-muted uppercase tracking-wider">
                    Name
                  </th>
                  <th class="text-left px-6 py-3 text-xs font-medium text-text-muted uppercase tracking-wider">
                    Size
                  </th>
                  <th class="text-left px-6 py-3 text-xs font-medium text-text-muted uppercase tracking-wider">
                    Type
                  </th>
                  <th class="text-left px-6 py-3 text-xs font-medium text-text-muted uppercase tracking-wider">
                    Modified
                  </th>
                  <th class="text-right px-6 py-3 text-xs font-medium text-text-muted uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody class="divide-y divide-border">
                <tr
                  v-for="folder in folders"
                  :key="folder.path"
                  @click="navigateTo(folder.path)"
                  class="hover:bg-surface-hover cursor-pointer"
                >
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="flex items-center">
                      <FolderOpen class="w-5 h-5 text-primary mr-3" />
                      <span class="text-sm font-medium text-text">{{ folder.name }}</span>
                    </div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-text-muted">-</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-text-muted">folder</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-text-muted">
                    {{ formatDate(folder.modified) }}
                  </td>
                  <td></td>
                </tr>
                <tr
                  v-for="file in files"
                  :key="file.path"
                  class="hover:bg-surface-hover"
                >
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="flex items-center">
                      <File class="w-5 h-5 text-text-muted mr-3" />
                      <span class="text-sm font-medium text-text">{{ file.name }}</span>
                    </div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-text-muted">
                    {{ formatSize(file.size) }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-text-muted">
                    {{ getFileType(file.mime_type || '') }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-text-muted">
                    {{ formatDate(file.modified) }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-right text-sm">
                    <div class="flex items-center justify-end gap-2">
                      <Button
                        variant="ghost"
                        size="xs"
                        @click="downloadFile(file)"
                        title="Download"
                      >
                        <template #leftIcon><Download class="w-4 h-4" /></template>
                      </Button>
                      <Button
                        variant="ghost"
                        size="xs"
                        @click="confirmDelete(file)"
                        title="Delete"
                      >
                        <template #leftIcon><Trash2 class="w-4 h-4 text-red-500" /></template>
                      </Button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </AppLayout>
</template>
