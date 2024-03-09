<script lang="ts">
	import '../app.pcss';
  import Header from '$lib/components/Header.svelte';
</script>

<Header/>
<div class="m-10 p-5 bg-white rounded-lg space-y-3">
  <slot />
</div>
