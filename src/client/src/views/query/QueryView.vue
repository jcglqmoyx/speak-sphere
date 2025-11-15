<template>
  <ContentBase @keydown="handleKeydown" tabindex="0">
    <div class="query-container">
      <div class="search-section">
        <el-input
            v-model="searchVocabulary"
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
        <el-card class="vocabulary-card">
          <template #header>
            <div class="card-header">
              <span class="vocabulary-title">{{ currentVocabulary }}</span>
              <div class="header-actions">
                <el-tooltip content="使用LLM搜索单词" placement="top">
                  <el-button circle @click="toggleLLMPanel">
                    <el-icon>
                      <ChatDotRound/>
                    </el-icon>
                  </el-button>
                </el-tooltip>
                <el-tooltip content="收藏单词" placement="top">
                  <el-button circle @click="handleAddToVocabularySet" :disabled="!searchVocabulary.trim()"
                             class="favorite-btn">
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

          <div class="vocabulary-content" @click="handleVocabularyContentClick">
            <div v-if="isLoading" class="loading-container">
              <el-skeleton :rows="5" animated/>
            </div>
            <div v-else-if="vocabularyMeaning" class="meaning" v-html="vocabularyMeaning"></div>
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
          <a href="#" @click.prevent="handleDictionaryClick(dictionary)"
             :class="{'wiktionary-link': dictionary.title.toLowerCase().includes('wiktionary')}">
            <strong>{{ index + 1 }}. {{ dictionary.title }}</strong>
          </a>
          <el-divider/>
        </p>
      </div>
    </el-drawer>

    <!-- 全屏Wiktionary模态框 -->
    <div v-if="showWiktionaryFullscreen" class="wiktionary-fullscreen-modal">
      <div class="wiktionary-fullscreen-content">
        <button class="close-fullscreen-btn" @click="showWiktionaryFullscreen = false">
          <el-icon>
            <Close/>
          </el-icon>
        </button>
        <iframe
            :src="currentWiktionaryUrl"
            class="wiktionary-fullscreen-iframe"
            loading="lazy"
        ></iframe>
      </div>
    </div>

    <!-- 收藏单词对话框 -->
    <el-dialog v-model="showAddToVocabularySetDialog" title="收藏单词" width="500px">
      <div class="dialog-title">将单词 "<strong>{{ searchVocabulary }}</strong>" 添加到词书</div>
      <el-form>
        <el-form-item label="选择词书">
          <el-radio-group v-model="selectedVocabularySetId" class="vocabulary-set-radio-group">
            <el-radio v-for="vocabularySet in vocabularySets" :key="vocabularySet.id" :label="vocabularySet.id"
                      class="vocabulary-set-radio">
              <div class="vocabulary-set-radio-content">
                <span>{{ vocabularySet.title }}</span>
                <el-tag v-if="vocabularySetVocabularyExistsMap[vocabularySet.id]" type="success" size="small"
                        class="exists-tag">
                  已存在
                </el-tag>
              </div>
            </el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showAddToVocabularySetDialog = false">取消</el-button>
          <el-button
              type="primary"
              @click="confirmAddToVocabularySet"
              :disabled="!selectedVocabularySetId"
              :class="{ 'warning-btn': selectedVocabularySetId && isVocabularyInSelectedVocabularySet() }"
          >
            {{ selectedVocabularySetId && isVocabularyInSelectedVocabularySet() ? '无需重复添加' : '确认添加' }}
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- LLM查询面板 -->
    <LLMQueryPanel
        v-if="showLLMPanel"
        :current-vocabulary="currentVocabulary"
        @close="showLLMPanel = false"
    />
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
import {ChatDotRound, Close, Search, Star} from '@element-plus/icons-vue';
import ContentBase from '@/components/ContentBase.vue';
import LLMQueryPanel from '@/components/LLMQueryPanel.vue';
import {getDictionaryList} from '@/assets/js/module/dictionary/query';
import {getVocabularySetList} from '@/assets/js/module/vocabulary_set/query';
import {AddVocabulary} from '@/assets/js/module/vocabulary/add';
import {checkVocabularyInVocabularySet} from '@/assets/js/module/vocabulary/query';
import {formatDefinition, getVocabularyDefinition} from '@/assets/js/util/dictionary_api';

const searchVocabulary = ref('');
const currentVocabulary = ref('');
const vocabularyMeaning = ref('');
const isLoading = ref(false);
const showResults = ref(false);
const showSearchDrawer = ref(false);
const showEditNoteDrawer = ref(false);
const showAddToVocabularySetDialog = ref(false);
const showLLMPanel = ref(false);
const showWiktionaryFullscreen = ref(false);
const currentWiktionaryUrl = ref('');
const dictionaries = ref([]);
const vocabularySets = ref([]);
const selectedVocabularySetId = ref(null);
const vocabularySetVocabularyExistsMap = ref({}); // 存储每个词书是否包含当前单词

