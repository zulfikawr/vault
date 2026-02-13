# Storage Web UI Design

## Overview
A file browser interface for managing uploaded files in `vault_data/storage/`.

## Route
- `/storage` - Main storage browser view

## Layout Structure

```
┌─────────────────────────────────────────────────────────┐
│ Storage Browser                          [Upload File]  │
├─────────────────────────────────────────────────────────┤
│                                                          │
│ 📁 Breadcrumb: Storage / posts / abc123                │
│                                                          │
│ ┌──────────────────────────────────────────────────┐   │
│ │ Search files...                            [🔍]  │   │
│ └──────────────────────────────────────────────────┘   │
│                                                          │
│ ┌────────────────────────────────────────────────────┐ │
│ │ Name              Size      Type      Modified     │ │
│ ├────────────────────────────────────────────────────┤ │
│ │ 📁 posts          -         folder    2h ago      │ │
│ │ 📁 users          -         folder    1d ago      │ │
│ └────────────────────────────────────────────────────┘ │
│                                                          │
│ Inside folder:                                           │
│ ┌────────────────────────────────────────────────────┐ │
│ │ 🖼️  image.jpg     2.3 MB    image     2h ago  [⋮] │ │
│ │ 📄 doc.pdf        1.1 MB    pdf       1d ago  [⋮] │ │
│ │ 🎵 audio.mp3      5.2 MB    audio     3d ago  [⋮] │ │
│ └────────────────────────────────────────────────────┘ │
│                                                          │
│ Total: 3 files, 8.6 MB                                  │
└─────────────────────────────────────────────────────────┘
```

## Features

### 1. File Browser
- **Folder Navigation**
  - Click folders to navigate into them
  - Breadcrumb navigation at top
  - Back button
  - Show folder hierarchy: storage → collection → recordID → files

### 2. File List
- **Display Info**
  - Icon based on file type (image, pdf, audio, video, document)
  - Filename
  - File size (human readable: KB, MB, GB)
  - MIME type
  - Last modified date
  - Actions menu (⋮)

- **Sorting**
  - By name (A-Z, Z-A)
  - By size (smallest/largest)
  - By date (newest/oldest)
  - By type

### 3. File Actions (⋮ menu)
- **View/Download** - Open file in new tab or download
- **Copy URL** - Copy file URL to clipboard
- **Rename** - Rename file
- **Delete** - Delete file (with confirmation)
- **Details** - Show full metadata

### 4. Upload
- **Upload Button** - Opens upload modal
- **Drag & Drop** - Drop files anywhere on the page
- **Upload Modal**
  - Select collection (dropdown)
  - Select/create recordID (input)
  - Choose files (multiple)
  - Progress bar for each file
  - Cancel upload

### 5. Search & Filter
- **Search** - Filter by filename
- **Filter by type** - Images, Documents, Audio, Video, Other
- **Filter by size** - < 1MB, 1-10MB, > 10MB
- **Filter by date** - Today, This week, This month, Older

### 6. Bulk Actions
- **Select multiple files** - Checkboxes
- **Bulk delete** - Delete selected files
- **Bulk download** - Download as ZIP

### 7. Storage Stats (Top Cards)
```
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│ Total Files  │ │ Total Size   │ │ Collections  │
│    1,234     │ │   2.4 GB     │ │      12      │
└──────────────┘ └──────────────┘ └──────────────┘
```

## API Endpoints Needed

### Backend (New)
- `GET /api/admin/storage` - List all files/folders
  - Query params: `path`, `sort`, `filter`
  - Returns: `{ files: [], folders: [], stats: {} }`
- `GET /api/admin/storage/stats` - Storage statistics
- `DELETE /api/admin/storage` - Delete file
  - Body: `{ path: "collection/recordID/filename" }`
- `PUT /api/admin/storage/rename` - Rename file
  - Body: `{ oldPath, newPath }`

### Existing
- `POST /api/files` - Upload (already exists)
- `GET /api/files/{collection}/{id}/{filename}` - Serve (already exists)

## UI Components

### StorageBrowser.vue (Main View)
- Breadcrumb navigation
- File/folder list
- Upload button
- Search/filter bar
- Stats cards

### FileUploadModal.vue
- Collection selector
- RecordID input
- File picker
- Upload progress
- Success/error messages

### FileActionsMenu.vue
- Dropdown menu with actions
- Confirmation dialogs

## Styling
- Gruvbox theme (consistent with existing UI)
- Icons for file types
- Hover effects on rows
- Responsive grid/list view toggle
- Empty state when no files

## User Flow

1. **Navigate to Storage**
   - Click "Storage" in sidebar
   - See all collections as folders

2. **Browse Files**
   - Click collection folder → see recordIDs
   - Click recordID folder → see files
   - Use breadcrumb to go back

3. **Upload File**
   - Click "Upload File" button
   - Select collection and recordID
   - Choose files
   - See upload progress
   - Files appear in list

4. **Manage Files**
   - Click ⋮ menu on file
   - View, download, rename, or delete
   - Confirm destructive actions

5. **Search**
   - Type in search box
   - Results filter in real-time
   - Clear search to see all

## Implementation Priority

**Phase 1 (MVP)**
- Basic file browser with folder navigation
- File list with name, size, date
- Upload modal
- Download/delete actions
- Storage stats

**Phase 2 (Enhanced)**
- Search and filters
- Bulk actions
- Rename functionality
- Drag & drop upload
- Image previews/thumbnails

**Phase 3 (Advanced)**
- Grid view toggle
- Bulk download as ZIP
- File sharing/public URLs
- Storage quotas/limits
