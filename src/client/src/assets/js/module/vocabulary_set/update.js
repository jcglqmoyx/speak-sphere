import axios from "axios";

export async function updateVocabularySet(id, title, category, created_at) {
    try {
        const token = localStorage.getItem("token");
        const serverLink = localStorage.getItem("server_link");
        const url = serverLink + '/vocabulary_set/update'
        const response = await axios.put(url, {
            id, title, category, created_at,
        }, {
            headers: {
                Authorization: `Bearer ${token}`, 'Content-Type': 'application/json',
            },
        });
        return response.data
    } catch (error) {
        console.error('Error fetching data from backend:', error);
    }
}
