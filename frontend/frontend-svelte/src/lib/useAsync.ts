export async function useAsync<T>(callback: () => T) {
	return new Promise<T>((resolve) => setTimeout(() => resolve(callback()), 0));
}
