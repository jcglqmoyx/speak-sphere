<template>
  <NavBarView/>
  <router-view/>
</template>

<script setup>
import "element-plus/dist/index.full.min.js";
import "element-plus/theme-chalk/dark/css-vars.css"
import NavBarView from "@/components/NavBar.vue";
import { getThemeInstance } from "@/store/theme.js";

// 初始化主题
getThemeInstance();
</script>

<style>
</style>
