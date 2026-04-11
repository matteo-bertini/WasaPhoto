<template>
  <div class="wasaphoto-bg">
    <header class="glass-header">
      <div class="header-content">
        <h1 class="logo" @click="$router.push('/home')">WASA<span>PHOTO</span></h1>
        
        <div class="search-bar">
          <i class="fa-solid fa-magnifying-glass"></i>
          <input type="text" v-model="searchQuery" placeholder="Cerca utenti..." @keyup.enter="handleSearch">
        </div>

        <div class="nav-actions">
          <button class="nav-icon" @click="$router.push('/upload')" title="Carica foto"><i class="fa-solid fa-circle-plus"></i></button>
          <button class="nav-icon" @click="$router.push('/profile')" title="Il mio profilo"><i class="fa-solid fa-user"></i></button>
          <button class="nav-icon logout" @click="logout" title="Esci"><i class="fa-solid fa-right-from-bracket"></i></button>
        </div>
      </div>
    </header>

    <main class="feed">
      <div v-if="loading" class="status-msg">Caricamento photostream...</div>
      
      <PhotoComponent 
        v-for="photo in photoStream" 
        :key="photo.PhotoStreamComponentPhotoId"
        :owner="photo.PhotoStreamComponentUsername"
        :photoid="photo.PhotoStreamComponentPhotoId"
        :likesnumber="photo.PhotoStreamComponentLikesNumber"
        :commentsnumber="photo.PhotoStreamComponentCommentsNumber"
        :dateofupload="photo.PhotoStreamComponentDateOfUpload"
        :mockImageUrl="photo.mockUrl" @photo_deleted_from_database="handleDeletion"
      />

      <div v-if="!loading && photoStream.length === 0" class="status-msg">
        Nessuna foto da mostrare. Inizia a seguire qualcuno!
      </div>
    </main>
  </div>
</template>

<script>
import PhotoComponent from '../components/PhotoComponent.vue';
// import axios from 'axios'; // Commentato per il test senza backend

export default {
  components: { PhotoComponent },
  data() {
    return {
      searchQuery: "",
      photoStream: [], // Inizialmente vuoto, lo popoliamo nel mounted
      loading: true
    }
  },
  mounted() {
    // Invece di chiamare il backend, carichiamo i dati mock
    this.loadMockStream();
  },
  methods: {
    // FUNZIONE TEMPORANEA PER IL TEST
    loadMockStream() {
      // Simuliamo un leggero ritardo di rete (500ms)
      setTimeout(() => {
        this.photoStream = [
          {
            PhotoStreamComponentUsername: "Eren_Yeager",
            PhotoStreamComponentPhotoId: "photo_1",
            PhotoStreamComponentLikesNumber: 125,
            PhotoStreamComponentCommentsNumber: 18,
            PhotoStreamComponentDateOfUpload: "2024-05-15T10:30:00Z",
            // Aggiungiamo un URL reale per il test visivo
            mockUrl: "https://images.unsplash.com/photo-1518770660439-4636190af475?q=80&w=600"
          },
          {
            PhotoStreamComponentUsername: "Mikasa_Ackermann",
            PhotoStreamComponentPhotoId: "photo_2",
            PhotoStreamComponentLikesNumber: 242,
            PhotoStreamComponentCommentsNumber: 35,
            PhotoStreamComponentDateOfUpload: "2024-05-14T15:45:00Z",
            mockUrl: "https://images.unsplash.com/photo-1563207151-5813920f0190?q=80&w=600"
          },
          {
            PhotoStreamComponentUsername: "Armin_Arlert",
            PhotoStreamComponentPhotoId: "photo_3",
            PhotoStreamComponentLikesNumber: 95,
            PhotoStreamComponentCommentsNumber: 12,
            PhotoStreamComponentDateOfUpload: "2024-05-13T09:12:00Z",
            mockUrl: "https://images.unsplash.com/photo-1498050108023-c5249f4df085?q=80&w=600"
          },
          {
            PhotoStreamComponentUsername: "Levi_Ackermann",
            PhotoStreamComponentPhotoId: "photo_4",
            PhotoStreamComponentLikesNumber: 512,
            PhotoStreamComponentCommentsNumber: 64,
            PhotoStreamComponentDateOfUpload: "2024-05-12T18:20:00Z",
            mockUrl: "https://images.unsplash.com/photo-1531297484001-80022131f5a1?q=80&w=600"
          }
        ];
        this.loading = false;
      }, 500); // Ritardo simulato
    },

    // La funzione reale è commentata per ora
    /*
    async loadStream() {
      try {
        const user = localStorage.getItem('Username');
        const response = await axios.get(`/users/${user}/`);
        this.photoStream = response.data.PhotoStream || [];
      } catch (e) {
        console.error("Errore stream", e);
      } finally {
        this.loading = false;
      }
    },
    */
    handleSearch() {
      if (this.searchQuery.trim()) {
        alert("Ricerca per: " + this.searchQuery); // Placeholder
        // this.$router.push(`/search?q=${this.searchQuery}`);
      }
    },
    handleDeletion(id) {
      this.photoStream = this.photoStream.filter(p => p.PhotoStreamComponentPhotoId !== id);
    },
    logout() {
      // localStorage.clear();
      alert("Logout effettuato");
      this.$router.push('/');
    }
  }
}
</script>

<style scoped>
/* Lo stile rimane invariato, lo riporto per completezza */
.wasaphoto-bg {
  min-height: 100vh;
  background: radial-gradient(circle at center, #1a1a1a 0%, #0a0a0a 100%);
  padding-top: 100px; /* Spazio per la header fissa */
}

.glass-header {
  position: fixed;
  top: 0;
  width: 100%;
  height: 80px;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  z-index: 1000;
  display: flex;
  align-items: center;
}

.header-content {
  max-width: 1000px;
  margin: 0 auto;
  width: 90%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo { color: white; cursor: pointer; font-weight: 800; font-size: 1.4rem; }
.logo span { color: #00d2ff; }

.search-bar {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 8px 15px;
  display: flex;
  align-items: center;
  width: 40%;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.search-bar input {
  background: none;
  border: none;
  color: white;
  margin-left: 10px;
  width: 100%;
}

.search-bar input:focus { outline: none; }

.nav-icon {
  background: none;
  border: none;
  color: white;
  font-size: 1.5rem;
  margin-left: 20px;
  cursor: pointer;
  transition: 0.3s;
}

.nav-icon:hover { color: #00d2ff; transform: scale(1.1); }
.logout:hover { color: #ff4757; }

.feed {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-bottom: 50px;
}

.status-msg { color: #888; margin-top: 40px; }
</style>