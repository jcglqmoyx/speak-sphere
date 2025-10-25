import axios from "axios";

export async function updateDictionary(id, title, prefix, suffix) {
    try {
        const token = localStorage.getItem("token");
        const serverLink = localStorage.getItem("server_link");
        const url = serverLink + '/dictionary/update'
        const response = await axios.put(url, {
            id,
            title,
            prefix,
            suffix,
        }, {
            headers: {
                Authorization: `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
        return response.data
    } catch (error) {
        console.error('Error fetching data from backend:', error);
    }
}
