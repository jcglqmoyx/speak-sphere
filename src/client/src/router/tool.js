import VocabularyQueryView from "@/views/vocabulary_query/VocabularyQueryView.vue";
import AudioRecorderView from "@/views/audio-recorder/AudioRecorderView.vue";

const ToolRoutes = [
    {
        path: '/tool',
        name: 'tool',
        children: [
            {
                path: "query",
                name: 'query',
                component: VocabularyQueryView,
            },
            {
                path: 'record',
                name: 'record',
                component: AudioRecorderView,
            }
        ]
    }
]

export default ToolRoutes;