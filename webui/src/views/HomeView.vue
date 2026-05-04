<script>
/**
 * @file HomeView.vue
 * @description Main feed component displaying followed and suggested posts with infinite scroll.
 */

import ErrorMsg from '../components/ErrorMsg.vue';
import PostCard from '../components/PostCard.vue';

export default {
  name: 'HomeView',
  components: {
    PostCard,
    ErrorMsg
  },
  data() {
    return {
      /** @type {Array} List of posts in the stream */
      posts: [],
      
      /** @type {boolean} State for initial loading */
      loading: false,
      /** @type {boolean} State for infinite scroll loading */
      loadingMore: false,
      /** @type {boolean} Flag indicating if all posts have been fetched */
      allLoaded: false,
      /** @type {number} Offset for pagination */
      offset: 0,
      /** @type {number} Number of posts to fetch per request */
      limit: 10, 
      
      /** @type {string} Global error message to display in the UI */
      errorMessage: "",
      /** @type {string} Query for the user search input */
      searchQuery: "",
      
      /** @type {boolean} State management for the post upload modal */
      showUploadModal: false,
      /** @type {boolean} State for post upload progress */
      uploadLoading: false,
      /** @type {File|null} The selected image file for upload */
      selectedFile: null,
      /** @type {string|null} Local URL for image preview */
      uploadPreview: null,
      /** @type {string} Caption text for the new post */
      uploadCaption: ""
    };
  },
  methods: {
    /**
     * Fetches the hybrid stream from the API (Following + Suggested).
     * @param {boolean} isInitial - If true, resets the stream and starts from offset 0.
     */
    async fetchStream(isInitial = false) {
      if (this.loading || (this.allLoaded && !isInitial)) return;

      if (isInitial) {
        this.loading = true;
        this.offset = 0;
        this.allLoaded = false;
        this.posts = [];
      } else {
        this.loadingMore = true;
      }

      try {
        const token = localStorage.getItem("SessionToken");
        const response = await this.$axios.get('/stream', {
          headers: { Authorization: `Bearer ${token}` },
          params: {
            limit: this.limit,
            offset: this.offset
          }
        });

        const newPosts = response.data || [];
        
        if (newPosts.length < this.limit) {
          this.allLoaded = true;
        }

        if (isInitial) {
          this.posts = newPosts;
        } else {
          this.posts = [...this.posts, ...newPosts];
        }

        this.offset += this.limit;
      } catch (e) {
        this.handleApiError(e);
      } finally {
        this.loading = false;
        this.loadingMore = false;
      }
    },

    /**
     * Detects when user scrolls near the bottom of the page to trigger infinite scroll.
     */
    handleScroll() {
      const scrollHeight = document.documentElement.scrollHeight;
      const scrollTop = document.documentElement.scrollTop;
      const clientHeight = document.documentElement.clientHeight;

      if (scrollTop + clientHeight >= scrollHeight - 200) {
        if (!this.loading && !this.loadingMore && !this.allLoaded) {
          this.fetchStream();
        }
      }
    },

    /**
     * Shared error handling for stream API calls.
     * @param {Object} e - The error object caught from the Axios request.
     */
    handleApiError(e) {
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
        this.errorMessage = "Impossibile connettersi al server. Controlla la tua connessione.";
      }
    },

    /**
     * Redirects to a specified user profile.
     * @param {string} username - Target username or "self" for the logged-in user.
     */
    goToUserProfile(username) {
      const target = username === "self" ? localStorage.getItem("Username") : username;
      this.$router.push(`/users/${target}`);
    },

    /**
     * Redirects to the settings page.
     */
    goToSettingsPage() {
      this.$router.push("/settings");
    },

    /**
     * Clears session token and redirects to login page.
     */
    async doLogout() {
      const token = localStorage.getItem("SessionToken");
      try {
        await this.$axios.delete("/session", { headers: { Authorization: `Bearer ${token}` } });
      } catch (e) { 
        console.error("Logout error", e); 
      }
      localStorage.clear();
      this.$router.push("/login");
    },

    /**
     * Handles the user search input and redirects to the target profile.
     */
    handleSearch() {
      if (this.searchQuery.trim()) {
        this.$router.push(`/users/${this.searchQuery.trim()}`);
      }
    },

    /**
     * Opens the post upload modal.
     */
    openUploadModal() { 
      this.showUploadModal = true; 
    },

    /**
     * Closes the post upload modal and resets local state.
     */
    closeUploadModal() {
      if (this.uploadPreview) URL.revokeObjectURL(this.uploadPreview);
      this.showUploadModal = false;
      this.selectedFile = null;
      this.uploadPreview = null;
      this.uploadCaption = "";
    },

    /**
     * Handles file selection for post upload, generating a local preview.
     * @param {Event} event - The file input change event.
     */
    handleFileSelect(event) {
      const file = event.target.files[0];
      if (file) {
        this.selectedFile = file;
        this.uploadPreview = URL.createObjectURL(file);
      }
    },

    /**
     * Submits new post data to the server using FormData.
     */
    async submitUpload() {
      if (!this.selectedFile) return;
      this.uploadLoading = true;
      const formData = new FormData();
      formData.append("file", this.selectedFile);
      formData.append("caption", this.uploadCaption);

      try {
        const username = localStorage.getItem("Username");
        const token = localStorage.getItem("SessionToken");
        const response = await this.$axios.post(`/users/${username}/posts`, formData, {
          headers: { Authorization: `Bearer ${token}`, 'Content-Type': 'multipart/form-data' }
        });
        
        this.posts.unshift({ ...response.data, isSuggested: false });
        this.closeUploadModal();
      } catch (e) {
        alert("Caricamento fallito. Riprova.");
      } finally { 
        this.uploadLoading = false; 
      }
    }
  },
  mounted() {
    this.fetchStream(true);
    window.addEventListener('scroll', this.handleScroll);
  },
  unmounted() {
    window.removeEventListener('scroll', this.handleScroll);
  }
};
</script>

