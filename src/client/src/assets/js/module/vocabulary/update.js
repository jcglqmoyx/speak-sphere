import axios from "axios";

async function updateVocabulary(id, vocabulary, meaning, vocabulary_set_id, note, unwanted, study_count, date_to_review, created_at) {
    try {
        const token = localStorage.getItem("token");
        const serverLink = localStorage.getItem("server_link");
        const url = serverLink + '/vocabulary/update'
        console.log(id,
            vocabulary,
            meaning,
            vocabulary_set_id,
            note,
            date_to_review,
            created_at,)
        ;
        const response = await axios.put(url, {
            id,
            vocabulary: vocabulary,
            meaning,
            vocabulary_set_id,
            note,
            unwanted,
            study_count,
            date_to_review,
            created_at,
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

async function setUnwanted(id) {
    try {
        const token = localStorage.getItem("token");
        const serverLink = localStorage.getItem("server_link");
        const url = serverLink + '/vocabulary/update/unwanted/' + id;
        const response = await axios.put(url, null, {
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

async function updateVocabularyStudyCount(id) {
    try {
        const token = localStorage.getItem("token");
        const serverLink = localStorage.getItem("server_link");
        const url = serverLink + '/vocabulary/update/study/count/' + id;
        const response = await axios.put(url, null,
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

async function resetVocabularyStudyCountToZero(id) {
    try {
        const token = localStorage.getItem("token");
        const serverLink = localStorage.getItem("server_link");
        const url = serverLink + '/vocabulary/update/reset/' + id;
        const response = await axios.put(url, null,
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

export {updateVocabulary, setUnwanted, updateVocabularyStudyCount, resetVocabularyStudyCountToZero};