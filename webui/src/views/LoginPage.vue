<script>
import ErrorMsg from '../components/ErrorMsg.vue';

export default {
  components: {
    ErrorMsg
  },
	data() {
		return {
			Username: "",   
			Password: "",   
			IsSignup: false,
      errorMessage: null,
      loading: false,
      isSuccess: false
		}
	},
	methods: {
		async handleAuth() { //  metodo chiamato dal @submit.prevent
      this.loading = true;
			this.errorMessage = null;
      this.isSuccess = false;
			
			try {
				let response = await this.$axios.post("/session", {
					Username: this.Username,
					Password: this.Password,
					IsSignup: this.IsSignup
				});

				// Saving session data.
				localStorage.setItem("SessionToken", response.data.SessionToken);
				localStorage.setItem("Username", this.Username);

        this.isSuccess = true;
        this.loading = false;
    
        setTimeout(() => {
          this.$router.push("/users/"+this.Username);
        }, 1200);

			} catch (e) {
        this.loading = false;
        this.isSuccess = false;
        if (e.response) {
          if (e.response.status === 409) {
            this.errorMessage = "Username già esistente!";
          } 
          else if (e.response.status === 401) {
            this.errorMessage = "Credenziali non valide.";
          } 
          else if (e.response.status === 400) {
            this.errorMessage = "Dati non validi. Controlla i campi inseriti.";
          } 
          else {
            this.errorMessage = "Si è verificato un errore sul server. Riprova.";
          }
        } 
        else {
          this.errorMessage = "Impossibile collegarsi al server. Controlla la tua connessione.";
        }
    }
}
		}
	}

</script>

<template>
  <div class="wasaphoto-bg">
    
    <header class="main-header">
      <div class="logo-container">
        <i class="fas fa-camera camera-icon"></i>
        <h1 class="brand-name">WASAPHOTO</h1>
      </div>
    </header>

    <div class="glass-card-container">
      <div class="glass-card">
        
        <h2 class="card-title">{{ isSignup ? 'Registrazione' : 'Login' }}</h2>
        <ErrorMsg :msg="errorMessage" />

        <form @submit.prevent="handleAuth">
          
          <div class="input-group">
            <input 
              type="text" 
              v-model="Username" 
              placeholder="Username" 
              required
              class="glass-input"
            />
          </div>

          <div class="input-group">
            <input 
              type="password" 
              v-model="Password" 
              placeholder="Password" 
              required
              class="glass-input"
            />
          </div>

          <div class="switch-container">
            <span :class="{ 'active-label': !IsSignup }">Accedi</span>
            <label class="switch">
              <input type="checkbox" v-model="IsSignup">
              <span class="slider round"></span>
            </label>
            <span :class="{ 'active-label': IsSignup }">Registrati</span>
          </div>

          <button type="submit" class="gradient-button" :disabled="loading || isSuccess" :class="{ 'btn-success': isSuccess }">
            <template v-if="loading">
              <i class="fas fa-circle-notch fa-spin"></i>
              <span>VERIFICA...</span>
            </template>
            
            <template v-else-if="isSuccess">
              <i class="fas fa-check-circle anim-pop"></i>
              <span>BENVENUTO!</span>
            </template>

            <template v-else>
              <i :class="IsSignup ? 'fas fa-user-plus' : 'fas fa-sign-in-alt'"></i>
              <span>{{ IsSignup ? 'CREA ACCOUNT' : 'ACCEDI' }}</span>
            </template>
          </button>
        </form>

        <div class="card-footer">
          <a href="#" class="forgot-password">Password dimenticata?</a>
        </div>
      </div>
    </div>

    <i class="fas fa-star decorative-star"></i>
  </div>
</template>

