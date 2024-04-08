import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ cookies, url }) => {
    if (!cookies.get('session')) {
        throw redirect(307, `/login?redirectTo=${url.pathname}`);
    }
};
