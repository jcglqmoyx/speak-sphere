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
                <el-tooltip content="收藏单词" placement="top">
                  <el-button circle @click="handleAddToBook" :disabled="!searchWord.trim()" class="favorite-btn">
                    <el-icon><Star /></el-icon>
                  </el-button>
                </el-tooltip>
                <el-tooltip content="搜索词典 (快捷键: S)" placement="top">
                  <el-button circle @click="showSearchDrawer = !showSearchDrawer">
                    <el-icon><Search /></el-icon>
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

    <!-- 收藏单词对话框 -->
    <el-dialog v-model="showAddToBookDialog" title="收藏单词" width="500px">
      <div class="dialog-title">将单词 "<strong>{{ searchWord }}</strong>" 添加到词书</div>
      <el-form>
        <el-form-item label="选择词书">
          <el-radio-group v-model="selectedBookId" class="book-radio-group">
            <el-radio v-for="book in books" :key="book.id" :label="book.id" class="book-radio">
              <div class="book-radio-content">
                <span>{{ book.title }}</span>
                <el-tag v-if="bookWordExistsMap[book.id]" type="success" size="small" class="exists-tag">
                  已存在
                </el-tag>
              </div>
            </el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showAddToBookDialog = false">取消</el-button>
          <el-button 
            type="primary" 
            @click="confirmAddToBook" 
            :disabled="!selectedBookId"
            :class="{ 'warning-btn': selectedBookId && isWordInSelectedBook() }"
          >
            {{ selectedBookId && isWordInSelectedBook() ? '无需重复添加' : '确认添加' }}
          </el-button>
        </span>
      </template>
    </el-dialog>
  </ContentBase>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { ElButton, ElCard, ElDialog, ElDivider, ElDrawer, ElForm, ElFormItem, ElIcon, ElInput, ElMessage, ElRadio, ElRadioGroup, ElTag, ElTooltip } from 'element-plus';
import { Search, Star } from '@element-plus/icons-vue';
import ContentBase from '@/components/ContentBase.vue';
import { getDictionaryList } from '@/assets/js/module/dictionary/query';
import { getBookList } from '@/assets/js/module/book/query';
import { AddEntry } from '@/assets/js/module/entry/add';
import { checkWordInBook } from '@/assets/js/module/entry/query';

const searchWord = ref('');
const currentWord = ref('');
const wordMeaning = ref('');
const showResults = ref(false);
const showSearchDrawer = ref(false);
const showEditNoteDrawer = ref(false);
const showAddToBookDialog = ref(false);
const dictionaries = ref([]);
const books = ref([]);
const selectedBookId = ref(null);
const bookWordExistsMap = ref({}); // 存储每个词书是否包含当前单词

onMounted(async () => {
  // 加载词典列表
  const dictionaryResponse = await getDictionaryList();
  if (dictionaryResponse && dictionaryResponse.data) {
    dictionaries.value = dictionaryResponse.data;
  }

  // 加载词书列表
  const bookResponse = await getBookList(100000000, 1);
  if (bookResponse && bookResponse.data) {
    books.value = bookResponse.data;
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

const handleAddToBook = async () => {
  if (!searchWord.value.trim()) {
    return;
  }
  
  selectedBookId.value = null;
  bookWordExistsMap.value = {}; // 清空之前的记录
  showAddToBookDialog.value = true;
  
  // 检查每个词书是否已包含当前单词
  if (books.value.length > 0) {
    const word = searchWord.value.trim();
    
    for (const book of books.value) {
      try {
        const response = await checkWordInBook(word, book.id);
        const exists = response && response.code === 0 && response.data.exists;
        bookWordExistsMap.value[book.id] = exists;
        
        if (exists && !selectedBookId.value) {
          selectedBookId.value = book.id;
        }
      } catch (error) {
        console.error(`检查词书 ${book.title} 失败:`, error);
        bookWordExistsMap.value[book.id] = false;
      }
    }
  }
};

const isWordInSelectedBook = () => {
  if (!selectedBookId.value) {
    return false;
  }
  return !!bookWordExistsMap.value[selectedBookId.value];
};

const confirmAddToBook = async () => {
  if (!selectedBookId.value || !searchWord.value.trim()) {
    return;
  }

  try {
    const response = await AddEntry(
      selectedBookId.value,
      searchWord.value.trim(),
      wordMeaning.value || '暂无释义', // 使用当前显示的释义，如果没有则使用默认值
      ''
    );

    if (response && response.code === 0) {
      ElMessage.success('添加单词成功');
      showAddToBookDialog.value = false;
    } else {
      const errorMsg = response ? response.message : '添加失败，请重试';
      ElMessage.error(errorMsg);
    }
  } catch (error) {
    console.error('添加单词失败:', error);
    ElMessage.error('添加单词失败，请重试');
  }
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

:deep(.el-input-group__append) {
  display: flex;
  gap: 0;
}

:deep(.el-input-group__append .el-button) {
  border-radius: 0;
  margin: 0;
}

:deep(.el-input-group__append .el-button:first-child) {
  border-top-left-radius: 0;
  border-bottom-left-radius: 0;
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
  border-right: 1px solid var(--el-border-color);
}

:deep(.el-input-group__append .el-button:last-child) {
  border-top-left-radius: 0;
  border-bottom-left-radius: 0;
  border-top-right-radius: var(--el-border-radius-base);
  border-bottom-right-radius: var(--el-border-radius-base);
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

.favorite-btn {
  /* 移除自定义样式，使用与其他按钮一致的默认样式 */
}

.favorite-btn:disabled {
  background-color: #f5f5f5;
  border-color: #dcdfe6;
  color: #c0c4cc;
}

.book-radio-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.book-radio {
  display: flex;
  align-items: center;
  padding: 8px;
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  transition: all 0.3s;
}

.book-radio:hover {
  border-color: var(--el-color-primary);
  background-color: var(--el-fill-color-lighter);
}

:deep(.el-radio__input.is-checked + .el-radio__label) {
  color: var(--el-color-primary);
}

:deep(.el-radio__input.is-checked .el-radio__inner) {
  background-color: var(--el-color-primary);
  border-color: var(--el-color-primary);
}

.dialog-title {
  margin-bottom: 20px;
  font-size: 16px;
  color: var(--el-text-color-primary);
}

.book-radio-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.exists-tag {
  margin-left: 10px;
}

.warning-btn {
  background-color: var(--el-color-warning) !important;
  border-color: var(--el-color-warning) !important;
}

.warning-btn:hover {
  background-color: var(--el-color-warning-light-3) !important;
  border-color: var(--el-color-warning-light-3) !important;
}
</style>
