import axios from "axios";

/**
 * 查询AI提示词列表
 * @returns {Promise<Object>}
 */
export const queryAIPrompts = async () => {
    try {
        const token = localStorage.getItem("token");
        const serverLink = localStorage.getItem("server_link");
        const url = serverLink + '/aiprompt/list';
        
        const response = await axios.get(url, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        
        return response.data;
    } catch (error) {
        console.error('查询AI提示词出错:', error);
        return { code: 1, message: '网络错误，请稍后重试' };
    }
};

/**
 * 查询默认AI提示词
 * @returns {Promise<Object>}
 */
export const queryDefaultAIPrompt = async () => {
    try {
        const token = localStorage.getItem("token");
        const serverLink = localStorage.getItem("server_link");
        const url = serverLink + '/aiprompt/default';
        
        const response = await axios.get(url, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        
        return response.data;
    } catch (error) {
        console.error('查询默认AI提示词出错:', error);
        return { code: 1, message: '网络错误，请稍后重试' };
    }
};
