<template>
    <div id="SettingsPageContainer">
        <ErrorMsg v-if="ErrorMessage" :Message="ErrorMessage"></ErrorMsg>

        <div class="settings-wrapper">
            
            <header class="settings-header">
                <button class="btn-glass" @click="BackButtonPressed">
                    <i class="fa-solid fa-arrow-left"></i> Back
                </button>
                <button class="btn-ban small-btn" @click="LogoutButtonPressed">
                    <i class="fa-solid fa-right-from-bracket"></i> Logout
                </button>
            </header>

            <main class="settings-main">
                <h1 class="page-title">Impostazioni</h1>

                <section class="glass-card">
                    <div class="card-header">
                        <i class="fa-solid fa-user-pen"></i>
                        <h3>Modifica Username</h3>
                    </div>
                    <div class="input-group-custom">
                        <input 
                            type="text" 
                            v-model="NewUsername" 
                            :placeholder="'Attuale: ' + Username" 
                            class="custom-input"
                        >
                        <button 
                            class="btn-confirm" 
                            :disabled="isUsernameInvalid" 
                            @click="ConfirmButtonPressed"
                        >
                            Conferma
                        </button>
                    </div>
                    <div class="helper-text">
                        L'username deve avere 3-30 caratteri (solo lettere e numeri).
                    </div>
                </section>

                <section class="glass-card">
                    <div class="flex-row-between">
                        <div class="card-info">
                            <h3>Utenti Bannati</h3>
                            <p class="helper-text">Visualizza e sblocca gli utenti che hai limitato.</p>
                        </div>
                        <button class="btn-utility" data-bs-toggle="modal" data-bs-target="#BannedUsersModal">
                            Gestisci
                        </button>
                    </div>
                </section>

                <section class="glass-card danger-zone">
                    <div class="flex-row-between">
                        <div class="card-info">
                            <h3 class="text-danger-soft">Elimina Account</h3>
                            <p class="helper-text">L'eliminazione del profilo è definitiva e irreversibile.</p>
                        </div>
                        <button class="btn-ban" data-bs-toggle="modal" data-bs-target="#DeleteUserModal">
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
                        <h5 class="modal-title">Lista Nera</h5>
                        <button type="button" class="btn-close btn-close-white" data-bs-dismiss="modal"></button>
                    </div>
                    <div class="modal-body">
                        <div v-if="BannedUsers.length === 0" class="empty-state">
                            Non ci sono utenti bannati.
                        </div>
                        <div v-for="BannedUser in BannedUsers" :key="BannedUser" class="banned-item">
                            <span class="user-tag">@{{BannedUser}}</span>
                            <button class="unban-btn" title="Sblocca" @click="UnbanButtonPressed(BannedUser)">
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
                            <button class="btn-ban-solid" @click="DeleteUserButtonPressed" data-bs-dismiss="modal">Conferma Eliminazione</button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    data() {
        return {
            ErrorMessage: null,
            Username: "",
            NewUsername: "",
            BannedUsers: []
        }
    },
    computed: {
        isUsernameInvalid() {
            const clean = this.NewUsername.trim();
            return clean.length < 3 || clean.length > 30 || clean.includes(" ");
        }
    },
    methods: {
        BackButtonPressed() {
            this.$router.replace("/users/" + this.Username + "/");
        },
        LogoutButtonPressed() {
            localStorage.clear();
            this.$router.replace("/login");
        },
        async ConfirmButtonPressed() {
            const config = {
                headers: { Authorization: `Bearer ${localStorage.getItem("Authstring")}` }
            };
            try {
                const response = await this.$axios.put("/users/" + this.Username + "/username", { Username: this.NewUsername }, config);
                this.Username = response.data.Username;
                this.NewUsername = "";
                localStorage.setItem("Username", this.Username);
                this.$router.replace("/users/" + this.Username + "/settings");
                // Opzionale: aggiungi un feedback di successo qui
            } catch (e) {
                this.NewUsername = "";
                if (e.response && e.response.status === 500) {
                    this.ErrorMessage = "Oops, esiste già un utente con questo Username.";
                }
            }
        },
        async UnbanButtonPressed(BannedId) {
            const config = {
                headers: { Authorization: `Bearer ${localStorage.getItem("Authstring")}` }
            };
            try {
                await this.$axios.delete("/users/" + this.Username + "/bannedusers/" + BannedId, config);
                this.BannedUsers = this.BannedUsers.filter(u => u != BannedId);
            } catch (e) {
                console.error("Errore unban:", e);
            }
        },
        async DeleteUserButtonPressed() {
            const config = {
                headers: { Authorization: `Bearer ${localStorage.getItem("SessionToken")}` }
            };
            try {
                await this.$axios.delete("/users/" + this.Username , config);
                localStorage.clear();
                this.$router.replace("/login");
            } catch (e) {
                console.error("Errore eliminazione:", e);
            }
        }
    },
    async mounted() {
        this.Username =localStorage.getItem("Username");
        const config = {
            headers: { Authorization: `Bearer ${localStorage.getItem("Authstring")}` }
        };
        try {
            const response = await this.$axios.get("/users/" + this.Username + "/bannedusers/", config);
            this.BannedUsers = response.data.BannedUsers.map(x => x.BannedId);
        } catch (e) {
            console.error("Errore caricamento ban:", e);
        }
    }
}
</script>

<style scoped>
#SettingsPageContainer {
    background-color: #0a0a0a; /* Sfondo nero profondo */
    min-height: 100vh;
    width: 100%;
    color: #f8f9fa;
    font-family: 'Inter', sans-serif;
    padding-bottom: 50px;
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

.btn-ban {
    background: rgba(102, 0, 26, 0.2);
    color: #ff4757;
    border: 1px solid rgba(102, 0, 26, 0.4);
    padding: 8px 20px;
    border-radius: 10px;
    cursor: pointer;
    transition: all 0.3s;
}

.btn-ban:hover {
    background: #66001a;
    color: white;
}

.btn-ban-solid {
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
    font-size: 0.85rem;
    color: rgba(255, 255, 255, 0.4);
    margin-top: 10px;
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