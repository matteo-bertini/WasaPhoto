<script>
export default {
	props: ['LikesList'],
}
</script>

<template>
	<div v-if="this.LikesList.length == 0">
		<p>Ancora nessun mi piace</p>
		
	</div>
	<div v-else>
		<p v-for="Like in LikesList" :key ="Like">{{Like}}</p>


	</div>
</template>

<style>
</style>
