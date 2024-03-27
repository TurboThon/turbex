<script lang="ts">
	import Header from '$lib/components/Header.svelte';
	import { userStore } from '$lib/store';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	onMount(() => {
		if (!$userStore) {
			goto('/login');
		}
	});
</script>

<Header />
<div class="m-10 space-y-3 rounded-lg bg-white p-5">
	<slot />
</div>
