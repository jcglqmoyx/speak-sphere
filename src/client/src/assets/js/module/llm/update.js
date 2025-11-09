import axios from "axios";

/**
 * 更新LLM服务
 * @param {number} id - 服务ID
 * @param {Object} updateData - 更新数据
 * @returns {Promise<Object>}
 */
export const updateLLMService = async (id, updateData) => {
    try {
        const token = localStorage.getItem("token");
        const serverLink = localStorage.getItem("server_link");
        const url = serverLink + `/llm/update/${id}`;
        
        const response = await axios.put(url, updateData, {
            headers: {
                Authorization: `Bearer ${token}`,
                'Content-Type': 'application/json'
            }
        });
        
        return response.data;
    } catch (error) {
        console.error('更新LLM服务出错:', error);
        return { code: 1, message: '网络错误，请稍后重试' };
    }
};
