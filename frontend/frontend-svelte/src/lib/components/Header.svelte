<script lang="ts">
	import { page } from "$app/stores";
	import {
		Navbar,
		NavBrand,
		NavLi,
		NavUl,
		NavHamburger,
		Avatar,
		Dropdown,
		DropdownHeader,
		DropdownItem,
		DropdownDivider,
		Tooltip,
	} from "flowbite-svelte";
	import { handleExpiredSession, userStore } from "$lib/store";
	import { getLogout } from "$lib/query";

	$: activeUrl = $page.url.pathname;

	const handleLogout = () => {
        getLogout()()
        handleExpiredSession()
    };
</script>

<Navbar>
	<NavBrand href="/">
		<p1 class="text-xl"><b>Turbex</b> - Turbo Extreme file sharing system</p1>
	</NavBrand>
	<NavHamburger />
	<NavUl {activeUrl}>
		<NavLi href="/">My files</NavLi>
		<NavLi href="/upload">Upload a file</NavLi>
		<NavLi href="/admin">Admin panel</NavLi>
	</NavUl>
	<div>
		<Avatar id="avatar-menu" style="cursor: pointer;" />
	</div>
	<Dropdown placement="bottom" triggeredBy="#avatar-menu">
		<DropdownHeader>
			<span class="block text-sm">{$userStore?.username ?? "Anonymous"}</span>
		</DropdownHeader>
		<DropdownItem class="not-implemented">Rotate keys</DropdownItem>
		<DropdownItem class="not-implemented">Change Password</DropdownItem>
		<DropdownDivider />
		<DropdownItem on:click={handleLogout}>Log out</DropdownItem>
		<Tooltip triggeredBy=".not-implemented">Feature soon to be implemented!</Tooltip>
	</Dropdown>
</Navbar>
