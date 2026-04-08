<template>
	<div class="error-container" v-if="msg">
		<div class="error-glass">
			<i class="fas fa-exclamation-circle error-icon"></i>
			<span class="error-text">{{ msg }}</span>
		</div>
	</div>
</template>

<script>
export default {
	props: ['msg'] 
}
</script>

<style scoped>
.error-container {
	margin-bottom: 20px;
	width: 100%;
}

.error-glass {
	background: rgba(255, 71, 87, 0.1); 
	backdrop-filter: blur(8px);
	-webkit-backdrop-filter: blur(8px);
	border: 1px solid rgba(255, 71, 87, 0.3);
	border-radius: 12px;
	padding: 12px 15px;
	display: flex;
	align-items: center;
	gap: 10px;
	
	animation: fadeIn 0.3s ease-in;
}

.error-icon {
	color: #ff4757;
	font-size: 1.1rem;
}

.error-text {
	color: #ffeef0;
	font-size: 0.9rem;
	font-weight: 500;
}

@keyframes fadeIn {
	from { opacity: 0; transform: translateY(-5px); }
	to { opacity: 1; transform: translateY(0); }
}
</style>