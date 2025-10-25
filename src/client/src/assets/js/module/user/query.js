import axios from "axios";

async function getUserProfile() {
    try {
        const serverLink = localStorage.getItem("server_link");
        const token = localStorage.getItem("token");
        const url = serverLink + '/user/profile';
        const response = await axios.get(url, {
            headers: {
                Authorization: `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
        return response.data.data
    } catch (error) {
        console.error('Error fetching data from backend:', error);
    }
}

export {getUserProfile};