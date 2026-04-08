<script>
    export default{
        data(){
            return{
                Username:"",
                PhotoStream: []
            }
        },
        async mounted(){

            // Setto l'Username
            this.Username=this.$route.params.Username;

            // Richiesta del photostream da mostrare
            let getMyStream_config = {
                headers: {
                    Authorization: `Bearer ${localStorage.getItem("Authstring")}`
                }
            };
            try{
                let getMyStream_response = await this.$axios.get("/users/"+this.Username+"/",getMyStream_config);
                this.PhotoStream = getMyStream_response.data.PhotoStream;
                return;

            }catch(e){
                console.log(e);
                return;
            }

            

        },
        methods: {
            // Click sul pulsante Profile
			ProfileButtonPressed(){
				this.$router.replace("/users/"+this.Username+"/");
				return;
			},

			// Click sul pulsante Logout
			LogoutButtonPressed(){

				// Pulizia del localstorage
				localStorage.clear();

				// Ritorno alla schermata di Login
				this.$router.replace("/login");
				return;

			},
        }
    }

</script>
<template>
  <div class="wasaphoto-bg">
    <nav class="glass-nav">
      <div class="nav-content">
        <h2 class="brand-logo">WASA<span>PHOTO</span></h2>
        <div class="nav-actions">
          <button @click="triggerUpload" class="icon-btn" title="Carica">
            <i class="fas fa-plus-circle"></i>
          </button>
          <button @click="goToProfile" class="icon-btn" title="Profilo">
            <i class="fas fa-user-circle"></i>
          </button>
          <button @click="handleLogout" class="icon-btn logout">
            <i class="fas fa-sign-out-alt"></i>
          </button>
        </div>
      </div>
    </nav>

    <main class="feed-container">
      <div v-if="photos.length === 0" class="empty-state">
        <p>Non ci sono ancora foto nel tuo feed. Inizia a seguire qualcuno!</p>
      </div>

      <div v-for="photo in photos" :key="photo.id" class="glass-card photo-card">
        <div class="card-header">
          <div class="user-info">
            <div class="user-avatar">{{ photo.username.charAt(0).toUpperCase() }}</div>
            <span class="username">@{{ photo.username }}</span>
          </div>
          <span class="photo-date">{{ formatDate(photo.createdAt) }}</span>
        </div>

        <div class="photo-wrapper">
          <img :src="photo.url" :alt="photo.description" class="main-photo" />
        </div>

        <div class="card-footer">
          <div class="interactions">
            <button @click="likePhoto(photo.id)" class="like-btn">
              <i class="far fa-heart"></i>
            </button>
            <span class="likes-count">12 likes</span>
          </div>
          <p class="description">
            <strong>{{ photo.username }}</strong> {{ photo.description }}
          </p>
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>/* Sfondo globale coerente col Login */
.wasaphoto-bg {
  min-height: 100vh;
  background: radial-gradient(circle at center, #1a1a1a 0%, #0a0a0a 100%);
  color: white;
  font-family: 'Inter', sans-serif;
}

/* Navbar Glassmorphism */
.glass-nav {
  position: fixed;
  top: 0;
  width: 100%;
  height: 70px;
  background: rgba(255, 255, 255, 0.03);
  backdrop-filter: blur(15px);
  -webkit-backdrop-filter: blur(15px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  z-index: 1000;
}

.nav-content {
  max-width: 800px;
  margin: 0 auto;
  width: 90%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.brand-logo { font-weight: 800; letter-spacing: -1px; }
.brand-logo span { color: #a53b59; }

/* Icone Navbar */
.icon-btn {
  background: none;
  border: none;
  color: #fff;
  font-size: 1.4rem;
  margin-left: 20px;
  cursor: pointer;
  transition: 0.3s;
}

.icon-btn:hover { color: #726fb4; transform: translateY(-2px); }
.logout:hover { color: #ff4757; }

/* Container del Feed */
.feed-container {
  padding: 100px 20px 50px;
  max-width: 600px;
  margin: 0 auto;
}

/* Photo Card Elegante */
.photo-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  margin-bottom: 40px;
  overflow: hidden;
  box-shadow: 0 10px 30px rgba(0,0,0,0.5);
}

.card-header {
  padding: 15px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.user-info { display: flex; align-items: center; gap: 12px; }

.user-avatar {
  width: 35px;
  height: 35px;
  background: linear-gradient(45deg, #a53b59, #726fb4);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 0.9rem;
}

.username { font-weight: 600; font-size: 0.95rem; }
.photo-date { font-size: 0.8rem; color: #888; }

/* Foto */
.photo-wrapper { width: 100%; line-height: 0; }
.main-photo {
  width: 100%;
  max-height: 700px;
  object-fit: cover;
}

/* Footer Card */
.card-footer { padding: 15px 20px; }

.interactions { display: flex; align-items: center; gap: 10px; margin-bottom: 10px; }
.like-btn { background: none; border: none; color: white; font-size: 1.3rem; cursor: pointer; }
.likes-count { font-size: 0.9rem; font-weight: 600; }

.description { font-size: 0.95rem; line-height: 1.4; color: #ddd; }
.description strong { color: white; margin-right: 5px; }

</style>