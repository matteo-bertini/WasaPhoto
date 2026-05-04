<script>
/**
 * @file ProfileView.vue
 * @description Component for managing and displaying user profiles, posts, and social interactions.
 */

import ErrorMsg from '../components/ErrorMsg.vue';
import PostCard from '../components/PostCard.vue';

export default {
  name: 'ProfileView',
  components: {
    PostCard,
    ErrorMsg
  },
  data() {
    return {
      /** @type {Object} User profile information and statistics */
      profileData: {
        username: "",
        bio: "",
        followersCount: 0,
        followingCount: 0,
        postsCount: 0,
        isFollowing: false,
        isBannedByMe: false,
        userPosts: [] 
      },
      /** @type {Array<string>} List of user followers */
      followers: [],
      /** @type {Array<string>} List of users followed by the profile */
      following: [],
      searchQuery: "",
      loading: false,
      errorMessage: "",
      
      /** @type {boolean} State management for the post upload modal */
      showUploadModal: false,
      uploadLoading: false,
      selectedFile: null,
      uploadPreview: null,
      uploadCaption: ""
    };
  },
  computed: {
    /**
     * Determines if the current user owns the profile being viewed.
     * @returns {boolean}
     */
    isOwner() {
      return this.$route.params.username === localStorage.getItem("Username");
    },
    /**
     * Extracts the first character of the username for the avatar placeholder.
     * @returns {string}
     */
    userInitial() {
      return this.profileData.username ? this.profileData.username.charAt(0).toUpperCase() : "?";
    }
  },
  methods: {
    /**
     * Fetches general profile information from the API.
     */
    async fetchProfile() {
      this.loading = true;
      this.errorMsg = "";
      try {
        const targetUsername = this.$route.params.username;
        const token = localStorage.getItem("SessionToken");
        
        const response = await this.$axios.get(`/users/${targetUsername}`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        
        this.profileData = {
          username: response.data.username,
          bio: response.data.bio,
          followersCount: response.data.followersCount || 0,
          followingCount: response.data.followingCount || 0,
          postsCount: response.data.postsCount || 0,
          isFollowing: response.data.isFollowing || false,
          isBannedByMe: response.data.isBannedByMe || false,
          userPosts: response.data.userPosts || []
        };
        this.loading = false; 
     } catch (e) {
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
              this.$router.push(`/users/${localStorage.getItem("Username")}`);
              break;
            default:
              this.errorMessage = "Si è verificato un errore sul server. Riprova più tardi.";
          }
        } else {
          this.errorMessage = "Impossibile connettersi al server. Controlla la tua connessione.";
        }
     }
    },

    /**
     * Fetches the followers list for the target user.
     */
    async getFollowers() {
      this.loading = true;
      try {
        const targetUsername = this.$route.params.username;
        const token = localStorage.getItem("SessionToken");
        const response = await this.$axios.get(`/users/${targetUsername}/followers`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        this.followers = response.data.map(element => element.username);
        this.loading = false;
      } catch (e) {
        this.loading = false;
        this.handleApiError(e);
      }
    },

    /**
     * Fetches the following list for the target user.
     */
    async getFollowing() {
      this.loading = true;
      try {
        const targetUsername = this.$route.params.username;
        const token = localStorage.getItem("SessionToken");
        const response = await this.$axios.get(`/users/${targetUsername}/following`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        this.following = response.data.map(element => element.username);
        this.loading = false;
      } catch (e) {
        this.loading = false;
        this.handleApiError(e);
      }
    },
    
    /**
     * Shared error handling for social API calls.
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

    goToSettingsPage() {
      this.$router.push("/settings");
    },

    /**
     * Redirects to a specified user profile or logs out if no username is provided.
     * @param {string} username
     */
    goToUserProfile(username) {
      if (username) {
        if(username === "self") {
          this.$router.push(`/users/${localStorage.getItem("Username")}`);
        } else {
          this.$router.push(`/users/${username}`);
        }
      } else {
        localStorage.clear();
        this.$router.push(`/login`);
      }
    },

    /**
     * Local state cleanup after a post deletion.
     * @param {number} postId
     */
    removePostFromList(postId) {
      this.profileData.userPosts = this.profileData.userPosts.filter(p => p.postId !== postId);
      this.profileData.postsCount--;
    },

    /**
     * Logic for following or unfollowing the profile user.
     */
    async toggleFollow() {
      const token = localStorage.getItem("SessionToken");
      const target = this.profileData.username;
      try {
        if (this.profileData.isFollowing) {
          await this.$axios.delete(`/users/${target}/followers`, {headers: { Authorization: `Bearer ${token}` }});
          this.profileData.followersCount--;
        } else {
          await this.$axios.put(`/users/${target}/followers`, {}, {headers: { Authorization: `Bearer ${token}` }});
          this.profileData.followersCount++;
        }
        this.profileData.isFollowing = !this.profileData.isFollowing;
      } catch (e) {
        this.handleApiError(e);
      }
    },

    /**
     * Logic for banning or unbanning the profile user.
     */
    async toggleBan() {
      const token = localStorage.getItem("SessionToken");
      const target = this.profileData.username;
      try {
        if(!this.profileData.isBannedByMe) {
          await this.$axios.put(`/users/${target}/ban`, {}, {headers: { Authorization: `Bearer ${token}` }});
          this.profileData.isBannedByMe = true;
          this.profileData.isFollowing = false;
        } else {
          await this.$axios.delete(`/users/${target}/ban`, {headers: { Authorization: `Bearer ${token}` }});
          this.profileData.isBannedByMe = false;
        }
      } catch (e) {
        this.handleApiError(e);
      }
    },

    openUploadModal() { 
      this.showUploadModal = true; 
    },

    closeUploadModal() {
      if (this.uploadPreview) URL.revokeObjectURL(this.uploadPreview);
      this.showUploadModal = false;
      this.selectedFile = null;
      this.uploadPreview = null;
      this.uploadCaption = "";
    },

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
          headers: { 
            Authorization: `Bearer ${token}`,
            'Content-Type': 'multipart/form-data' 
          }
        });
        this.profileData.userPosts.unshift(response.data);
        this.profileData.postsCount++;
        this.closeUploadModal();
      } catch (e) {
        alert("Upload failed.");
      } finally { 
        this.uploadLoading = false; 
      }
    },

    handleSearch() {
      this.loading = true;
      if (this.searchQuery.trim()) {
        this.$router.push(`/users/${this.searchQuery.trim()}`);
        this.searchQuery = "";
      }
      this.loading = false;
    },

    /**
     * Clears session token and redirects to login page.
     */
    async doLogout() {
      this.loading = true;
      const token = localStorage.getItem("SessionToken");
      try {
        await this.$axios.delete("/session", {headers: { Authorization: `Bearer ${token}` }});
        localStorage.clear();
        this.loading = false;
        this.$router.push("/login");
      } catch (e) {
        this.loading = false;
        this.handleApiError(e);
      }
    }
  },
  mounted() {
    this.fetchProfile();
  },
  watch: {
    /**
     * Refetches data when the route parameter changes.
     */
    '$route.params.username'(newVal) {
      if (newVal && localStorage.getItem("SessionToken")) {
        this.fetchProfile(newVal);
      }
    }
  }
};
</script>

