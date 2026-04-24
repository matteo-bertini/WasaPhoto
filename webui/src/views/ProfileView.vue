<script>
// 1. Import the PostCard component
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
      searchQuery: "",
      loading: false,
      errorMsg: "",
      
      // Upload State
      showUploadModal: false,
      uploadLoading: false,
      selectedFile: null,
      uploadPreview: null,
      uploadCaption: ""
    };
  },
  computed: {
    isOwner() {
      return this.$route.params.username === localStorage.getItem("Username");
    },
    // Get initial for avatar placeholder
    userInitial() {
      return this.profileData.username ? this.profileData.username.charAt(0).toUpperCase() : "?";
    }
  },
  methods: {
 
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
      } catch (e) {
        this.handleFetchError(e);
      } finally {
        this.loading = false;
      }
    },

    removePostFromList(postId) {
      this.profileData.userPosts = this.profileData.userPosts.filter(p => p.postId !== postId);
      this.profileData.postsCount--;
      console.log(`[Profile] Post ${postId} removed from local UI.`);
    },

    async toggleFollow() {
      const token = localStorage.getItem("SessionToken");
      const target = this.profileData.username;

      try {
        if (this.profileData.isFollowing) {
          await this.$axios.delete(`/users/${target}/followers`, {
            headers: { Authorization: `Bearer ${token}` }
          });
          this.profileData.followersCount--;
        } else {
          await this.$axios.put(`/users/${target}/followers`, {}, {
            headers: { Authorization: `Bearer ${token}` }
          });
          this.profileData.followersCount++;
        }
        this.profileData.isFollowing = !this.profileData.isFollowing;
      } catch (e) {
        console.error("[Profile] Follow error:", e);
      }
    },

    async toggleBan() {
      const token = localStorage.getItem("SessionToken");
      const myUsername = localStorage.getItem("Username");
      const target = this.profileData.username;

      try {
        await this.$axios.put(`/users/${target}/ban`, {}, {
          headers: { Authorization: `Bearer ${token}` }
        });
        this.profileData.isBannedByMe = true;
        // Redirect to own profile after banning
        this.$router.push(`/users/${myUsername}`); 
      } catch (e) {
        console.error("[Profile] Ban error:", e);
      }
    },

    // --- Upload Logic ---
    openUploadModal() { this.showUploadModal = true; },
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
        // Add new post to top of list
        this.profileData.userPosts.unshift(response.data);
        this.profileData.postsCount++;
        this.closeUploadModal();
      } catch (e) {
        alert("Upload failed.");
      } finally { this.uploadLoading = false; }
    },

    handleSearch() {
      if (this.searchQuery.trim()) {
        this.$router.push(`/users/${this.searchQuery.trim()}`);
        this.searchQuery = "";
      }
    },

    async handleLogout() {
      const token = localStorage.getItem("SessionToken");
      try {
        await this.$axios.delete("/session", {
          headers: { Authorization: `Bearer ${token}` }
        });
      } catch (e) {
        console.warn("[Profile] Logout session already invalid.");
      } finally {
        localStorage.clear();
        this.$router.push("/login");
      }
    },

    handleFetchError(e) {
      if (e.response) {
        const status = e.response.status;
        if (status === 401) {
          localStorage.clear();
          this.$router.push('/login');
        } else if (status === 404) this.errorMsg = "User not found.";
        else if (status === 403) this.errorMsg = "Access forbidden.";
        else this.errorMsg = "An error occurred loading the profile.";
      } else {
        this.errorMsg = "Connection failed.";
      }
      console.error("[Profile] Fetch Error:", e);
    }
  },
  mounted() {
    this.fetchProfile();
  },
  watch: {
    '$route.params.username': 'fetchProfile'
  }
};
</script>

