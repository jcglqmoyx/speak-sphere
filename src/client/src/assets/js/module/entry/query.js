import axios from "axios";

async function fetchWords(path) {
    try {
        const serverLink = localStorage.getItem("server_link");
        const token = localStorage.getItem("token");
        const response = await axios.get(serverLink + '/entry/' + path, {
            headers: {
                Authorization: `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
        return response.data;
    } catch (error) {
        console.error('Error fetching data from backend:', error);
    }
}

async function getWordCount(bookId) {
    try {
        const serverLink = localStorage.getItem("server_link");
        const token = localStorage.getItem("token");
        const response = await axios.get(serverLink + '/entry/count/' + bookId, {
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

async function getWordList(bookId, pageSize, currentPage) {
    try {
        const serverLink = localStorage.getItem("server_link");
        const token = localStorage.getItem("token");
        const response = await axios.get(
            serverLink + '/entry/list?book_id=' + bookId + '&pageSize=' + pageSize + '&currentPage=' + currentPage,
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

async function checkWordInBook(word, bookId) {
    try {
        const serverLink = localStorage.getItem("server_link");
        const token = localStorage.getItem("token");
        const response = await axios.get(serverLink + '/entry/check?word=' + encodeURIComponent(word) + '&book_id=' + bookId, {
            headers: {
                Authorization: `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
        });
        return response.data
    } catch (error) {
        console.error('Error fetching data from backend:', error);
        return { code: 1, message: '查询失败' };
    }
}

export {fetchWords, getWordCount, getWordList, checkWordInBook};