import axios from "axios";

async function fetchVocabularies(path) {
    try {
        const serverLink = localStorage.getItem("server_link");
        const token = localStorage.getItem("token");
        const response = await axios.get(serverLink + '/vocabulary/' + path, {
            headers: {
                Authorization: `Bearer ${token}`, 'Content-Type': 'application/json',
            },
        });
        return response.data;
    } catch (error) {
        console.error('Error fetching data from backend:', error);
    }
}

async function getVocabularyCount(vocabularySetId) {
    try {
        const serverLink = localStorage.getItem("server_link");
        const token = localStorage.getItem("token");
        const response = await axios.get(serverLink + '/vocabulary/count/' + vocabularySetId, {
            headers: {
                Authorization: `Bearer ${token}`, 'Content-Type': 'application/json',
            },
        });
        console.log('response.data: ', response.data)
        return response.data
    } catch (error) {
        console.error('Error fetching data from backend:', error);
    }
}

async function getVocabularyList(vocabularySetId, pageSize, currentPage) {
    try {
        const serverLink = localStorage.getItem("server_link");
        const token = localStorage.getItem("token");
        const response = await axios.get(serverLink + '/vocabulary/list?vocabulary_set_id=' + vocabularySetId + '&pageSize=' + pageSize + '&currentPage=' + currentPage, {
            headers: {
                Authorization: `Bearer ${token}`, 'Content-Type': 'application/json',
            },
        });
        return response.data
    } catch (error) {
        console.error('Error fetching data from backend:', error);
    }
}

async function checkVocabularyInVocabularySet(vocabulary, vocabularySetId) {
    try {
        const serverLink = localStorage.getItem("server_link");
        const token = localStorage.getItem("token");
        const response = await axios.get(serverLink + '/vocabulary/check?vocabulary=' + encodeURIComponent(vocabulary) + '&vocabulary_set_id=' + vocabularySetId, {
            headers: {
                Authorization: `Bearer ${token}`, 'Content-Type': 'application/json',
            },
        });
        return response.data
    } catch (error) {
        console.error('Error fetching data from backend:', error);
        return {code: 1, message: '查询失败'};
    }
}

export {fetchVocabularies, getVocabularyCount, getVocabularyList, checkVocabularyInVocabularySet};