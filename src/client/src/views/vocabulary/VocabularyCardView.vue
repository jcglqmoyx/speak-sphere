<template>
  <ContentBase @keydown="handleKeydown" tabindex="0">
    <!-- Katex CSS for math formula rendering -->
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.css" crossorigin="anonymous">
    <div v-if="dataStatus === 1" class="vocabulary-card">
      <div class="header">
        <el-tooltip content="使用LLM搜索单词" placement="top">
          <el-button v-if="showHeader" class="centered-button" circle @click="toggleLLMPanel">
            <el-icon>
              <ChatDotRound/>
            </el-icon>
          </el-button>
        </el-tooltip>
        &nbsp;
        <el-tooltip content="搜索单词 (快捷键: S)" placement="top">
          <el-button v-if="showHeader" class="centered-button" circle @click="showSearchDrawer = !showSearchDrawer">
            <el-icon>
              <Search/>
            </el-icon>
          </el-button>
        </el-tooltip>
        &nbsp;
        <el-tooltip content="编辑笔记 (支持Markdown语法)" placement="top">
          <el-button v-if="showHeader" class="centered-button" circle @click="showEditNoteDrawerButtonClicked">
            <el-icon>
              <EditPen/>
            </el-icon>
          </el-button>
        </el-tooltip>
        &nbsp;
        <el-tooltip content="标记为不想学 (快捷键: D)" placement="top">
          <el-button v-if="showHeader" class="centered-button" circle @click="markAsUnwanted">
            <el-icon>
              <Delete/>
            </el-icon>
          </el-button>
        </el-tooltip>
      </div>

      <!-- 单词显示区域 -->
      <div class="vocabulary-display">
        <div class="vocabulary-title-section">
          <h1 class="vocabulary">{{ vocabularies[idx].vocabulary }}</h1>
        </div>

        <!-- 释义显示区域 -->
        <div v-if="status === 2 || status === 3" class="meaning-container">
          <div v-if="isFetchingDefinition" class="fetching-meaning">
            正在查询释义...
          </div>
          <div v-else-if="hasMeaning" class="meaning-content" v-html="currentVocabularyMeaning"
               @click="handleVocabularyContentClick"></div>
          <div v-else class="no-meaning">
            暂无释义
          </div>
        </div>
      </div>

      <div class="footer">
        <el-button v-if="status === 1" @click="handleFamiliar">认识(J)</el-button>
        <el-button v-if="status === 1" @click="handleUnfamiliar">不认识(K)</el-button>
        <el-button v-if="status === 2" @click="gotItWrong">记错了(U)</el-button>
        <el-button v-if="status === 2 || status === 3" @click="showNextVocabulary">下一词(N)</el-button>
      </div>
    </div>
    <div v-else-if="dataStatus === 2" class="vocabulary-card">
      <h1 class="vocabulary">Loading data...</h1>
    </div>
    <div v-else-if="dataStatus === 3" class="vocabulary-card">
      <h1 class="vocabulary">当前任务已完成</h1>
    </div>

    <!-- Search Drawer 必须置于ContentBase中，如果不放在ContentBase中，则按下s键（搜索单词）后，ContentBase将失去焦点，再次按快捷键将无效 -->
    <el-drawer v-model="showSearchDrawer" title="搜索" :with-header="true">
      <p v-for="(dictionary, index) in dictionaries" :key="dictionary.id">
        <a href="#" @click.prevent="handleDictionaryClick(dictionary)"
           :class="{'wiktionary-link': dictionary.title.toLowerCase().includes('wiktionary')}">
          <b>{{ index + 1 }}. {{ dictionary.title }}</b>
        </a>
        <el-divider/>
      </p>
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
  </ContentBase>

  <el-drawer v-model="showEditNoteDrawer" title="笔记 (支持Markdown语法)" :with-header="true" size="50%">
    <div class="note-editor">
      <el-tabs v-model="activeTab" type="border-card">
        <el-tab-pane label="编辑" name="edit">
          <el-input
              type="textarea"
              placeholder="添加笔记 (支持Markdown语法)"
              v-model="currentVocabularyNote"
              :rows="20"
              maxlength="1000"
              show-vocabulary-limit
          ></el-input>
        </el-tab-pane>
        <el-tab-pane label="预览" name="preview">
          <div class="markdown-preview" v-html="renderedMarkdown"></div>
        </el-tab-pane>
      </el-tabs>
      <el-divider/>
      <div class="editor-actions">
        <el-button type="primary" @click="updateNote">保存</el-button>
        <el-button @click="showEditNoteDrawer = false">取消</el-button>
      </div>
    </div>
  </el-drawer>

  <!-- LLM查询面板 -->
  <LLMQueryPanel
      v-if="showLLMPanel"
      :current-vocabulary="getCurrentVocabulary()"
      @close="showLLMPanel = false"
  />
