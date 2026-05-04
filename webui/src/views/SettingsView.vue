<script>
import ErrorMsg from '../components/ErrorMsg.vue';

export default {
    data() {
        return {
            errorMessage: "",
            loading: false,
            username: "",
            newUsername: "",
            bannedUsers: []

        }
    },
    computed: {
        isUsernameInvalid() {
            const regex = /^[a-zA-Z0-9_]{3,30}$/;
            return !regex.test(this.newUsername);
        }
    },
    methods: {
        goBackToMyProfile() {
            this.$router.push(`/users/${this.username}`);

        },
        async doLogout() {
            this.loading = true;
            const token = localStorage.getItem("SessionToken");
            try {
                await this.$axios.delete("/session", {headers: { Authorization: `Bearer ${token}` }});
                localStorage.clear();
                this.loading=false;
                this.$router.push("/login");
            } 
            catch (e) {
                this.loading = false;
                if (e.response) {
                switch (e.response.status) {
                    case 401:
                    localStorage.clear();
                    this.$router.push("/login");
                    break;
                    default:
                    this.errorMessage = "Si è verificato un errore sul server. Riprova più tardi.";
                }
                } else {
                    // Handling network or connectivity issues
                    this.errorMessage = "Impossibile connettersi al server. Controlla la tua connessione.";
                }
            }
        },
        async getBannedUsers(){
            this.loading = true;
            try {
                const targetUsername = localStorage.getItem("Username");
                const token = localStorage.getItem("SessionToken");
                const response = await this.$axios.get(`/users/${targetUsername}/ban`, {headers: { Authorization: `Bearer ${token}` }});
                this.bannedUsers=response.data.map(element => element.username);
                this.loading = false;
            }
            catch (e) {
                this.loading = false;
                if (e.response) {
                    switch (e.response.status) {
                        case 401:
                            localStorage.clear();
                            this.$router.push("/login");
                            break;
                        case 403:
                            this.errorMessage = "Azione non consentita.";
                            break;
                        case 404:
                            this.errorMessage = "Risorsa non trovata.";
                            break;
                        default:
                            this.errorMessage = "Si è verificato un errore sul server. Riprova più tardi.";
                        }
                } else {
                    // Handling network or connectivity issues
                    this.errorMessage = "Impossibile connettersi al server. Controlla la tua connessione.";
                }
            }
            

        },
        async unbanUser(bannedUser) {
            this.loading = true;
            const token = localStorage.getItem("SessionToken");
            try {
                await this.$axios.delete("/users/"+bannedUser+"/ban", {headers: { Authorization: `Bearer ${token}` }});
                this.bannedUsers = this.bannedUsers.filter(user => user !== bannedUser);
                this.loading=false;
            } 
            catch (e) {
                this.loading = false;
                if (e.response) {
                switch (e.response.status) {
                    case 401:
                        localStorage.clear();
                        this.$router.push("/login");
                        break;
                    case 404:
                        this.errorMessage = "Risorsa non trovata."
                    default:
                        this.errorMessage = "Si è verificato un errore sul server. Riprova più tardi.";
                
                    }
                } else {
                    // Handling network or connectivity issues
                    this.errorMessage = "Impossibile connettersi al server. Controlla la tua connessione.";
                }
            }
        },
        async updateUsername(){
            this.loading = true;
            const token = localStorage.getItem("SessionToken");
            try {
                await this.$axios.put(`/users/${this.username}`, { username: this.newUsername }, { headers: { Authorization: `Bearer ${token}` } });           
                this.username = this.newUsername;
                this.newUsername="";
                localStorage.setItem("Username",this.username);
                this.loading=false;
            } 
            catch (e) {
                this.loading = false;
                if (e.response) {
                switch (e.response.status) {
                    case 400:
                        this.errorMessage = "Formato dell'username non valido.";
                        break;
                    case 401:
                        localStorage.clear();
                        this.$router.push("/login");
                        break;
                    case 403:
                        this.errorMessage = "Puoi cambiare solo il tuo username.";
                        break;
                    case 409:
                        this.errorMessage = "Username non disponibile.";
                        return;
                    default:
                    this.errorMessage = "Si è verificato un errore sul server. Riprova più tardi.";
                }
                } else {
                    // Handling network or connectivity issues
                    this.errorMessage = "Impossibile connettersi al server. Controlla la tua connessione.";
                }
            }
            

        },
        async deleteUser() {
            this.loading = true;
            const token = localStorage.getItem("SessionToken");
            try {
                await this.$axios.delete(`/users/${this.username}`, { headers: { Authorization: `Bearer ${token}` } });           
                localStorage.clear();
                this.loading=false;
                this.$router.push("/login");
            } 
            catch (e) {
                this.loading = false;
                if (e.response) {
                    switch (e.response.status) {
                        case 401:
                            localStorage.clear();
                            this.$router.push("/login");
                            break;
                        case 403:
                            this.errorMessage = "Operazione non permessa, puoi eliminare solo il tuo account."
                        case 404:
                            this.errorMessage = "Utente non trovato."
                        default:
                            this.errorMessage = "Si è verificato un errore sul server. Riprova più tardi.";
                        }
                } else {
                    // Handling network or connectivity issues
                    this.errorMessage = "Impossibile connettersi al server. Controlla la tua connessione.";
                }
            }
        }
        
    },
    async mounted() {
        this.username=localStorage.getItem("Username");
        this.bannedUsers = [];
    }
}
</script>

