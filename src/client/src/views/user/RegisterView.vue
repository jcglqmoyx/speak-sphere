<template>
  <ContentBase>
    <div class="register-container">
      <el-form ref="registerForm" label-width="80px">
        <el-form-item label="服务器" prop="serverLink" class="form-item">
          <el-input v-model="serverLink" placeholder="请输入后台服务器地址"></el-input>
        </el-form-item>
        <el-form-item label="用户名" prop="username" class="form-item">
          <el-input v-model="username" placeholder="请输入用户名"></el-input>
        </el-form-item>
        <el-form-item label="邮箱" prop="email" class="form-item">
          <el-input v-model="email" placeholder="请输入邮箱"></el-input>
        </el-form-item>
        <el-form-item label="头像" prop="avatar" class="form-item">
          <el-input v-model="avatar" placeholder="请输入邮箱"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password" class="form-item">
          <el-input type="password" v-model="password" placeholder="请输入密码"></el-input>
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword" class="form-item">
          <el-input type="password" v-model="confirmPassword" @keyup.enter="handleRegister"
                    placeholder="请确认密码"></el-input>
        </el-form-item>

        <div v-if="registerError" style="color: red; margin-bottom: 10px; margin-left: 100px;">
          {{ errorMessage }}
        </div>

        <el-form-item class="form-item">
          <el-button type="primary" @click="handleRegister" class="register-button">注册</el-button>
        </el-form-item>

        <p v-if="!registerSuccess" @click="login" style="margin-left: 80px; color: #04eefa">已有账号？点击登录</p>
        <p v-if="registerSuccess" style="margin-left: 80px; color: #0477fa">{{ registerSuccessMessage }}</p>
      </el-form>
    </div>
  </ContentBase>
</template>

<script setup>
import "element-plus/dist/index.css"
import {ref} from 'vue';
import {ElButton, ElForm, ElFormItem, ElInput} from "element-plus";
import {useRouter} from "vue-router";
import {useStore} from "vuex";
import ContentBase from "@/components/ContentBase.vue";

const store = useStore();

const serverLink = ref('http://127.0.0.1:8080');
const username = ref('');
const email = ref('');
const avatar = ref('');
const password = ref('');
const confirmPassword = ref('');
const registerError = ref(false);
const errorMessage = ref('');
const registerSuccess = ref(false);
const registerSuccessMessage = ref('');

const handleRegister = async () => {
  if (password.value !== confirmPassword.value) {
    registerError.value = true;
    errorMessage.value = '两次输入的密码不一致';
    return;
  }
  const response = await store.dispatch('register', {
    'serverLink': serverLink.value,
    'username': username.value,
    'email': email.value,
    'avatar': avatar.value,
    'password': password.value,
    'confirm_password': confirmPassword.value,
  });
  if (response === true) {
    registerError.value = false;
    registerSuccess.value = true;
    registerSuccessMessage.value = '注册成功, 即将自动跳转到登录页面';
    setTimeout(() => {
      login();
    }, 2000);
  } else {
    registerError.value = true;
    errorMessage.value = response;
  }
};

const router = useRouter();
const login = () => {
  router.push({name: 'login'});
};
</script>

<style scoped>
.register-container {
  max-width: 400px;
  margin: 0 auto;
  padding: 20px;
}

.form-item {
  width: 100%;
}

.register-button {
  width: 100%;
}
</style>
