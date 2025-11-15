import axios from "axios";

async function getVocabularySetCount() {
    try {
        const serverLink = localStorage.getItem("server_link");
        const token = localStorage.getItem("token");
        const response = await axios.get(serverLink + '/vocabulary_set/count', {
            headers: {
                Authorization: `Bearer ${token}`, 'Content-Type': 'application/json',
            },
        });
        return response.data
    } catch (error) {
        console.error('Error fetching data from backend:', error);
    }
}

async function getVocabularySetList(pageSize, currentPage) {
    try {
        const serverLink = localStorage.getItem("server_link");
        const token = localStorage.getItem("token");
        const url = serverLink + '/vocabulary_set/list/' + pageSize + '/' + currentPage;
        const response = await axios.get(url, {
            headers: {
                Authorization: `Bearer ${token}`, 'Content-Type': 'application/json',
            },
        });
        return response.data
    } catch (error) {
        console.error('Error fetching data from backend:', error);
    }
}

export {getVocabularySetCount, getVocabularySetList}