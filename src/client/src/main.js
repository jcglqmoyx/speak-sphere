import {createApp} from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'

const app = createApp(App);
app.use(store);
router.beforeEach((to, from, next) => {
    const requiresAuth = to.meta.requiresAuth;
    if (requiresAuth && !store.state.user.isAuthenticated) {
        next('/login');
    } else {
        next();
    }
});
app.use(router);
app.mount("#app");