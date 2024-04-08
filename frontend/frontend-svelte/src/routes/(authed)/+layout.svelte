<script lang="ts">
	import Header from '$lib/components/Header.svelte';
	import { userStore } from '$lib/store';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { getMe } from '$lib/query';
	import { useQuery, type QueryResult } from '$lib/useQuery';
    import {Spinner} from "flowbite-svelte"

    let userQuery: QueryResult<typeof getMe> | undefined

	onMount(() => {
        // Only check session cookies if userStore is empty
		if ($userStore == undefined) {
            userQuery = useQuery(getMe())
		}
	});

    $: if ($userQuery!=undefined && !$userQuery.isLoading) {
        if ($userQuery?.isError) {
            userStore.set(undefined)
            goto('/login')
        } else if ($userQuery?.isSuccess && $userQuery.data != undefined) {
            userStore.set({
                username: $userQuery.data.userName,
                publicKey: $userQuery.data.publicKey,
                privateKey: $userQuery.data.privateKey
            })
        }
    }
</script>

<Header />
<div class="m-10 space-y-3 rounded-lg bg-white p-5">
    {#if $userQuery?.isLoading || $userStore == undefined}
    <Spinner size={10}/>
    {:else if $userQuery?.isError} 
    <h3>Vous devez être connecté pour accéder à cette page</h3>
    <h4>Si vous n'êtes pas automatiquement redirigé, cliquez <a href="/login">ici</a> pour vous connecter</h4>
    {:else}
	<slot />
    {/if}
</div>
