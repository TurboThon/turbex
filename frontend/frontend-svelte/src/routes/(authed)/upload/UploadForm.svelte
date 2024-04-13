<script lang="ts">
	import { onMount } from "svelte";
	import { Label, Toggle, Hr, MultiSelect, Tooltip, Spinner, Alert, Button } from "flowbite-svelte";
	import ConfirmModal from "$lib/components/ConfirmModal.svelte";
	import { Dropzone } from "flowbite-svelte";
	import { useQuery, type QueryResult } from "$lib/useQuery";
	import { getUsers, getUser, postFile, postShare } from "$lib/query";
	import { userStore } from "$lib/store";
	import { useCrypt } from "$lib/useCrypt";
	import { useAsync } from "$lib/useAsync";

	let crypt: typeof import("turbex-crypt") | undefined;

	let file: File | undefined;
	let encryptedFile: ArrayBuffer | undefined;
	let encryptionLoading = false;
	let encryptionLoadingMessage = "";
	let key: string | undefined;
	let recipients: string[] | undefined;
	let uploadLoading = false;
	let uploadSuccessAlert = false;
	let confirmModal = false;

	const dropHandle = (event: DragEvent) => {
		event.preventDefault();
		let uploadedFile: File | undefined | null;
		if (!event.dataTransfer) return;
		if (event.dataTransfer.items) {
			let item = event.dataTransfer.items[0];
			if (item.kind === "file") {
				uploadedFile = item.getAsFile();
			}
		} else {
			uploadedFile = event.dataTransfer.files[0];
		}
		if (uploadedFile) {
			file = uploadedFile;
		}
	};

	const handleChange = async (event: Event) => {
		// The following type was infered based on runtime tests
		const eventTarget = event.target as HTMLInputElement;
		const files = eventTarget.files;
		if (files && files.length > 0) {
			file = files[0];
		}
	};

	let usersQuery: QueryResult<typeof getUsers> | undefined;

	$: if (file) encryptFile();

	const encryptFile = async () => {
		encryptionLoading = true;
		encryptionLoadingMessage = "Encryption key generation...";
		key = await useAsync(() => crypt?.generate_aes_key());
		if (key) {
			encryptionLoadingMessage = "Loading file...";
			const fileContent = await file!.arrayBuffer();
			encryptionLoadingMessage = "Protecting the file with TurbexEncryption...";
			encryptedFile = await useAsync(() => crypt?.encrypt_file(new Uint8Array(fileContent), key!));
		}
		encryptionLoading = false;
	};

	const uploadFile = async () => {
		uploadLoading = true;
		if (!file || !key || !encryptedFile || !crypt || !$userStore) return;
		const recipientKey = await encryptFileKey($userStore.publicKey);
		if (!recipientKey) return;
		let docId = await postFile({
			fileContent: encryptedFile,
			filename: file.name,
			ephemeralPubKey: recipientKey.ephemeral_pub_key,
			encryptedFileKey: recipientKey.encrypted_pfk,
		})();
		let shares = recipients?.map((recipient) => shareFileWithUser(docId.fileid, recipient));
		if (shares) await Promise.all(shares);
		uploadLoading = false;
		uploadSuccessAlert = true;
	};

	const shareFileWithUser = async (docId: string, username: string) => {
		if (key && encryptedFile && crypt) {
			const recipient = await getUser(username)();
			const recipientKey = await encryptFileKey(recipient.publicKey);
			if (recipientKey) {
				postShare({
					docId,
					username,
					request: {
						encryptionKey: recipientKey.encrypted_pfk,
						ephemeralPubKey: recipientKey.ephemeral_pub_key,
					},
				})();
			}
		}
	};

	const encryptFileKey = async (publicKey: string) => {
		return useAsync(() => crypt!.encrypt_pfk(key!, publicKey));
	};

	onMount(() => {
		usersQuery = useQuery(getUsers());
		useCrypt((crypt_lib) => {
			crypt = crypt_lib;
		});
	});
</script>

<div class="px-5">
	<Dropzone
		id="dropzone"
		on:drop={dropHandle}
		on:dragover={(event) => {
			event.preventDefault();
		}}
		on:change={handleChange}
	>
		<svg
			aria-hidden="true"
			class="mb-3 h-10 w-10 text-gray-400"
			fill="none"
			stroke="currentColor"
			viewBox="0 0 24 24"
			xmlns="http://www.w3.org/2000/svg"
			><path
				stroke-linecap="round"
				stroke-linejoin="round"
				stroke-width="2"
				d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
			/></svg
		>
		{#if !file}
			<p class="mb-2 text-sm text-gray-500 dark:text-gray-400">
				<span class="font-semibold">Click to upload</span> or drag and drop
			</p>
			<p class="text-xs text-gray-500 dark:text-gray-400">Anything you want, limited to one file</p>
		{:else}
			<p>{file.name}</p>
		{/if}
		{#if encryptionLoading}
			<div class="flex flex-row items-center">
				<Spinner />
				<p class="ml-4">{encryptionLoadingMessage}</p>
			</div>
		{/if}
	</Dropzone>
	<Hr />
	<Label for="people-select">Select people with whom the file will be shared</Label>
	{#if !usersQuery || $usersQuery?.isLoading}
		<Spinner />
	{:else if $usersQuery?.isSuccess}
		<MultiSelect
			id="people-select"
			items={$usersQuery?.data?.users
				.filter((user) => user.userName != $userStore?.username)
				.map((user) => {
					return { value: user.userName, name: user.userName };
				})}
			bind:value={recipients}
		/>
	{/if}
	<Hr />
	<div class="flex flex-row justify-between">
		<div>
			<Label>Sharing options</Label>
			<Toggle class="mt-3" checked disabled>Turbex encryption</Toggle>
			<Tooltip>For safety reasons, you cannot turn off Turbex encryption.</Tooltip>
			<Toggle class="mt-3" checked disabled>Keep a file for me</Toggle>
			<Tooltip>Save an encrypted copy of the file for yourself.</Tooltip>
		</div>
		<Button
			class="h-14 w-64 self-center justify-self-center text-base font-bold"
			on:click={() => (confirmModal = true)}
			disabled={uploadLoading}
		>
			{#if uploadLoading}
				<Spinner />
			{/if}
        <p>Send file</p>
		</Button>
		<ConfirmModal
			content="Do you want to upload?"
      showModal={confirmModal}
			onConfirm={() => {
				uploadFile();
			}}
		/>
	</div>
    <div class="h-16">
	{#if uploadSuccessAlert}
		<Alert class="mt-4" dismissable color="green">Your file has been uploaded successfully</Alert>
  {/if}
    </div>
</div>
