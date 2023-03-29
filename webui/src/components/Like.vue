<script>
export default {
	props: ['LikeUsername'],
	
		
}

</script>

<template>
	<div class="LikeUsername" >
		{{LikeUsername}}
	</div>

</template>

<style>
	.LikeUsername:hover{
		transform:scale(1.1);
		color:#007bff;
		cursor: pointer;

	}
	
</style>
