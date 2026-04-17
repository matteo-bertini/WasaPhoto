<template>
  <div id="ProfilePageContainer">
    <nav class="glass-nav">
      <div class="nav-content">
        <h1 class="brand" @click="$router.push('/home')">WASA<span>PHOTO</span></h1>
        
        <div class="search-container">
          <i class="fa-solid fa-magnifying-glass"></i>
          <input type="text" v-model="searchQuery" placeholder="Search users..." @keyup.enter="handleSearch">
        </div>

        <div class="nav-actions">
          <button @click="openUploadModal" class="icon-btn" title="Upload"><i class="fa-solid fa-circle-plus"></i></button>
          <button @click="goToMyProfile" class="icon-btn" title="My Profile"><i class="fa-solid fa-user"></i></button>
          <button @click="$router.push('/settings')" class="icon-btn" title="Settings"><i class="fa-solid fa-gear"></i></button>
          <button @click="handleLogout" class="icon-btn logout" title="Logout"><i class="fa-solid fa-right-from-bracket"></i></button>
        </div>
      </div>
    </nav>

    <main class="content-container">
      <div v-if="loading" class="state-container">
        <div class="loader"></div>
        <p>Loading profile...</p>
      </div>

      <div v-else-if="errorMsg" class="state-container glass-card error-box">
        <i class="fa-solid fa-circle-exclamation"></i>
        <p>{{ errorMsg }}</p>
        <button @click="$router.push('/home')" class="gradient-button">Back to Home</button>
      </div>

      <template v-else>
        <section class="profile-card glass-card">
          <div class="avatar-circle">
            {{ userInitial }}
          </div>
          
          <h2 class="profile-username">@{{ profileData.Username }}</h2>
          
          <p class="profile-bio">{{ profileData.Bio || 'No bio available.' }}</p>
          
          <div class="stats-row">
            <div class="stat-item">
              <span class="stat-value">{{ profileData.PostsCount }}</span>
              <span class="stat-label">Posts</span>
            </div>
            <div class="stat-item clickable">
              <span class="stat-value">{{ profileData.FollowersCount }}</span>
              <span class="stat-label">Followers</span>
            </div>
            <div class="stat-item clickable">
              <span class="stat-value">{{ profileData.FollowingCount }}</span>
              <span class="stat-label">Following</span>
            </div>
          </div>

          <div v-if="!isOwner" class="profile-actions">
            <button 
              :class="['gradient-button', profileData.IsFollowing ? 'btn-unfollow' : '']" 
              @click="toggleFollow">
              {{ profileData.IsFollowing ? 'Unfollow' : 'Follow' }}
            </button>
            
            <button 
              :class="['btn-ban', profileData.IsBannedByMe ? 'active-ban' : '']" 
              @click="toggleBan" 
              :title="profileData.IsBannedByMe ? 'Unban User' : 'Ban User'">
              <i class="fa-solid" :class="profileData.IsBannedByMe ? 'fa-user-check' : 'fa-user-slash'"></i>
              {{ profileData.IsBannedByMe ? ' Unban' : '' }}
            </button>
          </div>
        </section>

        <section class="posts-feed">
          <h3 class="section-title">Latest Posts</h3>
          <div v-if="profileData.UserPosts && profileData.UserPosts.length > 0" class="stream-list">
            <div v-for="post in profileData.UserPosts" :key="post.PostId" class="post-card glass-card">
              <img :src="'/api/images/' + post.PostId" class="post-image" alt="Post">
              
              <div class="post-info">
                <p v-if="post.Caption" class="post-caption">{{ post.Caption }}</p>
                <div class="post-meta">
                  <span :class="{'liked': post.IsLikedByMe}">
                    <i class="fa-solid fa-heart"></i> {{ post.LikesNumber }}
                  </span>
                  <span><i class="fa-solid fa-comment"></i> {{ post.CommentsNumber }}</span>
                  <button v-if="isOwner" class="btn-delete" @click="deletePost(post.PostId)">
                    <i class="fa-solid fa-trash"></i>
                  </button>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="empty-state">
            <p>No posts to show yet.</p>
          </div>
        </section>
      </template>
    </main>

    <div v-if="showUploadModal" class="modal-overlay" @click.self="closeUploadModal">
      <div class="glass-card upload-modal">
        <h3>Create New Post</h3>
        
        <div class="upload-zone" @click="$refs.fileInput.click()">
          <template v-if="!uploadPreview">
            <i class="fa-solid fa-cloud-arrow-up"></i>
            <p>Click to select a photo</p>
          </template>
          <img v-else :src="uploadPreview" class="preview-img">
        </div>
        
        <input type="file" ref="fileInput" @change="handleFileSelect" accept="image/jpeg,image/png" hidden>
        
        <textarea v-model="uploadCaption" placeholder="Write a caption..." rows="3"></textarea>
        
        <div class="modal-actions">
          <button class="btn-cancel" @click="closeUploadModal">Cancel</button>
          <button class="gradient-button" @click="submitUpload" :disabled="!selectedFile || uploadLoading">
            {{ uploadLoading ? 'Uploading...' : 'Share Post' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      profileData: {
        Username: "",
        Bio: "",
        FollowersCount: 0,
        FollowingCount: 0,
        PostsCount: 0,
        IsFollowing: false,
        IsBannedByMe: false,
        UserPosts: []
      },
      searchQuery: "",
      loading: false,
      errorMsg: "",
      
      // Upload States
      showUploadModal: false,
      uploadLoading: false,
      selectedFile: null,
      uploadPreview: null,
      uploadCaption: ""
    };
  },
  computed: {
    isOwner() {
      return this.$route.params.Username === localStorage.getItem("Username");
    },
    userInitial() {
      return this.profileData.Username ? this.profileData.Username.charAt(0).toUpperCase() : "?";
    }
  },
  methods: {
    async fetchProfile() {
      this.loading = true;
      this.errorMsg = "";
      try {
        const targetUsername = this.$route.params.Username;
        const token = localStorage.getItem("SessionToken");
        
        const response = await this.$axios.get(`/users/${targetUsername}`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        
        this.profileData = response.data;
      } catch (e) {
        if (e.response && e.response.status === 404) {
          this.errorMsg = "User not found.";
        } else if (e.response && e.response.status === 403) {
          this.errorMsg = "You are banned or cannot view this profile.";
        } else {
          this.errorMsg = "An error occurred while loading the profile.";
        }
      } finally {
        this.loading = false;
      }
    },

    // --- Upload Methods ---
    openUploadModal() {
      if (!this.isOwner) {
        alert("You can only upload posts to your own profile.");
        return;
      }
      this.showUploadModal = true;
    },
    closeUploadModal() {
      this.showUploadModal = false;
      this.resetUpload();
    },
    resetUpload() {
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
      const token = localStorage.getItem("SessionToken");
      const username = localStorage.getItem("Username");

      const formData = new FormData();
      formData.append("file", this.selectedFile);
      formData.append("caption", this.uploadCaption);

      try {
        const response = await this.$axios.post(`/users/${username}/posts`, formData, {
          headers: { 
            Authorization: `Bearer ${token}`,
            'Content-Type': 'multipart/form-data'
          }
        });

        // Add the new post to the grid immediately
        this.profileData.UserPosts.unshift(response.data);
        this.profileData.PostsCount++;
        
        this.closeUploadModal();
      } catch (e) {
        console.error("Upload error:", e);
        alert("Failed to upload post.");
      } finally {
        this.uploadLoading = false;
      }
    },

    // --- Other Methods ---
    async toggleFollow() {
      const token = localStorage.getItem("SessionToken");
      const myUsername = localStorage.getItem("Username");
      const target = this.profileData.Username;

      try {
        if (this.profileData.IsFollowing) {
          await this.$axios.delete(`/users/${target}/followers/${myUsername}`, {
            headers: { Authorization: `Bearer ${token}` }
          });
          this.profileData.FollowersCount--;
        } else {
          await this.$axios.put(`/users/${target}/followers/${myUsername}`, {}, {
            headers: { Authorization: `Bearer ${token}` }
          });
          this.profileData.FollowersCount++;
        }
        this.profileData.IsFollowing = !this.profileData.IsFollowing;
      } catch (e) {
        console.error("Toggle follow error:", e);
      }
    },

    async toggleBan() {
      const token = localStorage.getItem("SessionToken");
      const myUsername = localStorage.getItem("Username");
      const target = this.profileData.Username;

      try {
        if (this.profileData.IsBannedByMe) {
          await this.$axios.delete(`/users/${myUsername}/bans/${target}`, {
            headers: { Authorization: `Bearer ${token}` }
          });
          this.profileData.IsBannedByMe = false;
        } else {
          await this.$axios.put(`/users/${myUsername}/bans/${target}`, {}, {
            headers: { Authorization: `Bearer ${token}` }
          });
          this.profileData.IsBannedByMe = true;
          this.profileData.IsFollowing = false;
        }
      } catch (e) {
        console.error("Toggle ban error:", e);
      }
    },

    handleSearch() {
      if (this.searchQuery.trim()) {
        this.$router.push(`/users/${this.searchQuery.trim()}`);
        this.searchQuery = "";
      }
    },
    goToMyProfile() {
      this.$router.push(`/users/${localStorage.getItem("Username")}`);
    },
    handleLogout() {
      localStorage.clear();
      this.$router.push('/login');
    },
    async deletePost(postId) {
      if (confirm("Delete this post?")) {
        try {
          const token = localStorage.getItem("SessionToken");
          await this.$axios.delete(`/posts/${postId}`, {
            headers: { Authorization: `Bearer ${token}` }
          });
          this.profileData.UserPosts = this.profileData.UserPosts.filter(p => p.PostId !== postId);
          this.profileData.PostsCount--;
        } catch (e) {
          console.error("Delete post error:", e);
        }
      }
    }
  },
  mounted() {
    this.fetchProfile();
  },
  watch: {
    '$route.params.Username': 'fetchProfile'
  }
};
</script>

