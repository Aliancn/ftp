<template>
  <div>
    <LoginPage v-if="!isLoggedIn" @login-success="handleLoginSuccess" />
    <FileManager v-else />
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import LoginPage from "./components/LoginPage.vue";
import FileManager from "./components/FileManager.vue";

export default defineComponent({
  components: { LoginPage, FileManager },
  setup() {
    const isLoggedIn = ref(false);

    const handleLoginSuccess = () => {
      isLoggedIn.value = true;
    };

    return { isLoggedIn, handleLoginSuccess };
  },
});
</script>