<template>
    <div id="SettingsPageContainer">
      
        <ErrorMsg v-if="errorMessage" :message="errorMessage" @close="errorMessage = ''" />

        <div class="settings-wrapper">
            
            <header class="settings-header">
                
                <button class="btn-glass" @click="goBackToMyProfile">
                    <i class="fa-solid fa-arrow-left"></i> Back
                </button>

                <button :disabled="loading" class="btn-logout small-btn" @click="doLogout">
                    <i class="fa-solid fa-right-from-bracket"></i> Logout
                </button>

            </header>

            <main class="settings-main">
                <h1 class="page-title">Impostazioni</h1>

                <section class="glass-card">
                    <div class="card-header">
                        <i class="fa-solid fa-user-pen"></i>
                        <h3>Modifica username</h3>
                    </div>
                    <div class="input-group-custom">
                        <input 
                            type="text" 
                            v-model="newUsername" 
                            :placeholder="'Attuale: ' + username" 
                            class="custom-input"
                            minlength="3"
                            maxlength="30"
                            pattern="^[a-zA-Z0-9_]{3,30}$"
                            @keydown.space.prevent
                        >
                        <button 
                            class="btn-confirm" 
                            :disabled="isUsernameInvalid" 
                            @click="updateUsername"
                        >
                            Conferma
                        </button>
                    </div>
                    <div class="helper-text" :class="{ 'text-error': newUsername.length > 0 && isUsernameInvalid }">
                        <i v-if="newUsername.length > 0" 
                            :class="isUsernameInvalid ? 'fa-solid fa-circle-xmark' : 'fa-solid fa-circle-check'" class="me-1">
                        </i>
    
                        <span>
                            {{ isUsernameInvalid && newUsername.length > 0 ? 'Username non valido' : 'L\'username deve avere 3-30 caratteri (lettere, numeri e _ )' }}
                        </span>
                    </div>
                    
                </section>

                <section class="glass-card">
                    <div class="flex-row-between">
                        <div class="card-info">
                            <h3>Gestione ban</h3>
                            <p class="helper-text">Visualizza e sblocca gli utenti che hai bannato.</p>
                        </div>
                        <button @click ="getBannedUsers()" class="btn-utility" data-bs-toggle="modal" data-bs-target="#BannedUsersModal">
                            Gestisci
                        </button>
                    </div>
                </section>

                <section class="glass-card danger-zone">
                    <div class="flex-row-between">
                        <div class="card-info">
                            <h3 class="text-danger-soft">Eliminazione account</h3>
                            <p class="helper-text">L'eliminazione è definitiva e irreversibile.</p>
                        </div>
                        <button class="btn-logout" data-bs-toggle="modal" data-bs-target="#DeleteUserModal">
                            Elimina
                        </button>
                    </div>
                </section>
            </main>

        </div>

        <div class="modal fade custom-modal" id="BannedUsersModal" tabindex="-1">
            <div class="modal-dialog modal-dialog-centered modal-dialog-scrollable">
                <div class="modal-content">
                    <div class="modal-header">
                        
                        <h5 class="modal-title">Lista utenti bannati</h5>
                        
                        
                        <button type="button" class="btn-close btn-close-white" data-bs-dismiss="modal"></button>
                    </div>
                    <div class="modal-body">
                        
                        <div v-if="bannedUsers.length === 0" class="empty-state">
                            Non ci sono utenti bannati.
                        </div>

                        <div v-for="bannedUser in bannedUsers" :key="bannedUser" class="banned-item">
                            <span class="user-tag">{{bannedUser}}</span>
                            
                            <button class="unban-btn" title="Sblocca" @click="unbanUser(bannedUser)">
                                <i class="fa-regular fa-trash-can"></i>
                            </button>
                        
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="modal fade custom-modal" id="DeleteUserModal" tabindex="-1">
            <div class="modal-dialog modal-dialog-centered">
                <div class="modal-content border-danger-glow">
                    <div class="modal-body text-center p-5">
                        <i class="fa-solid fa-circle-exclamation danger-icon"></i>
                        <h4 class="mt-3">Sei proprio sicuro?</h4>
                        <p class="mb-4 opacity-75">Tutti i tuoi dati, i post e i follower verranno eliminati istantaneamente.</p>
                        <div class="modal-actions-flex">
                            <button class="btn-glass" data-bs-dismiss="modal">Annulla</button>
                            <button class="btn-logout-solid" @click="deleteUser" data-bs-dismiss="modal">Conferma Eliminazione</button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
