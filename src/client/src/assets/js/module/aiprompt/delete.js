import axios from "axios";

/**
 * 删除AI提示词
 * @param {number} id - 提示词ID
 * @returns {Promise<Object>}
 */
export const deleteAIPrompt = async (id) => {
    try {
        const token = localStorage.getItem("token");
        const serverLink = localStorage.getItem("server_link");
        const url = serverLink + `/aiprompt/delete/${id}`;
        
        const response = await axios.delete(url, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        
        return response.data;
    } catch (error) {
        console.error('删除AI提示词出错:', error);
        return { code: 1, message: '网络错误，请稍后重试' };
    }
};
