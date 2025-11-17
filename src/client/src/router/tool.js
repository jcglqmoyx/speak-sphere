import WordQueryView from "@/views/word_query/WordQueryView.vue";
import AudioRecorderView from "@/views/audio-recorder/AudioRecorderView.vue";

const ToolRoutes = [
    {
        path: '/tool',
        name: 'tool',
        children: [
            {
                path: "query",
                name: 'query',
                component: WordQueryView,
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