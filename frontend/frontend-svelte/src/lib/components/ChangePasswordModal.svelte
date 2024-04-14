<script lang="ts">
	import { Alert, Button, Hr, Input, Label, Modal, Spinner, Toast } from "flowbite-svelte";
	import { onMount } from "svelte";
	import { useCrypt } from "$lib/useCrypt";
	import { useAsync } from "$lib/useAsync";
	import { userStore } from "$lib/store";
	import { putUser } from "$lib/query";
	import { useQuery } from "$lib/useQuery";
	import { CheckCircleSolid } from "flowbite-svelte-icons";
	import { type Writable } from "svelte/store";

	let crypt: typeof import("turbex-crypt") | undefined;

	let oldPassword = "";
	let newPassword = "";
	let repeatNewPassword = "";
	let alertMessage = "";
	let passwordChangedToast = false;
	let loading = false;

	export let showModal: Writable<boolean>;
	export let onConfirm = async () => {
		if (!newPasswordAreEquals) {
			return;
		}
		if (!crypt) {
			console.error("Crypt is not loaded");
			alertMessage = "Turbex encryption is not loaded, try again later";
			return;
		}
		alertMessage = "";
		loading = true;
		try {
			const priv_key_and_password = await useAsync(() =>
				crypt?.change_password(oldPassword, newPassword, $userStore!.privateKey),
			);

			if (!priv_key_and_password) {
				console.error("Error in crypt, cannot recover");
				alertMessage = "Unexpected error";
				loading = false;
				return;
			}
			useQuery(
				putUser({
					username: $userStore!.username,
					request: {
						privateKey: priv_key_and_password.encrypted_key,
						password: priv_key_and_password.api_password,
					},
				}),
				{
					onSuccess: () => {
						userStore.set({
							...$userStore!,
							privateKey: priv_key_and_password.encrypted_key,
						});
						setTimeout(() => (passwordChangedToast = false), 5000);
						clearFields();
						passwordChangedToast = true;
						loading = false;
						showModal.set(false);
					},

					onError: (err) => {
						throw err;
					},
				},
			);
		} catch (error) {
			alertMessage = typeof error === "string" ? error : "Unexpected error";
			loading = false;
			return;
		}
	};

	$: newPasswordAreEquals = newPassword === repeatNewPassword;

	const clearFields = () => {
		oldPassword = "";
		newPassword = "";
		repeatNewPassword = "";
	};

	onMount(async () => {
		// Populates crypt with the functions from turbex-crypt
		useCrypt((crypt_lib) => {
			crypt = crypt_lib;
		});
	});
</script>

<Modal title="Change password" bind:open={$showModal} outsideclose on:close={clearFields}>
	<form>
		<p>
			Changing your password is a smooth process, we will re-encrypt your private keys with your new
			password.
		</p>
		<p>As of now, Turbex is not configured to invalidate your other sessions.</p>
		<Hr classHr="h-px my-5 bg-gray-200 border-0 dark:bg-gray-700" />

		{#if alertMessage !== ""}
			<Alert>{alertMessage}</Alert>
		{/if}
		<Label for="oldPassword" color={false ? "red" : undefined}>Old password</Label>
		<Input
			bind:value={oldPassword}
			type="password"
			id="oldPassword"
			placeholder="••••••••••••"
			color={false ? "red" : undefined}
			required
			class="mb-2"
		/>
		<Label for="newPassword" color={!newPasswordAreEquals ? "red" : undefined}>New password</Label>
		<Input
			bind:value={newPassword}
			type="password"
			id="newPassword"
			placeholder="••••••••••••"
			color={false ? "red" : undefined}
			required
			class="mb-2"
		/>
		<Label for="repeatNewPassword" color={!newPasswordAreEquals ? "red" : undefined}>
			Repeat new password
		</Label>
		<Input
			bind:value={repeatNewPassword}
			type="password"
			id="repeatNewPassword"
			placeholder="••••••••••••"
			color={false ? "red" : undefined}
			required
			class="mb-2"
		/>
	</form>

	<svelte:fragment slot="footer">
		<Button on:click={onConfirm} disabled={!newPasswordAreEquals || newPassword === "" || loading}>
			{#if loading}
				<Spinner class="mr-3 h-3 w-3" />
			{/if}
			Continue
		</Button>
		<Button color="alternative" on:click={() => showModal.set(false)}>Cancel</Button>
	</svelte:fragment>
</Modal>

<Toast position="bottom-right" dismissable={true} bind:open={passwordChangedToast}>
	<CheckCircleSolid slot="icon" class="h-5 w-5" />
	Your password has been changed!
</Toast>
