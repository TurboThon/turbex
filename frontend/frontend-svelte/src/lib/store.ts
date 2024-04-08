import { writable } from 'svelte/store';

type UserData = {
	username: string;
	privateKey: string;
	publicKey: string;
};

export const userStore = writable<UserData|undefined>(undefined);
