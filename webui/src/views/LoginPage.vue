<script>
export default {
	data() {
		return {
			ErrorMessage: null,
			username: "",   
			password: "",   
			isSignup: false 
		}
	},
	methods: {
		async handleAuth() { //  metodo chiamato dal @submit.prevent
			this.ErrorMessage = null;
			
			try {
				let response = await this.$axios.post("/session", {
					Username: this.username,
					Password: this.password,
					IsSignup: this.isSignup
				});

				// Salvataggio dati sessione
				localStorage.setItem("SessionToken", response.data.SessionToken);
				localStorage.setItem("Username", this.username);

				// Reindirizzamento al profilo dell'utente
				this.$router.push("/users/" + this.username + "/");

			} catch (e) {
				console.error("Errore durante l'autenticazione:", e);
				
				// Gestione errori specifica basata sugli status code del backend
				if (e.response && e.response.status === 409) {
					this.ErrorMessage = "Username già esistente!";
				} else if (e.response && e.response.status === 401) {
					this.ErrorMessage = "Credenziali non valide.";
				} else {
					this.ErrorMessage = "Errore del server. Riprova più tardi.";
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

        <form @submit.prevent="handleAuth">
          
          <div class="input-group">
            <input 
              type="text" 
              v-model="username" 
              placeholder="username" 
              required
              class="glass-input"
            />
          </div>

          <div class="input-group">
            <input 
              type="password" 
              v-model="password" 
              placeholder="password" 
              required
              class="glass-input"
            />
          </div>

          <div class="switch-container">
            <span :class="{ 'active-label': !isSignup }">Accedi</span>
            <label class="switch">
              <input type="checkbox" v-model="isSignup">
              <span class="slider round"></span>
            </label>
            <span :class="{ 'active-label': isSignup }">Registrati</span>
          </div>

          <button type="submit" class="gradient-button">
            <i :class="isSignup ? 'fas fa-user-plus' : 'fas fa-sign-in-alt'"></i>
            {{ isSignup ? 'CREA ACCOUNT' : 'ACCEDI' }}
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
  background-image: url('/path/to/my/camera-lens-bg.jpg'); 
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
  padding-bottom: 100px; /* Spazio per la stellina */
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
</style>