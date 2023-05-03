import { InjectionKey } from 'vue';
import { createStore, useStore as baseUseStore, Store } from 'vuex';

export interface State {
    logged: boolean
    sid: string
    token: string
    user: any
};

export const key: InjectionKey<Store<State>> = Symbol();

export const store = createStore<State>({
    state: {
        logged: false,
        sid: '',
        token: '',
        user: {},
    }
});

export function useStore() {
    return baseUseStore(key)
};