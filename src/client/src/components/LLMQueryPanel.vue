<template>
  <!-- 只保留侧边面板，按钮移到WordCardView -->
  <div v-if="drawerVisible" class="custom-drawer-overlay" @click="closePanel">
    <div class="custom-drawer" @click.stop>
      <div class="drawer-header">
        <h3>LLM 单词查询</h3>
        <el-button type="text" @click="closePanel" class="close-btn">
          <el-icon><Close /></el-icon>
        </el-button>
      </div>
      
      <div class="drawer-content">
        <div v-if="configStatus !== 'ready'" class="config-status">
          <el-alert v-if="configStatus === 'loading'" title="正在加载配置..." type="info" show-icon />
          <el-alert v-if="configStatus === 'no_word'" title="请先选择要查询的单词" type="warning" show-icon />
          <el-alert v-if="configStatus === 'no_prompt'" title="请先设置默认AI提示词" type="warning" show-icon />
          <el-alert v-if="configStatus === 'no_llm'" title="请先设置默认LLM服务" type="warning" show-icon />
        </div>

        <div v-if="configStatus === 'ready'" class="llm-content">
          <!-- 提示词区域 -->
          <div class="prompt-section">
            <div class="section-header">
              <h4>AI 提示词</h4>
              <div class="prompt-actions">
                <el-button class="test-button-copy" size="small" @click="copyPrompt" style="background-color: #409eff; color: white; border: 2px solid #409eff; padding: 8px 16px; border-radius: 4px;">复制</el-button>
                <el-button class="test-button-reset" size="small" @click="resetPrompt" style="background-color: #909399; color: white; border: 2px solid #909399; padding: 8px 16px; border-radius: 4px;">重置</el-button>
              </div>
            </div>
            <textarea
              v-model="editablePrompt"
              :rows="6"
              placeholder="正在加载提示词..."
              class="prompt-textarea-native"
              @input="onPromptChange"
            />
            <div class="prompt-tips">
              <small>提示：点击文本框可以编辑提示词内容，使用 <code>#word#</code> 作为当前单词的占位符</small>
            </div>
          </div>

          <!-- 查询按钮 -->
          <div class="query-section">
            <el-button 
              class="test-button-query"
              @click="queryLLM" 
              :loading="isQuerying"
              :disabled="isQuerying"
              style="background-color: #67c23a; color: white; border: 2px solid #67c23a; padding: 10px 20px; border-radius: 4px; min-width: 100px;"
            >
              {{ isQuerying ? '查询中...' : '开始查询' }}
            </el-button>
            <el-button v-if="hasResponse" @click="clearResponse" style="background-color: #f56c6c; color: white; border: 2px solid #f56c6c; padding: 10px 20px; border-radius: 4px;">清除回复</el-button>
          </div>

          <!-- 回复区域 -->
          <div v-if="hasResponse || isQuerying" class="response-section">
            <h4>LLM 回复</h4>
            <div ref="responseContentRef" class="response-content" v-html="renderedResponse" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { Close } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { queryDefaultAIPrompt } from '@/assets/js/module/aiprompt/query'
import { queryDefaultLLMService } from '@/assets/js/module/llm/query'
import MarkdownIt from 'markdown-it'
import markdownItAbbr from 'markdown-it-abbr'
import markdownItFootnote from 'markdown-it-footnote'
import markdownItIns from 'markdown-it-ins'
import markdownItMark from 'markdown-it-mark'
import markdownItSub from 'markdown-it-sub'
import markdownItSup from 'markdown-it-sup'
import markdownItTexmath from 'markdown-it-texmath'

