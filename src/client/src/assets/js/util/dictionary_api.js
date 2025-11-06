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
 * @param {string} word - 要查询的单词
 * @returns {Promise<Object>} 返回API响应数据
 */
async function getWordDefinition(word) {
    try {
        const response = await fetch(`https://api.dictionaryapi.dev/api/v2/entries/en/${encodeURIComponent(word)}`);

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const data = await response.json();
        return {
            success: true,
            data: data
        };
    } catch (error) {
        console.error('Error fetching word definition:', error);
        return {
            success: false,
            error: error.message || 'Failed to fetch word definition'
        };
    }
}

/**
 * 格式化字典API返回的释义数据
 * @param {Object} apiData - 字典API返回的数据
 * @returns {string} 格式化后的释义HTML字符串
 */
function formatDefinition(apiData) {
    if (!apiData || !Array.isArray(apiData) || apiData.length === 0) {
        return '<div class="no-meaning">No definition found</div>';
    }

    const wordData = apiData[0];
    let formattedHtml = '';

    // 单词基本信息
    if (wordData.word) {
        formattedHtml += `<div class="word-title">${wordData.word}</div>`;
    }

    // 音标 - 优先显示有文本音标的内容
    if (wordData.phonetics && Array.isArray(wordData.phonetics)) {
        const phonetic = wordData.phonetics.find(p => p.text)?.text;
        if (phonetic) {
            formattedHtml += `<div class="phonetic">Phonetic: ${phonetic}`;

            // 添加发音按钮
            const validAudioEntries = wordData.phonetics.filter(p => p.audio && p.audio.trim() !== '');
            if (validAudioEntries.length > 0) {
                formattedHtml += ' <span class="audio-buttons">';

                validAudioEntries.forEach((phoneticItem, index) => {
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
    if (wordData.meanings && Array.isArray(wordData.meanings)) {
        formattedHtml += '<div class="meanings-container">';
        wordData.meanings.forEach((meaning, index) => {
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

            if (index < wordData.meanings.length - 1) {
                formattedHtml += '<div class="meaning-separator"></div>';
            }
        });
        formattedHtml += '</div>';
    }

    return formattedHtml || '<div class="no-meaning">No definition found</div>';
}

export { getWordDefinition, formatDefinition };