<style scoped>
#ProfilePageContainer {
  min-height: 100vh;
  background: radial-gradient(circle at center, #1a1a1a 0%, #0a0a0a 100%);
  color: white;
  padding-top: 100px;
  padding-bottom: 50px;
}

/* Nav Styles */
.glass-nav {
  position: fixed; top: 0; left: 0; width: 100%; height: 75px;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(15px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  display: flex; align-items: center; justify-content: center; z-index: 2000;
}
.nav-content { width: 90%; max-width: 1100px; display: flex; justify-content: space-between; align-items: center; }
.brand { cursor: pointer; font-size: 1.5rem; }
.brand span { color: #a53b59; font-weight: bold; }

.search-container {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 12px; padding: 8px 15px; display: flex; align-items: center; gap: 10px; width: 300px;
}
.search-container input { background: transparent; border: none; color: white; outline: none; width: 100%; }

.icon-btn { 
  background: transparent; border: none; color: white; font-size: 1.4rem; 
  cursor: pointer; margin-left: 15px; transition: 0.3s; 
}
.icon-btn:hover { color: #a53b59; transform: scale(1.1); }

/* Card & Content Styles */
.glass-card {
  background: rgba(255, 255, 255, 0.03);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 25px;
}

.profile-card { max-width: 650px; margin: 0 auto; padding: 40px; text-align: center; }

.avatar-circle {
  width: 110px; height: 110px; border-radius: 50%;
  background: linear-gradient(135deg, #a53b59 0%, #726fb4 100%);
  margin: 0 auto 20px; display: flex; align-items: center; justify-content: center;
  font-size: 3rem; font-weight: bold; box-shadow: 0 0 20px rgba(165, 59, 89, 0.4);
}

.profile-bio { margin: 15px 0; opacity: 0.8; line-height: 1.4; }

.stats-row { display: flex; justify-content: center; gap: 50px; margin: 25px 0; }
.stat-item { display: flex; flex-direction: column; align-items: center; }
.stat-value { font-size: 1.4rem; font-weight: 800; }
.stat-label { font-size: 0.85rem; opacity: 0.6; text-transform: uppercase; }

.gradient-button {
  background: linear-gradient(135deg, #a53b59 0%, #726fb4 100%);
  border: none; border-radius: 12px; color: white; font-weight: 600;
  cursor: pointer; transition: 0.3s; padding: 10px 25px;
}
.gradient-button:disabled { opacity: 0.5; cursor: not-allowed; }

.btn-unfollow { background: transparent !important; border: 2px solid #a53b59 !important; }

.btn-ban { 
  background: rgba(255, 71, 87, 0.1); color: #ff4757; border: 1px solid rgba(255, 71, 87, 0.2);
  padding: 10px 15px; border-radius: 12px; cursor: pointer; margin-left: 10px; transition: 0.3s;
}
.active-ban { background: #ff4757; color: white; }

/* Posts Feed */
.posts-feed { max-width: 600px; margin: 50px auto; }
.section-title { margin-bottom: 25px; font-weight: 300; opacity: 0.7; }
.post-card { margin-bottom: 40px; padding: 15px; overflow: hidden; }
.post-image { width: 100%; border-radius: 15px; display: block; min-height: 200px; background: #222; }

.post-info { padding: 15px 5px 5px; }
.post-caption { text-align: left; font-size: 0.95rem; margin-bottom: 12px; color: #eee; }

.post-meta { display: flex; gap: 20px; align-items: center; }
.post-meta i { color: #a53b59; }
.liked i { color: #ff4757; }

.btn-delete { background: transparent; border: none; color: #ff4757; cursor: pointer; margin-left: auto; opacity: 0.6; }

/* Modal Styles */
.modal-overlay {
  position: fixed; top: 0; left: 0; width: 100%; height: 100%;
  background: rgba(0, 0, 0, 0.85); backdrop-filter: blur(10px);
  display: flex; align-items: center; justify-content: center; z-index: 3000;
}
.upload-modal { width: 90%; max-width: 500px; padding: 30px; }
.upload-zone {
  border: 2px dashed rgba(255, 255, 255, 0.2); border-radius: 20px;
  padding: 40px; margin: 20px 0; text-align: center; cursor: pointer; transition: 0.3s;
}
.upload-zone:hover { border-color: #a53b59; background: rgba(165, 59, 89, 0.05); }
.upload-zone i { font-size: 3rem; margin-bottom: 10px; opacity: 0.5; color: #a53b59; }
.preview-img { width: 100%; border-radius: 15px; max-height: 300px; object-fit: cover; }

textarea {
  width: 100%; background: rgba(255, 255, 255, 0.05); border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px; color: white; padding: 12px; margin-bottom: 20px; outline: none; font-family: inherit;
}

.modal-actions { display: flex; gap: 15px; justify-content: flex-end; align-items: center; }
.btn-cancel { background: transparent; border: none; color: white; cursor: pointer; opacity: 0.6; }

/* State & Loaders */
.state-container { text-align: center; padding: 100px 20px; }
.loader {
  border: 4px solid rgba(255, 255, 255, 0.1);
  border-top: 4px solid #a53b59;
  border-radius: 50%; width: 40px; height: 40px;
  animation: spin 1s linear infinite; margin: 0 auto 20px;
}
@keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }
</style>