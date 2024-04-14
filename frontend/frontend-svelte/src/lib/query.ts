import type { HttpError } from "@sveltejs/kit";
import { handleExpiredSession } from "$lib/store";
import { type File } from "$lib/types/file";

const BACKEND_ROOT = "";

const default_options: RequestInit = {
    mode: "cors",
    credentials: "include",
    headers: new Headers({ "Content-Type": "application/json" }),
};

type GetLoginParams = {
    userName: string;
    password: string;
};

type GetLoginResponse = {
    privateKey: string;
    publicKey: string;
    userName: string;
};

export function getLogin(params: GetLoginParams): () => Promise<GetLoginResponse> {
    return async () => {
        const res = await fetch(`${BACKEND_ROOT}/api/v1/login`, {
            ...default_options,
            method: "POST",
            body: JSON.stringify(params),
        });
        if (!res.ok) {
            let message = (await res.json()).error;
            throw { status: res.status, body: { message } } as HttpError;
        }
        return res.json();
    };
}

type GetMeResponse = { data: GetLoginResponse };

export function getMe(): () => Promise<GetMeResponse> {
    return async () => {
        const res = await fetch(`${BACKEND_ROOT}/api/v1/me`, {
            ...default_options,
            method: "GET",
        });
        if (!res.ok) {
            let message = (await res.json()).error;
            if (res.status == 401) {
                handleExpiredSession();
            }
            throw { status: res.status, body: { message } } as HttpError;
        }
        return res.json();
    };
}

type GetLogoutResponse = string;

export function getLogout(): () => Promise<GetLogoutResponse> {
    return async () => {
        const res = await fetch(`${BACKEND_ROOT}/api/v1/logout`, {
            ...default_options,
            method: "GET",
        });
        if (!res.ok) {
            let message = (await res.json()).error;
            if (res.status == 401) {
                handleExpiredSession();
            }
            throw { status: res.status, body: { message } } as HttpError;
        }
        return res.json();
    };
}

type PostUserParams = {
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
            method: "POST",
            body: JSON.stringify(params),
        });
        if (!res.ok) {
            let message = (await res.json()).error;
            throw { status: res.status, body: { message } } as HttpError;
        }
        return null;
    };
}

type GetUsersResponse = {
    users: {
        userName: string;
    }[];
};

export function getUsers(): () => Promise<GetUsersResponse> {
    return async () => {
        const res = await fetch(`${BACKEND_ROOT}/api/v1/user`, {
            ...default_options,
            method: "GET",
        });
        if (!res.ok) {
            let message = (await res.json()).error;
            if (res.status == 401) {
                handleExpiredSession();
            }
            throw { status: res.status, body: { message } } as HttpError;
        }
        return res.json();
    };
}

type GetUserResponse = {
    publicKey: string;
    userName: string;
};

export function getUser(username: string): () => Promise<GetUserResponse> {
    return async () => {
        const res = await fetch(`${BACKEND_ROOT}/api/v1/user/${username}`, {
            ...default_options,
            method: "GET",
        });
        if (!res.ok) {
            let message = (await res.json()).error;
            if (res.status == 401) {
                handleExpiredSession();
            }
            throw { status: res.status, body: { message } } as HttpError;
        }
        return res.json();
    };
}

type PutUserParams = {
    username: string;
    request: {
        password?: string;
        privateKey?: string;
        publicKey?: string;
    };
};

type PutUserResponse = {
    data: string;
};

export function putUser(params: PutUserParams): () => Promise<PutUserResponse> {
    return async () => {
        const res = await fetch(`${BACKEND_ROOT}/api/v1/user/${params.username}`, {
            ...default_options,
            method: "PUT",
            body: JSON.stringify(params.request),
        });
        if (!res.ok) {
            let message = (await res.json()).error;
            if (res.status == 401) {
                handleExpiredSession();
            }
            throw { status: res.status, body: { message } } as HttpError;
        }
        return res.json();
    };
}

type PostFileParams = {
    fileContent: ArrayBuffer;
    filename: string;
    encryptedFileKey: string;
    ephemeralPubKey: string;
};

type PostFileResponse = {
    fileid: string;
};

export function postFile(params: PostFileParams): () => Promise<PostFileResponse> {
    return async () => {
        const data = new FormData();
        const file = new File([params.fileContent], params.filename);
        data.append("encrypted_file_key", params.encryptedFileKey);
        data.append("ephemeral_pub_key", params.ephemeralPubKey);
        data.append("encrypted_file", file);
        const res = await fetch(`${BACKEND_ROOT}/api/v1/file`, {
            ...default_options,
            method: "POST",
            headers: new Headers(),
            body: data,
        });
        if (!res.ok) {
            let message = (await res.json()).error;
            if (res.status == 401) {
                handleExpiredSession();
            }
            throw { status: res.status, body: { message } } as HttpError;
        }
        return res.json();
    };
}

type PostShareParams = {
    docId: string;
    username: string;
    request: {
        encryptionKey: string;
        ephemeralPubKey: string;
    };
};

type PostShareResponse = string;

export function postShare(params: PostShareParams): () => Promise<PostShareResponse> {
    return async () => {
        const res = await fetch(`${BACKEND_ROOT}/api/v1/share/${params.docId}/${params.username}`, {
            ...default_options,
            method: "POST",
            body: JSON.stringify(params.request),
        });
        if (!res.ok) {
            let message = (await res.json()).error;
            if (res.status == 401) {
                handleExpiredSession();
            }
            throw { status: res.status, body: { message } } as HttpError;
        }
        return res.json();
    };
}

type GetFilesResponse = File[];

export function getFiles(): () => Promise<GetFilesResponse> {
    return async () => {
        const res = await fetch(`${BACKEND_ROOT}/api/v1/file`, {
            ...default_options,
            method: "GET",
        });
        if (!res.ok) {
            let message = (await res.json()).error;
            if (res.status == 401) {
                handleExpiredSession();
            }
            throw { status: res.status, body: { message } } as HttpError;
        }
        return (await res.json()).data;
    };
}

type GetFileResponse = Blob;

export function getFile(fileId: string): () => Promise<GetFileResponse> {
    return async () => {
        const res = await fetch(`${BACKEND_ROOT}/api/v1/file/${fileId}`, {
            ...default_options,
            method: "GET",
        });
        if (!res.ok) {
            let message = (await res.json()).error;
            if (res.status == 401) {
                handleExpiredSession();
            }
            throw { status: res.status, body: { message } } as HttpError;
        }
        return res.blob();
    };
}
