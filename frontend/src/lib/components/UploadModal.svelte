<script lang="ts">
  import { Label, Fileupload, Toggle } from 'flowbite-svelte';
</script>

<div class="container px-5">
  <div class="item item-1"><Label class="pb-2">Choose a file to upload</Label></div>
  <div class="item item-2"><Toggle checked disabled>Turbex encryption</Toggle></div>
  <div class="item item-3"><Fileupload /></div>
</div>

<style>
  .container {
    display: grid;
    grid-template-columns: 30% 70%;
    grid-template-rows: auto;
    grid-template-areas: 
      "label option"
      "main main";
  }
  .item-1 {
    grid-area: label;
  }
  .item-2 {
    grid-area: option;
  }
  .item-3 {
    grid-area: main;
  }
</style>