<style scoped>
.wasaphoto-bg {
  min-height: 100vh;
  background:radial-gradient(circle at center, #1a1a1a 0%, #0a0a0a 100%);
  background-size: cover;
  background-position: center;
  background-attachment: fixed;
  display: flex;
  flex-direction: column;
  font-family: 'Poppins', sans-serif; 
  color: white;
}


.main-header {
  padding: 40px 0;
  display: flex;
  justify-content: center;
}

.logo-container {
  display: flex;
  align-items: center;
  gap: 15px;
}

.camera-icon {
  font-size: 2.5rem;
  color: #f0e6d2; 
}

.brand-name {
  font-size: 2.2rem;
  font-weight: 700;
  letter-spacing: 2px;
  margin: 0;
  color: white;
}


.glass-card-container {
  flex-grow: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  padding-bottom: 100px; 
}

.glass-card {
  /* Trasparenza e Sfumatura sfondo */
  background: rgba(255, 255, 255, 0.05);
  
  /* SFOCATURA DELLO SFONDO (Backdrop Filter) - Fondamentale! */
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  
  /* Bordo sottile e luminoso per definire la forma */
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  

  box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.37);
  
  width: 100%;
  max-width: 400px;
  padding: 40px;
  text-align: center;
}

.card-title {
  font-size: 1.8rem;
  margin-bottom: 30px;
  font-weight: 600;
  color: white;
}

/* 4. INPUT STILIZZATI */
.input-group {
  margin-bottom: 20px;
}

.glass-input {
  width: 100%;
  padding: 15px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  color: white;
  font-size: 1rem;
  transition: all 0.3s ease;
}

.glass-input::placeholder {
  color: rgba(255, 255, 255, 0.5);
}

.glass-input:focus {
  outline: none;
  background: rgba(255, 255, 255, 0.15);
  border-color: rgba(255, 255, 255, 0.3);
  box-shadow: 0 0 10px rgba(255, 255, 255, 0.1);
}

.switch-container {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 15px;
  margin: 25px 0;
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.9rem;
}

.active-label {
  color: white;
  font-weight: 600;
}

.switch {
  position: relative;
  display: inline-block;
  width: 50px;
  height: 26px;
}

.switch input { opacity: 0; width: 0; height: 0; }

.slider {
  position: absolute;
  cursor: pointer;
  top: 0; left: 0; right: 0; bottom: 0;
  background-color: rgba(255, 255, 255, 0.1);
  transition: .4s;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.slider:before {
  position: absolute;
  content: "";
  height: 18px; width: 18px;
  left: 4px; bottom: 3px;
  background-color: white;
  transition: .4s;
}

input:checked + .slider {
  background-color: #a53b59; 
}

input:focus + .slider {
  box-shadow: 0 0 1px #a53b59;
}

input:checked + .slider:before {
  transform: translateX(24px);
}

.slider.round { border-radius: 34px; }
.slider.round:before { border-radius: 50%; }


.gradient-button {
  width: 100%;
  padding: 15px;
  border: none;
  border-radius: 10px;
  background: linear-gradient(135deg, #a53b59 0%, #726fb4 100%);
  color: white;
  font-size: 1.1rem;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
}

.gradient-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(165, 59, 89, 0.4);
}

.gradient-button:active {
  transform: translateY(1px);
}

/* 7. FOOTER CARD E DECORAZIONI */
.card-footer {
  margin-top: 25px;
}

.forgot-password {
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.85rem;
  text-decoration: none;
  transition: color 0.3s ease;
}

.forgot-password:hover {
  color: white;
  text-decoration: underline;
}

.decorative-star {
  position: fixed;
  bottom: 30px;
  right: 30px;
  font-size: 2rem;
  color: rgba(255, 255, 255, 0.3);

}

/* Feedback di Successo */
.btn-success {
  background: linear-gradient(135deg, #2ed573 0%, #7bed9f 100%) !important;
  box-shadow: 0 5px 15px rgba(46, 213, 115, 0.4) !important;
}

.anim-pop {
  animation: pop 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

@keyframes pop {
  0% { transform: scale(0.5); opacity: 0; }
  100% { transform: scale(1.2); opacity: 1; }
}

/* Stato disabilitato per evitare click multipli */
.gradient-button:disabled {
  cursor: not-allowed;
  opacity: 0.8;
  transform: none;
}

/* Rotazione icona loading */
.fa-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  100% { transform: rotate(360deg); }
}
</style>