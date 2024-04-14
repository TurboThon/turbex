<script lang="ts">
  import { Table, TableBody, TableBodyCell, TableBodyRow, TableHead, TableHeadCell, Button } from 'flowbite-svelte';
  import { DownloadSolid, TrashBinSolid } from 'flowbite-svelte-icons';
  import { type File } from '../../types/file';

  const doDownload = () => {
    console.log("downloading");
  };
  
  const doDelete = () => {
    console.log("deleting");
  };

  export let files: File[];
</script>

<Table>
  <TableHead>
    <TableHeadCell class="grow">File name</TableHeadCell>
    <TableHeadCell class="grow-0">Owner</TableHeadCell>
    <TableHeadCell class="grow-0">Shared date</TableHeadCell>
    <TableHeadCell class="grow-0">Actions</TableHeadCell>
  </TableHead>
  <TableBody>
    {#each files as file}
    <TableBodyRow>
      <TableBodyCell class="grow">{file.name}</TableBodyCell>
      <TableBodyCell class="grow-0">{file.owner}</TableBodyCell>
      <TableBodyCell class="grow-0">{file.date}</TableBodyCell>
      <TableBodyCell class="grow-0 flex space-x-4">
        <Button on:click={doDownload} size="xs" color="green"><DownloadSolid size="md"/></Button>
        <Button on:click={doDelete} size="xs"><TrashBinSolid size="md"/></Button>
      </TableBodyCell>
    </TableBodyRow>
    {/each}
  </TableBody>
</Table>

