<template>
  <ContentBase @keydown="handleKeydown" tabindex="0">
    <!-- Katex CSS for math formula rendering -->
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/katex@0.16.9/dist/katex.min.css" crossorigin="anonymous">
    <div v-if="dataStatus === 1" class="word-card">
      <div class="header">
        <el-tooltip content="搜索单词 (快捷键: S)" placement="top">
          <el-button v-if="showHeader" class="centered-button" circle @click="showSearchDrawer = !showSearchDrawer">
            <el-icon>
              <el-icon>
                <Search/>
              </el-icon>
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

      <div>
        <div class="word">{{ words[idx].word }}</div>
        <div v-if="status === 2 || status === 3" class="word">{{ words[idx].meaning }}</div>
      </div>

      <div class="footer">
        <el-button v-if="status === 1" @click="handleFamiliar">认识(J)</el-button>
        <el-button v-if="status === 1" @click="handleUnfamiliar">不认识(K)</el-button>
        <el-button v-if="status === 2" @click="gotItWrong">记错了(U)</el-button>
        <el-button v-if="status === 2 || status === 3" @click="showNextWord">下一词(N)</el-button>
      </div>
    </div>
    <div v-else-if="dataStatus === 2" class="word-card">
      <h1 class="word">Loading data...</h1>
    </div>
    <div v-else-if="dataStatus === 3" class="word-card">
      <h1 class="word">当前任务已完成</h1>
    </div>

    <!-- Search Drawer 必须置于ContentBase中，如果不放在ContentBase中，则按下s键（搜索单词）后，ContentBase将失去焦点，再次按快捷键将无效 -->
    <el-drawer v-model="showSearchDrawer" title="搜索" :with-header="true">
      <p v-for="(dictionary, index) in dictionaries" :key="dictionary.id">
        <a :href="dictionary.prefix + words[idx].word + (dictionary.suffix || '') " target="_blank">
          <b>{{ index + 1 }}. {{ dictionary.title }}</b>
        </a>
        <el-divider/>
      </p>
    </el-drawer>
  </ContentBase>

  <el-drawer v-model="showEditNoteDrawer" title="笔记 (支持Markdown语法)" :with-header="true" size="50%">
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
        <el-button type="primary" @click="updateNote">保存</el-button>
        <el-button @click="showEditNoteDrawer = false">取消</el-button>
      </div>
    </div>
  </el-drawer>
</template>

<script setup>
import {ElButton, ElDivider, ElDrawer, ElIcon, ElInput, ElTabPane, ElTabs, ElTooltip} from 'element-plus';
import {Delete, EditPen, Search} from "@element-plus/icons-vue";
import ContentBase from "@/components/ContentBase.vue";
import MarkdownIt from 'markdown-it';

import {useStore} from "vuex";
import {computed, defineProps, onMounted, reactive, ref} from "vue";

import {getUserProfile} from "@/assets/js/module/user/query";
import {getDictionaryList} from "@/assets/js/module/dictionary/query";
import {fetchWords} from "@/assets/js/module/entry/query";
import {
  resetEntryStudyCountToZero,
  setUnwanted,
  updateEntry,
  updateEntryStudyCount
} from "@/assets/js/module/entry/update";
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
const activeTab = ref('edit');

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  // 启用表格支持
  tables: true
})
    .use(MarkdownItMark) // 支持标记文本 ==标记==
    .use(MarkdownItFootnote) // 支持脚注
    .use(MarkdownItAbbr) // 支持缩写
    .use(MarkdownItIns) // 支持下划线 ++下划线++
    .use(MarkdownItSub) // 支持下标 H~2~O
    .use(MarkdownItSup) // 支持上标 x^2^
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
  setCurrentWordNote();
  activeTab.value = 'edit';
  showEditNoteDrawer.value = !showEditNoteDrawer.value;
}

// 计算属性：渲染markdown内容
const renderedMarkdown = computed(() => {
  if (!currentWordNote.value) {
    return '<p class="empty-note">暂无笔记内容</p>';
  }
  return md.render(currentWordNote.value);
});

let dictionaries = reactive([]);

let words = reactive([]);
let idx = ref(0);
let countRecognitionTime = {};
const timesCountedAsKnown = ref(0);

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
      const fetchWordsResponse = await fetchWords(props.type);
      checkResponse(fetchWordsResponse)
      words = fetchWordsResponse.data;
      if (words == null || words.length === 0) {
        dataStatus.value = 3;
        return;
      } else {
        dataStatus.value = 1;
        for (let i = 0; i < words.length; i++) {
          countRecognitionTime[i] = 0;
        }
      }

      const getDictionaryListResponse = await getDictionaryList();
      checkResponse(getDictionaryListResponse);
      dictionaries = getDictionaryListResponse.data;

      const getUserProfileResponse = await getUserProfile();
      checkResponse(getUserProfileResponse);
      timesCountedAsKnown.value = getUserProfileResponse.times_counted_as_known;
    }
);