<template>
  <div id="HomeViewContainer">
    
    <nav class="glass-nav">
      <div class="nav-content">
        <div  @click="$router.push('/home')" class="logo-container">
          <i class="fas fa-camera camera-icon"></i>
          <h2 class="brand-name">WASAPHOTO</h2>
          <button  class="icon-btn home" title="Home"><i class="fa-solid fa-house"></i></button>
        </div>
        
        <div class="search-container">
          <i class="fa-solid fa-magnifying-glass"></i>
          <input type="text" v-model="searchQuery" placeholder="Cerca tra gli utenti..." @keyup.enter="handleSearch">
        </div>

        <div class="nav-actions">
          <button @click="openUploadModal" class="icon-btn" title="Nuovo Post"><i class="fa-solid fa-circle-plus"></i></button>
          <button @click='goToUserProfile("self")' class="icon-btn" title="Il mio Profilo"><i class="fa-solid fa-user"></i></button>
          <button @click="goToSettingsPage" class="icon-btn settings" title="Impostazioni"><i class="fa-solid fa-gear"></i></button>
          <button @click="doLogout" class="icon-btn logout" title="Logout"><i class="fa-solid fa-right-from-bracket"></i></button>
        </div>
      </div>
    </nav>

    <main class="content-container">
      <div v-if="loading" class="state-container">
        <div class="loader"></div>
        <p>Preparando il tuo feed...</p>
      </div>

      <ErrorMsg v-if="errorMessage" :message="errorMessage" @close="errorMessage = ''" />

      <section v-if="!loading" class="posts-feed">
        <h3 class="section-title">Home Feed</h3>
        
        <div v-if="posts.length > 0" class="stream-list">
          <div v-for="post in posts" :key="post.postId" class="post-wrapper">
            <PostCard :post="post"/>
          </div>
        </div>

        <div v-else class="empty-state">
          <i class="fa-solid fa-photo-film" style="font-size: 3rem; margin-bottom: 20px; opacity: 0.3;"></i>
          <p>Il tuo feed è vuoto. Inizia a seguire qualcuno o cerca nuovi utenti!</p>
        </div>

        <div v-if="loadingMore" class="loading-more">
          <div class="loader small-loader"></div>
        </div>

        <div v-if="allLoaded && posts.length > 0" class="end-message">
          <p>Hai visualizzato tutti i post. Torna più tardi!</p>
        </div>
      </section>
    </main>

    <!-- Upload Modal -->
     <div v-if="showUploadModal" class="modal-overlay" @click.self="closeUploadModal">
      <div class="glass-card upload-modal">
        <h3>Crea un nuovo post</h3>
        <div class="upload-zone" @click="$refs.fileInput.click()">
          <template v-if="!uploadPreview">
            <i class="fa-solid fa-cloud-arrow-up"></i>
            <p>Clicca per selezionare una foto (JPG)</p>
          </template>
          <img v-else :src="uploadPreview" class="preview-img">
        </div>
        <input type="file" ref="fileInput" @change="handleFileSelect" accept="image/jpeg" hidden>
        <textarea v-model="uploadCaption" placeholder="Scrivi una didascalia..." rows="3" maxlength="2200"></textarea>
        <div class="char-counter">{{ uploadCaption.length }}/2200</div>
        
        <div class="modal-actions">
          <button class="gradient-button btn-cancel" @click="closeUploadModal">Annulla</button>
          <button class="gradient-button" @click="submitUpload" :disabled="!selectedFile || uploadLoading">
            {{ uploadLoading ? 'Caricamento...' : 'Condividi' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
#HomeViewContainer {
  min-height: 100vh;
  background: radial-gradient(circle at center, #1a1a1a 0%, #0a0a0a 100%);
  color: white; padding-top: 100px; padding-bottom: 50px; font-family: 'Inter', sans-serif;
}
.main-header {
  padding: 40px 0;
  display: flex;
  justify-content: center;
}

.logo-container {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 15px;
}

.camera-icon {
  font-size: 1.5rem;
  color: white; 
}

.brand-name {
  font-size: 1.2rem;
  font-weight: 700;
  letter-spacing: 2px;
  margin: 0;
  color: white;
}

.glass-nav {
  position: fixed; top: 0; left: 0; width: 100%; height: 75px;
  background: rgba(255, 255, 255, 0.05); backdrop-filter: blur(15px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  display: flex; align-items: center; justify-content: center; z-index: 2000;
}

.nav-content { width: 90%; max-width: 1100px; display: flex; justify-content: space-between; align-items: center; }
.brand { cursor: pointer; font-size: 1.5rem; }
.brand span { color: hsl(0, 0%, 100%); font-weight: bold; }

.search-container {
  background: rgba(255, 255, 255, 0.1); border-radius: 12px; padding: 8px 15px;
  display: flex; align-items: center; gap: 10px; width: 300px;
}
.search-container input { background: transparent; border: none; color: white; outline: none; width: 100%; }

.icon-btn { background: transparent; border: none; color: white; font-size: 1.4rem; cursor: pointer; margin-left: 15px; transition: 0.3s; }
.icon-btn:hover { color: #003366; transform: scale(1.1); }
.logout:hover { color: #66001a; }

.posts-feed { max-width: 600px; margin: 0 auto; }
.section-title { font-size: 1.2rem; margin-bottom: 25px; font-weight: 700; opacity: 0.7; text-align: center; text-transform: uppercase; letter-spacing: 2px; }

.suggested-badge {
  font-size: 0.7rem; color: #003366; background: rgba(0, 51, 102, 0.1);
  display: inline-block; padding: 4px 10px; border-radius: 20px;
  margin-bottom: 8px; border: 1px solid rgba(0, 51, 102, 0.3);
  text-transform: uppercase; font-weight: bold;
}

.post-wrapper { margin-bottom: 40px; }

.loading-more { display: flex; justify-content: center; padding: 20px; }
.small-loader { width: 25px !important; height: 25px !important; }

.end-message { text-align: center; padding: 40px; opacity: 0.5; font-style: italic; }

.glass-card { background: rgba(255, 255, 255, 0.03); backdrop-filter: blur(20px); border: 1px solid rgba(255, 255, 255, 0.08); border-radius: 25px; }
.modal-overlay { position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0, 0, 0, 0.85); backdrop-filter: blur(10px); display: flex; align-items: center; justify-content: center; z-index: 3000; }
.upload-modal { width: 90%; max-width: 500px; padding: 30px; }
.upload-zone { border: 2px dashed rgba(255, 255, 255, 0.2); border-radius: 20px; padding: 40px; margin: 20px 0; text-align: center; cursor: pointer; }
.preview-img { width: 100%; border-radius: 15px; max-height: 300px; object-fit: cover; }
textarea { width: 100%; background: rgba(255, 255, 255, 0.05); border: 1px solid rgba(255, 255, 255, 0.1); border-radius: 12px; color: white; padding: 12px; margin-bottom: 20px; outline: none; resize: none; }

.gradient-button { background: linear-gradient(135deg,#003366 100%); color: white; border: none; border-radius: 12px; padding: 10px 25px; cursor: pointer; font-weight: 600; transition: 0.3s; }

/* Stili specifici per allineare i pulsanti identici alla ProfileView */
.modal-actions { display: flex; justify-content: space-between; align-items: center; width: 100%; margin-top: 25px; padding: 0 10px; }
.btn-cancel {
  background: #66001a !important; color: white; border: 1px solid rgba(255, 255, 255, 0.1); 
  padding: 6px 12px; font-size: 0.65rem; border-radius: 8px; cursor: pointer;
  font-weight: 500; transition: all 0.2s ease; display: inline-flex; align-items: center;
}

.loader { border: 4px solid rgba(255, 255, 255, 0.1); border-top: 4px solid #003366; border-radius: 50%; width: 40px; height: 40px; animation: spin 1s linear infinite; margin: 0 auto 20px; }
@keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }
.state-container { text-align: center; margin-top: 50px; }
.empty-state { text-align: center; padding: 60px; opacity: 0.5; }
</style>