</template>

<script setup>
import {
  ElButton,
  ElDivider,
  ElDrawer,
  ElIcon,
  ElInput,
  ElNotification,
  ElTabPane,
  ElTabs,
  ElTooltip
} from 'element-plus';
import {ChatDotRound, Close, Delete, EditPen, Search} from "@element-plus/icons-vue";
import ContentBase from "@/components/ContentBase.vue";
import LLMQueryPanel from "@/components/LLMQueryPanel.vue";
import MarkdownIt from 'markdown-it';

import {useStore} from "vuex";
import {computed, defineProps, onMounted, reactive, ref} from "vue";

import {getUserProfile} from "@/assets/js/module/user/query";
import {getDictionaryList} from "@/assets/js/module/dictionary/query";
import {fetchVocabularies} from "@/assets/js/module/vocabulary/query";
import {formatDefinition, getVocabularyDefinition} from "@/assets/js/util/dictionary_api";
import {
  resetVocabularyStudyCountToZero,
  setUnwanted,
  updateVocabulary,
  updateVocabularyStudyCount
} from "@/assets/js/module/vocabulary/update";
// 初始化markdown渲染器
import MarkdownItMark from 'markdown-it-mark';
import MarkdownItFootnote from 'markdown-it-footnote';
import MarkdownItAbbr from 'markdown-it-abbr';
import MarkdownItIns from 'markdown-it-ins';
import MarkdownItSub from 'markdown-it-sub';
import MarkdownItSup from 'markdown-it-sup';
import MarkdownItTexmath from 'markdown-it-texmath';
import katex from 'katex';

const props = defineProps(["type"]);
/*
 * 1: 加载成功
 * 2: 加载中
 * 3: 学习或复习任务已完成
 */
const dataStatus = ref(2);

const showHeader = ref(true);
const showSearchDrawer = ref(false);
const showEditNoteDrawer = ref(false);
const showLLMPanel = ref(false);
const showWiktionaryFullscreen = ref(false);
const currentWiktionaryUrl = ref('');
const activeTab = ref('edit');
const isFetchingDefinition = ref(false);

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  tables: true
})
    .use(MarkdownItMark)
    .use(MarkdownItFootnote)
    .use(MarkdownItAbbr)
    .use(MarkdownItIns)
    .use(MarkdownItSub)
    .use(MarkdownItSup)
    .use(MarkdownItTexmath, {
      engine: katex,
      delimiters: 'dollars',
      katexOptions: {
        macros: {
          "\\RR": "\\mathbb{R}"
        }
      }
    });

const showEditNoteDrawerButtonClicked = () => {
  setCurrentVocabularyNote();
  activeTab.value = 'edit';
  showEditNoteDrawer.value = !showEditNoteDrawer.value;
}

// 计算属性：渲染markdown内容
const renderedMarkdown = computed(() => {
  if (!currentVocabularyNote.value) {
    return '<p class="empty-note">暂无笔记内容</p>';
  }
  return md.render(currentVocabularyNote.value);
});

let dictionaries = reactive([]);

let vocabularies = reactive([]);
let idx = ref(0);
let countRecognitionTime = {};
const timesCountedAsKnown = ref(0);

// 获取当前单词
const getCurrentVocabulary = () => {
  if (vocabularies.length > 0 && idx.value < vocabularies.length) {
    return vocabularies[idx.value].vocabulary || '';
  }
  return '';
}

const toggleLLMPanel = () => {
  showLLMPanel.value = !showLLMPanel.value;
}

// 处理词典点击事件
const handleDictionaryClick = (dictionary) => {
  const url = dictionary.prefix + vocabularies[idx.value].vocabulary + (dictionary.suffix || '');

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
}

// 全局函数：打开全屏Wiktionary
window.openWiktionaryFullscreen = (url) => {
  showWiktionaryFullscreen.value = true;
  currentWiktionaryUrl.value = url;
};

/*
 * 1: 刚开始的页面
 * 2: 选择了认识
 * 3: 选择了不认识
 */
const status = ref(1);

const checkResponse = (response) => {
  if (response == null || response.code === 2) {
    const store = useStore();
    store.dispatch("logout")
    location.reload();
  }
}

onMounted(
    async () => {
      const fetchVocabularysResponse = await fetchVocabularies(props.type);
      checkResponse(fetchVocabularysResponse)
      vocabularies = fetchVocabularysResponse.data;
      if (vocabularies == null || vocabularies.length === 0) {
        dataStatus.value = 3;
        return;
      } else {
        dataStatus.value = 1;
        for (let i = 0; i < vocabularies.length; i++) {
          countRecognitionTime[i] = 0;
        }
      }

      const getDictionaryListResponse = await getDictionaryList();
      checkResponse(getDictionaryListResponse);
      dictionaries = getDictionaryListResponse.data;

      const getUserProfileResponse = await getUserProfile();
      checkResponse(getUserProfileResponse);
      timesCountedAsKnown.value = getUserProfileResponse.times_counted_as_known;

      // 初始化词典查询
      if (dataStatus.value === 1 && vocabularies.length > 0) {
        await checkAndFetchDefinition();
      }
    }
);

