<template>
  <div>
    <h2>File Manager</h2>
    <button @click="refreshFiles">Refresh</button>
    <button @click="createFolder">Create Folder</button>
    <button @click="uploadFile">Upload File</button>
    <ul>
      <li v-for="file in directories" :key="file">
        {{ file }}
        <button @click="downloadFile(file)">Download</button>
        <button @click="deleteFile(file)">Delete</button>
      </li>
    </ul>
  </div>
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
export default defineComponent({
  setup() {
    const currentPath = ref("./");
    const directories = ref<string[]>([]);
    const uploadedFiles = ref<string[]>([]); // 用于存储已上传的文件路径

    const refreshFiles = async () => {
      try {
        directories.value = await List(currentPath.value);
        console.log(directories.value);
      } catch (error: any) {
        alert("Failed to list files: " + error.message);
      }
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
      directories,
      uploadedFiles,
      refreshFiles,
      uploadFile,
      downloadFile,
      createFolder,
      deleteFile,
    };
  },
});
</script>
