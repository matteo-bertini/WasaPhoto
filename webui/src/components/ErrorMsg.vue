<script>
export default {
    props: ['message'] 
}
</script>
<template>
    <div class="error-container" v-if="message">
        <div class="error-glass">
            <i class="fas fa-exclamation-circle error-icon"></i>
            <span class="error-text">{{ message }}</span>
            <button class="close-btn" @click="$emit('close')">
                <i class="fas fa-times"></i>
            </button>
        </div>
    </div>
</template>



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
    position: relative; /* Necessario per posizionare bene la X */
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
    flex: 1; /* Spinge la X tutto a destra */
}

/* Stile per la X di chiusura */
.close-btn {
    background: none;
    border: none;
    color: rgba(255, 238, 240, 0.6);
    cursor: pointer;
    font-size: 1rem;
    padding: 0 5px;
    transition: color 0.2s;
    display: flex;
    align-items: center;
}

.close-btn:hover {
    color: #ffffff;
}

@keyframes fadeIn {
    from { opacity: 0; transform: translateY(-5px); }
    to { opacity: 1; transform: translateY(0); }
}
</style>