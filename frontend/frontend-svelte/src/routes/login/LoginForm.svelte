<script lang="ts">
	import { Label, Input, Button, Alert, Spinner } from 'flowbite-svelte';
	import { useCrypt } from '$lib/useCrypt';
	import { useQuery, type QueryResult } from '$lib/useQuery';
	import { useAsync } from '$lib/useAsync';
	import { onMount } from 'svelte';
	import { getLogin, postUser } from '$lib/query';
	import { userStore } from '$lib/store';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';

	let showAlert = false;
	let alertMessage = '';

	let username = '';
	let password = '';

	let loginQuery: QueryResult<typeof getLogin> | undefined;
	let signupQuery: QueryResult<typeof postUser> | undefined;

	let loginSpinner = false;
	let signupSpinner = false;

	let crypt: typeof import('turbex-crypt') | undefined;

	$: if ($loginQuery?.isSuccess && $loginQuery.data) {
		userStore.set({
			username: $loginQuery.data.userName,
			privateKey: $loginQuery.data.privateKey,
			publicKey: $loginQuery.data.publicKey
		});
	}

	$: if ($loginQuery?.isError || $signupQuery?.isError) {
		let error = $loginQuery?.error ?? $signupQuery?.error;
		if (error != undefined) {
			switch (error.status) {
				case 401:
					alertMessage =
						"Nom d'utilisateur/Mot de passe incorrects, si vous n'avez pas de compte cliquez sur \"S'inscrire\"";
					break;
				case 409:
					alertMessage =
						"Ce nom d'utilisateur est déjà utilisé. Choisissez un autre nom ou connectez-vous.";
					break;
				case 400:
					alertMessage = "Ce nom d'utilisateur est invalide.";
			}
			showAlert = true;
		}
	}

	$: if ($userStore) {
		goto($page.url.searchParams.get('redirectTo') ?? '/');
	}

	let loginHandle = async () => {
		loginSpinner = true;

		// Invalidate previous queries result
		loginQuery = undefined;
		signupQuery = undefined;

		console.log('a');

		let cryptPasswords = await useAsync(() => crypt?.get_api_password_and_key(password));
		console.log('b');
		if (cryptPasswords == undefined) {
			console.error('turbex-crypt is not loaded yet, please try again later');
			return;
		}
		loginQuery = useQuery(getLogin({ userName: username, password: cryptPasswords.api_password }));
		console.log('finished');
		loginSpinner = false;
	};

	$: if ($userStore != undefined) {
	}

	let signupHandle = async () => {
		signupSpinner = true;

		// Invalidate previous queries result
		loginQuery = undefined;
		signupQuery = undefined;

		let new_keys = await useAsync(() => crypt?.get_new_keys_and_password(password));
		if (new_keys == undefined) {
			console.error('turbex-crypt is not loaded yet, please try again later');
			return;
		}
		signupQuery = useQuery(
			postUser({
				firstName: username,
				lastName: username,
				userName: username,
				password: new_keys.api_password,
				privateKey: new_keys.encrypted_key,
				publicKey: new_keys.pub_key
			}),
			{
				onSuccess: () => {
					userStore.set({
						username,
						privateKey: new_keys!.encrypted_key,
						publicKey: new_keys!.pub_key
					});
				}
			}
		);
		signupSpinner = false;
	};

	onMount(async () => {
		// Populates crypt with the functions from turbex-crypt
		useCrypt((crypt_lib) => {
			crypt = crypt_lib;
		});
	});
</script>

<form class="grid gap-4 md:grid-cols-2">
	<h3 class="col-span-2 mb-3 text-center text-3xl">Bienvenue sur Turbex</h3>
	<p class="col-span-2 text-slate-700">
		Pour utiliser Turbex vous devez vous connecter ou créer un compte
	</p>
	{#if showAlert}
		<Alert class="col-span-2">{alertMessage}</Alert>
	{/if}
	<div class="col-span-2">
		<Label for="username" color={$loginQuery?.isError || $signupQuery?.isError ? 'red' : undefined}
			>Nom d'utilisateur</Label
		>
		<Input
			bind:value={username}
			type="text"
			id="username"
			placeholder="Nom d'utilisateur"
			color={$loginQuery?.isError || $signupQuery?.isError ? 'red' : undefined}
			required
			class="mb-2"
		/>
	</div>
	<div class="col-span-2">
		<Label for="password" color={$loginQuery?.isError ? 'red' : undefined}>Mot de passe</Label>
		<Input
			bind:value={password}
			type="password"
			id="password"
			placeholder="••••••••••••"
			color={$loginQuery?.isError ? 'red' : undefined}
			required
			class="mb-2"
		/>
	</div>
	<Button on:click={loginHandle} type="submit" disabled={loginSpinner || $loginQuery?.isLoading}>
		{#if loginSpinner || $loginQuery?.isLoading}
			<Spinner class="mr-3 h-3 w-3" />
		{/if}
		Se connecter
	</Button>
	<Button
		on:click={signupHandle}
		color="alternative"
		disabled={signupSpinner || $signupQuery?.isLoading}
	>
		{#if signupSpinner || $signupQuery?.isLoading}
			<Spinner class="mr-3 h-3 w-3" />
		{/if}
		S'inscrire
	</Button>
</form>