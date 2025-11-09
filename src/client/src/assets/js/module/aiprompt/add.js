import axios from "axios";

/**
 * 添加AI提示词
 * @param {Object} aiPrompt - AI提示词对象
 * @returns {Promise<Object>}
 */
export const addAIPrompt = async (aiPrompt) => {
    try {
        const token = localStorage.getItem("token");
        const serverLink = localStorage.getItem("server_link");
        const url = serverLink + '/aiprompt/add';
        
        const response = await axios.post(url, aiPrompt, {
            headers: {
                Authorization: `Bearer ${token}`,
                'Content-Type': 'application/json',
            }
        });
        
        return response.data;
    } catch (error) {
        console.error('添加AI提示词出错:', error);
        return { code: 1, message: '网络错误，请稍后重试' };
    }
};
