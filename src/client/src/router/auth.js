import LoginView from "@/views/user/LoginView.vue";
import RegisterView from "@/views/user/RegisterView.vue";

const AuthRoutes = [
    {
        path: '/login',
        name: 'login',
        component: LoginView,
        beforeEnter: (to, from, next) => {
            // 在进入登录页前，等待页面渲染完成后设置焦点
            next(vm => {
                vm.$nextTick(() => {
                    // 例如，如果登录菜单项有 ref 为 'loginMenuItem'
                    vm.$refs.loginMenuItem.focus();
                });
            });
        },
    },
    {
        path: '/register',
        name: 'register',
        component: RegisterView
    },
]

export default AuthRoutes;