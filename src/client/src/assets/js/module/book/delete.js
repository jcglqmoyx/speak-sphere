import axios from "axios";

export async function deleteBook(id) {
    try {
        const serverLink = localStorage.getItem("server_link");
        const token = localStorage.getItem("token");
        const url = serverLink + '/book/delete/' + id;
        const response = await axios.delete(url,
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
