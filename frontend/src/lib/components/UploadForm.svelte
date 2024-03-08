<script lang="ts">
  import { Label, Toggle, Hr, MultiSelect, Tooltip } from 'flowbite-svelte';
	import ConfirmModal from './ConfirmModal.svelte';
  import { Dropzone } from 'flowbite-svelte';

  let value: string[] = [];
  const dropHandle = (event: DragEvent) => {
    value = [];
    event.preventDefault();
    if (!event.dataTransfer) return;
    if (event.dataTransfer.items) {
      [...event.dataTransfer.items].forEach((item, _) => {
        if (item.kind === 'file') {
          const file = item.getAsFile();
          if (file) {
            value.push(file.name);
            value = value;
          }
        }
      });
    } else {
      [...event.dataTransfer.files].forEach((file, _) => {
        value.push(file.name);
        value = value;
      });
    }
  };

  const handleChange = (event: Event) => {
    // The following type was infered based on runtime tests
    const eventTarget = event.target as HTMLInputElement;
    const files = eventTarget.files;
    if (files && files.length > 0) {
      value.push(files[0].name);
      value = value;
    }
  };

  const showFiles = (files: string[]) => {
    if (files.length === 1) return files[0];
    let concat = '';
    files.map((file) => {
      concat += file;
      concat += ',';
      concat += ' ';
    });

    if (concat.length > 40) concat = concat.slice(0, 40);
    concat += '...';
    return concat;
  };

  // Mocks
  const people = [
    { value: 'jc', name: 'John Cena' },
    { value: 'jd', name: 'John Doe' },
    { value: 'mb', name: 'Michel Bolt' }
  ]
</script>

<div class="container px-5">
  <div class="item item-dropzone">
    <Dropzone
      id="dropzone"
      on:drop={dropHandle}
      on:dragover={(event) => {
        event.preventDefault();
      }}
      on:change={handleChange}>
      <svg aria-hidden="true" class="mb-3 w-10 h-10 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" /></svg>
      {#if value.length === 0}
        <p class="mb-2 text-sm text-gray-500 dark:text-gray-400"><span class="font-semibold">Click to upload</span> or drag and drop</p>
        <p class="text-xs text-gray-500 dark:text-gray-400">Anything you want, limited to one file</p>
      {:else}
        <p>{showFiles(value)}</p>
      {/if}
    </Dropzone>
  </div>
  <div class="item item-divider-1"><Hr/></div>
  <div class="item item-people">
    <Label for="people-select">Select people with whom the file will be shared</Label>
    <MultiSelect id="people-select" items={people} />
  </div>
  <div class="item item-divider-2"><Hr/></div>
  <div class="item item-options flex-col space-y-4">
    <Label>Sharing options</Label>
    <Toggle checked disabled>Turbex encryption</Toggle>
    <Tooltip>
      For safety reasons, you cannot turn off Turbex encryption.
    </Tooltip>
    <Toggle checked disabled>Keep a file for me</Toggle>
    <Tooltip>
      Enabling this option will allow you to manage the file remotely.
    </Tooltip>

  </div>
  <div class="item item-upload"><ConfirmModal buttonText="Upload" content="Do you want to upload?"/></div>
</div>

<style>
  .container {
    display: grid;
    grid-template-columns: 70% 30%;
    grid-template-rows: auto;
    grid-template-areas: 
      "dropzone dropzone"
      "divider-1 divider-1"
      "people people"
      "divider-2 divider-2"
      "options upload";
  }
  .item-dropzone {
    grid-area: dropzone;
  }
  .item-divider-1 {
    grid-area: divider-1;
  }
  .item-divider-2 {
    grid-area: divider-2;
  }
  .item-options {
    grid-area: options;
  }
  .item-upload {
    grid-area: upload;
  }
  .item-people {
    grid-area: people;
  }
</style>