onMounted(async () => {
  // 加载词典列表
  const dictionaryResponse = await getDictionaryList();
  if (dictionaryResponse && dictionaryResponse.data) {
    dictionaries.value = dictionaryResponse.data;
  }

  // 加载词书列表
  const vocabularySetResponse = await getVocabularySetList(100000000, 1);
  if (vocabularySetResponse && vocabularySetResponse.data) {
    vocabularySets.value = vocabularySetResponse.data;
  }
});

// 处理单词内容区域的点击事件
const handleVocabularyContentClick = (event) => {
  // 检查点击的是否是发音按钮
  if (event.target.classList.contains('audio-btn')) {
    const audioUrl = event.target.getAttribute('data-audio');
    if (audioUrl) {
      playAudio(audioUrl);
    }
  }
};

const handleSearch = async () => {
  if (searchVocabulary.value.trim()) {
    currentVocabulary.value = searchVocabulary.value.trim();
    isLoading.value = true;
    showResults.value = true;

    try {
      const result = await getVocabularyDefinition(currentVocabulary.value);

      if (result.success) {
        vocabularyMeaning.value = formatDefinition(result);
      } else {
        vocabularyMeaning.value = 'Sorry, the vocabulary was not found. Please check the spelling and try again.';
      }
    } catch (error) {
      console.error('Failed to fetch vocabulary definition:', error);
      vocabularyMeaning.value = 'Search failed. Please check your network connection and try again.';
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
      window.open(dictionary.prefix + (currentVocabulary.value || '') + (dictionary.suffix || ''), '_blank');
    }
  }
};

const handleAddToVocabularySet = async () => {
  if (!searchVocabulary.value.trim()) {
    return;
  }

  selectedVocabularySetId.value = null;
  vocabularySetVocabularyExistsMap.value = {}; // 清空之前的记录
  showAddToVocabularySetDialog.value = true;

  // 检查每个词书是否已包含当前单词
  if (vocabularySets.value.length > 0) {
    const vocabulary = searchVocabulary.value.trim();

    for (const vocabularySet of vocabularySets.value) {
      try {
        const response = await checkVocabularyInVocabularySet(vocabulary, vocabularySet.id);
        const exists = response && response.code === 0 && response.data.exists;
        vocabularySetVocabularyExistsMap.value[vocabularySet.id] = exists;

        if (exists && !selectedVocabularySetId.value) {
          selectedVocabularySetId.value = vocabularySet.id;
        }
      } catch (error) {
        console.error(`检查词书 ${vocabularySet.title} 失败:`, error);
        vocabularySetVocabularyExistsMap.value[vocabularySet.id] = false;
      }
    }
  }
};

const isVocabularyInSelectedVocabularySet = () => {
  if (!selectedVocabularySetId.value) {
    return false;
  }
  return !!vocabularySetVocabularyExistsMap.value[selectedVocabularySetId.value];
};

const confirmAddToVocabularySet = async () => {
  if (!selectedVocabularySetId.value || !searchVocabulary.value.trim()) {
    return;
  }

  try {
    const response = await AddVocabulary(
        selectedVocabularySetId.value,
        searchVocabulary.value.trim(),
        vocabularyMeaning.value || '暂无释义', // 使用当前显示的释义，如果没有则使用默认值
        ''
    );

    if (response && response.code === 0) {
      ElMessage.success('添加单词成功');
      showAddToVocabularySetDialog.value = false;
    } else {
      const errorMsg = response ? response.message : '添加失败，请重试';
      ElMessage.error(errorMsg);
    }
  } catch (error) {
    console.error('添加单词失败:', error);
    ElMessage.error('添加单词失败，请重试');
  }
};

// 切换LLM查询面板
const toggleLLMPanel = () => {
  if (showResults.value && currentVocabulary.value) {
    showLLMPanel.value = !showLLMPanel.value;
  }
};

// 处理词典点击事件
const handleDictionaryClick = (dictionary) => {
  const url = dictionary.prefix + (currentVocabulary.value || '') + (dictionary.suffix || '');

  // 如果是Wiktionary并且在移动端，使用全屏模态框
  const isWiktionary = dictionary.title.toLowerCase().includes('wiktionary');
  const isMobile = window.innerWidth <= 768;

  if (isWiktionary && isMobile) {
    currentWiktionaryUrl.value = url;
    showWiktionaryFullscreen.value = true;
    showSearchDrawer.value = false; // 关闭抽屉
  } else {
    // 非移动端或非Wiktionary，在新标签页打开
    window.open(url, '_blank');
  }
};

