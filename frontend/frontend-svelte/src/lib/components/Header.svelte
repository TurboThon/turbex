<script lang="ts">
	import { page } from "$app/stores";
	import { useCrypt } from "$lib/useCrypt";
	import {
		Navbar,
		NavBrand,
		NavLi,
		NavUl,
		NavHamburger,
		Avatar,
		Dropdown,
		DropdownHeader,
		DropdownItem,
		DropdownDivider,
		Toast,
	} from "flowbite-svelte";
	import { handleExpiredSession, userStore } from "$lib/store";
	import { putUser, getLogout } from "$lib/query";
	import ConfirmModal from "./ConfirmModal.svelte";
	import { useQuery } from "$lib/useQuery";
	import { useAsync } from "$lib/useAsync";
	import { onMount } from "svelte";
	import { CheckCircleSolid } from "flowbite-svelte-icons";
	import ChangePasswordModal from "./ChangePasswordModal.svelte";
    import { writable } from "svelte/store";

	$: activeUrl = $page.url.pathname;

	let crypt: typeof import("turbex-crypt") | undefined;
	let confirmRotateKeys = writable(false);
	let rotateKeysToast = false;
	let changePasswordModal = writable(false);

	const rotateKeys = async () => {
		if (!$userStore) {
			console.error("Cannot find key_password, cannot rotate keys without it");
			return;
		}
		let new_keys = await useAsync(() =>
			crypt?.get_new_keys_from_key_password($userStore!.keyPassword),
		);
		if (new_keys) {
			useQuery(
				putUser({
					username: $userStore!.username,
					request: {
						privateKey: new_keys.encrypted_key,
						publicKey: new_keys.pub_key,
					},
				}),
				{
					onSuccess: () => {
						userStore.set({
							...$userStore!,
							privateKey: new_keys!.encrypted_key,
							publicKey: new_keys!.pub_key,
						});
					},
				},
			);
		}
	};

	onMount(async () => {
		// Populates crypt with the functions from turbex-crypt
		useCrypt((crypt_lib) => {
			crypt = crypt_lib;
		});
	});

	const handleLogout = () => {
		getLogout()();
		handleExpiredSession();
	};
</script>

<Navbar>
	<NavBrand href="/">
		<p1 class="text-xl"><b>Turbex</b> - Turbo Extreme file sharing system</p1>
	</NavBrand>
	<NavHamburger />
	<NavUl {activeUrl}>
		<NavLi href="/">My files</NavLi>
		<NavLi href="/upload">Upload a file</NavLi>
		<NavLi href="/admin">Admin panel</NavLi>
	</NavUl>
	<div>
		<Avatar id="avatar-menu" style="cursor: pointer;" />
	</div>
	<Dropdown placement="bottom" triggeredBy="#avatar-menu">
		<DropdownHeader>
			<span class="block text-sm">{$userStore?.username ?? "Anonymous"}</span>
		</DropdownHeader>
		<DropdownItem
			on:click={() => {
				confirmRotateKeys.set(true);
			}}>Rotate keys</DropdownItem
		>
		<DropdownItem
			on:click={() => {
				changePasswordModal.set(true);
			}}>Change Password</DropdownItem
		>
		<DropdownDivider />
		<DropdownItem on:click={handleLogout}>Log out</DropdownItem>
	</Dropdown>
</Navbar>

<ConfirmModal
	title="Rotate keys"
	bind:showModal={confirmRotateKeys}
	onConfirm={async () => {
		await rotateKeys();
		rotateKeysToast = true;
		setTimeout(() => {
			rotateKeysToast = false;
		}, 5000);
	}}
>
	Renewing your keys is a good practice if you have used them for a long time.
	<br />
	However as of now, Turbex is not able to ensure you a smooth renewing, thus you will be unable to download
	files shared with your previous keys.
	<br />
	<b>Therefore we highly recommend you to download all your files before renewing your keys.</b>
	<br />
	Are you sure you want to rotate your keys?
</ConfirmModal>

<ChangePasswordModal bind:showModal={changePasswordModal}></ChangePasswordModal>

<Toast position="bottom-right" dismissable={true} bind:open={rotateKeysToast}>
	<CheckCircleSolid slot="icon" class="h-5 w-5" />
	Your keys have been renewed!
</Toast>
