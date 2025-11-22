const { defineConfig } = require("@vue/cli-service");
module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    port: 3000,
  },
  // 对于 GitHub Pages，设置相对路径以确保资源正确加载
  publicPath: process.env.NODE_ENV === "production" ? "./" : "/",
  outputDir: "dist",
  assetsDir: "static",
  filenameHashing: true,
  // 确保生产环境构建时正确处理路由
  chainWebpack: (config) => {
    // 在开发环境不需要特殊处理
    if (process.env.NODE_ENV === "development") {
      return;
    }

    // 为生产环境配置正确的路径
    config.plugin("html").tap((args) => {
      args[0].minify = {
        removeComments: true,
        collapseWhitespace: true,
        removeAttributeQuotes: true,
        collapseBooleanAttributes: true,
        removeScriptTypeAttributes: true,
      };
      return args;
    });
  },
});
