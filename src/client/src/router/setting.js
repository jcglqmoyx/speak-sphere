import VocabularySetSettingView from "@/views/vocabulary_set/VocabularySetSettingView.vue";
import VocabularySettingView from "@/views/vocabulary/VocabularySettingView.vue";
import DictionarySettingView from "@/views/dictionary/DictionarySettingView.vue";
import LLMServiceSettingView from "@/views/llm/LLMServiceSettingView.vue";
import AIPromptSettingView from "@/views/aiprompt/AIPromptSettingView.vue";
import UserSettingView from "@/views/user/UserSettingView.vue";

const SettingRoutes = [{
    path: '/setting', name: 'setting', redirect: '/setting/vocabulary_set', meta: {
        requiresAuth: true
    }, children: [{
        path: "vocabulary_set", name: 'vocabulary_set_setting', component: VocabularySetSettingView,
    }, {
        path: "vocabulary", name: "vocabulary_setting", component: VocabularySettingView, props: true,
    }, {
        path: 'dictionary', name: 'dictionary_setting', component: DictionarySettingView,
    }, {
        path: 'llm', name: 'llm_setting', component: LLMServiceSettingView,
    }, {
        path: 'aiprompt', name: 'aiprompt_setting', component: AIPromptSettingView,
    }, {
        path: 'user', name: 'user_setting', component: UserSettingView,
    }]
}]

export default SettingRoutes;