<template>
  <el-menu-item v-if="store.state.user.isAuthenticated" index="4" @click="logout">退出登录</el-menu-item>
  <el-menu-item v-else index="4" @click="login">登录</el-menu-item>
</template>

<script setup>
import {ElMenuItem} from "element-plus";
import {useStore} from "vuex";
import {useRouter} from "vue-router";

const store = useStore();
const router = useRouter();

const login = async () => {
  await router.push({name: 'login'});
}

const logout = async () => {
  await store.dispatch('logout');
  await router.push({name: 'login'});
}
</script>

<style scoped>
</style>