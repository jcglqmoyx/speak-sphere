import WordModuleHomeView from "@/views/word/WordModuleView.vue";
import WordCardView from "@/views/word/WordCardView.vue";
import LearnWordView from "@/views/word/LearnWordView.vue";
import ReviewWordView from "@/views/word/ReviewWordView.vue";

const WordRoutes = [
    {
        path: '/word',
        name: 'word',
        component: WordModuleHomeView,
        meta: {
            requiresAuth: true
        }
    },
    {
        path: '/card',
        name: 'card',
        component: WordCardView,
        meta: {
            requiresAuth: true
        }
    },
    {
        path: '/learn',
        name: 'learn',
        component: LearnWordView,
        meta: {
            requiresAuth: true
        }
    }, {
        path: '/review',
        name: 'review',
        component: ReviewWordView,
        meta: {
            requiresAuth: true
        }
    }
];
export default WordRoutes;