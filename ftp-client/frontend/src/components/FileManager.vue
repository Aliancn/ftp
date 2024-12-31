<template>
  <n-space justify="space-around" align="center">
    <n-button @click="showFilePage = true" type="primary" size="large"
      >FilePage</n-button
    >
    <n-button @click="showFilePage = false" type="primary" size="large"
      >DownloadPage</n-button
    >
  </n-space>

  <n-space vertical align="start" class="container" v-if="showFilePage">
    <!-- 标题 -->
    <h2 class="title">File Manager</h2>
    <h3 class="title">Current Path: {{ currentPath }}</h3>
    <n-space>
      <n-button @click="refreshFiles" type="primary" size="large"
        >Refresh</n-button
      >
      <n-button @click="returnToDir" type="primary" size="large">
        Return
      </n-button>
      <n-button @click="showCreateFolder" type="success" size="large"
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
              <td>{{ item.Type }}</td>
              <td>{{ formatSize(item.Size) }}</td>
              <td>{{ formatTime(item.Time) }}</td>
              <td>
                <n-space>
                  <n-button
                    v-if="item.Type === 'Directory'"
                    @click="openDir(item.Name)"
                    type="primary"
                    size="small"
                  >
                    enter
                  </n-button>
                  <n-button
                    v-else
                    @click="downloadFile(item)"
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
    <n-modal v-model:show="showCreateFolderModal" title="新建文件夹">
      <n-card
        style="width: 600px"
        title="新建文件夹"
        :bordered="false"
        size="huge"
        role="dialog"
        aria-modal="true"
        ><n-input
          v-model:value="newFolderName"
          placeholder="请输入文件夹名称"
          style="margin: 10px; padding: 10px"
        />
        <n-space justify="space-between" style="width: 200px; margin: auto"
          ><n-button @click="showCreateFolderModal = false">取消</n-button>
          <n-button @click="createFolder" type="primary">确定</n-button>
        </n-space>
      </n-card>
    </n-modal>
  </n-space>
  <n-message-provider v-else>
    <DownloadPage :downloads="exampleDownloads" />
  </n-message-provider>