const handleKeydown = async (event) => {
  if (!showEditNoteDrawer.value) {
    if ('0123456789'.includes(event.key)) {
      let number = +event.key;
      if (number === 0) number = 10;
      if (number <= dictionaries.length && showSearchDrawer.value) {
        window.open(dictionaries[number - 1].prefix + vocabularies[idx.value].vocabulary + (dictionaries[number - 1].suffix || ""), '_blank');
      }
    } else if (event.key === 's') {
      showSearchDrawer.value = !showSearchDrawer.value;
    } else if (event.key === "n") {
      if (status.value !== 2 && status.value !== 3) {
        return true;
      }
      await showNextVocabulary();
    } else if (event.key === "u") {
      if (status.value !== 2) {
        return true;
      }
      await gotItWrong();
    } else if (event.key === "j") {
      if (status.value !== 1) {
        return true;
      }
      await handleFamiliar();
    } else if (event.key === "k") {
      if (status.value !== 1) {
        return true;
      }
      await handleUnfamiliar();
    } else if (event.key === "d" || event.key === "D") {
      await markAsUnwanted();
    } else {
      return true;
    }
  }
};

const markAsUnwanted = async () => {
  countRecognitionTime[idx.value] = timesCountedAsKnown.value;
  const markedAsUnwantedResponse = await setUnwanted(vocabularies[idx.value].id);
  checkResponse(markedAsUnwantedResponse);
  await showNextVocabulary();
};

const handleUnfamiliar = async () => {
  countRecognitionTime[idx.value] = 0;
  const resetVocabularyStudyCountToZeroResponse = await resetVocabularyStudyCountToZero(vocabularies[idx.value].id);
  checkResponse(resetVocabularyStudyCountToZeroResponse);
  status.value = 3;
};

const handleFamiliar = async () => {
  countRecognitionTime[idx.value]++;
  if (countRecognitionTime[idx.value] >= timesCountedAsKnown.value) {
    const updateVocabularyStudyCountResponse = await updateVocabularyStudyCount(vocabularies[idx.value].id);
    checkResponse(updateVocabularyStudyCountResponse);
  }
  status.value = 2;
};

const gotItWrong = async () => {
  countRecognitionTime[idx.value] = 0;
  const resetVocabularyStudyCountToZeroResponse = await resetVocabularyStudyCountToZero(vocabularies[idx.value].id);
  checkResponse(resetVocabularyStudyCountToZeroResponse);
  await showNextVocabulary();
}

let currentVocabularyNote = ref('');
let currentVocabularyMeaning = ref('');
const setCurrentVocabularyNote = () => {
  currentVocabularyNote.value = vocabularies[idx.value].note;
};

// 计算当前单词是否有释义
const hasMeaning = computed(() => {
  // 检查当前显示的释义，包括API查询到的
  const vocabularyMeaning = vocabularies[idx.value]?.meaning || currentVocabularyMeaning.value;
  return vocabularyMeaning && vocabularyMeaning.trim() !== '';
});

const checkAndFetchDefinition = async () => {
  const currentVocabulary = vocabularies[idx.value];
  if (!currentVocabulary) return;

  const meaning = currentVocabulary.meaning || '';
  // 如果没有释义，调用API查询
  if (!meaning || meaning.trim() === '' && !isFetchingDefinition.value) {
    isFetchingDefinition.value = true;
    currentVocabularyMeaning.value = '';

    try {
      const result = await getVocabularyDefinition(currentVocabulary.vocabulary);

      if (result.success) {
        currentVocabularyMeaning.value = formatDefinition(result);
      } else {
        currentVocabularyMeaning.value = '暂无释义';
      }
    } catch (error) {
      currentVocabularyMeaning.value = '暂无释义';
    } finally {
      isFetchingDefinition.value = false;
    }
  } else {
    currentVocabularyMeaning.value = meaning;
    console.log(`单词 "${currentVocabulary.vocabulary}" 已有释义:`, meaning);
  }
};

// 初始化当前单词的释义查询
const showNextVocabulary = async () => {
  let tempIdx = idx.value;
  do {
    tempIdx = (tempIdx + 1) % vocabularies.length;
  } while (tempIdx !== idx.value && countRecognitionTime[tempIdx] >= timesCountedAsKnown.value);
  if (tempIdx === idx.value && countRecognitionTime[tempIdx] >= timesCountedAsKnown.value) {
    dataStatus.value = 3;
  } else {
    status.value = 1;
    idx.value = tempIdx;

    // 在切换单词后检查是否需要查询释义
    await checkAndFetchDefinition();
  }
}

