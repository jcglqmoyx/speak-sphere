<template>
  <ContentBase @keydown="handleKeydown" tabindex="0">
    <div class="query-container">
      <div class="search-section">
        <el-input
          v-model="searchWord"
          placeholder="输入要查询的单词"
          size="large"
          @keyup.enter="handleSearch"
          class="search-input"
        >
          <template #append>
            <el-button @click="handleSearch" :icon="Search">搜索</el-button>
          </template>
        </el-input>
      </div>

      <div v-if="showResults" class="results-section">
        <el-card class="word-card">
          <template #header>
            <div class="card-header">
              <span class="word-title">{{ currentWord }}</span>
              <div class="header-actions">
                <el-tooltip content="搜索词典 (快捷键: S)" placement="top">
                  <el-button circle @click="showSearchDrawer = !showSearchDrawer">
                    <el-icon><Search /></el-icon>
                  </el-button>
                </el-tooltip>
                <el-tooltip content="添加笔记 (Markdown支持)" placement="top">
                  <el-button circle @click="showEditNoteDrawer = true">
                    <el-icon><EditPen /></el-icon>
                  </el-button>
                </el-tooltip>
              </div>
            </div>
          </template>

          <div class="word-content">
            <div v-if="wordMeaning" class="meaning">{{ wordMeaning }}</div>
            <div v-else class="no-meaning">该单词暂无释义</div>
          </div>
        </el-card>
      </div>

      <div v-else class="welcome-section">
        <el-card class="welcome-card">
          <h2>单词查询工具</h2>
          <p>在输入框中输入要查询的单词，按回车或点击搜索按钮进行查询</p>
          <p class="shortcuts">快捷键说明：</p>
          <ul>
            <li><strong>Enter</strong> - 搜索单词</li>
            <li><strong>S</strong> - 打开词典搜索</li>
            <li><strong>数字键 1-9</strong> - 在对应词典中搜索</li>
          </ul>
        </el-card>
      </div>
    </div>

    <!-- 词典搜索抽屉 -->
    <el-drawer v-model="showSearchDrawer" title="词典搜索" :with-header="true">
      <div class="dictionary-list">
        <p v-for="(dictionary, index) in dictionaries" :key="dictionary.id" class="dictionary-item">
          <a :href="dictionary.prefix + (currentWord || '') + (dictionary.suffix || '')" target="_blank">
            <strong>{{ index + 1 }}. {{ dictionary.title }}</strong>
          </a>
          <el-divider />
        </p>
      </div>
    </el-drawer>

    <!-- 笔记编辑抽屉 -->
    <el-drawer v-model="showEditNoteDrawer" title="笔记 (Markdown支持)" :with-header="true" size="50%">
      <div class="note-editor">
        <el-tabs v-model="activeTab" type="border-card">
          <el-tab-pane label="编辑" name="edit">
            <el-input
              type="textarea"
              placeholder="添加笔记 (支持Markdown语法)"
              v-model="currentWordNote"
              :rows="20"
              maxlength="1000"
              show-word-limit
            ></el-input>
          </el-tab-pane>
          <el-tab-pane label="预览" name="preview">
            <div class="markdown-preview" v-html="renderedMarkdown"></div>
          </el-tab-pane>
        </el-tabs>
        <el-divider/>
        <div class="editor-actions">
          <el-button type="primary" @click="saveNote">保存</el-button>
          <el-button @click="showEditNoteDrawer = false">取消</el-button>
        </div>
      </div>
    </el-drawer>
  </ContentBase>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { ElButton, ElCard, ElDivider, ElDrawer, ElIcon, ElInput, ElTabPane, ElTabs, ElTooltip } from 'element-plus';
import { EditPen, Search } from '@element-plus/icons-vue';
import ContentBase from '@/components/ContentBase.vue';
import MarkdownIt from 'markdown-it';
import { getDictionaryList } from '@/assets/js/module/dictionary/query';

const searchWord = ref('');
const currentWord = ref('');
const wordMeaning = ref('');
const showResults = ref(false);
const showSearchDrawer = ref(false);
const showEditNoteDrawer = ref(false);
const activeTab = ref('edit');
const currentWordNote = ref('');
const dictionaries = ref([]);

// 初始化markdown渲染器
const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true
});

// 计算属性：渲染markdown内容
const renderedMarkdown = computed(() => {
  if (!currentWordNote.value) {
    return '<p class="empty-note">暂无笔记内容</p>';
  }
  return md.render(currentWordNote.value);
});

onMounted(async () => {
  // 加载词典列表
  const response = await getDictionaryList();
  if (response && response.data) {
    dictionaries.value = response.data;
  }
});

const handleSearch = () => {
  if (searchWord.value.trim()) {
    currentWord.value = searchWord.value.trim();
    wordMeaning.value = ''; // 这里只做前端演示，实际应用中可以从本地存储或其他来源获取释义
    showResults.value = true;

    // 模拟获取释义（实际应用中这里可以调用API）
    setTimeout(() => {
      wordMeaning.value = `这是 ${currentWord.value} 的示例释义。在实际应用中，这里会显示单词的真实释义。`;
    }, 300);
  }
};

