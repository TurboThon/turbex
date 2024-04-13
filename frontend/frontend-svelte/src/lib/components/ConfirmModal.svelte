<script lang="ts">
	import { Modal, Button } from "flowbite-svelte";

	export let showModal = false;
	export let title = "Confirm";
	export let content: any;
	export let confirmText = "Continue";
	export let cancelText = "Cancel";
	export let onConfirm = () => {};
</script>

<Modal {title} bind:open={showModal} autoclose outsideclose on:close>
	<p><slot>{content}</slot></p>
	<svelte:fragment slot="footer">
		<Button on:click={onConfirm}>{confirmText}</Button>
		<Button color="alternative">{cancelText}</Button>
	</svelte:fragment>
</Modal>
