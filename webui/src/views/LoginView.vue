<script>
import ErrorMsg from '../components/ErrorMsg.vue';

export default {
   name: 'LoginView',
  components: {
    ErrorMsg
  },
  data() {
    return {
      // User input fields
      Username: "",   
      Password: "",   
      IsSignup: false,
      
      // UI State management
      errorMessage: null,
      passwordVisible: false,
      loading: false,
      isSuccess: false
    }
  },
  computed: {
    usernameError() {
      return !(this.Username.length >= 3 && this.Username.length <= 30);
    },
    passwordError() {
      return !(this.Password.length >= 8 && this.Password.length <= 30);
    },
    isFormInvalid() {
      return this.Username.length < 3 || this.Username.length > 30 || this.Password.length < 8 || this.Password.length >30;
    }
  },
  methods: {
    /**
     * Handles the authentication process (Login or Signup).
     * Invoked by the @submit.prevent event on the form.
     */
    async handleAuth() {
      // Reset UI state before starting the request
      this.loading = true;
      this.errorMessage = null;
      this.isSuccess = false;
      
      try {
        // API call to the backend session endpoint
        let response = await this.$axios.post("/session", {
          username: this.Username, 
          password: this.Password,
          isSignUp: this.IsSignup
        });

        /* Session Persistence:
           Store the token and username in localStorage for subsequent 
           authenticated requests 
        */
        localStorage.setItem("SessionToken", response.data.sessionToken);
        localStorage.setItem("Username", this.Username);

        this.isSuccess = true;
        this.loading = false;
    
        // Graceful redirect after a short delay to allow the user to see the success state
        setTimeout(() => {
          this.$router.push("/users/" + this.Username);
        },800);

      } catch (e) {
        this.loading = false;
        this.isSuccess = false;

        // Error handling based on HTTP status codes defined in OpenAPI specs
        if (e.response) {
          switch (e.response.status) {
            case 409:
              this.errorMessage = "Username già esistente!";
              break;
            case 401:
              this.errorMessage = "Credenziali non valide.";
              break;
            case 400:
              this.errorMessage = "Richiesta non valida. Controlla i dati inseriti.";
              break;
            default:
              this.errorMessage = "Si è verificato un errore sul server. Riprova più tardi.";
          }
        } else {
          // Handling network or connectivity issues
          this.errorMessage = "Impossibile connettersi al server. Controlla la tua connessione.";
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
        
        <h2 class="card-title">{{ IsSignup ? 'Registrazione' : 'Login' }}</h2>
        
            <ErrorMsg v-if="errorMessage" :message="errorMessage" @close="errorMessage = ''" />


        <form @submit.prevent="handleAuth">
          
          <div class="input-group">
            <input 
              type="text" 
              v-model="Username" 
              maxlength="30"
              placeholder="Username" 
              required
              class="glass-input"
            />
            
          </div>
          <p v-if="usernameError" class="input-hint">Username 3-30 caratteri</p>

          <div class="input-group">
            <div class="password-wrapper">
              <input 
              :type="passwordVisible ? 'text' : 'password'"
              v-model="Password" 
              maxlength="72"
              placeholder="Password" 
              required
              class="glass-input"
              />
              <span class="password-toggle" @click="passwordVisible = !passwordVisible">
                <i :class="passwordVisible ? 'fas fa-eye' : 'fas fa-eye-slash'"></i>
              </span>
            </div>
          </div>
          <p v-if="passwordError" class="input-hint">Minimo 8 caratteri</p>

          <div class="switch-container">
            <span :class="{ 'active-label': !IsSignup }">Accedi</span>
            <label class="switch">
              <input type="checkbox" v-model="IsSignup">
              <span class="slider round"></span>
            </label>
            <span :class="{ 'active-label': IsSignup }">Registrati</span>
          </div>

          <button type="submit" class="gradient-button" :disabled="loading || isSuccess || isFormInvalid" :class="{ 'btn-success': isSuccess,'btn-invalid': isFormInvalid && !loading && !isSuccess }">
            
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


      </div>
    </div>

  </div>
</template>

<style scoped>

/* --- BASE LAYOUT --- */
.wasaphoto-bg {
  min-height: 100vh;
  background: radial-gradient(circle at center, #1a1a1a 0%, #0a0a0a 100%);
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
  color: white; 
}

.brand-name {
  font-size: 2.2rem;
  font-weight: 700;
  letter-spacing: 2px;
  margin: 0;
  color: white;
}

/* --- CARD STRUCTURE --- */
.glass-card-container {
  flex-grow: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  padding-bottom: 100px; 
}

.glass-card {
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
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

/* --- INPUTS --- */
.input-group {
  margin-bottom: 20px;
  display: flex;
  flex-direction: column;
  align-items: flex-start; 
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

.input-hint {
    margin-left: 15px; 
    margin-top: 6px;
    font-size: 0.75rem;
    color: rgba(255, 255, 255, 0.6); 
    text-align: left;
}
.password-wrapper {
  position: relative;
  width: 100%;
  display: flex;
  align-items: center;
}

.password-toggle {
  position: absolute;
  right: 15px; 
  cursor: pointer;
  color: rgba(255, 255, 255, 0.5);
  transition: color 0.3s ease;
  z-index: 10;
  padding: 5px; 
}

.password-toggle:hover {
  color: white;
}


.password-wrapper .glass-input {
  padding-right: 45px;
}

/* --- TOGGLE SWITCH --- */
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

input:checked + .slider { background-color: #003366; }
input:checked + .slider:before { transform: translateX(24px); }
.slider.round { border-radius: 34px; }
.slider.round:before { border-radius: 50%; }

/* --- BUTTONS & STATES --- */
.gradient-button {
  width: 100%;
  padding: 15px;
  border: none;
  border-radius: 10px;
  background: linear-gradient(135deg,#003366 100%);  color: white;
  font-size: 1.1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
}

.gradient-button:hover:not(:disabled) {
  transform: translateY(-2px);
box-shadow: 0 8px 15px rgba(0, 51, 102, 0.4);}

.gradient-button:active:not(:disabled) {
  transform: translateY(1px);
}

.gradient-button:disabled {
  cursor: not-allowed;
  transition: none; 
}


.gradient-button.btn-invalid:disabled {
  background: #535353 !important;
  filter: grayscale(1) !important;
  opacity: 0.5 !important;
  box-shadow: none !important;
  transform: none !important;
}


.gradient-button.btn-success:disabled {
  background: #2fab3966 !important;
  filter: grayscale(0) !important;
  opacity: 1 !important;
  box-shadow: 0 5px 15px rgba(47, 171, 57, 0.4) !important;
  transform: none !important;
}

/* --- ANIMATIONS --- */
.anim-pop {
  animation: pop 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

@keyframes pop {
  0% { transform: scale(0.5); opacity: 0; }
  100% { transform: scale(1.2); opacity: 1; }
}

.fa-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  100% { transform: rotate(360deg); }
}
</style>