<script lang="ts">
	import {
		Table,
		TableBody,
		TableBodyCell,
		TableBodyRow,
		TableHead,
		TableHeadCell,
		Button,
	} from "flowbite-svelte";
	import { DownloadSolid, TrashBinSolid } from "flowbite-svelte-icons";
	import { type File } from "$lib/types/file";
	import { getFiles, getFile } from "$lib/query";
	import { useQuery, type QueryResult } from "$lib/useQuery";
	import { onMount } from "svelte";
	import { useCrypt } from "$lib/useCrypt";
	import { useAsync } from "$lib/useAsync";
	import { userStore } from "$lib/store";

	let crypt: typeof import("turbex-crypt") | undefined;

	const doDownload = async (file: File) => {
		const encryptedContent = await (await getFile(file.id)()).arrayBuffer();
		const content = await useAsync(() => {
			if ($userStore && crypt) {
				return crypt.decrypt_file(
					new Uint8Array(encryptedContent),
					file.encryptionKey,
					file.ephemeralPubKey,
					$userStore.privateKey,
					$userStore.keyPassword,
				);
			}
		});
		if (content) {
            const decryptedFile = new File([content], file.filename)
            const fileUrl = URL.createObjectURL(decryptedFile)
            const dowloadLink = document.createElement("a")
            dowloadLink.href = fileUrl
            dowloadLink.download = file.filename
            dowloadLink.click()
		}
	};

	const doDelete = () => {
		console.log("deleting");
	};

	let filesQuery: QueryResult<typeof getFiles> | undefined;

	onMount(() => {
		filesQuery = useQuery(getFiles());
		useCrypt((crypt_lib) => {
			crypt = crypt_lib;
		});
	});
</script>

<Table>
	<TableHead>
		<TableHeadCell class="grow">File name</TableHeadCell>
		<TableHeadCell class="grow-0">Owner</TableHeadCell>
		<TableHeadCell class="grow-0">Shared date</TableHeadCell>
		<TableHeadCell class="grow-0">Actions</TableHeadCell>
	</TableHead>
	<TableBody>
		{#if $filesQuery?.isSuccess && $filesQuery?.data}
			{#each $filesQuery.data as file}
				<TableBodyRow>
					<TableBodyCell class="grow">{file.filename}</TableBodyCell>
					<TableBodyCell class="grow-0">{file.senderUserName != $userStore?.username ? file.senderUserName : "You"}</TableBodyCell>
					<TableBodyCell class="grow-0">{new Date(file.uploadDate).toLocaleDateString()}</TableBodyCell>
					<TableBodyCell class="flex grow-0 space-x-4">
						<Button on:click={() => doDownload(file)} size="xs" color="green"
							><DownloadSolid size="md" /></Button
						>
						<Button on:click={doDelete} size="xs"><TrashBinSolid size="md" /></Button>
					</TableBodyCell>
				</TableBodyRow>
			{/each}
		{/if}
	</TableBody>
</Table>