const handleKeydown = async (event) => {
  if (!showEditNoteDrawer.value) {
    if ('0123456789'.includes(event.key)) {
      let number = +event.key;
      if (number === 0) number = 10;
      if (number <= dictionaries.length) {
        window.open(dictionaries[number - 1].prefix + words[idx.value].word + (dictionaries[number - 1].suffix || ""), '_blank');
      }
    } else if (event.key === 's') {
      showSearchDrawer.value = !showSearchDrawer.value;
    } else if (event.key === "n") {
      if (status.value !== 2 && status.value !== 3) {
        return true;
      }
      await showNextWord();
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
  const markedAsUnwantedResponse = await setUnwanted(words[idx.value].id);
  checkResponse(markedAsUnwantedResponse);
  await showNextWord();
};

const handleUnfamiliar = async () => {
  countRecognitionTime[idx.value] = 0;
  const resetEntryStudyCountToZeroResponse = await resetEntryStudyCountToZero(words[idx.value].id);
  checkResponse(resetEntryStudyCountToZeroResponse);
  status.value = 3;
};

const handleFamiliar = async () => {
  countRecognitionTime[idx.value]++;
  if (countRecognitionTime[idx.value] >= timesCountedAsKnown.value) {
    const updateEntryStudyCountResponse = await updateEntryStudyCount(words[idx.value].id);
    checkResponse(updateEntryStudyCountResponse);
  }
  status.value = 2;
};

const gotItWrong = async () => {
  countRecognitionTime[idx.value] = 0;
  const resetEntryStudyCountToZeroResponse = await resetEntryStudyCountToZero(words[idx.value].id);
  checkResponse(resetEntryStudyCountToZeroResponse);
  await showNextWord();
}

const showNextWord = async () => {
  let tempIdx = idx.value;
  do {
    tempIdx = (tempIdx + 1) % words.length;
  } while (tempIdx !== idx.value && countRecognitionTime[tempIdx] >= timesCountedAsKnown.value);
  if (tempIdx === idx.value && countRecognitionTime[tempIdx] >= timesCountedAsKnown.value) {
    dataStatus.value = 3;
  } else {
    status.value = 1;
    idx.value = tempIdx;
  }
}

let currentWordNote = ref('');
const setCurrentWordNote = () => {
  currentWordNote.value = words[idx.value].note;
};

const updateNote = async () => {
  const updateNoteResponse = await updateEntry(words[idx.value].id, words[idx.value].word, words[idx.value].meaning, words[idx.value].book_id, currentWordNote.value, words[idx.value].unwanted, words[idx.value].study_count, words[idx.value].date_to_review, words[idx.value].created_at);
  checkResponse(updateNoteResponse);
}
</script>

<style scoped>
.word-card {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 80vh;
  max-width: 400px;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: flex-end;
  padding: 10px;
}

.word {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.footer {
  display: flex;
  justify-content: space-around;
  padding: 10px;
}

el-icon {
  height: 100%;
  width: 100%;
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

.markdown-preview :deep(h4) {
  color: var(--el-text-color-primary);
}

.markdown-preview :deep(h5) {
  color: var(--el-text-color-primary);
}

.markdown-preview :deep(h6) {
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

.markdown-preview :deep(strong) {
  color: var(--el-text-color-primary);
  font-weight: 600;
}

.markdown-preview :deep(em) {
  color: var(--el-text-color-primary);
  font-style: italic;
}

/* Katex数学公式样式 */
.markdown-preview :deep(.katex) {
  font-size: 1.1em;
}

.markdown-preview :deep(.katex-display) {
  margin: 1em 0;
  overflow-x: auto;
  overflow-y: hidden;
}

.markdown-preview :deep(.katex-display > .katex) {
  display: block;
  text-align: center;
}

/* 增强表格样式 */
.markdown-preview :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 10px 0;
  border: 1px solid var(--el-border-color);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  background-color: var(--el-bg-color);
}

.markdown-preview :deep(th), .markdown-preview :deep(td) {
  border: 1px solid var(--el-border-color);
  padding: 12px;
  text-align: left;
  line-height: 1.4;
}

.markdown-preview :deep(th) {
  background-color: var(--el-fill-color-light);
  color: var(--el-text-color-primary);
  font-weight: 600;
  text-align: center;
}

.markdown-preview :deep(td) {
  color: var(--el-text-color-primary);
  background-color: var(--el-bg-color);
}

.markdown-preview :deep(tr:nth-child(even)) {
  background-color: var(--el-fill-color-lighter);
}

.markdown-preview :deep(tr:hover) {
  background-color: var(--el-fill-color);
}

.markdown-preview :deep(mark) {
  background-color: var(--el-color-warning-light-9);
  color: var(--el-color-black);
  padding: 0 2px;
}

.markdown-preview :deep(abbr) {
  border-bottom: 1px dotted var(--el-text-color-primary);
  cursor: help;
}

.markdown-preview :deep(ins) {
  background-color: var(--el-color-success-light-9);
  color: var(--el-color-black);
  text-decoration: none;
  padding: 0 2px;
}

.markdown-preview :deep(sub) {
  font-size: 0.8em;
  vertical-align: sub;
}

.markdown-preview :deep(sup) {
  font-size: 0.8em;
  vertical-align: super;
}

.markdown-preview :deep(.footnotes) {
  font-size: 0.9em;
  margin-top: 20px;
  border-top: 1px solid var(--el-border-color);
  padding-top: 10px;
}

.markdown-preview :deep(.footnote-ref) {
  font-size: 0.8em;
  vertical-align: super;
  text-decoration: none;
}

.markdown-preview :deep(.footnote-backref) {
  font-size: 0.8em;
  text-decoration: none;
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