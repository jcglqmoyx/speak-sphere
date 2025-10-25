import axios from "axios";

export async function AddEntry(bookId, word, meaning, note) {
    try {
        const token = localStorage.getItem("token");
        const serverLink = localStorage.getItem("server_link");
        const url = serverLink + '/entry/add';
        const response = await axios.post(url, {
                book_id: bookId,
                word,
                meaning,
                note,
            },
            {
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
