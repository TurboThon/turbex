<script lang="ts">
	import FileTable from '$lib/components/FileTable.svelte';
	import type { File } from '../../types/file';

	let files: File[] = [
		{ id: 1, name: 'Mock file 1.pdf', owner: 'John Doe', date: '2024-03-05' },
		{ id: 2, name: 'Mocked.md', owner: 'Satoshi Nakamoto', date: '2024-02-20' },
		{
			id: 3,
			name: 'Lord of the rings - Return of the King - xvdi_with_audio_with_subtitles_fr_FR-the-pirate-bay.mp4',
			owner: 'Not Me',
			date: '2023-07-28'
		}
	];
</script>

<FileTable {files} />
