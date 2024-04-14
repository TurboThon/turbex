<script lang="ts">
	import { Modal, Button } from "flowbite-svelte";
    import { type Writable } from "svelte/store";

	export let showModal: Writable<boolean>;
	export let title = "Confirm";
	export let content = "";
	export let confirmText = "Continue";
	export let cancelText = "Cancel";
	export let onConfirm = async () => {};
</script>

<Modal {title} bind:open={$showModal} autoclose outsideclose>
	<p><slot>{content}</slot></p>
	<svelte:fragment slot="footer">
		<Button on:click={onConfirm}>{confirmText}</Button>
		<Button color="alternative">{cancelText}</Button>
	</svelte:fragment>
</Modal>
