**简体中文** | [English](README.md)

# Speak Sphere 🌍 语言学习平台

**Speak Sphere** 是一款集成了"不背单词"和"轻听英语"功能的现代语言学习应用，使用Rust和egui框架构建，提供流畅的桌面学习体验。

## 🌟 功能特性

### 📚 背单词模块（集成不背单词）
- **智能记忆系统**：基于艾宾浩斯遗忘曲线
- **多模式学习**：单词卡片、拼写练习、听写测试
- **丰富词库**：涵盖各阶段英语学习词汇
- **发音支持**：标准英语发音练习
- **进度跟踪**：学习进度可视化

### 🎧 阅读有声书（集成轻听英语）
- **沉浸式阅读**：边听边读提升语感
- **同步字幕**：实时显示文本内容
- **语音控制**：支持播放/暂停/调速
- **生词标记**：一键标记生词并加入学习计划
- **多语言支持**：英语、日语、法语等多种语言

## 🚀 快速开始

### 系统要求
- **操作系统**: Windows 10/11, macOS 10.15+, 或 Linux
- **内存**: 至少 4GB RAM
- **存储**: 至少 100MB 可用空间

### 安装方法

#### 从源码构建
```bash
# 克隆项目
git clone https://github.com/jcglqmoyx/speak-sphere.git
cd speak-sphere

# 构建项目
cargo build --release

# 运行应用
cargo run
```

#### 下载预编译版本
访问 [Release页面](https://github.com/jcglqmoyx/speak-sphere/releases) 下载对应平台的预编译版本。

### 使用指南

1. **启动应用**：双击应用图标或运行 `target/release/speak-sphere`
2. **选择功能**：主界面选择"背单词"或"阅读有声书"
3. **开始学习**：根据指引进入相应学习模块
4. **个性化设置**：根据需要调整学习参数

## 🛠️ 开发指南

### 技术栈
- **后端**: Rust
- **前端**: egui (即时模式GUI)
- **音频**: 待集成音频播放库
- **数据存储**: SQLite (计划中)

### 开发环境设置
```bash
# 安装Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# 克隆项目
git clone https://github.com/jcglqmoyx/speak-sphere.git
cd speak-sphere

# 运行开发版本
cargo run

# 运行测试
cargo test
```

### 项目结构
```
speak-sphere/
├── src/
│   ├── main.rs          # 主应用入口
│   ├── config.rs        # 配置管理
│   └── utils/           # 工具模块
│       ├── mod.rs
│       ├── font.rs      # 字体管理
│       └── display_size.rs
├── Cargo.toml          # 依赖配置
└── README.md          # 项目文档
```

## 📦 依赖库

主要依赖：
- `egui` - 即时模式GUI框架
- `eframe` - egui应用框架
- `rdev` - 系统交互功能
- `once_cell` - 单次初始化容器

## 🤝 贡献指南

我们欢迎各种形式的贡献！

### 如何贡献
1. Fork 本仓库
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 开发规范
- 遵循Rust官方编码规范
- 提交信息使用英文描述
- 新功能需要包含测试用例
- 确保所有测试通过

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 📞 联系方式

- **项目主页**: [GitHub Repository](https://github.com/jcglqmoyx/speak-sphere)
- **问题反馈**: [Issues](https://github.com/jcglqmoyx/speak-sphere/issues)
- **邮箱**: jcglqmoyx@gmail.com

## 🙏 致谢

感谢以下开源项目的支持：
- [egui](https://github.com/emilk/egui) - 优秀的即时模式GUI框架
- [不背单词](https://www.bbdc.net/) - 单词学习理念参考
- [轻听英语](https://www.qt-ing.com/) - 有声学习理念参考

---

**Speak Sphere** - 让语言学习更加高效有趣！ 🚀
