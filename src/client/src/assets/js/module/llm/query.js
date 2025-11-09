import axios from "axios";

/**
 * 查询LLM服务列表
 * @returns {Promise<Object>}
 */
export const queryLLMServices = async () => {
    try {
        const token = localStorage.getItem("token");
        const serverLink = localStorage.getItem("server_link");
        const url = serverLink + '/llm/list';
        
        const response = await axios.get(url, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        
        return response.data;
    } catch (error) {
        console.error('查询LLM服务出错:', error);
        return { code: 1, message: '网络错误，请稍后重试' };
    }
};

/**
 * 查询默认LLM服务
 * @returns {Promise<Object>}
 */
export const queryDefaultLLMService = async () => {
    try {
        const token = localStorage.getItem("token");
        const serverLink = localStorage.getItem("server_link");
        const url = serverLink + '/llm/default';
        
        const response = await axios.get(url, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        
        return response.data;
    } catch (error) {
        console.error('查询默认LLM服务出错:', error);
        return { code: 1, message: '网络错误，请稍后重试' };
    }
};
