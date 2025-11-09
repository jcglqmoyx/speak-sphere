import axios from "axios";

/**
 * 添加LLM服务
 * @param {Object} llmService - LLM服务对象
 * @returns {Promise<Object>}
 */
export const addLLMService = async (llmService) => {
    try {
        const token = localStorage.getItem("token");
        const serverLink = localStorage.getItem("server_link");
        const url = serverLink + '/llm/add';
        
        // 确保字段名与后端一致
        const requestData = {
            ...llmService,
            endpoint: llmService.endpoint,
            api_key: llmService.api_key
        };
        
        const response = await axios.post(url, requestData, {
            headers: {
                Authorization: `Bearer ${token}`,
                'Content-Type': 'application/json',
            }
        });
        
        return response.data;
    } catch (error) {
        console.error('添加LLM服务出错:', error);
        return { code: 1, message: '网络错误，请稍后重试' };
    }
};
