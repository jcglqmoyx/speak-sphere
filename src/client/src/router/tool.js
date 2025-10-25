import QueryView from "@/views/query/QueryView.vue";
import AudioRecorderView from "@/views/audio-recorder/AudioRecorderView.vue";

const ToolRoutes = [
    {
        path: '/tool',
        name: 'tool',
        children: [
            {
                path: "query",
                name: 'query',
                component: QueryView,
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