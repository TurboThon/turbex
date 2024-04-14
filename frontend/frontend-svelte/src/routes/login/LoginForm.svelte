<script lang="ts">
	import { Label, Input, Button, Alert, Spinner } from "flowbite-svelte";
	import { useCrypt } from "$lib/useCrypt";
	import { useQuery, type QueryResult } from "$lib/useQuery";
	import { useAsync } from "$lib/useAsync";
	import { onMount } from "svelte";
	import { getLogin, postUser } from "$lib/query";
	import { userStore } from "$lib/store";
	import { goto } from "$app/navigation";
	import { page } from "$app/stores";

	let showAlert = false;
	let alertMessage = "";

	let username = "";
	let password = "";

	let key_password: string | undefined;
	let loginQuery: QueryResult<typeof getLogin> | undefined;
	let signupQuery: QueryResult<typeof postUser> | undefined;

	let loginSpinner = false;
	let signupSpinner = false;

	let crypt: typeof import("turbex-crypt") | undefined;

	$: if ($loginQuery?.isError || $signupQuery?.isError) {
		let error = $loginQuery?.error ?? $signupQuery?.error;
		if (error != undefined) {
			switch (error.status) {
				case 401:
					alertMessage =
						"Username or password invalid, if you don't have a account please sign up first.";
					break;
				case 409:
					alertMessage = "Username already taken, please choose another username.";
					break;
				case 400:
					alertMessage = "This username is invalid.";
			}
			showAlert = true;
		}
	}

	$: if ($userStore) {
		goto($page.url.searchParams.get("redirectTo") ?? "/");
	}

	let loginHandle = async () => {
		loginSpinner = true;

		// Invalidate previous queries result
		loginQuery = undefined;
		signupQuery = undefined;

		let cryptPasswords = await useAsync(() => crypt?.get_api_password_and_key(password));
		if (cryptPasswords == undefined) {
			console.error("turbex-crypt is not loaded yet, please try again later");
			return;
		}
		loginQuery = useQuery(getLogin({ userName: username, password: cryptPasswords.api_password }), {
			onSuccess: (data) => {
				userStore.set({
					username: data.userName,
					privateKey: data.privateKey,
					publicKey: data.publicKey,
					keyPassword: cryptPasswords!.key_password,
				});
			},
		});
		loginSpinner = false;
	};

	let signupHandle = async () => {
		signupSpinner = true;

		// Invalidate previous queries result
		loginQuery = undefined;
		signupQuery = undefined;

		let new_keys = await useAsync(() => crypt?.get_new_keys_and_password(password));
		if (new_keys == undefined) {
			console.error("turbex-crypt is not loaded yet, please try again later");
			return;
		}
		signupQuery = useQuery(
			postUser({
				userName: username,
				password: new_keys.api_password,
				privateKey: new_keys.encrypted_key,
				publicKey: new_keys.pub_key,
			}),
			{
				onSuccess: () => {
					// fetch session token and login
					loginHandle();
				},
			},
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
	<h3 class="col-span-2 mb-3 text-center text-3xl">Welcome to Turbex</h3>
	<p class="col-span-2 text-slate-700">To use turbex you need to log in or create an account</p>
	{#if showAlert}
		<Alert class="col-span-2">{alertMessage}</Alert>
	{/if}
	<div class="col-span-2">
		<Label for="username" color={$loginQuery?.isError || $signupQuery?.isError ? "red" : undefined}
			>Username</Label
		>
		<Input
			bind:value={username}
			type="text"
			id="username"
			placeholder="Username"
			color={$loginQuery?.isError || $signupQuery?.isError ? "red" : undefined}
			required
			class="mb-2"
		/>
	</div>
	<div class="col-span-2">
		<Label for="password" color={$loginQuery?.isError ? "red" : undefined}>Mot de passe</Label>
		<Input
			bind:value={password}
			type="password"
			id="password"
			placeholder="••••••••••••"
			color={$loginQuery?.isError ? "red" : undefined}
			required
			class="mb-2"
		/>
	</div>
	<Button on:click={loginHandle} type="submit" disabled={loginSpinner || $loginQuery?.isLoading}>
		{#if loginSpinner || $loginQuery?.isLoading}
			<Spinner class="mr-3 h-3 w-3" />
		{/if}
		Log in
	</Button>
	<Button
		on:click={signupHandle}
		color="alternative"
		disabled={signupSpinner || $signupQuery?.isLoading}
	>
		{#if signupSpinner || $signupQuery?.isLoading}
			<Spinner class="mr-3 h-3 w-3" />
		{/if}
		Sign up
	</Button>
</form>
