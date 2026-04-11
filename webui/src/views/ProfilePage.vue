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
          <button @click="$router.push('/upload')" class="icon-btn" title="Upload"><i class="fa-solid fa-circle-plus"></i></button>
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
            <div class="stat-item clickable" @click="showFollowers">
              <span class="stat-value">{{ profileData.FollowersCount }}</span>
              <span class="stat-label">Followers</span>
            </div>
            <div class="stat-item clickable" @click="showFollowing">
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
            <button class="btn-ban" @click="handleBan" title="Ban User">
              <i class="fa-solid fa-user-slash"></i>
            </button>
          </div>
        </section>

        <section class="posts-feed">
          <h3 class="section-title">Latest Posts</h3>
          <div v-if="profileData.UserPosts && profileData.UserPosts.length > 0" class="stream-list">
            <div v-for="post in profileData.UserPosts" :key="post.PhotoId" class="post-card glass-card">
              <img :src="'/images/' + post.PhotoId + '.jpg'" class="post-image" alt="Post">
              
              <div class="post-info">
                <p v-if="post.Caption" class="post-caption">{{ post.Caption }}</p>
                <div class="post-meta">
                  <span :class="{'liked': post.IsLikedByMe}">
                    <i class="fa-solid fa-heart"></i> {{ post.LikesNumber }}
                  </span>
                  <span><i class="fa-solid fa-comment"></i> {{ post.CommentsNumber }}</span>
                  <button v-if="isOwner" class="btn-delete" @click="deletePost(post.PhotoId)">
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
      errorMsg: ""
    };
  },
  computed: {
    isOwner() {
      // Comparison between route param and stored identity
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
        const username = this.$route.params.Username;
        const token = localStorage.getItem("SessionToken");
        
        // API call to our Go Backend
        const response = await this.$axios.get(`/users/${username}`, {
          headers: { Authorization: `Bearer ${token}` }
        });
        
        this.profileData = response.data;
      } catch (e) {
        if (e.response && (e.response.status === 404 || e.response.status === 403)) {
          this.errorMsg = "User not found or profile unavailable.";
        } else {
          this.errorMsg = "An error occurred while loading the profile.";
        }
        console.error("Fetch profile error:", e);
      } finally {
        this.loading = false;
      }
    },
    async toggleFollow() {
      const token = localStorage.getItem("SessionToken");
      const myUsername = localStorage.getItem("Username");
      const target = this.profileData.Username;

      try {
        if (this.profileData.IsFollowing) {
          await this.$axios.delete(`/users/${myUsername}/following/${target}`, {
            headers: { Authorization: `Bearer ${token}` }
          });
          this.profileData.FollowersCount--;
        } else {
          await this.$axios.put(`/users/${myUsername}/following/${target}`, {}, {
            headers: { Authorization: `Bearer ${token}` }
          });
          this.profileData.FollowersCount++;
        }
        this.profileData.IsFollowing = !this.profileData.IsFollowing;
      } catch (e) {
        console.error("Toggle follow error:", e);
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
      // Placeholder for delete logic
      if (confirm("Delete this post?")) {
        console.log("Deleting post:", postId);
      }
    }
  },
  mounted() {
    this.fetchProfile();
  },
  watch: {
    // Detects when searching for a new user while already on a profile
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

/* NAVBAR */
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

/* CARDS & CONTAINERS */
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

/* BUTTONS */
.gradient-button {
  width: 100%; max-width: 180px; height: 42px;
  background: linear-gradient(135deg, #a53b59 0%, #726fb4 100%);
  border: none; border-radius: 12px; color: white; font-weight: 600;
  cursor: pointer; transition: 0.3s;
}
.btn-unfollow { background: transparent; border: 2px solid #a53b59; }

.btn-ban { 
  background: rgba(255, 71, 87, 0.1); color: #ff4757; border: 1px solid rgba(255, 71, 87, 0.2);
  padding: 10px 15px; border-radius: 12px; cursor: pointer; margin-left: 10px;
}

/* POSTS FEED */
.posts-feed { max-width: 600px; margin: 50px auto; }
.section-title { margin-bottom: 25px; font-weight: 300; opacity: 0.7; }
.post-card { margin-bottom: 40px; padding: 15px; overflow: hidden; }
.post-image { width: 100%; border-radius: 15px; display: block; }

.post-info { padding: 15px 5px 5px; }
.post-caption { text-align: left; font-size: 0.95rem; margin-bottom: 12px; color: #eee; }

.post-meta { display: flex; gap: 20px; align-items: center; }
.post-meta i { color: #a53b59; }
.liked i { color: #ff4757; }

.btn-delete { background: transparent; border: none; color: #ff4757; cursor: pointer; margin-left: auto; opacity: 0.6; }
.btn-delete:hover { opacity: 1; }

/* STATES */
.state-container { text-align: center; padding: 100px 20px; }
.error-box i { font-size: 3rem; color: #a53b59; margin-bottom: 20px; }
.error-box p { margin-bottom: 25px; }

.loader {
  border: 4px solid rgba(255, 255, 255, 0.1);
  border-top: 4px solid #a53b59;
  border-radius: 50%; width: 40px; height: 40px;
  animation: spin 1s linear infinite; margin: 0 auto 20px;
}
@keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }
</style>