import axios from "axios";

export async function updateUser(username, email, avatar, current_book_id, daily_count, times_counted_as_known, review_frequency_formula, llm_service_provider, llm_token, llm_model) {
    try {
        const token = localStorage.getItem("token");
        const serverLink = localStorage.getItem("server_link");
        const url = serverLink + '/user/update'
        const response = await axios.put(url, {
            username,
            email,
            avatar,
            current_book_id,
            daily_count,
            times_counted_as_known,
            review_frequency_formula,
            llm_service_provider,
            llm_token,
            llm_model,
        }, {
            headers: {
                Authorization: `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });

        console.log('更新用户响应:', response.data);
        return response.data
    } catch (error) {
        console.error('更新用户时发生错误:', error);
        if (error.response) {
            console.error('错误响应数据:', error.response.data);
            console.error('错误状态码:', error.response.status);
        }
        return error.response ? error.response.data : {code: 1, message: '网络错误'};
    }
}
