<template>
  <ContentBase>
    <div class="login-container">
      <el-form ref="loginForm" label-width="80px">
        <el-form-item label="Server Link" prop="serverLink" class="form-item">
          <el-input v-model="serverLink" placeholder="请输入服务器地址"></el-input>
        </el-form-item>
        <el-form-item label="Username" prop="username" class="form-item">
          <el-input v-model="username" placeholder="请输入用户名"></el-input>
        </el-form-item>
        <el-form-item label="Password" prop="password" class="form-item">
          <el-input type="password" v-model="password" @keyup.enter="handleLogin" placeholder="请输入密码"></el-input>
        </el-form-item>
        <div v-if="loginError" style="color: red; margin-bottom: 10px; margin-left: 100px;">
          {{ errorMessage }}
        </div>
        <el-form-item class="form-item">
          <el-button type="primary" @click="handleLogin" class="login-button">Login</el-button>
        </el-form-item>

        <p @click="register" style="margin-left: 80px; color: #0477fa">没有账号？点击注册</p>

      </el-form>
    </div>
  </ContentBase>
</template>

<script setup>
import {ref} from 'vue';
import "element-plus/dist/index.css"
import {ElButton, ElForm, ElFormItem, ElInput} from "element-plus";
import {useStore} from "vuex";
import {useRouter} from "vue-router";
import ContentBase from "@/components/ContentBase.vue";

const serverLink = ref('http://127.0.0.1:8080')
if (localStorage.getItem('server_link') !== null) {
  serverLink.value = localStorage.getItem('server_link');
}

const username = ref('')
const password = ref('')
const loginError = ref(false)
const errorMessage = ref('')

const store = useStore();
const router = useRouter();

const handleLogin = async () => {
  const response = await store.dispatch('login', {
    'serverLink': serverLink.value,
    'username': username.value,
    'password': password.value,
  });
  if (response === true) {
    loginError.value = false;
    await router.push({name: 'home'});
  } else {
    loginError.value = true;
    errorMessage.value = response;
  }
};

const register = async () => {
  await router.push({name: 'register'});
}
</script>

<style scoped>
.login-container {
  max-width: 400px;
  margin: 0 auto;
  padding: 20px;
}

.form-item {
  width: 100%;
}

.login-button {
  width: 100%;
}
</style>
