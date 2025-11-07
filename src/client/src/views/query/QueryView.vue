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
                    <el-icon>
                      <Star/>
                    </el-icon>
                  </el-button>
                </el-tooltip>
                <el-tooltip content="搜索词典 (快捷键: S)" placement="top">
                  <el-button circle @click="showSearchDrawer = !showSearchDrawer">
                    <el-icon>
                      <Search/>
                    </el-icon>
                  </el-button>
                </el-tooltip>
              </div>
            </div>
          </template>

          <div class="word-content" @click="handleWordContentClick">
            <div v-if="isLoading" class="loading-container">
              <el-skeleton :rows="5" animated/>
            </div>
            <div v-else-if="wordMeaning" class="meaning" v-html="wordMeaning"></div>
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
          <el-divider/>
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
import {onMounted, ref} from 'vue';
import {
  ElButton,
  ElCard,
  ElDialog,
  ElDivider,
  ElDrawer,
  ElForm,
  ElFormItem,
  ElIcon,
  ElInput,
  ElMessage,
  ElRadio,
  ElRadioGroup,
  ElTag,
  ElTooltip
} from 'element-plus';
import {Search, Star} from '@element-plus/icons-vue';
import ContentBase from '@/components/ContentBase.vue';
import {getDictionaryList} from '@/assets/js/module/dictionary/query';
import {getBookList} from '@/assets/js/module/book/query';
import {AddEntry} from '@/assets/js/module/entry/add';
import {checkWordInBook} from '@/assets/js/module/entry/query';
import {formatDefinition, getWordDefinition} from '@/assets/js/util/dictionary_api';

const searchWord = ref('');
const currentWord = ref('');
const wordMeaning = ref('');
const isLoading = ref(false);
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

// 处理单词内容区域的点击事件
const handleWordContentClick = (event) => {
  // 检查点击的是否是发音按钮
  if (event.target.classList.contains('audio-btn')) {
    const audioUrl = event.target.getAttribute('data-audio');
    if (audioUrl) {
      playAudio(audioUrl);
    }
  }
};

const handleSearch = async () => {
  if (searchWord.value.trim()) {
    currentWord.value = searchWord.value.trim();
    isLoading.value = true;
    showResults.value = true;

    try {
      const result = await getWordDefinition(currentWord.value);

      if (result.success) {
        wordMeaning.value = formatDefinition(result);
      } else {
        wordMeaning.value = 'Sorry, the word was not found. Please check the spelling and try again.';
      }
    } catch (error) {
      console.error('Failed to fetch word definition:', error);
      wordMeaning.value = 'Search failed. Please check your network connection and try again.';
    } finally {
      isLoading.value = false;
    }
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

// 播放音频函数
const playAudio = (audioUrl) => {
  try {
    const audio = new Audio(audioUrl);
    audio.play().catch(error => {
      console.error('音频播放失败:', error);
    });
  } catch (error) {
    console.error('音频加载失败:', error);
  }
};
</script>

<style scoped>
/* 基础样式 - 移动端优先 */
.query-container {
  width: 100%;
  max-width: 100%;
  margin: 0 auto;
  padding: 20px;
  box-sizing: border-box;
}

.search-section {
  margin-bottom: 30px;
  display: flex;
  justify-content: center;
}

.search-input {
  width: 100%;
  max-width: 100%;
}

.results-section {
  margin-top: 20px;
}

.welcome-section {
  text-align: center;
  margin-top: 50px;
}

.welcome-card,
.word-card {
  width: 100%;
  max-width: 100%;
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

/* 平板端 - 768px 及以上 */
@media (min-width: 768px) {
  .query-container {
    max-width: 90%;
    padding: 30px;
  }

  .search-input {
    width: 80%;
  }

  .welcome-card,
  .word-card {
    width: 80%;
    max-width: 80%;
    margin: 0 auto;
  }

  .word-title {
    font-size: 28px;
  }

  :deep(.meaning) {
    font-size: 16px;
    line-height: 1.8;
  }

  :deep(.meaning .word-title) {
    font-size: 40px;
  }
}

/* 桌面端 - 1024px 及以上 */
@media (min-width: 1024px) {
  .query-container {
    max-width: 1200px;
    padding: 40px;
  }

  .search-input {
    width: 70%;
    max-width: 800px;
  }

  .welcome-card,
  .word-card {
    width: 70%;
    max-width: 70%;
    margin: 0 auto;
  }

  .word-title {
    font-size: 32px;
  }

  .header-actions {
    gap: 15px;
  }

  :deep(.meaning) {
    font-size: 18px;
    line-height: 1.8;
  }

  :deep(.meaning .word-title) {
    font-size: 48px;
    margin-bottom: 15px;
  }
}

/* ===== 释义内容样式 ===== */
:deep(.meaning) {
  color: var(--el-text-color-regular);
}

:deep(.meaning .word-title) {
  font-weight: bold;
  color: var(--el-color-primary);
  margin-bottom: 10px;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

/* 响应式 iframe 高度调整 */
@media (max-width: 768px) {
}

@media (min-width: 1024px) {
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
