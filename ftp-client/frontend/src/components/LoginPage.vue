<template>
  <div class="login-page">
    <div class="background">
      <div class="shape"></div>
      <div class="shape"></div>
    </div>
    <form @submit.prevent="login">
      <h3>Login Here</h3>

      <label for="server">Server Address</label>
      <input v-model="server" placeholder="e.g., 127.0.0.1" id="server" />

      <label for="username">Username</label>
      <input v-model="username" placeholder="Username" id="username" />

      <label for="password">Password</label>
      <input
        type="password"
        v-model="password"
        placeholder="Password"
        id="password"
      />

      <button type="submit" :disabled="isLoading">
        <span v-if="isLoading">Logging in...</span>
        <span v-else>Login</span>
      </button>
    </form>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from "vue";
import { Connect } from "../../wailsjs/go/main/app";

export default defineComponent({
  emits: ["login-success"],
  setup(_, { emit }) {
    const server = ref("127.0.0.1:2121");
    const username = ref("rw");
    const password = ref("123");
    const isLoading = ref(false);

    const login = async () => {
      isLoading.value = true;
      try {
        console.log("login", server.value, username.value, password.value);
        await Connect(server.value, username.value, password.value);
        emit("login-success");
      } catch (error: any) {
        alert("Login failed: " + error.message);
      } finally {
        isLoading.value = false;
      }
    };

    return { server, username, password, login ,isLoading};
  },
});
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Poppins:wght@300;500;600&display=swap');

*,
*:before,
*:after {
  padding: 0;
  margin: 0;
  box-sizing: border-box;
}

body {
  background-color: #080710;
}

.background {
  width: 430px;
  height: 520px;
  position: absolute;
  transform: translate(-50%, -50%);
  left: 50%;
  top: 50%;
}

.background .shape {
  height: 200px;
  width: 200px;
  position: absolute;
  border-radius: 50%;
}

.shape:first-child {
  background: linear-gradient(#1845ad, #23a2f6);
  left: -80px;
  top: -80px;
}

.shape:last-child {
  background: linear-gradient(to right, #ff512f, #f09819);
  right: -30px;
  bottom: -80px;
}

form {
  height: 520px;
  width: 400px;
  background-color: rgba(255, 255, 255, 0.13);
  position: absolute;
  transform: translate(-50%, -50%);
  top: 50%;
  left: 50%;
  border-radius: 10px;
  backdrop-filter: blur(20px);
  border: 2px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 0 40px rgba(8, 7, 16, 0.6);
  padding: 50px 35px;
}

form * {
  font-family: 'Poppins', sans-serif;
  color: #ffffff;
  letter-spacing: 0.5px;
  outline: none;
  border: none;
}

form h3 {
  font-size: 32px;
  font-weight: 500;
  line-height: 42px;
  text-align: center;
}

label {
  display: block;
  margin-top: 30px;
  font-size: 16px;
  font-weight: 500;
}

input {
  display: block;
  height: 50px;
  width: 100%;
  background-color: rgba(255, 255, 255, 0.07);
  border-radius: 3px;
  padding: 0 10px;
  margin-top: 8px;
  font-size: 14px;
  font-weight: 300;
}

input::placeholder {
  color: #e5e5e5;
}

button {
  margin-top: 50px;
  width: 100%;
  background-color: #007bff;
  color: white;
  padding: 15px 0;
  font-size: 18px;
  font-weight: 600;
  border-radius: 5px;
  cursor: pointer;
  border: none;
  transition: background-color 0.3s;
}

button:hover {
  background-color: rgb(74, 89, 105);
}

button:disabled {
  background-color: #c0c0c0;
  cursor: not-allowed;
}
</style>
