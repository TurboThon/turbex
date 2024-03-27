import { readable, type Readable } from 'svelte/store';
import type { HttpError } from '@sveltejs/kit';

type QueryContext<T> = {
	isLoading: boolean;
	isSuccess: boolean;
	isError: boolean;
	data: T | undefined;
	error: HttpError | undefined;
};

export type QueryContextStore<T> = Readable<QueryContext<T>>;

// Type to get the type of the store returned by a specific query
// Use it as follows QueryResult<typeof query>
export type QueryResult<T extends (args: any) => () => Promise<any>> = QueryContextStore<
	Awaited<ReturnType<ReturnType<T>>>
>;

type QueryHooks<T> = {
	onSuccess: (data: T) => void;
	onError: (err: HttpError) => void;
};

// This hooks runs the async queryCallback
// and returns a store which updates as the query resolves
// It also runs the passed hooks when the promise resolves or fails
export function useQuery<T>(
	queryCallback: () => Promise<T>,
	hooks?: Partial<QueryHooks<T>>
): QueryContextStore<T> {
	// Start query
	const query = queryCallback();
	// Populate store with query result
	const queryContext = readable<QueryContext<T>>(
		{
			isLoading: true,
			isSuccess: false,
			isError: false,
			data: undefined,
			error: undefined
		},
		(_, update) => {
			query
				.then((data) => {
					update((context) => {
						return { ...context, isSuccess: true, data };
					});
				})
				.catch((err: HttpError) => {
					update((context) => {
						return { ...context, isError: true, error: err };
					});
				})
				.finally(() => {
					update((context) => {
						return { ...context, isLoading: false };
					});
				});
		}
	);
	// Execute hooks
	if (hooks?.onSuccess != undefined) {
		query.then((data) => hooks.onSuccess!(data));
	}
	if (hooks?.onError != undefined) {
		query.catch((err) => hooks.onError!(err));
	}
	return queryContext;
}
