import BookSettingView from "@/views/book/BookSettingView.vue";
import WordSettingView from "@/views/word/WordSettingView.vue";
import DictionarySettingView from "@/views/dictionary/DictionarySettingView.vue";
import UserSettingView from "@/views/user/UserSettingView.vue";

const SettingRoutes = [
    {
        path: '/setting',
        name: 'setting',
        redirect: '/setting/book',
        meta: {
            requiresAuth: true
        },
        children: [
            {
                path: "book",
                name: 'book_setting',
                component: BookSettingView,
            },
            {
                path: "word",
                name: "word_setting",
                component: WordSettingView,
                props: true,
            },
            {
                path: 'dictionary',
                name: 'dictionary_setting',
                component: DictionarySettingView,
            },
            {
                path: 'user',
                name: 'user_setting',
                component: UserSettingView,
            }
        ]
    }
]

export default SettingRoutes;