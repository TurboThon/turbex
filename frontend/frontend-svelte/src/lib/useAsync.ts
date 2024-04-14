export async function useAsync<T>(callback: () => T) {
	return new Promise<T>((resolve, reject) => setTimeout(() => {
    try {
      resolve(callback())
    } catch (error) {
      reject(error)
    }
  }, 0));
}