<template>
  <div id="ProfilePageContainer">
    
    <nav class="glass-nav">
      <div class="nav-content">
        <div @click="$router.push('/home')" class="logo-container">
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
        <p>Caricamento profilo...</p>
      </div>

      <ErrorMsg v-if="errorMessage" :message="errorMessage" @close="errorMessage = ''" />

      <template v-else>
        <section class="profile-card glass-card">
          <div class="avatar-circle">{{ userInitial }}</div>
          <h2 class="profile-username">@{{ profileData.username }}</h2>
          
          <template v-if="!profileData.isBannedByMe">
            <p class="profile-bio">{{ profileData.bio || 'Nessuna biografia impostata.' }}</p>
            
            <div class="stats-row">
              <div class="stat-item">
                <span class="stat-value">{{ profileData.postsCount }}</span>
                <span class="stat-label">Post</span>
              </div>
              <div @click="getFollowers" data-bs-toggle="modal" data-bs-target="#followersModal" class="stat-item clickable">
                <span class="stat-value">{{ profileData.followersCount }}</span>
                <span class="stat-label">Followers</span>
              </div>
              <div @click="getFollowing" data-bs-toggle="modal" data-bs-target="#followingModal" class="stat-item clickable">
                <span class="stat-value">{{ profileData.followingCount }}</span>
                <span class="stat-label">Following</span>
              </div>
            </div>
          </template>

          <div v-if="!isOwner" class="profile-actions">
            <button 
              v-if="!profileData.isBannedByMe" 
              :class="['gradient-button', profileData.isFollowing ? 'btn-unfollow' : '']" 
              @click="toggleFollow"
            >
              {{ profileData.isFollowing ? 'Unfollow' : 'Follow' }}
            </button>
            
            <button :class="['btn-ban', profileData.isBannedByMe ? 'active-ban' : '']" @click="toggleBan">
              <i class="fa-solid" :class="profileData.isBannedByMe ? 'fa-user-check' : 'fa-user-slash'"></i>
              {{ profileData.isBannedByMe ? ' Unban' : ' Ban' }}
            </button>
          </div>
        </section>

        <section v-if="!profileData.isBannedByMe" class="posts-feed">
          <h3 class="section-title">Post</h3>
          <div v-if="profileData.userPosts && profileData.userPosts.length > 0" class="stream-list">
            <PostCard 
              v-for="post in profileData.userPosts" 
              :key="post.postId" 
              :post="post" 
              @post-deleted="removePostFromList"
            />
          </div>
          <div v-else class="empty-state">
            <p>Nessun post da mostrare.</p>
          </div>
        </section>

        <section v-else class="empty-state banned-notice">
          <i class="fa-solid fa-eye-slash" style="font-size: 2rem; margin-bottom: 10px; opacity: 0.5;"></i>
          <p>Hai bloccato questo utente. Sbloccalo per vedere i contenuti.</p>
        </section>
      </template>
    </main>

    <div class="modal fade custom-modal" data-bs-dismiss="modal" id="followersModal" tabindex="-1">
        <div class="modal-dialog modal-dialog-centered modal-dialog-scrollable">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title">Followers</h5>
                    <button type="button" class="btn-close btn-close-white" data-bs-dismiss="modal"></button>
                </div>
                <div class="modal-body">
                    <div v-if="followers.length === 0" class="empty-state">
                        Nessun utente trovato.
                    </div>
                    <div v-for="follower in followers" :key="follower" class="social-item" @click="goToUserProfile(follower)">
                        <span class="social-username">@{{follower}}</span>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <div class="modal fade custom-modal" data-bs-dismiss="modal" id="followingModal" tabindex="-1">
        <div class="modal-dialog modal-dialog-centered modal-dialog-scrollable">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title">Following</h5>
                    <button type="button" class="btn-close btn-close-white" data-bs-dismiss="modal"></button>
                </div>
                <div class="modal-body">
                    <div v-if="following.length === 0" class="empty-state">
                        Nessun utente trovato.
                    </div>
                    <div v-for="followed in following" :key="followed" class="social-item" @click="goToUserProfile(followed)">
                        <span class="social-username">@{{followed}}</span>
                    </div>
                </div>
            </div>
        </div>
    </div>

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
/* Page main layout container */
#ProfilePageContainer {
  min-height: 100vh;
  background: radial-gradient(circle at center, #1a1a1a 0%, #0a0a0a 100%);
  color: white; padding-top: 100px; padding-bottom: 50px; font-family: 'Inter', sans-serif;
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

/* Glassmorphism navigation bar styling */
.glass-nav {
  position: fixed; top: 0; left: 0; width: 100%; height: 75px;
  background: rgba(255, 255, 255, 0.05); backdrop-filter: blur(15px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  display: flex; align-items: center; justify-content: center; z-index: 2000;
}

.nav-content { width: 90%; max-width: 1100px; display: flex; justify-content: space-between; align-items: center; }
.brand { cursor: pointer; font-size: 1.5rem; }
.brand span { color: hsl(0, 0%, 100%); font-weight: bold; }

/* User search input container */
.search-container {
  background: rgba(255, 255, 255, 0.1); border-radius: 12px; padding: 8px 15px;
  display: flex; align-items: center; gap: 10px; width: 300px;
}
.search-container input { background: transparent; border: none; color: white; outline: none; width: 100%; }

/* Navigation buttons and interaction hover effects */
.icon-btn { background: transparent; border: none; color: white; font-size: 1.4rem; cursor: pointer; margin-left: 15px; transition: 0.3s; }
.icon-btn:hover { color: #003366; transform: scale(1.1); }
.logout:hover { color: #66001a; }


/* Profile card main glassmorphism effect */
.glass-card { background: rgba(255, 255, 255, 0.03); backdrop-filter: blur(20px); border: 1px solid rgba(255, 255, 255, 0.08); border-radius: 25px; }
.profile-card { max-width: 650px; margin: 0 auto 40px; padding: 40px; text-align: center; }

/* Profile avatar circle placeholder styling */
.avatar-circle {
  width: 110px; height: 110px; border-radius: 50%;
  background: linear-gradient(135deg,#003366 100%);  color: white;
  margin: 0 auto 20px; display: flex; align-items: center; justify-content: center;
  font-size: 3rem; font-weight: bold; box-shadow: 0 0 20px rgba(0, 51, 102, 0.4);
}

/* Grid layout for user statistics (Followers/Following/Posts) */
.stats-row { display: flex; justify-content: center; gap: 50px; margin: 25px 0; }
.stat-item { display: flex; flex-direction: column; align-items: center; }
.stat-item.clickable {
  cursor: pointer; 
  transition: opacity 0.2s ease;
}
.stat-value { font-size: 1.4rem; font-weight: 800; }
.stat-label { font-size: 0.85rem; opacity: 0.6; text-transform: uppercase; }

/* Standard buttons and variant styling */
.gradient-button {
  background: linear-gradient(135deg,#003366 100%);  color: white;
  border: none; border-radius: 12px; color: white; font-weight: 600;
  cursor: pointer; transition: 0.3s; padding: 10px 25px;
}
.btn-cancel {
  background: #66001a !important; color: white; border: 1px solid rgba(255, 255, 255, 0.1); 
  padding: 6px 12px; font-size: 0.65rem; border-radius: 8px; cursor: pointer;
  font-weight: 500; transition: all 0.2s ease; display: inline-flex; align-items: center;
}
.btn-unfollow { background: transparent !important; border: 2px solid #003366 !important; }
.btn-ban { background: #66001a; color: #ffffff; border: 1px solid rgba(255, 71, 87, 0.2); padding: 10px 15px; border-radius: 12px; cursor: pointer; margin-left: 10px; }
.active-ban { background: #66001a; color: white; }

/* Main posts feed layout */
.posts-feed { max-width: 600px; margin: 0 auto; }
.section-title { margin-bottom: 25px; font-weight: 300; opacity: 0.7; }
.banned-notice { text-align: center; padding: 40px; opacity: 0.6; }

/* Modal overlay and dark theme customization */
.modal-overlay { position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0, 0, 0, 0.85); backdrop-filter: blur(10px); display: flex; align-items: center; justify-content: center; z-index: 3000; }
.upload-modal { width: 90%; max-width: 500px; padding: 30px; }

.custom-modal .modal-content {
    background: rgba(10, 10, 10, 0.98); backdrop-filter: blur(20px);
    border: 1px solid rgba(255, 255, 255, 0.08); border-radius: 20px; color: white;
}
.custom-modal .modal-header { border-bottom: 1px solid rgba(255, 255, 255, 0.1); padding: 1.5rem; }
.custom-modal .modal-title { font-weight: 600; letter-spacing: 0.5px; }

.btn-close-white { transition: all 0.2s ease-in-out; }
.btn-close-white:hover { color: #66001a; opacity: 1; transform: scale(1.1); }

/* Faster transitions for modals */
.modal.fade .modal-dialog { transition: transform 0.15s ease-out; }
.modal-backdrop.fade { transition: opacity 0.15s linear; }

/* Social interaction items in modals */
.social-item {
    padding: 15px 25px; cursor: pointer; transition: all 0.25s ease;
    display: flex; align-items: center; border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}
.social-item:hover .social-username { color: #003366; }
.social-username { font-size: 1rem; font-weight: 500; color: rgba(255, 255, 255, 0.9); transition: color 0.2s ease; }

/* Empty states and scrollbar styling */
.empty-state { padding: 3rem; text-align: center; opacity: 0.5; font-size: 0.9rem; }
.modal-body::-webkit-scrollbar { width: 5px; }
.modal-body::-webkit-scrollbar-thumb { background: rgba(255, 255, 255, 0.1); border-radius: 10px; }

/* Upload logic specific components */
.modal-actions { display: flex; justify-content: space-between; align-items: center; width: 100%; margin-top: 25px; padding: 0 10px; }
.upload-zone { border: 2px dashed rgba(255, 255, 255, 0.2); border-radius: 20px; padding: 40px; margin: 20px 0; text-align: center; cursor: pointer; }
.preview-img { width: 100%; border-radius: 15px; max-height: 300px; object-fit: cover; }
textarea { width: 100%; background: rgba(255, 255, 255, 0.05); border: 1px solid rgba(255, 255, 255, 0.1); border-radius: 12px; color: white; padding: 12px; margin-bottom: 20px; outline: none; resize: none; }

/* Loading spinner animation */
.loader { border: 4px solid rgba(255, 255, 255, 0.1); border-top: 4px solid #a53b59; border-radius: 50%; width: 40px; height: 40px; animation: spin 1s linear infinite; margin: 0 auto 20px; }
@keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }
</style>