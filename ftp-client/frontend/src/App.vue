<template>
  <div>
    <LoginPage v-if="!isLoggedIn" @login-success="handleLoginSuccess" />
    <FileManager v-else ref="fileManagerRef" />
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, watch, nextTick } from "vue";
import LoginPage from "./components/LoginPage.vue";
import FileManager from "./components/FileManager.vue";

export default defineComponent({
  components: { LoginPage, FileManager },
  setup() {
    const fileManagerRef = ref<InstanceType<typeof FileManager> | null>(null);
    const isLoggedIn = ref(false);

    const handleLoginSuccess = () => {
      isLoggedIn.value = true;
    };
    watch(isLoggedIn, async (newValue) => {
      if (newValue) {
        await nextTick(); // 等待 DOM 更新完成
        if (fileManagerRef.value) {
          console.log("refresh files");
          fileManagerRef.value.refreshFiles();
        }
      }
    });

    return { isLoggedIn, handleLoginSuccess, fileManagerRef };
  },
});
</script>