const handleKeydown = (event) => {
  if (showEditNoteDrawer.value) return;
  
  if (event.key === 'Enter') {
    handleSearch();
  } else if (event.key === 's' || event.key === 'S') {
    if (showResults.value) {
      showSearchDrawer.value = !showSearchDrawer.value;
    }
  } else if ('0123456789'.includes(event.key)) {
    let number = +event.key;
    if (number === 0) number = 10;
    if (number <= dictionaries.value.length && showResults.value) {
      const dictionary = dictionaries.value[number - 1];
      window.open(dictionary.prefix + (currentWord.value || '') + (dictionary.suffix || ''), '_blank');
    }
  }
};

const saveNote = () => {
  // 这里可以添加保存笔记的逻辑
  console.log('保存笔记:', currentWordNote.value);
  showEditNoteDrawer.value = false;
};
</script>

<style scoped>
.query-container {
  max-width: 600px;
  margin: 0 auto;
  padding: 20px;
}

.search-section {
  margin-bottom: 30px;
}

.search-input {
  width: 100%;
}

.results-section {
  margin-top: 20px;
}

.welcome-section {
  text-align: center;
  margin-top: 50px;
}

.welcome-card {
  max-width: 500px;
  margin: 0 auto;
}

.word-card {
  margin: 20px 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.word-title {
  font-size: 24px;
  font-weight: bold;
  color: var(--el-text-color-primary);
}

.header-actions {
  display: flex;
  gap: 10px;
}

.word-content {
  padding: 20px 0;
}

.meaning {
  font-size: 16px;
  line-height: 1.6;
  color: var(--el-text-color-regular);
}

.no-meaning {
  font-style: italic;
  color: var(--el-text-color-secondary);
}

.shortcuts {
  margin-top: 20px;
  font-weight: bold;
}

.dictionary-list {
  padding: 10px;
}

.dictionary-item {
  margin: 10px 0;
}

.note-editor {
  height: calc(100vh - 60px);
  display: flex;
  flex-direction: column;
}

:deep(.el-tabs__content) {
  flex: 1;
  overflow-y: auto;
}

:deep(.el-textarea__inner) {
  height: 100%;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 14px;
}

.markdown-preview {
  height: 100%;
  overflow-y: auto;
  padding: 10px;
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  background-color: var(--el-bg-color);
  color: var(--el-text-color-primary);
}

.markdown-preview :deep(h1) {
  color: var(--el-text-color-primary);
  border-bottom: 1px solid var(--el-border-color);
  padding-bottom: 5px;
  margin-top: 0;
}

.markdown-preview :deep(h2) {
  color: var(--el-text-color-primary);
  border-bottom: 1px solid var(--el-border-color);
  padding-bottom: 5px;
}

.markdown-preview :deep(h3) {
  color: var(--el-text-color-primary);
}

.markdown-preview :deep(p) {
  color: var(--el-text-color-primary);
  line-height: 1.6;
}

.markdown-preview :deep(code) {
  background-color: var(--el-fill-color-light);
  padding: 2px 4px;
  border-radius: 3px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  color: var(--el-text-color-primary);
}

.markdown-preview :deep(pre) {
  background-color: var(--el-fill-color-light);
  padding: 10px;
  border-radius: 4px;
  overflow-x: auto;
  border: 1px solid var(--el-border-color);
}

.markdown-preview :deep(pre code) {
  background-color: transparent;
  color: var(--el-text-color-primary);
}

.markdown-preview :deep(blockquote) {
  border-left: 4px solid var(--el-color-primary);
  padding-left: 10px;
  margin-left: 0;
  color: var(--el-text-color-primary);
  background-color: var(--el-fill-color-lighter);
}

.markdown-preview :deep(ul), .markdown-preview :deep(ol) {
  padding-left: 20px;
  color: var(--el-text-color-primary);
}

.markdown-preview :deep(li) {
  margin-bottom: 4px;
  color: var(--el-text-color-primary);
}

.markdown-preview :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 10px 0;
  border: 1px solid var(--el-border-color);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.markdown-preview :deep(th), .markdown-preview :deep(td) {
  border: 1px solid var(--el-border-color);
  padding: 8px;
  text-align: left;
}

.markdown-preview :deep(th) {
  background-color: var(--el-fill-color-light);
  color: var(--el-text-color-primary);
  font-weight: 600;
}

.markdown-preview :deep(td) {
  color: var(--el-text-color-primary);
  background-color: var(--el-bg-color);
}

.markdown-preview :deep(tr:nth-child(even)) {
  background-color: var(--el-fill-color-lighter);
}

.markdown-preview :deep(a) {
  color: var(--el-color-primary);
  text-decoration: none;
}

.markdown-preview :deep(a:hover) {
  text-decoration: underline;
}

.empty-note {
  color: var(--el-text-color-placeholder);
  font-style: italic;
  text-align: center;
  padding: 20px;
}

.editor-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 10px;
}
</style>
