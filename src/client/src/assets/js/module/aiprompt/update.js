import axios from "axios";

/**
 * 更新AI提示词
 * @param {number} id - 提示词ID
 * @param {Object} updateData - 更新数据
 * @returns {Promise<Object>}
 */
export const updateAIPrompt = async (id, updateData) => {
    try {
        const token = localStorage.getItem("token");
        const serverLink = localStorage.getItem("server_link");
        const url = serverLink + `/aiprompt/update/${id}`;
        
        const response = await axios.put(url, updateData, {
            headers: {
                Authorization: `Bearer ${token}`,
                'Content-Type': 'application/json'
            }
        });
        
        return response.data;
    } catch (error) {
        console.error('更新AI提示词出错:', error);
        return { code: 1, message: '网络错误，请稍后重试' };
    }
};
