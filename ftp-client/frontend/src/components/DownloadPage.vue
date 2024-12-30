<template>
  <n-card title="文件下载列表">
    <!-- 文件表格 -->
    <n-table :bordered="true">
      <thead>
        <tr>
          <th>文件名</th>
          <th>状态</th>
          <th>进度</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <template v-for="row in downloads" :key="row.fileName">
          <tr>
            <td>{{ row.fileName }}</td>
            <td>
              <n-tag :type="statusType(row.status)">
                {{ statusText(row.status) }}
              </n-tag>
            </td>
            <td>
              <n-progress
                type="line"
                :percentage="computeProgress(row)"
                indicator-placement="inside"
                processing
              />
            </td>
            <td>
              <n-button
                size="small"
                type="primary"
                v-if="row.status === 'downloading'"
                @click="Stop(row)"
                >暂停</n-button
              >
              <n-button
                size="small"
                type="primary"
                v-else-if="row.status === 'paused'"
                @click="Continue(row)"
                >继续</n-button
              >
              <n-button
                size="small"
                type="error"
                @click="deleteDownload(row.fileName)"
                >删除</n-button
              >
            </td>
          </tr>
        </template>
      </tbody>
    </n-table>
  </n-card>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted, PropType } from "vue";
import { NCard, NTable, NProgress, NButton, NTag, useMessage } from "naive-ui";
import { StopDownload, Download } from "../../wailsjs/go/main/app";
export default defineComponent({
  name: "FileList",
  components: { NCard, NTable, NProgress, NButton, NTag },
  props: {
    downloads: {
      type: Array as PropType<
        {
          fileName: string;
          fileSize: number;
          downloaded: number;
          status: string;
        }[]
      >,
      required: true,
    },
  },
  setup(props) {
    const message = useMessage();

    // 状态标签的样式
    const statusType = (status: string) => {
      switch (status) {
        case "downloading":
          return "info";
        case "paused":
          return "warning";
        case "completed":
          return "success";
        default:
          return "default";
      }
    };

    const statusText = (status: string) => {
      switch (status) {
        case "downloading":
          return "下载中";
        case "paused":
          return "已暂停";
        case "completed":
          return "已完成";
        default:
          return "未知";
      }
    };

    // 计算进度百分比
    const computeProgress = (row: any) => {
      if (row.fileSize === 0) return 0;
      return (row.downloaded / row.fileSize) * 100;
    };

    const Stop = async (row: any) => {
      const file = props.downloads.find((d) => d.fileName === row.fileName);
      if (file) {
        file.status = "paused";
        await StopDownload();
        message.warning(`暂停下载: ${row.fileName}`);
      }
    };

    const Continue = async (row: any) => {
      const file = props.downloads.find((d) => d.fileName === row.fileName);
      if (file) {
        file.status = "downloading";
        await Download(row.localPath, row.remotePath, row.Size);
        message.success(`继续下载: ${row.fileName}`);
      }
    };

    // 删除下载任务
    const deleteDownload = async (fileName: string) => {
      const index = props.downloads.findIndex((d) => d.fileName === fileName);
      if (index !== -1) {
        props.downloads.splice(index, 1);
      }
      message.error(`删除下载: ${fileName}`);
    };

    return {
      statusType,
      statusText,
      computeProgress,
      Stop,
      Continue,
      deleteDownload,
    };
  },
});
</script>

<style scoped>
.n-card {
  margin-top: 16px;
  width: 90vh;
}
.n-table {
  margin-top: 16px;
}
</style>
