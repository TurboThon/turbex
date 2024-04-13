import { writable } from 'svelte/store';
import { goto } from '$app/navigation';
import { browser } from '$app/environment';

type UserData = {
    username: string;
    keyPassword: string;
    privateKey: string;
    publicKey: string;
};

export const userStore = writable<UserData | undefined>(undefined);

export const handleExpiredSession = () => {
    // Only executes on browser
    if (browser) {
      userStore.set(undefined);
      goto(`/login?redirectTo=${window.location.pathname}`);
    }
};
