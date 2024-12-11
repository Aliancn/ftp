<!-- <template>
    <div>
      <h1>FTP Client</h1>
      <div>
        <input v-model="server" placeholder="Server Address" />
        <input v-model="port" placeholder="Port" />
        <input v-model="username" placeholder="Username" />
        <input v-model="password" placeholder="Password" type="password" />
        <button @click="connect">Connect</button>
      </div>
  
      <div v-if="isConnected">
        <h2>Files</h2>
        <ul>
          <li v-for="file in files" :key="file">{{ file }}</li>
        </ul>
        <input type="file" @change="uploadFile" />
      </div>
    </div>
  </template>
  
  <script>
  import { Connect, List, Upload, Disconnect } from "../wailsjs/go/backend/FTPClient";
  
  export default {
    data() {
      return {
        server: "",
        port: "21",
        username: "",
        password: "",
        isConnected: false,
        files: [],
      };
    },
    methods: {
      async connect() {
        try {
          await Connect(`${this.server}:${this.port}`, this.username, this.password);
          this.isConnected = true;
          this.files = await List("/");
        } catch (error) {
          alert("Failed to connect: " + error.message);
        }
      },
      async uploadFile(event) {
        const file = event.target.files[0];
        if (!file) return;
  
        try {
          await Upload(file.path, `/${file.name}`);
          alert("File uploaded successfully!");
          this.files = await List("/");
        } catch (error) {
          alert("Failed to upload file: " + error.message);
        }
      },
      async disconnect() {
        await Disconnect();
        this.isConnected = false;
      },
    },
  };
  </script>
   -->