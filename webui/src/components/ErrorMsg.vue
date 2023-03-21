<script>
export default {
	props: ["Title","Message"]
}
</script>

<template>
	<div class="alert alert-danger alert-dismissible fade show" role="alert">
		<strong>
			<i class="fa-solid fa-triangle-exclamation"></i>
		</strong>


		<strong>{{Title}}</strong> {{Message}}
		<button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
	 </div>
</template>

<style>
</style>
