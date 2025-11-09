import axios from "axios";

/**
 * 删除LLM服务
 * @param {number} id - 服务ID
 * @returns {Promise<Object>}
 */
export const deleteLLMService = async (id) => {
    try {
        const token = localStorage.getItem("token");
        const serverLink = localStorage.getItem("server_link");
        const url = serverLink + `/llm/delete/${id}`;
        
        const response = await axios.delete(url, {
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        
        return response.data;
    } catch (error) {
        console.error('删除LLM服务出错:', error);
        return { code: 1, message: '网络错误，请稍后重试' };
    }
};
