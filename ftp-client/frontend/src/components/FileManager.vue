<template>
  <n-space vertical align="start" class="container">
    <!-- 标题 -->
    <h2 class="title">File Manager</h2>

    <!-- 操作按钮 -->
    <n-space>
      <n-button @click="refreshFiles" type="primary" size="large"
        >Refresh</n-button
      >
      <n-button @click="returnToDir" type="primary" size="large">
        Return
      </n-button>
      <n-button @click="createFolder" type="success" size="large"
        >Create Folder</n-button
      >
      <n-button @click="uploadFile" type="info" size="large"
        >Upload File</n-button
      >
    </n-space>

    <!-- 文件列表容器 -->
    <n-card class="file-list-card" bordered>
      <n-table bordered v-if="directories.length > 0">
        <thead>
          <tr>
            <th>文件名</th>
            <th>类型</th>
            <th>文件大小</th>
            <th>日期</th>
            <th>操作</th>
          </tr>
        </thead>

        <tbody>
          <template v-for="item in directories" :key="item.Name">
            <tr>
              <td>{{ item.Name }}</td>
              <td>{{ getFileType(item.Type) }}</td>
              <td>{{ formatSize(item.Size) }}</td>
              <td>{{ formatTime(item.Time) }}</td>
              <td>
                <n-space>
                  <n-button
                    v-if="item.Type === 1"
                    @click="openDir(item.Name)"
                    type="primary"
                    size="small"
                  >
                    enter
                  </n-button>
                  <n-button
                    @click="downloadFile(item.Name)"
                    type="primary"
                    size="small"
                  >
                    Download
                  </n-button>
                  <n-button
                    @click="deleteFile(item.Name)"
                    type="error"
                    size="small"
                  >
                    Delete
                  </n-button>
                </n-space>
              </td>
            </tr>
          </template>
        </tbody>
      </n-table>
      <n-empty v-else description="No files available" />
    </n-card>
  </n-space>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import {
  List,
  Download,
  CreateFolder,
  Delete,
  Upload,
} from "../../wailsjs/go/main/FTPClient";
import { OpenAndUploadFile } from "../../wailsjs/go/main/app";
import {
  NButton,
  NSpace,
  NList,
  NCard,
  NEmpty,
  NListItem,
  NTable,
} from "naive-ui";
export default defineComponent({
  components: { NButton, NSpace, NList, NCard, NEmpty, NListItem, NTable },
  setup() {
    const currentPath = ref("./");
    // const directories = ref<string[]>([]);
    const directories = ref<any[]>([]);
    const uploadedFiles = ref<string[]>([]); // 用于存储已上传的文件路径

    const getFileType = (type: number) => {
      switch (type) {
        case 0:
          return "File";
        case 1:
          return "Directory";
        default:
          return "Unknown";
      }
    };

    const formatSize = (size: number) => {
      if (size < 1024) return `${size} B`;
      else if (size < 1024 * 1024) return `${(size / 1024).toFixed(2)} KB`;
      else if (size < 1024 * 1024 * 1024)
        return `${(size / (1024 * 1024)).toFixed(2)} MB`;
      else return `${(size / (1024 * 1024 * 1024)).toFixed(2)} GB`;
    };

    const formatTime = (time: string) => {
      const date = new Date(time);
      return date.toLocaleString();
    };

    const refreshFiles = async () => {
      try {
        directories.value = await List(currentPath.value);
        console.log(directories.value);
      } catch (error: any) {
        alert("Failed to list files: " + error.message);
      }
    };

    const openDir = (dir: string) => {
      currentPath.value = `${currentPath.value}/${dir}`;
      refreshFiles();
    };

    const returnToDir = () => {
      let path = currentPath.value.split("/");
      path.pop();
      currentPath.value = path.join("/");
      refreshFiles();
    };

    const uploadFile = async () => {
      try {
        let filePath = await OpenAndUploadFile();
        console.log("filepath", filePath);
        let filename = filePath.split("/").pop();
        await Upload(filePath, `${currentPath.value}/${filename}`);
        refreshFiles();
      } catch (error: any) {
        alert("Failed to upload file: " + error.message);
      }
    };

    const downloadFile = async (file: string) => {
      try {
        await Download(
          `${currentPath.value}/${file}`,
          `/Users/aliancn/Downloads/${file}`
        );
      } catch (error: any) {
        alert("Failed to download file: " + error.message);
      }
    };

    const createFolder = async () => {
      const folderName = prompt("Enter folder name:");
      if (folderName) {
        try {
          await CreateFolder(`${currentPath.value}/${folderName}`);
          refreshFiles();
        } catch (error: any) {
          alert("Failed to create folder: " + error.message);
        }
      }
    };

    const deleteFile = async (file: string) => {
      try {
        await Delete(`${currentPath.value}/${file}`);
        refreshFiles();
      } catch (error: any) {
        alert("Failed to delete file: " + error.message);
      }
    };

    return {
      getFileType,
      formatSize,
      formatTime,
      directories,
      uploadedFiles,
      returnToDir,
      openDir,
      refreshFiles,
      uploadFile,
      downloadFile,
      createFolder,
      deleteFile,
    };
  },
});
</script>

<style scoped>
/* 主容器样式 */
.container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

/* 标题样式 */
.title {
  font-size: 2rem;
  color: #333;
  font-weight: bold;
  margin-bottom: 20px;
}

/* 按钮间隔 */
n-space {
  margin-bottom: 20px;
}

/* 文件列表卡片样式 */
.file-list-card {
  width: 100%;
  background-color: #f9f9f9;
  padding: 20px;
  border-radius: 8px;
}

/* 文件名样式 */
.file-name {
  font-size: 1.1rem;
  color: #333;
}

/* 响应式布局调整 */
@media (max-width: 768px) {
  .container {
    padding: 15px;
  }
  .title {
    font-size: 1.5rem;
  }
}
</style>