export default {
  name: 'LLMQueryPanel',
  components: {
    Close
  },
  props: {
    currentWord: {
      type: String,
      default: ''
    }
  },
  data() {
    const md = new MarkdownIt({
      html: true,
      linkify: true,
      typographer: true
    })
      .use(markdownItAbbr)
      .use(markdownItFootnote)
      .use(markdownItIns)
      .use(markdownItMark)
      .use(markdownItSub)
      .use(markdownItSup)
      .use(markdownItTexmath)
    
    return {
      drawerVisible: true,
      configStatus: 'loading',
      aiPrompt: null,
      llmConfig: null,
      isQuerying: false,
      responseText: '',
      controller: null,
      md,
      editablePrompt: '',
      originalPrompt: ''
    }
  },
  computed: {
    hasResponse() {
      return this.responseText.length > 0
    },
    renderedResponse() {
      if (!this.md || !this.responseText) return ''
      return this.md.render(this.responseText)
    }
  },
  methods: {
    closePanel() {
      this.drawerVisible = false
      this.$emit('close')
    },
    
    onPromptChange(value) {
      // 提示词改变时的处理
      console.log('提示词已更新:', value)
    },
    
    copyPrompt() {
      try {
        navigator.clipboard.writeText(this.editablePrompt)
        ElMessage.success('提示词已复制到剪贴板')
      } catch (error) {
        console.error('复制失败:', error)
        ElMessage.error('复制失败')
      }
    },
    
    resetPrompt() {
      if (this.originalPrompt) {
        // 重置时重新替换#word#占位符
        this.editablePrompt = this.originalPrompt.replace(/#word#/g, this.currentWord)
        ElMessage.success('提示词已重置为默认')
      }
    },

    async loadConfigurations() {
      if (!this.currentWord) {
        this.configStatus = 'no_word'
        return
      }

      this.configStatus = 'loading'

      try {
        const promptResponse = await queryDefaultAIPrompt()
        if (promptResponse.code !== 0 || !promptResponse.data) {
          this.configStatus = 'no_prompt'
          return
        }
        this.aiPrompt = promptResponse.data

        // 设置原始模板和可编辑提示词（已经替换#word#）
        this.originalPrompt = this.aiPrompt.content || ''
        this.editablePrompt = this.originalPrompt.replace(/#word#/g, this.currentWord)

        console.log('提示词处理:', {
          original: this.originalPrompt,
          currentWord: this.currentWord,
          editable: this.editablePrompt
        })

        const llmResponse = await queryDefaultLLMService()
        if (llmResponse.code !== 0 || !llmResponse.data) {
          this.configStatus = 'no_llm'
          return
        }
        this.llmConfig = llmResponse.data

        this.configStatus = 'ready'
        console.log('配置加载完成，editablePrompt:', this.editablePrompt)
      } catch (error) {
        console.error('加载配置失败:', error)
        this.configStatus = 'no_prompt'
      }
    },

    async queryLLM() {
      if (!this.aiPrompt || !this.llmConfig) {
        ElMessage.error('配置信息不完整')
        return
      }

      if (!this.currentWord) {
        ElMessage.error('请先选择要查询的单词')
        return
      }

      this.isQuerying = true
      this.responseText = ''
      
      try {
        this.controller = new AbortController()
        
        // 确保endpoint格式正确
        let baseUrl = this.llmConfig.endpoint
        if (!baseUrl) {
          throw new Error('LLM服务配置中的endpoint为空')
        }
        
        // 如果endpoint不是以http开头，添加https://前缀
        if (!baseUrl.startsWith('http://') && !baseUrl.startsWith('https://')) {
          baseUrl = 'https://' + baseUrl
        }
        
        // 确保endpoint包含完整的API路径
        let apiUrl = baseUrl
        if (!apiUrl.includes('/chat/completions')) {
          if (!apiUrl.endsWith('/v1')) {
            apiUrl = apiUrl.replace(/\/+$/, '') + '/v1'
          }
          apiUrl += '/chat/completions'
        }
        
        const response = await fetch(apiUrl, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${this.llmConfig.api_key}`
          },
          body: JSON.stringify({
            model: this.llmConfig.model,
            messages: [
              {
                role: 'user',
                content: this.editablePrompt
              }
            ],
            stream: true,
            max_tokens: 1000,
            temperature: 0.7
          }),
          signal: this.controller.signal
        })

        if (!response.ok) {
          const errorText = await response.text()
          throw new Error(`HTTP ${response.status}: ${errorText}`)
        }

        const reader = response.body.getReader()
        const decoder = new TextDecoder()
        
        let reading = true
        while (reading) {
          const { done, value } = await reader.read()
          if (done) {
            reading = false
            break
          }

          const chunk = decoder.decode(value, { stream: true })
          const lines = chunk.split('\n').filter(line => line.trim() !== '')

          for (const line of lines) {
            if (line.startsWith('data: ')) {
              const data = line.slice(6)
              if (data === '[DONE]') {
                reading = false
                break
              }

              try {
                const parsed = JSON.parse(data)
                const content = parsed.choices[0]?.delta?.content
                if (content) {
                  this.responseText += content
                  // 滚动到底部
                  this.$nextTick(() => {
                    const el = this.$refs.responseContentRef
                    if (el) {
                      el.scrollTop = el.scrollHeight
                    }
                  })
                }
              } catch (e) {
                console.warn('解析流数据失败:', e, '数据:', data)
              }
            }
          }
        }
      } catch (error) {
        if (error.name === 'AbortError') {
          ElMessage.info('查询已取消')
        } else {
          console.error('查询失败:', error)
          
          // 提供更友好的错误信息
          let errorMessage = '查询失败'
          if (error.message.includes('Failed to fetch')) {
            errorMessage = '网络连接失败，请检查LLM服务地址是否正确'
          } else if (error.message.includes('401')) {
            errorMessage = 'API密钥无效，请检查LLM服务配置'
          } else if (error.message.includes('404')) {
            errorMessage = 'API端点未找到，请检查LLM服务地址'
          } else {
            errorMessage = `查询失败: ${error.message}`
          }
          
          ElMessage.error(errorMessage)
        }
      } finally {
        this.isQuerying = false
        this.controller = null
      }
    },
    
    clearResponse() {
      this.responseText = ''
    }
  },
  watch: {
    currentWord: {
      immediate: true,
      handler(newWord) {
        console.log('当前单词变化:', newWord)
        if (newWord) {
          this.loadConfigurations()
        }
      }
    }
  }
}
</script>

<style scoped>
/* 自定义侧边面板样式 */
.custom-drawer-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 2000;
}

.custom-drawer {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  width: 50%;
  min-width: 500px;
  background: var(--el-bg-color);
  color: var(--el-text-color-primary);
  box-shadow: -2px 0 8px rgba(0, 0, 0, 0.15);
  display: flex;
  flex-direction: column;
}

.drawer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid var(--el-border-color);
  background-color: var(--el-fill-color-light);
}

.drawer-header h3 {
  margin: 0;
  color: var(--el-text-color-primary);
}

.close-btn {
  font-size: 16px;
  color: var(--el-text-color-secondary);
}

.close-btn:hover {
  color: var(--el-color-primary);
}

.drawer-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  background-color: var(--el-bg-color);
  color: var(--el-text-color-primary);
}

.config-status {
  margin-bottom: 20px;
}

.llm-content {
  display: flex;
  flex-direction: column;
  height: 100%;
  gap: 20px;
}

.prompt-section,
.query-section,
.response-section {
  width: 100%;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.section-header h4 {
  margin: 0;
  color: var(--el-text-color-primary);
}

.prompt-actions {
  display: flex;
  gap: 8px;
}

.prompt-actions .el-button {
  border: 1px solid var(--el-border-color);
  background-color: var(--el-bg-color);
  color: var(--el-text-color-primary);
  font-weight: normal;
  padding: 8px 16px;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.prompt-actions .el-button:hover {
  background-color: var(--el-color-primary);
  color: white;
  border-color: var(--el-color-primary);
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.prompt-actions .el-button:active {
  transform: translateY(0);
  box-shadow: none;
}

.prompt-textarea-native {
  width: 100%;
  padding: 12px;
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.5;
  background-color: var(--el-bg-color);
  color: var(--el-text-color-primary);
  resize: vertical;
  min-height: 120px;
}

.prompt-textarea-native:focus {
  outline: none;
  border-color: var(--el-color-primary);
  box-shadow: 0 0 0 2px var(--el-color-primary-light-5);
}

.prompt-textarea-native::placeholder {
  color: var(--el-text-color-placeholder);
}

.prompt-tips {
  margin-top: 8px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.prompt-tips code {
  background-color: var(--el-fill-color);
  color: var(--el-color-primary);
  padding: 2px 4px;
  border-radius: 3px;
  font-family: 'Courier New', monospace;
}

.query-section {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.query-section .el-button {
  border: 1px solid var(--el-border-color);
  background-color: var(--el-bg-color);
  color: var(--el-text-color-primary);
  font-weight: normal;
  padding: 10px 20px;
  border-radius: 4px;
  transition: all 0.3s ease;
  min-width: 100px;
}

.query-section .el-button:hover {
  background-color: var(--el-color-primary);
  color: white;
  border-color: var(--el-color-primary);
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.query-section .el-button:active {
  transform: translateY(0);
  box-shadow: none;
}

/* 特殊处理查询按钮的样式 */
.query-section .el-button--success {
  background-color: var(--el-color-success);
  color: white;
  border-color: var(--el-color-success);
}

.query-section .el-button--success:hover {
  background-color: var(--el-color-success-light-3);
  border-color: var(--el-color-success-light-3);
}

.response-content {
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  padding: 12px;
  min-height: 300px;
  max-height: 400px;
  overflow-y: auto;
  background-color: var(--el-fill-color-light);
  line-height: 1.6;
  color: var(--el-text-color-primary);
}

/* Markdown内容样式 - 深色主题兼容 */
.response-content :deep(h1),
.response-content :deep(h2),
.response-content :deep(h3),
.response-content :deep(h4) {
  margin: 1.2em 0 0.6em 0;
  color: var(--el-text-color-primary);
  line-height: 1.25;
}

.response-content :deep(p) {
  margin: 0.8em 0;
  color: var(--el-text-color-primary);
}

.response-content :deep(blockquote) {
  margin: 1em 0;
  padding: 0.5em 1em;
  border-left: 4px solid var(--el-color-primary);
  background-color: var(--el-fill-color-lighter);
  color: var(--el-text-color-regular);
  font-style: italic;
}

.response-content :deep(ul),
.response-content :deep(ol) {
  margin: 0.8em 0;
  padding-left: 2em;
  color: var(--el-text-color-primary);
}

.response-content :deep(li) {
  margin: 0.4em 0;
  color: var(--el-text-color-primary);
}

.response-content :deep(code) {
  background-color: var(--el-fill-color);
  color: var(--el-color-danger);
  padding: 2px 6px;
  border-radius: 3px;
  font-family: 'Monaco', 'Menlo', 'Consolas', 'Courier New', monospace;
  font-size: 0.9em;
}

.response-content :deep(pre) {
  background-color: var(--el-fill-color-light);
  border: 1px solid var(--el-border-color);
  border-radius: 6px;
  padding: 1em;
  overflow-x: auto;
  margin: 1em 0;
}

.response-content :deep(pre code) {
  background: none;
  color: var(--el-text-color-primary);
  padding: 0;
  border-radius: 0;
}

.response-content :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 1em 0;
}

.response-content :deep(th),
.response-content :deep(td) {
  border: 1px solid var(--el-border-color);
  padding: 0.6em 1em;
  text-align: left;
  color: var(--el-text-color-primary);
}

.response-content :deep(th) {
  background-color: var(--el-fill-color-light);
  font-weight: 600;
}

.response-content :deep(a) {
  color: var(--el-color-primary);
  text-decoration: none;
}

.response-content :deep(a:hover) {
  text-decoration: underline;
  color: var(--el-color-primary-light-3);
}

.response-content :deep(strong) {
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.response-content :deep(em) {
  font-style: italic;
  color: var(--el-text-color-regular);
}

/* 深色主题下的特殊调整 */
:deep(.el-alert) {
  background-color: var(--el-bg-color-overlay) !important;
  color: var(--el-text-color-primary) !important;
}

:deep(.el-alert__title) {
  color: var(--el-text-color-primary) !important;
}

:deep(.el-textarea__inner) {
  background-color: var(--el-bg-color) !important;
  color: var(--el-text-color-primary) !important;
  border-color: var(--el-border-color) !important;
}
</style>
