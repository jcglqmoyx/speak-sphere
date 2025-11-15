/**
 * 分析音频URL并推断国家/地区标签
 * @param {string} audioUrl - 音频URL
 * @returns {string} 国家/地区标签
 */
function getAudioRegionLabel(audioUrl) {
    if (!audioUrl) return '';

    const url = audioUrl.toLowerCase();
    let index = url.lastIndexOf('/');
    let identifier = url.substring(index);
    if (identifier.includes('uk')) return 'UK';
    if (identifier.includes('us')) return 'US';
    if (identifier.includes('au')) return 'AU';
    if (identifier.includes('ca')) return 'CA';

    return '';
}

/**
 * 调用字典API获取单词释义
 * @param {string} vocabulary - 要查询的单词
 * @returns {Promise<Object>} 返回API响应数据
 */
async function getVocabularyDefinition(vocabulary) {
    try {
        const response = await fetch(`https://api.dictionaryapi.dev/api/v2/entries/en/${encodeURIComponent(vocabulary)}`);

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();

        // 检查是否返回有效数据
        if (!data || !Array.isArray(data) || data.length === 0 || !data[0].vocabulary) {
            throw new Error('Empty or invalid data returned from dictionary API');
        }

        return {
            success: true, data: data, source: 'dictionaryapi'
        };
    } catch (error) {
        console.error('Error fetching vocabulary definition from dictionaryapi.dev:', error);

        // 返回wiktionary作为备用方案
        return {
            success: true, // 仍然返回success=true，因为wiktionary是有效的备选方案
            data: null,
            source: 'wiktionary',
            wiktionaryUrl: `https://en.wiktionary.org/wiki/${encodeURIComponent(vocabulary)}`
        };
    }
}

/**
 * 格式化字典API返回的释义数据
 * @param {Object} apiResponse - 字典API返回的响应数据，包含source字段
 * @returns {string} 格式化后的释义HTML字符串
 */
