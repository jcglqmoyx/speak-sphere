import {createStore} from 'vuex'
import userModule from './user';

const store = createStore({
    modules: {
        user: userModule
    }
});
export default store;
