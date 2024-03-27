import { onMount } from 'svelte';

type CryptCallback = (crypt: typeof import('turbex-crypt')) => any;

const useCrypt = async (callback: CryptCallback) => {
	let crypt = await import('turbex-crypt');
	await crypt.default();
	callback(crypt);
};

export { useCrypt };
