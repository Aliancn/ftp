<template>
  <div class="login-page">
    <h2>FTP Client Login</h2>
    <form @submit.prevent="login">
      <div>
        <label>Server Address:</label>
        <input v-model="server" placeholder="e.g., 127.0.0.1" />
      </div>
      <div>
        <label>Username:</label>
        <input v-model="username" />
      </div>
      <div>
        <label>Password:</label>
        <input type="password" v-model="password" />
      </div>
      <button type="submit">Login</button>
    </form>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import { Connect } from "../../wailsjs/go/main/FTPClient";

export default defineComponent({
  emits: ["login-success"],
  setup(_, { emit }) {
    const server = ref("");
    const username = ref("");
    const password = ref("");

    const login = async () => {
      try {
        await Connect(server.value, username.value, password.value);
        emit("login-success");
      } catch (error: any) {
        alert("Login failed: " + error.message);
      }
    };

    return { server, username, password, login };
  },
});
</script>

<style scoped>
.login-page {
  max-width: 400px;
  margin: 25vh auto;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 4px;
  /* 修改页面居中 */
  position: relative;
}

.login-page h2 {
  text-align: center;
}

.login-page form div {
  margin-bottom: 15px;
}

.login-page label {
  display: block;

  margin-bottom: 5px;
}

.login-page input {
  width: 100%;
  padding: 8px;
  box-sizing: border-box;
}

.login-page button {
  width: 100%;
  padding: 10px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.login-page button:hover {
  background-color: #0056b3;
}
</style>
