<script lang="ts">
  import { Modal, Button } from 'flowbite-svelte';
  import { createEventDispatcher } from 'svelte';

  export let showModal = false;
  export let buttonText: string;
  export let title = "Confirm";
  export let content;
  export let confirmText = "Continue";
  export let cancelText = "Cancel";

  const dispatch = createEventDispatcher();
  const onConfirm = () => {
    dispatch("confirm");
  }
</script>

<Button on:click={() => showModal = true}>{buttonText}</Button>
<Modal title={title} bind:open={showModal} autoclose outsideclose on:close>
  <p>{content}</p>
  <svelte:fragment slot="footer">
    <Button on:click={onConfirm}>{confirmText}</Button>
    <Button color="alternative">{cancelText}</Button>
  </svelte:fragment>
</Modal>

