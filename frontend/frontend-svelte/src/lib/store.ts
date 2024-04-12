import { writable } from 'svelte/store';
import { goto } from '$app/navigation';

type UserData = {
    username: string;
    privateKey: string;
    publicKey: string;
};

export const userStore = writable<UserData | undefined>(undefined);

export const handleExpiredSession = () => {
    // Only executes on browser
    if (window) {
        userStore.set(undefined);
        goto(`/login?redirectTo${window.location.pathname}`);
    }
};
