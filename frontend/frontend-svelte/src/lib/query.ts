import type { HttpError } from '@sveltejs/kit';

const BACKEND_ROOT = 'http://localhost:8000';

const default_options: RequestInit = {
	mode: 'cors',
	credentials: 'include',
	headers: { content: 'application/json' }
};

type GetLoginParams = {
	userName: string;
	password: string;
};

type GetLoginResponse = {
	firstName: string;
	lastName: string;
	privateKey: string;
	publicKey: string;
	userName: string;
};

export function getLogin(params: GetLoginParams): () => Promise<GetLoginResponse> {
	return async () => {
		const res = await fetch(`${BACKEND_ROOT}/api/v1/login`, {
			...default_options,
			method: 'POST',
			body: JSON.stringify(params)
		});
		if (!res.ok) {
			let message = (await res.json()).error;
			throw { status: res.status, body: { message } } as HttpError;
		}
		return res.json();
	};
}

type PostUserParams = {
	firstName: string;
	lastName: string;
	password: string;
	privateKey: string;
	publicKey: string;
	userName: string;
};

type PostUserResponse = null;

export function postUser(params: PostUserParams): () => Promise<PostUserResponse> {
	return async () => {
		const res = await fetch(`${BACKEND_ROOT}/api/v1/user`, {
			...default_options,
			method: 'POST',
			body: JSON.stringify(params)
		});
		if (!res.ok) {
			let message = (await res.json()).error;
			throw { status: res.status, body: { message } } as HttpError;
		}
		return null;
	};
}