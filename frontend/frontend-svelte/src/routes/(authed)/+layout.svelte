<script lang="ts">
	import Header from '$lib/components/Header.svelte';
	import { userStore } from '$lib/store';
	import { onMount } from 'svelte';
	import { getMe } from '$lib/query';
	import { useQuery, type QueryResult } from '$lib/useQuery';
	import { Spinner } from 'flowbite-svelte';

	let userQuery: QueryResult<typeof getMe> | undefined;

	onMount(() => {
		// Only check session cookies if userStore is empty
		if ($userStore == undefined) {
			userQuery = useQuery(getMe(), {
				onSuccess: (data) => {
					userStore.set({
						username: data.data.userName,
						publicKey: data.data.publicKey,
						privateKey: data.data.privateKey
					});
				}
			});
		}
	});

</script>

<Header />
<div class="m-10 space-y-3 rounded-lg bg-white p-5">
	{#if $userQuery?.isLoading || $userStore == undefined}
		<Spinner size={10} />
	{:else if $userQuery?.isError}
		<h3>You must log in to access this page</h3>
		<h4>
			If you are not automatically redirected, click <a href="/login">here</a> to log in.
		</h4>
	{:else if $userQuery?.isSuccess || $userStore}
		<slot />
	{/if}
</div>