// 全局函数：打开全屏Wiktionary
window.openWiktionaryFullscreen = (url) => {
  showWiktionaryFullscreen.value = true;
  currentWiktionaryUrl.value = url;
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
.vocabulary-card {
  width: 100%;
  max-width: 100%;
  margin: 0 auto;
}

.vocabulary-card {
  margin: 20px 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.vocabulary-title {
  font-size: 24px;
  font-weight: bold;
  color: var(--el-text-color-primary);
}

.header-actions {
  display: flex;
  gap: 10px;
}

.vocabulary-content {
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
  .vocabulary-card {
    width: 80%;
    max-width: 80%;
    margin: 0 auto;
  }

  .vocabulary-title {
    font-size: 28px;
  }

  :deep(.meaning) {
    font-size: 16px;
    line-height: 1.8;
  }

  :deep(.meaning .vocabulary-title) {
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
  .vocabulary-card {
    width: 70%;
    max-width: 70%;
    margin: 0 auto;
  }

  .vocabulary-title {
    font-size: 32px;
  }

  .header-actions {
    gap: 15px;
  }

  :deep(.meaning) {
    font-size: 18px;
    line-height: 1.8;
  }

  :deep(.meaning .vocabulary-title) {
    font-size: 48px;
    margin-bottom: 15px;
  }
}

/* ===== 释义内容样式 ===== */
:deep(.meaning) {
  color: var(--el-text-color-regular);
}

:deep(.meaning .vocabulary-title) {
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

/* 全屏Wiktionary模态框样式 */
.wiktionary-fullscreen-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(0, 0, 0, 0.8);
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
}

.wiktionary-fullscreen-content {
  position: relative;
  width: 95%;
  height: 95%;
  background: var(--el-bg-color);
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
}

.close-fullscreen-btn {
  position: absolute;
  top: 15px;
  right: 15px;
  width: 40px;
  height: 40px;
  background: var(--el-color-primary);
  color: white;
  border: none;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  cursor: pointer;
  z-index: 10000;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
  transition: all 0.3s ease;
}

.close-fullscreen-btn:hover {
  background: var(--el-color-primary-light-3);
  transform: scale(1.1);
}

.wiktionary-fullscreen-iframe {
  width: 100%;
  height: 100%;
  border: none;
  border-radius: 8px;
  -webkit-overflow-scrolling: touch;
}

/* 平板端适配 */
@media (max-width: 1024px) {
  .wiktionary-fullscreen-content {
    width: 98%;
    height: 98%;
    border-radius: 6px;
  }

  .close-fullscreen-btn {
    width: 45px;
    height: 45px;
    font-size: 22px;
  }
}

/* 移动端适配 */
@media (max-width: 768px) {
  .wiktionary-fullscreen-content {
    width: 100%;
    height: 100%;
    border-radius: 0;
  }

  .close-fullscreen-btn {
    width: 50px;
    height: 50px;
    font-size: 24px;
    top: 10px;
    right: 10px;
  }
}

/* Wiktionary iframe 样式 */
:deep(.wiktionary-container) {
  width: 100%;
  height: 400px;
  margin-top: 10px;
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  overflow: hidden;
}

:deep(.wiktionary-notice) {
  background-color: var(--el-color-warning-light-8);
  color: var(--el-color-warning-dark-2);
  padding: 8px 12px;
  font-size: 12px;
  font-weight: 500;
  border-bottom: 1px solid var(--el-border-color);
}

:deep(.wiktionary-loading) {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  flex-direction: column;
  color: var(--el-text-color-secondary);
  gap: 10px;
}

:deep(.wiktionary-loading .loading-spinner) {
  width: 40px;
  height: 40px;
  border: 3px solid var(--el-border-color-light);
  border-top: 3px solid var(--el-color-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

:deep(.wiktionary-iframe) {
  width: 100%;
  height: calc(100% - 40px);
  border: none;
}

/* 响应式调整 */
@media (min-width: 768px) {
  :deep(.wiktionary-container) {
    height: 500px;
  }
}

@media (min-width: 1024px) {
  :deep(.wiktionary-container) {
    height: 600px;
  }
}

@media (min-width: 768px) {
  :deep(.wiktionary-container) {
    height: 500px;
  }
}

@media (min-width: 1024px) {
  :deep(.wiktionary-container) {
    height: 600px;
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

.vocabulary-set-radio-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.vocabulary-set-radio {
  display: flex;
  align-items: center;
  padding: 8px;
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  transition: all 0.3s;
}

.vocabulary-set-radio:hover {
  border-color: var(--el-color-primary);
  background-color: var(--el-fill-color-lighter);
}

.dialog-title {
  margin-bottom: 20px;
  font-size: 16px;
  color: var(--el-text-color-primary);
}

.vocabulary-set-radio-content {
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