</template>
<script lang="ts">
import { defineComponent, ref, onMounted } from "vue";
import {
  OpenAndUploadFile,
  List,
  Download,
  CreateFolder,
  Delete,
  Upload,
} from "../../wailsjs/go/main/app";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import {
  NButton,
  NSpace,
  NList,
  NCard,
  NEmpty,
  NListItem,
  NTable,
  NModal,
  NInput,
  NMessageProvider,
} from "naive-ui";
import DownloadPage from "./DownloadPage.vue";
export default defineComponent({
  components: {
    NButton,
    NSpace,
    NList,
    NCard,
    NEmpty,
    NListItem,
    NTable,
    NModal,
    NInput,
    DownloadPage,
    NMessageProvider,
  },
  setup() {
    const currentPath = ref(".");
    // const directories = ref<string[]>([]);
    const directories = ref<any[]>([]);
    const uploadedFiles = ref<string[]>([]); // 用于存储已上传的文件路径
    const showCreateFolderModal = ref(false);
    const newFolderName = ref("");
    const showFilePage = ref(true);
    const exampleDownloads = ref([
      {
        fileName: "example1.zip",
        fileSize: 1000000, // 1 MB
        downloaded: 500000, // 500 KB
        status: "downloading", // 下载中
        remotePath: "/Users/aliancn/Downloads/example1.zip",
      },
      {
        fileName: "example2.pdf",
        fileSize: 2000000, // 2 MB
        downloaded: 2000000, // 已完成
        status: "completed", // 已完成
        remotePath: "/Users/aliancn/Downloads/example2.pdf",
      },
      {
        fileName: "example3.mp4",
        fileSize: 500000000, // 500 MB
        downloaded: 100000000, // 100 MB
        status: "paused", // 已暂停
        remotePath: "/Users/aliancn/Downloads/example3.mp4",
      },
    ]);

    const formatSize = (size: number) => {
      if (size < 1024) return `${size} B`;
      else if (size < 1024 * 1024) return `${(size / 1024).toFixed(2)} KB`;
      else if (size < 1024 * 1024 * 1024)
        return `${(size / (1024 * 1024)).toFixed(2)} MB`;
      else return `${(size / (1024 * 1024 * 1024)).toFixed(2)} GB`;
    };

    const formatTime = (dateTime: string) => {
      const months = {
        Jan: "01",
        Feb: "02",
        Mar: "03",
        Apr: "04",
        May: "05",
        Jun: "06",
        Jul: "07",
        Aug: "08",
        Sep: "09",
        Oct: "10",
        Nov: "11",
        Dec: "12",
      };

      const parts = dateTime.split(" ");
      const year = parts[0];
      const month = months[parts[1] as keyof typeof months];
      const day = parts[2].padStart(2, "0");
      const time = parts[3];

      return `${year}/${month}/${day} ${time}`;
    };

    const parseListInfo = (listInfo: string[]) => {
      let result: any[] = [];
      if (listInfo == null) return result;
      if (listInfo.length === 0) return result;
      for (const line of listInfo) {
        // 使用正则表达式解析 LIST 信息
        const parts = line.match(
          /^(?<permissions>[drwx\-]+)\s+\d+\s+\w+\s+\w+\s+(?<size>\d+)\s+(?<month>\w+)\s+(?<day>\d+)\s+(?<time>\d+:\d+|\d{4})\s+(?<name>.+)$/
        );

        if (parts?.groups) {
          const { permissions, size, month, day, time, name } = parts.groups;

          // 判断文件类型
          const type: "File" | "Directory" =
            permissions[0] === "d" ? "Directory" : "File";

          // 拼接时间字符串
          const currentYear = new Date().getFullYear();
          const formattedTime = time.includes(":") // 如果有时间，说明是当年
            ? `${currentYear} ${month} ${day} ${time}`
            : `${month} ${day} ${time}`; // 如果没有冒号，则是年份

          // 添加到结果数组
          result.push({
            Type: type,
            Size: parseInt(size, 10), // 转换为数字
            Time: formattedTime,
            Name: name,
          });
        }
      }

      return result;
    };

    const refreshFiles = async () => {
      try {
        const listResult = await List(currentPath.value);
        directories.value = parseListInfo(listResult);
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
      if (path.length < 2) return;
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

    const downloadFile = async (file: any) => {
      try {
        const index = exampleDownloads.value.findIndex(
          (d) => d.fileName === currentPath.value + "/" + file.Name
        );
        if (index !== -1) {
          alert("File is already downloaded");
          return;
        }
        let localPath = currentPath.value + "/" + file.Name;
        let remotePath = `/Users/aliancn/Downloads/${file.Name}`;
        exampleDownloads.value.push({
          fileName: localPath,
          fileSize: file.Size,
          downloaded: 0,
          status: "downloading",
          remotePath: remotePath,
        });

        await Download(
          localPath,
          remotePath, // TODO change to download path
          file.Size
        );
      } catch (error: any) {
        alert("Failed to download file: " + error.message);
      }
    };

    const showCreateFolder = () => {
      showCreateFolderModal.value = !showCreateFolderModal.value;
      console.log("showCreateFolderModal", showCreateFolderModal.value);
    };

    const createFolder = async () => {
      const folderName = newFolderName.value;
      if (folderName) {
        try {
          await CreateFolder(`${currentPath.value}/${folderName}`);
          showCreateFolderModal.value = !showCreateFolderModal.value;
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

    onMounted(() => {
      // 监听后端推送的下载进度
      EventsOn("download-progress", (progress: any) => {
        console.log("download-progress", progress);
        // 查找目标文件
        const index = exampleDownloads.value.findIndex(
          (d) => d.fileName === progress.fileName
        );
        if (index !== -1) {
          console.log("index", index);
          exampleDownloads.value[index].downloaded = progress.downloaded;
          if (progress.downloaded === exampleDownloads.value[index].fileSize) {
            exampleDownloads.value[index].status = "completed";
          } 
        }
        console.log("exampleDownloads", exampleDownloads.value);
      });
    });

    return {
      showCreateFolderModal,
      showCreateFolder,
      showFilePage,
      currentPath,
      exampleDownloads,
      newFolderName,
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
