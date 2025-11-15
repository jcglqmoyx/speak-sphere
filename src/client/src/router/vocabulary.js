import VocabularyModuleHomeView from "@/views/vocabulary/VocabularyModuleView.vue";
import VocabularyCardView from "@/views/vocabulary/VocabularyCardView.vue";
import LearnVocabularyView from "@/views/vocabulary/LearnVocabularyView.vue";
import ReviewVocabularyView from "@/views/vocabulary/ReviewVocabularyView.vue";

const VocabularyRoutes = [
    {
        path: '/vocabulary',
        name: 'vocabulary',
        component: VocabularyModuleHomeView,
        meta: {
            requiresAuth: true
        }
    },
    {
        path: '/card',
        name: 'card',
        component: VocabularyCardView,
        meta: {
            requiresAuth: true
        }
    },
    {
        path: '/learn',
        name: 'learn',
        component: LearnVocabularyView,
        meta: {
            requiresAuth: true
        }
    }, {
        path: '/review',
        name: 'review',
        component: ReviewVocabularyView,
        meta: {
            requiresAuth: true
        }
    }
];
export default VocabularyRoutes;