const updateNote = async () => {
  const updateNoteResponse = await updateVocabulary(vocabularies[idx.value].id, vocabularies[idx.value].vocabulary, vocabularies[idx.value].meaning, vocabularies[idx.value].vocabulary_set_id, currentVocabularyNote.value, vocabularies[idx.value].unwanted, vocabularies[idx.value].study_count, vocabularies[idx.value].date_to_review, vocabularies[idx.value].created_at);
  checkResponse(updateNoteResponse);
  if (updateNoteResponse != null && updateNoteResponse.code === 0) {
    ElNotification({
      title: 'Success',
      message: '更新成功',
      type: 'success',
      duration: 1000,
    });
  }
}

// 处理单词内容区域的点击事件 - 用于发音按钮
const handleVocabularyContentClick = (event) => {
  // 检查点击的是否是发音按钮
  if (event.target.classList.contains('audio-btn')) {
    const audioUrl = event.target.getAttribute('data-audio');
    if (audioUrl) {
      playAudio(audioUrl);
    }
  }
};

// 播放音频
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
.vocabulary-card {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 85vh;
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.header {
  display: flex;
  justify-content: flex-end;
  padding: 10px;
}

.vocabulary-display {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-start;
  padding: 20px;
}

.vocabulary-title-section {
  margin-bottom: 30px;
  text-align: center;
}

.vocabulary {
  font-size: 48px;
  font-weight: bold;
  color: var(--el-color-primary);
  margin: 0;
  text-align: center;
  line-height: 1.2;
}

.meaning-container {
  width: 100%;
  /* 使用calc是为了使字体在桌面和手机上的大小能比较均衡 */
  font-size: calc(1vw + 1.5vh);
  min-height: 400px;
  max-height: 600px;
  overflow-y: auto;
  padding: 25px;
  background: var(--el-bg-color-page);
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  border: 1px solid var(--el-border-color-light);

  /* 优化滚动条样式 */
  scrollbar-width: thin;
  scrollbar-color: var(--el-color-primary-light-5) var(--el-bg-color-page);
}

/* Webkit浏览器滚动条样式 */
.meaning-container::-webkit-scrollbar {
  width: 8px;
}

.meaning-container::-webkit-scrollbar-track {
  background: var(--el-bg-color-page);
  border-radius: 4px;
}

.meaning-container::-webkit-scrollbar-thumb {
  background: var(--el-color-primary-light-5);
  border-radius: 4px;
}

.meaning-container::-webkit-scrollbar-thumb:hover {
  background: var(--el-color-primary);
}

.meaning-content {
  width: 100%;
}

.fetching-meaning {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 200px;
  font-size: 18px;
  color: var(--el-text-color-secondary);
  font-style: italic;
}

.no-meaning {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 200px;
  font-size: 18px;
  color: var(--el-text-color-placeholder);
  font-style: italic;
}

.footer {
  display: flex;
  justify-content: space-around;
  padding: 10px;
}

.centered-button {
  display: grid;
  align-items: center;
  justify-content: center;
  height: 5vh;
  width: 5vh;
}

.note-editor {
  height: calc(100vh - 60px);
  display: flex;
  flex-direction: column;
}

.empty-note {
  color: var(--el-text-color-placeholder);
  font-style: italic;
  text-align: center;
  padding: 20px;
}

.markdown-preview {
  height: 100%;
  overflow-y: auto;
  padding: 10px;
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  background-color: var(--el-bg-color);
}

.editor-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 10px;
}

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

/* 移动端Wiktionary适配 */
@media (max-width: 768px) {
  :deep(.wiktionary-container) {
    height: 300px;
  }

  /* 为Wiktionary iframe添加移动端优化 */
  :deep(.wiktionary-container iframe) {
    width: 100%;
    height: 100%;
    border: none;
    -webkit-overflow-scrolling: touch; /* iOS平滑滚动 */
  }
}

/* 小屏手机Wiktionary适配 */
@media (max-width: 480px) {
  :deep(.wiktionary-container) {
    height: 250px;
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

/* 移动端 (≤768px) */
@media (max-width: 768px) {
  .meaning-container {
    min-height: 350px;
    max-height: 500px;
    padding: 15px;
  }

  .vocabulary {
    font-size: 36px;
  }
}

/* 桌面端 (≥1024px) */
@media (min-width: 1024px) {
  .meaning-container {
    max-height: 700px;
  }
}

/* 大桌面端 (≥1200px) */
@media (min-width: 1200px) {
  .meaning-container {
    max-height: 800px;
  }
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
</style>