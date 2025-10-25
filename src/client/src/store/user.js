import axios from "axios";

export default {
    state: {
        isAuthenticated: !!localStorage.getItem('token'),
        token: localStorage.getItem('token') || '',
    },
    mutations: {
        setServerLink(state, serverLink) {
            localStorage.setItem('server_link', serverLink);
        },
        setAuth(state, isAuthenticated) {
            state.isAuthenticated = isAuthenticated;
        },
        setToken(state, token) {
            state.token = token;
            localStorage.setItem('token', token);
        },
        logout(state) {
            state.isAuthenticated = false;
            state.token = null;
            localStorage.removeItem("book_id");
            localStorage.removeItem('server_link');
            localStorage.removeItem('token');
        }
    },
    actions: {
        async register(_, {serverLink, username, email, avatar, password, confirm_password}) {
            try {
                const response = await axios.post(serverLink + '/user/register', {
                    'username': username,
                    'email': email,
                    'avatar': avatar,
                    'password': password,
                    'confirm_password': confirm_password,
                });
                const data = response.data;
                if (data.code === 0) {
                    return true;
                } else {
                    return data.message;
                }
            } catch (error) {
                return "发生了未知错误, 注册失败";
            }
        },
        async login({commit}, {serverLink, username, password}) {
            try {
                const response = await axios.post(serverLink + '/user/login', {
                    'username': username,
                    'password': password,
                });
                const data = response.data;
                if (data.code === 0) {
                    commit('setServerLink', serverLink);
                    commit('setAuth', true);
                    commit('setToken', data.data.token);
                    return true;
                } else {
                    return data.message;
                }
            } catch (error) {
                return '登录失败，请检查链接是否正确，后台服务是否开启';
            }
        },
        logout({commit}) {
            commit('logout');
        },
    },
    getters: {
        isAuthenticated: (state) => {
            return state.isAuthenticated;
        }
    }
};