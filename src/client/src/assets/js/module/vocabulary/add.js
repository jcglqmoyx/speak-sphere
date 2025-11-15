import axios from "axios";

export async function AddVocabulary(vocabularySetId, vocabulary, meaning, note) {
    try {
        const token = localStorage.getItem("token");
        const serverLink = localStorage.getItem("server_link");
        const url = serverLink + '/vocabulary/add';
        const response = await axios.post(url, {
                vocabulary_set_id: vocabularySetId,
                vocabulary: vocabulary,
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