<template>
  <div id="ProfilePageContainer">
    <nav class="glass-nav">
      <div class="nav-content">
        <h2 class="brand" @click="$router.push('/home')">WASAPHOTO</h2>
        
        <div class="search-container">
          <i class="fa-solid fa-magnifying-glass"></i>
          <input type="text" v-model="searchQuery" placeholder="Cerca tra gli utenti..." @keyup.enter="handleSearch">
        </div>

        <div class="nav-actions">
          <button @click="openUploadModal" class="icon-btn" title="Nuovo Post"><i class="fa-solid fa-circle-plus"></i></button>
          <button @click="$router.push(`/users/${localStorage.getItem('Username')}`)" class="icon-btn" title="Mio Profilo"><i class="fa-solid fa-user"></i></button>
          <button @click="handleLogout" class="icon-btn logout" title="Logout"><i class="fa-solid fa-right-from-bracket"></i></button>
        </div>
      </div>
    </nav>

    <main class="content-container">
      <div v-if="loading" class="state-container">
        <div class="loader"></div>
        <p>Caricamento profilo...</p>
      </div>

      <div v-else-if="errorMsg" class="state-container glass-card error-box">
        <i class="fa-solid fa-circle-exclamation"></i>
        <p>{{ errorMsg }}</p>
        <button @click="$router.push('/home')" class="gradient-button">Torna alla Home</button>
      </div>

      <template v-else>
        <section class="profile-card glass-card">
          <div class="avatar-circle">{{ userInitial }}</div>
          <h2 class="profile-username">@{{ profileData.username }}</h2>
          <p class="profile-bio">{{ profileData.bio || 'Nessuna biografia impostata.' }}</p>
          
          <div class="stats-row">
            <div class="stat-item">
              <span class="stat-value">{{ profileData.postsCount }}</span>
              <span class="stat-label">Post</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ profileData.followersCount }}</span>
              <span class="stat-label">Followers</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ profileData.followingCount }}</span>
              <span class="stat-label">Following</span>
            </div>
          </div>

          <div v-if="!isOwner" class="profile-actions">
            <button :class="['gradient-button', profileData.isFollowing ? 'btn-unfollow' : '']" @click="toggleFollow">
              {{ profileData.isFollowing ? 'Unfollow' : 'Follow' }}
            </button>
            <button :class="['btn-ban', profileData.isBannedByMe ? 'active-ban' : '']" @click="toggleBan">
              <i class="fa-solid" :class="profileData.isBannedByMe ? 'fa-user-check' : 'fa-user-slash'"></i>
              {{ profileData.isBannedByMe ? ' Unban' : ' Ban' }}
            </button>
          </div>
        </section>

        <section class="posts-feed">
          <h3 class="section-title">Post </h3>
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
      </template>
    </main>

    <div v-if="showUploadModal" class="modal-overlay" @click.self="closeUploadModal">
      <div class="glass-card upload-modal">
        <h3>Crea nuovo post</h3>
        <div class="upload-zone" @click="$refs.fileInput.click()">
          <template v-if="!uploadPreview">
            <i class="fa-solid fa-cloud-arrow-up"></i>
            <p>Clicca per selezionare una foto (JPG)</p>
          </template>
          <img v-else :src="uploadPreview" class="preview-img">
        </div>
        <input type="file" ref="fileInput" @change="handleFileSelect" accept="image/jpeg" hidden>
        <textarea v-model="uploadCaption" placeholder="Scrivi una didascalia..." rows="3"></textarea>
        <div class="modal-actions">
          <button class="btn-cancel" @click="closeUploadModal">Annulla</button>
          <button class="gradient-button" @click="submitUpload" :disabled="!selectedFile || uploadLoading">
            {{ uploadLoading ? 'Caricamento...' : 'Condividi' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
#ProfilePageContainer {
  min-height: 100vh;
  background: radial-gradient(circle at center, #1a1a1a 0%, #0a0a0a 100%);
  color: white; padding-top: 100px; padding-bottom: 50px; font-family: 'Inter', sans-serif;
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
.icon-btn:hover { color: #a53b59; transform: scale(1.1); }
.glass-card { background: rgba(255, 255, 255, 0.03); backdrop-filter: blur(20px); border: 1px solid rgba(255, 255, 255, 0.08); border-radius: 25px; }
.profile-card { max-width: 650px; margin: 0 auto 40px; padding: 40px; text-align: center; }
.avatar-circle {
  width: 110px; height: 110px; border-radius: 50%;
  background: linear-gradient(135deg, #a53b59 0%, #726fb4 100%);
  margin: 0 auto 20px; display: flex; align-items: center; justify-content: center;
  font-size: 3rem; font-weight: bold; box-shadow: 0 0 20px rgba(165, 59, 89, 0.4);
}
.stats-row { display: flex; justify-content: center; gap: 50px; margin: 25px 0; }
.stat-item { display: flex; flex-direction: column; align-items: center; }
.stat-value { font-size: 1.4rem; font-weight: 800; }
.stat-label { font-size: 0.85rem; opacity: 0.6; text-transform: uppercase; }
.gradient-button {
  background: linear-gradient(135deg, #a53b59 0%, #726fb4 100%);
  border: none; border-radius: 12px; color: white; font-weight: 600;
  cursor: pointer; transition: 0.3s; padding: 10px 25px;
}
.btn-unfollow { background: transparent !important; border: 2px solid #a53b59 !important; }
.btn-ban { background: rgba(255, 71, 87, 0.1); color: #ff4757; border: 1px solid rgba(255, 71, 87, 0.2); padding: 10px 15px; border-radius: 12px; cursor: pointer; margin-left: 10px; }
.active-ban { background: #ff4757; color: white; }
.posts-feed { max-width: 600px; margin: 0 auto; }
.section-title { margin-bottom: 25px; font-weight: 300; opacity: 0.7; }
.modal-overlay { position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0, 0, 0, 0.85); backdrop-filter: blur(10px); display: flex; align-items: center; justify-content: center; z-index: 3000; }
.upload-modal { width: 90%; max-width: 500px; padding: 30px; }
.upload-zone { border: 2px dashed rgba(255, 255, 255, 0.2); border-radius: 20px; padding: 40px; margin: 20px 0; text-align: center; cursor: pointer; }
.preview-img { width: 100%; border-radius: 15px; max-height: 300px; object-fit: cover; }
textarea { width: 100%; background: rgba(255, 255, 255, 0.05); border: 1px solid rgba(255, 255, 255, 0.1); border-radius: 12px; color: white; padding: 12px; margin-bottom: 20px; outline: none; resize: none; }
.loader { border: 4px solid rgba(255, 255, 255, 0.1); border-top: 4px solid #a53b59; border-radius: 50%; width: 40px; height: 40px; animation: spin 1s linear infinite; margin: 0 auto 20px; }
@keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }
</style>