function formatDefinition(apiResponse) {
    // 如果是wiktionary源，返回iframe嵌入的HTML
    if (apiResponse.source === 'wiktionary') {
        const isMobile = window.innerWidth <= 768;

        if (isMobile) {
            // 移动端：显示链接而不是内嵌iframe
            return `
        <div class="wiktionary-mobile-notice">
          <div class="wiktionary-mobile-content">
            <div class="notice-title">使用 Wiktionary 查询</div>
            <div class="notice-message">当前使用 Wiktionary 作为备选词典源</div>
            <div class="action-buttons">
              <button class="open-wiktionary-btn" onclick="openWiktionaryFullscreen('${apiResponse.wiktionaryUrl}')">
                全屏打开 Wiktionary
              </button>
              <button class="open-new-tab-btn" onclick="window.open('${apiResponse.wiktionaryUrl}', '_blank')">
                新标签页打开
              </button>
            </div>
          </div>
        </div>
        <style>
          .wiktionary-mobile-notice {
            background: var(--el-bg-color);
            border: 1px solid var(--el-border-color);
            border-radius: 8px;
            padding: 20px;
            margin: 10px 0;
          }
          .wiktionary-mobile-content {
            display: flex;
            flex-direction: column;
            gap: 15px;
          }
          .notice-title {
            font-size: 18px;
            font-weight: bold;
            color: var(--el-text-color-primary);
          }
          .notice-message {
            color: var(--el-text-color-regular);
          }
          .action-buttons {
            display: flex;
            gap: 10px;
            flex-wrap: wrap;
          }
          .open-wiktionary-btn, .open-new-tab-btn {
            padding: 10px 15px;
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-size: 14px;
            transition: all 0.3s ease;
          }
          .open-wiktionary-btn {
            background: var(--el-color-primary);
            color: white;
          }
          .open-wiktionary-btn:hover {
            background: var(--el-color-primary-light-3);
          }
          .open-new-tab-btn {
            background: var(--el-fill-color-light);
            color: var(--el-text-color-primary);
            border: 1px solid var(--el-border-color);
          }
          .open-new-tab-btn:hover {
            background: var(--el-fill-color);
          }
        </style>
      `;
        } else {
            // 桌面端：保持原有的内嵌iframe
            return `
        <div class="wiktionary-container">
          <div class="wiktionary-notice">
            当前使用 Wiktionary 作为备选词典源
          </div>
          <div class="wiktionary-loading" id="wiktionary-loading-${Date.now()}">
            <div class="loading-spinner"></div>
            <span>正在加载 Wiktionary 页面...</span>
          </div>
          <iframe 
            src="${apiResponse.wiktionaryUrl}" 
            class="wiktionary-iframe"
            title="Wiktionary - ${apiResponse.wiktionaryUrl.split('/').pop()}"
            sandbox="allow-same-origin allow-scripts allow-popups allow-forms"
            referrerpolicy="no-referrer"
            onload="document.getElementById('wiktionary-loading-${Date.now()}').style.display='none'; this.style.display='block'"
            style="display:none"
          ></iframe>
        </div>
      `;
        }
    }

    // 原有的dictionaryapi.dev数据处理逻辑
    const apiData = apiResponse.data;
    if (!apiData || !Array.isArray(apiData) || apiData.length === 0) {
        return '<div class="no-meaning">No definition found</div>';
    }

    const vocabularyData = apiData[0];
    let formattedHtml = '';

    // 单词基本信息
    if (vocabularyData.vocabulary) {
        formattedHtml += `<div class="vocabulary-title">${vocabularyData.vocabulary}</div>`;
    }

    // 音标 - 优先显示有文本音标的内容
    if (vocabularyData.phonetics && Array.isArray(vocabularyData.phonetics)) {
        const phonetic = vocabularyData.phonetics.find(p => p.text)?.text;
        if (phonetic) {
            formattedHtml += `<div class="phonetic">Phonetic: ${phonetic}`;

            // 添加发音按钮
            const validAudioVocabularies = vocabularyData.phonetics.filter(p => p.audio && p.audio.trim() !== '');
            if (validAudioVocabularies.length > 0) {
                formattedHtml += ' <span class="audio-buttons">';

                validAudioVocabularies.forEach((phoneticItem, index) => {
                    const audioUrl = phoneticItem.audio;
                    let regionLabel = getAudioRegionLabel(audioUrl);

                    // 如果无法推断出地区标签，使用序号
                    if (!regionLabel) {
                        regionLabel = `Audio ${index + 1}`;
                    }

                    formattedHtml += `<button class="audio-btn" data-audio="${audioUrl}" title="Play ${regionLabel} pronunciation">🔊${regionLabel}</button>`;
                });

                formattedHtml += '</span>';
            }

            formattedHtml += '</div>';
        }
    }

    // 遍历所有词性释义
    if (vocabularyData.meanings && Array.isArray(vocabularyData.meanings)) {
        formattedHtml += '<div class="meanings-container">';
        vocabularyData.meanings.forEach((meaning, index) => {
            if (meaning.partOfSpeech) {
                formattedHtml += `<div class="part-of-speech">【${meaning.partOfSpeech}】</div>`;
            }

            if (meaning.definitions && Array.isArray(meaning.definitions)) {
                formattedHtml += '<div class="definitions-list">';
                meaning.definitions.forEach((def, defIndex) => {
                    if (def.definition) {
                        formattedHtml += `<div class="definition-item">
              <span class="definition-number">${defIndex + 1}.</span>
              <span class="definition-text">${def.definition}</span>`;

                        // 如果有例句
                        if (def.example) {
                            formattedHtml += `<div class="example">Example: <em>${def.example}</em></div>`;
                        }

                        // 如果有同义词
                        if (def.synonyms && def.synonyms.length > 0) {
                            formattedHtml += `<div class="synonyms">Synonyms: ${def.synonyms.join(', ')}</div>`;
                        }

                        formattedHtml += '</div>';
                    }
                });
                formattedHtml += '</div>';
            }

            // 如果有整体同义词
            if (meaning.synonyms && meaning.synonyms.length > 0) {
                formattedHtml += `<div class="meaning-synonyms">[Synonyms] ${meaning.synonyms.join(', ')}</div>`;
            }

            // 如果有整体反义词
            if (meaning.antonyms && meaning.antonyms.length > 0) {
                formattedHtml += `<div class="meaning-antonyms">[Antonyms] ${meaning.antonyms.join(', ')}</div>`;
            }

            if (index < vocabularyData.meanings.length - 1) {
                formattedHtml += '<div class="meaning-separator"></div>';
            }
        });
        formattedHtml += '</div>';
    }

    return formattedHtml || '<div class="no-meaning">No definition found</div>';
}

export {getVocabularyDefinition, formatDefinition};