#SettingsPageContainer {
    background: radial-gradient(circle at center, #1a1a1a 0%, #0a0a0a 100%);
    color: white; padding-bottom: 50px; font-family: 'Inter', sans-serif;
    min-height: 100vh;
    width: 100%;
}

.settings-wrapper {
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
}

.settings-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px 0 40px 0;
}

.page-title {
    font-size: 2.2rem;
    font-weight: 800;
    margin-bottom: 30px;
    background: linear-gradient(to right, #fff, #666);
    background-clip: unsets;
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
}

/* Glass Cards Style */
.glass-card {
    background: rgba(255, 255, 255, 0.03);
    backdrop-filter: blur(12px);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 20px;
    padding: 25px;
    margin-bottom: 20px;
    transition: all 0.3s ease;
}

.glass-card:hover {
    background: rgba(255, 255, 255, 0.05);
    border-color: rgba(255, 255, 255, 0.15);
}

.card-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 15px;
}

.card-header h3 {
    margin: 0;
    font-size: 1.2rem;
    font-weight: 600;
}

.danger-zone {
    border-left: 4px solid #66001a;
}

/* Form Elements */
.input-group-custom {
    display: flex;
    gap: 10px;
}

.custom-input {
    flex: 1;
    background: rgba(0, 0, 0, 0.3);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    color: white;
    padding: 12px 18px;
    outline: none;
    transition: border 0.3s;
}

.custom-input:focus {
    border-color: #003366;
}

/* Buttons */
.btn-confirm {
    background: #003366;
    color: white;
    border: none;
    border-radius: 12px;
    padding: 0 25px;
    font-weight: 600;
    cursor: pointer;
    box-shadow: 0 4px 15px rgba(0, 51, 102, 0.4);
    transition: all 0.3s;
}

.btn-confirm:hover:not(:disabled) {
    transform: translateY(-2px);
    background: #004080;
}

.btn-confirm:disabled {
    opacity: 0.3;
    cursor: not-allowed;
}

.btn-logout {
    background: rgba(102, 0, 26, 0.2);
    color: #ff4757;
    border: 1px solid rgba(102, 0, 26, 0.4);
    padding: 8px 20px;
    border-radius: 10px;
    cursor: pointer;
    transition: all 0.3s;
}

.btn-logout:hover {
    background: #66001a;
    color: white;
}

.btn-logout-solid {
    background: #66001a;
    color: white;
    border: none;
    padding: 12px 25px;
    border-radius: 12px;
    font-weight: 600;
}

.btn-glass {
    background: rgba(255, 255, 255, 0.05);
    border: 1px solid rgba(255, 255, 255, 0.1);
    color: white;
    padding: 8px 18px;
    border-radius: 10px;
    cursor: pointer;
}

.btn-utility {
    background: #003366;
    color: white;
    border: none;
    padding: 8px 18px;
    border-radius: 10px;
}

/* Modal Customization */
.custom-modal .modal-content {
    background: #121212;
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 24px;
    color: white;
}

.banned-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px;
    background: rgba(255,255,255,0.02);
    border-radius: 12px;
    margin-bottom: 8px;
}

.user-tag {
    font-weight: 500;
    color: #f8f9fa;
}

.unban-btn {
    background: none;
    border: none;
    color: rgba(255,255,255,0.3);
    font-size: 1.1rem;
    cursor: pointer;
    transition: color 0.3s;
}

.unban-btn:hover {
    color: #ff4757;
}

.danger-icon {
    font-size: 3rem;
    color: #66001a;
}

.modal-actions-flex {
    display: flex;
    justify-content: center;
    gap: 15px;
}

.helper-text {
    font-size: 0.95rem;
    color: rgba(255, 255, 255, 0.4); 
    margin-top: 10px;
    transition: color 0.3s ease; 
}

.text-error {
    color: #66001a !important; 
    font-weight: 8000;
}

.flex-row-between {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.empty-state {
    text-align: center;
    padding: 20px;
    opacity: 0.5;
}
</style>